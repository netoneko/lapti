package main

import (
	"strconv"

	"github.com/orbs-network/orbs-contract-sdk/go/sdk/v1"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/v1/state"
)

var PUBLIC = sdk.Export(saveRevision, getLastRevision)
var SYSTEM = sdk.Export(_init)

type Revision struct {
	ID          uint64
	Timestamp   uint64

	Author      string
	Name string
	Text     string
}

func _init() {

}

const REVISIONS_COUNTER_KEY = "revisions_counter"

func saveRevision(name string, text string) (messageID uint64) {
	revisionId := newRevision(name)
	state.WriteString(textKey(revisionId), text)

	return
}

func getLastRevision(name string) Revision {
	revisionId := lastRevisionId(name)
	return Revision{
		ID: revisionId,
		Name: name,
		Text: state.ReadString(textKey(revisionId)),
	}
}


func key(revisionId uint64, postfix string) []byte {
	return []byte(strconv.FormatUint(revisionId, 10) + "_" + postfix)
}

func newRevision(name string) (revisionId uint64) {
	revisionId = state.ReadUint64([]byte(REVISIONS_COUNTER_KEY)) + 1
	state.WriteUint64([]byte(REVISIONS_COUNTER_KEY), revisionId)

	// Save latest revision by that name
	state.WriteUint64([]byte(name), revisionId)

	return
}

func textKey(revisionId uint64) []byte {
	return key(revisionId, "t")
}

func lastRevisionId(name string) (revisionId uint64) {
	return state.ReadUint64([]byte(name))
}