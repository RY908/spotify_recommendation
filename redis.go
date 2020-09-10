package main

import(
    "github.com/gomodule/redigo/redis"
)

// Connection
func Connection() redis.Conn {
    const Addr = "127.0.0.1:6379"

    c, err := redis.Dial("tcp", Addr)
    if err != nil {
        panic(err)
    }
    return c
}

func Set(key, value string, c redis.Conn) string{
    res, err := redis.String(c.Do("SET", key, value))
    if err != nil {
        panic(err)
    }
    return res
}

func Get(key string, c redis.Conn) string {
    res, err := redis.String(c.Do("GET", key))
    if err != nil {
        panic(err)
    }
    return res
}