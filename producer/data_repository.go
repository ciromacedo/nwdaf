package producer

import (
	"encoding/json"
	"fmt"
	"free5gc/lib/MongoDBLibrary"
	"free5gc/lib/http_wrapper"
	"free5gc/lib/openapi/models"
	nwdaf_context "free5gc/src/nwdaf/context"
	"free5gc/src/nwdaf/logger"
	"free5gc/src/nwdaf/util"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
)

var CurrentResourceUri string

func HandleCreateAccessAndMobilityData(request *http_wrapper.Request) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusOK, nil, map[string]interface{}{})
}

func HandleDeleteAccessAndMobilityData(request *http_wrapper.Request) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusOK, nil, map[string]interface{}{})
}

func HandleQueryAccessAndMobilityData(request *http_wrapper.Request) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusOK, nil, map[string]interface{}{})
}

func HandleQueryAmData(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QueryAmData")

	collName := "subscriptionData.provisionedData.amData"
	ueId := request.Params["ueId"]
	servingPlmnId := request.Params["servingPlmnId"]
	response, problemDetails := QueryAmDataProcedure(collName, ueId, servingPlmnId)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
}

func QueryAmDataProcedure(collName string, ueId string, servingPlmnId string) (*map[string]interface{},
	*models.ProblemDetails) {
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
	accessAndMobilitySubscriptionData := MongoDBLibrary.RestfulAPIGetOne(collName, filter)
	if accessAndMobilitySubscriptionData != nil {
		return &accessAndMobilitySubscriptionData, nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}
}

func HandleAmfContext3gpp(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle AmfContext3gpp")
	collName := "subscriptionData.contextData.amf3gppAccess"
	patchItem := request.Body.([]models.PatchItem)
	ueId := request.Params["ueId"]

	problemDetails := AmfContext3gppProcedure(collName, ueId, patchItem)
	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
}

func AmfContext3gppProcedure(collName string, ueId string, patchItem []models.PatchItem) *models.ProblemDetails {
	filter := bson.M{"ueId": ueId}
	origValue := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	patchJSON, err := json.Marshal(patchItem)
	if err != nil {
		logger.DataRepoLog.Error(err)
	}
	success := MongoDBLibrary.RestfulAPIJSONPatch(collName, filter, patchJSON)

	if success {
		newValue := MongoDBLibrary.RestfulAPIGetOne(collName, filter)
		PreHandleOnDataChangeNotify(ueId, CurrentResourceUri, patchItem, origValue, newValue)
		return nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "MODIFY_NOT_ALLOWED",
			Status: http.StatusNotFound,
		}
		return problemDetails
	}
}

func HandleCreateAmfContext3gpp(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle CreateAmfContext3gpp")

	Amf3GppAccessRegistration := request.Body.(models.Amf3GppAccessRegistration)
	ueId := request.Params["ueId"]
	collName := "subscriptionData.contextData.amf3gppAccess"

	CreateAmfContext3gppProcedure(collName, ueId, Amf3GppAccessRegistration)

	return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
}

func CreateAmfContext3gppProcedure(collName string, ueId string,
	Amf3GppAccessRegistration models.Amf3GppAccessRegistration) {

	filter := bson.M{"ueId": ueId}
	putData := util.ToBsonM(Amf3GppAccessRegistration)
	putData["ueId"] = ueId

	MongoDBLibrary.RestfulAPIPutOne(collName, filter, putData)
}

func HandleQueryAmfContext3gpp(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QueryAmfContext3gpp")

	ueId := request.Params["ueId"]
	collName := "subscriptionData.contextData.amf3gppAccess"

	response, problemDetails := QueryAmfContext3gppProcedure(collName, ueId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(http.StatusForbidden, nil, problemDetails)
}

func QueryAmfContext3gppProcedure(collName string, ueId string) (*map[string]interface{}, *models.ProblemDetails) {

	filter := bson.M{"ueId": ueId}
	amf3GppAccessRegistration := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	if amf3GppAccessRegistration != nil {
		return &amf3GppAccessRegistration, nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}
}

func HandleAmfContextNon3gpp(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle AmfContextNon3gpp")

	ueId := request.Params["ueId"]
	collName := "subscriptionData.contextData.amfNon3gppAccess"
	patchItem := request.Body.([]models.PatchItem)
	filter := bson.M{"ueId": ueId}

	problemDetails := AmfContextNon3gppProcedure(ueId, collName, patchItem, filter)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
}

func AmfContextNon3gppProcedure(ueId string, collName string, patchItem []models.PatchItem,
	filter bson.M) *models.ProblemDetails {
	origValue := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	patchJSON, err := json.Marshal(patchItem)
	if err != nil {
		logger.DataRepoLog.Error(err)
	}
	success := MongoDBLibrary.RestfulAPIJSONPatch(collName, filter, patchJSON)
	if success {
		newValue := MongoDBLibrary.RestfulAPIGetOne(collName, filter)
		PreHandleOnDataChangeNotify(ueId, CurrentResourceUri, patchItem, origValue, newValue)
		return nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "MODIFY_NOT_ALLOWED",
			Status: http.StatusForbidden,
		}
		return problemDetails
	}
}

func HandleCreateAmfContextNon3gpp(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle CreateAmfContextNon3gpp")

	AmfNon3GppAccessRegistration := request.Body.(models.AmfNon3GppAccessRegistration)
	collName := "subscriptionData.contextData.amfNon3gppAccess"
	ueId := request.Params["ueId"]

	CreateAmfContextNon3gppProcedure(AmfNon3GppAccessRegistration, collName, ueId)

	return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
}

func CreateAmfContextNon3gppProcedure(AmfNon3GppAccessRegistration models.AmfNon3GppAccessRegistration,
	collName string, ueId string) {
	putData := util.ToBsonM(AmfNon3GppAccessRegistration)
	putData["ueId"] = ueId
	filter := bson.M{"ueId": ueId}

	MongoDBLibrary.RestfulAPIPutOne(collName, filter, putData)

}

func HandleQueryAmfContextNon3gpp(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QueryAmfContextNon3gpp")

	collName := "subscriptionData.contextData.amfNon3gppAccess"
	ueId := request.Params["ueId"]

	response, problemDetails := QueryAmfContextNon3gppProcedure(collName, ueId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(http.StatusForbidden, nil, problemDetails)
}

func QueryAmfContextNon3gppProcedure(collName string, ueId string) (*map[string]interface{}, *models.ProblemDetails) {
	filter := bson.M{"ueId": ueId}
	response := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	if response != nil {
		return &response, nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}
}

func HandleModifyAuthentication(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle ModifyAuthentication")

	collName := "subscriptionData.authenticationData.authenticationSubscription"
	ueId := request.Params["ueId"]
	patchItem := request.Body.([]models.PatchItem)

	problemDetails := ModifyAuthenticationProcedure(collName, ueId, patchItem)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
}

func ModifyAuthenticationProcedure(collName string, ueId string, patchItem []models.PatchItem) *models.ProblemDetails {
	filter := bson.M{"ueId": ueId}
	origValue := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	patchJSON, err := json.Marshal(patchItem)
	if err != nil {
		logger.DataRepoLog.Error(err)
	}
	success := MongoDBLibrary.RestfulAPIJSONPatch(collName, filter, patchJSON)

	if success {
		newValue := MongoDBLibrary.RestfulAPIGetOne(collName, filter)
		PreHandleOnDataChangeNotify(ueId, CurrentResourceUri, patchItem, origValue, newValue)
		return nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "MODIFY_NOT_ALLOWED",
			Status: http.StatusForbidden,
		}
		return problemDetails
	}
}

func HandleQueryAuthSubsData(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QueryAuthSubsData")

	collName := "subscriptionData.authenticationData.authenticationSubscription"
	ueId := request.Params["ueId"]

	response, problemDetails := QueryAuthSubsDataProcedure(collName, ueId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func QueryAuthSubsDataProcedure(collName string, ueId string) (map[string]interface{}, *models.ProblemDetails) {
	filter := bson.M{"ueId": ueId}

	authenticationSubscription := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	if authenticationSubscription != nil {
		return authenticationSubscription, nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}
}

func HandleCreateAuthenticationSoR(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle CreateAuthenticationSoR")
	putData := util.ToBsonM(request.Body)
	ueId := request.Params["ueId"]
	collName := "subscriptionData.ueUpdateConfirmationData.sorData"

	CreateAuthenticationSoRProcedure(collName, ueId, putData)

	return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
}

func CreateAuthenticationSoRProcedure(collName string, ueId string, putData bson.M) {
	filter := bson.M{"ueId": ueId}
	putData["ueId"] = ueId

	MongoDBLibrary.RestfulAPIPutOne(collName, filter, putData)
}

func HandleQueryAuthSoR(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QueryAuthSoR")

	ueId := request.Params["ueId"]
	collName := "subscriptionData.ueUpdateConfirmationData.sorData"

	response, problemDetails := QueryAuthSoRProcedure(collName, ueId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func QueryAuthSoRProcedure(collName string, ueId string) (map[string]interface{}, *models.ProblemDetails) {
	filter := bson.M{"ueId": ueId}

	sorData := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	if sorData != nil {
		return sorData, nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}
}

func HandleCreateAuthenticationStatus(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle CreateAuthenticationStatus")

	putData := util.ToBsonM(request.Body)
	ueId := request.Params["ueId"]
	collName := "subscriptionData.authenticationData.authenticationStatus"

	CreateAuthenticationStatusProcedure(collName, ueId, putData)

	return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
}

func CreateAuthenticationStatusProcedure(collName string, ueId string, putData bson.M) {
	filter := bson.M{"ueId": ueId}
	putData["ueId"] = ueId

	MongoDBLibrary.RestfulAPIPutOne(collName, filter, putData)
}

func HandleQueryAuthenticationStatus(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QueryAuthenticationStatus")

	ueId := request.Params["ueId"]
	collName := "subscriptionData.authenticationData.authenticationStatus"

	response, problemDetails := QueryAuthenticationStatusProcedure(collName, ueId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func QueryAuthenticationStatusProcedure(collName string, ueId string) (*map[string]interface{},
	*models.ProblemDetails) {
	filter := bson.M{"ueId": ueId}

	authEvent := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	if authEvent != nil {
		return &authEvent, nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}
}

func HandleApplicationDataInfluenceDataGet(request *http_wrapper.Request) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusOK, nil, map[string]interface{}{})
}

func HandleApplicationDataInfluenceDataInfluenceIdDelete(request *http_wrapper.Request) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusOK, nil, map[string]interface{}{})
}

func HandleApplicationDataInfluenceDataInfluenceIdPatch(request *http_wrapper.Request) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusOK, nil, map[string]interface{}{})
}

func HandleApplicationDataInfluenceDataInfluenceIdPut(request *http_wrapper.Request) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusOK, nil, map[string]interface{}{})
}

func HandleApplicationDataInfluenceDataSubsToNotifyGet(request *http_wrapper.Request) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusOK, nil, map[string]interface{}{})
}

func HandleApplicationDataInfluenceDataSubsToNotifyPost(request *http_wrapper.Request) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusOK, nil, map[string]interface{}{})
}

func HandleApplicationDataInfluenceDataSubsToNotifySubscriptionIdDelete(
	request *http_wrapper.Request) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusOK, nil, map[string]interface{}{})
}

func HandleApplicationDataInfluenceDataSubsToNotifySubscriptionIdGet(
	request *http_wrapper.Request) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusOK, nil, map[string]interface{}{})
}

func HandleApplicationDataInfluenceDataSubsToNotifySubscriptionIdPut(
	request *http_wrapper.Request) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusOK, nil, map[string]interface{}{})
}

func HandleApplicationDataPfdsAppIdDelete(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle ApplicationDataPfdsAppIdDelete")

	collName := "applicationData.pfds"
	pfdsAppId := request.Params["pfdsAppId"]

	ApplicationDataPfdsAppIdDeleteProcedure(collName, pfdsAppId)

	return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
}

func ApplicationDataPfdsAppIdDeleteProcedure(collName string, pfdsAppId string) {
	filter := bson.M{"applicationId": pfdsAppId}
	MongoDBLibrary.RestfulAPIDeleteOne(collName, filter)
}

func HandleApplicationDataPfdsAppIdGet(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle ApplicationDataPfdsAppIdGet")

	collName := "applicationData.pfds"
	pfdsAppId := request.Params["pfdsAppId"]

	response, problemDetails := ApplicationDataPfdsAppIdGetProcedure(collName, pfdsAppId)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
}

func ApplicationDataPfdsAppIdGetProcedure(collName string, pfdsAppId string) (*map[string]interface{},
	*models.ProblemDetails) {
	filter := bson.M{"applicationId": pfdsAppId}

	getData := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	if getData != nil {
		delete(getData, "_id")
		return &getData, nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "DATA_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}
}

func HandleApplicationDataPfdsAppIdPut(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle ApplicationDataPfdsAppIdPut")

	collName := "applicationData.pfds"
	pfdsAppId := request.Params["pfdsAppId"]
	pfdDataForApp := request.Body.(models.PfdDataForApp)

	response, status := ApplicationDataPfdsAppIdPutProcedure(collName, pfdsAppId, pfdDataForApp)
	if status == http.StatusOK {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else {
		return http_wrapper.NewResponse(http.StatusCreated, nil, response)
	}
}

func ApplicationDataPfdsAppIdPutProcedure(collName string, pfdsAppId string,
	PfdDataForApp models.PfdDataForApp) (bson.M, int) {
	putData := util.ToBsonM(PfdDataForApp)
	filter := bson.M{"applicationId": pfdsAppId}

	isExisted := MongoDBLibrary.RestfulAPIPutOne(collName, filter, putData)

	if isExisted {
		return putData, http.StatusOK
		//PreHandlePolicyDataChangeNotification("", pfdsAppId, body)
	} else {
		return putData, http.StatusCreated
	}
}

func HandleApplicationDataPfdsGet(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle ApplicationDataPfdsGet")

	pfdsAppIdArray := request.Query["appId"]
	collName := "applicationData.pfds"

	response := ApplicationDataPfdsGetProcedure(collName, pfdsAppIdArray)

	return http_wrapper.NewResponse(http.StatusOK, nil, response)
}

func ApplicationDataPfdsGetProcedure(collName string, pfdsAppIdArray []string) (response *[]map[string]interface{}) {
	filter := bson.M{}

	var pfdsArray []map[string]interface{}
	if len(pfdsAppIdArray) == 0 {
		pfdsArray = MongoDBLibrary.RestfulAPIGetMany(collName, filter)
		for i := 0; i < len(pfdsArray); i++ {
			delete(pfdsArray[i], "_id")
		}
	} else {
		for _, e := range pfdsAppIdArray {
			filter["applicationId"] = e
			getData := MongoDBLibrary.RestfulAPIGetOne(collName, filter)
			if getData != nil {
				delete(getData, "_id")
				pfdsArray = append(pfdsArray, getData)
			}
		}
	}
	return &pfdsArray
}

func HandleExposureDataSubsToNotifyPost(request *http_wrapper.Request) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusOK, nil, map[string]interface{}{})
}

func HandleExposureDataSubsToNotifySubIdDelete(request *http_wrapper.Request) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusOK, nil, map[string]interface{}{})
}

func HandleExposureDataSubsToNotifySubIdPut(request *http_wrapper.Request) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusOK, nil, map[string]interface{}{})
}

func HandlePolicyDataBdtDataBdtReferenceIdDelete(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle PolicyDataBdtDataBdtReferenceIdDelete")

	collName := "policyData.bdtData"
	bdtReferenceId := request.Params["bdtReferenceId"]

	PolicyDataBdtDataBdtReferenceIdDeleteProcedure(collName, bdtReferenceId)
	return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
}

func PolicyDataBdtDataBdtReferenceIdDeleteProcedure(collName string, bdtReferenceId string) {
	filter := bson.M{"bdtReferenceId": bdtReferenceId}
	MongoDBLibrary.RestfulAPIDeleteOne(collName, filter)
}

func HandlePolicyDataBdtDataBdtReferenceIdGet(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle PolicyDataBdtDataBdtReferenceIdGet")

	collName := "policyData.bdtData"
	bdtReferenceId := request.Params["bdtReferenceId"]

	response, problemDetails := PolicyDataBdtDataBdtReferenceIdGetProcedure(collName, bdtReferenceId)
	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func PolicyDataBdtDataBdtReferenceIdGetProcedure(collName string, bdtReferenceId string) (*map[string]interface{},
	*models.ProblemDetails) {
	filter := bson.M{"bdtReferenceId": bdtReferenceId}

	bdtData := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	if bdtData != nil {
		return &bdtData, nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "DATA_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}
}

func HandlePolicyDataBdtDataBdtReferenceIdPut(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle PolicyDataBdtDataBdtReferenceIdPut")

	collName := "policyData.bdtData"
	bdtReferenceId := request.Params["bdtReferenceId"]
	bdtData := request.Body.(models.BdtData)

	response := PolicyDataBdtDataBdtReferenceIdPutProcedure(collName, bdtReferenceId, bdtData)
	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	}
	problemDetails := models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func PolicyDataBdtDataBdtReferenceIdPutProcedure(collName string, bdtReferenceId string,
	bdtData models.BdtData) bson.M {
	putData := util.ToBsonM(bdtData)
	putData["bdtReferenceId"] = bdtReferenceId
	filter := bson.M{"bdtReferenceId": bdtReferenceId}

	isExisted := MongoDBLibrary.RestfulAPIPutOne(collName, filter, putData)

	if isExisted {
		PreHandlePolicyDataChangeNotification("", bdtReferenceId, bdtData)
		return putData
	} else {
		return putData
	}
}

func HandlePolicyDataBdtDataGet(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle PolicyDataBdtDataGet")

	collName := "policyData.bdtData"

	response := PolicyDataBdtDataGetProcedure(collName)
	return http_wrapper.NewResponse(http.StatusOK, nil, response)
}

func PolicyDataBdtDataGetProcedure(collName string) (response *[]map[string]interface{}) {
	filter := bson.M{}
	bdtDataArray := MongoDBLibrary.RestfulAPIGetMany(collName, filter)
	return &bdtDataArray
}

func HandlePolicyDataPlmnsPlmnIdUePolicySetGet(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle PolicyDataPlmnsPlmnIdUePolicySetGet")

	collName := "policyData.plmns.uePolicySet"
	plmnId := request.Params["plmnId"]

	response, problemDetails := PolicyDataPlmnsPlmnIdUePolicySetGetProcedure(collName, plmnId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func PolicyDataPlmnsPlmnIdUePolicySetGetProcedure(collName string,
	plmnId string) (*map[string]interface{}, *models.ProblemDetails) {
	filter := bson.M{"plmnId": plmnId}
	uePolicySet := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	if uePolicySet != nil {
		return &uePolicySet, nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}
}

func HandlePolicyDataSponsorConnectivityDataSponsorIdGet(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle PolicyDataSponsorConnectivityDataSponsorIdGet")

	collName := "policyData.sponsorConnectivityData"
	sponsorId := request.Params["sponsorId"]

	response, status := PolicyDataSponsorConnectivityDataSponsorIdGetProcedure(collName, sponsorId)

	if status == http.StatusOK {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if status == http.StatusNoContent {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	}
	problemDetails := models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(http.StatusForbidden, nil, problemDetails)
}

func PolicyDataSponsorConnectivityDataSponsorIdGetProcedure(collName string,
	sponsorId string) (*map[string]interface{}, int) {
	filter := bson.M{"sponsorId": sponsorId}

	sponsorConnectivityData := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	if sponsorConnectivityData != nil {
		return &sponsorConnectivityData, http.StatusOK
	} else {
		return nil, http.StatusNoContent
	}
}

func HandlePolicyDataSubsToNotifyPost(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle PolicyDataSubsToNotifyPost")

	PolicyDataSubscription := request.Body.(models.PolicyDataSubscription)

	locationHeader := PolicyDataSubsToNotifyPostProcedure(PolicyDataSubscription)

	headers := http.Header{
		"Location": {locationHeader},
	}
	return http_wrapper.NewResponse(http.StatusCreated, headers, PolicyDataSubscription)
}

func PolicyDataSubsToNotifyPostProcedure(PolicyDataSubscription models.PolicyDataSubscription) string {
	nwdafSelf := nwdaf_context.NWDAF_Self()

	newSubscriptionID := strconv.Itoa(nwdafSelf.PolicyDataSubscriptionIDGenerator)
	nwdafSelf.PolicyDataSubscriptions[newSubscriptionID] = &PolicyDataSubscription
	nwdafSelf.PolicyDataSubscriptionIDGenerator++

	/* Contains the URI of the newly created resource, according
	   to the structure: {apiRoot}/subscription-data/subs-to-notify/{subsId} */
	locationHeader := fmt.Sprintf("%s/policy-data/subs-to-notify/%s", nwdafSelf.GetIPv4GroupUri(nwdaf_context.NNWDAF_DR),
		newSubscriptionID)

	return locationHeader
}

func HandlePolicyDataSubsToNotifySubsIdDelete(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle PolicyDataSubsToNotifySubsIdDelete")

	subsId := request.Params["subsId"]

	problemDetails := PolicyDataSubsToNotifySubsIdDeleteProcedure(subsId)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
}

func PolicyDataSubsToNotifySubsIdDeleteProcedure(subsId string) (problemDetails *models.ProblemDetails) {
	nwdafSelf := nwdaf_context.NWDAF_Self()
	_, ok := nwdafSelf.PolicyDataSubscriptions[subsId]
	if !ok {
		problemDetails = &models.ProblemDetails{
			Cause:  "SUBSCRIPTION_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return problemDetails
	}
	delete(nwdafSelf.PolicyDataSubscriptions, subsId)

	return nil
}

func HandlePolicyDataSubsToNotifySubsIdPut(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle PolicyDataSubsToNotifySubsIdPut")

	subsId := request.Params["subsId"]
	policyDataSubscription := request.Body.(models.PolicyDataSubscription)

	response, problemDetails := PolicyDataSubsToNotifySubsIdPutProcedure(subsId, policyDataSubscription)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
}

func PolicyDataSubsToNotifySubsIdPutProcedure(subsId string,
	policyDataSubscription models.PolicyDataSubscription) (*models.PolicyDataSubscription, *models.ProblemDetails) {
	nwdafSelf := nwdaf_context.NWDAF_Self()
	_, ok := nwdafSelf.PolicyDataSubscriptions[subsId]
	if !ok {
		problemDetails := &models.ProblemDetails{
			Cause:  "SUBSCRIPTION_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}

	nwdafSelf.PolicyDataSubscriptions[subsId] = &policyDataSubscription

	return &policyDataSubscription, nil
}

func HandlePolicyDataUesUeIdAmDataGet(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle PolicyDataUesUeIdAmDataGet")

	collName := "policyData.ues.amData"
	ueId := request.Params["ueId"]

	response, problemDetails := PolicyDataUesUeIdAmDataGetProcedure(collName, ueId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func PolicyDataUesUeIdAmDataGetProcedure(collName string,
	ueId string) (*map[string]interface{}, *models.ProblemDetails) {
	filter := bson.M{"ueId": ueId}

	amPolicyData := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	if amPolicyData != nil {
		return &amPolicyData, nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}
}

func HandlePolicyDataUesUeIdOperatorSpecificDataGet(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle PolicyDataUesUeIdOperatorSpecificDataGet")

	collName := "policyData.ues.operatorSpecificData"
	ueId := request.Params["ueId"]

	response, problemDetails := PolicyDataUesUeIdOperatorSpecificDataGetProcedure(collName, ueId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func PolicyDataUesUeIdOperatorSpecificDataGetProcedure(collName string,
	ueId string) (*interface{}, *models.ProblemDetails) {
	filter := bson.M{"ueId": ueId}

	operatorSpecificDataContainerMapCover := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	if operatorSpecificDataContainerMapCover != nil {
		operatorSpecificDataContainerMap := operatorSpecificDataContainerMapCover["operatorSpecificDataContainerMap"]
		return &operatorSpecificDataContainerMap, nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}
}

func HandlePolicyDataUesUeIdOperatorSpecificDataPatch(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle PolicyDataUesUeIdOperatorSpecificDataPatch")

	collName := "policyData.ues.operatorSpecificData"
	ueId := request.Params["ueId"]
	patchItem := request.Body.([]models.PatchItem)

	problemDetails := PolicyDataUesUeIdOperatorSpecificDataPatchProcedure(collName, ueId, patchItem)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
}

func PolicyDataUesUeIdOperatorSpecificDataPatchProcedure(collName string, ueId string,
	patchItem []models.PatchItem) *models.ProblemDetails {
	filter := bson.M{"ueId": ueId}

	patchJSON, err := json.Marshal(patchItem)
	if err != nil {
		logger.DataRepoLog.Warnln(err)
	}

	success := MongoDBLibrary.RestfulAPIJSONPatchExtend(collName, filter, patchJSON,
		"operatorSpecificDataContainerMap")

	if success {
		return nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "MODIFY_NOT_ALLOWED",
			Status: http.StatusForbidden,
		}
		return problemDetails
	}
}

func HandlePolicyDataUesUeIdOperatorSpecificDataPut(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle PolicyDataUesUeIdOperatorSpecificDataPut")

	// json.NewDecoder(c.Request.Body).Decode(&operatorSpecificDataContainerMap)

	collName := "policyData.ues.operatorSpecificData"
	ueId := request.Params["ueId"]
	OperatorSpecificDataContainer := request.Body.(map[string]models.OperatorSpecificDataContainer)

	PolicyDataUesUeIdOperatorSpecificDataPutProcedure(collName, ueId, OperatorSpecificDataContainer)

	return http_wrapper.NewResponse(http.StatusOK, nil, map[string]interface{}{})
}

func PolicyDataUesUeIdOperatorSpecificDataPutProcedure(collName string, ueId string,
	OperatorSpecificDataContainer map[string]models.OperatorSpecificDataContainer) {
	filter := bson.M{"ueId": ueId}

	putData := map[string]interface{}{"operatorSpecificDataContainerMap": OperatorSpecificDataContainer}
	putData["ueId"] = ueId

	MongoDBLibrary.RestfulAPIPutOne(collName, filter, putData)
}

func HandlePolicyDataUesUeIdSmDataGet(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle PolicyDataUesUeIdSmDataGet")

	collName := "policyData.ues.smData"
	ueId := request.Params["ueId"]
	sNssai := models.Snssai{}
	sNssaiQuery := request.Query.Get("snssai")
	err := json.Unmarshal([]byte(sNssaiQuery), &sNssai)
	if err != nil {
		logger.DataRepoLog.Warnln(err)
	}
	dnn := request.Query.Get("dnn")

	response, problemDetails := PolicyDataUesUeIdSmDataGetProcedure(collName, ueId, sNssai, dnn)
	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func PolicyDataUesUeIdSmDataGetProcedure(collName string, ueId string, snssai models.Snssai,
	dnn string) (*models.SmPolicyData, *models.ProblemDetails) {
	filter := bson.M{"ueId": ueId}

	if !reflect.DeepEqual(snssai, models.Snssai{}) {
		filter["smPolicySnssaiData."+util.SnssaiModelsToHex(snssai)] = bson.M{"$exists": true}
	}
	if !reflect.DeepEqual(snssai, models.Snssai{}) && dnn != "" {
		filter["smPolicySnssaiData."+util.SnssaiModelsToHex(snssai)+".smPolicyDnnData."+dnn] = bson.M{"$exists": true}
	}

	smPolicyData := MongoDBLibrary.RestfulAPIGetOne(collName, filter)
	if smPolicyData != nil {
		var smPolicyDataResp models.SmPolicyData
		err := json.Unmarshal(util.MapToByte(smPolicyData), &smPolicyDataResp)
		if err != nil {
			logger.DataRepoLog.Warnln(err)
		}
		{
			collName := "policyData.ues.smData.usageMonData"
			filter := bson.M{"ueId": ueId}
			usageMonDataMapArray := MongoDBLibrary.RestfulAPIGetMany(collName, filter)

			if !reflect.DeepEqual(usageMonDataMapArray, []map[string]interface{}{}) {
				var usageMonDataArray []models.UsageMonData
				err = json.Unmarshal(util.MapArrayToByte(usageMonDataMapArray), &usageMonDataArray)
				if err != nil {
					logger.DataRepoLog.Warnln(err)
				}
				smPolicyDataResp.UmData = make(map[string]models.UsageMonData)
				for _, element := range usageMonDataArray {
					smPolicyDataResp.UmData[element.LimitId] = element
				}
			}
		}
		return &smPolicyDataResp, nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}
}

func HandlePolicyDataUesUeIdSmDataPatch(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle PolicyDataUesUeIdSmDataPatch")

	collName := "policyData.ues.smData.usageMonData"
	ueId := request.Params["ueId"]
	usageMonData := request.Body.(map[string]models.UsageMonData)

	problemDetails := PolicyDataUesUeIdSmDataPatchProcedure(collName, ueId, usageMonData)
	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
}

func PolicyDataUesUeIdSmDataPatchProcedure(collName string, ueId string,
	UsageMonData map[string]models.UsageMonData) *models.ProblemDetails {
	filter := bson.M{"ueId": ueId}

	successAll := true
	for k, usageMonData := range UsageMonData {
		limitId := k
		filterTmp := bson.M{"ueId": ueId, "limitId": limitId}
		success := MongoDBLibrary.RestfulAPIMergePatch(collName, filterTmp, util.ToBsonM(usageMonData))
		if !success {
			successAll = false
		} else {
			var usageMonData models.UsageMonData
			usageMonDataBsonM := MongoDBLibrary.RestfulAPIGetOne(collName, filter)
			err := json.Unmarshal(util.MapToByte(usageMonDataBsonM), &usageMonData)
			if err != nil {
				logger.DataRepoLog.Warnln(err)
			}
			PreHandlePolicyDataChangeNotification(ueId, limitId, usageMonData)
		}
	}

	if successAll {
		smPolicyDataBsonM := MongoDBLibrary.RestfulAPIGetOne(collName, filter)
		var smPolicyData models.SmPolicyData
		err := json.Unmarshal(util.MapToByte(smPolicyDataBsonM), &smPolicyData)
		if err != nil {
			logger.DataRepoLog.Warnln(err)
		}
		{
			collName := "policyData.ues.smData.usageMonData"
			filter := bson.M{"ueId": ueId}
			usageMonDataMapArray := MongoDBLibrary.RestfulAPIGetMany(collName, filter)

			if !reflect.DeepEqual(usageMonDataMapArray, []map[string]interface{}{}) {
				var usageMonDataArray []models.UsageMonData
				err = json.Unmarshal(util.MapArrayToByte(usageMonDataMapArray), &usageMonDataArray)
				if err != nil {
					logger.DataRepoLog.Warnln(err)
				}
				smPolicyData.UmData = make(map[string]models.UsageMonData)
				for _, element := range usageMonDataArray {
					smPolicyData.UmData[element.LimitId] = element
				}
			}
		}
		PreHandlePolicyDataChangeNotification(ueId, "", smPolicyData)
		return nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "MODIFY_NOT_ALLOWED",
			Status: http.StatusForbidden,
		}
		return problemDetails
	}
}

func HandlePolicyDataUesUeIdSmDataUsageMonIdDelete(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle PolicyDataUesUeIdSmDataUsageMonIdDelete")

	collName := "policyData.ues.smData.usageMonData"
	ueId := request.Params["ueId"]
	usageMonId := request.Params["usageMonId"]

	PolicyDataUesUeIdSmDataUsageMonIdDeleteProcedure(collName, ueId, usageMonId)
	return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
}

func PolicyDataUesUeIdSmDataUsageMonIdDeleteProcedure(collName string, ueId string, usageMonId string) {
	filter := bson.M{"ueId": ueId, "usageMonId": usageMonId}
	MongoDBLibrary.RestfulAPIDeleteOne(collName, filter)
}

func HandlePolicyDataUesUeIdSmDataUsageMonIdGet(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle PolicyDataUesUeIdSmDataUsageMonIdGet")

	collName := "policyData.ues.smData.usageMonData"
	ueId := request.Params["ueId"]
	usageMonId := request.Params["usageMonId"]

	response := PolicyDataUesUeIdSmDataUsageMonIdGetProcedure(collName, usageMonId, ueId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	}
}

func PolicyDataUesUeIdSmDataUsageMonIdGetProcedure(collName string, usageMonId string,
	ueId string) *map[string]interface{} {
	filter := bson.M{"ueId": ueId, "usageMonId": usageMonId}

	usageMonData := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	return &usageMonData
}

func HandlePolicyDataUesUeIdSmDataUsageMonIdPut(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle PolicyDataUesUeIdSmDataUsageMonIdPut")

	ueId := request.Params["ueId"]
	usageMonId := request.Params["usageMonId"]
	usageMonData := request.Body.(models.UsageMonData)
	collName := "policyData.ues.smData.usageMonData"

	response := PolicyDataUesUeIdSmDataUsageMonIdPutProcedure(collName, ueId, usageMonId, usageMonData)

	return http_wrapper.NewResponse(http.StatusCreated, nil, response)
}

func PolicyDataUesUeIdSmDataUsageMonIdPutProcedure(collName string, ueId string, usageMonId string,
	usageMonData models.UsageMonData) *bson.M {
	putData := util.ToBsonM(usageMonData)
	putData["ueId"] = ueId
	putData["usageMonId"] = usageMonId
	filter := bson.M{"ueId": ueId, "usageMonId": usageMonId}

	MongoDBLibrary.RestfulAPIPutOne(collName, filter, putData)
	return &putData
}

func HandlePolicyDataUesUeIdUePolicySetGet(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle PolicyDataUesUeIdUePolicySetGet")

	ueId := request.Params["ueId"]
	collName := "policyData.ues.uePolicySet"

	response, problemDetails := PolicyDataUesUeIdUePolicySetGetProcedure(collName, ueId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func PolicyDataUesUeIdUePolicySetGetProcedure(collName string, ueId string) (*map[string]interface{},
	*models.ProblemDetails) {
	filter := bson.M{"ueId": ueId}

	uePolicySet := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	if uePolicySet != nil {
		return &uePolicySet, nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}
}

func HandlePolicyDataUesUeIdUePolicySetPatch(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle PolicyDataUesUeIdUePolicySetPatch")

	collName := "policyData.ues.uePolicySet"
	ueId := request.Params["ueId"]
	UePolicySet := request.Body.(models.UePolicySet)

	problemDetails := PolicyDataUesUeIdUePolicySetPatchProcedure(collName, ueId, UePolicySet)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
}

func PolicyDataUesUeIdUePolicySetPatchProcedure(collName string, ueId string,
	UePolicySet models.UePolicySet) *models.ProblemDetails {
	patchData := util.ToBsonM(UePolicySet)
	patchData["ueId"] = ueId
	filter := bson.M{"ueId": ueId}

	success := MongoDBLibrary.RestfulAPIMergePatch(collName, filter, patchData)

	if success {
		var uePolicySet models.UePolicySet
		uePolicySetBsonM := MongoDBLibrary.RestfulAPIGetOne(collName, filter)
		err := json.Unmarshal(util.MapToByte(uePolicySetBsonM), &uePolicySet)
		if err != nil {
			logger.DataRepoLog.Warnln(err)
		}
		PreHandlePolicyDataChangeNotification(ueId, "", uePolicySet)
		return nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "MODIFY_NOT_ALLOWED",
			Status: http.StatusForbidden,
		}
		return problemDetails
	}
}

func HandlePolicyDataUesUeIdUePolicySetPut(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle PolicyDataUesUeIdUePolicySetPut")

	collName := "policyData.ues.uePolicySet"
	ueId := request.Params["ueId"]
	UePolicySet := request.Body.(models.UePolicySet)

	response, status := PolicyDataUesUeIdUePolicySetPutProcedure(collName, ueId, UePolicySet)

	if status == http.StatusNoContent {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else if status == http.StatusCreated {
		return http_wrapper.NewResponse(http.StatusCreated, nil, response)
	}
	problemDetails := &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func PolicyDataUesUeIdUePolicySetPutProcedure(collName string, ueId string,
	UePolicySet models.UePolicySet) (bson.M, int) {
	putData := util.ToBsonM(UePolicySet)
	putData["ueId"] = ueId
	filter := bson.M{"ueId": ueId}

	isExisted := MongoDBLibrary.RestfulAPIPutOne(collName, filter, putData)
	if !isExisted {
		return putData, http.StatusCreated
	} else {
		return nil, http.StatusNoContent
	}
}

func HandleCreateAMFSubscriptions(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle CreateAMFSubscriptions")

	ueId := request.Params["ueId"]
	subsId := request.Params["subsId"]
	AmfSubscriptionInfo := request.Body.([]models.AmfSubscriptionInfo)

	problemDetails := CreateAMFSubscriptionsProcedure(subsId, ueId, AmfSubscriptionInfo)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
}

func CreateAMFSubscriptionsProcedure(subsId string, ueId string,
	AmfSubscriptionInfo []models.AmfSubscriptionInfo) *models.ProblemDetails {
	nwdafSelf := nwdaf_context.NWDAF_Self()
	value, ok := nwdafSelf.UESubsCollection.Load(ueId)
	if !ok {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return problemDetails
	}
	UESubsData := value.(*nwdaf_context.UESubsData)

	_, ok = UESubsData.EeSubscriptionCollection[subsId]
	if !ok {
		problemDetails := &models.ProblemDetails{
			Cause:  "SUBSCRIPTION_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return problemDetails
	}

	UESubsData.EeSubscriptionCollection[subsId].AmfSubscriptionInfos = AmfSubscriptionInfo
	return nil
}

func HandleRemoveAmfSubscriptionsInfo(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle RemoveAmfSubscriptionsInfo")

	ueId := request.Params["ueId"]
	subsId := request.Params["subsId"]

	problemDetails := RemoveAmfSubscriptionsInfoProcedure(subsId, ueId)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
}

func RemoveAmfSubscriptionsInfoProcedure(subsId string, ueId string) *models.ProblemDetails {
	nwdafSelf := nwdaf_context.NWDAF_Self()
	value, ok := nwdafSelf.UESubsCollection.Load(ueId)
	if !ok {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return problemDetails
	}

	UESubsData := value.(*nwdaf_context.UESubsData)
	_, ok = UESubsData.EeSubscriptionCollection[subsId]

	if !ok {
		problemDetails := &models.ProblemDetails{
			Cause:  "SUBSCRIPTION_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return problemDetails
	}

	if UESubsData.EeSubscriptionCollection[subsId].AmfSubscriptionInfos == nil {
		problemDetails := &models.ProblemDetails{
			Cause:  "AMFSUBSCRIPTION_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return problemDetails
	}

	UESubsData.EeSubscriptionCollection[subsId].AmfSubscriptionInfos = nil

	return nil
}

func HandleModifyAmfSubscriptionInfo(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle ModifyAmfSubscriptionInfo")

	patchItem := request.Body.([]models.PatchItem)
	ueId := request.Params["ueId"]
	subsId := request.Params["subsId"]

	problemDetails := ModifyAmfSubscriptionInfoProcedure(ueId, subsId, patchItem)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
}

func ModifyAmfSubscriptionInfoProcedure(ueId string, subsId string,
	patchItem []models.PatchItem) *models.ProblemDetails {
	nwdafSelf := nwdaf_context.NWDAF_Self()
	value, ok := nwdafSelf.UESubsCollection.Load(ueId)
	if !ok {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return problemDetails
	}
	UESubsData := value.(*nwdaf_context.UESubsData)

	_, ok = UESubsData.EeSubscriptionCollection[subsId]

	if !ok {
		problemDetails := &models.ProblemDetails{
			Cause:  "SUBSCRIPTION_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return problemDetails
	}

	if UESubsData.EeSubscriptionCollection[subsId].AmfSubscriptionInfos == nil {
		problemDetails := &models.ProblemDetails{
			Cause:  "AMFSUBSCRIPTION_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return problemDetails
	}
	var patchJSON []byte
	if patchJSONtemp, err := json.Marshal(patchItem); err != nil {
		logger.DataRepoLog.Errorln(err)
	} else {
		patchJSON = patchJSONtemp
	}
	var patch jsonpatch.Patch
	if patchtemp, err := jsonpatch.DecodePatch(patchJSON); err != nil {
		logger.DataRepoLog.Errorln(err)
		problemDetails := &models.ProblemDetails{
			Cause:  "MODIFY_NOT_ALLOWED",
			Detail: "PatchItem attributes are invalid",
			Status: http.StatusForbidden,
		}
		return problemDetails
	} else {
		patch = patchtemp
	}
	original, err := json.Marshal((UESubsData.EeSubscriptionCollection[subsId]).AmfSubscriptionInfos)
	if err != nil {
		logger.DataRepoLog.Warnln(err)
	}

	modified, err := patch.Apply(original)
	if err != nil {
		problemDetails := &models.ProblemDetails{
			Cause:  "MODIFY_NOT_ALLOWED",
			Detail: "Occur error when applying PatchItem",
			Status: http.StatusForbidden,
		}
		return problemDetails
	}
	var modifiedData []models.AmfSubscriptionInfo
	err = json.Unmarshal(modified, &modifiedData)
	if err != nil {
		logger.DataRepoLog.Error(err)
	}

	UESubsData.EeSubscriptionCollection[subsId].AmfSubscriptionInfos = modifiedData
	return nil
}

func HandleGetAmfSubscriptionInfo(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle GetAmfSubscriptionInfo")

	ueId := request.Params["ueId"]
	subsId := request.Params["subsId"]

	response, problemDetails := GetAmfSubscriptionInfoProcedure(subsId, ueId)
	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func GetAmfSubscriptionInfoProcedure(subsId string, ueId string) (*[]models.AmfSubscriptionInfo,
	*models.ProblemDetails) {
	nwdafSelf := nwdaf_context.NWDAF_Self()

	value, ok := nwdafSelf.UESubsCollection.Load(ueId)
	if !ok {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}

	UESubsData := value.(*nwdaf_context.UESubsData)
	_, ok = UESubsData.EeSubscriptionCollection[subsId]

	if !ok {
		problemDetails := &models.ProblemDetails{
			Cause:  "SUBSCRIPTION_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}

	if UESubsData.EeSubscriptionCollection[subsId].AmfSubscriptionInfos == nil {
		problemDetails := &models.ProblemDetails{
			Cause:  "AMFSUBSCRIPTION_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}
	return &UESubsData.EeSubscriptionCollection[subsId].AmfSubscriptionInfos, nil
}

func HandleQueryEEData(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QueryEEData")

	ueId := request.Params["ueId"]
	collName := "subscriptionData.eeProfileData"

	response, problemDetails := QueryEEDataProcedure(collName, ueId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func QueryEEDataProcedure(collName string, ueId string) (*map[string]interface{}, *models.ProblemDetails) {
	filter := bson.M{"ueId": ueId}
	eeProfileData := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	if eeProfileData != nil {
		return &eeProfileData, nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}
}

func HandleRemoveEeGroupSubscriptions(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle RemoveEeGroupSubscriptions")

	ueGroupId := request.Params["ueGroupId"]
	subsId := request.Params["subsId"]

	problemDetails := RemoveEeGroupSubscriptionsProcedure(ueGroupId, subsId)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
}

func RemoveEeGroupSubscriptionsProcedure(ueGroupId string, subsId string) *models.ProblemDetails {
	nwdafSelf := nwdaf_context.NWDAF_Self()
	value, ok := nwdafSelf.UEGroupCollection.Load(ueGroupId)
	if !ok {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return problemDetails
	}

	UEGroupSubsData := value.(*nwdaf_context.UEGroupSubsData)
	_, ok = UEGroupSubsData.EeSubscriptions[subsId]

	if !ok {
		problemDetails := &models.ProblemDetails{
			Cause:  "SUBSCRIPTION_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return problemDetails
	}
	delete(UEGroupSubsData.EeSubscriptions, subsId)

	return nil
}

func HandleUpdateEeGroupSubscriptions(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle UpdateEeGroupSubscriptions")

	ueGroupId := request.Params["ueGroupId"]
	subsId := request.Params["subsId"]
	EeSubscription := request.Body.(models.EeSubscription)

	problemDetails := UpdateEeGroupSubscriptionsProcedure(ueGroupId, subsId, EeSubscription)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
}

func UpdateEeGroupSubscriptionsProcedure(ueGroupId string, subsId string,
	EeSubscription models.EeSubscription) *models.ProblemDetails {
	nwdafSelf := nwdaf_context.NWDAF_Self()
	value, ok := nwdafSelf.UEGroupCollection.Load(ueGroupId)
	if !ok {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return problemDetails
	}

	UEGroupSubsData := value.(*nwdaf_context.UEGroupSubsData)
	_, ok = UEGroupSubsData.EeSubscriptions[subsId]

	if !ok {
		problemDetails := &models.ProblemDetails{
			Cause:  "SUBSCRIPTION_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return problemDetails
	}
	UEGroupSubsData.EeSubscriptions[subsId] = &EeSubscription

	return nil
}

func HandleCreateEeGroupSubscriptions(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle CreateEeGroupSubscriptions")

	ueGroupId := request.Params["ueGroupId"]
	EeSubscription := request.Body.(models.EeSubscription)

	locationHeader := CreateEeGroupSubscriptionsProcedure(ueGroupId, EeSubscription)

	headers := http.Header{
		"Location": {locationHeader},
	}
	return http_wrapper.NewResponse(http.StatusCreated, headers, EeSubscription)
}

func CreateEeGroupSubscriptionsProcedure(ueGroupId string, EeSubscription models.EeSubscription) string {
	nwdafSelf := nwdaf_context.NWDAF_Self()

	value, ok := nwdafSelf.UEGroupCollection.Load(ueGroupId)
	if !ok {
		nwdafSelf.UEGroupCollection.Store(ueGroupId, new(nwdaf_context.UEGroupSubsData))
		value, _ = nwdafSelf.UEGroupCollection.Load(ueGroupId)
	}
	UEGroupSubsData := value.(*nwdaf_context.UEGroupSubsData)
	if UEGroupSubsData.EeSubscriptions == nil {
		UEGroupSubsData.EeSubscriptions = make(map[string]*models.EeSubscription)
	}

	newSubscriptionID := strconv.Itoa(nwdafSelf.EeSubscriptionIDGenerator)
	UEGroupSubsData.EeSubscriptions[newSubscriptionID] = &EeSubscription
	nwdafSelf.EeSubscriptionIDGenerator++

	/* Contains the URI of the newly created resource, according
	   to the structure: {apiRoot}/nnwdaf-dr/v1/subscription-data/group-data/{ueGroupId}/ee-subscriptions */
	locationHeader := fmt.Sprintf("%s/nnwdaf-dr/v1/subscription-data/group-data/%s/ee-subscriptions/%s",
		nwdafSelf.GetIPv4GroupUri(nwdaf_context.NNWDAF_DR), ueGroupId, newSubscriptionID)

	return locationHeader
}

func HandleQueryEeGroupSubscriptions(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QueryEeGroupSubscriptions")

	ueGroupId := request.Params["ueGroupId"]

	response, problemDetails := QueryEeGroupSubscriptionsProcedure(ueGroupId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func QueryEeGroupSubscriptionsProcedure(ueGroupId string) ([]models.EeSubscription, *models.ProblemDetails) {
	nwdafSelf := nwdaf_context.NWDAF_Self()

	value, ok := nwdafSelf.UEGroupCollection.Load(ueGroupId)
	if !ok {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}

	UEGroupSubsData := value.(*nwdaf_context.UEGroupSubsData)
	var eeSubscriptionSlice []models.EeSubscription

	for _, v := range UEGroupSubsData.EeSubscriptions {
		eeSubscriptionSlice = append(eeSubscriptionSlice, *v)
	}
	return eeSubscriptionSlice, nil
}

func HandleRemoveeeSubscriptions(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle RemoveeeSubscriptions")

	ueId := request.Params["ueId"]
	subsId := request.Params["subsId"]

	problemDetails := RemoveeeSubscriptionsProcedure(ueId, subsId)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
}

func RemoveeeSubscriptionsProcedure(ueId string, subsId string) *models.ProblemDetails {
	nwdafSelf := nwdaf_context.NWDAF_Self()
	value, ok := nwdafSelf.UESubsCollection.Load(ueId)
	if !ok {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return problemDetails
	}

	UESubsData := value.(*nwdaf_context.UESubsData)
	_, ok = UESubsData.EeSubscriptionCollection[subsId]

	if !ok {
		problemDetails := &models.ProblemDetails{
			Cause:  "SUBSCRIPTION_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return problemDetails
	}
	delete(UESubsData.EeSubscriptionCollection, subsId)
	return nil
}

func HandleUpdateEesubscriptions(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle UpdateEesubscriptions")

	ueId := request.Params["ueId"]
	subsId := request.Params["subsId"]
	EeSubscription := request.Body.(models.EeSubscription)

	problemDetails := UpdateEesubscriptionsProcedure(ueId, subsId, EeSubscription)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
}

func UpdateEesubscriptionsProcedure(ueId string, subsId string,
	EeSubscription models.EeSubscription) *models.ProblemDetails {
	nwdafSelf := nwdaf_context.NWDAF_Self()
	value, ok := nwdafSelf.UESubsCollection.Load(ueId)
	if !ok {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return problemDetails
	}

	UESubsData := value.(*nwdaf_context.UESubsData)
	_, ok = UESubsData.EeSubscriptionCollection[subsId]

	if !ok {
		problemDetails := &models.ProblemDetails{
			Cause:  "SUBSCRIPTION_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return problemDetails
	}
	UESubsData.EeSubscriptionCollection[subsId].EeSubscriptions = &EeSubscription

	return nil
}

func HandleCreateEeSubscriptions(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle CreateEeSubscriptions")

	ueId := request.Params["ueId"]
	EeSubscription := request.Body.(models.EeSubscription)

	locationHeader := CreateEeSubscriptionsProcedure(ueId, EeSubscription)

	headers := http.Header{
		"Location": {locationHeader},
	}
	return http_wrapper.NewResponse(http.StatusCreated, headers, EeSubscription)
}

func CreateEeSubscriptionsProcedure(ueId string, EeSubscription models.EeSubscription) string {
	nwdafSelf := nwdaf_context.NWDAF_Self()

	value, ok := nwdafSelf.UESubsCollection.Load(ueId)
	if !ok {
		nwdafSelf.UESubsCollection.Store(ueId, new(nwdaf_context.UESubsData))
		value, _ = nwdafSelf.UESubsCollection.Load(ueId)
	}
	UESubsData := value.(*nwdaf_context.UESubsData)
	if UESubsData.EeSubscriptionCollection == nil {
		UESubsData.EeSubscriptionCollection = make(map[string]*nwdaf_context.EeSubscriptionCollection)
	}

	newSubscriptionID := strconv.Itoa(nwdafSelf.EeSubscriptionIDGenerator)
	UESubsData.EeSubscriptionCollection[newSubscriptionID] = new(nwdaf_context.EeSubscriptionCollection)
	UESubsData.EeSubscriptionCollection[newSubscriptionID].EeSubscriptions = &EeSubscription
	nwdafSelf.EeSubscriptionIDGenerator++

	/* Contains the URI of the newly created resource, according
	   to the structure: {apiRoot}/subscription-data/{ueId}/context-data/ee-subscriptions/{subsId} */
	locationHeader := fmt.Sprintf("%s/subscription-data/%s/context-data/ee-subscriptions/%s",
		nwdafSelf.GetIPv4GroupUri(nwdaf_context.NNWDAF_DR), ueId, newSubscriptionID)

	return locationHeader
}

func HandleQueryeesubscriptions(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle Queryeesubscriptions")

	ueId := request.Params["ueId"]

	response, problemDetails := QueryeesubscriptionsProcedure(ueId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func QueryeesubscriptionsProcedure(ueId string) ([]models.EeSubscription, *models.ProblemDetails) {
	nwdafSelf := nwdaf_context.NWDAF_Self()

	value, ok := nwdafSelf.UESubsCollection.Load(ueId)
	if !ok {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}

	UESubsData := value.(*nwdaf_context.UESubsData)
	var eeSubscriptionSlice []models.EeSubscription

	for _, v := range UESubsData.EeSubscriptionCollection {
		eeSubscriptionSlice = append(eeSubscriptionSlice, *v.EeSubscriptions)
	}
	return eeSubscriptionSlice, nil
}

func HandlePatchOperSpecData(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle PatchOperSpecData")

	collName := "subscriptionData.operatorSpecificData"
	ueId := request.Params["ueId"]
	patchItem := request.Body.([]models.PatchItem)

	problemDetails := PatchOperSpecDataProcedure(collName, ueId, patchItem)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
}

func PatchOperSpecDataProcedure(collName string, ueId string, patchItem []models.PatchItem) *models.ProblemDetails {
	filter := bson.M{"ueId": ueId}

	origValue := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	patchJSON, err := json.Marshal(patchItem)
	if err != nil {
		logger.DataRepoLog.Errorln(err)
	}

	success := MongoDBLibrary.RestfulAPIJSONPatch(collName, filter, patchJSON)

	if success {
		newValue := MongoDBLibrary.RestfulAPIGetOne(collName, filter)
		PreHandleOnDataChangeNotify(ueId, CurrentResourceUri, patchItem, origValue, newValue)
		return nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "MODIFY_NOT_ALLOWED",
			Status: http.StatusForbidden,
		}
		return problemDetails
	}
}

func HandleQueryOperSpecData(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QueryOperSpecData")

	ueId := request.Params["ueId"]
	collName := "subscriptionData.operatorSpecificData"

	response, problemDetails := QueryOperSpecDataProcedure(collName, ueId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func QueryOperSpecDataProcedure(collName string, ueId string) (*map[string]interface{}, *models.ProblemDetails) {
	filter := bson.M{"ueId": ueId}

	operatorSpecificDataContainer := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	// The key of the map is operator specific data element name and the value is the operator specific data of the UE.

	if operatorSpecificDataContainer != nil {
		return &operatorSpecificDataContainer, nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}
}

func HandleGetppData(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle GetppData")

	collName := "subscriptionData.ppData"
	ueId := request.Params["ueId"]

	response, problemDetails := GetppDataProcedure(collName, ueId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func GetppDataProcedure(collName string, ueId string) (*map[string]interface{}, *models.ProblemDetails) {
	filter := bson.M{"ueId": ueId}

	ppData := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	if ppData != nil {
		return &ppData, nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}
}

func HandleCreateSessionManagementData(request *http_wrapper.Request) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusOK, nil, map[string]interface{}{})
}

func HandleDeleteSessionManagementData(request *http_wrapper.Request) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusOK, nil, map[string]interface{}{})
}

func HandleQuerySessionManagementData(request *http_wrapper.Request) *http_wrapper.Response {
	return http_wrapper.NewResponse(http.StatusOK, nil, map[string]interface{}{})
}

func HandleQueryProvisionedData(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QueryProvisionedData")

	var provisionedDataSets models.ProvisionedDataSets
	ueId := request.Params["ueId"]
	servingPlmnId := request.Params["servingPlmnId"]

	response, problemDetails := QueryProvisionedDataProcedure(ueId, servingPlmnId, provisionedDataSets)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func QueryProvisionedDataProcedure(ueId string, servingPlmnId string,
	provisionedDataSets models.ProvisionedDataSets) (*models.ProvisionedDataSets, *models.ProblemDetails) {
	{
		collName := "subscriptionData.provisionedData.amData"
		filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
		accessAndMobilitySubscriptionData := MongoDBLibrary.RestfulAPIGetOne(collName, filter)
		if accessAndMobilitySubscriptionData != nil {
			var tmp models.AccessAndMobilitySubscriptionData
			err := mapstructure.Decode(accessAndMobilitySubscriptionData, &tmp)
			if err != nil {
				panic(err)
			}
			provisionedDataSets.AmData = &tmp
		}
	}

	{
		collName := "subscriptionData.provisionedData.smfSelectionSubscriptionData"
		filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
		smfSelectionSubscriptionData := MongoDBLibrary.RestfulAPIGetOne(collName, filter)
		if smfSelectionSubscriptionData != nil {
			var tmp models.SmfSelectionSubscriptionData
			err := mapstructure.Decode(smfSelectionSubscriptionData, &tmp)
			if err != nil {
				panic(err)
			}
			provisionedDataSets.SmfSelData = &tmp
		}
	}

	{
		collName := "subscriptionData.provisionedData.smsData"
		filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
		smsSubscriptionData := MongoDBLibrary.RestfulAPIGetOne(collName, filter)
		if smsSubscriptionData != nil {
			var tmp models.SmsSubscriptionData
			err := mapstructure.Decode(smsSubscriptionData, &tmp)
			if err != nil {
				panic(err)
			}
			provisionedDataSets.SmsSubsData = &tmp
		}
	}

	{
		collName := "subscriptionData.provisionedData.smData"
		filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
		sessionManagementSubscriptionDatas := MongoDBLibrary.RestfulAPIGetMany(collName, filter)
		if sessionManagementSubscriptionDatas != nil {
			var tmp []models.SessionManagementSubscriptionData
			err := mapstructure.Decode(sessionManagementSubscriptionDatas, &tmp)
			if err != nil {
				panic(err)
			}
			provisionedDataSets.SmData = tmp
		}
	}

	{
		collName := "subscriptionData.provisionedData.traceData"
		filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
		traceData := MongoDBLibrary.RestfulAPIGetOne(collName, filter)
		if traceData != nil {
			var tmp models.TraceData
			err := mapstructure.Decode(traceData, &tmp)
			if err != nil {
				panic(err)
			}
			provisionedDataSets.TraceData = &tmp
		}
	}

	{
		collName := "subscriptionData.provisionedData.smsMngData"
		filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
		smsManagementSubscriptionData := MongoDBLibrary.RestfulAPIGetOne(collName, filter)
		if smsManagementSubscriptionData != nil {
			var tmp models.SmsManagementSubscriptionData
			err := mapstructure.Decode(smsManagementSubscriptionData, &tmp)
			if err != nil {
				panic(err)
			}
			provisionedDataSets.SmsMngData = &tmp
		}
	}

	if !reflect.DeepEqual(provisionedDataSets, models.ProvisionedDataSets{}) {
		return &provisionedDataSets, nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}
}

func HandleModifyPpData(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle ModifyPpData")

	collName := "subscriptionData.ppData"
	patchItem := request.Body.([]models.PatchItem)
	ueId := request.Params["ueId"]

	problemDetails := ModifyPpDataProcedure(collName, ueId, patchItem)
	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
}

func ModifyPpDataProcedure(collName string, ueId string, patchItem []models.PatchItem) *models.ProblemDetails {
	filter := bson.M{"ueId": ueId}

	origValue := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	patchJSON, err := json.Marshal(patchItem)
	if err != nil {
		logger.DataRepoLog.Errorln(err)
	}

	success := MongoDBLibrary.RestfulAPIJSONPatch(collName, filter, patchJSON)

	if success {
		newValue := MongoDBLibrary.RestfulAPIGetOne(collName, filter)
		PreHandleOnDataChangeNotify(ueId, CurrentResourceUri, patchItem, origValue, newValue)
		return nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "MODIFY_NOT_ALLOWED",
			Status: http.StatusForbidden,
		}
		return problemDetails
	}
}

func HandleGetIdentityData(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle GetIdentityData")

	ueId := request.Params["ueId"]
	collName := "subscriptionData.identityData"

	response, problemDetails := GetIdentityDataProcedure(collName, ueId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func GetIdentityDataProcedure(collName string, ueId string) (*map[string]interface{}, *models.ProblemDetails) {
	filter := bson.M{"ueId": ueId}

	identityData := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	if identityData != nil {
		return &identityData, nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}
}

func HandleGetOdbData(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle GetOdbData")

	ueId := request.Params["ueId"]
	collName := "subscriptionData.operatorDeterminedBarringData"

	response, problemDetails := GetOdbDataProcedure(collName, ueId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func GetOdbDataProcedure(collName string, ueId string) (*map[string]interface{}, *models.ProblemDetails) {
	filter := bson.M{"ueId": ueId}

	operatorDeterminedBarringData := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	if operatorDeterminedBarringData != nil {
		return &operatorDeterminedBarringData, nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}
}

func HandleGetSharedData(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle GetSharedData")

	var sharedDataIds []string
	if len(request.Query["shared-data-ids"]) != 0 {
		sharedDataIds = request.Query["shared-data-ids"]
		if strings.Contains(sharedDataIds[0], ",") {
			sharedDataIds = strings.Split(sharedDataIds[0], ",")
		}
	}
	collName := "subscriptionData.sharedData"

	response, problemDetails := GetSharedDataProcedure(collName, sharedDataIds)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func GetSharedDataProcedure(collName string, sharedDataIds []string) (*[]map[string]interface{},
	*models.ProblemDetails) {
	var sharedDataArray []map[string]interface{}
	for _, sharedDataId := range sharedDataIds {
		filter := bson.M{"sharedDataId": sharedDataId}
		sharedData := MongoDBLibrary.RestfulAPIGetOne(collName, filter)
		if sharedData != nil {
			sharedDataArray = append(sharedDataArray, sharedData)
		}
	}

	if sharedDataArray != nil {
		return &sharedDataArray, nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "DATA_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}

}

func HandleRemovesdmSubscriptions(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle RemovesdmSubscriptions")

	ueId := request.Params["ueId"]
	subsId := request.Params["subsId"]

	problemDetails := RemovesdmSubscriptionsProcedure(ueId, subsId)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
}

func RemovesdmSubscriptionsProcedure(ueId string, subsId string) *models.ProblemDetails {
	nwdafSelf := nwdaf_context.NWDAF_Self()
	value, ok := nwdafSelf.UESubsCollection.Load(ueId)
	if !ok {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return problemDetails
	}

	UESubsData := value.(*nwdaf_context.UESubsData)
	_, ok = UESubsData.SdmSubscriptions[subsId]

	if !ok {
		problemDetails := &models.ProblemDetails{
			Cause:  "SUBSCRIPTION_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return problemDetails
	}
	delete(UESubsData.SdmSubscriptions, subsId)

	return nil
}

func HandleUpdatesdmsubscriptions(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle Updatesdmsubscriptions")

	ueId := request.Params["ueId"]
	subsId := request.Params["subsId"]
	SdmSubscription := request.Body.(models.SdmSubscription)

	problemDetails := UpdatesdmsubscriptionsProcedure(ueId, subsId, SdmSubscription)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
}

func UpdatesdmsubscriptionsProcedure(ueId string, subsId string,
	SdmSubscription models.SdmSubscription) *models.ProblemDetails {
	nwdafSelf := nwdaf_context.NWDAF_Self()
	value, ok := nwdafSelf.UESubsCollection.Load(ueId)
	if !ok {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return problemDetails
	}

	UESubsData := value.(*nwdaf_context.UESubsData)
	_, ok = UESubsData.SdmSubscriptions[subsId]

	if !ok {
		problemDetails := &models.ProblemDetails{
			Cause:  "SUBSCRIPTION_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return problemDetails
	}
	SdmSubscription.SubscriptionId = subsId
	UESubsData.SdmSubscriptions[subsId] = &SdmSubscription

	return nil
}

func HandleCreateSdmSubscriptions(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle CreateSdmSubscriptions")

	SdmSubscription := request.Body.(models.SdmSubscription)
	collName := "subscriptionData.contextData.amfNon3gppAccess"
	ueId := request.Params["ueId"]

	locationHeader, SdmSubscription := CreateSdmSubscriptionsProcedure(SdmSubscription, collName, ueId)
	headers := http.Header{
		"Location": {locationHeader},
	}

	return http_wrapper.NewResponse(http.StatusCreated, headers, SdmSubscription)
}

func CreateSdmSubscriptionsProcedure(SdmSubscription models.SdmSubscription,
	collName string, ueId string) (string, models.SdmSubscription) {
	nwdafSelf := nwdaf_context.NWDAF_Self()

	value, ok := nwdafSelf.UESubsCollection.Load(ueId)
	if !ok {
		nwdafSelf.UESubsCollection.Store(ueId, new(nwdaf_context.UESubsData))
		value, _ = nwdafSelf.UESubsCollection.Load(ueId)
	}
	UESubsData := value.(*nwdaf_context.UESubsData)
	if UESubsData.SdmSubscriptions == nil {
		UESubsData.SdmSubscriptions = make(map[string]*models.SdmSubscription)
	}

	newSubscriptionID := strconv.Itoa(nwdafSelf.SdmSubscriptionIDGenerator)
	SdmSubscription.SubscriptionId = newSubscriptionID
	UESubsData.SdmSubscriptions[newSubscriptionID] = &SdmSubscription
	nwdafSelf.SdmSubscriptionIDGenerator++

	/* Contains the URI of the newly created resource, according
	   to the structure: {apiRoot}/subscription-data/{ueId}/context-data/sdm-subscriptions/{subsId}' */
	locationHeader := fmt.Sprintf("%s/subscription-data/%s/context-data/sdm-subscriptions/%s",
		nwdafSelf.GetIPv4GroupUri(nwdaf_context.NNWDAF_DR), ueId, newSubscriptionID)

	return locationHeader, SdmSubscription
}

func HandleQuerysdmsubscriptions(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle Querysdmsubscriptions")

	ueId := request.Params["ueId"]

	response, problemDetails := QuerysdmsubscriptionsProcedure(ueId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func QuerysdmsubscriptionsProcedure(ueId string) (*[]models.SdmSubscription, *models.ProblemDetails) {
	nwdafSelf := nwdaf_context.NWDAF_Self()

	value, ok := nwdafSelf.UESubsCollection.Load(ueId)
	if !ok {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}

	UESubsData := value.(*nwdaf_context.UESubsData)
	var sdmSubscriptionSlice []models.SdmSubscription

	for _, v := range UESubsData.SdmSubscriptions {
		sdmSubscriptionSlice = append(sdmSubscriptionSlice, *v)
	}
	return &sdmSubscriptionSlice, nil
}

func HandleQuerySmData(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QuerySmData")

	collName := "subscriptionData.provisionedData.smData"
	ueId := request.Params["ueId"]
	servingPlmnId := request.Params["servingPlmnId"]
	singleNssai := models.Snssai{}
	singleNssaiQuery := request.Query.Get("single-nssai")
	err := json.Unmarshal([]byte(singleNssaiQuery), &singleNssai)
	if err != nil {
		logger.DataRepoLog.Warnln(err)
	}

	dnn := request.Query.Get("dnn")
	response := QuerySmDataProcedure(collName, ueId, servingPlmnId, singleNssai, dnn)

	return http_wrapper.NewResponse(http.StatusOK, nil, response)
}

func QuerySmDataProcedure(collName string, ueId string, servingPlmnId string,
	singleNssai models.Snssai, dnn string) *[]map[string]interface{} {
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}

	if !reflect.DeepEqual(singleNssai, models.Snssai{}) {
		if singleNssai.Sd == "" {
			filter["singleNssai.sst"] = singleNssai.Sst
		} else {
			filter["singleNssai.sst"] = singleNssai.Sst
			filter["singleNssai.sd"] = singleNssai.Sd
		}
	}

	if dnn != "" {
		filter["dnnConfigurations."+dnn] = bson.M{"$exists": true}
	}

	sessionManagementSubscriptionDatas := MongoDBLibrary.RestfulAPIGetMany(collName, filter)

	return &sessionManagementSubscriptionDatas
}

func HandleCreateSmfContextNon3gpp(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle CreateSmfContextNon3gpp")

	SmfRegistration := request.Body.(models.SmfRegistration)
	collName := "subscriptionData.contextData.smfRegistrations"
	ueId := request.Params["ueId"]
	pduSessionId, err := strconv.ParseInt(request.Params["pduSessionId"], 10, 64)
	if err != nil {
		logger.DataRepoLog.Warnln(err)
	}

	response, status := CreateSmfContextNon3gppProcedure(SmfRegistration, collName, ueId, pduSessionId)

	if status == http.StatusCreated {
		return http_wrapper.NewResponse(http.StatusCreated, nil, response)
	} else if status == http.StatusOK {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	}
	problemDetails := &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func CreateSmfContextNon3gppProcedure(SmfRegistration models.SmfRegistration,
	collName string, ueId string, pduSessionIdInt int64) (bson.M, int) {
	putData := util.ToBsonM(SmfRegistration)
	putData["ueId"] = ueId
	putData["pduSessionId"] = int32(pduSessionIdInt)

	filter := bson.M{"ueId": ueId, "pduSessionId": pduSessionIdInt}
	isExisted := MongoDBLibrary.RestfulAPIPutOne(collName, filter, putData)

	if !isExisted {
		return putData, http.StatusCreated
	} else {
		return putData, http.StatusOK
	}
}

func HandleDeleteSmfContext(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle DeleteSmfContext")

	collName := "subscriptionData.contextData.smfRegistrations"
	ueId := request.Params["ueId"]
	pduSessionId := request.Params["pduSessionId"]

	DeleteSmfContextProcedure(collName, ueId, pduSessionId)
	return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
}

func DeleteSmfContextProcedure(collName string, ueId string, pduSessionId string) {
	pduSessionIdInt, err := strconv.ParseInt(pduSessionId, 10, 32)
	if err != nil {
		logger.DataRepoLog.Error(err)
	}
	filter := bson.M{"ueId": ueId, "pduSessionId": pduSessionIdInt}

	MongoDBLibrary.RestfulAPIDeleteOne(collName, filter)
}

func HandleQuerySmfRegistration(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QuerySmfRegistration")

	ueId := request.Params["ueId"]
	pduSessionId := request.Params["pduSessionId"]
	collName := "subscriptionData.contextData.smfRegistrations"

	response, problemDetails := QuerySmfRegistrationProcedure(collName, ueId, pduSessionId)
	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func QuerySmfRegistrationProcedure(collName string, ueId string,
	pduSessionId string) (*map[string]interface{}, *models.ProblemDetails) {
	pduSessionIdInt, err := strconv.ParseInt(pduSessionId, 10, 32)
	if err != nil {
		logger.DataRepoLog.Error(err)
	}

	filter := bson.M{"ueId": ueId, "pduSessionId": pduSessionIdInt}

	smfRegistration := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	if smfRegistration != nil {
		return &smfRegistration, nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}
}

func HandleQuerySmfRegList(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QuerySmfRegList")

	collName := "subscriptionData.contextData.smfRegistrations"
	ueId := request.Params["ueId"]
	response := QuerySmfRegListProcedure(collName, ueId)

	if response == nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, []map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	}
}

func QuerySmfRegListProcedure(collName string, ueId string) *[]map[string]interface{} {
	filter := bson.M{"ueId": ueId}
	smfRegList := MongoDBLibrary.RestfulAPIGetMany(collName, filter)

	if smfRegList != nil {
		return &smfRegList
	} else {
		// Return empty array instead
		return nil
	}
}

func HandleQuerySmfSelectData(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QuerySmfSelectData")

	collName := "subscriptionData.provisionedData.smfSelectionSubscriptionData"
	ueId := request.Params["ueId"]
	servingPlmnId := request.Params["servingPlmnId"]
	response, problemDetails := QuerySmfSelectDataProcedure(collName, ueId, servingPlmnId)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
}

func QuerySmfSelectDataProcedure(collName string, ueId string,
	servingPlmnId string) (*map[string]interface{}, *models.ProblemDetails) {
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
	smfSelectionSubscriptionData := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	if smfSelectionSubscriptionData != nil {
		return &smfSelectionSubscriptionData, nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}
}

func HandleCreateSmsfContext3gpp(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle CreateSmsfContext3gpp")

	SmsfRegistration := request.Body.(models.SmsfRegistration)
	collName := "subscriptionData.contextData.smsf3gppAccess"
	ueId := request.Params["ueId"]

	CreateSmsfContext3gppProcedure(collName, ueId, SmsfRegistration)

	return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
}

func CreateSmsfContext3gppProcedure(collName string, ueId string, SmsfRegistration models.SmsfRegistration) {
	putData := util.ToBsonM(SmsfRegistration)
	putData["ueId"] = ueId
	filter := bson.M{"ueId": ueId}

	MongoDBLibrary.RestfulAPIPutOne(collName, filter, putData)
}

func HandleDeleteSmsfContext3gpp(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle DeleteSmsfContext3gpp")

	collName := "subscriptionData.contextData.smsf3gppAccess"
	ueId := request.Params["ueId"]

	DeleteSmsfContext3gppProcedure(collName, ueId)
	return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
}

func DeleteSmsfContext3gppProcedure(collName string, ueId string) {
	filter := bson.M{"ueId": ueId}
	MongoDBLibrary.RestfulAPIDeleteOne(collName, filter)
}

func HandleQuerySmsfContext3gpp(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QuerySmsfContext3gpp")

	collName := "subscriptionData.contextData.smsf3gppAccess"
	ueId := request.Params["ueId"]

	response, problemDetails := QuerySmsfContext3gppProcedure(collName, ueId)
	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func QuerySmsfContext3gppProcedure(collName string, ueId string) (*map[string]interface{}, *models.ProblemDetails) {
	filter := bson.M{"ueId": ueId}

	smsfRegistration := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	if smsfRegistration != nil {
		return &smsfRegistration, nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}
}

func HandleCreateSmsfContextNon3gpp(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle CreateSmsfContextNon3gpp")

	SmsfRegistration := request.Body.(models.SmsfRegistration)
	collName := "subscriptionData.contextData.smsfNon3gppAccess"
	ueId := request.Params["ueId"]

	CreateSmsfContextNon3gppProcedure(SmsfRegistration, collName, ueId)

	return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
}

func CreateSmsfContextNon3gppProcedure(SmsfRegistration models.SmsfRegistration, collName string, ueId string) {
	putData := util.ToBsonM(SmsfRegistration)
	putData["ueId"] = ueId
	filter := bson.M{"ueId": ueId}

	MongoDBLibrary.RestfulAPIPutOne(collName, filter, putData)
}

func HandleDeleteSmsfContextNon3gpp(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle DeleteSmsfContextNon3gpp")

	collName := "subscriptionData.contextData.smsfNon3gppAccess"
	ueId := request.Params["ueId"]

	DeleteSmsfContextNon3gppProcedure(collName, ueId)
	return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
}

func DeleteSmsfContextNon3gppProcedure(collName string, ueId string) {
	filter := bson.M{"ueId": ueId}
	MongoDBLibrary.RestfulAPIDeleteOne(collName, filter)
}

func HandleQuerySmsfContextNon3gpp(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QuerySmsfContextNon3gpp")

	ueId := request.Params["ueId"]
	collName := "subscriptionData.contextData.smsfNon3gppAccess"

	response, problemDetails := QuerySmsfContextNon3gppProcedure(collName, ueId)
	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func QuerySmsfContextNon3gppProcedure(collName string, ueId string) (*map[string]interface{}, *models.ProblemDetails) {
	filter := bson.M{"ueId": ueId}

	smsfRegistration := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	if smsfRegistration != nil {
		return &smsfRegistration, nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}
}

func HandleQuerySmsMngData(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QuerySmsMngData")

	collName := "subscriptionData.provisionedData.smsMngData"
	ueId := request.Params["ueId"]
	servingPlmnId := request.Params["servingPlmnId"]
	response, problemDetails := QuerySmsMngDataProcedure(collName, ueId, servingPlmnId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func QuerySmsMngDataProcedure(collName string, ueId string,
	servingPlmnId string) (*map[string]interface{}, *models.ProblemDetails) {
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
	smsManagementSubscriptionData := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	if smsManagementSubscriptionData != nil {
		return &smsManagementSubscriptionData, nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}
}

func HandleQuerySmsData(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QuerySmsData")

	ueId := request.Params["ueId"]
	servingPlmnId := request.Params["servingPlmnId"]
	collName := "subscriptionData.provisionedData.smsData"

	response, problemDetails := QuerySmsDataProcedure(collName, ueId, servingPlmnId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func QuerySmsDataProcedure(collName string, ueId string,
	servingPlmnId string) (*map[string]interface{}, *models.ProblemDetails) {
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}

	smsSubscriptionData := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	if smsSubscriptionData != nil {
		return &smsSubscriptionData, nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}
}

func HandlePostSubscriptionDataSubscriptions(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle PostSubscriptionDataSubscriptions")

	SubscriptionDataSubscriptions := request.Body.(models.SubscriptionDataSubscriptions)

	locationHeader := PostSubscriptionDataSubscriptionsProcedure(SubscriptionDataSubscriptions)

	headers := http.Header{
		"Location": {locationHeader},
	}
	return http_wrapper.NewResponse(http.StatusCreated, headers, SubscriptionDataSubscriptions)
}

func PostSubscriptionDataSubscriptionsProcedure(
	SubscriptionDataSubscriptions models.SubscriptionDataSubscriptions) string {
	nwdafSelf := nwdaf_context.NWDAF_Self()

	newSubscriptionID := strconv.Itoa(nwdafSelf.SubscriptionDataSubscriptionIDGenerator)
	nwdafSelf.SubscriptionDataSubscriptions[newSubscriptionID] = &SubscriptionDataSubscriptions
	nwdafSelf.SubscriptionDataSubscriptionIDGenerator++

	/* Contains the URI of the newly created resource, according
	   to the structure: {apiRoot}/subscription-data/subs-to-notify/{subsId} */
	locationHeader := fmt.Sprintf("%s/subscription-data/subs-to-notify/%s",
		nwdafSelf.GetIPv4GroupUri(nwdaf_context.NNWDAF_DR), newSubscriptionID)

	return locationHeader
}

func HandleRemovesubscriptionDataSubscriptions(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle RemovesubscriptionDataSubscriptions")

	subsId := request.Params["subsId"]

	problemDetails := RemovesubscriptionDataSubscriptionsProcedure(subsId)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, map[string]interface{}{})
	} else {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
}

func RemovesubscriptionDataSubscriptionsProcedure(subsId string) *models.ProblemDetails {
	nwdafSelf := nwdaf_context.NWDAF_Self()
	_, ok := nwdafSelf.SubscriptionDataSubscriptions[subsId]
	if !ok {
		problemDetails := &models.ProblemDetails{
			Cause:  "SUBSCRIPTION_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return problemDetails
	}
	delete(nwdafSelf.SubscriptionDataSubscriptions, subsId)
	return nil
}

func HandleQueryTraceData(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("Handle QueryTraceData")

	collName := "subscriptionData.provisionedData.traceData"
	ueId := request.Params["ueId"]
	servingPlmnId := request.Params["servingPlmnId"]

	response, problemDetails := QueryTraceDataProcedure(collName, ueId, servingPlmnId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func QueryTraceDataProcedure(collName string, ueId string,
	servingPlmnId string) (*map[string]interface{}, *models.ProblemDetails) {
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}

	traceData := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	if traceData != nil {
		return &traceData, nil
	} else {
		problemDetails := &models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusNotFound,
		}
		return nil, problemDetails
	}
}
