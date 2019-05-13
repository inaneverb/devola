// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package chat

//
type CIM struct {
	getter FCIMGetChatInfo
	setter FCIMSaveChatInfo
}

//
type FCIMGetChatInfo func(id ID) (*Info, error)

//
type FCIMSaveChatInfo func(ci *Info) error

//
type FCIMGetSession func()
