package v1

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/service/operator"
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
	request := model.GetSupportEntrypointRequest{}
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
	result, err := operator.GetSupportEntrypointExecute(&request)
	if err != nil {
		errStr := fmt.Sprintf("%v", err)
		response.SetHttpCode(http.StatusInternalServerError).FailCode(c, http.StatusInternalServerError, errStr)
		return
	}
	response.WithData(result).Success(c)

}
