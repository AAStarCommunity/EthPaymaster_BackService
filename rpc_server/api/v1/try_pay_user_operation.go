package v1

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/service/executor"
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
// @Param tryPay body model.TryPayUserOpRequest true "UserOp Request"
// @Success 200
// @Security JWT
func TryPayUserOperation(c *gin.Context) {
	//1.TODO API validate
	//2. recall service
	request := model.TryPayUserOpRequest{}
	response := model.GetResponse()

	if err := c.ShouldBindJSON(&request); err != nil {
		errStr := fmt.Sprintf("Request Error [%v]", err)
		response.SetHttpCode(http.StatusBadRequest).FailCode(c, http.StatusBadRequest, errStr)
	}

	if err := request.Validate(); err != nil {
		errStr := fmt.Sprintf("Request Error [%v]", err)
		response.SetHttpCode(http.StatusBadRequest).FailCode(c, http.StatusBadRequest, errStr)
	}

	result, err := executor.TryPayUserOpExecute(request)
	if err != nil {
		errStr := fmt.Sprintf("TryPayUserOpExecute ERROR [%v]", err)
		response.SetHttpCode(http.StatusInternalServerError).FailCode(c, http.StatusInternalServerError, errStr)
	}
	response.WithData(result).Success(c)
}
