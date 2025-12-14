package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var r = gin.Default()

func main() {
	group := r.Group("/users")
	{
		// 查询列表
		group.GET("/queryUsers", queryUserList)
		group.POST("/create", create)
		group.GET("/getUser", getUser)
		group.GET("/deleteByName", deleteByName)
	}
	err := r.Run(":8083")
	if err != nil {
		fmt.Println(err)
	}
}

func getUser(c *gin.Context) {
	name := c.Query("name")
	//var user UserInfo
	//c.ShouldBindQuery(&user)
	userInfo := getUserByName(name)
	c.JSON(http.StatusOK, userInfo)
}

func deleteByName(c *gin.Context) {
	name := c.Query("name")
	err := deleteUserByName(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func queryUserList(c *gin.Context) {
	users := queryUsers()
	c.JSON(http.StatusOK, users)
}

func create(c *gin.Context) {
	var user UserInfo
	_ = c.ShouldBindJSON(&user)

	id, err := createUser(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}
