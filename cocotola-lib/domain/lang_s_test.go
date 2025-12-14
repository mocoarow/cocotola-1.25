//go:build small

package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"
)

func TestNewLang2_shouldReturnLang_whenArgumentValid(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		args string
		want *domain.Lang2
	}{
		{name: "en", args: "en", want: domain.Lang2EN},
		{name: "ja", args: "ja", want: domain.Lang2JA},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := domain.NewLang2(tt.args)
			require.NoError(t, err)
			assert.Equal(t, tt.want.String(), got.String())
		})
	}
}

func TestNewLang2_shouldReturnError_whenArgumentInvalid(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		args    string
		wantErr error
	}{
		{name: "empty string", args: "", wantErr: domain.ErrInvalidArgument},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := domain.NewLang2(tt.args)
			require.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestLang5_ToLang2_shouldReturnLang2_whenSupportedLang(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		lang5 *domain.Lang5
		lang2 *domain.Lang2
	}{
		{
			name:  domain.Lang5ENUS.String(),
			lang5: domain.Lang5ENUS,
			lang2: domain.Lang2EN,
		},
		{
			name:  domain.Lang5JAJP.String(),
			lang5: domain.Lang5JAJP,
			lang2: domain.Lang2JA,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.lang5.ToLang2(), tt.lang2)
		})
	}
}
