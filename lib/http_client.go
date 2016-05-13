package lib

import (
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func timeoutDialer(connTimeout int, rwTimeout int) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, time.Duration(connTimeout)*time.Second)
		if err != nil {
			log.Printf("Failed to connect to [%s]. Timed out after %d seconds\n", addr, rwTimeout)
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(time.Duration(rwTimeout) * time.Second))
		return conn, nil
	}
}

func HttpClientGetToLocal(u string, l string) (err error) {
	var (
		req  *http.Request
		resp *http.Response
		out  *os.File
	)

	os.MkdirAll(filepath.Dir(l), 0755)

	transport := http.Transport{
		Dial: timeoutDialer(ConnTimeout, ReadWriteTimeout),
		ResponseHeaderTimeout: time.Duration(ReadWriteTimeout) * time.Second,
		DisableCompression:    false,
		DisableKeepAlives:     true,
		MaxIdleConnsPerHost:   20,
	}

	defer transport.CloseIdleConnections()

	client := &http.Client{
		Transport: &transport,
	}

	if req, err = http.NewRequest("GET", u, nil); err != nil {
		return
	}

	req.Close = true

	if resp, err = client.Do(req); err != nil {
		return
	}

	if resp.StatusCode >= 400 {
		err = errors.New(resp.Status)
		return
	}

	defer resp.Body.Close()

	out, err = os.Create(l)
	if err != nil {
		return
	}
	_, err = io.Copy(out, resp.Body)
	return
}
