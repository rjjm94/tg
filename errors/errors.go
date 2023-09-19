//errors/errors.go

package errors

import (
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net"
	"strings"
)

type ErrorCheckFunc func(err error) bool

var errorChecks = []ErrorCheckFunc{
	IsBotKicked,
	IsMongoNoDocuments,
	IsMongoDuplicateKey,
	IsTelegramAPIError,
	IsRateLimited,
	IsTimeout,
	IsDatabaseConnectionError,
	IsInvalidMessage,
	IsInvalidCommand,
	IsUserNotFound,
	IsGroupNotFound,
	IsFailedToSendMessage,
	IsFailedToUpdateDatabase,
}

func HandleError(err error) error {
	if err != nil {
		for _, check := range errorChecks {
			if check(err) {
				log.Printf("Handled error: %v", err)
				break
			}
		}
		return errors.New(err.Error())
	}
	return nil
}

func IsBotKicked(err error) bool {
	return err != nil && err.Error() == "Forbidden: bot was kicked from the group chat"
}

func IsMongoNoDocuments(err error) bool {
	return err == mongo.ErrNoDocuments
}

func IsMongoDuplicateKey(err error) bool {
	we, ok := err.(mongo.WriteException)
	return ok && we.HasErrorCode(11000)
}

func IsTelegramAPIError(err error) bool {
	return err != nil && err.Error() == "TELEGRAM_ERROR"
}

func IsRateLimited(err error) bool {
	return strings.Contains(err.Error(), "rate limit")
}

func IsTimeout(err error) bool {
	_, ok := err.(net.Error)
	return ok && err.(net.Error).Timeout()
}

func IsDatabaseConnectionError(err error) bool {
	return strings.Contains(err.Error(), "failed to connect to database")
}

func IsInvalidMessage(err error) bool {
	return strings.Contains(err.Error(), "invalid message")
}

func IsInvalidCommand(err error) bool {
	return strings.Contains(err.Error(), "invalid command")
}

func IsUserNotFound(err error) bool {
	return strings.Contains(err.Error(), "user not found")
}

func IsGroupNotFound(err error) bool {
	return strings.Contains(err.Error(), "group not found")
}

func IsFailedToSendMessage(err error) bool {
	return strings.Contains(err.Error(), "failed to send message")
}

func IsFailedToUpdateDatabase(err error) bool {
	return strings.Contains(err.Error(), "failed to update database")
}
