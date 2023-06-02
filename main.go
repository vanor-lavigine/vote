package main

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	_ "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"net/http"
	"time"
)

const (
	configFile  = "conf.yaml"
	initialized = false
	EduCC       = "mycc"
)

type ApiResponse struct {
	Code      int         `json:"code"`
	Data      interface{} `json:"data"`
	ErrorCode string      `json:"errorCode"`
	ErrorMsg  string      `json:"ErrorMsg"`
}

var db *sql.DB

type Session struct {
	ID        string
	Username  string
	ExpiresAt time.Time
}

var sessions map[string]Session

func main() {
	var err error

	/*initInfo := &sdkInit.InitInfo{

		ChannelID:      "hustgym",
		ChannelConfig:  "/home/u/go/src/fixturesPIC/channel-artifacts/HUSTgym.tx",
		OrgAdmin:       "Admin",
		OrgName:        "HUST",
		OrdererOrgName: "orderer.test.com",

		ChaincodeID:     EduCC,
		ChaincodeGoPath: "/home/u/go/",
		ChaincodePath:   "voteList/chaincode/",
		UserName:        "User1",
	}

	sdk, err := sdkInit.SetupSDK(configFile, initialized)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	defer sdk.Close()

	err = sdkInit.CreateChannel(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	channelClient, err := sdkInit.InstallAndInstantiateCC(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(channelClient)

	//===========================================//

	serviceSetup := service.ServiceSetup{
		ChaincodeID: EduCC,
		Client:      channelClient,
	}

	msg, err := serviceSetup.SetInfo("van", "123")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(msg)
	}
	msg, err = serviceSetup.SetInfo("van", "345")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(msg)
	}
	msg, err = serviceSetup.GetInfo("van")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("van value:", msg)
	}

	// Setup MySQL connection
	*/

	sessions = make(map[string]Session)
	db, err = sql.Open("mysql", "debian-sys-maint:XMQWnyGB6Or12Oxk@tcp(localhost:3306)/vote")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fs := http.FileServer(http.Dir("view3/my-app/build"))
	//http.Handle("/static/", http.StripPrefix("view3/my-app/build/static/", fs))
	http.Handle("/", fs)

	http.HandleFunc("/register", errorHandler(registerHandler))
	http.HandleFunc("/login", errorHandler(loginHandler))
	http.HandleFunc("/logout", errorHandler(logoutHandler))
	http.HandleFunc("/createCandidate", withAuth(adminOnly(errorHandler(createCandidateHandler))))
	http.HandleFunc("/deleteCandidate", withAuth(adminOnly(errorHandler(deleteCandidateHandler))))
	http.HandleFunc("/listCandidates", errorHandler(listCandidatesHandler))
	fmt.Println("服务开启成功：地址为", "http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func allowCORS(w http.ResponseWriter) {
	// Set CORS headers
	// "Access-Control-Allow-Origin": "*"
	// Or you could be more specific:
	// "Access-Control-Allow-Origin": "http://localhost:3000"
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Allow methods
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	// Allow headers
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func errorHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allowCORS(w)
		defer func() {
			if err, ok := recover().(error); ok {
				json.NewEncoder(w).Encode(ApiResponse{
					Code:      500,
					ErrorCode: "internal_error",
					ErrorMsg:  err.Error(),
				})
			}
		}()
		fn(w, r)
	}
}

func withAuth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allowCORS(w)
		if _, err := checkLogin(r); err != nil {
			json.NewEncoder(w).Encode(ApiResponse{
				Code:      401,
				ErrorCode: "not_logged_in",
				ErrorMsg:  "You need to log in to perform this action.",
			})
			return
		}
		fn(w, r)
	}
}

func adminOnly(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allowCORS(w)
		username, _ := checkLogin(r)
		if username != "admin" {
			json.NewEncoder(w).Encode(ApiResponse{
				Code:      403,
				ErrorCode: "forbidden",
				ErrorMsg:  "Only admin can perform this action.",
			})
			return
		}
		fn(w, r)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	allowCORS(w)
	// Create a new struct to hold the request body
	type RegisterRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var request RegisterRequest

	// Decode the request body into the struct
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		panic(err)
	}

	username := request.Username
	password := request.Password

	// Hash the password
	hasher := md5.New()
	hasher.Write([]byte(password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM user_new WHERE username=?)", username).Scan(&exists)
	if err != nil {
		panic(err)
	}
	if exists {
		json.NewEncoder(w).Encode(ApiResponse{
			Code:      400,
			ErrorCode: "user_exists",
			ErrorMsg:  "Username is already taken.",
		})
		return
	}
	hash := md5.New()
	hash.Write([]byte(username))
	hashValue := hex.EncodeToString(hash.Sum(nil))

	_, err = db.Exec("INSERT INTO user_new (username, password, hash) VALUES (?, ?, ?)", username, hashedPassword, hashValue)
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(ApiResponse{
		Code: 200,
		Data: "User registered successfully.",
	})
}

func generateSessionID() string {
	_uuid := uuid.New()
	return _uuid.String()
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	allowCORS(w)
	type RegisterRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var request RegisterRequest

	// Decode the request body into the struct
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		panic(err)
	}

	username := request.Username
	password := request.Password

	// Hash the password
	hasher := md5.New()
	hasher.Write([]byte(password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM user_new WHERE username=? AND password=?)", username, hashedPassword).Scan(&exists)
	if err != nil {
		panic(err)
	}
	if !exists {
		json.NewEncoder(w).Encode(ApiResponse{
			Code:      400,
			ErrorCode: "invalid_credentials",
			ErrorMsg:  "Invalid username or password.",
		})
		return
	}

	sessionID := generateSessionID()
	expiration := time.Now().Add(14 * 24 * time.Hour)
	session := Session{
		ID:        sessionID,
		Username:  username,
		ExpiresAt: expiration,
	}
	sessions[sessionID] = session

	cookie := http.Cookie{
		Name:     "session-id",
		Value:    sessionID,
		Path:     "/",
		Expires:  expiration,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	fmt.Println("session", session.ID)

	json.NewEncoder(w).Encode(ApiResponse{
		Code: 200,
		Data: "Login successful.",
	})
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	allowCORS(w)
	cookie, err := r.Cookie("session-id")
	if err != nil {
		json.NewEncoder(w).Encode(ApiResponse{
			Code:      401,
			ErrorCode: "not_logged_in",
			ErrorMsg:  "You are not logged in.",
		})
		return
	}

	sessionID := cookie.Value
	delete(sessions, sessionID)

	// Clear cookie on the client side
	cookie = &http.Cookie{
		Name:    "session-id",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(w, cookie)
	json.NewEncoder(w).Encode(ApiResponse{
		Code: 200,
		Data: "Logout successful.",
	})
}

func createCandidateHandler(w http.ResponseWriter, r *http.Request) {
	allowCORS(w)
	//TODO: zengjia privatekey he publicKey
	type RegisterRequest struct {
		Username string `json:"username"`
	}
	var request RegisterRequest

	// Decode the request body into the struct
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		panic(err)
	}

	username := request.Username

	privateKey, publicKey := GenRsaKey()

	//fmt.Println(string(privateKey))
	//fmt.Println(string(publicKey))

	fmt.Println("createCandidateHandler after login")

	_, err = db.Exec("INSERT INTO candidates (username, valid, publicKey, privateKey) VALUES (?, true, ?, ?)", username, publicKey, privateKey)
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(ApiResponse{
		Code: 200,
		Data: "Candidate created successfully.",
	})
}

func deleteCandidateHandler(w http.ResponseWriter, r *http.Request) {
	allowCORS(w)
	type RegisterRequest struct {
		Username string `json:"username"`
	}
	var request RegisterRequest

	// Decode the request body into the struct
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		panic(err)
	}

	username := request.Username

	_, err = db.Exec("UPDATE candidates SET valid = false WHERE username = ?", username)
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(ApiResponse{
		Code: 200,
		Data: "Candidate deleted successfully.",
	})
}

func listCandidatesHandler(w http.ResponseWriter, r *http.Request) {
	allowCORS(w)
	rows, err := db.Query("SELECT id, username FROM candidates WHERE valid=true")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var candidates []map[string]interface{}
	for rows.Next() {
		var id int
		var username string
		err := rows.Scan(&id, &username)
		if err != nil {
			panic(err)
		}
		candidates = append(candidates, map[string]interface{}{
			"id":       id,
			"username": username,
		})
	}

	json.NewEncoder(w).Encode(ApiResponse{
		Code: 200,
		Data: candidates,
	})
}

func checkLogin(r *http.Request) (string, error) {
	//session, err := store.Get(r, "session-name")
	//if err != nil {
	//	return "", err
	//}
	//
	//username, ok := session.Values["username"].(string)
	//if !ok {
	//	return "", errors.New("not logged in")
	//}
	//
	//return username, nil

	cookie, err := r.Cookie("session-id")
	if err != nil {
		return "", err
	}

	sessionID := cookie.Value
	for _, session := range sessions {
		if session.ID == sessionID {
			return session.Username, nil
		}
	}
	return "", nil
}

func GenRsaKey() (prvkey, pubkey []byte) {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	prvkey = pem.EncodeToMemory(block)
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		panic(err)
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	pubkey = pem.EncodeToMemory(block)
	return
}

// gongyao jiami
func RsaEncrypt(data, keyBytes []byte) []byte {
	//解密pem格式的公钥
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("public key error"))
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pub, data)
	if err != nil {
		panic(err)
	}
	return ciphertext
}

// 私钥解密
func RsaDecrypt(ciphertext, keyBytes []byte) []byte {
	//获取私钥
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("private key error!"))
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	// 解密
	data, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	if err != nil {
		panic(err)
	}
	return data
}
