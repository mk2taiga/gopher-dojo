package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	server := Server{config: DefaultConfig()}
	rand.Seed(server.config.Now().UnixNano())
	server.Run()
}

type Server struct {
	config *Config
}

type Config struct {
	Port string
	// 外から時間特定の関数を渡せるようにして、テストしやすくしている。
	Now timerFunc
}

func DefaultConfig() *Config {
	return &Config{
		Port: "9090",
		Now:  time.Now,
	}
}

type timerFunc func() time.Time

func (s *Server) Handler() http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		result := drawOmikuji(s.config.Now)
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Println("Error: ", err)
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", f)

	return mux
}

func (s *Server) Run() {
	httpServer := &http.Server{
		Addr:    ":" + s.config.Port,
		Handler: s.Handler(),
	}
	httpServer.ListenAndServe()
}

type omikuji struct {
	Result string `json:"result"`
}

func drawOmikuji(t timerFunc) omikuji {
	if isSanganiti(t) {
		return omikuji{Result: "大吉"}
	}

	var result string
	switch rand.Intn(7) {
	case 6:
		result = "大吉"
	case 5, 4:
		result = "吉"
	case 3, 2, 1:
		result = "中吉"
	default:
		result = "凶"
	}

	return omikuji{Result: result}
}

// 正月の1-3日かどうかを判断する
func isSanganiti(t timerFunc) bool {
	now := t()
	return now.Month() == 1 && now.Day() <= 3
}
