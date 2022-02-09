package rangedownloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

const tempDir = "dlTmp"

var wg sync.WaitGroup

type Downloader struct {
	Argv  []string
	procs int
	url   string
	name  string
}

type cliOptions struct {
	Name  string `short:"n" long:"name" description:"output file name with extension. if not provided, rangedownloader will guess a file name based on URL"`
	Procs int    `short:"p" long:"procs" description:"number of parallel" default:"1"`
	Args  struct {
		URL string
	} `positional-args:"yes"`
}

func New() *Downloader {
	return &Downloader{Argv: os.Args[1:]}
}

func (d *Downloader) Run() int {
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		fmt.Println(err)
		return 1
	}

	if err := d.parseCommandLine(); err != nil {
		fmt.Println(err)
		return 1
	}

	l, err := d.getContentLength()
	if err != nil {
		fmt.Println(err)
		return 1
	}
	subFileLen := l / d.procs
	remaining := l % d.procs

	eg, ctx := errgroup.WithContext(context.Background())
	for i := 0; i < d.procs; i++ {
		// goroutine に i を渡すので、再度インスタンスを生成している。
		i := i

		from := subFileLen * i
		to := subFileLen * (i + 1)

		// 最後のループの時だけ、余る分も追加で読み込んでやる。
		if i == d.procs-1 {
			to += remaining
		}

		// goroutine を用いて範囲リクエストを行う。
		eg.Go(func() error {
			return d.rangeRequest(ctx, from, to, i)
		})
	}

	// error group が終わるまで待つ。
	if err := eg.Wait(); err != nil {
		fmt.Println(err)
		return 1
	}

	// 全ての goroutine が終わったら、ファイルを結合する。
	if err := d.createFile(); err != nil {
		fmt.Println(err)
		return 1
	}

	// 一時ディレクトリを削除する。
	if err := os.RemoveAll(tempDir); err != nil {
		fmt.Println(err)
		return 1
	}

	return 0
}

func (d *Downloader) createFile() error {
	// コンソールで指定されて名前でファイルを作成する。
	file, err := os.Create(d.name)
	if err != nil {
		return errors.Wrap(err, "failed to create output file")
	}
	defer file.Close()

	// 一時ディレクトリ内のファイルを全て読み込んで、一つのファイルに結合する。
	for i := 0; i < d.procs; i++ {
		subFile, err := os.Open(path.Join(tempDir, fmt.Sprint(i)))
		if err != nil {
			return errors.Wrap(err, "failed to generate output file")
		}
		io.Copy(file, subFile)
		subFile.Close()
	}

	return nil
}

func (d *Downloader) rangeRequest(ctx context.Context, from int, to int, i int) error {
	client := http.Client{}

	// リクエストを生成
	req, err := http.NewRequest("GET", d.url, nil)
	if err != nil {
		return errors.Wrapf(err, "failed to access the site you provided: %s", d.url)
	}

	// リクエストヘッダに Range を設定
	rangeHeader := fmt.Sprintf("bytes=%d-%d", from, to-1)
	req.Header.Add("Range", rangeHeader)
	// errgroup.WithContext wraps context by calling context.WithCancel
	// cf. https://github.com/golang/sync/blob/master/errgroup/errgroup.go#L34
	req = req.WithContext(ctx)
	// リクエスト実行
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to get response. please try again later ")
	}
	fmt.Printf("Range: %v, %v bytes\n", rangeHeader, resp.ContentLength)
	defer resp.Body.Close()

	// ファイルが存在しなかった場合に生成
	file, err := os.OpenFile(path.Join(tempDir, fmt.Sprint(i)), os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "failed to open temp file: %d/%d", i, d.procs-1)
	}
	defer file.Close()

	// レスポンスの内容をファイルにコピーする
	io.Copy(file, resp.Body)

	return nil
}

func (d *Downloader) getContentLength() (int, error) {
	res, err := http.Head(d.url)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to access the site you provided: %s", d.url)
	}

	if res.Header.Get("Accept-Ranges") != "bytes" {
		return 0, errors.New("this site doesn't support a range request")
	}

	l, err := strconv.Atoi(res.Header.Get("Content-Length"))
	if err != nil {
		return 0, errors.Wrap(err, "failed to get Content-Length")
	}
	fmt.Printf("total length: %d bytes\n", l)

	return l, nil
}

func (d *Downloader) parseCommandLine() error {
	ops := cliOptions{}
	// パーサーを作成する。
	p := flags.NewParser(&ops, flags.HelpFlag)
	// 入力内容をパースする。
	_, err := p.ParseArgs(d.Argv)
	if err != nil {
		return errors.Wrap(err, "failed to parse command line")
	}

	// URL を設定
	if ops.Args.URL == "" {
		p.WriteHelp(os.Stdout)
		return fmt.Errorf("\n please check usage above")
	}
	d.url = ops.Args.URL

	// Name を設定
	if ops.Name != "" {
		d.name = ops.Name
	} else {
		if name := guessFileName(d.url); name == "" {
			return errors.Wrap(err, "please provide output file name")
		} else {
			d.name = name
		}
	}

	// Procs を設定
	d.procs = ops.Procs

	return nil
}

// ファイル名を推測する。
func guessFileName(URL string) string {
	s := strings.Split(URL, "/")
	return s[len(s)-1]
}
