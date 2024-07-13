package sshs

import (
	"GinProject12/model"
	"bytes"
	"github.com/gorilla/websocket"
	"io"
	"sync"
)

//本文件定义了与虚拟机进行交互的相关对象

// 线程安全的缓冲区
type OutputDataBuffer struct {
	buffer bytes.Buffer
	mu     sync.Mutex
}

func (w *OutputDataBuffer) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer.Write(p)
}

func (w *OutputDataBuffer) Flush(ws *websocket.Conn, t string, connectID string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.buffer.Len() != 0 {
		str := w.buffer.String()
		err := ws.WriteJSON(model.ShellResponse{Code: 200, ConnectID: connectID, Type: t, Message: "", Data: str})
		if err != nil {
			return err
		}

		w.buffer.Reset()
	}
	return nil
}

func (w *OutputDataBuffer) WriteToHost(stdinPipe io.WriteCloser) {
	w.mu.Lock()
	defer w.mu.Unlock()
	stdinPipe.Write(w.buffer.Bytes())
	w.buffer.Reset()
}
