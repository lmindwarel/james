package http

import (
	"github.com/gin-gonic/gin"
)

func (a *API) setupRoutes(e *gin.Engine) {

	// public router
	e.GET("/accounts", a.GetAccounts)
	e.POST("/accounts", a.PostAccount)

	// authenticated only
	// Authorization group
	authenticated := e.Group("/")
	// per group middleware! in this case we use the custom created
	// AuthRequired() middleware just in the "authorized" group.
	authenticated.Use(a.AuthenticatedMiddleware())
}
