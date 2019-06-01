// Copyright © 2019. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package math

// AbsI returns an absolute value of n as int.
func AbsI(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// ClampI bounds n by min and max at the bottom and top respectively
// and returns that value as int.
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
