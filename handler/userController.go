package handler

import (
	chachingservice "First/chachingService"
	"First/model"
	"First/service"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service     *service.UserService
	postService *service.PostService
}

func NewUserHandler(service *service.UserService, postService *service.PostService) *UserHandler {
	return &UserHandler{service: service, postService: postService}
}

func (h *UserHandler) CreateUser(ctx *gin.Context) {
	var user model.User
	if err := ctx.BindJSON(&user); err != nil {
		JSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if strings.TrimSpace(user.Email) == "" || strings.TrimSpace(user.Password) == "" {
		JSONError(ctx, http.StatusBadRequest, "Email and password are required")
		return
	}

	if err := h.service.RegisterUser(&user); err != nil {
		JSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	user.Password = ""

	if user.Role == "" {
		user.Role = "user"
	}

	ctx.IndentedJSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	log.Println(id)
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, posts := chachingservice.GetChachedUserProfile(id, ctx)
	if user == nil {
		user, err = h.service.GetUser(id)
		if err != nil {
			JSONError(ctx, http.StatusNotFound, "User dose not exits!")
			return
		}
	}
	user.Password = ""

	if posts == nil {
		posts, err = h.postService.UserPosts(id)
		if err != nil {
			msg := "can't fetch user's posts"
			ctx.IndentedJSON(http.StatusOK, gin.H{
				"user":  user,
				"posts": msg,
			})

			return
		}
	}

	chachingservice.ChachedUserProfile(user.Id, ctx, user, posts)

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"user":  user,
		"posts": posts,
	})

}

func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, err.Error())
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
	roleRaw, _ := ctx.Get("role")
	role, _ := roleRaw.(string)
	isAdmin := role == "admin"

	if err := h.service.DeleteUser(id, userID, isAdmin); err != nil {
		JSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	chachingservice.InvalUserIDateUserProfileChahe(userID, ctx)

	ctx.Status(http.StatusOK)
}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var user model.User
	if err := ctx.BindJSON(&user); err != nil {
		JSONError(ctx, http.StatusBadRequest, "Invalid request body")
		return
	}
	user.Id = id
	user.Password = "" // Do not allow password change

	if err := h.service.UpdateUser(&user); err != nil {
		JSONError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	chachingservice.InvalUserIDateUserProfileChahe(id, ctx)
	ctx.Status(http.StatusOK)
}

func (h *UserHandler) GetFeed(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, err.Error())
	}

	uidVal, exists := ctx.Get("userID")
	if !exists {
		JSONError(ctx, http.StatusUnauthorized, "Unauthorized user")
		return
	}
	id, ok := uidVal.(int)
	if !ok {
		JSONError(ctx, http.StatusInternalServerError, "Failed to extract user ID")
		return
	}
	userID = id

	page_ := ctx.Query("page")
	page, err := strconv.Atoi(page_)
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, "Invalid limit")
		return
	}

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * 10

	if feed := chachingservice.GetChachedUserFeed(userID, page, ctx); feed != nil {
		ctx.IndentedJSON(http.StatusOK, feed)
		return
	}

	feed, err := h.service.UserFeed(userID, offset)
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	chachingservice.ChachedUserFeed(userID, page, ctx, feed)

	ctx.IndentedJSON(http.StatusOK, feed)
}
