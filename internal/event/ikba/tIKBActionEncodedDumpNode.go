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

// tIKBActionEncodedDumpNode is type for method tIKBActionEncoded.dump.
//
// This type represents one node of encoded IKB action tIKBActionEncoded.
// All fields has JSON tags and it's easy to JSON dump output.
//
// Object of this type fills by tIKBActionEncoded.dump method.
//
// More info: tIKBActionEncoded, tIKBActionEncoded.dump.
type tIKBActionEncodedDumpNode struct {

	// Type is a description of IKB encoded node.
	// If node represents encoded argument, Type is a string description
	// of argument's type.
	Type string `json:"type"`

	// Position in encoded IKB action where this node byte view starts.
	Pos byte `json:"pos"`

	// Position in encoded IKB action where encoded argument's header is placed.
	// If node is not about encoded argument, this field is zero.
	PosType byte `json:"pos_type,omitempty"`

	// Position in encoded IKB action where encoded argument's content is placed.
	// If node is not about encoded argument, this field is zero or the same
	// as Pos field.
	// Anyway better use Pos field for non-argument's nodes.
	PosContent byte `json:"pos_value,omitempty"`

	// This is RAW type header of encoded argument if node is about it.
	// Otherwise it is zero.
	TypeHeader byte `json:"type_encoded,omitempty"`

	// Typed RAW value stored in interface{} variable that represents by node.
	// Thus, for example, if node about encoded int8 argument, Value is this
	// int8 argument.
	Value interface{} `json:"value"`
}
