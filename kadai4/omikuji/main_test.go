package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_Normal(t *testing.T) {
	t.Parallel()
	tests := []string{"大吉", "吉", "中吉", "凶"}
	for _, tt := range tests {
		tt = tt
		count := 0
		for {
			// 30回くらいやれば、期待したおみくじが返ってくるだろう。
			if count == 30 {
				t.Fatalf("%s doesn't appear in %d requests", tt, count)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)
			server := Server{config: DefaultConfig()}
			rand.Seed(server.config.Now().UnixNano())
			server.Handler().ServeHTTP(w, r)

			rw := w.Result()
			defer rw.Body.Close()

			if rw.StatusCode != http.StatusOK {
				t.Fatalf("status code is %d", rw.StatusCode)
			}

			b, err := ioutil.ReadAll(rw.Body)
			if err != nil {
				t.Fatalf("failed to read body: %s", err)
			}

			re := fmt.Sprintf("{\"result\":\"%s\"}\n", tt)
			if string(b) == re {
				break
			}
			count++
		}
	}
}

func Test_Sanganiti(t *testing.T) {
	t.Parallel()
	tests := []struct {
		t      time.Time
		result string
	}{
		{
			t:      time.Date(2018, 1, 1, 0, 0, 0, 0, time.Local),
			result: "大吉",
		},
		{
			t:      time.Date(2018, 1, 2, 3, 4, 0, 0, time.Local),
			result: "大吉",
		},
		{
			t:      time.Date(2018, 1, 3, 23, 59, 59, 0, time.Local),
			result: "大吉",
		},
	}
	for _, tt := range tests {
		server := Server{config: &Config{
			Port: "9090",
			Now: func() time.Time {
				return tt.t
			},
		}}

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		rand.Seed(server.config.Now().UnixNano())
		server.Handler().ServeHTTP(w, r)

		rw := w.Result()
		defer rw.Body.Close()

		if rw.StatusCode != http.StatusOK {
			t.Fatalf("status code is %d", rw.StatusCode)
		}

		b, err := ioutil.ReadAll(rw.Body)
		if err != nil {
			t.Fatalf("failed to read body: %s", err)
		}

		re := fmt.Sprintf("{\"result\":\"%s\"}\n", tt.result)
		if string(b) != re {
			t.Fatalf("want: %s, got: %s", re, string(b))
		}
	}
}
