package gorails

import (
	"github.com/gin-gonic/gin"
	"github.com/mail2fish/gorails/route"
)

type Params = route.Params
type Response = route.Response
type GoRailsHandlerFun = route.GoRailsHandlerFun

func Wrap[T route.Params, U route.Response](f GoRailsHandlerFun) gin.HandlerFunc {
	return route.Wrap[T, U](f)
}
