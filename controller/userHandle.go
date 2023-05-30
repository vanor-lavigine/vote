package controller

import (
	"crypto/elliptic"
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"voteList/dao"
	"voteList/response"
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
	//users := dao.QueryAllUsers()
	//var dList []big.Int
	//var xList []big.Int
	//var yList []big.Int
	var pkList []urs.PublicKey
	var skList []urs.PrivateKey
	//TODO: 投票流程待修改
	//for i := 0; i < len(users); i++ {
	//dstring := new(big.Int)
	//dstring, _ = dstring.SetString(users[i].D, 10)
	//xstring := new(big.Int)
	//xstring, _ = xstring.SetString(users[i].X, 10)
	//ystring := new(big.Int)
	//ystring, _ = ystring.SetString(users[i].Y, 10)
	//pk := urs.PublicKey{X: xstring, Y: ystring, Curve: elliptic.P256()}
	//pkList = append(pkList, pk)
	//sk := urs.PrivateKey{D: dstring, PublicKey: pk}
	//skList = append(skList, sk)
	//}

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
	//fmt.Printf(r.Body)
	fmt.Printf("abc", r.Method)
	//username := r.FormValue("username")
	//fmt.Printf(username)1

}

// 创建候选者
type Response struct {
	Code      int    `json:"code"`
	Data      *Data  `json:"data"`
	ErrorCode string `json:"errorCode"`
}

type Data struct {
	Username string `json:"username"`
	ID       string `json:"id"`
	Pbk      string `json:"pbk""`
	Prk      string `json:"prk""`
}

func checkLoginAndAdminStatus() {
	return
}

func CreateCandidateHandle(w http.ResponseWriter, r *http.Request) {
	//resp := &Response{}
	//if r.Method == "POST" {
	//	username := r.FormValue("username")
	//	id := r.FormValue("id")
	//	pubk := r.FormValue("pubk")
	//	prk := r.FormValue("prk")
	//	loggedIn, isAdmin := checkLoginAndAdminStatus(username) // Implement this function according to your login and admin status checking logic
	//	if !loggedIn {
	//		resp.Code = 401
	//		resp.ErrorCode = "User not logged in"
	//	} else if !isAdmin {
	//		resp.Code = 403
	//		resp.ErrorCode = "User not an admin"
	//	} else {
	//		exists, _ := dao.CheckCandidateExists(username) // Implement this function to check if the candidate exists
	//		if exists {
	//			resp.Code = 400
	//			resp.ErrorCode = "canxxxxExist"
	//		} else {
	//			//pubKey, privKey := generateKeys()    */                 // Implement this function to generate public and private keys
	//			err := dao.CreateCandidate(username, id, pubk, prk) // Implement this function to create a candidate in your database
	//			if err != nil {
	//				resp.Code = 500
	//				resp.ErrorCode = "Backend error"
	//			} else {
	//				resp.Code = 200
	//				resp.Data = &Data{
	//					Username: username,
	//					Prk:      prk, // You may want to return something else as the ID
	//				}
	//			}
	//		}
	//	}
	//} else {
	//	resp.Code = 405
	//	resp.ErrorCode = "Invalid method"
	//}
	//bodyBytes, err := ioutil.ReadAll(r.Body)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//bodyString := string(bodyBytes)
	//
	//fmt.Printf("%s\n", bodyString)
	//fmt.Printf("abc")
	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(resp)
}

// 候选人列表
type CandidateListResponse struct {
	Code      int       `json:"code"`
	Data      *CandData `json:"data"`
	ErrorCode string    `json:"errorCode"`
}

type CandData struct {
	List []CandidateInfo `json:"list"`
}

type CandidateInfo struct {
	Username string `json:"username"`
	ID       string `json:"id"`
}

func GetCandidateListHandle(w http.ResponseWriter, r *http.Request) {
	resp := &CandidateListResponse{}

	/*if r.Method == "GET" {
		loggedIn, isAdmin := checkLoginAndAdminStatus(r) // Implement this function according to your login and admin status checking logic
		if !loggedIn {
			resp.Code = 401
			resp.ErrorCode = "User not logged in"
		} else if !isAdmin {
			resp.Code = 403
			resp.ErrorCode = "User not an admin"
		} else {
			candidates, err := dao.GetCandidateList() // Implement this function to get a list of candidates from your database
			if err != nil {
				resp.Code = 500
				resp.ErrorCode = "Backend error"
			} else {
				resp.Code = 200
				resp.Data = &CandData{
					List: candidates,
				}
			}
		}
	} else {
		resp.Code = 405
		resp.ErrorCode = "Invalid method"
	}*/

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	bodyString := string(bodyBytes)

	fmt.Printf("%s\n", bodyString)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
	fmt.Printf("abc")
}

func handleError(w http.ResponseWriter, statusCode int, errorCode string, errorMessage string) {
	// 构造错误响应结构体
	_response := response.ApiResponse{
		Code:      statusCode,
		Data:      "",
		ErrorCode: errorCode,
		ErrorMsg:  errorMessage,
	}

	// 将错误响应转换为JSON格式
	responseJSON, err := json.Marshal(_response)
	if err != nil {
		http.Error(w, "Failed to marshal error response", http.StatusInternalServerError)
		return
	}

	// 设置响应头并写入错误响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(responseJSON)
}

func CreateUserHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		username := r.FormValue("username")
		password := r.FormValue("password")

		fmt.Println("username", username)
		fmt.Println("password", password)

		exists, err := dao.CheckUserExists(username)
		if err != nil {
			handleError(w, http.StatusInternalServerError, "服务出错", "服务出错")
			return
		}

		// 检查用户是否存在
		if exists {
			fmt.Sprintf(`user %s is already exists`, username)
			handleError(w, http.StatusInternalServerError, response.UserExists, "服务出错")
			return
		}

		// 调用插入用户方法
		insertErr := dao.InsertUser(username, password)
		if insertErr != nil {
			handleError(w, http.StatusInternalServerError, response.CreateUserFailed, insertErr.Error())
			return
		}

		// 构造响应结构体
		_response := response.ApiResponse{
			Code:      http.StatusOK,
			Data:      "User inserted successfully",
			ErrorCode: "",
			ErrorMsg:  "",
		}

		// 将响应转换为JSON格式
		responseJSON, err := json.Marshal(_response)
		if err != nil {
			handleError(w, http.StatusInternalServerError, "Failed to marshal response", "")
			return
		}

		// 设置响应头并写入响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	}

}
