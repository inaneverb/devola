// Copyright © 2019. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package sys

import (
	"unsafe"
)

// GoInterface represents what "interface{}" means in internal Golang parts.
type GoInterface struct {
	Type uintptr        // pointer to the type definition struct
	Word unsafe.Pointer // pointer to the value
}
