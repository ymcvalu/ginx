package ginx

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var _ error = BindError{}

type BindError struct {
	err error
}

func (e BindError) Error() string {
	return e.err.Error()
}

type Renderer interface {
	Render(ctx *gin.Context, data interface{})
}

var _ Renderer = defRenderer{}

type defRenderer struct{}

func (d defRenderer) Render(ctx *gin.Context, data interface{}) {
	switch v := data.(type) {
	case nil:
		ctx.JSON(http.StatusOK, gin.H{
			"code": "0",
			"msg":  "success",
		})

	case BindError:
		log.Printf("failed to bind request parameters: %s", v.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "1",
			"msg":  v.Error(),
		})

	case error:
		log.Printf("failed to handle requst for api[%s]: %s", ctx.Request.URL.Path, v.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": "500",
			"msg":  v.Error(),
		})

	default:
		ctx.JSON(http.StatusOK, v)
	}
}
