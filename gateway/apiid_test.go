package gateway

import (
	"testing"

	"github.com/42wim/matterbridge/bridge"
	"github.com/42wim/matterbridge/bridge/config"
	lru "github.com/hashicorp/golang-lru"
	"github.com/stretchr/testify/assert"
)

// The API bridge's Send() always returns an empty message ID, so handleMessage
// never records a BrMsgID for it. These tests pin down what that means for
// anything trying to address a message through the API.
func TestGetDestMsgIDForAPIDest(t *testing.T) {
	cache, err := lru.New(100)
	assert.NoError(t, err)

	whatsapp := &bridge.Bridge{Protocol: "whatsapp", Name: "wa"}
	apiBr := &bridge.Bridge{Protocol: "api", Name: "local"}

	// What the router caches for an inbound whatsapp message relayed onward:
	// entries for every bridge that returned an ID. The API never does.
	cache.Add("whatsapp 3EB0ABC", []*BrMsgID{
		{br: whatsapp, ID: "whatsapp 3EB0ABC", ChannelID: "groupwhatsapp.wa"},
	})

	gw := &Gateway{Messages: cache}
	channel := &config.ChannelInfo{ID: "apiapi.local"}

	got := gw.getDestMsgID("whatsapp 3EB0ABC", apiBr, channel)

	// No "api" entry exists, so the API bridge is handed an empty ID —
	// meaning every message the API bridge receives arrives with id == "".
	assert.Equal(t, "", got)
}

func TestGetDestMsgIDForAPIPostedMessage(t *testing.T) {
	cache, err := lru.New(100)
	assert.NoError(t, err)

	whatsapp := &bridge.Bridge{Protocol: "whatsapp", Name: "wa"}

	// A message the bot POSTed through the API now carries a generated ID,
	// so the router caches it under an "api ..." key with the downstream IDs.
	cache.Add("api d1n2o3p4", []*BrMsgID{
		{br: whatsapp, ID: "whatsapp 3EB0XYZ", ChannelID: "groupwhatsapp.wa"},
	})

	gw := &Gateway{Messages: cache}
	channel := &config.ChannelInfo{ID: "groupwhatsapp.wa"}

	got := gw.getDestMsgID("api d1n2o3p4", whatsapp, channel)

	// Deleting a message the bot itself sent does resolve to the native ID.
	assert.Equal(t, "3EB0XYZ", got)
}
