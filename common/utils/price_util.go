package utils

import (
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"golang.org/x/xerrors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var (
	URLMap = map[types.TokenType]string{}
)

type Price struct {
}

func init() {
	URLMap = make(map[types.TokenType]string)
	URLMap[types.ETH] = "https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd"
	URLMap[types.OP] = "https://api.coingecko.com/api/v3/simple/price?ids=optimism&vs_currencies=usd"
}

func GetPriceUsd(tokenType types.TokenType) (float64, error) {
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
