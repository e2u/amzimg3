package lib

import "testing"

func TestAllowSource(t *testing.T) {
	src := []string{
		"127.0.0.1:8080",
		"211.95.79.211",
		"www.funguide.com.cn",
	}
	as := NewAllowSourceByArray(src)

	if as.Check("127.0.0.1:8080") != true {
		t.Fail()
	}
	if as.Check("noexists.com") == true {
		t.Fail()
	}
}
