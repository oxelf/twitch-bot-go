package main

import (
	"fmt"

	twitch "github.com/gempir/go-twitch-irc/v4"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	clientUsername            = "oxe1f"
	clientAuthenticationToken = "oauth:m3xiknq3f6kdxszlkktq6xtwwmx1zq"
)

type User struct {
	gorm.Model
	EpicName     string
	EpicId       string
	TwitchName   string
	TwitchId     string
	DeviceId     string
	Secret       string
	Verified     bool
	CommandsUsed int
	Lookups      int
}

func main() {
	router := gin.Default()
	router.GET("/user/delete/:method/:user", ginDeleteUser)
	router.GET("/user/get/twitch/:user", ginGetUserByTwitch)
	router.GET("/user/get/epic/:user", ginGetUserByEpic)
	router.GET("/user/get/id/:user", ginGetUserById)
	router.GET("/user/get/all", ginGetAllUsers)
	go func() {
		router.Run("localhost:8082")
	}()
	getAllUsers()
	//deleteUser("twitch", "testitttd")
	//deltest := deleteUser("twitch", "testitttd")

	//fmt.Println("del: " + deltest)
	db, err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&User{})
	client := twitch.NewClient(clientUsername, clientAuthenticationToken)
	client.Join("oxe1f")
	println("test")
	client.OnConnect(func() {
		client.Say("oxe1f", "bot is now running.")
	})

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		if message.User.Name == clientUsername {
			return
		}
		runCmd(message, *client)
	})

	cerr := client.Connect()
	if cerr != nil {
		panic(cerr)
	}

}

func ginDeleteUser(c *gin.Context) {
	println("delete request.")
	c.Header("Access-Control-Allow-Origin", "*")
	username := c.Param("user")
	method := c.Param("method")
	delRes := deleteUser(method, username)
	if delRes == "200" {
		c.JSON(200, gin.H{
			"status": "Successfully deleted.",
		})
	}
	if delRes == "401" {
		c.JSON(401, gin.H{
			"status": "Method not supported.",
		})
	}
	if delRes == "404" {
		c.JSON(404, gin.H{
			"status": "User not found/User doesnt exist.",
		})
	}
	c.JSON(500, gin.H{
		"status": "Undefined Error.",
	})
}
func ginGetUserByTwitch(c *gin.Context) {
	println("delete request.")
	c.Header("Access-Control-Allow-Origin", "*")
	username := c.Param("user")
	fmt.Println("username: " + username)
	user := getUserByTwitch(username)
	switch user.EpicName {
	default:
		c.JSON(200, gin.H{
			"status":       "Success",
			"epicName":     user.EpicName,
			"epicId":       user.EpicId,
			"twitchName":   user.TwitchName,
			"twitchId":     user.TwitchId,
			"deviceId":     user.DeviceId,
			"secret":       user.Secret,
			"CommandsUsed": user.CommandsUsed,
			"lookups":      user.Lookups,
		})
	case "":
		c.JSON(500, gin.H{
			"status": "User not found/User doesnt exist.",
		})
	}
}
func ginGetUserById(c *gin.Context) {
	println("delete request.")
	c.Header("Access-Control-Allow-Origin", "*")
	username := c.Param("user")
	fmt.Println("username: " + username)
	user := getUserById(username)
	switch user.EpicName {
	default:
		c.JSON(200, gin.H{
			"status":       "Success",
			"epicName":     user.EpicName,
			"epicId":       user.EpicId,
			"twitchName":   user.TwitchName,
			"twitchId":     user.TwitchId,
			"deviceId":     user.DeviceId,
			"secret":       user.Secret,
			"CommandsUsed": user.CommandsUsed,
			"lookups":      user.Lookups,
		})
	case "":
		c.JSON(500, gin.H{
			"status": "User not found/User doesnt exist.",
		})
	}
}
func ginGetUserByEpic(c *gin.Context) {
	println("delete request.")
	c.Header("Access-Control-Allow-Origin", "*")
	username := c.Param("user")
	fmt.Println("username: " + username)
	userList := getUserByEpic(username)
	switch userList[0].EpicName {
	default:
		c.JSON(200, gin.H{
			"status":      "Success",
			"linkedUsers": userList,
		})
	case "":
		c.JSON(500, gin.H{
			"status": "User not found/User doesnt exist.",
		})
	}
}
func ginGetAllUsers(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userList := getAllUsers()
	switch userList[0].EpicName {
	default:
		c.JSON(200, gin.H{
			"count":  len(userList),
			"status": "Success",
			"users":  userList,
		})
	case "":
		c.JSON(500, gin.H{
			"status": "Error.",
		})
	}

}
