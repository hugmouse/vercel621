package handler

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
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
	var (
		getStatic bool
		err       error
	)
	e621Resp := &http.Response{}

	if len(r.URL.Path) > 4 && r.URL.Path[:5] == "/data" {
		getStatic = true
	}

	if getStatic {
		e621Resp, err = http.Get(E621StaticURL + r.URL.Path + url.QueryEscape(r.URL.RawPath))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(CombinedError(ErrRequestFailed.Error(), err.Error()))
			return
		}
	} else {
		e621Resp, err = http.Get(E621Url + r.URL.Path + url.QueryEscape(r.URL.RawPath))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(CombinedError(ErrRequestFailed.Error(), err.Error()))
			return
		}
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

	e621InfoMirrored := bytes.ReplaceAll(e621Info, []byte(E621StaticURL), []byte("https://"+r.Host))
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(e621InfoMirrored)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(CombinedError(ErrResponseWriterFailed.Error(), err.Error()))
		return
	}
}
