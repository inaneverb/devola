// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package view

import (
	"../errors"
)

// IDConverter is the encoder from readable View ID format (ID)
// to internal (used in SDK) format (IDEnc).
//
// Works using a two maps:
// A view from ID to encoded value (IDEnc) and vice-versa.
// Encode, Decode methods using these maps and Register method just to add
// entries in these maps.
type IDConverter struct {

	// When some new View ID registers, the entry of this link
	// View ID -> View ID encoded places here.
	mEncodeStorage map[ID]IDEnc

	// When some new View ID registers, the entry of this link
	// View ID encoded -> View ID places here.
	mDecodeStorage map[IDEnc]ID

	// Generator for encoded View ID.
	// Increases by one for each new View ID by Register method.
	encodedIDGenerator IDEnc
}

// Encode tries to encode passed View ID.
// Returns an encoded ID and error code of that operation.
// (See error codes in error_codes.go file).
//
// So, it will be successfully only if passed id is already registered
// by Register method and all internal checks are passed.
func (idc *IDConverter) Encode(id ID) (IDEnc, errors.Code) {

	if !id.IsValid() {
		return CIDEncNil, EInvalidID
	}

	idenc, found := idc.mEncodeStorage[id]
	if !found {
		return CIDEncNil, ENotRegistered
	}

	if !idenc.IsValid() {
		return CIDEncNil, EInvalidIDEnc
	}

	return idenc, EOK
}

// Decode tries to decode passed encoded View ID.
// Returns a decoded ID and error code of that operation.
// (See error codes in error_codes.go file).
//
// So, it will be successfully only if some View ID is already registered
// by Register method and Register method links idenc with that some View ID.
func (idc *IDConverter) Decode(idenc IDEnc) (ID, errors.Code) {

	if !idenc.IsValid() {
		return CIDNil, EInvalidIDEnc
	}

	id, found := idc.mDecodeStorage[idenc]
	if !found {
		return CIDNil, ENotRegistered
	}

	if !id.IsValid() {
		return CIDNil, EInvalidID
	}

	return id, EOK
}

// Register tries to register View ID (with encoding it).
// Returns an encoded View ID and error code of registering operation.
// (See error codes in error_codes.go file).
//
// So, it will be successfully only if the same View ID is not already registered
// and if limit of registered View IDs is not reached.
func (idc *IDConverter) Register(id ID) (IDEnc, errors.Code) {

	if !id.IsValid() {
		return CIDEncNil, EInvalidID
	}

	if _, found := idc.mEncodeStorage[id]; found {
		return CIDEncNil, EAlreadyRegistered
	}

	if !idc.encodedIDGenerator.IsValid() {
		return CIDEncNil, ELimitReached
	}

	idc.encodedIDGenerator++
	idc.mEncodeStorage[id] = idc.encodedIDGenerator
	idc.mDecodeStorage[idc.encodedIDGenerator] = id

	return idc.encodedIDGenerator, EOK
}

// MakeIDConverter is the IDConverter constructor.
//
// Allocates memory for internal parts and initializes the default state of
// IDEnc generator.
func MakeIDConverter() *IDConverter {
	return &IDConverter{
		mEncodeStorage:     make(map[ID]IDEnc),
		mDecodeStorage:     make(map[IDEnc]ID),
		encodedIDGenerator: CIDEncStartValue,
	}
}
