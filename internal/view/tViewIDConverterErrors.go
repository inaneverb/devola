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

// tViewIDConverterErrors is the type of storage of tViewIDConverter errors.
// Is a part of tErrors type.
type tViewIDConverterErrors struct {

	// Invalid View ID error.
	// Returned:
	// - From Encode method if invalid View ID is passed as argument
	// - From Register method if invalid View ID is passed as argument
	// - From Decode method if after decoding encoded View ID, the decoded
	//   View ID is invalid.
	InvalidViewID *tViewIDConverterError

	// Invalid encoded View ID error.
	// Returned:
	// - From Decode method if invalid encoded View ID is passed as argument
	// - From Encode method if after encoding View ID, the encoded View ID
	//   is invalid.
	InvalidEncodedViewID *tViewIDConverterError

	// Unregistered View ID error.
	// Returned:
	// - From Encode method if an unregistered by Register method View ID
	//   is passed as argument.
	// - From Decode method if passed an encoded View ID without decoded View ID
	//   pair (as a conclusion, there is no registered View ID with received
	//   encoded View ID).
	UnregisteredViewID *tViewIDConverterError

	// Already registered View ID error.
	// Returned:
	// - From Register method if passed View ID is already registered by
	//   one of previous Register calls.
	AlreadyRegisteredViewID *tViewIDConverterError

	// Registered View IDs limit is reached error.
	// Returned:
	// - From Register method if there is limit of registered View IDs
	//   is reached and no one View ID can be registered anymore.
	LimitReached *tViewIDConverterError
}

// Initializes ViewIDConverter field of Errors object.
//
// Direct initialization tViewIDConverterErrors errors storage as a part of
// tErrors object.
func init() {

	Errors.ViewIDConverter.InvalidViewID = &tViewIDConverterError{
		What: "View ID Converter: Invalid View ID.",
	}

	Errors.ViewIDConverter.InvalidEncodedViewID = &tViewIDConverterError{
		What: "View ID Converter: Invalid encoded View ID.",
	}

	Errors.ViewIDConverter.UnregisteredViewID = &tViewIDConverterError{
		What: "View ID Converter: This View ID is not registered.",
	}

	Errors.ViewIDConverter.AlreadyRegisteredViewID = &tViewIDConverterError{
		What: "View ID Converter: This View ID is already registered.",
	}

	Errors.ViewIDConverter.LimitReached = &tViewIDConverterError{
		What: "View ID Converter: Limit of registered View IDs is reached.",
	}
}
