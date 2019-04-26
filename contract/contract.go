package main

import (
	"encoding/hex"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/v1"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/v1/address"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/v1/state"
	"strconv"
)

var PUBLIC = sdk.Export()
var SYSTEM = sdk.Export(_init)

type Revision struct {
	ID        uint64
	Timestamp uint64

	Author string
	Name   string
	Text   string
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
	return getRevisionById(revisionId)
}

func getRevisionById(revisionId uint64) Revision {
	return Revision{
		ID:     revisionId,
		Author: hex.EncodeToString(state.ReadBytes(authorKey(revisionId))),
		Name:   state.ReadString(nameKey(revisionId)),
		Text:   state.ReadString(textKey(revisionId)),
	}
}

func getRevisions(name string) []Revision {
	return nil
}

func newRevision(name string) (revisionId uint64) {
	revisionId = state.ReadUint64([]byte(REVISIONS_COUNTER_KEY)) + 1
	state.WriteUint64([]byte(REVISIONS_COUNTER_KEY), revisionId)

	// Save latest revision by that name
	state.WriteUint64([]byte(name), revisionId)
	state.WriteString(nameKey(revisionId), name)
	state.WriteBytes(authorKey(revisionId), address.GetCallerAddress())

	return
}

func key(revisionId uint64, postfix string) []byte {
	return []byte(int10(revisionId) + "_" + postfix)
}

func textKey(revisionId uint64) []byte {
	return key(revisionId, "t")
}

func nameKey(revisionId uint64) []byte {
	return key(revisionId, "n")
}

func lastRevisionId(name string) (revisionId uint64) {
	return state.ReadUint64([]byte(name))
}

func authorKey(revisionId uint64) []byte {
	return key(revisionId, "a")
}

type ListSerializer func(compositeKey []byte, id uint64, params ...interface{})
type ListDeserializer func(compositeKey []byte, id uint64) interface{}

type List interface {
	Count() uint64
	Add(...interface{}) uint64
	Get(id uint64) interface{}
	// stop iterating if false
	Iterate(func (id uint64, item interface{}) bool)
}

type list struct {
	prefix       string
	serializer   ListSerializer
	deserializer ListDeserializer
}

func NewList(prefix string, serializer ListSerializer, deserializer ListDeserializer) List {
	return list{prefix, serializer, deserializer}
}

func (l list) Count() uint64 {
	return state.ReadUint64([]byte(l.prefix + "_counter"))
}

func (l list) Add(params ...interface{}) uint64 {
	count := l.Count() + 1
	state.WriteUint64([]byte(l.prefix+"_counter"), count)
	l.serializer([]byte(l.prefix+"_"+int10(count)), count, params...)
	return count
}

func (l list) Get(id uint64) interface{} {
	return l.deserializer([]byte(l.prefix+"_"+int10(id)), id)
}

func (l list) Iterate(f func (id uint64, item interface{}) bool) {
	for i := uint64(1); i <= l.Count(); i++ {
		if !f(i, l.Get(i)) {
			break
		}
	}
}

func int10(i uint64) string {
	return strconv.FormatUint(i, 10)
}