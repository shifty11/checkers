package keeper_test

import (
	"context"
	"errors"
	keepertest "github.com/alice/checkers/testutil/keeper"
	"github.com/alice/checkers/testutil/mock_types"
	"github.com/alice/checkers/x/checkers"
	"github.com/alice/checkers/x/checkers/keeper"
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func setupKeeperForWagerHandler(t testing.TB) (keeper.Keeper, context.Context,
	*gomock.Controller, *mock_types.MockBankEscrowKeeper) {
	ctrl := gomock.NewController(t)
	bankMock := mock_types.NewMockBankEscrowKeeper(ctrl)
	k, ctx := keepertest.CheckersKeeperWithMocks(t, bankMock)
	checkers.InitGenesis(ctx, *k, *types.DefaultGenesis())
	context := sdk.WrapSDKContext(ctx)
	return *k, context, ctrl, bankMock
}

func TestWagerHandlerCollectWrongNoBlack(t *testing.T) {
	keeper, context, ctrl, _ := setupKeeperForWagerHandler(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	defer func() {
		r := recover()
		require.NotNil(t, r, "The code did not panic")
		require.Equal(t, "black address is invalid: : empty address string is not allowed", r)
	}()
	keeper.CollectWager(ctx, &types.StoredGame{
		MoveCount: 0,
	})
}

func TestWagerHandlerCollectFailedNoMove(t *testing.T) {
	keeper, context, ctrl, escrow := setupKeeperForWagerHandler(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	black, _ := sdk.AccAddressFromBech32(alice)
	escrow.EXPECT().
		SendCoinsFromAccountToModule(ctx, black, types.ModuleName, gomock.Any()).
		Return(errors.New("Oops"))
	err := keeper.CollectWager(ctx, &types.StoredGame{
		Black:     alice,
		MoveCount: 0,
		Wager:     45,
	})
	require.NotNil(t, err)
	require.EqualError(t, err, "black cannot pay the wager: Oops")
}

func TestWagerHandlerCollectNoMove(t *testing.T) {
	keeper, context, ctrl, escrow := setupKeeperForWagerHandler(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	escrow.ExpectPay(context, alice, 45)
	err := keeper.CollectWager(ctx, &types.StoredGame{
		Black:     alice,
		MoveCount: 0,
		Wager:     45,
	})
	require.Nil(t, err)
}

func TestWagerHandlerPayWrongEscrowFailed(t *testing.T) {
	keeper, context, ctrl, escrow := setupKeeperForWagerHandler(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	black, _ := sdk.AccAddressFromBech32(alice)
	escrow.EXPECT().
		SendCoinsFromModuleToAccount(ctx, types.ModuleName, black, gomock.Any()).
		Times(1).
		Return(errors.New("Oops"))
	defer func() {
		r := recover()
		require.NotNil(t, r, "The code did not panic")
		require.Equal(t, r, "cannot pay winnings to winner: Oops")
	}()
	keeper.MustPayWinnings(ctx, &types.StoredGame{
		Black:     alice,
		Red:       bob,
		Winner:    "b",
		MoveCount: 1,
		Wager:     45,
	})
}

func TestWagerHandlerPayEscrowCalledTwoMoves(t *testing.T) {
	keeper, context, ctrl, escrow := setupKeeperForWagerHandler(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	escrow.ExpectRefund(context, alice, 90)
	keeper.MustPayWinnings(ctx, &types.StoredGame{
		Black:     alice,
		Red:       bob,
		Winner:    "b",
		MoveCount: 2,
		Wager:     45,
	})
}
