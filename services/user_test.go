package services

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUsersName(t *testing.T) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root",
		"openIM123",
		"34.192.77.144:13306",
		"testdb")
	InitMysql(dsn)

	uids := []int64{6450102772, 6374129665}
	users := GetUserNames(uids)
	assert.Equal(t, len(users), len(uids))
	assert.Equal(t, "goat2023_bot", users[6374129665].Username)
}
