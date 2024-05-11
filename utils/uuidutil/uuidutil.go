package uuidutil

import (
	"github.com/google/uuid"
	"math/rand"
	"strconv"
	"time"
)

// 获取length长度的随机序列
func randStr(length int) string {
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := []byte(str)
	var result []byte
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

/*
UUID ...
基于时间戳的uuid生成

	@param:
		length：uuid长度
	@return:
		length <= 7 ：""
		length >  7 : length长度uuid (length-7随机符号:length)
*/
func UUID(length int) string {
	if length <= 7 {
		return ""
	}
	timeStamp := time.Now().Unix()
	formatTimeStr := strconv.Itoa(int(timeStamp))
	seq := randStr(length - 7)
	return seq + formatTimeStr[3:10]
}

func UUID4(length int) string {
	// 生成 UUID4
	u := uuid.New()

	// 去掉 "-"
	s := u.String()
	s = s[:8] + s[9:13] + s[14:18] + s[19:23] + s[24:]

	// 获取前 length 位
	if len(s) < length {
		length = len(s)
	}
	s = s[:length]
	return s
}
