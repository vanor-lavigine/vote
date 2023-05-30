package main

import (
	"fmt"
	"net/http"
	"voteList/controller"
)

const (
	configFile  = "conf.yaml"
	initialized = false
	EduCC       = "mycc"
)

func main() {

	/*initInfo := &sdkInit.InitInfo{

		ChannelID:      "hustgym",
		ChannelConfig:  "/home/u/go/src/fixturesPIC/channel-artifacts/HUSTgym.tx",
		OrgAdmin:       "Admin",
		OrgName:        "HUST",
		OrdererOrgName: "orderer.test.com",

		ChaincodeID:     "mycc",
		ChaincodeGoPath: "/home/u/go/",
		ChaincodePath:   "voteList/chaincode/",
		UserName:        "User1",
	}

	sdk, err := sdkInit.SetupSDK(configFile, initialized)
	if err != nil {
		fmt.Printf(err.Error())
	}

	defer sdk.Close()
	err = sdkInit.CreateChannel(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
	}
	//安装实例化链码

	sdkInit.InstallAndInstantiateCC(sdk, initInfo)*/

	// 处理静态资源路径
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("views/static/"))))
	http.Handle("/pages/", http.StripPrefix("/pages/", http.FileServer(http.Dir("views/pages/"))))
	http.HandleFunc("/login", controller.LoginHandle)
	http.HandleFunc("/register", controller.RegisterHandle)
	http.HandleFunc("/vote", controller.VoteHandle)
	http.HandleFunc("/check", controller.Votecheck)

	http.HandleFunc("/GetCandidateListHandle", controller.GetCandidateListHandle)
	http.HandleFunc("/CreateCandidateHandle", controller.CreateCandidateHandle)
	http.HandleFunc("/abc", controller.Abc)
	fmt.Println("服务开启成功：地址为", "http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}
