/*
 * NRF AnalyticsInfo
 *
 * NRF Analytics Info Service
 */

package analyticsinfo

import (
	"github.com/free5gc/logger_util"
	"github.com/free5gc/nrf/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Route is the information for every URI.
type Route struct {
	// Name is the name of this Route.
	Name string
	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method string
	// Pattern is the pattern of the URI.
	Pattern string
	// HandlerFunc is the handler function of this route.
	HandlerFunc gin.HandlerFunc
}

// Routes is the list of the generated Route.
type Routes []Route

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

// Index is the index handler.
func Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello World! WWWWWwwwwwwWwwwwwwwwwWwWwWWWwWWwwwwWwWWWWWWWww")
}

// CachedPageHandler is an example handler that returns HTML.
// Reference: https://github.com/gin-gonic/gin/issues/628
func CachedPageHandler(c *gin.Context) {
	//Do what you need to get the cached html
	yourHtmlString := "<h1 style=\"color: #5e9ca0;\">NWDAF</h1>\n<h2 style=\"color: #2e6c80;\"><strong>Service:</strong> Analytics Info</h2>\n<p><strong><a title=\"Back\" href=\"http://localhost:29599/\">&lt;&lt; Index</a></strong></p>\n<p><strong>URL:</strong> <a href=\"http://localhost:29599/analyticsinfo\">http://localhost:29599/analyticsinfo</a></p>\n<p><strong>NF Port:</strong> 29599</p>\n<p><strong>Route file: </strong>src/nwdaf/analyticsinfo/<span style=\"text-decoration: underline;\">routers.go</span></p>\n<p><strong>Route:</strong></p>\n<pre style=\"background-color: #2b2b2b; color: #a9b7c6; font-family: 'JetBrains Mono',monospace; font-size: 9,9pt;\"><span style=\"color: #c7773e;\">var </span>routes = <span style=\"color: #6fafbd;\">Routes</span>{<br />   <span style=\"color: #787878;\">// Route Analytics Info:</span><br />   {<br />      <span style=\"color: #6a8759;\">\"AnalyticsInfo\"</span><span style=\"color: #cc7832;\">,       </span><span style=\"color: #999999;\">// Name</span><span style=\"color: #cc7832;\"><br /></span>      <span style=\"color: #6a8759;\">\"GET\"</span><span style=\"color: #cc7832;\">,                 <span style=\"color: #999999;\">// Method</span><br /></span>      <span style=\"color: #6a8759;\">\"/analyticsinfo\"</span><span style=\"color: #cc7832;\">,      <span style=\"color: #999999;\">// Pattern</span><br /></span>      <span style=\"color: #e6b163;\">CachedPageHandler</span><span style=\"color: #cc7832;\">,     <span style=\"color: #999999;\">// HandlerFunc</span><br /></span>   }<span style=\"color: #cc7832;\">,<br /></span></pre>\n<p><strong>Init file:</strong> src/nwdaf/service/<span style=\"text-decoration: underline;\">init.go</span></p>\n<pre style=\"background-color: #2b2b2b; color: #a9b7c6; font-family: 'JetBrains Mono',monospace; font-size: 9,9pt;\"><span style=\"color: #c7773e;\">func </span>(<span style=\"color: #4eade5;\">nwdaf </span>*<span style=\"color: #6fafbd;\">NWDAF</span>) <span style=\"color: #e6b163;\">Start</span>() {<br /><span style=\"color: #787878;\">   ...<br /><br />   // Order is important for the same route pattern.<br /></span>   <span style=\"color: #afbf7e;\">analyticsinfo</span>.<span style=\"color: #b09d79;\">AddService</span>(router)<strong><br /><br /><span style=\"color: #787878;\">   ...</span><br /><br /></strong></pre>\n<p>&nbsp;</p>"

	//Write your 200 header status (or other status codes, but only WriteHeader once)
	c.Writer.WriteHeader(http.StatusOK)
	//Convert your cached html string to byte array
	c.Writer.Write([]byte(yourHtmlString))
	return
}

var routes = Routes{
	{
		"AnalyticsInfo",
		"GET",
		"/analyticsinfo",
		CachedPageHandler,
	},

	//{
	//	"AccessTokenRequest",
	//	strings.ToUpper("Post"),
	//	"/oauth2/token",
	//	HTTPAccessTokenRequest,
	//},
}
