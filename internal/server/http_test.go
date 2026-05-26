package server

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// TestCORSHeaders 验证 CORS filter 正确设置响应头，OPTIONS 返回 204。
func TestCORSHeaders(t *testing.T) {
	handler := corsFilter("http://localhost:5173")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// 验证 OPTIONS 预检请求返回 204
	t.Run("OPTIONS preflight returns 204", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/api/health", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", rec.Code)
		}
		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
			t.Errorf("expected CORS origin header, got %q", got)
		}
		if got := rec.Header().Get("Access-Control-Allow-Methods"); got == "" {
			t.Error("Access-Control-Allow-Methods header missing")
		}
	})

	// 验证普通 GET 请求携带 CORS 头
	t.Run("GET request has CORS headers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
			t.Errorf("expected CORS origin header, got %q", got)
		}
	})
}

// TestSPAFallback 验证请求不存在的路径时回退到 index.html。
func TestSPAFallback(t *testing.T) {
	handler := spaHandler()

	req := httptest.NewRequest(http.MethodGet, "/some/non-existent-page", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 for SPA fallback, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "<div id=\"root\">") {
		t.Errorf("expected index.html content in SPA fallback, got: %s", body)
	}
}

// TestStaticFile 验证请求 /index.html 返回 200 和 HTML 内容。
func TestStaticFile(t *testing.T) {
	handler := spaHandler()

	req := httptest.NewRequest(http.MethodGet, "/index.html", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 for /index.html, got %d", rec.Code)
	}
	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Errorf("expected text/html Content-Type, got %q", contentType)
	}
}

// TestHealthCheck 启动真实 Kratos HTTP Server，验证 /api/health 端点。
func TestHealthCheck(t *testing.T) {
	// 用 :0 让 OS 分配随机端口，Endpoint() 会创建监听器并返回实际地址
	srv := NewHTTPServer("localhost:0", log.DefaultLogger)

	endpoint, err := srv.Endpoint()
	if err != nil {
		t.Fatalf("failed to get server endpoint: %v", err)
	}
	addr := endpoint.Host // e.g. "127.0.0.1:54321"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() { _ = srv.Start(ctx) }()

	// 等待服务器就绪
	time.Sleep(100 * time.Millisecond)
	defer func() { _ = srv.Stop(ctx) }()

	resp, err := http.Get("http://" + addr + "/api/health")
	if err != nil {
		t.Fatalf("request to /api/health failed: %v", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)
	if !strings.Contains(bodyStr, `"status"`) || !strings.Contains(bodyStr, `"ok"`) {
		t.Errorf("unexpected health response body: %s", bodyStr)
	}
	if !strings.Contains(bodyStr, `"service"`) {
		t.Errorf("missing 'service' field in health response: %s", bodyStr)
	}
}
