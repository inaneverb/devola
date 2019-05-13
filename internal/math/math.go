// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package math

// todo: comment

//
func AbsI(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

//
func ClampI(n, min, max int) int {
	if min > max {
		min, max = max, min
	}
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}