package enums

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"HTTPs-Golang/logger"
)

type Config struct {
	IP        string `json:"ip"`
	Port      int    `json:"port"`
	LoginURL  string `json:"loginurl"`
	RateLimit int    `json:"ratelimit"`
	CDN       string `json:"cdn"`
}

type Server struct {
	Config            Config
	allowedUserAgents map[string]struct{}
	blockedIPPrefixes []string
	rateLimit         sync.Map
	serverDataResp    string
}

var loggers = logger.NewLogger("")

type rateLimitInfo struct {
	count     int
	timestamp time.Time
}

func NewServer() (*Server, error) {
	s := &Server{
		allowedUserAgents: map[string]struct{}{
			"UbiServices_SDK_2022.Release.9_PC64_ansi_static": {},
			"UbiServices_SDK_2022.Release.9_ANDROID64_static": {},
			"UbiServices_SDK_2022.Release.9_ANDROID32_static": {},
			"UbiServices_SDK_2022.Release.9_IOS64":            {},
		},
		blockedIPPrefixes: []string{"35.", "52.", "169.", "198.", "199.", "200.", "216.", "47."},
	}

	if err := s.loadConfig(); err != nil {
		loggers.Error("Failed to load configuration", map[string]interface{}{
			"error": err,
		})
		return nil, err
	}

	s.initServerData()
	return s, nil
}

func (s *Server) loadConfig() error {
	data, err := os.ReadFile(filepath.Join("config", "config.json"))
	if err != nil {
		return fmt.Errorf("reading config: %w", err)
	}

	if err := json.Unmarshal(data, &s.Config); err != nil {
		return fmt.Errorf("parsing config: %w", err)
	}

	return nil
}

func (s *Server) initServerData() {
	s.serverDataResp = strings.Join([]string{
		fmt.Sprintf("server|%s", s.Config.IP),
		fmt.Sprintf("port|%d", s.Config.Port),
		"type|1",
		fmt.Sprintf("loginurl|%s", s.Config.LoginURL),
		"beta_server|127.0.0.1",
		"beta_port|17091",
		"beta_type|1",
		fmt.Sprintf("meta|%d", time.Now().Unix()),
		"RTENDMARKERBS1001",
	}, "\n")
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clientIP := s.normalizeIP(r.Header.Get("X-Forwarded-For"))
	if clientIP == "" {
		clientIP = s.normalizeIP(r.RemoteAddr)
	}

	loggers.Info("Incoming request", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
		"client": clientIP,
	})

	if s.isBlockedIP(clientIP) {
		loggers.Warn("Blocked IP address", map[string]interface{}{
			"client": clientIP,
		})
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if !s.checkRateLimit(clientIP) {
		loggers.Warn("Rate limit exceeded", map[string]interface{}{
			"client": clientIP,
		})
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	switch r.URL.Path {
	case "/growtopia/server_data.php":
		s.handleServerData(w, r)
	default:
		if strings.HasPrefix(r.URL.Path, "/cache") || strings.Contains(r.URL.Path, "/0098") {
			s.handleCacheRequest(w, r)
		}
	}
}

func (s *Server) normalizeIP(ip string) string {
	ip = strings.Split(ip, ":")[0]
	if strings.HasPrefix(ip, "::ffff:") {
		return strings.TrimPrefix(ip, "::ffff:")
	}
	if ip == "::1" {
		return "127.0.0.1"
	}
	return ip
}

func (s *Server) isBlockedIP(ip string) bool {
	for _, prefix := range s.blockedIPPrefixes {
		if strings.HasPrefix(ip, prefix) {
			return true
		}
	}
	return false
}

func (s *Server) checkRateLimit(ip string) bool {
	now := time.Now()
	value, ok := s.rateLimit.LoadOrStore(ip, &rateLimitInfo{count: 1, timestamp: now})
	info := value.(*rateLimitInfo)

	if !ok {
		return true
	}

	if now.Sub(info.timestamp) > 30*time.Second {
		info.count = 1
		info.timestamp = now
		return true
	}

	info.count++
	return info.count <= s.Config.RateLimit
}

func (s *Server) handleServerData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if _, ok := s.allowedUserAgents[r.Header.Get("User-Agent")]; !ok {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	fmt.Fprint(w, s.serverDataResp)
}

func (s *Server) handleCacheRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if _, ok := s.allowedUserAgents[r.Header.Get("User-Agent")]; !ok {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	filePath := "." + r.URL.Path
	data, err := os.ReadFile(filePath)
	if err != nil {
		redirectURL := fmt.Sprintf("https://ubistatic-a.akamaihd.net/%s%s", s.Config.CDN, r.URL.Path)
		http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
		return
	}

	mimeType := getMIMEType(filePath)
	w.Header().Set("Content-Type", mimeType)
	w.Write(data)
}

func getMIMEType(filePath string) string {
	mimeTypes := map[string]string{
		".ico":  "image/x-icon",
		".html": "text/html",
		".js":   "text/javascript",
		".json": "application/json",
		".css":  "text/css",
		".png":  "image/png",
		".jpg":  "image/jpeg",
		".wav":  "audio/wav",
		".mp3":  "audio/mpeg",
		".svg":  "image/svg+xml",
		".pdf":  "application/pdf",
		".doc":  "application/msword",
	}

	ext := filepath.Ext(filePath)
	if mimeType, ok := mimeTypes[ext]; ok {
		return mimeType
	}
	return "application/octet-stream"
}
