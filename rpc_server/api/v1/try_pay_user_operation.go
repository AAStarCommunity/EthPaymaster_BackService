package v1

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/service/operator"
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
	request := model.TryPayUserOpRequest{}
	response := model.GetResponse()

	//1. API validate
	if err := c.ShouldBindJSON(&request); err != nil {
		errStr := fmt.Sprintf("Request Error [%v]", err)
		response.SetHttpCode(http.StatusBadRequest).FailCode(c, http.StatusBadRequest, errStr)
		return
	}

	if err := request.Validate(); err != nil {
		errStr := fmt.Sprintf("Request Error [%v]", err)
		response.SetHttpCode(http.StatusBadRequest).FailCode(c, http.StatusBadRequest, errStr)
		return
	}

	//2. recall service
	if result, err := operator.TryPayUserOpExecute(&request); err != nil {
		errStr := fmt.Sprintf("TryPayUserOpExecute ERROR [%v]", err)
		response.SetHttpCode(http.StatusInternalServerError).FailCode(c, http.StatusInternalServerError, errStr)
		return
	} else {
		response.WithDataSuccess(c, result)
		return
	}
}
