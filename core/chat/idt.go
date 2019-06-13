// Copyright © 2019. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package chat

// ID represents the type of backend chat's ID (and only ID not a chat's type).
type ID int64

// Type represents the type of backend chat's type.
//
// The number of types depended on backend internal struct
// and it can be determined up to MaxTypeValue chat types.
type Type uint8

// IDT is the combination of ID and Type of chat.
// This is the way to store ID and type as one entity.
//
// At this moment IDT is an alias to 64 bit variable and
// allows to use up to 52 bits to storing chat's ID value  +1 signed bit (53 total),
// and allows to use up to 4 bits to storing chat's Type.
//
// internal view of IDT:
//
//    |+|    |       |                52 bits of chat's ID                |
//     ↑  ↑      ↑
//     |  |      7 reserved bits
//     |  4 bits of chat's Type
//     1 sign bit of chat's ID
//
type IDT uint64

// MaxTypeValue represents the upper bound value of chat's Type.
// Backend can determine not more than MaxTypeValue types.
//
// WARNING!
// The chat's type Type value more than MaxTypeValue will be considered incorrect
// and will cause UB.
const MaxTypeValue Type = 15

// Constants to perform bitwise operation of extracting/storing ID/Type subvalues
// from/in IDT.
const (
	idtMaskID       IDT = 0x800FFFFFFFFFFFFF // Mask of ID (high sign bit + low 52 bits)
	idtMaskType     IDT = 0x7800000000000000 // Mask of Type (high 8 bits after one)
	idtMaskReserved IDT = 0x07F0000000000000 // Mask of reserved bits (high 3 bits after 8+1 bits)

	idtTypeSHR2uint8 uint8 = 59 // SHR IDT to make Type subvalue as uint8
)

// ID extracts and returns ID from IDT - a combination of ID and Type.
func (idt IDT) ID() ID {
	return ID(idt & idtMaskID)
}

// Type extracts and returns Type from IDT - a combination of ID and Type.
func (idt IDT) Type() Type {
	return Type((idt & idtMaskType) >> idtTypeSHR2uint8)
}

// NewIDT creates a new IDT value by combining passed chat's ID and chat's Type.
func NewIDT(id ID, typ Type) IDT {
	return (IDT(typ) << idtTypeSHR2uint8) & IDT(id)
}
