//handlers/handlers.go

package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
	beta "tg/beta"
	db "tg/db"
	errors "tg/errors"
	help "tg/help"
	"time"
)

// HandleMessage logs the chat message and user profile in the database.
// It also returns a response based on the content of the message.
func HandleMessage(update *tgbotapi.Update) tgbotapi.Chattable {
	var response tgbotapi.Chattable
	var betaInfo db.Beta // Create a variable to store the Beta information

	// Check if the update is a callback query or a message
	if update.CallbackQuery != nil {
		response, betaInfo = handleCallbackQuery(update, betaInfo)
	} else if update.Message != nil {
		response, betaInfo = handleTextMessage(update, betaInfo)
	}

	// Log the message and user profile if there is a response
	if response != nil {
		logMessageAndUserProfile(update, response)
	}

	return response
}

// handleCallbackQuery handles a callback query from a user.
func handleCallbackQuery(update *tgbotapi.Update, betaInfo db.Beta) (tgbotapi.Chattable, db.Beta) {
	var response tgbotapi.Chattable // Define response here

	// Handle callback queries here
	switch update.CallbackQuery.Data {
	case "yes":
		betaInfo.APIKey = true
		response, betaInfo = beta.HandleProvider(update, betaInfo)
	case "no":
		betaInfo.APIKey = false
		response = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Please obtain an API key.")
	case "azure", "openai":
		response, betaInfo = beta.HandleModel(update, betaInfo)
	case "gpt3.5", "gpt4", "gpt4-32k":
		response, betaInfo = beta.HandleEmail(update, betaInfo) // Adjusted call to HandleEmail

	case "name":
		response, betaInfo = beta.HandleContact(update, betaInfo)
	case "contact":
		betaInfo = db.Beta{
			Username:      update.CallbackQuery.From.UserName,
			UserID:        int64(update.CallbackQuery.From.ID),
			GroupID:       update.CallbackQuery.Message.Chat.ID,
			APIKey:        betaInfo.APIKey,
			Provider:      betaInfo.Provider,
			Model:         betaInfo.Model,
			Email:         betaInfo.Email,
			Name:          betaInfo.Name,
			ContactTime:   betaInfo.ContactTime,
			ContactMethod: betaInfo.ContactMethod,
			Created:       time.Now(),
		}
		response = beta.HandleSummary(update, betaInfo)
	case "submit":
		err := db.SaveBeta(betaInfo) // Save the Beta information to the database
		if err != nil {
			log.Fatal(errors.HandleError(err))
		}
	case "reset":
		response, betaInfo = beta.Handle(int64(update.CallbackQuery.From.ID), update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.From.UserName)
	}
	return response, betaInfo
}

// handleTextMessage handles a text message from a user.
func handleTextMessage(update *tgbotapi.Update, betaInfo db.Beta) (tgbotapi.Chattable, db.Beta) {
	var response tgbotapi.Chattable // Define response here

	// Handle text messages here
	switch update.Message.Text {
	case "/beta":
		response, _ = beta.Handle(int64(update.Message.From.ID), update.Message.Chat.ID, update.Message.From.UserName)
	default:
		if strings.Contains(update.Message.Text, "@") && betaInfo.Email == "" {
			// Assume the message is an email if it contains "@"
			betaInfo.Email = update.Message.Text
			response, _ = beta.HandleName(update)
		} else if betaInfo.Email != "" && betaInfo.Name == "" {
			// If the bot is waiting for the user's name, store the incoming message as the name
			betaInfo.Name = update.Message.Text
			response, _ = beta.HandleContact(update, betaInfo)
		} else if betaInfo.Name != "" && betaInfo.ContactMethod == "" {
			// If the bot is waiting for the user's contact information, store the incoming message as the contact method
			betaInfo.ContactMethod = update.Message.Text
			response = beta.HandleSummary(update, betaInfo)
		} else {
			// Handle other messages
			response = handleOtherMessages(update)
		}
	}
	return response, betaInfo
}

// logMessageAndUserProfile logs a chat message and user profile in the database.
func logMessageAndUserProfile(update *tgbotapi.Update, response tgbotapi.Chattable) {
	var message db.Message
	var user *db.User

	// Check if the update is a message or a callback query
	if update.Message != nil {
		message = db.Message{
			MessageID: update.Message.MessageID,
			UserID:    update.Message.From.ID,
			Username:  update.Message.From.UserName,
			GroupID:   update.Message.Chat.ID,
			Text:      update.Message.Text,
			Timestamp: update.Message.Time(),
		}

		// Determine the message type based on the response
		switch v := response.(type) {
		case tgbotapi.MessageConfig:
			message.MessageType = v.Text
		case tgbotapi.EditMessageTextConfig:
			message.MessageType = v.Text
		}

		user = &db.User{
			User: tgbotapi.User{
				ID:        update.Message.From.ID,
				FirstName: update.Message.From.FirstName,
				LastName:  update.Message.From.LastName,
				UserName:  update.Message.From.UserName,
			},
			IsInGroup:   true,
			LastUpdated: time.Now(),
		}
	} else if update.CallbackQuery != nil {
		message = db.Message{
			MessageID:   update.CallbackQuery.Message.MessageID,
			UserID:      update.CallbackQuery.From.ID,
			Username:    update.CallbackQuery.From.UserName,
			GroupID:     update.CallbackQuery.Message.Chat.ID,
			Text:        update.CallbackQuery.Data,
			MessageType: update.CallbackQuery.Data,
			Timestamp:   update.CallbackQuery.Message.Time(),
		}

		user = &db.User{
			User: tgbotapi.User{
				ID:        update.CallbackQuery.From.ID,
				FirstName: update.CallbackQuery.From.FirstName,
				LastName:  update.CallbackQuery.From.LastName,
				UserName:  update.CallbackQuery.From.UserName,
			},
			IsInGroup:   true,
			LastUpdated: time.Now(),
		}
	}

	// Log the chat message and user profile
	err := db.LogChatMessage(message)
	if err != nil {
		log.Fatal(errors.HandleError(err))
	}

	err = db.LogUserProfile(*user)
	if err != nil {
		log.Fatal(errors.HandleError(err))
	}
}

// handleOtherMessages handles other types of messages from a user.
func handleOtherMessages(update *tgbotapi.Update) tgbotapi.Chattable {
	var response tgbotapi.Chattable

	// Other cases
	switch update.Message.Text {
	case "/help":
		response = tgbotapi.NewMessage(update.Message.Chat.ID, help.Handle())

	}
	return response
}
