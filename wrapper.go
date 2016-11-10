package main

import (
	"github.com/gorilla/websocket"
	"io"
	"sync"
)

type WebSocketWrapper struct {
	ws     *websocket.Conn
	reader io.Reader
	writer io.WriteCloser
	// websockets support only concurrent reader and one writer
	readMu  sync.Mutex
	writeMu sync.Mutex
}

func (wrapper *WebSocketWrapper) Read(p []byte) (n int, err error) {
	wrapper.readMu.Lock()
	defer wrapper.readMu.Unlock()

	if wrapper.reader == nil {
		_, wrapper.reader, err = wrapper.ws.NextReader()
		if err != nil {
			return 0, err
		}
	}
	for n = 0; n < len(p); {
		var m int
		m, err = wrapper.reader.Read(p[n:])
		n += m
		if err == io.EOF {
			wrapper.reader = nil
			break
		}
		if err != nil {
			break
		}
	}
	return
}

func (wrapper *WebSocketWrapper) Write(p []byte) (n int, err error) {
	wrapper.writeMu.Lock()
	defer wrapper.writeMu.Unlock()

	if wrapper.writer == nil {
		wrapper.writer, err = wrapper.ws.NextWriter(websocket.TextMessage)
		if err != nil {
			return 0, err
		}
	}
	for n = 0; n < len(p); {
		var m int
		m, err = wrapper.writer.Write(p)
		n += m
		if err != nil {
			break
		}
	}
	if err != nil || n == len(p) {
		err = wrapper.Close()
	}
	return
}

func (wrapper *WebSocketWrapper) Close() (err error) {
	if wrapper.writer != nil {
		err = wrapper.writer.Close()
		wrapper.writer = nil
	}
	return err
}
