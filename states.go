package main

import (
	"context"
	"fmt"
)

func drainService(ctx context.Context, args Args) (Args, State[Args], error) {
	l, err := args.Services.List(ctx)
	if err != nil {
		return args, nil, err
	}

	found := false
	for _, entry := range l {
		if entry == args.Name {
			found = true
			break
		}
	}
	if !found {
		return args, nil, fmt.Errorf("the service was not found")
	}

	if err := args.Services.Drain(ctx, args.Name); err != nil {
		return args, nil, fmt.Errorf("problem draining the service: %w", err)
	}

	return args, removeService, nil
}

func removeService(ctx context.Context, args Args) (Args, State[Args], error) {
	if err := args.Services.Remove(ctx, args.Name); err != nil {
		return args, nil, fmt.Errorf("could not remove the service: %w", err)
	}

	hasStorage, err := args.Services.HasStorage(ctx, args.Name)
	if err != nil {
		return args, nil, fmt.Errorf("HasStorage() failed: %w", err)
	}
	if hasStorage {
		return args, removeBackups, nil
	}

	return args, nil, nil
}

func removeBackups(ctx context.Context, args Args) (Args, State[Args], error) {
	return args, nil, nil
}
