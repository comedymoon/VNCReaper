package gui

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/comedymoon/VNCReaper/internal/types"
)

// StartGUI starts the web interface to view scan results
func StartGUI(outputFile string, port int) {
	tmpl := template.Must(template.ParseFiles("internal/gui/web.tmpl"))

	// Main HTML page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	// Data endpoint
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile(outputFile)
		if err != nil {
			// File doesn't exist yet — return empty results
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("[]"))
			return
		}

		if len(data) == 0 {
			// Empty file — return empty results
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("[]"))
			return
		}

		var results []types.ScanResult
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if len(strings.TrimSpace(line)) == 0 {
				continue
			}
			var res types.ScanResult
			if err := json.Unmarshal([]byte(line), &res); err == nil {
				results = append(results, res)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	})

	fmt.Printf("[GUI MODE] Web server running on http://localhost:%d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		fmt.Println("Failed to start web server:", err)
	}
}