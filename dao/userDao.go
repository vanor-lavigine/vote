package dao

import (
	"fmt"
	"voteList/model"
	"voteList/util"
)

// CheckUserNameAndPassword 检查用户名和密码
func CheckUserNameAndPassword(username string) (*model.User, error) {
	if util.InitDb() != true {
		fmt.Println("db false!")
	}
	sqlStr := "select * from user where username=?"
	row := util.Db.QueryRow(sqlStr, username)
	u := &model.User{}
	err := row.Scan(&u.Uid, &u.Username, &u.Password, &u.D, &u.X, &u.Y)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// SaveUser 保存用户信息
func SaveUser(username string, password string, dstirng string, xstring string, ystring string) error {
	if util.InitDb() != true {
		fmt.Println("db false!")
	}
	sqlStr := "insert into user(username,password,d,x,y) values(?,?,?,?,?)"
	_, err := util.Db.Exec(sqlStr, username, password, dstirng, xstring, ystring)
	fmt.Println(err)
	if err != nil {
		return err
	}
	return nil
}

//QueryAllUsers query all users
func QueryAllUsers() []*model.User {
	if util.InitDb() != true {
		fmt.Println("db false!")
	}
	sqlStr := "select * from user"
	rows, err := util.Db.Query(sqlStr)
	if err != nil {
		fmt.Println(err)
	}
	//u := &model.User{}
	//i := 0
	var users []*model.User
	for rows.Next() {
		u := &model.User{}
		rows.Scan(&u.Uid, &u.Username, &u.Password, &u.D, &u.X, &u.Y)
		fmt.Println(&u.Uid, u.Username, u.Password, &u.D, &u.X, &u.Y)
		users = append(users, u)
	}
	return users
}

//func main() {
//	users := QueryAllUsers()
//	for i := 0; i < len(users); i++ {
//		fmt.Println(users[i])
//	}
//}
