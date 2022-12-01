package consumer

import (
	"context"
	"fmt"
	"github.com/free5gc/openapi/Nnrf_NFDiscovery"
	"github.com/free5gc/openapi/models"
	"net/http"
)

func SendSearchNFInstances(nrfUri string, targetNfType, requestNfType models.NfType,
	param Nnrf_NFDiscovery.SearchNFInstancesParamOpts) (models.SearchResult, error) {

	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath(nrfUri)
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	var res *http.Response
	result, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requestNfType, &param)
	if res != nil && res.StatusCode == http.StatusTemporaryRedirect {
		err = fmt.Errorf("Temporary Redirect For Non NRF Consumer")
	}
	return result, err
}
