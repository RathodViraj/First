package handler

import (
	chachingService "First/chachingservice"
	"First/model"
	"First/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	service           *service.PostService
	connectionService *service.ConnectionService
}

func NewPostHandler(service *service.PostService, connectionService *service.ConnectionService) *PostHandler {
	return &PostHandler{service, connectionService}
}

func (h *PostHandler) CreatePost(ctx *gin.Context) {
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

	var post model.Post
	if err := ctx.BindJSON(&post); err != nil {
		JSONError(ctx, http.StatusBadRequest, "Invalid request body")
		return
	}

	if post.Content == "" {
		JSONError(ctx, http.StatusBadRequest, "Content is required")
		return
	}

	post.Uid = userID

	if err := h.service.CreatePost(&post); err != nil {
		JSONError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	followers, err := h.connectionService.GetFollowers(userID)
	if err == nil {
		for _, user := range followers {
			chachingService.InvalidateUserFeedChahe(user.Id, ctx)
		}
	}

	ctx.IndentedJSON(http.StatusCreated, post)
}

func (h *PostHandler) DeletePost(ctx *gin.Context) {
	postID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, "Invalid post ID")
		return
	}

	userIDRaw, exists := ctx.Get("userID")
	if !exists {
		JSONError(ctx, http.StatusUnauthorized, "User not authenticated")
		return
	}
	userID := userIDRaw.(int)

	roleRaw, _ := ctx.Get("role")
	role, _ := roleRaw.(string)
	isAdmin := role == "admin"

	if err := h.service.DeletePost(postID, userID, isAdmin); err != nil {
		JSONError(ctx, http.StatusForbidden, err.Error())
		return
	}

	followers, err := h.connectionService.GetFollowers(userID)
	if err == nil {
		for _, user := range followers {
			chachingService.InvalidateUserFeedChahe(user.Id, ctx)
		}
	}

	ctx.Status(http.StatusNoContent)
}

func (h *PostHandler) GetPost(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, "Invalid post ID")
		return
	}

	post, err := h.service.GetPost(id)
	if err != nil {
		JSONError(ctx, http.StatusNotFound, err.Error())
		return
	}

	ctx.IndentedJSON(http.StatusOK, post)
}

func (h *PostHandler) RecentPosts(ctx *gin.Context) {
	page_ := ctx.DefaultQuery("page", "1")

	page, err := strconv.Atoi(page_)
	if err != nil || page < 1 {
		page = 1
	}

	offset := (page - 1) * 10
	posts, err := h.service.NewPosts(offset)
	if err != nil {
		JSONError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.IndentedJSON(http.StatusOK, posts)
}

func (h *PostHandler) UpdatePost(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
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

	var post model.Post
	if err := ctx.BindJSON(&post); err != nil {
		JSONError(ctx, http.StatusBadRequest, "Invalid request body")
		return
	}
	post.Id = id
	post.Uid = userID

	if err := h.service.UpdatePost(&post); err != nil {
		JSONError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	chachingService.InvalUserIDateUserProfileChahe(userID, ctx)
	ctx.Status(http.StatusOK)
}
