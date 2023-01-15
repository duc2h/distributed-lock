package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

func main() {

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second * 3)
		wg.Add(1)
		go func() {
			defer wg.Done()
			incr()
		}()
	}

	wg.Wait()

}

func incr() {
	var lockKey = "counter_lock"
	var counterKey = "counter_key"

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "redis_pass",
		DB:       0,
	})

	// make a lock
	resp := client.SetNX(lockKey, 1, time.Second*5)
	lockSuccess, err := resp.Result()

	if err != nil || !lockSuccess {
		fmt.Println(err, "lock result: ", lockSuccess)
		return
	}

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
	}

	fmt.Println("current value is: ", value)

	// release lock
	delResp := client.Del(lockKey)
	releaseResp, err := delResp.Result()
	if err == nil && releaseResp > 0 {
		fmt.Println("Release lock success!")
	} else {
		fmt.Println("Release lock failed: ", err.Error())
	}

}
