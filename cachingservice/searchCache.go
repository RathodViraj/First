package cachingservice

import (
	"First/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func GetUserSearchCache(key string, ctx context.Context) []model.User {
	cacheKey := fmt.Sprintf("search:users:%s", key)

	var users []model.User
	data, err := RDB.Get(ctx, cacheKey).Result()
	if err == nil && data != "" {
		if err = json.Unmarshal([]byte(data), &users); err != nil {
			log.Println(err.Error())
			return []model.User{}
		}
	}

	return users
}

func CachedSearchUser(key string, users []model.User, ctx context.Context) {
	cacheKey := fmt.Sprintf("search:users:%s", key)
	if data, err := json.Marshal(users); err == nil {
		RDB.Set(ctx, cacheKey, data, time.Minute*2)
	}
}

func InvalidateSearchUserCache(key string, ctx context.Context) {
	cacheKey := fmt.Sprintf("search:users:%s", key)
	if err := RDB.Del(ctx, cacheKey).Err(); err != nil {
		log.Println(err.Error())
	}
}

func GetPostSearchCache(key string, ctx context.Context) []model.Post {
	cacheKey := fmt.Sprintf("search:posts:%s", key)

	var posts []model.Post
	data, err := RDB.Get(ctx, cacheKey).Result()
	if err == nil && data != "" {
		if err = json.Unmarshal([]byte(data), &posts); err != nil {
			log.Println(err.Error())
			return []model.Post{}
		}
	}

	return posts
}

func CachedSearchPost(key string, posts []model.Post, ctx context.Context) {
	cacheKey := fmt.Sprintf("search:posts:%s", key)
	if data, err := json.Marshal(posts); err == nil {
		RDB.Set(ctx, cacheKey, data, time.Minute*2)
	}
}

func InvalidateSearchPostCache(key string, ctx context.Context) {
	cacheKey := fmt.Sprintf("search:posts:%s", key)
	if err := RDB.Del(ctx, cacheKey).Err(); err != nil {
		log.Println(err.Error())
	}
}
