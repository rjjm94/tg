// /help/help.go

package help

// Handle responds with the list of available commands and their descriptions.
func Handle() string {
	return `Here are the available commands:

/welcome - Welcome the new users
/broadcast - Broadcast a message to all users
/getstarted - Get started with the bot
/submit - Submit a ticket for support
/demo - Get a demo of the bot's functionalities
/beta - Participate in the beta testing of the bot
/mute - Mute the bot's notifications
/help - Get a list of available commands
/howto - Learn how to use the bot
/social - Connect with us on social media
/news - Get the latest news
/support - Get support for any issues`
}
