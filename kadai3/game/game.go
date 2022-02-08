package game

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Game struct {
	OutStream, ErrStream io.Writer
}

func (g *Game) Run() {
	// ゲームスタート
	sch := start(os.Stdin)
	<-sch

	// 入力値用のチャネルを用意する。
	ch := input(os.Stdin)
}

func input(r io.Reader) <-chan string {
	ch := make(chan string)
	go func() {
		s := bufio.NewScanner(r)
		for s.Scan() {
			ch <- s.Text()
		}
		close(ch)
	}()

	return ch
}

func start(r io.Reader) <-chan struct{} {
	fmt.Println("■ タイピングゲームを始めます。")
	fmt.Println("■ 制限時間は30秒です。")
	fmt.Println(">>> press any key to start <<<")

	ch := make(chan struct{})
	go func() {
		s := bufio.NewScanner(r)
		for s.Scan() {
			ch <- struct{}{}
			break
		}
		close(ch)
	}()

	return ch
}
