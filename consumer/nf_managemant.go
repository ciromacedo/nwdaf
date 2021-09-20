package consumer

import (
	"context"
	"fmt"
	"github.com/free5gc/openapi/Nnrf_NFManagement"
	"github.com/free5gc/openapi/models"
	nwdaf_context "github.com/free5gc/nwdaf/context"
	"github.com/free5gc/nwdaf/factory"
	"net/http"
	"strings"
	"time"
)

func BuildNFInstance(context *nwdaf_context.NWDAFContext) models.NfProfile {
	var profile models.NfProfile
	config := factory.NwdafConfig
	profile.NfInstanceId = context.NfId
	profile.NfType = models.NfType_NWDAF
	profile.NfStatus = models.NfStatus_REGISTERED
	version := config.Info.Version
	tmpVersion := strings.Split(version, ".")
	versionUri := "v" + tmpVersion[0]
	apiPrefix := fmt.Sprintf("%s://%s:%d", context.UriScheme, context.RegisterIPv4, context.SBIPort)
	services := []models.NfService{ //TODO: Outras funções usam um "for" para preencher os serviços.
		{
			ServiceInstanceId: "nwdafdatarepository",        //TODO: Renomear para o ID correto. E excluir o código do serviço de exemplo: ServiceName_NNWDAF_DR
			ServiceName:       "nnwdaf-dr",                  //TODO: Renomear para o serviço correto! ServiceName_NNWDAF_ANALYTICSINFO
			Versions: &[]models.NfServiceVersion{
				{
					ApiFullVersion:  version,
					ApiVersionInUri: versionUri,
				},
			},
			Scheme:          models.UriScheme(context.UriScheme),
			NfServiceStatus: models.NfServiceStatus_REGISTERED,
			ApiPrefix:       apiPrefix,
			IpEndPoints: &[]models.IpEndPoint{
				{
					Ipv4Address: context.RegisterIPv4,
					Transport:   models.TransportProtocol_TCP,
					Port:        int32(context.SBIPort),
				},
			},
		},
		{
			ServiceInstanceId: "analyticsinfo",
			ServiceName:       models.ServiceName_NNWDAF_ANALYTICSINFO,
			Versions: &[]models.NfServiceVersion{
				{
					ApiFullVersion:  version,
					ApiVersionInUri: versionUri,
				},
			},
			Scheme:          models.UriScheme(context.UriScheme),
			NfServiceStatus: models.NfServiceStatus_REGISTERED,
			ApiPrefix:       apiPrefix,
			IpEndPoints: &[]models.IpEndPoint{
				{
					Ipv4Address: context.RegisterIPv4,
					Transport:   models.TransportProtocol_TCP,
					Port:        int32(context.SBIPort),
				},
			},
		},
		{
			ServiceInstanceId: "eventssubscription",
			ServiceName:       models.ServiceName_NNWDAF_EVENTSSUBSCRIPTION,
			Versions: &[]models.NfServiceVersion{
				{
					ApiFullVersion:  version,
					ApiVersionInUri: versionUri,
				},
			},
			Scheme:          models.UriScheme(context.UriScheme),
			NfServiceStatus: models.NfServiceStatus_REGISTERED,
			ApiPrefix:       apiPrefix,
			IpEndPoints: &[]models.IpEndPoint{
				{
					Ipv4Address: context.RegisterIPv4,
					Transport:   models.TransportProtocol_TCP,
					Port:        int32(context.SBIPort),
				},
			},
		},
	}
	profile.NfServices = &services
	// TODO: finish the Nwdaf Info
	/*profile.NwdafInfo = &models.NwdafInfo{
		SupportedDataSets: []models.DataSetId{
			// models.DataSetId_APPLICATION,
			// models.DataSetId_EXPOSURE,
			// models.DataSetId_POLICY,
			models.DataSetId_SUBSCRIPTION,
		},
	}*/
	return profile
}

func SendRegisterNFInstance(nrfUri, nfInstanceId string, profile models.NfProfile) (string, string, error) {

	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath(nrfUri)
	client := Nnrf_NFManagement.NewAPIClient(configuration)
	var resouceNrfUri string
	var retrieveNfInstanceId string

	for {
		_, res, err := client.NFInstanceIDDocumentApi.RegisterNFInstance(context.TODO(), nfInstanceId, profile)
		if err != nil || res == nil {
			//TODO : add log
			fmt.Println(fmt.Errorf("NWDAF register to NRF Error[%s]", err.Error()))
			time.Sleep(2 * time.Second)
			continue
		}
		status := res.StatusCode
		if status == http.StatusOK {
			// NFUpdate
			return resouceNrfUri, retrieveNfInstanceId, err
		} else if status == http.StatusCreated {
			// NFRegister
			resourceUri := res.Header.Get("Location")
			resouceNrfUri = resourceUri[:strings.Index(resourceUri, "/nnrf-nfm/")]
			retrieveNfInstanceId = resourceUri[strings.LastIndex(resourceUri, "/")+1:]
			return resouceNrfUri, retrieveNfInstanceId, err
		} else {
			fmt.Println("handler returned wrong status code", status)
			fmt.Println("NRF return wrong status code", status)
		}
	}
}
