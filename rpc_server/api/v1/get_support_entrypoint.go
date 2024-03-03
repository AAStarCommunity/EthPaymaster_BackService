package v1

import (
	"AAStarCommunity/EthPaymaster_BackService/rpc_server/models"
	"AAStarCommunity/EthPaymaster_BackService/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetSupportEntrypoint
// @Tags Sponsor
// @Description get the support entrypoint
// @Accept json
// @Product json
// @Router /api/v1/get-support-entrypoint [get]
// @Success 200
func GetSupportEntrypoint(c *gin.Context) {
	//1.TODO API validate
	//2. recall service
	result, err := service.GetSupportEntrypointExecute()
	response := models.GetResponse()
	if err != nil {
		errStr := fmt.Sprintf("%v", err)
		response.SetHttpCode(http.StatusInternalServerError).FailCode(c, http.StatusInternalServerError, errStr)
	}
	response.WithData(result).Success(c)
}
