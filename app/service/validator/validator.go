package validator

type User struct {
	Username string `v:"required|length:2,16#帳號不能為空|帳號長度應當在2到16之間"`
	Password string `v:"required|length:6,16#密碼不能為空|密碼長度應當在6到16之間"`
	AvatarId string `v:"required|integer#請選擇大頭貼"`
}
