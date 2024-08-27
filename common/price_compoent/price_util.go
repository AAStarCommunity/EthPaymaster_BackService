package price_compoent

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/config"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"io"
	"math/big"
	"net/http"
	"time"
)

var (
	URLMap = map[global_const.TokenType]string{}
)

type Price struct {
}

func init() {
	URLMap = make(map[global_const.TokenType]string)
	template := "https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd"
	for token, tokenSymbol := range TokenSymbolMap {
		URLMap[token] = fmt.Sprintf(template, tokenSymbol)
	}

}

var globalPriceMap = make(map[global_const.TokenType]float64)

var TokenSymbolMap = map[global_const.TokenType]string{
	global_const.TokenTypeETH: "ethereum",
	global_const.TokenTypeOP:  "optimism",
}

func Init() {

	go func() {
		for {
			GetConfigTokenPrice()
			time.Sleep(60 * time.Second)
		}
	}()
}

func GetGlobalPriceMap() map[global_const.TokenType]float64 {
	return globalPriceMap
}

func GetConfigTokenPrice() {
	priceNetworkMap := make(map[global_const.TokenType]float64)
	for token, tokenUrl := range URLMap {

		req, err := http.NewRequest("GET", tokenUrl, nil)
		if err != nil {
			logrus.Error(xerrors.Errorf("[Price Thread ERROR] http request error: %w", err))
			continue
		}
		req.Header.Add("x-cg-demo-api-key", config.GetPriceOracleApiKey())

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			logrus.Error(xerrors.Errorf("[Price Thread ERROR] http request error: %w", err))
			continue
		}
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)
		priceMap := make(map[string]map[string]float64)
		err = json.Unmarshal(body, &priceMap)
		if err != nil {
			logrus.Error(xerrors.Errorf("[Price Thread ERROR] json unmarshal error: %w", err))
			continue
		}
		tokenSymPol := TokenSymbolMap[token]
		price := priceMap[tokenSymPol]["usd"]
		priceNetworkMap[token] = price

	}
	globalPriceMap = priceNetworkMap
	logrus.Infof("[Price Thread] Update price map: %v", priceNetworkMap)
}
func GetTokenCostInUsd(tokenType global_const.TokenType, amount *big.Float) (*big.Float, error) {
	price, err := GetPriceUsd(tokenType)
	if err != nil {
		return nil, xerrors.Errorf("get price error: %w", err)
	}
	amountInUsd := new(big.Float).Mul(new(big.Float).SetFloat64(price), amount)
	return amountInUsd, nil
}

func GetPriceUsd(tokenType global_const.TokenType) (float64, error) {
	if global_const.IsStableToken(tokenType) {
		return 1, nil
	}
	price, ok := globalPriceMap[tokenType]
	if !ok {
		return 0, xerrors.Errorf("tokens type [%w] not found", tokenType)
	}
	return price, nil
}

// GetToken Get The FromToken/ToToken Ratew
func GetToken(fromToken global_const.TokenType, toToken global_const.TokenType) (float64, error) {
	if toToken == global_const.TokenTypeUSDT {
		return GetPriceUsd(fromToken)
	}
	formTokenPrice, _ := GetPriceUsd(fromToken)
	toTokenPrice, _ := GetPriceUsd(toToken)

	return formTokenPrice / toTokenPrice, nil
}

//func GetCoinMarketPrice() {
//	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v2/tools/price-conversion", nil)
//	if err != nil {
//		log.Print(err)
//		os.Exit(1)
//	}
//	q := url.Values{}
//	q.Add("amount", "2")
//	q.Add("symbol", "BTC")
//	q.Add("convert", "USD")
//
//	req.Header.Set("Accepts", "application/json")
//	req.Header.Add("X-CMC_PRO_API_KEY", "a1441679-b8fd-49a0-aa47-51f88f7d3d52")
//	req.URL.RawQuery = q.Encode()
//	resp, err := httpClient.Do(req)
//	if err != nil {
//		fmt.Println("Error sending request to server")
//		os.Exit(1)
//	}
//	fmt.Println(resp.Status)
//	respBody, _ := ioutil.ReadAll(resp.Body)
//	fmt.Println(string(respBody))
//}
