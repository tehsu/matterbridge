package api

import (
	"testing"

	"github.com/42wim/matterbridge/bridge"
	"github.com/42wim/matterbridge/bridge/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func testAPI() *API {
	logger := logrus.New()
	logger.SetOutput(nil)
	return &API{
		Config: &bridge.Config{
			Bridge: &bridge.Bridge{
				Account: "api.local",
				Log:     logrus.NewEntry(logger),
			},
		},
	}
}

func TestStampMessageAssignsID(t *testing.T) {
	b := testAPI()

	msg := config.Message{Text: "hello", Username: "bot"}
	b.stampMessage(&msg)

	assert.NotEmpty(t, msg.ID, "a posted message needs an ID so the gateway caches it")
	assert.Equal(t, "api", msg.Channel)
	assert.Equal(t, "api", msg.Protocol)
	assert.Equal(t, "api.local", msg.Account)
	assert.False(t, msg.Timestamp.IsZero())
}

func TestStampMessageIDsAreUnique(t *testing.T) {
	b := testAPI()

	first := config.Message{Text: "one"}
	second := config.Message{Text: "two"}
	b.stampMessage(&first)
	b.stampMessage(&second)

	assert.NotEqual(t, first.ID, second.ID)
}

func TestStampMessageKeepsDeleteID(t *testing.T) {
	b := testAPI()

	msg := config.Message{Event: config.EventMsgDelete, ID: "abc123"}
	b.stampMessage(&msg)

	// the ID names the message to remove, so it must survive untouched
	assert.Equal(t, "abc123", msg.ID)
	// an empty text would be dropped by the gateway before reaching a bridge
	assert.Equal(t, config.EventMsgDelete, msg.Text)
}

func TestStampMessageKeepsDeleteText(t *testing.T) {
	b := testAPI()

	msg := config.Message{Event: config.EventMsgDelete, ID: "abc123", Text: "custom"}
	b.stampMessage(&msg)

	assert.Equal(t, "custom", msg.Text)
}

func TestStampMessageKeepsParentID(t *testing.T) {
	b := testAPI()

	msg := config.Message{Text: "a reply", ParentID: "parent-1"}
	b.stampMessage(&msg)

	assert.Equal(t, "parent-1", msg.ParentID)
	assert.NotEmpty(t, msg.ID)
}
