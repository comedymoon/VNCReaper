package scanner

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/comedymoon/VNCReaper/internal/types"
)

func ScanWorker(
    jobs <-chan types.Job,
    results chan<- types.ScanResult,
    wg *sync.WaitGroup,
    timeout time.Duration,
    httpOnly bool,
    noFavicon bool,
    disableNoVNC bool,
) {
    defer wg.Done()
    for job := range jobs {
        result := fastScanTarget(job.IP, job.Port, timeout, httpOnly, noFavicon)
        atomic.AddInt64(&types.ScannedCount, 1)
        results <- result

        if !disableNoVNC {
            if found, port, url := DetectNoVNC(job.IP); found {
                results <- types.ScanResult{
                    IP:        job.IP,
                    Port:      fmt.Sprintf("%d", port),
                    Protocol:  "HTTP(noVNC)",
                    Banner:    "noVNC Web Client",
                    Status:    "open",
                    Timestamp: time.Now().Format(time.RFC3339),
                    Raw:       url,
                }
            }
        }
    }
}


func fastScanTarget(ip, port string, timeout time.Duration, httpOnly bool, noFavicon bool) types.ScanResult {
	result := types.ScanResult{
		IP:        ip,
		Port:      port,
		Status:    "closed",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	address := net.JoinHostPort(ip, port)

	// TCP connection test
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return result
	}
	defer conn.Close()
	result.Status = "open"

	// Fast RFB check (if not HTTP-only mode)
	if !httpOnly {
		conn.SetReadDeadline(time.Now().Add(timeout))
		buf := make([]byte, 12) // RFB version is exactly 12 bytes
		n, _ := conn.Read(buf)
		if n >= 3 && strings.HasPrefix(string(buf), "RFB") {
			result.Protocol = "RFB"
			result.Banner = strings.TrimSpace(string(buf[:n]))
			result.Raw = hex.EncodeToString(buf[:n])
			return result
		}
	}

	// Fast HTTP check for noVNC
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			DisableKeepAlives:   true,
			MaxIdleConns:        0,
			IdleConnTimeout:     timeout,
			TLSHandshakeTimeout: timeout,
		},
	}

	resp, err := client.Get("http://" + address)
	if err != nil {
		return result
	}
	defer resp.Body.Close()

	// Read only first 2KB for speed
	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
	body := string(bodyBytes)

	result.Protocol = "HTTP"
	result.Banner = fmt.Sprintf("%d %s", resp.StatusCode, http.StatusText(resp.StatusCode))

	// Fast title extraction
	if titleStart := strings.Index(body, "<title>"); titleStart != -1 {
		if titleEnd := strings.Index(body[titleStart:], "</title>"); titleEnd != -1 {
			result.Title = strings.TrimSpace(body[titleStart+7 : titleStart+titleEnd])
		}
	}

	// Store minimal raw response
	if len(body) > 200 {
		result.Raw = body[:200] + "..."
	} else {
		result.Raw = body
	}

	// Fast favicon hash (if enabled)
	if !noFavicon {
		if faviconHash := getFaviconHashFast("http://"+address+"/favicon.ico", timeout); faviconHash != "" {
			result.Favicon = faviconHash
		}
	}

	if isLikelyNoVNC(body, result.Title, result.Favicon) {
		result.Protocol = "noVNC"
	}

	return result
}

func isLikelyNoVNC(body, title, favicon string) bool {
	body = strings.ToLower(body)
	title = strings.ToLower(title)

	// Strong indicators in body
	keywords := []string{
		"vnc_auto.html", "vnc.html", "websockify", "rfb.js",
		"novnc", "novnc_container", "rfb.connect", "disconnectbutton", "viewport-container",
	}
	for _, kw := range keywords {
		if strings.Contains(body, strings.ToLower(kw)) {
			return true
		}
	}

	// Title hints
	if strings.Contains(title, "vnc") || strings.Contains(title, "novnc") {
		return true
	}

	// Favicon hash checks (optional, weak indicator)
	knownFavicon := map[string]bool{
		"1531707515":  true,
		"-1454792475": true,
		"-1699606641": true,
	}
	if knownFavicon[favicon] {
		return true
	}

	return false
}

func getFaviconHashFast(url string, timeout time.Duration) string {
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			DisableKeepAlives: true,
			MaxIdleConns:      0,
		},
	}
	resp, err := client.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return ""
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 5120))
	if err != nil {
		return ""
	}
	hash := md5.Sum(data)
	return fmt.Sprintf("%d", int32(hashToInt(hash[:]))) // store as signed 32-bit
}

func hashToInt(b []byte) int32 {
	sum := int32(0)
	for i := 0; i < len(b); i++ {
		sum += int32(b[i]) << (uint(i%4) * 8)
	}
	return sum
}