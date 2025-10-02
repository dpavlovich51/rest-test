package storage

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	// Add logger
	// "github.com/rs/zerolog/log"

	"encoding/json"
	m "my_rest_server/model"
)

const (
	UserKeyPrefix = "user:"
	UserSetName   = "users"
)

type Cache struct {
	client redis.Client
}

func NewClient(addr string, password string, db int) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	ctx := context.Background()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis. %s", err)
	}
	return &Cache{
		client: *client,
	}, nil
}

func (c *Cache) SaveUser(ctx context.Context, user m.User) (string, error) {

	// to json
	data, err := json.Marshal(user)
	if err != nil {
		return "", fmt.Errorf("failed map user: %s to json. Error: %v", user, err)
	}

	key := userKey(user.Id)
	// Open transaction
	pipe := c.client.TxPipeline()

	// Save user
	pipe.HSet(ctx, key, "data", data)
	// Save user id
	pipe.SAdd(ctx, UserSetName, user.Id)

	// Commit transaction
	_, err = pipe.Exec(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to save user: %s. Error: %v", user, err)
	}
	return user.Id, nil
}

func userKey(uuid string) string {
	return UserKeyPrefix + uuid
}

func (c *Cache) GetUser(ctx context.Context, userId string) (m.User, error) {
	var user m.User

	data, err := c.client.HGet(ctx, userKey(userId), "data").Result()
	// 404
	if err == redis.Nil {
		return user, fmt.Errorf("user: %s not found. error: ", err)
	}
	// 500
	if err != nil {
		return user, fmt.Errorf("failed to get user: %s. error: ", err)
	}
	// 500
	err = json.Unmarshal([]byte(data), &user)
	if err != nil {
		return user, fmt.Errorf("failed to parse user json: %s from redis. error: %s", data, err)
	}
	// 200
	return user, nil
}

func (c *Cache) GetAllUsers(ctx context.Context) (*[]m.User, error) {
	// Get all user keys
	keys, err := c.client.SMembers(ctx, UserSetName).Result()
	log.Info().Msg(fmt.Sprintf("all keys: %s", keys))

	if err != nil {
		return &[]m.User{}, fmt.Errorf("failed to get all user keys. error: %s", err)
	}
	if len(keys) == 0 {
		return &[]m.User{}, nil
	}
	// Result from redis
	resultCmds := make([]*redis.StringCmd, len(keys))
	// Open pipe in order to send all cmds by one request
	pipe := c.client.Pipeline()

	for i, key := range keys {
		resultCmds[i] = pipe.Get(ctx, userKey(key))
	}
	// Execute pipe
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users from redis. error: %s", err)
	}

	users := make([]m.User, 0, len(keys))
	for _, cmd := range resultCmds {
		jsonData, err := cmd.Result()
		if err != nil {
			log.Warn().Msg(fmt.Sprintf("failed to get user data. error: %s", err))
			continue
		}
		var user m.User
		err = json.Unmarshal([]byte(jsonData), &user)
		if err != nil {
			log.Warn().Msg(fmt.Sprintf("failed to map json data: %s to user struct. error: %s", jsonData, err))
			continue
		}
		users = append(users, user)
	}
	return &users, nil
}

func (c *Cache) DeleteUser(ctx context.Context, userId string) (m.User, error) {
	var user m.User

	user, err := c.GetUser(ctx, userId)
	if err != nil {
		return user, fmt.Errorf("failed to find user: %s", userId)
	}
	userKey := userKey(userId)
	pipe := c.client.TxPipeline()

	pipe.Del(ctx, userKey)
	pipe.SRem(ctx, "users", userKey)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return user, fmt.Errorf("failed to delete user: %s. error: %s", userId, err)
	}
	return user, nil
}
