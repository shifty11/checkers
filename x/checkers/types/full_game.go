package types

import (
	"errors"
	"fmt"
	"github.com/alice/checkers/x/checkers/rules"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidBlack     = sdkerrors.Register(ModuleName, 1101, "black address is invalid: %s")
	ErrInvalidRed       = sdkerrors.Register(ModuleName, 1102, "red address is invalid: %s")
	ErrGameNotParseable = sdkerrors.Register(ModuleName, 1103, "game cannot be parsed")
)

func (m *StoredGame) GetBlackAddress() (black sdk.AccAddress, err error) {
	black, errBlack := sdk.AccAddressFromBech32(m.Black)
	return black, sdkerrors.Wrapf(errBlack, ErrInvalidBlack.Error(), m.Black)
}

func (m *StoredGame) GetRedAddress() (red sdk.AccAddress, err error) {
	red, errRed := sdk.AccAddressFromBech32(m.Red)
	return red, sdkerrors.Wrapf(errRed, ErrInvalidRed.Error(), m.Red)
}

func (m *StoredGame) ParseGame() (game *rules.Game, err error) {
	board, errBoard := rules.Parse(m.Board)
	if errBoard != nil {
		return nil, sdkerrors.Wrapf(errBoard, ErrGameNotParseable.Error())
	}
	board.Turn = rules.StringPieces[m.Turn].Player
	if board.Turn.Color == "" {
		return nil, sdkerrors.Wrapf(errors.New(fmt.Sprintf("Turn: %s", m.Turn)), ErrGameNotParseable.Error())
	}
	return board, nil
}

func (m *StoredGame) Validate() (err error) {
	_, err = m.GetBlackAddress()
	if err != nil {
		return err
	}
	_, err = m.GetRedAddress()
	if err != nil {
		return err
	}
	_, err = m.ParseGame()
	return err
}
