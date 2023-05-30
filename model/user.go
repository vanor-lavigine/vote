package model

// User 结构体
//type User struct {
//	Uid      int
//	Username string
//	Password string
//	D        string
//	X        string
//	Y        string
//}

type User struct {
	Username   string
	Password   string
	Hash       string
	ID         int
	PrivateKey string
	PublicKey  string
	Valid      bool
}
