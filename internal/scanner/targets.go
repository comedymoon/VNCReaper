package scanner

import (
	"bufio"
	"log"
	"net/netip"
	"os"
	"strings"
)

// ExpandTargets reads IPs or CIDRs from a file, expands CIDRs, and streams them to a channel.
// It respects the limit if set (>0) and skips huge ranges.
func ExpandTargets(targetFile string, limit int, verbose bool) <-chan string {
	out := make(chan string, 512)

	go func() {
		defer close(out)
		file, err := os.Open(targetFile)
		if err != nil {
			log.Fatalf("Failed to open target file: %v", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		count := 0

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			if strings.Contains(line, "/") { // CIDR
				prefix, err := netip.ParsePrefix(line)
				if err != nil {
					continue
				}

				// sanity check (optional): skip huge ranges
				ones := prefix.Bits()
				if 32-ones > 16 { // IPv4 only
					if verbose {
						log.Printf("[skip] CIDR %s too large (%d hosts)", prefix.String(), 1<<(32-ones))
					}
					continue
				}

				for ip := prefix.Masked().Addr(); prefix.Contains(ip); ip = ip.Next() {
					out <- ip.String()
					count++
					if limit > 0 && count >= limit {
						return
					}
				}
			} else { // Single IP
				if _, err := netip.ParseAddr(line); err == nil {
					out <- line
					count++
					if limit > 0 && count >= limit {
						return
					}
				}
			}
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Error reading target file: %v", err)
		}
	}()

	return out
}