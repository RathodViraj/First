package chachingService

import (
	"First/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func GetChachedUserProfile(userID int, ctx context.Context) (*model.User, *[]model.Post) {
	var (
		user  *model.User
		posts *[]model.Post
	)

	ukey := fmt.Sprintf("user:%d", userID)
	udata, err := RDB.Get(ctx, ukey).Result()
	if err == nil && udata != "" {
		if err = json.Unmarshal([]byte(udata), user); err != nil {
			user = nil
			log.Println(err.Error())
		}
	}

	pkey := fmt.Sprintf("user:post:%d", userID)
	pdata, err := RDB.Get(ctx, pkey).Result()
	if err == nil && udata != "" {
		if err = json.Unmarshal([]byte(pdata), posts); err != nil {
			posts = nil
			log.Println(err.Error())
		}
	}

	return user, posts
}

func ChachedUserProfile(UserID int, ctx context.Context, user *model.User, posts *[]model.Post) {
	ukey := fmt.Sprintf("user:%d", UserID)
	if data, err := json.Marshal(user); err == nil {
		RDB.Set(ctx, ukey, data, time.Minute*2)
	}
	pkey := fmt.Sprintf("user:post:%d", UserID)
	if data, err := json.Marshal(posts); err == nil {
		RDB.Set(ctx, pkey, data, time.Minute*2)
	}
}

func InvalUserIDateUserProfileChahe(UserID int, ctx context.Context) {
	ukey := fmt.Sprintf("user:*:%d", UserID)
	if err := RDB.Del(ctx, ukey).Err(); err != nil {
		log.Println(err.Error())
	}
}

func GetChachedUserFeed(userID, page int, ctx context.Context) *[]model.Post {
	var feed *[]model.Post

	key := fmt.Sprintf("home:user:%d:%d", userID, page)
	data, err := RDB.Get(ctx, key).Result()
	if err == nil && data != "" {
		if err = json.Unmarshal([]byte(data), &feed); err != nil {
			feed = nil
			log.Println(err.Error())
		}
	}

	return feed
}

func ChachedUserFeed(userID, page int, ctx context.Context, posts []model.Post) {
	key := fmt.Sprintf("user:home:%d:%d", userID, page)
	if data, err := json.Marshal(posts); err == nil {
		RDB.Set(ctx, key, data, time.Minute*2)
	}
}

func InvalidateUserFeedChahe(userID int, ctx context.Context) {
	ukey := fmt.Sprintf("user:home:%d:*", userID)
	if err := RDB.Del(ctx, ukey).Err(); err != nil {
		log.Println(err.Error())
	}
}
