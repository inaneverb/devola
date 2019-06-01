// Copyright © 2019. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package chat

//
type Chat struct {
	StartedUnixstamp      int64 `json:"started_unixstamp"`
	LastActivityUnixstamp int64 `json:"last_activity_unixstamp"`
	// currentSSID SessionID `json:"current_ssid"`
}
