package httpx

import (
	"fortunefinger/backend/internal/pkg/errorx"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Body struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type DataBody struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Ok(c *gin.Context, msg string, data ...any) {
	if len(data) > 0 {
		c.JSON(http.StatusOK, DataBody{
			Code: 0,
			Msg:  msg,
			Data: data[0],
		})
		return
	}
	c.JSON(http.StatusOK, Body{Code: 0, Msg: msg})
}

func Error(c *gin.Context, err error) {
	appErr, ok := errorx.As(err)
	if !ok {
		appErr = errorx.New(errorx.InternalServer, "")
	}
	writeJSON(c, errorx.HTTPStatus(appErr.Code()), int(appErr.Code()), appErr.Message())
}

func writeJSON(c *gin.Context, httpStatus, code int, msg string) {
	c.JSON(httpStatus, Body{Code: code, Msg: msg})
}
