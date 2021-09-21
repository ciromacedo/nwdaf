package datacollection

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/tools/go/analysis/passes/printf/testdata/src/a"
)


func HTTPAmfRegistrationAccept(c *gin.Context) {
	requestBody, err := c.GetRawData()
	if err != nil {
		return
	}
	a.Println(requestBody)

	/*err = openapi.Deserialize(&ueContextRelease, requestBody, "application/json")
	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.CommLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	req := http_wrapper.NewRequest(c.Request, ueContextRelease)
	req.Params["ueContextId"] = c.Params.ByName("ueContextId")
	rsp := producer.HandleReleaseUEContextRequest(req)

	responseBody, err := openapi.Serialize(rsp.Body, "application/json")
	if err != nil {
		logger.CommLog.Errorln(err)
		problemDetails := models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "SYSTEM_FAILURE",
			Detail: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, problemDetails)
	} else {
		c.Data(rsp.Status, "application/json", responseBody)
	}*/
}