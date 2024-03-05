package v1

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/service/executor"
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
	result, err := executor.GetSupportStrategyExecute()
	response := model.GetResponse()
	if err != nil {
		errStr := fmt.Sprintf("%v", err)
		response.SetHttpCode(http.StatusInternalServerError).FailCode(c, http.StatusInternalServerError, errStr)
	}
	response.WithData(result).Success(c)
}
