package session

import (
	"go-gofram-chat/app/models"
	"strconv"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gsession"
)

func EnableCookieSession(s *ghttp.Server) {
	s.SetConfigWithMap(g.Map{
		"SessionMaxAge":  360 * time.Minute,
		"SessionStorage": gsession.NewStorageMemory(),
	})
}

// save seesion信息
func SaveAuthSession(r *ghttp.Request, info interface{}) {
	r.Session.Set("uid", info)
}

func GetSessionUserInfo(r *ghttp.Request) map[string]interface{} {

	uid := r.Session.Get("uid")

	data := make(map[string]interface{})
	if uid != nil {
		user := models.FindUserByField("id", uid.(string))
		data["uid"] = user.ID
		data["username"] = user.Username
		data["avatar_id"] = user.AvatarId
	}
	return data
}

// 退出時清除session
func ClearAuthSession(r *ghttp.Request) {

	r.Session.Clear()

}

func HasSession(r *ghttp.Request) bool {

	if sessionValue := r.Session.Get("uid"); sessionValue == nil {
		return false
	}
	return true
}

func AuthSessionMiddle() ghttp.HandlerFunc {
	return func(r *ghttp.Request) {

		sessionValue := r.Session.Get("uid")
		if sessionValue == nil {
			r.Response.RedirectTo("/")
			return
		}

		uidInt, _ := strconv.Atoi(sessionValue.(string))

		if uidInt <= 0 {
			r.Response.RedirectTo("/")
			return
		}

		// set context uid
		r.SetCtxVar("uid", sessionValue)
		// c.Set("uid", sessionValue)
		r.Middleware.Next()

		return
	}
}
