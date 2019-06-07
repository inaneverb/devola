// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package lirester

import (
	"github.com/qioalice/devola/core/chat"
)

// cleanupConfig represents an object that contains two important things:
// - Chat id of chat, cleanup operation will perform over
// - Unixnano timestamp when cleanup operation will perform
type cleanupConfig struct {
	chatID chat.ID
	when   int64
}

// makeCleanupConfig creates a new cleanupConfig object with
// passed chat id and timestamp when created cleanup rule should be applied.
func makeCleanupConfig(chatID chat.ID, when int64) cleanupConfig {
	return cleanupConfig{
		chatID: chatID,
		when:   when,
	}
}
