package redis

// import (
// 	"context"
// 	"os"
// 	"strconv"

// 	"github.com/redis/go-redis/v9"

// 	log "github.com/alexshv/file-storage/logger"
// )

// type RedisClient struct {
// 	client *redis.Client
// 	ctx    context.Context
// }

// var redisClient *RedisClient

// func Init(ctx context.Context) {
// 	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))

// 	if err != nil {
// 		log.GetLogger().WithField("message", err.Error()).Info("redis.connect.error")
// 		os.Exit(1)
// 	}

// 	redisClient = &RedisClient{
// 		client: redis.NewClient(&redis.Options{
// 			Addr: os.Getenv("REDIS_DSN"),
// 			DB:   db,
// 		}),
// 		ctx: ctx,
// 	}

// 	log.GetLogger().Info("redis.connected")
// }

// func Shutdown() {
// 	client := redisClient.client

// 	if client == nil {
// 		return
// 	}

// 	log.GetLogger().Info("redis.shutdown")

// 	client.Shutdown(redisClient.ctx)
// 	client.Close()
// }
