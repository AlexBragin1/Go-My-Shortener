package main

import (
	"errors"
	"fmt"

	"io"
	"math/rand"
	"net/http"
	"net/url"
)

const AlphaBet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const host = "http://localhost:8080/"

type MyURL struct {
	LongURL  string
	ShortURL string
	Code     int
}

var Urls = map[string]MyURL{}

func Shorten() string {
	b := make([]byte, 7)
	for i := range b {
		b[i] = AlphaBet[rand.Intn(len(AlphaBet))]
	}
	return string(b)
}

func IsValueURL(LongUrl string) bool {
	_, err := url.ParseRequestURI(LongUrl)
	if err != nil {
		fmt.Println("1", "false", err)
		return false
	}
	u, err := url.Parse(LongUrl)
	if err != nil || u.Host == "" {
		fmt.Println("2", "false", err)
		return false
	}
	return true
}

func SearchID(id string) (string, error) {
	if Urls[id].ShortURL != "" {
		return Urls[id].LongURL, nil
	}
	return "", errors.New("NotFoundShortUrl")
}
func SearchLongURL(LUrl string) (MyURL, bool) {
	var t MyURL
	for _, a := range Urls {
		if a.LongURL == LUrl {
			return a, true
		}
	}
	t.LongURL = LUrl
	return t, false
}

func AddURL(id string, a MyURL) {
	Urls[id] = a

}

func GetOrPostHandler(w http.ResponseWriter, r *http.Request) {
	var id string
	var err error
	var LUrl string
	var NewURL MyURL
	flag := false
	switch r.Method {
	case http.MethodGet:
		{
			ShortID := string(r.URL.Path)
			if ShortID[1:] != "" {
				LUrl, err = SearchID(ShortID[1:])
				if err != nil {
					http.Error(w, "Missing1", http.StatusBadRequest)
					return
				}
				w.Header().Set("Location", LUrl)
				w.WriteHeader(307)
				return
			}
		}
	case http.MethodPost:
		{

			b, err2 := io.ReadAll(r.Body)
			// обрабатываем ошибку

			if err2 != nil {
				http.Error(w, err2.Error(), 500)
				NewURL.Code = 400
			}
			NewURL.LongURL = string(b)
			if !IsValueURL(NewURL.LongURL) {
				http.Error(w, "Missing2", 400)
				NewURL.Code = 400
			} else {
				NewURL, flag = SearchLongURL(NewURL.LongURL)
				if !flag && NewURL.Code == 0 {
					id = Shorten()
					NewURL.Code = http.StatusOK
					NewURL.ShortURL = host + id
				}
			}
			w.WriteHeader(NewURL.Code)
			w.Header().Set("Content-Type", "application/text")
			fmt.Fprintf(w, "%s", NewURL.ShortURL)

			if !flag {
				AddURL(id, NewURL)
			}

			return
		}
	}

}

func main() {

	http.HandleFunc("/", GetOrPostHandler)

	http.ListenAndServe(":8080", nil)
}
