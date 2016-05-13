package lib

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/disintegration/imaging"
)

type Image struct {
	Reload bool

	source    string
	width     int
	thumb     string
	localSrc  string
	imageType string
	Storage   Storager
	fixedURL  string
}

type Lock struct {
	m    map[string]time.Time
	lock sync.RWMutex
}

var AllowRemoteSource *AllowSource

func NewLock() *Lock {
	return &Lock{m: make(map[string]time.Time)}
}

func (l *Lock) Lock(s string) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	_, ok := l.m[s]
	if !ok {
		l.m[s] = time.Now()
	}
	return nil
}

func (l *Lock) Unlock(s string) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	_, ok := l.m[s]
	if ok {
		delete(l.m, s)
	}
	return nil
}

func (l *Lock) IsLock(s string) bool {
	l.lock.RLock()
	defer l.lock.RUnlock()
	_, ok := l.m[s]
	return ok
}

var ThumbLock *Lock = NewLock()

var re0 *regexp.Regexp = regexp.MustCompile("^/show\\?.+$")     //'/show?src'
var re1 *regexp.Regexp = regexp.MustCompile("^/([0-9]+)/(.+)$") //'/:width/:src'
var re2 *regexp.Regexp = regexp.MustCompile("^/r/(.+)$")        // /r/:src
var re3 *regexp.Regexp = regexp.MustCompile("^/(.+)$")          // /:src

func NewImage(s string) *Image {
	return &Image{fixedURL: s}
}

func (img *Image) parsePath(u string) error {
	var err error
	var purl *url.URL
	w := func(ws string) int {
		if i, err := strconv.Atoi(ws); err == nil {
			return i
		}
		return DefaultWidth
	}

	srcURL := func(ws string) string {
		if is := strings.IndexByte(ws, '?'); is > 0 {
			return string(ws[0:is])
		}
		return ws
	}

	switch {
	case re0.MatchString(u):
		purl, err = url.Parse(u)
		if err != nil {
			return err
		}
		for k, ps := range purl.Query() {
			switch k {
			case "w":
				img.width = w(ps[0])
			case "src":
				img.source = srcURL(ps[0])
			}
		}
	case re1.MatchString(u): //'/:width/:src'
		re1r := re1.FindStringSubmatch(u)
		img.width = w(re1r[1])
		img.source = srcURL(re1r[2])

	case re2.MatchString(u): // /r/:src
		re2r := re2.FindStringSubmatch(u)
		img.source = srcURL(re2r[1])
		img.width = DefaultWidth

	case re3.MatchString(u): // /:src
		re3r := re3.FindStringSubmatch(u)
		img.source = srcURL(re3r[1])
		img.width = RawWidth

	default:
		return errors.New("unknow error")
	}

	if img.width > MaxWidth {
		img.width = MaxWidth
	}

	img.localSrc = img.Storage.FullPath("src", img.getLocalFileName(img.source))

	if img.width > 0 {
		img.thumb = img.Storage.FullPath(strconv.Itoa(img.width), img.getLocalFileName(img.source))
	}

	ext := strings.ToLower(filepath.Ext(img.source))
	switch {
	case strings.HasPrefix(ext, ".jpg"):
		img.imageType = FormatJpg
	case strings.HasPrefix(ext, ".jpeg"):
		img.imageType = FormatJpg
	case strings.HasPrefix(ext, ".png"):
		img.imageType = FormatPng
	case strings.HasPrefix(ext, ".gif"):
		img.imageType = FormatGif
	default:
		img.imageType = FormatJpg
	}

	// last fixed source path
	if !strings.HasPrefix(img.source, "http") {
		img.source = "http://" + img.source
	}

	return nil

}

func (img *Image) getLocalFileName(path string) string {
	path = strings.Replace(path, "http://", "", -1)
	path = strings.Replace(path, "https://", "", -1)
	return path
}

func (img *Image) Output(w func(b []byte, contentType string)) error {
	if err := img.parsePath(img.fixedURL); err != nil {
		return err
	}

	if !AllowRemoteSource.Check(img.source) {
		return errors.New("Forbidden Access Remote Sources")
	}

	if img.Reload {
		img.Storage.Clean(img.source)
	}

	if !img.Storage.Exists(img.localSrc) {
		if err := img.Storage.CopyRemoteFile(img.source, img.localSrc); err != nil {
			return errors.New("CopyRemoteFile: " + err.Error())
		}
	}

	if img.width <= 0 {
		if outByte, err := img.Storage.ReadLocalFile(img.localSrc); err != nil {
			return err
		} else {
			w(outByte, img.imageType)
			return nil
		}
	}

	if img.Storage.Exists(img.thumb) {
		if outByte, err := img.Storage.ReadLocalFile(img.thumb); err == nil {
			w(outByte, img.imageType)
			return err
		}
	}

	createThumb := func() error {
		ThumbLock.Lock(img.thumb)
		defer ThumbLock.Unlock(img.thumb)

		// create thumb
		in, err := os.Open(img.localSrc)
		if err != nil {
			return err
		}
		m, _, err := image.Decode(in)
		if err != nil {
			return err
		}

		resized := imaging.Resize(m, img.width, 0, imaging.Box)
		img.Storage.MkdirAll(img.thumb)
		out, err := os.Create(img.thumb)
		if err != nil {
			return err
		}
		defer out.Close()

		switch img.imageType {
		case FormatJpg:
			err = jpeg.Encode(out, resized, nil)
		case FormatPng:
			err = png.Encode(out, resized)
		case FormatGif:
			err = gif.Encode(out, resized, &gif.Options{NumColors: 256})
		default:
			err = jpeg.Encode(out, resized, nil)
		}
		if err != nil {
			return err
		}
		outByte, err := img.Storage.ReadLocalFile(img.thumb)
		if err != nil {
			return err
		}
		w(outByte, img.imageType)
		return nil
	}

	tryCount := 0
try:
	tryCount++
	if ThumbLock.IsLock(img.thumb) && tryCount < 5 {
		log.Println("retry...", img.thumb)
		time.Sleep(100 * time.Millisecond)
		goto try
	}

	return createThumb()
}
