package dao

import (
	"fmt"
	"voteList/model"
	"voteList/util"
)

// CheckUserNameAndPassword 检查用户名和密码
func CheckUserNameAndPassword(username string) (*model.User, error) {
	//if util.InitDb() != true {
	//	fmt.Println("db false!")
	//}
	//sqlStr := "select * from user where username=?"
	//row := util.Db.QueryRow(sqlStr, username)
	//u := &model.User{}
	//err := row.Scan(&u.Uid, &u.Username, &u.Password, &u.D, &u.X, &u.Y)
	//if err != nil {
	//	return nil, err
	//}
	//return u, nil
	return nil, nil
}

func CheckUserExists(username string) (bool, error) {
	if util.InitDb() != true {
		fmt.Println("db false!")
		return false, nil
	}

	//TODO: Mock 不存在用户
	return false, nil

	// 准备查询语句
	query := "SELECT COUNT(*) FROM user WHERE username = ?"

	// 执行查询
	var count int
	err := util.Db.QueryRow(query, username).Scan(&count)
	if err != nil {
		return false, err
	}

	// 判断是否存在用户
	if count > 0 {
		return true, nil
	}

	return false, nil
}

func InsertUser(username, password string) error {
	hash := util.Hash(username)
	passwordHash := util.Hash(password)
	privateKey := "privateKey"
	publicKey := "publicKey"

	query := "INSERT INTO user (username, password, hash, privateKey, publicKey) VALUES (?, ?, ?)"
	_, err := util.Db.Exec(query, username, passwordHash, hash, privateKey, publicKey)
	if err != nil {
		return err
	}
	return nil
}

// SaveUser 保存用户信息
func SaveUser(username string, password string, dstirng string, xstring string, ystring string) error {
	//if util.InitDb() != true {
	//	fmt.Println("db false!")
	//}
	//sqlStr := "insert into user(username,password,d,x,y) values(?,?,?,?,?)"
	//_, err := util.Db.Exec(sqlStr, username, password, dstirng, xstring, ystring)
	//fmt.Println(err)
	//if err != nil {
	//	return err
	//}
	return nil
}

// QueryAllUsers query all users
//func QueryAllUsers() []*model.User {
//if util.InitDb() != true {
//	fmt.Println("db false!")
//}
//sqlStr := "select * from user"
//rows, err := util.Db.Query(sqlStr)
//if err != nil {
//	fmt.Println(err)
//}
////u := &model.User{}
////i := 0
//var users []*model.User
//for rows.Next() {
//	u := &model.User{}
//	rows.Scan(&u.Uid, &u.Username, &u.Password, &u.D, &u.X, &u.Y)
//	fmt.Println(&u.Uid, u.Username, u.Password, &u.D, &u.X, &u.Y)
//	users = append(users, u)
//}
//return users
//}

func checkLoginAndAdminStatus() {

}
func GetCandidateList() {

}

func CreateCandidate(username string, id string, pubk string, prk string) error {
	if util.InitDb() != true {
		fmt.Println("db false!")
	}
	sqlStr := "insert into candidate(username, id, pubk, prk) values(?,?,?,?)"
	_, err := util.Db.Exec(sqlStr, username, id, pubk, prk)
	fmt.Println(err)
	if err != nil {
		return err
	}
	return nil
}

//func main() {
//	users := QueryAllUsers()
//	for i := 0; i < len(users); i++ {
//		fmt.Println(users[i])
//	}
//}
