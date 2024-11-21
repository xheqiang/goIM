package service

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type assetService struct{}

var AssetService = new(assetService)

func (a *assetService) UpImage(ctx *gin.Context) {

	request := ctx.Request

	// 获取文件
	srcFile, head, err := request.FormFile("file")
	if err != nil {
		jsonOutPut(ctx, -1, "文件上传失败", string(err.Error()))
		return
	}

	// 检查文件后缀
	suffix := ".png"
	fileName := head.Filename
	fileNameArr := strings.Split(fileName, ".")
	if len(fileNameArr) > 1 {
		suffix = "." + fileNameArr[len(fileNameArr)-1]
	}

	upFileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)
	dstFile, err := os.Create("./asset/upload/" + upFileName)
	if err != nil {
		jsonOutPut(ctx, -1, "文件创建失败", string(err.Error()))
		return
	}
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		jsonOutPut(ctx, -1, "文件写入失败", string(err.Error()))
		return
	}

	url := "./asset/upload/" + upFileName
	jsonOutPut(ctx, 0, "上传成功", url)
}
