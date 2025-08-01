package handler

import (
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

type LikeHandler struct {
	service *service.LikeService
	Hub     *notification.Hub
	gr      repository.Graph
	nr      repository.NotificationRepository
}

func NewLikeHandler(service *service.LikeService, h *notification.Hub, gr repository.Graph, nr repository.NotificationRepository) *LikeHandler {
	return &LikeHandler{service, h, gr, nr}
}

func (h *LikeHandler) LikePost(ctx *gin.Context) {
	postID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, "Invalid post ID")
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

	like := model.Like{
		Pid: postID,
		Uid: userID,
	}

	if err := h.service.Like(&like); err != nil {
		JSONError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	//notifaction logic
	ownerID, _ := h.gr.GetUserIDByPostID(postID)
	if ownerID == -1 || ownerID == userID {
		ctx.Status(http.StatusCreated)
		return
	}

	notif := model.Notification{
		Type:      "like",
		FromUser:  userID,
		ToUser:    ownerID,
		PostID:    &postID,
		Message:   "Someone liked your post",
		Timestamp: time.Now().Unix(),
	}
	h.Hub.Broadcast <- notif
	log.Println("Sending notification via hub")

	_ = h.nr.SaveNotification(model.Notification{
		Type:      "like",
		FromUser:  userID,
		ToUser:    ownerID,
		PostID:    &postID,
		Message:   "Someone liked your post",
		Seen:      false,
		Timestamp: time.Now().Unix(),
	})

	ctx.Status(http.StatusCreated)
}

func (h *LikeHandler) UnlikePost(ctx *gin.Context) {
	postId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, "Invalid post ID")
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

	like := model.Like{
		Pid: postId,
		Uid: userID,
	}

	if err := h.service.Unlike(&like); err != nil {
		JSONError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusCreated)
}

func (h *LikeHandler) GetLikes(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	_, exists := ctx.Get("userID")
	if !exists {
		JSONError(ctx, http.StatusUnauthorized, "User not authenticated")
		return
	}

	users, err := h.service.GetLikes(id)
	if err != nil {
		JSONError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.IndentedJSON(http.StatusOK, users)
}
