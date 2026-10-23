//go:build unit

package service

import (
	"context"
	"math"
	"testing"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestAdminGroupHiddenRateMultiplierCreateAndValidation(t *testing.T) {
	for _, value := range []*float64{nil, groupDuplicateTestPointer(0.0001), groupDuplicateTestPointer(2.75), groupDuplicateTestPointer(999999.9999)} {
		repo := &groupRepoStubForAdmin{}
		svc := &adminServiceImpl{groupRepo: repo}
		created, err := svc.CreateGroup(context.Background(), &CreateGroupInput{
			Name: "hidden-rate", Platform: PlatformOpenAI, RateMultiplier: 1.25, HiddenRateMultiplier: value,
		})
		require.NoError(t, err)
		require.Equal(t, value, created.HiddenRateMultiplier)
		require.Equal(t, 1.25, created.RateMultiplier)
		if value != nil {
			require.NotSame(t, value, created.HiddenRateMultiplier)
		}
	}

	for _, value := range []float64{0, -1, 0.00001, 0.00009999, math.SmallestNonzeroFloat64, math.NaN(), math.Inf(1), math.Inf(-1), 1000000} {
		t.Run("reject invalid value", func(t *testing.T) {
			repo := &groupRepoStubForAdmin{getByID: &Group{ID: 1, Platform: PlatformOpenAI, RateMultiplier: 1}}
			svc := &adminServiceImpl{groupRepo: repo}
			_, err := svc.CreateGroup(context.Background(), &CreateGroupInput{RateMultiplier: 1, HiddenRateMultiplier: &value})
			require.Equal(t, "INVALID_HIDDEN_RATE_MULTIPLIER", infraerrors.Reason(err))
			require.Nil(t, repo.created)
			_, err = svc.UpdateGroup(context.Background(), 1, &UpdateGroupInput{HiddenRateMultiplier: &value})
			require.Equal(t, "INVALID_HIDDEN_RATE_MULTIPLIER", infraerrors.Reason(err))
			require.Nil(t, repo.updated)
		})
	}
}

func TestAdminGroupHiddenRateMultiplierUpdateTriState(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input UpdateGroupInput
		want  *float64
	}{
		{name: "omitted retains", want: groupDuplicateTestPointer(2.5)},
		{name: "null clears", input: UpdateGroupInput{HiddenRateMultiplierSet: true}},
		{name: "minimum value sets", input: UpdateGroupInput{HiddenRateMultiplier: groupDuplicateTestPointer(0.0001)}, want: groupDuplicateTestPointer(0.0001)},
		{name: "value sets", input: UpdateGroupInput{HiddenRateMultiplierSet: true, HiddenRateMultiplier: groupDuplicateTestPointer(3.5)}, want: groupDuplicateTestPointer(3.5)},
		{name: "direct value sets without flag", input: UpdateGroupInput{HiddenRateMultiplier: groupDuplicateTestPointer(4.5)}, want: groupDuplicateTestPointer(4.5)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			repo := &groupRepoStubForAdmin{getByID: &Group{ID: 1, Platform: PlatformOpenAI, RateMultiplier: 1.25, HiddenRateMultiplier: groupDuplicateTestPointer(2.5)}}
			svc := &adminServiceImpl{groupRepo: repo}
			got, err := svc.UpdateGroup(context.Background(), 1, &tt.input)
			require.NoError(t, err)
			require.Equal(t, tt.want, got.HiddenRateMultiplier)
			require.Equal(t, 1.25, got.RateMultiplier)
			require.Same(t, got, repo.updated)
		})
	}
}

func TestDuplicateGroupHiddenRateMultiplierIsIndependent(t *testing.T) {
	source := &Group{RateMultiplier: 1.25, HiddenRateMultiplier: groupDuplicateTestPointer(3.5)}
	cloned := cloneGroupForDuplicate(source, "")
	require.Equal(t, source.HiddenRateMultiplier, cloned.HiddenRateMultiplier)
	require.NotSame(t, source.HiddenRateMultiplier, cloned.HiddenRateMultiplier)
	*cloned.HiddenRateMultiplier = 9
	require.Equal(t, 3.5, *source.HiddenRateMultiplier)
	require.Nil(t, cloneGroupForDuplicate(&Group{}, "").HiddenRateMultiplier)
}
