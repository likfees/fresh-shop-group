package response

import (
	"github.com/silenceper/wechat/v2/pay/notify"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	ERROR         = 7
	SUCCESS       = 0
	TOKEN_EXPIRED = 401
)

func Result(code int, data interface{}, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "操作成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "查询成功", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "操作失败", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}

func FailWithTookenExpired(data interface{}, message string, c *gin.Context) {
	Result(TOKEN_EXPIRED, data, message, c)
}

// WxpayNotify 响应微信支付回调
func WxpayNotify(returnCode, returnMsg string, c *gin.Context) {
	c.JSON(http.StatusOK, notify.PaidResp{
		ReturnCode: returnCode,
		ReturnMsg:  returnMsg,
	})
}
