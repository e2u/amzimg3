package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var urlFixRegexp1 = regexp.MustCompile("^(.*)(http:|https:)\\/{0,}(.*)$")
var accessID uint64
var mux map[string]func(http.ResponseWriter, *http.Request) = map[string]func(http.ResponseWriter, *http.Request){"/favicon.ico": favicon}


type Server struct {
	Address string
	Port    uint
	BaseDir string //TODO 从参数中取 baseDir
}

type LogLine struct {
	Time         time.Time `json:"time"`
	AccessID     uint64    `json:"aid"`
	Remote       string    `json:"remote"`
	Method       string    `json:"method"`
	RequestURI   string    `json:"req"`
	UserAgent    string    `json:"ua"`
	Referer      string    `json:"refer"`
	ResponseCode int       `json:"resp_code"`
	Length       int       `json:"length"`
	During       string    `json:"during"`
}

func (l *LogLine) ToString() string {
	b, _ := json.Marshal(l)
	return string(b)
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start() error {
	tl := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", s.Address, s.Port),
		Handler:        http.HandlerFunc(s.handler),
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 12,
	}
	return tl.ListenAndServe()
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}

	var (
		fixedURL string
		err      error
		startAt  time.Time
	)
	accessID++

	remoteIP := func() string {
		ip, port, _ := net.SplitHostPort(r.RemoteAddr)
		return ip + ":" + port
	}
	startAt = time.Now()

	logline := LogLine{
		AccessID:   accessID,
		Time:       time.Now(),
		Remote:     remoteIP(),
		Method:     r.Method,
		RequestURI: r.RequestURI,
		UserAgent:  r.UserAgent(),
		Referer:    r.Referer(),
	}

	defer func() {
		logline.During = fmt.Sprintf("%fs", time.Now().Sub(startAt).Seconds())
		log.Println(logline.ToString())
	}()

	errorPage := func(err error) {
		ol := strconv.Itoa(len(err.Error()))

		logline.ResponseCode = 500
		logline.Length = len(err.Error())

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Expires", "Sun, 15 Jan 1970 08:18:55 GMT")
		w.Header().Set("Server", HeaderServer)
		w.Header().Set("Content-Length", ol)
		io.WriteString(w, err.Error())
	}

	imageOutput := func(b []byte, contentType string) {
		ol := len(b)

		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Expires", "Sun, 15 Jan 2100 08:18:55 GMT")
		w.Header().Set("Cache-Control", "max-age=31104000")
		w.Header().Set("Server", HeaderServer)
		w.Header().Set("Connection", "keep-alive")

		logline.ResponseCode = 200
		logline.Length = ol

		w.Header().Set("Content-Length", fmt.Sprintf("%d", ol))
		w.Write(b)
	}

	if fixedURL, err = s.fixURL(r.RequestURI); err != nil {
		log.Println(err)
		errorPage(err)
		return
	}

	img := NewImage(fixedURL)
	img.Reload = r.Header.Get("reload") == "true"
	img.Storage = NewStorage(s.BaseDir)
	if err := img.Output(imageOutput); err != nil {
		errorPage(err)
		return
	}
}

func (s *Server) fixURL(u string) (rs string, err error) {

	rs = strings.Replace(u, "..", "", -1)
	rs = filepath.Clean(rs)
	if rs, err = filepath.Abs(rs); err != nil {
		return
	}

	rs = urlFixRegexp1.ReplaceAllString(rs, "${1}${2}//${3}")
	return
}
