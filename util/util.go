package util

import (
	"encoding/json"
	"fmt"
	"crypto/tls"
	"github.com/ciromacedo/nwdaf/model"
	"github.com/free5gc/openapi/models"
	"golang.org/x/net/http2"
	"net/http"
	"net"
	"os"
)

func GetConfiguration()model.Config{
	file, fail := os.Open("config/config.json")
	if fail != nil {
		panic(fail.Error())
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := model.Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		panic(err.Error())
	}
	return configuration
}

func SearchNFServiceUri(nfProfile models.NfProfile, serviceName models.ServiceName,
	nfServiceStatus models.NfServiceStatus) (nfUri string, endpoint string, apiVersion string) {
	if nfProfile.NfServices != nil {
		for _, service := range *nfProfile.NfServices {
			if service.ServiceName == serviceName && service.NfServiceStatus == nfServiceStatus {
				if nfProfile.Fqdn != "" {
					nfUri = nfProfile.Fqdn
				} else if service.Fqdn != "" {
					nfUri = service.Fqdn
				} else if service.ApiPrefix != "" {
					nfUri = service.ApiPrefix
				} else if service.IpEndPoints != nil {
					point := (*service.IpEndPoints)[0]
					if point.Ipv4Address != "" {
						nfUri = getSbiUri(service.Scheme, point.Ipv4Address, point.Port)
					} else if len(nfProfile.Ipv4Addresses) != 0 {
						nfUri = getSbiUri(service.Scheme, nfProfile.Ipv4Addresses[0], point.Port)
					}
				}
			}
			if nfUri != "" {
				endpoint = string(service.ServiceName)
				apiVersion = getApiVersion(service.Versions, 1)
				break
			}
		}
	}
	return
}

func getApiVersion(versions *[]models.NfServiceVersion, position int) (version string) {
	for _, v := range *versions {
		version = v.ApiVersionInUri
		break
	}
	return
}

func getSbiUri(scheme models.UriScheme, ipv4Address string, port int32) (uri string) {
	if port != 0 {
		uri = fmt.Sprintf("%s://%s:%d", scheme, ipv4Address, port)
	} else {
		switch scheme {
		case models.UriScheme_HTTP:
			uri = fmt.Sprintf("%s://%s:80", scheme, ipv4Address)
		case models.UriScheme_HTTPS:
			uri = fmt.Sprintf("%s://%s:443", scheme, ipv4Address)
		}
	}
	return
}

func GetHttpConnection()(http.Client){
	client := http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}
	return client
}

func GetServerPort()string{
	ConfigPort := GetConfiguration().Port
	return fmt.Sprintf("%s%d", ":", ConfigPort)
}

func GetMongoDBUri()string{
	return GetConfiguration().MongoURI
}

func GetDBName()string{
	dbName := GetConfiguration().DBName
	return dbName
}
