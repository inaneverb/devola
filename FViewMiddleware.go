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

// FViewMiddleware is the function, callback, that registered by you as
// a function for some event and it will be called when that event is occurred,
// but BEFORE event handler (tHandlerCallback) calling!
//
// Moreover, your function, a middleware, should return true or false.
// It's an indicator.
// If event handler callback(s) must be called and event
// should be proceed, you must return true from ALL registered middleewares.
// If at least one middleware returns false, the rest of event middlewares
// and event handlers will never called for occurred event and event
// will be discarded.
//
// It guarantees, that when the middlewares is called, the TCtx object
// is considered completely created and TCtx object will be passed to the
// handler callback in the form in which it is after middleware call
// will finish.
// So, it means all changes in TCtx object made inside middleware will be saved
// and handler callback will receive modified by you TCtx object.
// But there is no internal modifies of TCtx object between middleware
// and handler calls.
//
// More info: tReceiver, TCtx, tBehaviourCreator, TBot.
type FViewMiddleware func(c *TCtx) bool
