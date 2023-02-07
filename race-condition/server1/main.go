package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

func main() {
	// incr()
	incr2()
	// incr3()
}

// this function will make race condition with server2.
// Server2 will increase the value by 10, then the value is replaced by server1 -> race condition.
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
		time.Sleep(time.Millisecond * 4800)
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

// this function will make race condition with server2.
// Server3 will increase the value by one, then the value is replaced by server1 -> race condition.
// incr3 will fix this problem.
func incr2() {
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
		// release lock due to timeout.
		time.Sleep(time.Millisecond * 8000)
		value++
		// Replacing the same value into `counterKey's value` is changed by server3 -> race-condition.
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
		fmt.Printf("Release lock failed: %s, releaseResp: %d", err, releaseResp)
	}

}

// we need to check: is lock still exist in this process? if yes -> change the value, no -> don't change.
func incr3() {
	var lockKey = "counter_lock"
	var counterKey = "counter_key"

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "redis_pass",
		DB:       0,
	})

	// make a lock
	// add server-name for lockey
	resp := client.SetNX(lockKey, "server1", time.Second*5)
	lockSuccess, err := resp.Result()

	if err != nil || !lockSuccess {
		fmt.Println(err, "lock result: ", lockSuccess)
		return
	}

	// increase the value by one.
	getResp := client.Get(counterKey)
	value, err := getResp.Int64()
	if err == nil || err == redis.Nil {
		// release lock due to timeout.
		time.Sleep(time.Millisecond * 8000)
		value++
		// Replacing the same value into `counterKey's value` is changed by server3 -> race-condition.
		lockResp := client.Get(lockKey)
		lockValue := lockResp.Val()
		if lockValue != "server1" {
			fmt.Println("lock is taken by other process.")
			return
		}
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
		fmt.Printf("Release lock failed: %s, releaseResp: %d", err, releaseResp)
	}

}
