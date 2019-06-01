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

package receiver

// 'tReceiverParam' is the alias to function that applying to the some
// receiver object.
// It uses as parameters for receiver constructor to initialize consts.
type tReceiverParam func(r *Receiver)

//
func ReceiverSingleExecutor() tReceiverParam {
	return func(r *Receiver) { r.consts.executorsCount = 1 }
}

// n must be in the range [0..16]
func ReceiverMultiExecutors(n ...int) tReceiverParam {
	if len(n) == 0 || n[0] <= 0 || n[0] > 16 {
		return ReceiverSingleExecutor()
	}
	return func(r *Receiver) { r.consts.executorsCount = n[0] }
}

//
func ReceiverForceServe(start bool) tReceiverParam {
	return func(r *Receiver) { r.consts.forceServe = start }
}
