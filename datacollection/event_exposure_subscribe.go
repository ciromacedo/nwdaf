package datacollection

import (
	"bytes"
	"fmt"
	"github.com/ciromacedo/nwdaf/consumer"
	nwdaf_context "github.com/ciromacedo/nwdaf/context"
	"github.com/ciromacedo/nwdaf/util"
	"github.com/free5gc/openapi/Nnrf_NFDiscovery"
	"github.com/free5gc/openapi/models"
	"io/ioutil"
	"log"
	"net/http"
)

func InitEventExposureSubscriber(self*nwdaf_context.NWDAFContext) {

	searchOpt := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{};
	// recupera todas as AMFs registradas na NRF
	resp, err := consumer.SendSearchNFInstances(self.NrfUri, models.NfType_AMF, models.NfType_NWDAF, searchOpt);
	if err != nil {
		fmt.Println(err)
	}

	//para cada uma das AMF's registrar no core realiza o subscriber de coleta
	for _, nfProfile := range resp.NfInstances {

		/* localiza a URL do end-point de subscriber com status de REGISTRADO */
		amfUri, endpoint, apiversion := util.SearchNFServiceUri(nfProfile, models.ServiceName_NAMF_EVTS, models.NfServiceStatus_REGISTERED)

		fmt.Println(endpoint)
		fmt.Println(apiversion)

		var buffer bytes.Buffer;

		buffer.WriteString(amfUri);
		buffer.WriteString("/");
		buffer.WriteString(endpoint);
		buffer.WriteString("/");
		buffer.WriteString(apiversion);
		buffer.WriteString("/");
		buffer.WriteString("subscriptions");

		url := buffer.String()

		/*
		 * 1 º os possiveis tipos de eventos p/ AMF estão em AmfEventType
		 */

		jsonData := `
    {	
		"Subscription" : { 	"EventList"	: 
										[{ "Type" : "REGISTRATION_ACCEPT",
                                           "ImmediateFlag" : true}], 
							"EventNotifyUri": "http://127.0.0.1:29599/datacollection/amf-contexts/registration-accept",
							"AnyUE" : true,
							"NfId"  : "NWDAF"
                          },
		"SupportedFeatures"	: "xx"
	}
		
	`

		var jsonStr = []byte(jsonData)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("X-Custom-Header", "myvalue")
		req.Header.Set("Content-Type", "application/json")

		client := util.GetHttpConnection()

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body) // response body is []byte
		fmt.Println(string(body))
	}



}
