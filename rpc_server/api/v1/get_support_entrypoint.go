package v1

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/service/operator"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
	"net/http"
)

// GetSupportEntrypoint
// @Tags Sponsor
// @Description get the support entrypoint
// @Accept json
// @Product json
// @Router /api/v1/get-support-entrypoint [GET]
// @Param network query string false "network"
// @Success 200
// @Security JWT
func GetSupportEntrypoint(c *gin.Context) {
	response := model.GetResponse()
	//1. API validate
	network := c.Query("network")
	if network == "" {
		errStr := fmt.Sprintf("Request Error [network is empty]")
		response.SetHttpCode(http.StatusBadRequest).FailCode(c, http.StatusBadRequest, errStr)
		return
	}

	//2. recall service
	result, err := operator.GetSupportEntrypointExecute(network)
	if err != nil {
		errStr := fmt.Sprintf("%v", err)
		response.SetHttpCode(http.StatusInternalServerError).FailCode(c, http.StatusInternalServerError, errStr)
		return
	}
	response.WithData(result).Success(c)
}

func GetSupportEntryPointFunc() MethodFunctionFunc {
	return func(ctx *gin.Context, jsonRpcRequest model.JsonRpcRequest) (result interface{}, err error) {
		if jsonRpcRequest.Params[0] == nil || jsonRpcRequest.Params[0].(string) == "" {
			return nil, xerrors.Errorf("Request Error [network is empty]")
		}
		return operator.GetSupportEntrypointExecute(jsonRpcRequest.Params[0].(string))
	}
}
