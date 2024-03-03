package v1

import (
	"AAStarCommunity/EthPaymaster_BackService/rpc_server/models"
	"AAStarCommunity/EthPaymaster_BackService/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// TryPayUserOperation
// @Tags Sponsor
// @Description sponsor the userOp
// @Accept json
// @Product json
// @Router /api/v1/try-pay-user-operation [POST]
// @Success 200
// @Security JWT
func TryPayUserOperation(c *gin.Context) {
	//1.TODO API validate
	//2. recall service
	result, err := service.TryPayUserOpExecute()
	response := models.GetResponse()
	if err != nil {
		errStr := fmt.Sprintf("%v", err)
		response.SetHttpCode(http.StatusInternalServerError).FailCode(c, http.StatusInternalServerError, errStr)
	}
	response.WithData(result).Success(c)
}
