package lib

import (
	"bufio"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
)

type AllowSource struct {
	m map[string]bool
	sync.RWMutex
}

func (as *AllowSource) Set(src string) bool {
	as.Lock()
	defer as.Unlock()
	as.m[src] = true
	_, ok := as.m[src]
	return ok
}

func (as *AllowSource) Check(src string) bool {
	as.RLock()
	defer as.RUnlock()

	host, port, err := splitURLHostPort(src)
	if err != nil {
		return false
	}
	if v, ok := as.m[host+":"+port]; ok {
		return v
	}
	return false
}

func NewAllowSourceByArray(src []string) *AllowSource {
	as := &AllowSource{m: make(map[string]bool)}
	for _, v := range src {
		as.Set(v)
	}
	return as
}

func NewAllowSourceByFile(f string) *AllowSource {
	as := &AllowSource{m: make(map[string]bool)}
	file, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	log.Println("read allow source list: ", f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		host, port, err := splitURLHostPort(line)
		if err != nil {
			log.Println("parse url err:", err.Error())
			continue
		}
		log.Println("add allow source: " + host)
		as.Set(host + ":" + port)
	}
	return as
}

func splitURLHostPort(s string) (string, string, error) {
	if !strings.HasPrefix(s, "http") {
		s = "http://" + s
	}
	purl, err := url.Parse(s)
	if err != nil {
		return "", "", err
	}
	rs := strings.Split(purl.Host, ":")
	if len(rs) <= 1 {
		return rs[0], "80", nil
	}
	return rs[0], rs[1], nil
}
