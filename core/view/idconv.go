// Copyright © 2019. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package view

import (
	"github.com/qioalice/devola/core/errors"
)

// IDConv is the encoder from readable View ID format (ID)
// to internal (used in SDK) format (IDEnc).
//
// Works using a two maps:
// A view from ID to encoded value (IDEnc) and vice-versa.
// Encode, Decode methods using these maps.
type IDConv struct {

	// When some new View ID registers, the entry of this link
	// View ID -> View ID encoded places here.
	mEncodeStorage map[ID]IDEnc

	// When some new View ID registers, the entry of this link
	// View ID encoded -> View ID places here.
	mDecodeStorage map[IDEnc]ID

	// Generator for encoded View ID.
	// Increases by one for each new View ID by Encode method.
	encodedIDGenerator IDEnc
}

// Encode tries to encode passed View ID.
// Also registers passed View ID if it is not. Returns an encoded ID.
func (idc *IDConv) Encode(id ID) IDEnc {

	if !id.IsValid() {
		return CIDEncNil
	}

	if alreadyRegistered, found := idc.mEncodeStorage[id]; found {
		return alreadyRegistered
	}

	if !idc.encodedIDGenerator.IsValid() {
		return CIDEncNil
	}

	idc.encodedIDGenerator++
	idc.mEncodeStorage[id] = idc.encodedIDGenerator
	idc.mDecodeStorage[idc.encodedIDGenerator] = id

	return idc.encodedIDGenerator
}

// Decode tries to decode passed encoded View ID.
//
// Returns a decoded ID and error code of that operation.
// (See error codes in ec.go file).
//
// So, it will be successfully only if some View ID is already registered
// by Encode method.
func (idc *IDConv) Decode(idenc IDEnc) (ID, errors.Code) {

	if !idenc.IsValid() {
		return CIDNil, ECInvalidIDEnc
	}

	id, found := idc.mDecodeStorage[idenc]
	if !found {
		return CIDNil, ECNotRegistered
	}

	if !id.IsValid() {
		return CIDNil, ECInvalidID
	}

	return id, errors.ECOK
}

// MakeIDConv is the IDConv constructor.
//
// Allocates memory for internal parts and initializes the default state of
// IDEnc generator.
func MakeIDConv() *IDConv {
	return &IDConv{
		mEncodeStorage:     make(map[ID]IDEnc),
		mDecodeStorage:     make(map[IDEnc]ID),
		encodedIDGenerator: cIDEncStartValue,
	}
}
