package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	// Add logger
	// "github.com/rs/zerolog/log"

	"encoding/json"
	"my_rest_server/client"
	e "my_rest_server/error"
	m "my_rest_server/model"
)

const (
	UserKeyPrefix = "user:"
	UserSetName   = "users"
)

type UserService struct {
	client        *redis.Client
	kafkaProducer *client.KafkaProducer
}

func NewUserService(
	addr string,
	password string,
	db int,
	kafkaProducer *client.KafkaProducer,
) (*UserService, error) {
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
	return &UserService{
		client:        client,
		kafkaProducer: kafkaProducer,
	}, nil
}

func (c *UserService) IsUserExists(ctx context.Context, userId string) (bool, error) {
	return c.client.HExists(ctx, userKey(userId), "data").Result()
}

func (c *UserService) SaveUser(ctx context.Context, user m.User) (string, error) {

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

	c.kafkaProducer.ProduceMessage(string(data))
	return user.Id, nil
}

func userKey(uuid string) string {
	return UserKeyPrefix + uuid
}

func (c *UserService) GetUser(ctx context.Context, userId string) (*m.User, error) {
	var user m.User

	data, err := c.client.HGet(ctx, userKey(userId), "data").Result()
	// 404
	if err == redis.Nil {
		return &user, e.NewError2(
			http.StatusNotFound,
			fmt.Sprintf("user: %s not found", userId),
			err,
		)
	}
	// 500
	if err != nil {
		return &user, fmt.Errorf("failed to get user: %s. error: %s", userId, err)
	}
	// 500
	err = json.Unmarshal([]byte(data), &user)
	if err != nil {
		return &user, fmt.Errorf("failed to parse user json: %s from redis. error: %s", data, err)
	}
	// 200
	return &user, nil
}

func (c *UserService) GetAllUsers(ctx context.Context) (*[]m.User, error) {
	// Get all user ids
	ids, err := c.client.SMembers(ctx, UserSetName).Result()

	if err != nil {
		return &[]m.User{}, fmt.Errorf("failed to get all user keys. error: %s", err)
	}
	if len(ids) == 0 {
		return &[]m.User{}, nil
	}
	// Result from redis
	resultCmds := make([]*redis.StringCmd, len(ids))

	// Open pipe in order to send all cmds by one request
	pipe := c.client.Pipeline()

	for i, id := range ids {
		resultCmds[i] = pipe.HGet(ctx, userKey(id), "data")
	}
	// Execute pipe
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users from redis. error: %s", err)
	}

	users := make([]m.User, 0, len(ids))
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

func (c *UserService) DeleteUser(ctx context.Context, userId string) (*m.User, error) {
	var user *m.User

	user, err := c.GetUser(ctx, userId)
	if err != nil {
		return user, fmt.Errorf("failed to find user: %s", userId)
	}
	userKey := userKey(userId)
	pipe := c.client.TxPipeline()

	pipe.Del(ctx, userKey)
	pipe.SRem(ctx, "users", userId)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return user, fmt.Errorf("failed to delete user: %s. error: %s", userId, err)
	}
	return user, nil
}

func (u *UserService) Close() error {
	err := u.client.Close()
	if err != nil {
		return fmt.Errorf("failed to close redis client. error: %s", err)
	}
	return nil
}
