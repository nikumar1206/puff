package puff_test

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ThePuffProject/puff"
)

type MockResponseWriter struct {
	httptest.ResponseRecorder
}

func (m *MockResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	conn := &net.TCPConn{}
	return conn, bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}

func TestNewContext(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	ctx := puff.NewContext(w, req)
	if ctx == nil {
		t.Log("expected non-nil context, got nil context")
		t.FailNow()
	}
	if ctx.Request != req {
		t.Error("unexpected request field")
	}
	if ctx.ResponseWriter != w {
		t.Error("unexpected response writer field")
	}
}

func TestGetSet(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	ctx := puff.NewContext(w, req)

	key := "hello"
	value := "world"
	ctx.Set(key, value)

	rv := ctx.Get("hello")
	if rv != value {
		t.Errorf("expected value: %v, got: %v", value, rv)
	}
}

func TestGetRequestHeader(t *testing.T) {
	w := httptest.NewRecorder()
	key := "X-TestHeader"
	value := "test_header_value"
	req := &http.Request{
		Method: "GET",
		Header: http.Header{},
	}
	req.Header.Set(key, value)

	ctx := puff.NewContext(w, req)
	if ctx.GetRequestHeader(key) != value {
		t.Errorf("expected header value %s for key %s, got %s", value, key, ctx.GetRequestHeader(key))
	}
}

func TestSetResponseHeader(t *testing.T) {
	w := httptest.NewRecorder()
	key := "X-TestResponseHeader"
	value := "test_responseheader_value"

	req := &http.Request{
		Method: "GET",
	}

	ctx := puff.NewContext(w, req)

	ctx.SetResponseHeader(key, value)

	if w.Result().Header.Get(key) != value {
		t.Errorf("expected header value %s for key %s, got %s", value, key, w.Result().Header.Get(key))
	}
}

func TestGetBody(t *testing.T) {
	w := httptest.NewRecorder()
	body := []byte("this is a test body!")
	req, _ := http.NewRequest("GET", "http://127.0.0.1/helloworld", bytes.NewBuffer(body))

	ctx := puff.NewContext(w, req)
	got, err := ctx.GetBody()
	if err != nil {
		t.Logf("unexpected get body error: %s", err.Error())
		t.FailNow()
	}

	if !bytes.Equal(got, body) {
		t.Errorf("expected body %v, got %v", body, got)
	}
}

func TestSendResponseGeneric(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://127.0.0.1/helloworld", nil)
	ctx := puff.NewContext(w, req)

	text := "hello world :|"
	ctx.SendResponse(puff.GenericResponse{
		StatusCode:  200,
		Content:     text,
		ContentType: "text/plain",
	})

	body, err := io.ReadAll(w.Body)
	if err != nil {
		t.Logf("unexpected body buf read error: %s", err.Error())
		t.FailNow()
	}
	if w.Code != 200 {
		t.Errorf("expected status code 200, got code %d", w.Code)
	}
	if !bytes.Equal(body, []byte(text)) {
		t.Logf("unexpected body buf read error: %s", err.Error())
	}
}

func TestHijack(t *testing.T) {
	w := new(MockResponseWriter)
	req := &http.Request{
		Method: "GET",
	}

	ctx := puff.NewContext(w, req)

	conn, buf, err := ctx.Hijack()
	if err != nil {
		t.Errorf("unexpected error during hijacking: %s", err.Error())
	}

	if conn == nil || buf == nil {
		t.Errorf("conn and/or buf is/are nil, expected non-nil values from hijack")
	}
}
