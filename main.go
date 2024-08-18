package main

import (
	"context"
	"fmt"
)

type State[T any] func(ctx context.Context, args T) (T, State[T], error)

func main() {

}

// Service removes a service from a cluster and associated storage.
// The last 3 storage backups are retained for whatever the storage retainment
// period is.
func Service(ctx context.Context, args Args) error {
	if err := args.validate(ctx); err != nil {
		return err
	}

	start := drainService
	_, err := Run[Args](ctx, args, start)
	if err != nil {
		return fmt.Errorf("problem removing service %q: %w", args.Name, err)
	}
	return nil
}

// Run is a state running for the state machine
func Run[T any](ctx context.Context, args T, start State[T]) (T, error) {
	var err error
	current := start
	for {
		if ctx.Err() != nil {
			return args, ctx.Err()
		}
		args, current, err = current(ctx, args)
		if err != nil {
			return args, err
		}
		if current == nil {
			return args, nil
		}
	}
}

// storageClient provides the methods on a storage service
// that must be provided to use Remove().
type storageClient interface {
	RemoveBackups(ctx context.Context, service string, mustKeep int) error
	RemoveContainer(ctx context.Context, service string) error
}

// serviceClient provides methods to do operations for services
// within a cluster.
type serviceClient interface {
	Drain(ctx context.Context, service string) error
	Remove(ctx context.Context, service string) error
	List(ctx context.Context) ([]string, error)
	HasStorage(ctx context.Context, service string) (bool, error)
}

// Args are arguments to Service().
type Args struct {
	// Name is the name of the service.
	Name string

	// Storage is a client that can remove storage backups and storage
	// containers for a service.
	Storage storageClient
	// Services is a client that allows the draining and removal of
	// a service from the cluster.
	Services serviceClient
}

func (a Args) validate(ctx context.Context) error {
	if a.Name == "" {
		return fmt.Errorf("Name cannot be an empty string")
	}

	if a.Storage == nil {
		return fmt.Errorf("Storage cannot be nil")
	}
	if a.Services == nil {
		return fmt.Errorf("Services cannot be nil")
	}
	return nil
}
