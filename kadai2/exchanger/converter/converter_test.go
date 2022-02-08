package converter

import (
	"gopher-dojo/kadai2/exchanger/testing/helper"
	"os"
	"testing"
)

func TestConvertEtx(t *testing.T) {
	tests := []struct {
		src         string
		from        string
		to          string
		resultCount int
		wantError   bool
	}{
		{"testdata/sample", "jpg", "png", 2, false},
		{"testdata/sample", "PNG", "jpeg", 7, false},
		{"testdata/sample", "png", "jpg", 7, false},
		{"testdata/sample", "hoge", "jpg", 0, true},
		{"testdata/sample", "jpg", "hoge", 0, true},
		{"testdata/sample", "jpg", "jpg", 0, true},
		{"testdata/sample", "", "", 0, true},
	}

	if err := os.MkdirAll("output", 0777); err != nil {
		t.Error("failed to make an output folder")
	}

	for _, tt := range tests {
		count, err := ConvertEtx(tt.src, tt.from, tt.to)
		// ちゃんとエラーが発生したかの確認
		helper.TestWantError(t, err, tt.wantError)
		// 期待した通りの結果が返ってきたかの確認
		if count != tt.resultCount {
			t.Errorf("ConvertExt(%v, %v, %v) = %d, got %d",
				tt.src, tt.from, tt.to, tt.resultCount, count)
		}
	}

	// テスト完了後の掃除
	if err := os.RemoveAll("output"); err != nil {
		t.Error("failed to delete an output folder")
	}
}
