package memoize

import (
	"context"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test1To1(t *testing.T) {
	calls := int64(0)
	wrapped := func(ctx context.Context, i1 string) (int64, error) {
		atomic.AddInt64(&calls, 1)
		return strconv.ParseInt(i1, 10, 64)
	}

	memoized := Memoize1To1(50*time.Millisecond, wrapped)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r, err := memoized(context.Background(), "5")
			require.NoError(t, err)
			require.EqualValues(t, 5, r)
			require.EqualValues(t, 1, atomic.LoadInt64(&calls))
		}()
	}
	wg.Wait()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r, err := memoized(context.Background(), "7")
			require.NoError(t, err)
			require.EqualValues(t, 7, r)
			require.EqualValues(t, 2, atomic.LoadInt64(&calls))
		}()
	}
	wg.Wait()

	_, err := memoized(context.Background(), "bubba")
	require.Error(t, err)
	require.EqualValues(t, 3, atomic.LoadInt64(&calls))

	time.Sleep(100 * time.Millisecond)
	r, err := memoized(context.Background(), "5")
	require.NoError(t, err)
	require.EqualValues(t, 5, r)
	require.EqualValues(t, 4, atomic.LoadInt64(&calls))
}
