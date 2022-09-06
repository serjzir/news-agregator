package api

import (
	"github.com/gin-gonic/gin"
	"github.com/serjzir/news-agregator/pkg/storage"
	"net/http"
	"strconv"
)

type API struct {
	posts storage.Repository
	r     *gin.Engine
}

// New конструктор API
func New(db *storage.Repository) *API {
	router := gin.Default()
	return &API{posts: *db, r: router}
}

// GetRouter конструктор API - содает router и запускает веб сервис.
func (api *API) GetRouter(ip, port string) *gin.Engine {
	api.r.GET("/news/:id", api.getNews)
	api.r.Run(ip + ":" + port)
	return api.r
}

func (api *API) getNews(c *gin.Context) {
	strID := c.Param("id")
	id, _ := strconv.Atoi(strID)
	posts, err := api.posts.News(c, id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
	}
	c.IndentedJSON(http.StatusOK, posts)
}
