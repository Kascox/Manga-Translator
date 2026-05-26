package server

import (
	"bytes"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
	webpkg "github.com/user/manga-translator/web"
)

const (
	MaxBodySize        = 50 << 20 // 50 MB
	TranslationTimeout = 120 * time.Second
)

// NewHTTPServer 创建并返回配置好的 Kratos HTTP Server。
func NewHTTPServer(addr string, logger log.Logger) *kratoshttp.Server {
	srv := kratoshttp.NewServer(
		kratoshttp.Address(addr),
		kratoshttp.Timeout(30*time.Second),
		kratoshttp.Logger(logger),
		kratoshttp.Filter(corsFilter("http://localhost:5173")),
		kratoshttp.Filter(maxBytesFilter(MaxBodySize)),
	)

	registerRoutes(srv)
	return srv
}

func registerRoutes(srv *kratoshttp.Server) {
	// 健康检查路由
	r := srv.Route("/api")
	r.GET("/health", healthHandler)

	// 静态资源 + SPA fallback
	srv.HandlePrefix("/", spaHandler())
}

func healthHandler(ctx kratoshttp.Context) error {
	return ctx.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "manga-translator",
	})
}

// corsFilter 返回一个允许指定 Origin 跨域请求的 Filter。
func corsFilter(allowedOrigin string) kratoshttp.FilterFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "86400")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// maxBytesFilter 限制请求体大小，防止超大上传。
func maxBytesFilter(maxBytes int64) kratoshttp.FilterFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next.ServeHTTP(w, r)
		})
	}
}

// spaHandler 提供嵌入的静态文件，不存在的路径回退到 index.html（支持 SPA 路由）。
// 使用 http.ServeContent 而非 http.FileServer，避免对 /index.html 的内置 301 重定向。
func spaHandler() http.Handler {
	distFS, err := fs.Sub(webpkg.Dist, "dist")
	if err != nil {
		panic("failed to get dist sub-fs: " + err.Error())
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 跳过 API 路由，由 Kratos 路由器处理
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		p := strings.TrimPrefix(r.URL.Path, "/")
		if p == "" {
			p = "index.html"
		}

		// 尝试打开目标文件
		f, openErr := distFS.Open(p)
		if openErr != nil {
			// 文件不存在，SPA 路由回退到 index.html
			serveEmbedContent(w, r, distFS, "index.html")
			return
		}
		info, statErr := f.Stat()
		_ = f.Close()
		if statErr != nil || info.IsDir() {
			serveEmbedContent(w, r, distFS, "index.html")
			return
		}
		serveEmbedContent(w, r, distFS, p)
	})
}

// serveEmbedContent 使用 http.ServeContent 服务嵌入 FS 中的文件。
// 与 http.FileServer 不同，不会对 /index.html 产生 301 重定向。
func serveEmbedContent(w http.ResponseWriter, r *http.Request, fsys fs.FS, name string) {
	data, err := fs.ReadFile(fsys, name)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	http.ServeContent(w, r, name, time.Time{}, bytes.NewReader(data))
}
