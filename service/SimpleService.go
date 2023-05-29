package service

import "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"

//向分类帐中添加状态
func (t *ServiceSetup) SetInfo(name, num string) (string, error) {
	eventID := "eventSetInfo"
	reg, notifier := regiterEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "set",
		Args: [][]byte{[]byte(name), []byte(num), []byte(eventID)}}
	response, err := t.Client.Execute(req)

	if err != nil {
		return "", err
	}
	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}
	return string(response.TransactionID), nil
}

//查询分类账中信息
func (t *ServiceSetup) GetInfo(name string) (string, error) {
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "get", Args: [][]byte{[]byte(name)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return "", err
	}
	return string(respone.Payload), nil
}
