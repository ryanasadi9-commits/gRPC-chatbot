package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"hamrahTask1/proto"
)

func main() {
	// client connecting to server port
	connection, err := grpc.NewClient("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer connection.Close()

	client := proto.NewChatbotServiceClient(connection)
	scanner := bufio.NewScanner(os.Stdin)
	var currentUserID string = ""
	var currentPage string = "auth"
	run := true

	fmt.Print("\033[2J\033[H")
	fmt.Printf("			--- auth ---\n\n")

	for run {
		if currentPage == "auth" || currentPage == "main" {
			fmt.Print("--> ")
		} else {
			fmt.Print("- ")
		}
		if !scanner.Scan() {
			run = false
		}

		input := scanner.Text()
		words := strings.Fields(input)

		if len(words) == 0 {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

		switch words[0] {
		case "commands":
			fmt.Println("- /commands")
			fmt.Println("- /register <username> <password> <age> <gender(true/false)>")
			fmt.Println("- /login <username> <password>")
			fmt.Println("- /show <contact_username> (<count>)")
			fmt.Println("- <message_text...>")
			fmt.Println("- /main")
			fmt.Println("- /close")

		case "/main":
			if currentPage == "auth" || currentPage == "main" {
				fmt.Println("invalid command")
				continue
			}
			currentPage = "main"
			response, err := client.ShowUserList(ctx, &proto.Empty{})
			if err != nil {
				log.Fatalf("Failed to get user list: %v", err)
			} else {
				fmt.Print("\033[2J\033[H")
				fmt.Println("			--- user list ---")
				for i := 0; i < len(response.UserList); i++ {
					fmt.Printf(". %s (%t)\n", response.UserList[i], response.Status[i])
				}
			}

		case "/register":
			if len(words) != 5 {
				fmt.Println("-----> /register <username> <password> <age> <gender(true/false)>")
				cancel()
				continue
			}
			age, _ := strconv.ParseInt(words[3], 10, 64)
			gender, _ := strconv.ParseBool(words[4])

			response, err := client.Register(ctx, &proto.RegisterRequest{
				Username: words[1],
				Password: words[2],
				Age:      age,
				Gender:   gender,
			})

			if err != nil {
				log.Fatalf("Failed to register: %v", err)
			} else {
				fmt.Printf("-%s\n", response.Message)
			}

		case "/login":
			if currentUserID != "" {
				fmt.Println("you are logged in. to switch account you need to log out first!")
				continue
			}

			if len(words) < 3 {
				log.Println("-----> /login <username> <password>")
				cancel()
				continue
			}
			response, err := client.Login(ctx, &proto.LoginRequest{
				Username: words[1],
				Password: words[2],
			})

			if err != nil {
				log.Printf("Failed to login: %v", err)
			} else {
				//fmt.Printf("-%s\n", response.Message)
				currentPage = "main"
				currentUserID = response.Id

				response, err := client.ShowUserList(ctx, &proto.Empty{})
				if err != nil {
					log.Fatalf("Failed to get user list: %v", err)
				} else {
					fmt.Print("\033[2J\033[H")
					fmt.Println("			--- contact list ---")
					for i := 0; i < len(response.UserList); i++ {
						fmt.Printf(". %s (%t)\n", response.UserList[i], response.Status[i])
					}
				}
			}

		case "/logout":
			if currentUserID == "" {
				fmt.Println("Error: You must login first!")
				cancel()
				continue
			}
			response, _ := client.Logout(ctx, &proto.LogoutRequest{Id: currentUserID})
			fmt.Printf("-%s\n", response.Message)
			currentUserID = ""
			currentPage = "auth"
			fmt.Print("\033[2J\033[H")
			fmt.Printf("			--- auth ---\n\n")

		case "/show":
			if currentUserID == "" {
				fmt.Println("Error: You must login first!")
				cancel()
				continue
			}
			if len(words) < 2 {
				fmt.Println("-----> /show <contact_username> (<count>)")
				cancel()
				continue
			}
			var count int64 = 50
			if len(words) == 3 {
				a, _ := strconv.ParseInt(words[2], 10, 64)
				count = a
			}

			contactUsername := words[1]

			response, err := client.ShowMessage(ctx, &proto.ShowMessagesRequest{
				Id:        currentUserID,
				ContactId: contactUsername,
				Count:     count,
			})

			if err != nil {
				log.Printf("Failed to show messages: %v", err)
			} else {
				if response.Success {
					fmt.Print("\033[2J\033[H")
					fmt.Printf("			--- %s ---\n\n", contactUsername)
					currentPage = contactUsername
				}

				for _, msg := range response.Message {
					fmt.Println(msg)
				}
			}

		case "/close":
			run = false
			if currentUserID != "" {
				client.Logout(ctx, &proto.LogoutRequest{Id: currentUserID})
				currentUserID = ""
			}

		default:
			if currentPage == "auth" || currentPage == "main" {
				fmt.Println("wrong command")
				continue
			}

			//message := strings.Join(words[0:], " ")
			message := input

			_, err := client.SendMessage(ctx, &proto.SendMessageRequest{
				SenderId:   currentUserID,
				ReceiverId: currentPage,
				Message:    message,
			})

			if err != nil {
				log.Printf("Failed to send message: %v", err)
			} else {
				//fmt.Printf("-%s\n", response.Message)
			}
		}
		cancel()
	}
}
