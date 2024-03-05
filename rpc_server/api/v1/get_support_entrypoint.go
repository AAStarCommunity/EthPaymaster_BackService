package v1

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/rpc_server/api/utils"
	"AAStarCommunity/EthPaymaster_BackService/service/executor"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetSupportEntrypoint
// @Tags Sponsor
// @Description get the support entrypoint
// @Accept json
// @Product json
// @Router /api/v1/get-support-entrypoint [GET]
// @Success 200
// @Security JWT
func GetSupportEntrypoint(c *gin.Context) {
	response := model.GetResponse()
	if ok, apiKey := utils.CurrentUser(c); ok {
		_ = apiKey

		//1.TODO API validate
		//2. recall service
		result, err := executor.GetSupportEntrypointExecute()
		if err != nil {
			errStr := fmt.Sprintf("%v", err)
			response.SetHttpCode(http.StatusInternalServerError).FailCode(c, http.StatusInternalServerError, errStr)
		}
		response.WithData(result).Success(c)
	} else {
		response.SetHttpCode(http.StatusUnauthorized)
	}
}
