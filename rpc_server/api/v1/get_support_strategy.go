package v1

import (
	"AAStarCommunity/EthPaymaster_BackService/rpc_server/models"
	"AAStarCommunity/EthPaymaster_BackService/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetSupportStrategy
// @Tags Sponsor
// @Description get the support strategy
// @Accept json
// @Produce json
// @Success 200
// @Router /api/v1/get-support-strategy [GET]
// @Security JWT
func GetSupportStrategy(c *gin.Context) {
	//1.TODO API validate
	//2. recall service
	result, err := service.GetSupportStrategyExecute()
	response := models.GetResponse()
	if err != nil {
		errStr := fmt.Sprintf("%v", err)
		response.SetHttpCode(http.StatusInternalServerError).FailCode(c, http.StatusInternalServerError, errStr)
	}
	response.WithData(result).Success(c)
}
