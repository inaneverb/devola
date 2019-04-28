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

// tParams is the type of storage all parameters that can be applied
// to the some types in that SDK.
// More info: Params.
type tParams struct {

	// Lirester have params to manage behaviour of "Limit Telegram Avoider" class.
	// Read more: tLirester, tLiresterParams.
	Lirester tLiresterParams
}

// Params is a storage of all SDK params grouped by SDK types.
//
// Param is a special type - an alias of function that takes some type
// by its pointer and changes behaviour into taken object.
//
// You can read more about what these params do from doc to the their types
// (all parameters of some type has its type's params storage).
// Read these field's types docs.
var Params tParams

// By expr above memory already allocated for tParams struct and all nested.
// Params' fields (nested structs) initialized by its init functions. 
