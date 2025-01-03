// Copyright 2024 New Vector Ltd.
// Copyright 2022 The Matrix.org Foundation C.I.C.
//
// SPDX-License-Identifier: AGPL-3.0-only OR LicenseRef-Element-Commercial
// Please see LICENSE files in the repository root for full details.

package gobind

import (
	"strings"
	"testing"

	"github.com/matrix-org/gomatrixserverlib/spec"
)

func TestMonolithStarts(t *testing.T) {
	monolith := DendriteMonolith{
		StorageDirectory: t.TempDir(),
		CacheDirectory:   t.TempDir(),
	}
	monolith.Start()
	monolith.PublicKey()
	monolith.Stop()
}

func TestMonolithSetRelayServers(t *testing.T) {
	testCases := []struct {
		name           string
		nodeID         string
		relays         string
		expectedRelays string
		expectSelf     bool
	}{
		{
			name:           "assorted valid, invalid, empty & self keys",
			nodeID:         "@valid:abcdef123456abcdef123456abcdef123456abcdef123456abcdef123456abcd",
			relays:         "@valid:123456123456abcdef123456abcdef123456abcdef123456abcdef123456abcd,@invalid:notakey,,",
			expectedRelays: "123456123456abcdef123456abcdef123456abcdef123456abcdef123456abcd",
			expectSelf:     true,
		},
		{
			name:           "invalid node key",
			nodeID:         "@invalid:notakey",
			relays:         "@valid:123456123456abcdef123456abcdef123456abcdef123456abcdef123456abcd,@invalid:notakey,,",
			expectedRelays: "",
			expectSelf:     false,
		},
		{
			name:           "node is self",
			nodeID:         "self",
			relays:         "@valid:123456123456abcdef123456abcdef123456abcdef123456abcdef123456abcd,@invalid:notakey,,",
			expectedRelays: "123456123456abcdef123456abcdef123456abcdef123456abcdef123456abcd",
			expectSelf:     false,
		},
	}

	for _, tc := range testCases {
		monolith := DendriteMonolith{
			StorageDirectory: t.TempDir(),
			CacheDirectory:   t.TempDir(),
		}
		monolith.Start()

		inputRelays := tc.relays
		expectedRelays := tc.expectedRelays
		if tc.expectSelf {
			inputRelays += "," + monolith.PublicKey()
			expectedRelays += "," + monolith.PublicKey()
		}
		nodeID := tc.nodeID
		if nodeID == "self" {
			nodeID = monolith.PublicKey()
		}

		monolith.SetRelayServers(nodeID, inputRelays)
		relays := monolith.GetRelayServers(nodeID)
		monolith.Stop()

		if !containSameKeys(strings.Split(relays, ","), strings.Split(expectedRelays, ",")) {
			t.Fatalf("%s: expected %s got %s", tc.name, expectedRelays, relays)
		}
	}
}

func containSameKeys(expected []string, actual []string) bool {
	if len(expected) != len(actual) {
		return false
	}

	for _, expectedKey := range expected {
		hasMatch := false
		for _, actualKey := range actual {
			if actualKey == expectedKey {
				hasMatch = true
			}
		}

		if !hasMatch {
			return false
		}
	}

	return true
}

func TestParseServerKey(t *testing.T) {
	testCases := []struct {
		name        string
		serverKey   string
		expectedErr bool
		expectedKey spec.ServerName
	}{
		{
			name:        "valid userid as key",
			serverKey:   "@valid:abcdef123456abcdef123456abcdef123456abcdef123456abcdef123456abcd",
			expectedErr: false,
			expectedKey: "abcdef123456abcdef123456abcdef123456abcdef123456abcdef123456abcd",
		},
		{
			name:        "valid key",
			serverKey:   "abcdef123456abcdef123456abcdef123456abcdef123456abcdef123456abcd",
			expectedErr: false,
			expectedKey: "abcdef123456abcdef123456abcdef123456abcdef123456abcdef123456abcd",
		},
		{
			name:        "invalid userid key",
			serverKey:   "@invalid:notakey",
			expectedErr: true,
			expectedKey: "",
		},
		{
			name:        "invalid key",
			serverKey:   "@invalid:notakey",
			expectedErr: true,
			expectedKey: "",
		},
	}

	for _, tc := range testCases {
		key, err := getServerKeyFromString(tc.serverKey)
		if tc.expectedErr && err == nil {
			t.Fatalf("%s: expected an error", tc.name)
		} else if !tc.expectedErr && err != nil {
			t.Fatalf("%s: didn't expect an error: %s", tc.name, err.Error())
		}
		if tc.expectedKey != key {
			t.Fatalf("%s: keys not equal. expected: %s got: %s", tc.name, tc.expectedKey, key)
		}
	}
}
