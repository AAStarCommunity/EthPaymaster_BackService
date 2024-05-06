package v1

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/service/operator"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
	"net/http"
)

// GetSupportStrategy
// @Tags Sponsor
// @Description get the support strategy
// @Accept json
// @Produce json
// @Success 200
// @Param network query string false "network"
// @Router /api/v1/get-support-strategy [GET]
// @Security JWT
func GetSupportStrategy(c *gin.Context) {
	response := model.GetResponse()
	network := c.Query("network")
	if network == "" {
		errStr := fmt.Sprintf("Request Error [network is empty]")
		response.SetHttpCode(http.StatusBadRequest).FailCode(c, http.StatusBadRequest, errStr)
		return
	}

	if result, err := operator.GetSupportStrategyExecute(network); err != nil {
		errStr := fmt.Sprintf("%v", err)
		response.SetHttpCode(http.StatusInternalServerError).FailCode(c, http.StatusInternalServerError, errStr)
		return
	} else {
		response.WithData(result).Success(c)
		return
	}
}
func GetSupportStrategyFunc() MethodFunctionFunc {
	return func(ctx *gin.Context, jsonRpcRequest model.JsonRpcRequest) (result interface{}, err error) {
		if jsonRpcRequest.Params[0] == nil || jsonRpcRequest.Params[0].(string) == "" {
			return nil, xerrors.Errorf("Request Error [network is empty]")
		}
		return operator.GetSupportStrategyExecute(jsonRpcRequest.Params[0].(string))
	}
}
