package main

import (
	"context"
	"encoding/json"
	"fmt"
)

type Request struct {
	Body string `json:"body"`
}

type Update struct {
	UpdateId int32 `json:"update_id"`
	Message  Message
}

type Message struct {
	From User
	Chat Chat
	Text string `json:"text"`
}

type User struct {
	Id        int32  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type Chat struct {
	Id       int32  `json:"id"`
	Type     string `json:"type"`
	Username string `json:"username"`
}

type Response struct {
	StatusCode int32             `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       interface{}       `json:"body"`
}

type UpdateResponse struct {
	Method string `json:"method"`
	ChatId int32  `json:"chat_id"`
	Text   string `json:"text"`
}

func Handler(ctx context.Context, req []byte) (Response, error) {
	var request Request
	err := json.Unmarshal(req, &request)
	if err != nil {
		return Response{}, err
	}

	var update Update
	err = json.Unmarshal([]byte(request.Body), &update)
	if err != nil {
		return Response{}, err
	}

	InitDB(ctx)

	message := update.Message

	fmt.Println(fmt.Sprintf("Got and update from chat id %d, user %s (@%s) with text %s", message.Chat.Id, message.From.FirstName, message.From.Username, message.Text))

	command := extractCommandFrom(message.From.Id, update.Message.Text)

	respMessage, err := command.Execute()

	if err != nil {
		fmt.Println(err.Error())
		respMessage = fmt.Sprintf("Thank you for your message. Not sure what you mean. Please try to send me some commands by typing \"/\".")
	}

	response := Response{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body: UpdateResponse{
			Method: "sendMessage",
			ChatId: message.Chat.Id,
			Text:   respMessage,
		},
	}

	return response, nil
}
