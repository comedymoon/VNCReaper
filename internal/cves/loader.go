package cves

import "github.com/comedymoon/VNCReaper/internal/types"

// CVE is the interface every exploit module must follow
type CVE interface {
	Name() string
	Description() string
	Exploit(target types.ScanResult) (bool, string) // returns (success, details)
}

var registry []CVE

// Register adds a CVE to the registry
func Register(c CVE) {
	registry = append(registry, c)
}

// RunAll executes all registered CVEs against a target
func RunAll(target types.ScanResult) []string {
	var hits []string
	for _, c := range registry {
		if ok, info := c.Exploit(target); ok {
			hits = append(hits, c.Name()+" → "+info)
		}
	}
	return hits
}

// GetAll returns all registered CVEs
func GetAll() []CVE {
	return registry
}