package fastcheck

import (
	"fmt"
	"strings"
	"testing"
)

var check = NewFastCheck(true)

func init() {
	check.AddWord("fuck")
	check.AddWord("草泥马")
	check.AddWord("abc")
	check.AddWord("中国")
	fmt.Println("init done")
	fmt.Println(strings.ToUpper("abc中国"))
}

func Test_check(t *testing.T) {
	t.Log(check.HasWord("fckyou") == false)
	t.Log(check.HasWord("草泥马124") == true, check.ReplaceWith("草泥马124", '?'))
	t.Log(check.HasWord("草泥A") == false)
	t.Log(check.HasWord("abc中国") == true, check.ReplaceWith("abc中国", '*'))
	t.Log(check.HasWord(" GM的狗") == false)
	t.Log(check.HasWord("FUCK"), check.ReplaceWith("FUCK", '*'))
	t.Log(check.HasWord("fuck"))
	t.Log(check.HasWord("fuck") == true, check.ReplaceWith("fuck", '?'))

}

func BenchmarkHasword(b *testing.B) {
	b.Log(check.HasWord("FUCK"))
	b.Log(check.HasWord("fuck"))
}
