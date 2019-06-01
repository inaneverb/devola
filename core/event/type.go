// Copyright © 2019. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package event

import (
	"strings"
	"sync"
)

// Type represents type of event that is occurred.
// This is a part of Event object, which is a part of backend Ctx.
//
// Thus you can always to figure out inside handler which kind of event
// your handler handling. Use predefined constants, presented below for this.
type Type uint8

// Constants of Type.
// Use these constants to figure out what kind of event is occurred
// (by comparing Type).
const (

	// Marker of invalid type.
	CTypeInvalid Type = 0
)

// String returns a string representation of type.
// It uses atoi method if no registered alias detected.
func (t Type) String() string {

	if t == CTypeInvalid {
		return "Invalid type"
	}

	typeAliases.mu.RLock()
	defer typeAliases.mu.RUnlock()

	types, ok := typeAliases.m[t]
	if !ok {
		return "Invalid type"
	}

	return strings.Join(types, ", ")
}

// typeAliases is a map that allows to set comments to each type by TypeComment func
// and then get them by String method.
var typeAliases = struct {
	m  map[Type][]string
	mu sync.RWMutex
}{
	m: make(map[Type][]string),
}

// TypeComment creates a comment for type t that allows to get that comment
// by String method.
func TypeComment(t Type, comment string) {

	if t == CTypeInvalid || comment == "" {
		return
	}

	typeAliases.mu.Lock()
	defer typeAliases.mu.Unlock()

	v := typeAliases.m[t]

	// Don't save the same comments
	for _, vv := range v {
		if vv == comment {
			return
		}
	}

	typeAliases.m[t] = append(v, comment)
}
