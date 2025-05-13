package route

import (
	"github.com/gin-gonic/gin"
	"github.com/mail2fish/gorails/errors"
)

type EmptyParams struct {
}

func (p *EmptyParams) Parse(c *gin.Context) errors.Error {
	return nil
}
