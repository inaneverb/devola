// Copyright © 2019. All rights reserved.
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
type FCIMGetChatInfo func(id ID) (*Chat, error)

//
type FCIMSaveChatInfo func(ci *Chat) error

//
type FCIMGetSession func()
