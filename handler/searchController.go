package handler

import (
	"First/cachingservice"
	"First/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	service *service.SearchService
}

func NewSearchHandler(s *service.SearchService) *SearchHandler {
	return &SearchHandler{s}
}

func (h SearchHandler) SearchUser(ctx *gin.Context) {
	key := ctx.Query("search")

	if users := cachingservice.GetUserSearchCache(key, ctx); len(users) > 0 {
		ctx.IndentedJSON(http.StatusOK, users)
		return
	}

	users, err := h.service.SerachUser(key)
	if err != nil {
		JSONError(ctx, http.StatusInternalServerError, "Couldn't serach")
		return
	}

	if len(users) == 0 {
		msg := fmt.Sprintf("No user found with %s", key)
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": msg})
		return
	}

	cachingservice.CachedSearchUser(key, users, ctx)

	ctx.IndentedJSON(http.StatusOK, users)
}

func (h SearchHandler) SearchPost(ctx *gin.Context) {
	key := ctx.Query("search")

	if posts := cachingservice.GetPostSearchCache(key, ctx); len(posts) > 0 {
		ctx.IndentedJSON(http.StatusOK, posts)
		return
	}

	posts, err := h.service.SerachPost(key)
	if err != nil {
		JSONError(ctx, http.StatusInternalServerError, "Couldn't serach")
		return
	}

	if len(posts) == 0 {
		msg := fmt.Sprintf("No post found with %s", key)
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": msg})
		return
	}

	cachingservice.CachedSearchPost(key, posts, ctx)

	ctx.IndentedJSON(http.StatusOK, posts)
}
