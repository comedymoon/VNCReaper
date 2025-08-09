package scanner

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Common HTTP/WebSocket ports used by noVNC
var NoVNCPorts = []int{
	6080, 6081, 6082, 8888, 8080, 8081, 80, 443,
}

// DetectNoVNC scans known ports for HTML/JS/WebSocket code matching noVNC
func DetectNoVNC(ip string) (bool, int, string) {
	client := http.Client{
		Timeout: 2 * time.Second,
	}

	for _, port := range NoVNCPorts {
		url := fmt.Sprintf("http://%s:%d/", ip, port)
		resp, err := client.Get(url)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		body := strings.ToLower(string(data))

		if strings.Contains(body, "novnc") ||
			strings.Contains(body, "websockify") ||
			strings.Contains(body, "connect('ws") ||
			strings.Contains(body, "vnc_password") ||
			strings.Contains(body, "rfb") {

			return true, port, url
		}
	}
	return false, 0, ""
}