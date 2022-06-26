package redis

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

var redisConn *redis.Client

func init() {
	redisConn = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func NewSession(ctx *gin.Context, cookieKey, redisValue string) {
	b := make([]byte, 64)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}
	newRedisKey := base64.URLEncoding.EncodeToString(b)

	if err := redisConn.Set(ctx, newRedisKey, redisValue, 0).Err(); err != nil {
		panic("Session登録時にエラーが発生：" + err.Error())
	}
	ctx.SetCookie(cookieKey, newRedisKey, 0, "/", "localhost", false, false)
}

func GetSession(ctx *gin.Context, cookieKey string) (string, error) {
	redisKey, _ := ctx.Cookie(cookieKey)
	redisValue, err := redisConn.Get(ctx, redisKey).Result()
	switch {
	case err == redis.Nil:
		fmt.Println("SessionKeyが登録されていません。")
		return "", nil
	case err != nil:
		fmt.Println("Session取得時にエラー発生：" + err.Error())
		return "", nil
	}
	return redisValue, nil
}

func DeleteSession(ctx *gin.Context, cookieKey string) {
	redisKey, _ := ctx.Cookie(cookieKey)
	if err := redisConn.Del(ctx, redisKey).Err(); err != nil {
		panic("Session削除時にエラーが発生：" + err.Error())
	}
	ctx.SetCookie(cookieKey, "", -1, "/", "localhost", false, false)
}
