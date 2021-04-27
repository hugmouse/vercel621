package handler

import (
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	E621Url = "https://e621.net"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	e621Resp, err := http.Get(E621Url + r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("\"/\" E621 request failed: " + err.Error()))
		return
	}
	if e621Resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("\"/\" E621 not ok status code returned: " + strconv.Itoa(e621Resp.StatusCode)))
		return
	}
	e621Info, err := ioutil.ReadAll(e621Resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("\"/\" E621 body reading failed: " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(e621Info)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("An error occurred while trying to write to the ResponseWriter: " + err.Error()))
		return
	}
}
