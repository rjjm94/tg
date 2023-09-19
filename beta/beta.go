//beta/beta.go

package beta

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"sync"
	db "tg/db"
)

var (
	bot         *tgbotapi.BotAPI          // Declare the bot variable at the package level
	betaInfoMap = make(map[int64]db.Beta) // Map to store BetaInfo for each user
	mu          sync.Mutex                // Mutex to prevent data race
)

func SetBot(b *tgbotapi.BotAPI) {
	bot = b // Set the bot variable
}

func Handle(userID int64, groupID int64, userName string) (tgbotapi.MessageConfig, db.Beta) {
	mu.Lock()                       // Lock the mutex
	betaInfo := betaInfoMap[userID] // Retrieve the betaInfo for the user
	mu.Unlock()                     // Unlock the mutex

	betaInfo.Username = userName
	betaInfo.UserID = userID
	betaInfo.GroupID = groupID

	row := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Yes", "yes"),
		tgbotapi.NewInlineKeyboardButtonData("No", "no"),
	}

	var keyboard [][]tgbotapi.InlineKeyboardButton
	keyboard = append(keyboard, row)

	markup := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}

	msg := tgbotapi.NewMessage(userID, "Do you have an API Key?")
	msg.ReplyMarkup = markup

	mu.Lock()                      // Lock the mutex
	betaInfoMap[userID] = betaInfo // Update the betaInfo for the user
	mu.Unlock()                    // Unlock the mutex

	return msg, betaInfo
}

func HandleProvider(update *tgbotapi.Update, betaInfo db.Beta) (tgbotapi.EditMessageTextConfig, db.Beta) {
	betaInfo.APIKey = true

	row := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Azure", "azure"),
		tgbotapi.NewInlineKeyboardButtonData("OpenAI", "openai"),
	}

	var keyboard [][]tgbotapi.InlineKeyboardButton
	keyboard = append(keyboard, row)

	markup := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}

	msg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "Do you have Azure or OpenAI API key?")
	msg.ReplyMarkup = &markup

	return msg, betaInfo
}

func HandleModel(update *tgbotapi.Update, betaInfo db.Beta) (tgbotapi.EditMessageTextConfig, db.Beta) {
	betaInfo.Provider = update.CallbackQuery.Data

	row := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("GPT3.5", "gpt3.5"),
		tgbotapi.NewInlineKeyboardButtonData("GPT4", "gpt4"),
		tgbotapi.NewInlineKeyboardButtonData("GPT4-32k", "gpt4-32k"),
	}

	var keyboard [][]tgbotapi.InlineKeyboardButton
	keyboard = append(keyboard, row)

	markup := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}

	msg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "What model do you have access to?")
	msg.ReplyMarkup = &markup

	return msg, betaInfo
}

func HandleEmail(update *tgbotapi.Update, betaInfo db.Beta) (tgbotapi.Chattable, db.Beta) {
	var msg tgbotapi.Chattable

	mu.Lock()                                             // Lock the mutex
	betaInfo = betaInfoMap[int64(update.Message.From.ID)] // Retrieve the betaInfo for the user
	mu.Unlock()                                           // Unlock the mutex

	if update.Message != nil && update.Message.Text != "" {
		// Save the input as email
		betaInfo.Email = update.Message.Text
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "What is your name?")
	} else {
		// If update.Message is nil or text is empty, return a default message
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Please enter your email:")
	}

	mu.Lock()                                             // Lock the mutex
	betaInfoMap[int64(update.Message.From.ID)] = betaInfo // Update the betaInfo for the user
	mu.Unlock()                                           // Unlock the mutex

	return msg, betaInfo
}

func HandleName(update *tgbotapi.Update) (tgbotapi.Chattable, db.Beta) {
	var msg tgbotapi.Chattable

	mu.Lock()                                              // Lock the mutex
	betaInfo := betaInfoMap[int64(update.Message.From.ID)] // Retrieve the betaInfo for the user
	mu.Unlock()                                            // Unlock the mutex

	if update.Message != nil {
		betaInfo.Name = update.Message.Text
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "What is the best time and method of contacting you?")
	}

	mu.Lock()                                             // Lock the mutex
	betaInfoMap[int64(update.Message.From.ID)] = betaInfo // Update the betaInfo for the user
	mu.Unlock()                                           // Unlock the mutex

	return msg, betaInfo
}

func HandleContact(update *tgbotapi.Update, betaInfo db.Beta) (tgbotapi.Chattable, db.Beta) {
	var msg tgbotapi.Chattable

	if update.Message != nil {
		betaInfo.ContactMethod = update.Message.Text
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Please review your information:")
	}

	return msg, betaInfo
}

func HandleSummary(update *tgbotapi.Update, betaInfo db.Beta) tgbotapi.Chattable {
	summary := fmt.Sprintf("API Key: %v\nProvider: %s\nModel: %s\nEmail: %s\nName: %s\nContact: %s\n",
		betaInfo.APIKey, betaInfo.Provider, betaInfo.Model, betaInfo.Email, betaInfo.Name, betaInfo.ContactMethod)

	row := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Submit", "submit"),
		tgbotapi.NewInlineKeyboardButtonData("Reset", "reset"),
	}

	var keyboard [][]tgbotapi.InlineKeyboardButton
	keyboard = append(keyboard, row)

	markup := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, summary)
	msg.ReplyMarkup = &markup

	return msg
}
