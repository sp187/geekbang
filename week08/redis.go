package week08

import (
	"errors"
	"github.com/go-redis/redis"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

var (
	address = "localhost:6379"
	client  *redis.Client
)

func InitClient() {
	client = redis.NewClient(&redis.Options{
		Addr: address,
	})
	// test connection
	_, err := client.Ping().Result()
	if err != nil {
		panic("create redis client fail")
	}
}

func SetKey(key string, value interface{}) error {
	return client.Set(key, value, 600*time.Second).Err()
}

func GetKey(key string) (interface{}, error) {
	return client.Get(key).Result()
}

var memReg = regexp.MustCompile(`used_memory:(\d+)`)

func GetMemInfo() (int64, error) {
	mem := client.Info("memory").String()
	result := memReg.FindStringSubmatch(mem)
	if len(result) > 0 {
		return strconv.ParseInt(result[1], 10, 64)
	}
	return 0, errors.New("unexpected error")
}

func FlushAll() error {
	return client.FlushAll().Err()
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func randomStr(length int) string {
	b := make([]byte, 0, length)
	rand.Seed(time.Now().Unix())
	for i := 0; i < length; i++ {
		switch rand.Intn(3) {
		case 0:
			b = append(b, byte(randomInt(48, 57)))
		case 1:
			b = append(b, byte(randomInt(65, 90)))
		case 2:
			b = append(b, byte(randomInt(97, 122)))

		}
	}
	return string(b)
}
