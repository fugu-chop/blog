package views

import (
	"errors"
	"io/fs"
	"testing"

	"github.com/fugu-chop/blog/pkg/templates"
	"github.com/stretchr/testify/assert"
)

func TestMust(t *testing.T) {
	tests := []struct {
		name       string
		inputError error
	}{
		{
			name:       "no panic with no error",
			inputError: nil,
		},
		{
			name:       "panics with error",
			inputError: errors.New("some error"),
		},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.inputError == nil {
				assert.NotPanics(t, func() { Must(Template{}, tc.inputError) })
			} else {
				assert.Panics(t, func() { Must(Template{}, tc.inputError) })
			}
		})
	}
}

func TestParseFS(t *testing.T) {
	tpl, err := ParseFS(templates.FS, "layout.gohtml", "home.gohtml")

	assert.Nil(t, err)
	assert.NotNil(t, tpl.htmlTpl)
}

func TestParseFS_Errors(t *testing.T) {
	tests := []struct {
		name       string
		fileSystem fs.FS
		patterns   []string
	}{
		{
			name:       "no FS or templates provided",
			fileSystem: nil,
			patterns:   []string{},
		},
		{
			name:       "template does not exist",
			fileSystem: templates.FS,
			patterns:   []string{"some-random-template"},
		},
	}

	for _, testCase := range tests {
		tc := testCase

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tpl, err := ParseFS(tc.fileSystem, tc.patterns...)

			assert.Nil(t, tpl.htmlTpl)
			assert.NotNil(t, err)
		})
	}

}
