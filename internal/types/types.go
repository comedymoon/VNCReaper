package types

import "time"

type ScanResult struct {
	IP        string `json:"ip"`
	Port      string `json:"port"`
	Protocol  string `json:"protocol"`
	Banner    string `json:"banner"`
	Title     string `json:"title"`
	Favicon   string `json:"favicon_hash"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Raw       string `json:"raw_response"`
}

type Job struct {
	IP   string
	Port string
}

var CommonPorts = []string{
    "5900", "5901", "5902", "5903", // VNC direct
    "5800", "5801", "5802", "5803", // VNC over HTTP
    "6080", "6081", "6082",         // Common noVNC ports
    "6789",                         // Some embedded web VNC servers
    "8080", "8081", "8090", "8091", // Alternate HTTP ports
}

// Performance counters
var (
	ScannedCount int64
	FoundCount   int64
	StartTime    time.Time
)