package datacollection

import (
	"github.com/ciromacedo/nwdaf/model"
	"github.com/ciromacedo/nwdaf/util"
	"github.com/free5gc/openapi"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)


func HTTPAmfRegistrationAccept(c *gin.Context) {
	var registrationAccept model.RegistrationAccept
	requestBody, err := c.GetRawData()
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadGateway)
		c.Writer.Write([]byte("Internal Error"))
		return
	}

	err = openapi.Deserialize(&registrationAccept, requestBody, "application/json")
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadGateway)
		c.Writer.Write([]byte("Json Parser Error"))
		return
	}

	registrationAccept.Date = time.Now()
	/* registrar na base */
	util.AddRegistrationAccept(&registrationAccept);
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write([]byte("Ok"))
}