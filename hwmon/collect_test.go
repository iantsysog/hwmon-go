package hwmon

import (
	"context"
	"errors"
	"testing"

	"github.com/iantsysog/hwmon-go/internal/test"
)

type fakeBackend struct {
	name string
	rs   []Reading
	err  error
}

func (f fakeBackend) Name() string { return f.name }
func (f fakeBackend) Collect(ctx context.Context, emit func(Reading)) error {
	for _, r := range f.rs {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		emit(r)
	}
	return f.err
}

type staticRegistry []Backend

func (s staticRegistry) List() []Backend { return []Backend(s) }

func TestCollect_NoBackends(t *testing.T) {
	_, err := Collect(context.Background(), WithRegistry(staticRegistry(nil)))
	test.ErrorIs(t, err, ErrNoBackends)
}

func TestCollect_FilterSortAndJoinErrors(t *testing.T) {
	r := staticRegistry([]Backend{
		fakeBackend{
			name: "a",
			rs: []Reading{
				{Kind: KindTemp, Name: "T2", Source: "a", KeyOrID: "2", Value: 2.0},
				{Kind: KindTemp, Name: "T1", Source: "a", KeyOrID: "1", Value: 1.0},
			},
		},
		fakeBackend{
			name: "b",
			rs:   []Reading{{Kind: KindOther, Name: "skip", Source: "b", KeyOrID: "x"}},
			err:  errors.New("boom"),
		},
	})

	out, err := Collect(
		context.Background(),
		WithRegistry(r),
		WithFilter(func(r Reading) bool { return r.Kind != KindOther }),
	)

	test.Len(t, out, 2)
	test.Eq(t, "T1", out[0].Name)
	test.Eq(t, "T2", out[1].Name)
	test.Error(t, err)
}

func TestCollect_WithBackendsFiltersByName(t *testing.T) {
	r := staticRegistry([]Backend{
		fakeBackend{name: "smc", rs: []Reading{{Name: "x"}}},
		fakeBackend{name: "hid", rs: []Reading{{Name: "y"}}},
	})
	out, err := Collect(context.Background(), WithRegistry(r), WithBackends("hid"))
	test.NoError(t, err)
	test.Len(t, out, 1)
	test.Eq(t, "y", out[0].Name)
}
