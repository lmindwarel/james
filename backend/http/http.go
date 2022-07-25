package http

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/lmindwarel/james/backend/datastore"
	"github.com/lmindwarel/james/backend/models"
	"github.com/pkg/errors"

	"github.com/lmindwarel/james/backend/controller"

	"github.com/lmindwarel/james/backend/utils"

	"github.com/gin-gonic/gin"
)

// Account ID of the doer
const CtxDoerID = "doerID"

// Config is the config for the api
type Config struct {
	ServerPort string `json:"serverPort"`
}

// API is the api object
type API struct {
	config Config
	ctrl   *controller.Controller
	ds     *datastore.Datastore
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

var log = utils.GetLogger("http")

// New create new api with the given config
func New(config Config, ctrl *controller.Controller, ds *datastore.Datastore) *API {
	return &API{
		config: config,
		ctrl:   ctrl,
		ds:     ds,
	}
}

// StartServer start the server on core config port
func (a *API) StartServer() error {
	e := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("X-Doer")
	e.Use(cors.New(corsConfig))
	e.Use(gin.Logger())

	a.setupRoutes(e)

	log.Info("Starting server on port %s", a.config.ServerPort)

	return e.Run(":" + a.config.ServerPort)
}

func (a *API) AuthenticatedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		doerID := models.UUID(c.GetHeader("X-Doer"))

		account, err := a.ds.GetAccount(doerID)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, errors.Wrapf(err, "unknow account id %s", doerID))
			return
		}

		c.Set(CtxDoerID, account)
	}
}

func getPaginationQuery(c *gin.Context) (pagination models.Pagination) {
	c.ShouldBindQuery(&pagination)
	if pagination.Limit == 0 {
		pagination.Limit = models.MaxPaginationLimit
	}
	return
}
