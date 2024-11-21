package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func jsonOutPut(ctx *gin.Context, code int, msg string, data ...interface{}) {

	if len(data) > 0 && data[0] != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  msg,
			"data": data[0],
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": make(map[string]interface{}),
	})
}
