package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nanmu42/limitio"
	"github.com/seanhagen/cipherator/cipher"
	"github.com/spf13/cobra"
)

const webShortHelpText = "web help short"
const webLongHelpText = "web help looooooooong"

const (
	// read at most 1MB when generating an error from an HTTP request.
	maxPostBodySize           = 1 << (10 * 2)
	shouldRegardOversizeAsEOF = false
)

func getWebCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: webShortHelpText,
		Long:  webLongHelpText,
	}
}

// cipherRouteHandler ...
func cipherRouteHandler(w http.ResponseWriter, r *http.Request) {
	enc, err := getRequestEncoder(w, r)
	if returnIfErrorToHandle(w, err) {
		return
	}

	err = processRequest(r, enc)
	if returnIfErrorToHandle(w, err) {
		return
	}

	w.WriteHeader(http.StatusOK)
}

// processRequest ...
func processRequest(r *http.Request, enc cipher.Handler) error {
	// reader := limitio.NewReadCloser(r.Body, maxPostBodySize, shouldRegardOversizeAsEOF)
	reader := limitio.NewReadCloser(r.Body, maxPostBodySize, false)
	defer reader.Close()

	vars := mux.Vars(r)
	op, ok := vars["operation"]
	if !ok {
		return fmt.Errorf("'operation' not a valid key in request vars")
	}

	var err error
	switch op {
	case "encode":
		err = enc.Encode(reader)
	case "decode":
		err = enc.Decode(reader)
	default:
		err = fmt.Errorf("operation must be one of 'encode' or 'decode'")
	}

	return err
}

// getRequestEncoder ...
func getRequestEncoder(wr http.ResponseWriter, r *http.Request) (cipher.Handler, error) {
	vars := mux.Vars(r)
	c, ok := vars["cipher"]
	if !ok {
		return nil, fmt.Errorf("'cipher' not a valid key in request vars")
	}

	et, err := cipher.ParseEncoderType(c)
	if err != nil {
		return nil, fmt.Errorf("unable to parse encoder type: %w", err)
	}

	return cipher.New(et, wr)
}

// getWebRouter ...
func getWebRouter() (http.Handler, error) {
	r := mux.NewRouter()
	r.HandleFunc("/{operation}/{cipher}", cipherRouteHandler)
	return r, nil
}

// returnIfErrorToHandle  ...
func returnIfErrorToHandle(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	io.WriteString(w, "\n")

	if errors.Is(err, limitio.ErrThresholdExceeded) {
		http.Error(w, " error during operation", http.StatusRequestEntityTooLarge)
	} else {
		http.Error(w, " error during operation", http.StatusInternalServerError)
	}

	return true
}
