package chachingService

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

func GetChahcedMutal(userID int, ctx context.Context) []model.User {
	key := fmt.Sprintf("user:Mutal:%d", userID)
	data, err := RDB.Get(ctx, key).Result()
	if err == nil && data != "" {
		var mutals []model.User
		if err = json.Unmarshal([]byte(data), &mutals); err == nil {
			return mutals
		}
	}

	return []model.User{}
}

func ChachedMutals(userID int, mutals []model.User, ctx context.Context) {
	key := fmt.Sprintf("user:mutal:%d", userID)
	if data, err := json.Marshal(mutals); err == nil {
		RDB.Set(ctx, key, data, time.Minute*2)
	}
}

func InvalidateMutalCache(ctx context.Context, userId int) {
	key := fmt.Sprintf("user:mutal:%d", userId)

	if err := RDB.Del(ctx, key).Err(); err != nil {
		log.Println(err.Error())
	}
}
