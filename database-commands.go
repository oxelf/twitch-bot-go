package main

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getAllUsers() []User {
	db, err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if err != nil {
		return []User{}
	}
	db.AutoMigrate(&User{})
	var usersList []User
	db.Find(&usersList)
	for _, user := range usersList {
		fmt.Println(user.TwitchName + " " + user.TwitchId)
	}
	return usersList
}
func getUserById(id string) User {
	db, err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if err != nil {
		return User{}
	}
	db.AutoMigrate(&User{})
	var user User
	db.First(&user, "twitch_id = ?", id)
	fmt.Println(user.EpicName + user.EpicId + user.TwitchName)
	return user
}
func getUserByTwitch(twitchName string) User {
	db, err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if err != nil {
		return User{}
	}
	db.AutoMigrate(&User{})
	var user User
	db.First(&user, "twitch_name = ?", twitchName)
	fmt.Println(user.EpicName + user.EpicId + user.TwitchName)
	return user
}
func getUserByEpic(epicName string) []User {
	db, err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if err != nil {
		return []User{}
	}
	db.AutoMigrate(&User{})
	var userList []User
	db.Where("epic_name = ?", epicName).Find(&userList)
	//db.First(&user1, "epic_name = ?", response)
	uerr := db.Where("epic_name = ?", epicName).Find(&userList).Error
	//db.First(&user1, "epic_name = ?", response).Error
	println(uerr)
	if uerr != nil {
		return []User{}
	}
	if uerr == nil {
		return userList
	}
	return userList
}
func deleteUser(method string, keyval string) string {
	//methods: "twitch", "epic", "twitchId", "epicId"
	//string res:
	//200 = Successfull deleted.
	//401 = method not supported.
	//404 = User not found/Doesnt exist.
	//500 = undefined error
	db, err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if err != nil {
		return "error. Error occured while accessing db."
	}
	db.AutoMigrate(&User{})
	methodString := ""
	switch method {
	case "twitch":
		methodString = "twitch_name = ?"
	case "twitchId":
		methodString = "twitch_id = ?"
	case "epic":
		methodString = "epic_name = ?"
	case "epicId":
		methodString = "epic_id = ?"
	default:
		return "401"
	}
	var user1 User
	//db.First(&user1, "twitch_name = ?", message.User.Name).Delete(&user1)
	uerr := db.First(&user1, methodString, keyval).Delete(&user1).Error
	println(uerr)
	if uerr != nil {
		return "404"
	}
	if uerr == nil {
		return "200"
	}
	return "500"
}
func createUser(twitchName string, twitchId string, epicName string, epicId string, deviceId string, secret string) string {

	db, dberr := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if dberr != nil {
		return "Error. Error occurred while accessing db."
	}
	var user User
	db.First(&user, "twitch_name = ?", twitchName)
	uerr := db.First(&user, "twitch_name = ?", twitchName).Error
	println(uerr)
	if uerr == nil {
		return "already exists."
	}
	if uerr != nil {
		err := db.Create(&User{EpicName: epicName, EpicId: epicId, TwitchName: twitchName, TwitchId: twitchId, DeviceId: deviceId, Secret: secret, Verified: true, CommandsUsed: 0, Lookups: 0}).Error
		if err != nil {
			return "Error. Error while creating user."
		}
		if err == nil {
			return "Successfully created."
		}
	}
	return "undefined error"
}
