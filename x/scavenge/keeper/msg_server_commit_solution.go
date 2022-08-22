package keeper

import (
	"context"

	"scavenge/x/scavenge/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Tasks:
// 1. Check that commit with a given hash doesn't exist in the store
// 2. Write a new commit to the store

func (k msgServer) CommitSolution(goCtx context.Context, msg *types.MsgCommitSolution) (*types.MsgCommitSolutionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var commit = types.Commit{
		Index:                 msg.SolutionScavengerHash,
		SolutionHash:          msg.SolutionHash,
		SolutionScavengerHash: msg.SolutionScavengerHash,
	}

	// 2. We define the index as SolutionScavengerHash.
	// By using this index, we'll try to get the commit object
	// If it's found, it means that commit already exists in store
	_, isFound := k.GetScavenge(ctx, commit.SolutionScavengerHash)

	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Commit with that hash already exists!")
	}

	// write commit to the store
	k.SetCommit(ctx, commit)
	return &types.MsgCommitSolutionResponse{}, nil
}
