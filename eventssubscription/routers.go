/*
 * NRF AnalyticsInfo
 *
 * NRF Analytics Info Service
 */

package eventssubscription

import (
	"github.com/ciromacedo/nwdaf/commom"
	"github.com/ciromacedo/nwdaf/logger"
	"github.com/free5gc/logger_util"
	"github.com/gin-gonic/gin"
	"net/http"
)


// Routes is the list of the generated Route.
type Routes []commom.Route

// NewRouter returns a new router.
func NewRouter() *gin.Engine {
	router := logger_util.NewGinWithLogrus(logger.GinLog)
	AddService(router)
	return router
}

func AddService(engine *gin.Engine) *gin.RouterGroup {
	group := engine.Group("")

	for _, route := range routes {
		switch route.Method {
		case "GET":
			group.GET(route.Pattern, route.HandlerFunc)
		case "POST":
			group.POST(route.Pattern, route.HandlerFunc)
		case "PUT":
			group.PUT(route.Pattern, route.HandlerFunc)
		case "DELETE":
			group.DELETE(route.Pattern, route.HandlerFunc)
		case "PATCH":
			group.PATCH(route.Pattern, route.HandlerFunc)
		}
	}

	return group
}

func Index(c *gin.Context) {
	c.String(http.StatusOK, "Event Subscription - Hello World!")
}

var routes = Routes{
	{
		"Index",
		"GET",
		"/",
		IndexHandler,
	},
	{
		"EventsSubscription",
		"GET",
		"/eventssubscription",
		CachedPageHandler,
	},
}
