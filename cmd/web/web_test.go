package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWeb_Setup(t *testing.T) {
	args := []string{"-h"}
	output := bytes.NewBuffer(nil)
	expect := webLongHelpText

	cmd := getWebCommand()
	cmd.SetOutput(output)
	cmd.SetArgs(args)
	err := cmd.Execute()

	assert.NoError(t, err)
	assert.Contains(t, output.String(), expect)
}

func TestWeb_BuildRoutes(t *testing.T) {
	var expect *http.Handler
	router, err := getWebRouter()
	require.NoError(t, err)
	require.NotNil(t, router)

	assert.Implements(t, expect, router)
}

func TestWeb_TranslateRouteExistsAndWorks(t *testing.T) {
	router, err := getWebRouter()
	require.NoError(t, err)
	require.NotNil(t, router)

	srv := httptest.NewServer(router)
	require.NotNil(t, srv)

	input := "hello world"
	expect := "ellohay orldway"

	body := strings.NewReader(input)
	req, err := http.NewRequest(http.MethodPost, "/encode/piglatin", body)
	require.NoError(t, err)
	require.NotNil(t, req)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Equal(t, expect, w.Body.String())
}
