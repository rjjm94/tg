//db/db.go// Path: db/db.go

package db

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	connectionString = "mongodb://127.0.0.1:27017"
	dbName           = "cntrlTGtest"
)

// DB represents the database client.
type DB struct {
	client *mongo.Client
	ctx    context.Context
}

// Message represents a message in the database.
type Message struct {
	MessageID   int       // Unique identifier of the message
	UserID      int       // User ID of the user who sent the message
	Username    string    // Username of the user who sent the message
	GroupID     int64     // Group ID of the group where the message was sent
	Text        string    // Content of the message
	MessageType string    // Type of the message
	Timestamp   time.Time // Timestamp of when the message was sent
}

// User represents a user in the database.
type User struct {
	tgbotapi.User
	IsInGroup   bool      `bson:"is_in_group"`  // Whether the user is in a group or not
	LastUpdated time.Time `bson:"last_updated"` // Timestamp of when the user's information was last updated
}

// Group represents a group in the database.
type Group struct {
	GroupName string // Name of the group
	GroupID   int64  // Unique identifier of the group
	IsActive  bool   // Whether the group is active or not
}

// Beta represents a beta in the database.
type Beta struct {
	Username      string
	UserID        int64
	GroupID       int64
	APIKey        bool
	Provider      string
	Model         string
	Email         string
	Name          string
	ContactTime   string
	ContactMethod string
	Created       time.Time
}

// Connect initializes a new database client.
func Connect() (*DB, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &DB{client: client, ctx: ctx}, nil
}

// SaveGroup saves a group in the database.
func (db *DB) SaveGroup(group Group) (*mongo.InsertOneResult, error) {
	collection := db.client.Database(dbName).Collection("groups")
	return collection.InsertOne(db.ctx, group)
}

// GetGroup retrieves a group from the database.
func (db *DB) GetGroup(groupID int64) (*Group, error) {
	collection := db.client.Database(dbName).Collection("groups")
	group := &Group{}
	err := collection.FindOne(db.ctx, bson.M{"groupid": groupID}).Decode(group)
	return group, err
}

// UpdateGroup updates a group in the database.
func (db *DB) UpdateGroup(group Group) error {
	collection := db.client.Database(dbName).Collection("groups")
	_, err := collection.UpdateOne(db.ctx, bson.M{"groupid": group.GroupID}, bson.M{"$set": bson.M{"isactive": group.IsActive}})
	return err
}

// DeactivateGroup deactivates a group in the database.
func (db *DB) DeactivateGroup(groupID int64) error {
	collection := db.client.Database(dbName).Collection("groups")
	_, err := collection.UpdateOne(db.ctx, bson.M{"groupid": groupID}, bson.M{"$set": bson.M{"isactive": false}})
	return err
}

// LogChatMessage logs a chat message in the database.
// LogChatMessage logs a chat message in the database.
func (db *DB) LogChatMessage(chatMessage Message) error {
	collection := db.client.Database(dbName).Collection("messages")
	_, err := collection.InsertOne(db.ctx, chatMessage)
	return err
}

// LogUserProfile logs a user profile in the database.
func (db *DB) LogUserProfile(userProfile User) error {
	collection := db.client.Database(dbName).Collection("users")
	dbUser := User{
		User: tgbotapi.User{
			ID:           userProfile.User.ID,
			FirstName:    userProfile.User.FirstName,
			LastName:     userProfile.User.LastName,
			UserName:     userProfile.User.UserName,
			LanguageCode: userProfile.User.LanguageCode,
			IsBot:        userProfile.User.IsBot,
		},
		IsInGroup:   true,
		LastUpdated: time.Now(),
	}

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"user.id": userProfile.User.ID}
	update := bson.M{"$set": dbUser}

	_, err := collection.UpdateOne(db.ctx, filter, update, opts)
	return err
}

// SaveBeta saves a beta in the database.
func (db *DB) SaveBeta(betaInfo Beta) error {
	collection := db.client.Database(dbName).Collection("beta")
	_, err := collection.InsertOne(db.ctx, betaInfo)
	return err
}
