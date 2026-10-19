package parser

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

type libraryTransport func(*http.Request) (*http.Response, error)

func (f libraryTransport) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func TestLibraryClientIsolationAndSharingText(t *testing.T) {
	clients := make([]*Client, 2)
	for i := range clients {
		title := fmt.Sprintf("client-%d", i)
		clients[i] = NewClient(&http.Client{Transport: libraryTransport(func(req *http.Request) (*http.Response, error) {
			if req.URL.Hostname() != "xhslink.cn" {
				t.Errorf("unexpected host: %s", req.URL.Hostname())
			}
			body := fmt.Sprintf(`<script>window.__INITIAL_STATE__={"noteData":{"data":{"noteData":{"title":%q,"imageList":[{"url":"https://cdn.example/a.jpg"}]}}}}</script>`, title)
			return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body)), Request: req}, nil
		})})
	}
	var wg sync.WaitGroup
	for n := range 24 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			i := n % 2
			result, err := clients[i].Parse(context.Background(), "文字 https://xhslink.cn/o/example 复制打开")
			if err != nil {
				t.Error(err)
				return
			}
			if result.Title != fmt.Sprintf("client-%d", i) {
				t.Errorf("client configuration leaked: %q", result.Title)
			}
		}()
	}
	wg.Wait()
}

func TestLibraryCancellation(t *testing.T) {
	client := NewClient(&http.Client{Timeout: time.Second, Transport: libraryTransport(func(req *http.Request) (*http.Response, error) {
		<-req.Context().Done()
		return nil, req.Context().Err()
	})})
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	start := time.Now()
	_, err := client.Parse(ctx, "https://xhslink.cn/o/example")
	if err != context.DeadlineExceeded {
		t.Fatalf("expected caller deadline, got %v", err)
	}
	if time.Since(start) > 500*time.Millisecond {
		t.Fatal("request did not use caller context")
	}
}

func TestLibraryRejectsUnregisteredHostsBeforeIO(t *testing.T) {
	client := NewClient(&http.Client{Transport: libraryTransport(func(*http.Request) (*http.Response, error) {
		t.Error("untrusted input reached transport")
		return nil, fmt.Errorf("unexpected I/O")
	})})
	for _, url := range []string{"https://xhslink.cn.evil.test/o/x", "https://user:secret@xhslink.cn/o/x"} {
		if _, err := client.Parse(context.Background(), url); err == nil {
			t.Fatal("expected invalid source")
		}
	}
}
