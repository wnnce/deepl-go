package deepl

import (
	"fmt"
	"testing"
)

func TestAcquireTextTranslateParams(t *testing.T) {
	body := AcquireTextTranslateParams()
	body.SourceLang = "zh"
	body.Text = []string{"hello"}
	fmt.Println(body)
	RecycleParams(body)
	body = AcquireTextTranslateParams()
	fmt.Println(body)
}
