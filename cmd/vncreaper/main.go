package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/comedymoon/VNCReaper/internal/scanner"
	"github.com/comedymoon/VNCReaper/internal/types"
	"github.com/comedymoon/VNCReaper/internal/gui"
)

func main() {
	targetFile := flag.String("i", "targets.txt", "Input file with target IPs or CIDRs")
	outputFile := flag.String("o", "results.json", "Output file for scan results")
	guiMode := flag.Bool("gui", false, "Enable GUI mode (launches simple local web viewer)")
	limit := flag.Int("limit", 0, "Limit total number of IPs to scan (0 = unlimited)")
	port := flag.Int("port", 7777, "Port for GUI mode")
	verbose := flag.Bool("v", false, "Verbose output")
	threads := flag.Int("t", 2000, "Number of concurrent threads (default: 2000)")
	timeout := flag.Int("timeout", 800, "Connection timeout in milliseconds (default: 800ms)")
	httpOnly := flag.Bool("http-only", false, "Skip TCP banner checks, HTTP/noVNC only")
	noFavicon := flag.Bool("no-favicon", false, "Skip favicon hash calculation (faster)")
	novncDisabled := flag.Bool("novnc-disabled", false, "Disable noVNC detection")
	flag.Parse()

	if _, err := os.Stat(*targetFile); os.IsNotExist(err) {
		fmt.Printf("Error: target file %s not found.\n\n", *targetFile)
		flag.Usage()
		os.Exit(1)
	}

	if *guiMode {
		gui.StartGUI(*outputFile, *port)
		return
	}

	types.StartTime = time.Now()

	out, err := os.Create(*outputFile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer out.Close()

	resultChan := make(chan types.ScanResult, 1000)
	jobs := make(chan types.Job, *threads*4)
	var wg sync.WaitGroup

	for i := 0; i < *threads; i++ {
		wg.Add(1)
		go scanner.ScanWorker(
			jobs,
			resultChan,
			&wg,
			time.Duration(*timeout)*time.Millisecond,
			*httpOnly,
			*noFavicon,
			*novncDisabled,
		)
	}


	writerDone := make(chan bool)
	go func() {
		for result := range resultChan {
			if result.Status == "open" && (result.Protocol == "RFB" || (!*novncDisabled && result.Protocol == "noVNC")) {
				atomic.AddInt64(&types.FoundCount, 1)
				data, _ := json.Marshal(result)
				fmt.Fprintln(out, string(data))
				if *verbose {
					fmt.Printf("FOUND VNC: %s:%s (%s)\n", result.IP, result.Port, result.Protocol)
				}
			}
		}
		writerDone <- true
	}()

	if *verbose {
		go func() {
			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()
			for range ticker.C {
				scanned := atomic.LoadInt64(&types.ScannedCount)
				found := atomic.LoadInt64(&types.FoundCount)
				elapsed := time.Since(types.StartTime)
				rate := float64(scanned) / elapsed.Seconds()
				fmt.Printf("Progress: %d scanned, %d found, %.0f scans/sec\n", scanned, found, rate)
			}
		}()
	}

	go func() {
		for ip := range scanner.ExpandTargets(*targetFile, *limit, *verbose) {
			for _, port := range types.CommonPorts {
				jobs <- types.Job{IP: ip, Port: port}
			}
		}
		close(jobs)
	}()

	wg.Wait()
	close(resultChan)
	<-writerDone

	elapsed := time.Since(types.StartTime)
	scanned := atomic.LoadInt64(&types.ScannedCount)
	found := atomic.LoadInt64(&types.FoundCount)
	rate := float64(scanned) / elapsed.Seconds()

	fmt.Printf("\nScan complete! Scanned %d targets in %v\n", scanned, elapsed)
	fmt.Printf("Found %d VNC services. Rate: %.0f scans/sec\n", found, rate)
	fmt.Printf("Results saved to %s\n", *outputFile)
}