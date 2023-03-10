package main

import (
	"fmt"
	"strings"

	"github.com/gempir/go-twitch-irc/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func runCmd(message twitch.PrivateMessage, client twitch.Client) {
	db, err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	db.AutoMigrate(&User{})
	if err != nil {
		panic("failed to connect database")
	}
	switch {
	case strings.Contains(strings.ToLower(message.Message), "!verify"):
		fmt.Println("command: !verify")
		verify(message, client)
	case strings.Contains(strings.ToLower(message.Message), "!verifycheck"):
		fmt.Println("command: !verifycheck")
		verifyCheck(message, client, *db)
	case strings.Contains(strings.ToLower(message.Message), "!unlink"):
		fmt.Println("command: !unlink")
		unlinkCmd(message, client)
	case strings.Contains(strings.ToLower(message.Message), "!service"):
		fmt.Println("command: !service")
		serviceCmd(message, client, *db)
	case strings.Contains(strings.ToLower(message.Message), "!verifiedcount"):
		fmt.Sprintln("command: !verifiedcount")
		verifiedCount(message, client, *db)
	case strings.Contains(strings.ToLower(message.Message), "!epic"):
		fmt.Sprintln("command: !epic")
		epicCmd(message, client)
	case strings.Contains(strings.ToLower(message.Message), "!twitchuser"):
		fmt.Sprintln("command: !twitchuser")
		twitchCmd(message, client)
	}
}
func serviceCmd(message twitch.PrivateMessage, client twitch.Client, db gorm.DB) {
	var user User
	db.First(&user, "twitch_name = ?", message.User.Name)
	deviceAuth := &DeviceAuth{DeviceId: user.DeviceId, Secret: user.Secret, AccountId: user.EpicId}
	println("deviceauth formatted")
	bearer := getBearerWithDeviceAuth(*deviceAuth)
	//println(bearer)
	client.Say("oxe1f", bearer)
	status := getLightswitchStatus(bearer)
	client.Say("oxe1f", status)
}
func unlinkCmd(message twitch.PrivateMessage, client twitch.Client) {
	db, err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if err != nil {
		return
	}
	db.AutoMigrate(&User{})
	var user1 User
	//db.First(&user1, "twitch_name = ?", message.User.Name).Delete(&user1)
	uerr := db.First(&user1, "twitch_name = ?", message.User.Name).Delete(&user1).Error
	println(uerr)
	if uerr != nil {
		client.Say(message.Channel, "Error")
	}
	if uerr == nil {
		var ans string = fmt.Sprint("Success")
		client.Say(message.Channel, ans)
	}
}
func epicCmd(message twitch.PrivateMessage, client twitch.Client) {
	db, err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if err != nil {
		return
	}
	db.AutoMigrate(&User{})
	response := ""
	split := strings.Split(message.Message, " ")
	if len(split) > 1 && split[0] == "!epic" {
		response = split[1]
		fmt.Println(response)
		response = strings.ToLower(response)
		fmt.Println(response)
	}
	var user1 User
	db.First(&user1, "twitch_name = ?", response)
	uerr := db.First(&user1, "twitch_name = ?", response).Error
	println(uerr)
	if uerr != nil {
		client.Say(message.Channel, "This user hasn't verified his account.")
	}
	if uerr == nil {
		var ans string = fmt.Sprint(response + "s epic name: " + user1.EpicName)
		client.Say(message.Channel, ans)
	}
}
func twitchCmd(message twitch.PrivateMessage, client twitch.Client) {
	db, err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if err != nil {
		return
	}
	db.AutoMigrate(&User{})
	response := ""
	split := strings.Split(message.Message, "!twitchuser ")
	if len(split) > 1 {
		fmt.Println(split[1])
		response = split[1]
	}
	var user1 []User
	db.Where("epic_name = ?", response).Find(&user1)
	//db.First(&user1, "epic_name = ?", response)
	uerr := db.Where("epic_name = ?", response).Find(&user1).Error
	//db.First(&user1, "epic_name = ?", response).Error
	println(uerr)
	if uerr != nil {
		client.Say(message.Channel, "This user has no twitch account linked with his epic.")
	}
	if uerr == nil {
		// Initialize an empty string to hold the concatenated names
		var names []string

		// Iterate over the list of users and append their Twitch names to the string
		for _, user := range user1 {
			names = append(names, user.TwitchName)
		}

		// Join the names using commas
		result := strings.Join(names, ", ")

		// Print the result
		fmt.Println(result)
		var ans string = fmt.Sprint(response + "s Twitch user name(s): " + result)
		client.Say(message.Channel, ans)
	}
}
func verifiedCount(message twitch.PrivateMessage, client twitch.Client, db gorm.DB) {
	var count int64
	db.Model(&User{}).Where("verified = ?", true).Count(&count)
	//// SELECT count(*) FROM users WHERE name = 'jinzhu'; (count)
	countInt := fmt.Sprintln("verified people: ", count)
	client.Say(message.Channel, countInt)
}
func verify(message twitch.PrivateMessage, client twitch.Client) {
	db, dberr := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if dberr != nil {
		return
	}
	var user1 User
	db.First(&user1, "twitch_name = ?", message.User.Name)
	uerr := db.First(&user1, "twitch_name = ?", message.User.Name).Error
	println(uerr)
	if uerr == nil {
		var ans string = fmt.Sprint("You're already verified as: " + user1.EpicName + " , use !unlink to remove this account.")
		client.Say(message.Channel, ans)
	}
	if uerr != nil {
		verifylink := getVerifyLink()
		client.Say(message.Channel, verifylink.VerifyUri)
		println(verifylink.DeviceCode)
		deviceauth := createDeviceAuth(verifylink.DeviceCode)
		createtest := createUser(message.User.Name, message.User.ID, deviceauth.DisplayName, deviceauth.AccountId, deviceauth.DeviceId, deviceauth.Secret)
		fmt.Println("create test: " + createtest)
		client.Say(message.Channel, createtest)
	}

}
func verifyCheck(message twitch.PrivateMessage, client twitch.Client, db gorm.DB) {
	db.AutoMigrate(&User{})
	db.Create(&User{EpicName: "Mopafu", EpicId: "", TwitchName: "oxe1f", TwitchId: "", DeviceId: "", Secret: "", Verified: false})
	var user User
	db.First(&user, "twitch_name = ?", message.User.Name)
	if !user.Verified {
		client.Reply("oxe1f", message.ID, "Sorry, it seems like youre not verified")
	}
	if user.Verified {
		client.Reply("oxe1f", message.ID, "Youre verified")
	}
	fmt.Printf("DB Test: %s", user.EpicName)
}
