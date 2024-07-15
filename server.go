package main

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	leaderKey := "leaderPort"
	port := env("PORT", "8001")

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     env("REDIS_ADDR", "localhost:6379"),
		Password: env("REDIS_PASSWORD", ""),
		DB:       0,
	})

	// Ensure connection
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis connection failed: %v\n", err)
	}

	// Maintain leader: goroutine that continuously checks every 1-10s and maintains a leader value in Redis,
	// setting it to serviceId if it is not already set.
	go func() {
		for {
			leaderValue, err := redisClient.Get(ctx, leaderKey).Result()
			if err == redis.Nil {
				log.Printf("set %s to %s\n", leaderKey, port)
				redisClient.Set(ctx, leaderKey, port, time.Second*10)
			} else if err != nil {
				log.Fatalln(err)
			} else {
				log.Printf("leader is %s\n", leaderValue)
			}
			time.Sleep(time.Second * time.Duration(rand.IntN(10)+1))
		}
	}()

	// Setup server
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Setup handler to respond with what port the leader is running on.
	r.GET("/", func(c *gin.Context) {
		leaderValue, err := redisClient.Get(c.Request.Context(), leaderKey).Result()
		if err == redis.Nil || err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{leaderKey: leaderValue})
	})

	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalln(err)
	}
}

func env(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
