package v1tov2

import "github.com/alice/checkers/x/checkers/types"

const (
	UpgradeName         = "v1tov2"
	StoredGameChunkSize = 1_000
	PlayerInfoChunkSize = types.LeaderboardWinnerLength * 2
)
