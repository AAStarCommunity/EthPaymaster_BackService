package v1

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"AAStarCommunity/EthPaymaster_BackService/envirment"
	"AAStarCommunity/EthPaymaster_BackService/service/operator"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
	"net/http"
)

// TryPayUserOperation
// @Tags Sponsor
// @Description sponsor the userOp
// @Accept json
// @Product json
// @Router /api/v1/try-pay-user-operation [POST]
// @Param tryPay body model.UserOpRequest true "UserOp Request"
// @Success 200
// @Security JWT
func TryPayUserOperation(c *gin.Context) {
	request := model.UserOpRequest{}
	response := model.GetResponse()

	//1. API validate
	if err := c.ShouldBindJSON(&request); err != nil {
		errStr := fmt.Sprintf("Request Error [%v]", err)
		response.SetHttpCode(http.StatusBadRequest).FailCode(c, http.StatusBadRequest, errStr)
		return
	}

	if err := ValidateUserOpRequest(request); err != nil {
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
func ValidateUserOpRequest(request model.UserOpRequest) error {
	if len(request.ForceStrategyId) == 0 {
		if len(request.ForceNetwork) == 0 || len(request.ForceToken) == 0 || len(request.ForceEntryPointAddress) == 0 {
			return xerrors.Errorf("strategy configuration illegal")
		}
	}
	if request.ForceStrategyId == "" && (request.ForceToken == "" || request.ForceNetwork == "") {
		return xerrors.Errorf("Token And Network Must Set When ForceStrategyId Is Empty")
	}
	if envirment.Environment.IsDevelopment() && request.ForceNetwork != "" {
		if !conf.IsTestNet(request.ForceNetwork) {
			return xerrors.Errorf("ForceNetwork: [%s] is not test network", request.ForceNetwork)
		}
	}
	exist := conf.CheckEntryPointExist(request.ForceNetwork, request.ForceEntryPointAddress)
	if !exist {
		return xerrors.Errorf("ForceEntryPointAddress: [%s] not exist in [%s] network", request.ForceEntryPointAddress, request.ForceNetwork)
	}
	return nil
}
