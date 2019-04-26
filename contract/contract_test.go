package main

import (
	. "github.com/orbs-network/orbs-contract-sdk/go/testing/unit"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_saveRevision_Single(t *testing.T) {
	caller := AnAddress()

	InServiceScope(nil, caller, func(m Mockery) {
		saveRevision("David Bowie", "The greatest performer of the 20th century")

		bowieFirstRevision := getLastRevision("David Bowie")

		require.EqualValues(t, "David Bowie", bowieFirstRevision.Name)
		require.EqualValues(t, "The greatest performer of the 20th century", bowieFirstRevision.Text)
	})
}

func Test_saveRevision_Multiple(t *testing.T) {
	caller := AnAddress()

	InServiceScope(nil, caller, func(m Mockery) {
		saveRevision("David Bowie", "The greatest performer of the 20th century")

		bowieFirstRevision := getLastRevision("David Bowie")

		require.EqualValues(t, "David Bowie", bowieFirstRevision.Name)
		require.EqualValues(t, "The greatest performer of the 20th century", bowieFirstRevision.Text)


		saveRevision("Iggy Pop", "Another great performer of the 20th century")

		iggyFirstRevision := getLastRevision("Iggy Pop")

		require.EqualValues(t, "Iggy Pop", iggyFirstRevision.Name)
		require.EqualValues(t, "Another great performer of the 20th century", iggyFirstRevision.Text)
	})
}
