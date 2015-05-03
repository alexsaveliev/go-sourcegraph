package mock

import (
	"testing"

	"golang.org/x/net/context"

	"sourcegraph.com/sourcegraph/go-sourcegraph/sourcegraph"
)

func (s *BuildsServer) MockGet_Return(t *testing.T, want *sourcegraph.Build) (called *bool) {
	called = new(bool)
	s.Get_ = func(ctx context.Context, op *sourcegraph.BuildSpec) (*sourcegraph.Build, error) {
		*called = true
		return want, nil
	}
	return
}

func (s *BuildsServer) MockGetRepoBuildInfo(t *testing.T, info *sourcegraph.RepoBuildInfo) (called *bool) {
	called = new(bool)
	s.GetRepoBuildInfo_ = func(ctx context.Context, op *sourcegraph.BuildsGetRepoBuildInfoOp) (*sourcegraph.RepoBuildInfo, error) {
		*called = true
		return info, nil
	}
	return
}

func (s *BuildsServer) MockList(t *testing.T, want ...*sourcegraph.Build) (called *bool) {
	called = new(bool)
	s.List_ = func(ctx context.Context, op *sourcegraph.BuildListOptions) (*sourcegraph.BuildList, error) {
		*called = true
		return &sourcegraph.BuildList{Builds: want}, nil
	}
	return
}
