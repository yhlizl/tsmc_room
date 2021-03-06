package message_service

import (
	"go-gofram-chat/app/models"
)

func GetLimitMsg(roomId string, offset int) []map[string]interface{} {
	return models.GetLimitMsg(roomId, offset)
}

func GetLimitPrivateMsg(uid, toUId string, offset int) []map[string]interface{} {
	return models.GetLimitPrivateMsg(uid, toUId, offset)
}
