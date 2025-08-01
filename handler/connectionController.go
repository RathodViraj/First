package handler

import (
	cachingservice "First/cachingservice"
	"First/model"
	"First/notification"
	"First/repository"
	"First/service"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ConnectionHandler struct {
	service *service.ConnectionService
	hub     *notification.Hub
	gr      repository.Graph
	nr      repository.NotificationRepository
}

func NewConnectionHandler(service *service.ConnectionService, hub *notification.Hub, gr repository.Graph, nr repository.NotificationRepository) *ConnectionHandler {
	return &ConnectionHandler{service, hub, gr, nr}
}

func (h *ConnectionHandler) FollowUser(ctx *gin.Context) {
	followerID, err := strconv.Atoi(ctx.Param("follower_id"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	followingID, err := strconv.Atoi(ctx.Param("following_id"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	uidVal, exists := ctx.Get("userID")
	if !exists {
		JSONError(ctx, http.StatusUnauthorized, "User not authenticated")
		return
	}
	userID, ok := uidVal.(int)
	if !ok {
		JSONError(ctx, http.StatusInternalServerError, "Failed to extract user ID")
		return
	}
	followerID = userID

	if err := h.service.Follow(followerID, followingID); err != nil {
		JSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	cachingservice.InvalidateUserFollowersCache(ctx, followingID)
	cachingservice.InvalidateUserFollowingsCache(ctx, followerID)
	cachingservice.InvalidatemutualCache(ctx, followerID)

	if followingID != followerID {
		notif := model.Notification{
			Type:      "Follow",
			FromUser:  followerID,
			ToUser:    followingID,
			PostID:    nil,
			Message:   "Someone start following you",
			Timestamp: time.Now().Unix(),
		}

		h.hub.Broadcast <- notif
		log.Println("Sending notification via hub")
	}

	_ = h.nr.SaveNotification(model.Notification{
		Type:      "follow",
		FromUser:  followerID,
		ToUser:    followingID,
		PostID:    nil,
		Message:   "Someone start following you",
		Seen:      false,
		Timestamp: time.Now().Unix(),
	})

	ctx.Status(http.StatusCreated)
}

func (h *ConnectionHandler) UnfollowUser(ctx *gin.Context) {
	followerID, err := strconv.Atoi(ctx.Param("follower_id"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	followingID, err := strconv.Atoi(ctx.Param("following_id"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	uidVal, exists := ctx.Get("userID")
	if !exists {
		JSONError(ctx, http.StatusUnauthorized, "User not authenticated")
		return
	}
	userID, ok := uidVal.(int)
	if !ok {
		JSONError(ctx, http.StatusInternalServerError, "Failed to extract user ID")
		return
	}
	followerID = userID

	if err := h.service.Unfollow(followerID, followingID); err != nil {
		JSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	cachingservice.InvalidateUserFollowersCache(ctx, followingID)
	cachingservice.InvalidateUserFollowingsCache(ctx, followerID)
	cachingservice.InvalidatemutualCache(ctx, followerID)

	ctx.Status(http.StatusOK)
}

func (h *ConnectionHandler) GetFollowers(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	uidVal, exists := ctx.Get("userID")
	if !exists {
		JSONError(ctx, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userID, ok := uidVal.(int)
	if !ok {
		JSONError(ctx, http.StatusInternalServerError, "Failed to extract user ID")
		return
	}

	if followers := cachingservice.GetFollowersCached(userID, ctx); len(followers) > 0 {
		ctx.IndentedJSON(http.StatusOK, followers)
		return
	}

	followers, err := h.service.GetFollowers(userID)
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	cachingservice.CachedFollowers(userID, followers, ctx)

	ctx.IndentedJSON(http.StatusOK, followers)
}

func (h *ConnectionHandler) GetFollowings(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	uidVal, exists := ctx.Get("userID")
	if !exists {
		JSONError(ctx, http.StatusUnauthorized, "User not authenticated")
		return
	}
	userID, ok := uidVal.(int)
	if !ok {
		JSONError(ctx, http.StatusInternalServerError, "Failed to extract user ID")
		return
	}

	if followings := cachingservice.GetFollowingsCached(userID, ctx); len(followings) > 0 {
		ctx.IndentedJSON(http.StatusOK, followings)
		return
	}

	followings, err := h.service.GetFollowings(userID)
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	cachingservice.CachedFollowings(userID, followings, ctx)

	ctx.IndentedJSON(http.StatusOK, followings)
}

func (h *ConnectionHandler) GetMutual(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, "Invalid id")
	}

	uidVal, exists := ctx.Get("userID")
	if !exists {
		JSONError(ctx, http.StatusUnauthorized, "User not authenticated")
		return
	}
	userID, ok := uidVal.(int)
	if !ok {
		JSONError(ctx, http.StatusInternalServerError, "Failed to extract user ID")
		return
	}
	id = userID

	if mutual := cachingservice.GetFollowingsCached(userID, ctx); len(mutual) > 0 {
		ctx.IndentedJSON(http.StatusOK, mutual)
		return
	}

	mutual, err := h.service.Mutual(id)
	if err != nil {
		JSONError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	cachingservice.CachedMutuals(userID, mutual, ctx)

	ctx.IndentedJSON(http.StatusOK, mutual)
}
