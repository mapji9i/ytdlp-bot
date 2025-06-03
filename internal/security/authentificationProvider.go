package security

import (
	"strconv"
	"strings"
	Environment "ytdlp-bot/internal/environment"
)

func CheckAuth(UserId int64) bool {
	var acceptedUsersId = []int64{}
	acceptedUsersId = append(acceptedUsersId, Environment.Environment.RootUserId)
	botUserId, _ := strconv.ParseInt(strings.Split(Environment.Environment.TelegramToken, ":")[0], 10, 64)
	acceptedUsersId = append(acceptedUsersId, botUserId)
	for _, acceptedUserId := range acceptedUsersId {
		if UserId == acceptedUserId {
			return true
		}

	}

	return false
}
