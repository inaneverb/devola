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

// 'tSenderParam' is the aliass to function that applying to the some
// sender object.
// It uses as parameters for sender constructor to initialize consts.
type tSenderParam func(s *tSender)

//
func SenderAlwaysUseMD() tSenderParam {
	return func(s *tSender) {
		s.consts.isAlwaysUseMD, s.consts.isAlwaysUseHTML = true, false
	}
}

//
func SenderAlwaysUseHTML() tSenderParam {
	return func(s *tSender) {
		s.consts.isAlwaysUseHTML, s.consts.isAlwaysUseMD = true, false
	}
}

//
func SenderMaxSendAttempts(n int) tSenderParam {
	if n < -1 || n > 100 {
		return nil
	}
	return func(s *tSender) { s.consts.maxSendRetryAttempts = n }
}