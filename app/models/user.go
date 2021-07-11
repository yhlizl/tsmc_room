package models

import (
	"fmt"
	"time"
)

type User struct {
	ID        uint
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	AvatarId  string    `json:"avatar_id"`
	CreatedAt time.Time `time_format:"2006-01-02 15:04:05"`
	UpdatedAt time.Time `time_format:"2006-01-02 15:04:05"`
}

func AddUser(value interface{}) User {
	var u User
	fmt.Println(value.(map[string]interface{}))
	u.Username = value.(map[string]interface{})["username"].(string)
	u.Password = value.(map[string]interface{})["password"].(string)
	u.AvatarId = value.(map[string]interface{})["avatar_id"].(string)
	ChatDB.Model("users").Data(u).Insert()
	return u
}

func SaveAvatarId(AvatarId string, u *User) *User {
	u.AvatarId = AvatarId
	//ChatDB.Save("users",)

	ChatDB.Model("users").Save(*u)

	return u
}

func FindUserByField(field, value string) User {
	var u User

	if field == "id" || field == "username" {
		ChatDB.Model("users").Where(field+" = ?", value).Scan(&u)
	}

	return u
}

func GetOnlineUserList(uids []float64) []map[string]interface{} {
	var results []map[string]interface{}
	ChatDB.Model("users").Where("id IN ?", uids).Struct(&results)

	return results
}
