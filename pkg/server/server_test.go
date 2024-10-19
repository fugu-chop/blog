package server

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	ctx := context.Background()
	_, err := New(ctx, ":9999")

	assert.Nil(t, err)
}

func TestNewServer_StartandStop(t *testing.T) {
	ctx := context.Background()

	s, err := New(ctx, ":9999")

	assert.NotNil(t, s)
	assert.Nil(t, err)

	wg := &sync.WaitGroup{}

	wg.Add(1)

	go func() {
		defer wg.Done()

		err := s.Start(ctx)

		assert.Nil(t, err)
	}()

	err = s.Shutdown(ctx)

	assert.Nil(t, err)

	wg.Wait()
}
