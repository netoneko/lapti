package main

import (
	"encoding/hex"
	"fmt"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/v1/state"
	. "github.com/orbs-network/orbs-contract-sdk/go/testing/unit"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_saveRevision_Single(t *testing.T) {
	caller := AnAddress()

	InServiceScope(nil, caller, func(m Mockery) {
		saveRevision("David Bowie", "The greatest performer of the 20th century")

		bowieFirstRevision := getLastRevision("David Bowie")

		require.EqualValues(t, 1, bowieFirstRevision.ID)
		require.EqualValues(t, "David Bowie", bowieFirstRevision.Name)
		require.EqualValues(t, "The greatest performer of the 20th century", bowieFirstRevision.Text)
		require.EqualValues(t, hex.EncodeToString(caller), bowieFirstRevision.Author)
	})
}

func Test_saveRevision_Multiple(t *testing.T) {
	caller := AnAddress()

	InServiceScope(nil, caller, func(m Mockery) {
		saveRevision("David Bowie", "The greatest performer of the 20th century")

		bowieFirstRevision := getLastRevision("David Bowie")

		require.EqualValues(t, 1, bowieFirstRevision.ID)
		require.EqualValues(t, "David Bowie", bowieFirstRevision.Name)
		require.EqualValues(t, "The greatest performer of the 20th century", bowieFirstRevision.Text)
		require.EqualValues(t, hex.EncodeToString(caller), bowieFirstRevision.Author)

		saveRevision("Iggy Pop", "Another great performer of the 20th century")

		iggyFirstRevision := getLastRevision("Iggy Pop")

		require.EqualValues(t, 2, iggyFirstRevision.ID)
		require.EqualValues(t, "Iggy Pop", iggyFirstRevision.Name)
		require.EqualValues(t, "Another great performer of the 20th century", iggyFirstRevision.Text)
		require.EqualValues(t, hex.EncodeToString(caller), iggyFirstRevision.Author)
	})
}

func Test_listRevisions(t *testing.T) {
	t.Skip()
	caller := AnAddress()

	InServiceScope(nil, caller, func(m Mockery) {
		saveRevision("David Bowie", "The greatest performer of the 20th century")
		saveRevision("David Bowie", "The singer of the 20th century")

		revisions := getRevisions("David Bowie")
		require.Equal(t, 1, revisions[0].ID)
		require.Equal(t, 2, revisions[0].ID)
	})
}

func TestList(t *testing.T) {
	caller := AnAddress()

	InServiceScope(nil, caller, func(m Mockery) {
		s := func(id uint64, params ...interface{}) {
			state.WriteString([]byte(fmt.Sprintf("h_%d", id)), params[0].(string))
		}

		d := func(id uint64) interface{} {
			return state.ReadString([]byte(fmt.Sprintf("h_%d", id)))
		}

		l := NewList("some_list", s, d)
		l.Add("hello!")

		require.EqualValues(t, 1, l.Count())

		item := l.Get(1)

		require.EqualValues(t, "hello!", item)
	})
}
