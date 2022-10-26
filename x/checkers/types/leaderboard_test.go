package types_test

import (
	"errors"
	"github.com/alice/checkers/x/checkers/types"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestSortStringifiedWinners(t *testing.T) {
	tests := []struct {
		name     string
		unsorted []types.WinningPlayer
		sorted   []types.WinningPlayer
		err      error
	}{
		{
			name: "cannot parse date",
			unsorted: []types.WinningPlayer{
				{
					PlayerAddress: "alice",
					WonCount:      2,
					DateAdded:     "200T-01-02 15:05:05.999999999 +0000 UTC",
				},
			},
			sorted: []types.WinningPlayer{},
			err:    errors.New("date added cannot be parsed: 200T-01-02 15:05:05.999999999 +0000 UTC: parsing time \"200T-01-02 15:05:05.999999999 +0000 UTC\" as \"2006-01-02 15:04:05.999999999 +0000 UTC\": cannot parse \"-01-02 15:05:05.999999999 +0000 UTC\" as \"2006\""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leaderboard := types.Leaderboard{
				Winners: tt.unsorted,
			}
			parsed, err := leaderboard.ParseWinners()
			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
			}
			types.SortWinners(parsed)
			sorted := types.StringifyWinners(parsed)
			require.Equal(t, len(tt.sorted), len(sorted))
			require.EqualValues(t, tt.sorted, sorted)
		})
	}
}

func TestUpdatePlayerInfoAtNow(t *testing.T) {
	tests := []struct {
		name      string
		sorted    []types.WinningPlayer
		candidate types.PlayerInfo
		now       string
		expected  []types.WinningPlayer
	}{
		{
			name:   "add to empty",
			sorted: []types.WinningPlayer{},
			candidate: types.PlayerInfo{
				Index:    "alice",
				WonCount: 2,
			},
			now: "2006-01-02 15:05:05.999999999 +0000 UTC",
			expected: []types.WinningPlayer{
				{
					PlayerAddress: "alice",
					WonCount:      2,
					DateAdded:     "2006-01-02 15:05:05.999999999 +0000 UTC",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now, err := types.ParseDateAddedAsTime(tt.now)
			require.NoError(t, err)
			leaderboard := types.Leaderboard{
				Winners: tt.sorted,
			}
			err = leaderboard.UpdatePlayerInfoAtNow(now, tt.candidate)
			require.NoError(t, err)
			require.Equal(t, len(tt.expected), len(leaderboard.Winners))
			require.EqualValues(t, tt.expected, leaderboard.Winners)
			require.NoError(t, leaderboard.Validate())
		})
	}
}

func makeMaxLengthSortedWinningPlayers() []types.WinningPlayer {
	sorted := make([]types.WinningPlayer, 100)
	for i := uint64(0); i < 100; i++ {
		sorted[i] = types.WinningPlayer{
			PlayerAddress: strconv.FormatUint(i, 10),
			WonCount:      101 - i,
			DateAdded:     "2006-01-02 15:05:05.999999999 +0000 UTC",
		}
	}
	return sorted
}

func TestUpdatePlayerInfoAtNowTooLongNoAdd(t *testing.T) {
	beforeWinners := makeMaxLengthSortedWinningPlayers()
	now, err := types.ParseDateAddedAsTime("2006-01-02 15:05:05.999999999 +0000 UTC")
	require.NoError(t, err)
	leaderboard := types.Leaderboard{
		Winners: beforeWinners,
	}
	err = leaderboard.UpdatePlayerInfoAtNow(now, types.PlayerInfo{
		Index:    "100",
		WonCount: 1,
	})
	require.NoError(t, err)
	require.Equal(t, len(beforeWinners), len(leaderboard.Winners))
	require.EqualValues(t, beforeWinners, leaderboard.Winners)
	require.NoError(t, leaderboard.Validate())
}
