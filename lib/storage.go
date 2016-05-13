package lib

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	UnSafePath = errors.New("unsafe path")
)

type Storager interface {
	ReadLocalFile(f string) ([]byte, error)
	CopyRemoteFile(remote string, local string) error
	Clean(f string) error
	Exists(f string) bool
	FullPath(f ...string) string
	IsSafePath(path string) bool
	MkdirAll(path string) error
}

type Storage struct {
	BaseDir string
}

func NewStorage(baseDir string) Storage {
	return Storage{
		BaseDir: baseDir,
	}
}

func (s Storage) ReadLocalFile(f string) ([]byte, error) {
	if !s.IsSafePath(f) {
		return nil, UnSafePath
	}
	fs, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer fs.Close()
	return ioutil.ReadAll(fs)
}

func (s Storage) CopyRemoteFile(rf string, lf string) error {
	if !strings.HasPrefix(rf, "http") {
		rf = "http://" + rf
	}
	return HttpClientGetToLocal(rf, lf)
}

func (s Storage) Clean(f string) error {
	files, err := ioutil.ReadDir(s.BaseDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		fullPath := s.FullPath(file.Name(), f)
		if !s.IsSafePath(fullPath) {
			continue
		}
		os.RemoveAll(fullPath)
	}
	return nil
}

func (s Storage) Exists(f string) bool {
	if _, err := os.Stat(f); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (s Storage) FullPath(fs ...string) string {

	var vs []string = []string{s.BaseDir}

	for _, f := range fs {
		vs = append(vs, filepath.Clean(strings.Replace(f, "..", "", -1)))
	}
	return strings.Join(vs, string(filepath.Separator))
}

func (s Storage) IsSafePath(path string) bool {
	if strings.HasPrefix(path, s.BaseDir) {
		return true
	}
	return false
}

func (s Storage) MkdirAll(path string) error {
	return os.MkdirAll(filepath.Dir(path), 0755)
}
