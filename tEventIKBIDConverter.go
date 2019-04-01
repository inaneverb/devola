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

import (
	"errors"
	"fmt"
	"math"
)

// tEventIKBIDEncoded is the internal type that technically is just
// an encoded part of Telegram Bot Inline Keyboard Button identifier.
//
// It's done for saving readability of tEvent.Data but also optimise an usage
// of allowed len of Telegram Inline Keyboard Button Callback Data.
//
// At this moment (April, 2019) Telegram API (v4.1) allow represent
// action of Inline Keyboard Buttons as some string that must have length
// not more than 64 byte.
//
// As a result, each string-represented Inline Keyboard Button action,
// registered by this SDK, replaces by some tEventIKBIDEncoded value
// (it is a integer type) and then these bytes will be used as part of
// inline keyboard button identifier: The Callback Data.
type tEventIKBIDEncoded uint32

// Internal constants of encoded Inline Keyboard Button identifier.
const (
	// Internal constant.
	// Represents a nil encoded identifier and is a indicator of some error.
	cEventIKBDataEncodedNull tEventIKBIDEncoded = 0

	// Internal constant.
	// All encoded identifiers as numbers will be more than or equal to
	// that value. All less values are reserved for internal needs.
	cEventIKBDataEncodedStartValue tEventIKBIDEncoded = 100

	// Internal constant.
	// All encoded identifiers as numbers will be less than that value.
	// All more than or equal to that values are reserver for internal needs.
	cEventIKBDataEncodedOverMaxValue tEventIKBIDEncoded = math.MaxUint32 - 1
)

// tEventIKBIDConverter is the encoder from readable string-like
// Inline Keyboard Button identifier to the its number representation
// encoded value.
//
// Works using a two maps:
// A view from string to encoded value and vice-versa.
// Encode, Decode methods using these maps and Register method just add
// entries in these maps.
type tEventIKBIDConverter struct {
	mEncodeStorage map[string]tEventIKBIDEncoded
	mDecodeStorage map[tEventIKBIDEncoded]string
	idGenerator    tEventIKBIDEncoded
}

// isValidAction returns true only for strings (inline keyboard actions)
// that contains more than 2 characters and not starts from doubleunderscore.
// Doubleunderscore starting is reserved for internal needs.
func (*tEventIKBIDConverter) isValidAction(action string) bool {
	return len(action) > 2 && action[:2] != "__"
}

// isValidEncodedData returns true only for encoded inline keyboard actions
// that are not equal to the cEventIKBDataEncodedNull, and is in the
// allowed range.
func (*tEventIKBIDConverter) isValidEncodedData(
	encodedData tEventIKBIDEncoded,
) bool {
	return encodedData != cEventIKBDataEncodedNull &&
		encodedData >= cEventIKBDataEncodedStartValue &&
		encodedData < cEventIKBDataEncodedOverMaxValue
}

// Encode tries to represent and return an inline keyboard button identifier
// as some tEventIKBIDEncoded value.
// So, it will be successfully only if action is already registered
// by Register method.
func (eikbdc *tEventIKBIDConverter) Encode(
	action string,
) (
	tEventIKBIDEncoded,
	error,
) {
	if !eikbdc.isValidAction(action) {
		err := "EventIKBDataConverter: Trying to encode an invalid action."
		return cEventIKBDataEncodedNull, errors.New(err)
	}
	encodedData, found := eikbdc.mEncodeStorage[action]
	if !found {
		err := "EventIKBDataConverter: Trying to encode an unregistered action."
		return cEventIKBDataEncodedNull, errors.New(err)
	}
	if !eikbdc.isValidEncodedData(encodedData) {
		err := "EventIKBDataConverter: Encoded data is invalid. Internal error."
		return cEventIKBDataEncodedNull, errors.New(err)
	}
	return encodedData, nil
}

// Decode tries to represent encoded an inline keyboard button identifier
// as readable string.
// So, it will be successfully only if some action is already registered
// by Register method and Register method links encodedData with action.
func (eikbdc *tEventIKBIDConverter) Decode(
	encodedData tEventIKBIDEncoded,
) (
	string,
	error,
) {
	if !eikbdc.isValidEncodedData(encodedData) {
		err := "EventIKBDataConverter: Trying to decode an invalid encoded data."
		return "", errors.New(err)
	}
	action, found := eikbdc.mDecodeStorage[encodedData]
	if !found {
		err := "EventIKBDataConverter: " +
			"No registered action for passed encoded data."
		return "", errors.New(err)
	}
	if !eikbdc.isValidAction(action) {
		err := "EventIKBDataConverter: " +
			"Decoded action is invalid. Internal error."
		return "", errors.New(err)
	}
	return action, nil
}

// Register tries to register action as new Inline Keyboard Button readable
// identifier, provide a some tEventIKBIDEncoded identifier for this,
// link and associate it.
// So, it will be successfully only if the same action is not already registered
// and if limit of registered action is not reached.
// In this case encoded identifer and nil as error are returned.
func (eikbdc *tEventIKBIDConverter) Register(
	action string,
) (
	tEventIKBIDEncoded,
	error,
) {
	if !eikbdc.isValidAction(action) {
		err := "EventIKBDataConverter: Trying to register an invalid action."
		return cEventIKBDataEncodedNull, errors.New(err)
	}
	alreadyEncodedData, isAlreadyRegistered := eikbdc.mEncodeStorage[action]
	if isAlreadyRegistered {
		err := "EventIKBDataConverter: Trying to register already registered " +
			"(encoded to the: \"%d\") action."
		return cEventIKBDataEncodedNull, fmt.Errorf(err, alreadyEncodedData)
	}
	if !eikbdc.isValidEncodedData(eikbdc.idGenerator) {
		err := "EventIKBDataConverter: Action's limit to register is reached."
		return cEventIKBDataEncodedNull, errors.New(err)
	}
	eikbdc.idGenerator++
	eikbdc.mEncodeStorage[action] = eikbdc.idGenerator
	eikbdc.mDecodeStorage[eikbdc.idGenerator] = action
	return eikbdc.idGenerator, nil
}

// makeEventIKBDataConverter is the tEventIKBIDConverter constructor.
// Allocates memory for internal parts and initializes the default state of
// tEventIKBIDEncoded generator.
func makeEventIKBDataConverter() *tEventIKBIDConverter {
	return &tEventIKBIDConverter{
		mEncodeStorage: make(map[string]tEventIKBIDEncoded),
		mDecodeStorage: make(map[tEventIKBIDEncoded]string),
		idGenerator:    cEventIKBDataEncodedNull,
	}
}
