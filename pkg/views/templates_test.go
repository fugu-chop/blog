package views

import (
	"errors"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fugu-chop/blog/pkg/config"
	"github.com/fugu-chop/blog/pkg/templates"
	"github.com/fugu-chop/blog/test/pkg/templatetest"
	"github.com/stretchr/testify/assert"
)

func TestGenerateTemplate(t *testing.T) {
	assert.NotPanics(t, func() { GenerateTemplate("root/home.gohtml") })

	template := GenerateTemplate("root/home.gohtml")

	assert.NotNil(t, template)
	assert.IsType(t, Template{}, template)
}

func TestGenerateTemplate_Error(t *testing.T) {
	assert.Panics(t, func() { GenerateTemplate("non-existent-template.gohtml") })
}

func TestExecute(t *testing.T) {
	template := must(parseFS(templates.FS, config.LayoutTemplate, "root/home.gohtml"))

	// Request to home page
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	template.Execute(w, r, nil)
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "text/html; charset=utf-8", resp.Header.Get("Content-Type"))
	assert.Contains(t, string(body), "<h1>Dean Wan</h1>")
}

func TestExecute_CloneError(t *testing.T) {
	mock := templatetest.NewMockTemplateCloner(t)
	mock.EXPECT().Clone().Return(nil, errors.New("something went wrong!"))

	template := Template{
		htmlTpl: mock,
	}

	// Request to home page
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	template.Execute(w, r, nil)
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.NotEqual(t, "text/html; charset=utf-8", resp.Header.Get("Content-Type"))
	assert.Contains(t, string(body), "There was an error rendering the page")
}

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
				assert.NotPanics(t, func() { must(Template{}, tc.inputError) })
			} else {
				assert.Panics(t, func() { must(Template{}, tc.inputError) })
			}
		})
	}
}

func TestParseFS(t *testing.T) {
	tpl, err := parseFS(templates.FS, "layout.gohtml", "root/home.gohtml")

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
			tpl, err := parseFS(tc.fileSystem, tc.patterns...)

			assert.Nil(t, tpl.htmlTpl)
			assert.NotNil(t, err)
		})
	}

}
