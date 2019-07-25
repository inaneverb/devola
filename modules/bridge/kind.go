// Copyright © 2019. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package bridge

//
type Kind uint8

//
const (
	KindHandler Kind = 1 + iota

	KindMiddleware

	KindOnSuccessFinisher

	KindOnErrorFinisher
)

//
func (k Kind) String() string {

}
