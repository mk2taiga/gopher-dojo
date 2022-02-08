package word

import (
	"bufio"
	"os"
)

const questionCnt = 60

func GetWords() ([]string, error) {
	path := "words.txt"
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := make(map[string]struct{})
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		lines[sc.Text()] = struct{}{}
	}

	words := make([]string, 0, len(lines))
	cnt := 0
	for line := range lines {
		words = append(words, line)
		cnt++
		if cnt >= questionCnt {
			break
		}
	}

	return words, nil
}
