package db

import (
	"filestore-service/src/db/mysql"
	"fmt"
)

// 通过用户名密码完成注册操作
func UserSignup(username string, password string) bool {
	stmt, err := mysql.DBConn().Prepare("insert ignore into tbl_user (user_name,user_pwd) values (?,?)")
	if err != nil {
		fmt.Println("Failed to insert,err:" + err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(username, password)
	if err != nil {
		fmt.Println("Failed to insert,err:" + err.Error())
		return false
	}
	if rows, err := ret.RowsAffected(); err == nil && rows > 0 {
		return true
	}
	return false
}

// 判断密码是否一致
func UserSignIn(username string, encode_password string) bool {
	stmt, err := mysql.DBConn().Prepare("select * from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if rows == nil {
		fmt.Println("username not found:" + username)
		return false
	}

	prows := mysql.ParseRows(rows)
	if len(prows) > 0 && string(prows[0]["user_pwd"].([]byte)) == encode_password {
		return true
	}
	return false
}

// 刷新登录的token
func UpdateToken(username string, token string) bool {
	stmt, err := mysql.DBConn().Prepare("replace into tbl_user_token(user_name,user_token) values (?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

// 用户信息查询
func GetUserInfo(username string) (User, error) {
	user := User{}
	stmt, err := mysql.DBConn().Prepare("select user_name, signup_at from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Print(err.Error())
		return user, err
	}
	defer stmt.Close()
	fmt.Println(username)
	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		return user, err
	}
	return user, nil
}
