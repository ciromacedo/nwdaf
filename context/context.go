package context

import (
	"fmt"
	"sync"

	"github.com/free5gc/openapi/models"
)

var nwdafContext = NWDAFContext{}

type subsId = string

type NWDAFServiceType int

const (
	NNWDAF_DR NWDAFServiceType = iota
)

func init() {
	NWDAF_Self().Name = "nwdaf"
	NWDAF_Self().EeSubscriptionIDGenerator = 1
	NWDAF_Self().SdmSubscriptionIDGenerator = 1
	NWDAF_Self().SubscriptionDataSubscriptionIDGenerator = 1
	NWDAF_Self().PolicyDataSubscriptionIDGenerator = 1
	NWDAF_Self().SubscriptionDataSubscriptions = make(map[subsId]*models.SubscriptionDataSubscriptions)
	NWDAF_Self().PolicyDataSubscriptions = make(map[subsId]*models.PolicyDataSubscription)
}

type NWDAFContext struct {
	Name                                    string
	UriScheme                               models.UriScheme
	BindingIPv4                             string
	SBIPort                                 int
	RegisterIPv4                            string // IP register to NRF
	HttpIPv6Address                         string
	NfId                                    string
	NrfUri                                  string
	EeSubscriptionIDGenerator               int
	SdmSubscriptionIDGenerator              int
	PolicyDataSubscriptionIDGenerator       int
	UESubsCollection                        sync.Map //map[ueId]*UESubsData
	UEGroupCollection                       sync.Map //map[ueGroupId]*UEGroupSubsData
	SubscriptionDataSubscriptionIDGenerator int
	SubscriptionDataSubscriptions           map[subsId]*models.SubscriptionDataSubscriptions
	PolicyDataSubscriptions                 map[subsId]*models.PolicyDataSubscription
}

type UESubsData struct {
	EeSubscriptionCollection map[subsId]*EeSubscriptionCollection
	SdmSubscriptions         map[subsId]*models.SdmSubscription
}

type UEGroupSubsData struct {
	EeSubscriptions map[subsId]*models.EeSubscription
}

type EeSubscriptionCollection struct {
	EeSubscriptions      *models.EeSubscription
	AmfSubscriptionInfos []models.AmfSubscriptionInfo
}

// Reset NWDAF Context
func (context *NWDAFContext) Reset() {
	context.UESubsCollection.Range(func(key, value interface{}) bool {
		context.UESubsCollection.Delete(key)
		return true
	})
	context.UEGroupCollection.Range(func(key, value interface{}) bool {
		context.UEGroupCollection.Delete(key)
		return true
	})
	for key := range context.SubscriptionDataSubscriptions {
		delete(context.SubscriptionDataSubscriptions, key)
	}
	for key := range context.PolicyDataSubscriptions {
		delete(context.PolicyDataSubscriptions, key)
	}
	context.EeSubscriptionIDGenerator = 1
	context.SdmSubscriptionIDGenerator = 1
	context.SubscriptionDataSubscriptionIDGenerator = 1
	context.PolicyDataSubscriptionIDGenerator = 1
	context.UriScheme = models.UriScheme_HTTPS
	context.Name = "nwdaf"
}

func (context *NWDAFContext) GetIPv4Uri() string {
	return fmt.Sprintf("%s://%s:%d", context.UriScheme, context.RegisterIPv4, context.SBIPort)
}

func (context *NWDAFContext) GetIPv4GroupUri(nwdafServiceType NWDAFServiceType) string {
	var serviceUri string

	switch nwdafServiceType {
	case NNWDAF_DR:
		serviceUri = "/nnwdaf-dr/v1"
	default:
		serviceUri = ""
	}

	return fmt.Sprintf("%s://%s:%d%s", context.UriScheme, context.RegisterIPv4, context.SBIPort, serviceUri)
}

// Create new NWDAF context
func NWDAF_Self() *NWDAFContext {
	return &nwdafContext
}
