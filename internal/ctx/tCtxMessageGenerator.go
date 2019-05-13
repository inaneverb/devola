// Copyright Â© 2018. All rights reserved.
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

//
type tCtxMessageGenerator struct {
	isUpdateSession bool
	onSuccess       []FViewSuccessFinisher
	onError         []FViewErrorFinisher
	retryAttempts   int8
	isHTML          bool
	isMarkdown      bool
	text            string
	replyTo         int

	g *tCtxKeyboardGenerator
}

//
func (gen *tCtxMessageGenerator) genParseMode() string {
	if ctx.g.isHTML {
		nmsg.ParseMode = "HTML"
	}
	if ctx.g.isMarkdown {
		nmsg.ParseMode = "Markdown"
	}
}

//
func (gen *tCtxMessageGenerator) genKb() interface{} {
	// 	if !ctx.g.isDeleteKeyboard {
	// 		if ctx.g.keyboard != nil { nmsg.ReplyMarkup = ctx.g.keyboard }
	// 	} else {
	// 		nmsg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	// 	}
}