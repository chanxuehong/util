package http

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"sync"
	"time"
)

// TimeoutHandler returns a http.Handler that runs h with the given time limit.
//
// The new http.Handler calls h.ServeHTTP to handle each request, but if a
// call runs for longer than its time limit, the handler responds with
// the given http status code and the given message in its body.
// (If msg is empty, a suitable default message will be sent.)
// After such a timeout, writes by h to its http.ResponseWriter will return
// http.ErrHandlerTimeout.
//
// TimeoutHandler buffers all http.Handler writes to memory and does not
// support the http.Hijacker or http.Flusher interfaces.
func TimeoutHandler(h http.Handler, dt time.Duration, code int, msg string) http.Handler {
	return &timeoutHandler{
		handler: h,
		code:    code,
		body:    msg,
		dt:      dt,
	}
}

type timeoutHandler struct {
	handler http.Handler
	code    int
	body    string
	dt      time.Duration

	// When set, no context will be created and this context will
	// be used instead.
	testContext context.Context
}

func (h *timeoutHandler) errorCode() int {
	if h.code > 0 {
		return h.code
	}
	return http.StatusServiceUnavailable
}

func (h *timeoutHandler) errorBody() string {
	if h.body != "" {
		return h.body
	}
	return "<html><head><title>Timeout</title></head><body><h1>Timeout</h1></body></html>"
}

func (h *timeoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := h.testContext
	if ctx == nil {
		var cancelCtx context.CancelFunc
		ctx, cancelCtx = context.WithTimeout(r.Context(), h.dt)
		defer cancelCtx()
	}
	r = r.WithContext(ctx)
	done := make(chan struct{})
	tw := &timeoutWriter{
		w: w,
		h: make(http.Header),
	}
	go func() {
		h.handler.ServeHTTP(tw, r)
		close(done)
	}()
	select {
	case <-done:
		tw.mu.Lock()
		defer tw.mu.Unlock()
		dst := w.Header()
		for k, vv := range tw.h {
			dst[k] = vv
		}
		if !tw.wroteHeader {
			tw.code = http.StatusOK
		}
		w.WriteHeader(tw.code)
		w.Write(tw.wbuf.Bytes())
		tw.timedOut = true
		return
	case <-ctx.Done():
		tw.mu.Lock()
		defer tw.mu.Unlock()
		w.WriteHeader(h.errorCode())
		io.WriteString(w, h.errorBody())
		tw.timedOut = true
		return
	}
}

type timeoutWriter struct {
	w    http.ResponseWriter
	h    http.Header
	wbuf bytes.Buffer

	mu          sync.Mutex
	timedOut    bool
	wroteHeader bool
	code        int
}

func (tw *timeoutWriter) Header() http.Header { return tw.h }

func (tw *timeoutWriter) Write(p []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut {
		return 0, http.ErrHandlerTimeout
	}
	if !tw.wroteHeader {
		tw.writeHeader(http.StatusOK)
	}
	return tw.wbuf.Write(p)
}

func (tw *timeoutWriter) WriteHeader(code int) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut || tw.wroteHeader {
		return
	}
	tw.writeHeader(code)
}

func (tw *timeoutWriter) writeHeader(code int) {
	tw.wroteHeader = true
	tw.code = code
}

func (tw *timeoutWriter) WriteString(s string) (n int, err error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut {
		return 0, http.ErrHandlerTimeout
	}
	if !tw.wroteHeader {
		tw.writeHeader(http.StatusOK)
	}
	return tw.wbuf.WriteString(s)
}

func (tw *timeoutWriter) CloseNotify() <-chan bool {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if v, ok := tw.w.(http.CloseNotifier); ok {
		return v.CloseNotify()
	}
	return make(chan bool, 1)
}
