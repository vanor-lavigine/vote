package service

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"time"
)

type ServiceSetup struct {
	ChaincodeID string
	Client      *channel.Client
}

//注册事件函数
func regiterEvent(client *channel.Client, chaincodeID, eventID string) (fab.Registration, <-chan *fab.CCEvent) {
	//注册事件
	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeID, eventID)
	if err != nil {
		fmt.Println("注册链码事件失败！:%s", err)
	}
	fmt.Println("注册事件成功！")
	return reg, notifier

}

//接收链码事件
func eventResult(notifier <-chan *fab.CCEvent, eventID string) error {
	//接收链码事件
	select {
	case ccEvent := <-notifier:
		fmt.Println("接收到链码事件:%v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return fmt.Errorf("不能根据指定的事件ID接收到相应的链码事件(%s)", eventID)
	}
	return nil
}
