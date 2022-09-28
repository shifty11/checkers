package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/checkers module sentinel errors
var (
	ErrGameNotFound     = sdkerrors.Register(ModuleName, 1103, "game by id not found")
	ErrCreatorNotPlayer = sdkerrors.Register(ModuleName, 1104, "message creator is not a player")
	ErrNotPlayerTurn    = sdkerrors.Register(ModuleName, 1105, "player tried to play out of turn")
	ErrWrongMove        = sdkerrors.Register(ModuleName, 1106, "wrong move")
)
