# Bad words filter with golang
## implement fastcheck with golang

```
go get github.com/ZhiqinYang/fastcheck

var check = NewFastCheck(true)
check.AddWord("fuck")
check.AddWord("草泥马")
check.AddWord("abc")
check.AddWord("中国")
fmt.Println(check.HasWord("fuck"))
```
