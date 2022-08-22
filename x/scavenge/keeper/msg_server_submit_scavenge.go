package keeper

import (
	"context"

	"scavenge/x/scavenge/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto"
)

func (k msgServer) SubmitScavenge(goCtx context.Context, msg *types.MsgSubmitScavenge) (*types.MsgSubmitScavengeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// 1. create new scavenge store object with related type from the Msg (MsgSubmitScavenge)
	var scavenge = types.Scavenge{
		Index:        msg.SolutionHash,
		Description:  msg.Description,
		SolutionHash: msg.SolutionHash,
		Reward:       msg.Reward,
	}

	// 2. We define the index as SolutionHash.
	// By using this index, we'll try to get the scavenge object
	// If it's found, it means that scavenge already exists in store

	_, isFound := k.GetScavenge(ctx, scavenge.SolutionHash)

	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Scavenge with that solution hash already exists!")
	}

	// Getting the address of Scavenge module account
	moduleAcct := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))

	// convert the message creator address from a string into sdk.AccAddress
	scavenger, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	// convert tokens from string to sdk.Coins
	reward, err := sdk.ParseCoinsNormalized(scavenge.Reward)
	if err != nil {
		panic(err)
	}

	// send tokens from the scavenge creator to the module account
	sdkError := k.bankKeeper.SendCoins(ctx, scavenger, moduleAcct, reward)
	if sdkError != nil {
		return nil, sdkError
	}

	// write the scavenge to the store
	k.SetScavenge(ctx, scavenge)

	return &types.MsgSubmitScavengeResponse{}, nil
}
