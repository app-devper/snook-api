package errcode

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	ErrCode string `json:"errcode"`
	Error   string `json:"error"`
}

func Abort(ctx *gin.Context, httpStatus int, code string, msg string) {
	ctx.AbortWithStatusJSON(httpStatus, AppError{ErrCode: code, Error: msg})
}

func AbortByCode(ctx *gin.Context, code string, msg string) {
	info, ok := GetCodeInfo(code)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, AppError{ErrCode: code, Error: msg})
		return
	}
	if msg == "" {
		msg = info.Description
	}
	ctx.AbortWithStatusJSON(info.HttpStatus, AppError{ErrCode: code, Error: msg})
}
