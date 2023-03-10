package main

// db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
// 			if err != nil {
// 				panic("failed to connect database")
// 			}

// 			// Migrate the schema
// 			db.AutoMigrate(&Product{})

// 			// Create
// 			db.Create(&Product{Code: "MoinFortnite", Price: 187})

// 			// Read
// 			var product Product
// 			db.First(&product)
// 			fmt.Printf("DB Test: %s", product.Code) // find product with integer primary key++package main

// 			import (
// 				twitch "github.com/gempir/go-twitch-irc/v4"
// 				"gorm.io/gorm"
// 			)

// 			const (
// 				clientUsername            = "oxe1f"
// 				clientAuthenticationToken = "oauth:m3xiknq3f6kdxszlkktq6xtwwmx1zq"
// 			)

// 			type Product struct {
// 				gorm.Model
// 				Code  string
// 				Price uint
// 			}

// 			func main() {
// 				println("running main")

// 				client := twitch.NewClient(clientUsername, clientAuthenticationToken)
// 				defer client.Depart("moin")
// 				client.OnConnect(func() {
// 					println("connected")
// 				})
// 				client.Join("oxe1f")
// 				client.Say("oxe1f", "testing dm")
// 				client.OnPrivateMessage(func(message twitch.PrivateMessage) {
// 					if message.Message == "db read" {
// 						println("db read")
// 					}
// 					println(message.Raw)
// 					if message.Message == "whisper" {
// 						client.Say("oxe1f", "/w kunia_bot moooiiinnn")

// 					}

// 				})
// 			}
