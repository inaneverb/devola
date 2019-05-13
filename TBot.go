// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom
// the Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package tgbot

import (
	"reflect"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
)

// packageLoggerInit performs initialization of package logger.
// This function will be called from TBot constructor only if it not be
// called before and if TBot constructor doesn't receive special logger
// as argument which he should use as logger.
//func packageLoggerInit() {
//	log.Lock()
//	defer log.Unlock()
//	if log.Logger != nil {
//		return
//	}
//	// TODO: default logger initialization
//}

//
type TBot struct {
	consts struct {
		apitoken string
	}

	endpoint *api.BotAPI

	cim *tChatInfoManager

	converter *tViewIDConverter

	receiver *tReceiver

	sender *tSender

	prof *tProfiler
}

//
func (b *TBot) ExtendContext(extender interface{}) error {

	if b == nil {
		// return ??
	}

	return b.receiver.ExtendContext(extender)
}

//
func (b *TBot) Run() error {

}

//
func (b *TBot) Stop() error {

}

//
func (b *TBot) Restart() error {

}

// todo: TBot.LiveStatus method

//
func New() *TBot {
	reflect.ValueOf().TryRecv()
}
