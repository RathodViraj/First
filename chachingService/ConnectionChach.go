package chachingservice

import (
	"First/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func GetFollowersChached(userID int, ctx context.Context) []model.User {
	key := fmt.Sprintf("user:followers:%d", userID)
	data, err := RDB.Get(ctx, key).Result()
	if err == nil && data != "" {
		var followers []model.User
		if err = json.Unmarshal([]byte(data), &followers); err == nil {
			return followers
		}
	}

	return []model.User{}
}

func ChachedFollowers(userID int, followers []model.User, ctx context.Context) {
	key := fmt.Sprintf("user:followers:%d", userID)
	if data, err := json.Marshal(followers); err == nil {
		RDB.Set(ctx, key, data, time.Minute*2)
	}
}

func InvalidateUserFollowersCache(ctx context.Context, userId int) {
	key := fmt.Sprintf("user:followers:%d", userId)

	if err := RDB.Del(ctx, key).Err(); err != nil {
		log.Println(err.Error())
	}
}

func GetFollowingsChached(userID int, ctx context.Context) []model.User {
	key := fmt.Sprintf("user:followings:%d", userID)
	data, err := RDB.Get(ctx, key).Result()
	if err == nil && data != "" {
		var followings []model.User
		if err = json.Unmarshal([]byte(data), &followings); err == nil {
			return followings
		}
	}

	return []model.User{}
}

func ChachedFollowings(userID int, followings []model.User, ctx context.Context) {
	key := fmt.Sprintf("user:followings:%d", userID)
	if data, err := json.Marshal(followings); err == nil {
		RDB.Set(ctx, key, data, time.Minute*2)
	}
}

func InvalidateUserFollowingsCache(ctx context.Context, userId int) {
	key := fmt.Sprintf("user:followings:%d", userId)

	if err := RDB.Del(ctx, key).Err(); err != nil {
		log.Println(err.Error())
	}
}

func GetChahcedmutual(userID int, ctx context.Context) []model.User {
	key := fmt.Sprintf("user:mutual:%d", userID)
	data, err := RDB.Get(ctx, key).Result()
	if err == nil && data != "" {
		var mutuals []model.User
		if err = json.Unmarshal([]byte(data), &mutuals); err == nil {
			return mutuals
		}
	}

	return []model.User{}
}

func ChachedMutuals(userID int, mutuals []model.User, ctx context.Context) {
	key := fmt.Sprintf("user:mutual:%d", userID)
	if data, err := json.Marshal(mutuals); err == nil {
		RDB.Set(ctx, key, data, time.Minute*2)
	}
}

func InvalidatemutualCache(ctx context.Context, userId int) {
	key := fmt.Sprintf("user:mutual:%d", userId)

	if err := RDB.Del(ctx, key).Err(); err != nil {
		log.Println(err.Error())
	}
}
