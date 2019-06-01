// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package sender

import (
	"unsafe"

	api "github.com/go-telegram-bot-api/telegram-bot-api"

	"../chat"
)

// ToSend is a data struct that contains a rule of sending some Telegram message.
type ToSend struct {

	// Chat ID to which the message configs cfg and cfg2 will be sent.
	chatID chat.ID

	// True if the message configs cfg and cfg2 will be sent to the chat
	// with user, false if it's group chat or Telegram channel.
	isUserChat bool

	// Untyped pointer to ORIGINAL context object using which cfg, cfg2 are created.
	// ALWAYS POINTS TO *ctx.Ctx EVEN IF CONTEXT IS EXTENDED!
	originalCtx unsafe.Pointer

	// Untyped pointer to context object using which cfg, cfg2 are created.
	pass2finisherCtx unsafe.Pointer

	// Telegram Message API config.
	cfg api.Chattable

	// Addition Telegram Message API config.
	// Can contain only Delete message config.
	//
	// If one message should be deleted and the second be sent at the same time,
	// cfg2 will contain delete config, and cfg - regular config.
	// More info: Sender.reallySend.
	cfg2 api.Chattable

	// A set of callbacks that should be called when message
	// will be successfully sent (it applies only for cfg, not for cfg2!).
	onSuccess []unsafe.Pointer

	// A set of callbacks that should be called when message
	// is unsent (it applies only for cfg, not for cfg2!).
	onError []unsafe.Pointer

	// todo: comment
	isUpdateSession bool

	// todo: comment
	retryAttempts int8
}

//
func MakeToSend(originalCtx, pass2finisherCtx unsafe.Pointer, cfg, cfg2 api.Chattable, onSuccess, onError []unsafe.Pointer, isUpdateSession bool, retryAttempts int8) *ToSend {
	panic("implement me")
}
