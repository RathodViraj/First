package handler

import (
	chachingservice "First/chachingservice"
	"First/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ConnectionHandler struct {
	service *service.ConnectionService
}

func NewConnectionHandler(service *service.ConnectionService) *ConnectionHandler {
	return &ConnectionHandler{service: service}
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

	if err := h.service.Follow(followerID, followingID); err != nil {
		JSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	chachingservice.InvalidateUserFollowersCache(ctx, followingID)
	chachingservice.InvalidateUserFollowingsCache(ctx, followerID)
	chachingservice.InvalidatemutualCache(ctx, followerID)

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

	if err := h.service.Unfollow(followerID, followingID); err != nil {
		JSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	chachingservice.InvalidateUserFollowersCache(ctx, followingID)
	chachingservice.InvalidateUserFollowingsCache(ctx, followerID)
	chachingservice.InvalidatemutualCache(ctx, followerID)

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

	if followers := chachingservice.GetFollowersChached(userID, ctx); len(followers) > 0 {
		ctx.IndentedJSON(http.StatusOK, followers)
		return
	}

	followers, err := h.service.GetFollowers(userID)
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	chachingservice.ChachedFollowers(userID, followers, ctx)

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

	if followings := chachingservice.GetFollowingsChached(userID, ctx); len(followings) > 0 {
		ctx.IndentedJSON(http.StatusOK, followings)
		return
	}

	followings, err := h.service.GetFollowings(userID)
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	chachingservice.ChachedFollowings(userID, followings, ctx)

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

	if mutual := chachingservice.GetFollowingsChached(userID, ctx); len(mutual) > 0 {
		ctx.IndentedJSON(http.StatusOK, mutual)
		return
	}

	mutual, err := h.service.Mutual(id)
	if err != nil {
		JSONError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	chachingservice.ChachedMutuals(userID, mutual, ctx)

	ctx.IndentedJSON(http.StatusOK, mutual)
}
