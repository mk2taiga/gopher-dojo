package game

import (
	"bufio"
	"fmt"
	"gopher-dojo/kadai3/word"
	"io"
	"math"
	"os"
	"time"
)

type Game struct {
	OutStream, ErrStream io.Writer
}

const ErrorOccurred = 1

func (g *Game) Run() int {
	// ゲームスタート
	sch := start(os.Stdin)
	<-sch

	// 入力値用のチャネルを用意する。
	ch := input(os.Stdin)

	words, err := word.GetWords()
	if err != nil {
		fmt.Println(err)
		return ErrorOccurred
	}

	timer := time.NewTimer(time.Second * 30).C
	correctAnswerCnt := 0
	answerCnt := 0
Outer:
	for _, v := range words {
		fmt.Printf("type: %s\n", v)
	Inner:
		for {
			fmt.Print(">")
			select {
			case <-timer: // タイマーが切れる。
				break Outer
			case s := <-ch: // チャネルでキー入力を受け取る
				answerCnt++
				if s == v {
					fmt.Println("correct")
					// 正解したので、次のループへ。
					correctAnswerCnt++
					break Inner
				} else {
					fmt.Println("wrong")
					fmt.Printf("type: %s\n", v)
				}
			}
		}
	}

	fmt.Println("\nfinish")
	fmt.Printf("正解数: %d words\n", correctAnswerCnt)
	correctness := float64(0)
	if answerCnt != 0 {
		correctness = roundPlus(float64(correctAnswerCnt)/float64(answerCnt)*100, 2)
	}
	fmt.Printf("正確さ: %v %%\n", correctness)
	return 0
}

// 四捨五入
func round(f float64) float64 {
	return math.Floor(f + .5)
}

// 指定した桁数の小数点以下を残す
func roundPlus(f float64, places int) float64 {
	// 10 の冪乗をして、確実に指定した桁数分の小数点以下は残るようにして、四捨五入する。
	shift := math.Pow(10, float64(places))
	return round(f*shift) / shift
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
