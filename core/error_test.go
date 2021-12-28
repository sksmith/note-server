package core_test

import (
	"errors"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/sksmith/note-server/core"
)

func TestMain(m *testing.M) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	os.Exit(m.Run())
}

func TestIsErrNotFound(t *testing.T) {
	tests := []struct {
		input error
		want  bool
	}{
		{input: errors.New("some madeup error"), want: false},
		{input: &core.ErrNotFound{}, want: true},
	}

	for _, test := range tests {
		got := core.IsErrNotFound(test.input)
		if test.want != got {
			t.Errorf("want=[%v] got=[%v]", test.want, got)
		}
	}
}
