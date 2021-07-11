package user_service

import (
	"go-gofram-chat/app/models"
	"go-gofram-chat/app/service/helper"
	"go-gofram-chat/app/service/session"
	"go-gofram-chat/app/service/validator"
	"log"
	"strconv"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func Login(r *ghttp.Request) {

	username := r.GetString("username")
	pwd := r.GetString("password")
	avatarId := r.GetString("avatar_id")

	var u validator.User
	u_temp := validator.User{}
	u.Username = username
	u.Password = pwd
	u.AvatarId = avatarId

	if err := r.GetStruct(&u); err != nil {
		r.Response.WriteJson(g.Map{"code": 5000, "msg": err.Error()})
		return
	}

	user := models.FindUserByField("username", username)
	userInfo := user
	//log.Println("user info :", userInfo)
	md5Pwd := helper.Md5Encrypt(pwd)

	if userInfo.ID > 0 {
		// json check user
		// check pwd
		if userInfo.Password != md5Pwd {
			r.Response.WriteJson(g.Map{"code": 5000, "msg": "密碼錯誤！！"})
			return
		}

	} else {
		//check valid
		if err := g.Validator().Data(u).CheckStruct(u_temp); err != nil {
			r.Response.WriteJson(g.Map{"code": 5001, "msg": err.FirstString()})
			log.Println(err.Items())
			return
		}

		// if no in db, create new user
		userInfo = models.AddUser(map[string]interface{}{
			"username":  username,
			"password":  md5Pwd,
			"avatar_id": avatarId,
		})

		user = models.FindUserByField("username", username)
		userInfo = user
	}

	if userInfo.ID > 0 {
		session.SaveAuthSession(r, string(strconv.Itoa(int(userInfo.ID))))
		r.Response.WriteJson(g.Map{"code": 0, "msg": "save session"})
		return
	} else {
		//log.Println("user id : ", userInfo.ID)
		r.Response.WriteJson(g.Map{"code": 5001, "msg": "System error : save session fail"})

		return
	}
}

func GetUserInfo(r *ghttp.Request) map[string]interface{} {
	return session.GetSessionUserInfo(r)
}

func Logout(r *ghttp.Request) {
	session.ClearAuthSession(r)
	r.Response.RedirectTo("/")

}
