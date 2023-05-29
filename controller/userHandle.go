package controller

import (
	"crypto/elliptic"
	crand "crypto/rand"
	"fmt"
	"html/template"
	"math/big"
	"net/http"
	"voteList/dao"
	_ "voteList/ring"
	"voteList/sdkInit"
	"voteList/service"
	"voteList/urs"
)

const numOfKeys = 1000

//var (
//	DefaultCurve = elliptic.P256()
//	keyring      *ring.PublicKeyRing
//	testkey      *ring.PrivateKey
//	testmsg      []byte
//	testsig      *ring.RingSign
//)

const (
	configFile  = "conf.yaml"
	initialized = false
	EduCC       = "mycc"
)

// LoginHandle 用户登录
func LoginHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t := template.Must(template.ParseFiles("views/pages/user/login.html"))
		t.Execute(w, "")
	} else {
		username := r.FormValue("username")
		password := r.FormValue("password")
		fmt.Println(username, password)
		u, _ := dao.CheckUserNameAndPassword(username)
		if u != nil {
			t := template.Must(template.ParseFiles("views/pages/user/votelist.html"))
			t.Execute(w, "")
		} else {
			t := template.Must(template.ParseFiles("views/pages/user/login.html"))
			t.Execute(w, "")
		}
	}
}

// RegisterHandle 用户注册
func RegisterHandle(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println(username, password)
	sk, err := urs.GenerateKey(elliptic.P256(), crand.Reader)
	if err != nil {
		panic(err)
	}
	d := sk.D
	dSting := d.String()
	x := sk.X
	xString := x.String()
	y := sk.Y
	yString := y.String()

	err = dao.SaveUser(username, password, dSting, xString, yString)
	if err != nil {
		panic(err)
	}
}

// VoteHandle 用户投票
func VoteHandle(w http.ResponseWriter, r *http.Request) {
	voteid := r.FormValue("voteid")
	votename := r.FormValue("votename")
	fmt.Println(voteid, votename)
	users := dao.QueryAllUsers()
	//var dList []big.Int
	//var xList []big.Int
	//var yList []big.Int
	var pkList []urs.PublicKey
	var skList []urs.PrivateKey
	for i := 0; i < len(users); i++ {
		dstring := new(big.Int)
		dstring, _ = dstring.SetString(users[i].D, 10)
		xstring := new(big.Int)
		xstring, _ = xstring.SetString(users[i].X, 10)
		ystring := new(big.Int)
		ystring, _ = ystring.SetString(users[i].Y, 10)
		pk := urs.PublicKey{X: xstring, Y: ystring, Curve: elliptic.P256()}
		pkList = append(pkList, pk)
		sk := urs.PrivateKey{D: dstring, PublicKey: pk}
		skList = append(skList, sk)
	}

	ring := urs.PublicKeyRing{Ring: pkList}

	sign, _ := urs.Sign(crand.Reader, &skList[0], &ring, []byte(votename))
	fmt.Println(sign)

	initInfo := &sdkInit.InitInfo{

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
	//安装实例化链码
	channelClient, err := sdkInit.GetClient(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(channelClient)
	setup := service.ServiceSetup{
		ChaincodeID: EduCC,
		Client:      channelClient,
	}
	test := urs.Verify(&ring, []byte(votename), sign)
	fmt.Println(test)
	var msg string
	if test {
		msg, _ = setup.SetInfo(votename, "vote success!")
	} else {
		t := template.Must(template.ParseFiles("views/pages/user/votelist.html"))
		t.Execute(w, "")

	}
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(msg)
	}
}

func Votecheck(w http.ResponseWriter, r *http.Request) {
	votename := r.FormValue("votename")
	initInfo := &sdkInit.InitInfo{

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
	//安装实例化链码
	channelClient, err := sdkInit.GetClient(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(channelClient)
	setup := service.ServiceSetup{
		ChaincodeID: EduCC,
		Client:      channelClient,
	}
	message, _ := setup.GetInfo(votename)
	fmt.Println(message)
	t := template.Must(template.ParseFiles("views/pages/user/votesuccess.html"))
	t.Execute(w, message)

}
func Abc(w http.ResponseWriter, r *http.Request) {
	print(r.Body)
	fmt.Printf("abc", r.Method)
	//username := r.FormValue("username")
	//fmt.Printf(username)1

}
