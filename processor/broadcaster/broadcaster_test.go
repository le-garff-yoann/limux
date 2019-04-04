package broadcaster

import (
	"limux/processor"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBroadcaster(t *testing.T) {
	t.Parallel()

	br := New()
	require.NotNil(t, br)

	var (
		wg sync.WaitGroup

		n   = 5
		gor = 10
	)

	wg.Add(n * gor)

	for i := 1; i <= gor; i++ {
		recv, _ := br.Recv()

		go func(recv <-chan processor.Event) {
			for {
				select {
				case v, ok := <-recv:
					if !ok {
						return
					}

					t.Log(v.Message)

					wg.Done()
				}

			}
		}(recv)
	}

	go br.Run()

	for i := 1; i <= n; i++ {
		br.Pub <- processor.Event{Message: fmt.Sprintf("%d", i)}
	}

	wg.Wait()

	br.RemoveAllRecv()
}
