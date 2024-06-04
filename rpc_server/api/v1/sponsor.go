package v1

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/sponsor_manager"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// DepositSponsor
// @Tags Sponsor
// @Description Deposit Sponsor
// @Accept json
// @Product json
// @Param request body DepositSponsorRequest true "DepositSponsorRequest Model
// @Param is_test_net path boolean true "Is Test Net"
// @Router /api/v1/paymaster_sponsor/deposit [post]
// @Success 200
func DepositSponsor(ctx *gin.Context) {
	request := model.DepositSponsorRequest{}
	response := model.GetResponse()
	if err := ctx.ShouldBindJSON(&request); err != nil {
		errStr := fmt.Sprintf("Request Error [%v]", err)
		response.SetHttpCode(http.StatusBadRequest).FailCode(ctx, http.StatusBadRequest, errStr)
		return
	}
	//TODO Add Signature Verification
	result, err := sponsor_manager.DepositSponsor(&request)
	if err != nil {
		response.SetHttpCode(http.StatusInternalServerError).FailCode(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.WithDataSuccess(ctx, result)
	return
}

// WithdrawSponsor
// @Tags Sponsor
// @Description Withdraw Sponsor
// @Accept json
// @Product json
// @Param request body WithdrawSponsorRequest true "WithdrawSponsorRequest Model"
// @Param is_test_net path boolean true "Is Test Net"
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
	//TODO Add Signature Verification
	result, err := sponsor_manager.WithDrawSponsor(&request)
	if err != nil {
		response.SetHttpCode(http.StatusInternalServerError).FailCode(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.WithDataSuccess(ctx, result)
	return
}

type sponsorDepositTransaction struct {
	TxHash     string                  `json:"tx_hash"`
	Amount     string                  `json:"amount"`
	UpdateType global_const.UpdateType `json:"update_type"`
}

// GetSponsorDepositAndWithdrawTransactions
// @Tags Sponsor
// @Description Get Sponsor Deposit And Withdraw Transactions
// @Accept json
// @Product json
// @Param userId path string true "User Id"
// @Param is_test_net path boolean true "Is Test Net"
// @Router /api/v1/paymaster_sponsor/deposit_log
// @Success 200
func GetSponsorDepositAndWithdrawTransactions(ctx *gin.Context) {
	userId := ctx.Param("user_id")
	textNet := ctx.Param("is_test_net")
	// convertTOBool
	isTestNet, _ := strconv.ParseBool(textNet)
	response := model.GetResponse()
	models, err := sponsor_manager.GetDepositAndWithDrawLog(userId, isTestNet)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.FailCode(ctx, 400, "No Deposit Transactions")
		}
	}
	trans := make([]sponsorDepositTransaction, 0)
	for _, depositModel := range models {
		tran := sponsorDepositTransaction{
			TxHash: depositModel.TxHash,
			Amount: depositModel.Amount.String(),
		}
		trans = append(trans, tran)
	}
	response.WithDataSuccess(ctx, trans)
	return
}

// GetSponsorMetaData
// @Tags Sponsor
// @Description Get Sponsor Balance
// @Accept json
// @Product json
// @Param userId path string true "User Id"
// @Router /api/v1/paymaster_sponsor/balance/{userId}
// @Success 200
func GetSponsorMetaData(ctx *gin.Context) {
	userId := ctx.Param("userId")
	textNet := ctx.Param("is_test_net")
	isTestNet, _ := strconv.ParseBool(textNet)
	response := model.GetResponse()
	balance, err := sponsor_manager.FindUserSponsorBalance(userId, isTestNet)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.FailCode(ctx, 400, "No Balance")
		}
	}
	result := struct {
		AvailableBalance string `json:"available_balance"`
		SponsorAddress   string `json:"sponsor_address"`
	}{
		AvailableBalance: balance.AvailableBalance.String(),
		SponsorAddress:   balance.SponsorAddress,
	}
	response.WithDataSuccess(ctx, result)
	return
}
