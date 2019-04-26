package main

import (
	"encoding/hex"
	"encoding/json"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/v1"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/v1/address"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/v1/state"
	"strconv"
)

var PUBLIC = sdk.Export(GetLastRevision, SaveRevision)
var SYSTEM = sdk.Export(_init)

type Revision struct {
	ID        uint64
	Timestamp uint64

	Author string
	AuthorAddress []byte

	Name   string
	Text   string
}

func _init() {

}

func SaveRevision(name string, text string) (rawJSON string) {
	 saveRevision(name, text)
	 return GetLastRevision(name)
}

func GetLastRevision(name string) (rawJSON string) {
	data, _ := json.Marshal(getLastRevision(name))
	return string(data)
}

func saveRevision(name string, text string) (revisionId uint64) {
	revisions := revisionsList(name)
	revisions.Add(Revision{
		Name: name,
		Text: text,
		AuthorAddress: address.GetCallerAddress(),
	})

	return revisions.Count()
}

func getLastRevision(name string) Revision {
	revisions := revisionsList(name)
	return revisions.Last().(Revision)
}

func getRevisions(name string) []Revision {
	var revisions []Revision
	revisionsList(name).Iterate(func(id uint64, item interface{}) bool {
		revision := item.(Revision)
		revisions = append(revisions, revision)
		return true
	})

	return revisions
}

func serializeRevision(compositeKey []byte, id uint64, params ...interface{}) {
	revision := params[0].(Revision)

	// Save latest revision by that name
	state.WriteUint64([]byte(revision.Name), id)
	state.WriteString(nameKey(compositeKey), revision.Name)
	state.WriteString(textKey(compositeKey), revision.Text)
	state.WriteBytes(authorKey(compositeKey), address.GetCallerAddress())
}

func deserializeRevision(compositeKey []byte, id uint64) interface{} {
	author := state.ReadBytes(authorKey(compositeKey))

	return Revision{
		ID: id,
		Name: state.ReadString(nameKey(compositeKey)),
		Text: state.ReadString(textKey(compositeKey)),
		Author: hex.EncodeToString(author),
		AuthorAddress: author,
	}
}

func revisionsList(name string) List {
	return NewList(name+"_revisions", serializeRevision, deserializeRevision)
}

func textKey(compositeKey []byte) []byte {
	return []byte(string(compositeKey) + "_t")
}

func nameKey(compositeKey []byte) []byte {
	return []byte(string(compositeKey) + "_n")
}

func authorKey(compositeKey []byte) []byte {
	return []byte(string(compositeKey) + "_a")
}

type ListSerializer func(compositeKey []byte, id uint64, params ...interface{})
type ListDeserializer func(compositeKey []byte, id uint64) interface{}

type List interface {
	Count() uint64
	Add(...interface{}) uint64
	Get(id uint64) interface{}
	Last() interface{}
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

func (l list) Last() interface{} {
	return l.Get(l.Count())
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