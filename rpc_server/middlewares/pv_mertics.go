package middlewares

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/service/dashboard_service"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
	"time"
)

type PayMasterParam struct {
	ApiUserId    int64
	ApiKey       string
	Method       string
	SendTime     string
	Latency      time.Duration
	RequestBody  string
	ResponseBody string
	NetWork      string
	Status       int
}

func PvMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/v1/paymaster/") {
			responseWriter := &CustomResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
			c.Writer = responseWriter
			start := time.Now()
			// get Request Body
			requestBodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
				c.Abort()
				return
			}
			//Restore the request body to Request.Body for use by subsequent handlers.
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))

			c.Next()

			request := model.JsonRpcRequest{}
			_ = json.Unmarshal(requestBodyBytes, &request)
			requestJson, _ := json.Marshal(request)
			requestBodyStr := string(requestJson)

			responseWriter.body.Bytes()
			responseBody := responseWriter.body.String()
			endEnd := time.Now()
			network := c.Param("network")
			metricsParam := PayMasterParam{
				NetWork:      network,
				Method:       request.Method,
				SendTime:     start.Format("2006-01-02 15:04:05.MST"),
				Latency:      endEnd.Sub(start),
				Status:       c.Writer.Status(),
				RequestBody:  requestBodyStr,
				ResponseBody: responseBody,
			}
			apiKeyContext, exist := c.Get(global_const.ContextKeyApiMoDel)
			if exist {
				apiKeyModel := apiKeyContext.(*model.ApiKeyModel)
				metricsParam.ApiKey = apiKeyModel.ApiKey
				metricsParam.ApiUserId = apiKeyModel.UserId
			}
			MetricsPaymaster(c, metricsParam)
		} else {
			return
		}

	}
}

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func MetricsPaymaster(c *gin.Context, metricsParam PayMasterParam) {

	recallModel := dashboard_service.PaymasterRecallLogDbModel{
		ProjectApikey:   metricsParam.ApiKey,
		ProjectUserId:   metricsParam.ApiUserId,
		PaymasterMethod: metricsParam.Method,
		SendTime:        metricsParam.SendTime,
		Latency:         int64(metricsParam.Latency),
		RequestBody:     metricsParam.RequestBody,
		ResponseBody:    metricsParam.ResponseBody,
		Status:          metricsParam.Status,
		NetWork:         metricsParam.NetWork,
	}
	err := dashboard_service.CreatePaymasterCall(&recallModel)
	if err != nil {
		logrus.Error("CreatePaymasterCall error:", err)
	}
}
