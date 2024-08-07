package v1

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/price_compoent"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/config"
	"AAStarCommunity/EthPaymaster_BackService/sponsor_manager"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
	"log"
	"math/big"
	"net/http"
)

// DepositSponsor
// @Tags DepositSponsor
// @Description Deposit Sponsor
// @Accept json
// @Product json
// @Param request body model.DepositSponsorRequest true "DepositSponsorRequest Model"
// @Param relay_hash header string false "relay Request  Body Hash"
// @Param relay_signature header string false "relay Request  Body Hash"
// @Router /api/v1/paymaster_sponsor/deposit [post]
// @Success 200
func DepositSponsor(ctx *gin.Context) {
	request := model.DepositSponsorRequest{}
	response := model.GetResponse()
	if err := ctx.ShouldBindJSON(&request); err != nil {
		errStr := fmt.Sprintf("Request Error [%v]", err)
		logrus.Errorf("Request Error [%v]", err)
		response.SetHttpCode(http.StatusBadRequest).FailCode(ctx, http.StatusBadRequest, errStr)
		return
	}
	if request.DepositSource != "dashboard" {
		errStr := fmt.Sprintf("not Support Source")
		logrus.Errorf("Deposit Source Error :Not Support Source")
		response.SetHttpCode(http.StatusBadRequest).FailCode(ctx, http.StatusBadRequest, errStr)
		return
	}
	//validate Signature
	inputJson, err := json.Marshal(request)
	if err != nil {
		logrus.Errorf("Marshal Error [%v]", err)
		response.SetHttpCode(http.StatusInternalServerError).FailCode(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	var signerAddress string
	//TODO change to Enum
	if request.DepositSource == "dashboard" {
		signerAddress = config.GetSponsorConfig().DashBoardSignerAddress
	} else {
		logrus.Errorf("Deposit Source Error :Not Support Source")
		response.SetHttpCode(http.StatusBadRequest).FailCode(ctx, http.StatusBadRequest, "Deposit Source Error :Not Support Source")
		return
	}
	err = ValidateSignature(ctx.GetHeader("relay_hash"), ctx.GetHeader("relay_signature"), inputJson, signerAddress)
	if err != nil {
		logrus.Errorf("Validate Signature Error [%v]", err)
		response.SetHttpCode(http.StatusBadRequest).FailCode(ctx, http.StatusBadRequest, err.Error())
		return
	}
	//validate Deposit
	sender, amount, err := validateDeposit(&request)
	if err != nil {
		logrus.Errorf("Validate Deposit Error [%v]", err)
		response.SetHttpCode(http.StatusBadRequest).FailCode(ctx, http.StatusBadRequest, err.Error())
		return
	}

	depositInput := sponsor_manager.DepositSponsorInput{
		From:      sender.Hex(),
		Amount:    amount,
		TxHash:    request.TxHash,
		PayUserId: request.PayUserId,
		Source:    request.DepositSource,
		IsTestNet: request.IsTestNet,
	}
	result, err := sponsor_manager.DepositSponsor(&depositInput)
	if err != nil {
		logrus.Errorf("Deposit Sponsor Error [%v]", err)
		response.SetHttpCode(http.StatusInternalServerError).FailCode(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.WithDataSuccess(ctx, result)
	return
}
func ValidateSignature(originHash string, signatureHex string, inputJson []byte, signerAddress string) error {
	hash := sha256.New()
	hash.Write(inputJson)
	hashBytes := hash.Sum(nil)
	hashHex := hex.EncodeToString(hashBytes)
	if hashHex != originHash {
		return xerrors.Errorf("Hash Not Match")
	}
	signerAddressHex := common.HexToAddress(signerAddress)

	hashByte, _ := utils.DecodeStringWithPrefix(originHash)
	signatureByte, _ := utils.DecodeStringWithPrefix(signatureHex)
	pubKey, err := crypto.SigToPub(accounts.TextHash(hashByte), signatureByte)
	if err != nil {
		log.Fatalf("Failed to recover public key: %v", err)
	}
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	if signerAddressHex != recoveredAddr {
		return xerrors.Errorf("Signer Address Not Match")
	}
	return nil
}
func validateDeposit(request *model.DepositSponsorRequest) (sender *common.Address, amount *big.Float, err error) {
	txHash := request.TxHash
	client := config.GetPaymasterSponsorClient(request.IsTestNet)

	// check tx
	_, err = sponsor_manager.GetLogByTxHash(txHash, request.IsTestNet)
	if err == nil {
		return nil, nil, xerrors.Errorf("Transaction [%s] already exist", txHash)
	}
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, err
		}
	}
	tx, err := GetInfoByHash(txHash, client)
	if err != nil {
		return nil, nil, err
	}
	if tx.Type() != types.DynamicFeeTxType {
		return nil, nil, xerrors.Errorf("Tx Type is not DynamicFeeTxType")
	}
	txSender, err := types.Sender(types.NewLondonSigner(tx.ChainId()), tx)
	if err != nil {
		logrus.Errorf("Get Sender Error [%v]", err)
		return nil, nil, err
	}
	sender = &txSender
	if request.IsTestNet {
		//Only ETH
		if tx.Value().Uint64() == 0 {
			return nil, nil, xerrors.Errorf("Tx Value is 0")
		}
		if tx.To() == nil {
			return nil, nil, xerrors.Errorf("Tx To Address is nil")
		}
		if tx.To().Hex() != config.GetDepositerAddress().String() {
			return nil, nil, xerrors.Errorf("Tx To Address is not Sponsor Address")
		}
		value := tx.Value()
		valueEth := utils.ConvertBalanceToEther(value)
		logrus.Infof("ETH amount : %s", valueEth)

		amount, err = price_compoent.GetTokenCostInUsd(global_const.TokenTypeETH, valueEth)
		if err != nil {
			return nil, nil, err
		}
	} else {
		return nil, nil, xerrors.Errorf("not Support MainNet Right Now")
		//contractAddress := tx.To()
		//chain_service.CheckContractAddressAccess(contractAddress,"")
		//Only Usdt

	}
	return sender, amount, nil
}

func GetInfoByHash(txHash string, client *ethclient.Client) (*types.Transaction, error) {
	txHashHex := common.HexToHash(txHash)
	//TODO consider about pending
	tx, _, err := client.TransactionByHash(context.Background(), txHashHex)
	if err != nil {
		if errors.Is(err, ethereum.NotFound) {
			return nil, xerrors.Errorf("Transaction [%s] not found", txHash)
		}
		return nil, err
	}
	return tx, nil
}

// GetTokenPrice
// @Tags Sponsor
// @Description Get Token Price
// @Accept json
// @Product json
// @Router /api/v1/paymaster_sponsor/token_price [get]
// @Success 200
func GetTokenPrice(ctx *gin.Context) {
	response := model.GetResponse()
	result := price_compoent.GetGlobalPriceMap()
	response.WithDataSuccess(ctx, result)
}

// WithdrawSponsor
// @Tags Sponsor
// @Description Withdraw Sponsor
// @Accept json
// @Product json
// @Param request body model.WithdrawSponsorRequest true "WithdrawSponsorRequest Model"
// @Router /api/v1/paymaster_sponsor/withdraw [post]
// @Success 200
func WithdrawSponsor(ctx *gin.Context) {

	request := model.WithdrawSponsorRequest{}
	response := model.GetResponse()
	if err := ctx.ShouldBindJSON(&request); err != nil {
		errStr := fmt.Sprintf("Request Error [%v]", err)
		response.SetHttpCode(http.StatusBadRequest).FailCode(ctx, http.StatusBadRequest, errStr)
		return
	}
	client := config.GetPaymasterSponsorClient(request.IsTestNet)
	//validate Sign
	//validate Signature
	inputJson, err := json.Marshal(request)
	if err != nil {
		response.SetHttpCode(http.StatusInternalServerError).FailCode(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	var signerAddress string
	if request.DepositSource == "dashboard" {
		signerAddress = config.GetSponsorConfig().DashBoardSignerAddress
	} else {
		response.SetHttpCode(http.StatusBadRequest).FailCode(ctx, http.StatusBadRequest, "Deposit Source Error :Not Support Source")
		return
	}
	err = ValidateSignature(ctx.GetHeader("relay_hash"), ctx.GetHeader("relay_signature"), inputJson, signerAddress)
	if err != nil {
		response.SetHttpCode(http.StatusBadRequest).FailCode(ctx, http.StatusBadRequest, err.Error())
		return
	}
	//checkBalance
	balanceModel, err := sponsor_manager.FindUserSponsorBalance(request.PayUserId, request.IsTestNet)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.SetHttpCode(http.StatusBadRequest).FailCode(ctx, http.StatusBadRequest, "Balance Not Found")
			return
		}
		response.SetHttpCode(http.StatusInternalServerError).FailCode(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	amountFloat := big.NewFloat(request.Amount)

	if balanceModel.AvailableBalance.Cmp(amountFloat) < 0 {
		response.SetHttpCode(http.StatusBadRequest).FailCode(ctx, http.StatusBadRequest, "Insufficient Balance")
		return
	}
	// U -> ETH
	ethPrice, err := price_compoent.GetTokenCostInUsd(global_const.TokenTypeETH, amountFloat)
	if err != nil {
		response.SetHttpCode(http.StatusInternalServerError).FailCode(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	// U -> ETH -> Wei
	ethCount := new(big.Float)
	ethCount.Quo(ethPrice, amountFloat)

	ethWeiValue := utils.CounverEtherToWei(ethCount)

	toAddress := common.HexToAddress(request.RefundAddress)
	chainId := config.GetPaymasterSponsorChainId(request.IsTestNet)
	// Execute transfer
	tx, err := utils.TransEth(config.GetWithdrawerEoa().PrivateKey, &toAddress, client, ethWeiValue, chainId)
	// WithDrawSponsor
	result, err := sponsor_manager.WithDrawSponsor(&request, tx.Hash().Hex())
	if err != nil {
		response.SetHttpCode(http.StatusInternalServerError).FailCode(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.WithDataSuccess(ctx, result)
	return
}
