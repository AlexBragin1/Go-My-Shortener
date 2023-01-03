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

type MyUrl struct {
	LongUrl  string
	ShortUrl string
	Code     int
}

var Urls = map[string]MyUrl{}

func Shorten() string {
	b := make([]byte, 7)
	for i := range b {
		b[i] = AlphaBet[rand.Intn(len(AlphaBet))]
	}
	return string(b)
}

func IsValueUrl(LongUrl string) bool {
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

func SearchId(id string) (string, error) {
	if Urls[id].ShortUrl != "" {
		return Urls[id].LongUrl, nil
	}
	return "", errors.New("NotFoundShortUrl")
}
func SearchLongUrl(LUrl string) (MyUrl, bool) {
	var t MyUrl
	for _, a := range Urls {
		if a.LongUrl == LUrl {
			return a, true
		}
	}
	t.LongUrl = LUrl
	return t, false
}

func AddUrls(id string, a MyUrl) {
	Urls[id] = a
	return
}

func GetOrPostHandler(w http.ResponseWriter, r *http.Request) {
	var id string
	var err error
	var LUrl string
	var NewUrl MyUrl
	flag := false
	switch r.Method {
	case http.MethodGet:
		{
			ShortID := string(r.URL.Path)
			if ShortID[1:len(ShortID)] != "" {
				LUrl, err = SearchId(ShortID[1:len(ShortID)])
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
				NewUrl.Code = 400
			}
			NewUrl.LongUrl = string(b)
			if !IsValueUrl(NewUrl.LongUrl) {
				http.Error(w, "Missing2", 400)
				NewUrl.Code = 400
			} else {
				NewUrl, flag = SearchLongUrl(NewUrl.LongUrl)
				if !flag && NewUrl.Code == 0 {
					id = Shorten()
					NewUrl.Code = http.StatusOK
					NewUrl.ShortUrl = host + id
				}
			}
			w.WriteHeader(NewUrl.Code)
			w.Header().Set("Content-Type", "application/text")
			fmt.Fprintf(w, "%s", NewUrl.ShortUrl)

			if flag == false {
				AddUrls(id, NewUrl)
			}

			return
		}
	}
	return
}

func main() {

	http.HandleFunc("/", GetOrPostHandler)

	http.ListenAndServe(":8080", nil)
}
