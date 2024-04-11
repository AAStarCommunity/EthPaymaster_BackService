package v1

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/service/operator"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// EstimateUserOpGas
// @Tags Sponsor
// @Description get the estimate gas of the userOp
// @Accept json
// @Product json
// @Router /api/v1/estimate-user-operation-gas [POST]
// @Param tryPay body model.UserOpRequest true "UserOp Request"
// @Success 200
// @Security JWT
func EstimateUserOpGas(c *gin.Context) {
	request := model.UserOpRequest{}
	response := model.GetResponse()
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
	if result, err := operator.GetEstimateUserOpGas(&request); err != nil {
		errStr := fmt.Sprintf("GetEstimateUserOpGas ERROR [%v]", err)
		response.SetHttpCode(http.StatusInternalServerError).FailCode(c, http.StatusInternalServerError, errStr)
		return
	} else {
		response.WithDataSuccess(c, result)
		return
	}
}
