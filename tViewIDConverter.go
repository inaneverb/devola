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

// tViewIDConverter is the encoder from readable View ID format (tViewID)
// to internal (used in SDK) format (tViewIDEncoded).
//
// Works using a two maps:
// A view from tViewID to encoded value (tViewIDEncoded) and vice-versa.
// Encode, Decode methods using these maps and Register method just to add
// entries in these maps.
type tViewIDConverter struct {

	// When some new View ID registers, the entry of this link
	// View ID -> View ID encoded places here.
	mEncodeStorage map[tViewID]tViewIDEncoded

	// When some new View ID registers, the entry of this link
	// View ID encoded -> View ID places here.
	mDecodeStorage map[tViewIDEncoded]tViewID

	// Generator for encoded View ID.
	// Increases by one for each new View ID by Register method.
	encodedIDGenerator tViewIDEncoded
}

// Encode tries to encode passed View ID and return it.
//
// So, it will be successfully only if passed viewID is already registered
// by Register method and all internal checks are passed.
func (vidc *tViewIDConverter) Encode(viewID tViewID) (tViewIDEncoded, *tViewIDConverterError) {

	if !viewID.isValid() {
		return cViewIDEncodedNull, eViewIDConverter.InvalidViewID
	}

	viewIDEncoded, found := vidc.mEncodeStorage[viewID]
	if !found {
		return cViewIDEncodedNull, eViewIDConverter.UnregisteredViewID
	}

	if !viewIDEncoded.isValid() {
		return cViewIDEncodedNull, eViewIDConverter.InvalidEncodedViewID
	}

	return viewIDEncoded, nil
}

// Decode tries to decode passed encoded View ID and return it.
//
// So, it will be successfully only if some View ID is already registered
// by Register method and Register method links viewIDEncoded with that some View ID.
func (vidc *tViewIDConverter) Decode(viewIDEncoded tViewIDEncoded) (tViewID, *tViewIDConverterError) {

	if !viewIDEncoded.isValid() {
		return cViewIDNull, eViewIDConverter.InvalidEncodedViewID
	}

	viewID, found := vidc.mDecodeStorage[viewIDEncoded]
	if !found {
		return cViewIDNull, eViewIDConverter.UnregisteredViewID
	}

	if !viewID.isValid() {
		return cViewIDNull, eViewIDConverter.InvalidViewID
	}

	return viewID, nil
}

// Register tries to register View ID (with encoding it) and if it was successful,
// returns an encoded View ID.
//
// So, it will be successfully only if the same View ID is not already registered
// and if limit of registered View IDs is not reached.
func (vidc *tViewIDConverter) Register(viewID tViewID) (tViewIDEncoded, *tViewIDConverterError) {

	if !viewID.isValid() {
		return cViewIDEncodedNull, eViewIDConverter.InvalidViewID
	}

	alreadyEncodedViewID, isAlreadyRegistered := vidc.mEncodeStorage[viewID]
	if isAlreadyRegistered {
		return cViewIDEncodedNull, eViewIDConverter.AlreadyRegisteredViewID.arg(alreadyEncodedViewID)
	}

	if !vidc.encodedIDGenerator.isValid() {
		return cViewIDEncodedNull, eViewIDConverter.LimitReached.arg(vidc.encodedIDGenerator)
	}

	vidc.encodedIDGenerator++
	vidc.mEncodeStorage[viewID] = vidc.encodedIDGenerator
	vidc.mDecodeStorage[vidc.encodedIDGenerator] = viewID

	return vidc.encodedIDGenerator, nil
}

// makeViewIDConverter is the tViewIDConverter constructor.
//
// Allocates memory for internal parts and initializes the default state of
// tViewIDEncoded generator.
func makeViewIDConverter() *tViewIDConverter {
	return &tViewIDConverter{
		mEncodeStorage:     make(map[tViewID]tViewIDEncoded),
		mDecodeStorage:     make(map[tViewIDEncoded]tViewID),
		encodedIDGenerator: cViewIDEncodedStartValue,
	}
}
