package utils

import (
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"fmt"
	"golang.org/x/xerrors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var (
	URLMap     = map[types.TokenType]string{}
	httpClient = &http.Client{}
)

type Price struct {
}

func init() {
	URLMap = make(map[types.TokenType]string)
	URLMap[types.ETH] = "https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd"
	URLMap[types.OP] = "https://api.coingecko.com/api/v3/simple/price?ids=optimism&vs_currencies=usd"
}

func GetPriceUsd(tokenType types.TokenType) (float64, error) {

	if types.IsStableToken(tokenType) {
		return 1, nil
	}
	if tokenType == types.ETH {
		return 4000, nil
	}
	url, ok := URLMap[tokenType]
	if !ok {
		return 0, xerrors.Errorf("token type [%w] not found", tokenType)
	}
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-cg-demo-api-key", "CG-ioE6p8cmmSFBFwJnKECCbZ7U\t")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	bodystr := string(body)
	strarr := strings.Split(bodystr, ":")
	usdstr := strings.TrimRight(strarr[2], "}}")
	return strconv.ParseFloat(usdstr, 64)
}
func GetToken(fromToken types.TokenType, toToken types.TokenType) (float64, error) {
	if toToken == types.USDT {
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
