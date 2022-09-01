package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/seanhagen/cipherator/cipher"
	"github.com/spf13/cobra"
)

const webShortHelpText = "web help short"
const webLongHelpText = "web help looooooooong"

func getWebCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: webShortHelpText,
		Long:  webLongHelpText,
	}
}

// encodeHandler ...
func encodeHandler(w http.ResponseWriter, r *http.Request) {
	enc, err := getRequestEncoder(w, r)
	if returnIfErrorToHandle(w, err) {
		return
	}

	err = enc.Encode(r.Body)
	if returnIfErrorToHandle(w, err) {
		return
	}

	w.WriteHeader(http.StatusOK)
}

// returnIfErrorToHandle  ...
func returnIfErrorToHandle(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	w.WriteHeader(http.StatusInternalServerError)

	return true
}

// getRequestEncoder ...
func getRequestEncoder(wr http.ResponseWriter, r *http.Request) (cipher.Encoder, error) {
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
	// r.HandleFunc("/", homeHandler)
	r.HandleFunc("/encode/{cipher}", encodeHandler)

	return r, nil
}
