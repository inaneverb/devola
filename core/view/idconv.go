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
// Encode, Decode methods using these maps and Register method just to add
// entries in these maps.
type IDConv struct {

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
// (See error codes in ec.go file).
//
// So, it will be successfully only if passed id is already registered
// by Register method and all internal checks are passed.
func (idc *IDConv) Encode(id ID) (IDEnc, errors.Code) {

	if !id.IsValid() {
		return CIDEncNil, ECInvalidID
	}

	idenc, found := idc.mEncodeStorage[id]
	if !found {
		return CIDEncNil, ECNotRegistered
	}

	if !idenc.IsValid() {
		return CIDEncNil, ECInvalidIDEnc
	}

	return idenc, errors.ECOK
}

// Decode tries to decode passed encoded View ID.
// Returns a decoded ID and error code of that operation.
// (See error codes in ec.go file).
//
// So, it will be successfully only if some View ID is already registered
// by Register method and Register method links idenc with that some View ID.
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

// Register tries to register View ID (with encoding it).
// Returns an encoded View ID and error code of registering operation.
// (See error codes in ec.go file).
//
// So, it will be successfully only if the same View ID is not already registered
// and if limit of registered View IDs is not reached.
func (idc *IDConv) Register(id ID) (IDEnc, errors.Code) {

	if !id.IsValid() {
		return CIDEncNil, ECInvalidID
	}

	if _, found := idc.mEncodeStorage[id]; found {
		return CIDEncNil, ECAlreadyRegistered
	}

	if !idc.encodedIDGenerator.IsValid() {
		return CIDEncNil, ECLimitReached
	}

	idc.encodedIDGenerator++
	idc.mEncodeStorage[id] = idc.encodedIDGenerator
	idc.mDecodeStorage[idc.encodedIDGenerator] = id

	return idc.encodedIDGenerator, errors.ECOK
}

// MakeIDConv is the IDConv constructor.
//
// Allocates memory for internal parts and initializes the default state of
// IDEnc generator.
func MakeIDConv() *IDConv {
	return &IDConv{
		mEncodeStorage:     make(map[ID]IDEnc),
		mDecodeStorage:     make(map[IDEnc]ID),
		encodedIDGenerator: CIDEncStartValue,
	}
}
