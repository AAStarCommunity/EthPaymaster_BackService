package price_compoent

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/config"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var (
	URLMap     = map[global_const.TokenType]string{}
	httpClient = &http.Client{}
)

type Price struct {
}

func init() {
	URLMap = make(map[global_const.TokenType]string)
	URLMap[global_const.TokenTypeETH] = "https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd"
	URLMap[global_const.TokenTypeOP] = "https://api.coingecko.com/api/v3/simple/price?ids=optimism&vs_currencies=usd"
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
	//if tokenType == global_const.ETH {
	//	return 3100, nil
	//}
	tokenUrl, ok := URLMap[tokenType]
	if !ok {
		return 0, xerrors.Errorf("tokens type [%w] not found", tokenType)
	}
	req, _ := http.NewRequest("GET", tokenUrl, nil)
	//TODO remove APIKey
	req.Header.Add("x-cg-demo-api-key", config.GetPriceOracleApiKey())

	res, _ := http.DefaultClient.Do(req)

	defer func(Body io.ReadCloser) {
		if Body == nil {
			return
		}
		err := Body.Close()
		if err != nil {
			logrus.Error("close body error: ", err)
			return
		}
	}(res.Body)
	body, _ := io.ReadAll(res.Body)
	bodystr := string(body)
	strarr := strings.Split(bodystr, ":")
	usdstr := strings.TrimRight(strarr[2], "}}")
	return strconv.ParseFloat(usdstr, 64)
}

// GetToken Get The FromToken/ToToken Rate
func GetToken(fromToken global_const.TokenType, toToken global_const.TokenType) (float64, error) {
	if toToken == global_const.TokenTypeUSDT {
		return GetPriceUsd(fromToken)
	}
	formTokenPrice, _ := GetPriceUsd(fromToken)
	toTokenPrice, _ := GetPriceUsd(toToken)

	return formTokenPrice / toTokenPrice, nil
}

func GetCoinMarketPrice() {
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v2/tools/price-conversion", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	q := url.Values{}
	q.Add("amount", "2")
	q.Add("symbol", "BTC")
	q.Add("convert", "USD")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", "a1441679-b8fd-49a0-aa47-51f88f7d3d52")
	req.URL.RawQuery = q.Encode()
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
		os.Exit(1)
	}
	fmt.Println(resp.Status)
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respBody))
}
