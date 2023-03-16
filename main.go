package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Error connecting to Redis: ", err)
		return
	}

	fmt.Println("Connected to Redis: ", pong)

	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			fmt.Println("Error closing the connection no Redis: ", err)
		}
	}(client)
}
