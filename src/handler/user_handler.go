package handler

import (
	"filestore-service/src/db"
	"filestore-service/src/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	pwd_salt = "*#890"
)

// 处理用户注册请求
func SignupHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		data, err := ioutil.ReadFile("../static/view/signup.html")
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		response.Write(data)
		return
	} else if request.Method == "POST" {
		request.ParseForm()
		username := request.Form.Get("username")
		password := request.Form.Get("password")

		if len(username) < 3 || len(password) < 5 {
			response.Write([]byte("invalid parameter"))
			return
		}
		encode_password := util.Sha1([]byte(password + pwd_salt))
		success := db.UserSignup(username, encode_password)
		if success {
			response.Write([]byte("SUCCESS"))
		} else {
			response.Write([]byte("FAILED"))
		}
	}
}

// 登录
func SignInHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		data, err := ioutil.ReadFile("../static/view/signin.html")
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		response.Write(data)
		return
	}
	request.ParseForm()
	username := request.Form.Get("username")
	password := request.Form.Get("password")
	encode_password := util.Sha1([]byte(password + pwd_salt))
	// 校验用户名及密码
	pwdCheck := db.UserSignIn(username, encode_password)
	if !pwdCheck {
		response.Write([]byte("FAILED"))
		return
	}
	// 生成访问凭证（token）
	token := GenToken(username)
	success := db.UpdateToken(username, token)
	if !success {
		response.Write([]byte("FAILED"))
		return
	}
	//登录成功后重定向首页
	// response.Write([]byte("http://" + request.Host + "/static/view/home.html"))
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + request.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	response.Write(resp.JSONBytes())
}

// 查询用户信息
func UserInfoHandler(response http.ResponseWriter, request *http.Request) {
	// 解析请求参数
	request.ParseForm()
	username := request.Form.Get("username")
	token := request.Form.Get("token")
	//验证token是否有效
	ret := IsTokenVaild(token)
	fmt.Println(ret)
	if !ret {
		response.WriteHeader(http.StatusForbidden)
		return
	}
	//查询用户信息
	user, err := db.GetUserInfo(username)
	fmt.Println(err)
	if err != nil {
		response.WriteHeader(http.StatusForbidden)
		return
	}
	fmt.Println(user)
	//组装响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	response.Write(resp.JSONBytes())
}

func GenToken(username string) string {
	//40位token md5(username + timestamp + token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	token_prefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return token_prefix + ts[:8]
}

func IsTokenVaild(token string) bool {
	// TODO: 判断token时效性，是否过期

	// TODO: 从数据库查询用户对应的token信息

	//TODO: 对比两个token是否一致

	return true
}
