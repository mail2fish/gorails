package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mail2fish/gorails/errors"
	"github.com/mail2fish/gorails/route"
)

const (
	MODULE_DEMO errors.ErrorModule = 1
)

type Demo1Params struct {
	ID int `json:"id"`
}

func (p *Demo1Params) Parse(c *gin.Context) errors.Error {
	if err := c.ShouldBindJSON(p); err != nil {
		return errors.NewError(http.StatusBadRequest, errors.THIRD_PARTY, MODULE_DEMO, 1, err.Error(), err)
	}
	return nil
}

type Demo1Response struct {
	Message string `json:"message"`
}

func (r *Demo1Response) Render(c *gin.Context) {
	c.JSON(http.StatusOK, r)
}

func Demo1Handler(c *gin.Context, params route.Params) (route.Response, errors.Error) {
	p := params.(*Demo1Params)
	return &Demo1Response{Message: fmt.Sprintf("ID is %d", p.ID)}, nil
}

type Demo2Params struct {
	Name string `json:"name"`
}

func (p *Demo2Params) Parse(c *gin.Context) errors.Error {
	if err := c.ShouldBindJSON(p); err != nil {
		return errors.NewError(http.StatusBadRequest, errors.THIRD_PARTY, MODULE_DEMO, 1, err.Error(), err)
	}
	return nil
}

type Demo2Response struct {
	Message string `json:"message"`
}

func (r *Demo2Response) Render(c *gin.Context) {
	c.JSON(http.StatusOK, r)
}

func Demo2Handler(c *gin.Context, params route.Params) (route.Response, errors.Error) {
	p := params.(*Demo2Params)
	return &Demo2Response{Message: fmt.Sprintf("Name is %s", p.Name)}, nil
}

func DemoRoute() *gin.Engine {
	router := gin.Default()
	router.POST("/demo1", route.Wrap[*Demo1Params, *Demo1Response](Demo1Handler))
	router.POST("/demo2", route.Wrap[*Demo2Params, *Demo2Response](Demo2Handler))
	return router
}

func main() {
	router := DemoRoute()
	router.Run(":9090")
}
