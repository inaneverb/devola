// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package errors

// Is is a error filter.
// Returns a first Error object from candidates that the same as e.
func Is(e error, candidates ...Error) (matched Error) {
	return Is2(e, candidates)
}

// Is2 is the same as Is but takes candidates directly as slice.
// You can use this func when you need to avoid unnecessary copying slice you already have.
func Is2(e error, candidates []Error) (matched Error) {

	if e == nil || len(candidates) == 0 {
		return nil
	}

	if _, ok := e.(Error); !ok {
		return nil
	}

	for _, candidate := range candidates {
		if candidate.IsIt(e) {
			return candidate
		}
	}

	return nil
}