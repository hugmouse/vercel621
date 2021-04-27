package handler

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	E621Url       = "https://e621.net"
	E621StaticURL = "https://static1.e621.net"
)

var (
	ErrRequestFailed        = errors.New("e621 request failed")
	ErrNotOKCode            = errors.New("e621 not ok status code returned")
	ErrBodyReadingError     = errors.New("e621 body reading failed")
	ErrResponseWriterFailed = errors.New("an error occurred while trying to write to the ResponseWriter")
)

// CombinedError combines multiple errors into one
func CombinedError(errors ...string) []byte {
	return []byte(strings.Join(errors, ". "))
}

// Handler is what does all the work
func Handler(w http.ResponseWriter, r *http.Request) {
	e621Resp, err := http.Get(E621Url + r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(CombinedError(ErrRequestFailed.Error(), err.Error()))
		return
	}
	if e621Resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(CombinedError(ErrNotOKCode.Error(), strconv.Itoa(e621Resp.StatusCode)))
		return
	}
	e621Info, err := ioutil.ReadAll(e621Resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(CombinedError(ErrBodyReadingError.Error(), err.Error()))
		return
	}

	e621InfoMirrored := bytes.ReplaceAll(e621Info, []byte(E621StaticURL), []byte(r.Host))
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(e621InfoMirrored)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(CombinedError(ErrResponseWriterFailed.Error(), err.Error()))
		return
	}
}
