package callback

import (
	"context"
	"free5gc/lib/openapi/Nnwdaf_DataRepository"
	"free5gc/lib/openapi/models"
	nwdaf_context "free5gc/src/nwdaf/context"
	"free5gc/src/nwdaf/logger"
)

func SendOnDataChangeNotify(ueId string, notifyItems []models.NotifyItem) {
	nwdafSelf := nwdaf_context.NWDAF_Self()
	configuration := Nnwdaf_DataRepository.NewConfiguration()
	client := Nnwdaf_DataRepository.NewAPIClient(configuration)

	for _, subscriptionDataSubscription := range nwdafSelf.SubscriptionDataSubscriptions {
		if ueId == subscriptionDataSubscription.UeId {
			onDataChangeNotifyUrl := subscriptionDataSubscription.CallbackReference

			dataChangeNotify := models.DataChangeNotify{}
			dataChangeNotify.UeId = ueId
			dataChangeNotify.OriginalCallbackReference = []string{subscriptionDataSubscription.OriginalCallbackReference}
			dataChangeNotify.NotifyItems = notifyItems
			httpResponse, err := client.DataChangeNotifyCallbackDocumentApi.OnDataChangeNotify(context.TODO(),
				onDataChangeNotifyUrl, dataChangeNotify)
			if err != nil {
				if httpResponse == nil {
					logger.HttpLog.Errorln(err.Error())
				} else if err.Error() != httpResponse.Status {
					logger.HttpLog.Errorln(err.Error())
				}
			}
		}
	}
}

func SendPolicyDataChangeNotification(policyDataChangeNotification models.PolicyDataChangeNotification) {
	nwdafSelf := nwdaf_context.NWDAF_Self()

	for _, policyDataSubscription := range nwdafSelf.PolicyDataSubscriptions {
		policyDataChangeNotificationUrl := policyDataSubscription.NotificationUri

		configuration := Nnwdaf_DataRepository.NewConfiguration()
		client := Nnwdaf_DataRepository.NewAPIClient(configuration)
		httpResponse, err := client.PolicyDataChangeNotificationCallbackDocumentApi.PolicyDataChangeNotification(
			context.TODO(), policyDataChangeNotificationUrl, policyDataChangeNotification)
		if err != nil {
			if httpResponse == nil {
				logger.HttpLog.Errorln(err.Error())
			} else if err.Error() != httpResponse.Status {
				logger.HttpLog.Errorln(err.Error())
			}
		}
	}
}
