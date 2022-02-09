package rangedownloader

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
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

	return 0
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
