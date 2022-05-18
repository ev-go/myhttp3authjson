package My_redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type MyRedis struct {
	Value1 string
	Value2 string
}

var ctx = context.Background()

func Main() {
	//redis start
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	userid := "user#125"
	currentusertoken := "testtoken"

	node := rdb.Set(ctx, userid, currentusertoken, 0).Err()
	if node != nil {
		panic(node)
	}

	val, node := rdb.Get(ctx, userid).Result()
	if node == redis.Nil {
		fmt.Println("key1 does not exist")
	} else if node != nil {
		panic(node)
	} else {
		fmt.Println(userid, val)
	}

	val2, node := rdb.Get(ctx, "key2").Result()
	if node == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if node != nil {
		panic(node)
	} else {
		fmt.Println("key2", val2)
	}

	//var m = Message{"World", "Hello", Name{"Dmitry", "Victorovich"}, "79082706690", "393181839", 211}
	// Output: key value
	// key2 does not exist

	//redis end
}
