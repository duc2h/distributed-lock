package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

func main() {
	for i := 0; i < 10; i++ {
		incr()
	}
}

func incr() {
	var counterKey = "counter_key"

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "redis_pass",
		DB:       0,
	})

	// increase the value by one.
	getResp := client.Get(counterKey)
	value, err := getResp.Int64()
	if err == nil || err == redis.Nil {
		value++
		resp := client.Set(counterKey, value, 0)
		_, err := resp.Result()
		if err != nil {
			fmt.Println("set value error!")
		}
	} else {
		fmt.Println("Cannot get value from cache")
	}

	fmt.Println("current value is: ", value)

}
