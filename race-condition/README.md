Note:

We have 2 cases make race-condition happen when we using distributed-lock in redis.

1. Race-condition between server1 and server2.
- In order to test this case:
  - Start Redis
  - Start server1
  - start server2 immediately

- Result:
Server2
```
ducha@ucs-MacBook-Pro server2 % go run main.go     
current value is:  1
current value is:  2
current value is:  3
current value is:  4
current value is:  5
current value is:  6
current value is:  7
current value is:  8
current value is:  9
current value is:  10
ducha@ucs-MacBook-Pro server2 % 
```

Server1
```
ducha@ucs-MacBook-Pro server1 % go run main.go
current value is:  1
Release lock success!
ducha@ucs-MacBook-Pro server1 % 
```

=> So current value is 1, not 10. The value is replaced from 10 to 1 by server1, because server1 takes longer to update the value. -> race-condition.

2. Race-condition between server1 and server3.
- In order to test this case:
  - Start Redis
  - Start server1
  - start server3 immediately
- Result:
Server3:
```
ducha@ucs-MacBook-Pro server3 % go run main.go
current value is:  9
Release lock success!
ducha@ucs-MacBook-Pro server3 % 
```
Server1:
```
ducha@ucs-MacBook-Pro server1 % go run main.go
current value is:  9
Release lock failed: %!s(<nil>), releaseResp: 0%                                                      
ducha@ucs-MacBook-Pro server1 % 
```

As you can see, server1 cannot release lock due to timeout(over TTL). Althought we add lockKey's value for 2 servers, but at this time server3 checks the key is empty -> server3 takes the lock, server1 don't know it is timeout then it changes the value as well.


=====> The solution for 2 cases above is: we will add value for `lockKey`, the lockKey's value will relevant with server's name (unique). So we will check the lockKey's value in the logic like `server3's incr()` and `server1's incr3()`.

