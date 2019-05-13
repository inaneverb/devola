// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package chat

import (
	api "github.com/go-telegram-bot-api/telegram-bot-api"
	//"github.com/qioalice/i18n"
)

// todo: Add flag if bot is locked in chat
type Info struct {
	*api.Chat             `json:",inline"`
	StartedUnixstamp      int64 `json:"started_unixstamp"`
	LastActivityUnixstamp int64 `json:"last_activity_unixstamp"`
	//UsedLocale i18n.LocaleName
	currentSSID tSessionID `json:"current_ssid"`
}

//
func (ci *Info) TrFinish() error {
	panic("implement me")
}
