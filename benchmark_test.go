package puff_test

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"testing"
)

func BenchmarkGET(b *testing.B) {
	oncepuffserver()

	u, _ := url.Parse("http://127.0.0.1:7465/test")
	req := &http.Request{
		Method: "GET",
		URL:    u,
	}

	b.ResetTimer()
	for range b.N {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			b.Logf("unexpected error: %s", err.Error())
			b.FailNow()
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			b.Logf("unexpected error: %s", err.Error())
			resp.Body.Close()
			b.FailNow()
		}
		_ = body // preventing compiler optimizations; imagine if this was actually used
		resp.Body.Close()
	}
}

func BenchmarkLargePOST(b *testing.B) {
	oncepuffserver()

	u, _ := url.Parse("http://127.0.0.1:7465/data")
	req := &http.Request{
		Method: "POST",
		URL:    u,
		Body:   io.NopCloser(bytes.NewBuffer(randomdata())),
	}
	b.SetBytes(5e6)
	b.ResetTimer()
	for range b.N {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			b.Logf("unexpected error: %s", err.Error())
			b.FailNow()
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			b.Logf("unexpected error: %s", err.Error())
			resp.Body.Close()
			b.FailNow()
		}
		_ = body // preventing compiler optimizations; imagine if this was actually used
		resp.Body.Close()
	}
}

func BenchmarkNetHTTPLargePOST(b *testing.B) {
	oncenethttpserver()

	u, _ := url.Parse("http://127.0.0.1:7467/data")
	req := &http.Request{
		Method: "POST",
		URL:    u,
		Body:   io.NopCloser(bytes.NewBuffer(randomdata())),
	}
	b.SetBytes(5e6)
	b.ResetTimer()

	for range b.N {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			b.Logf("unexpected error: %s", err.Error())
			b.FailNow()
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			b.Logf("unexpected error: %s", err.Error())
			resp.Body.Close()
			b.FailNow()
		}
		_ = body // preventing compiler optimizations; imagine if this was actually used
		resp.Body.Close()
	}
}

func BenchmarkNetHTTPGet(b *testing.B) {
	oncenethttpserver()
	u, _ := url.Parse("http://127.0.0.1:7467/test")
	req := &http.Request{
		Method: "GET",
		URL:    u,
	}

	b.ResetTimer()
	for range b.N {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			b.Logf("unexpected error: %s", err.Error())
			b.FailNow()
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			b.Logf("unexpected error: %s", err.Error())
			resp.Body.Close()
			b.FailNow()
		}
		_ = body // preventing compiler optimizations; imagine if this was actually used
		resp.Body.Close()
	}
}
