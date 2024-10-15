package graph

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"runtime"
)

type Node struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

type Edge struct {
	Source int `json:"source"`
	Target int `json:"target"`
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

func graphHandler(w http.ResponseWriter, r *http.Request) {
	// Example graph data
	graphData := Graph{
		Nodes: []Node{
			{ID: 1, URL: "https://example.com"},
			{ID: 2, URL: "https://example.com/about"},
			{ID: 3, URL: "https://example.com/contact"},
		},
		Edges: []Edge{
			{Source: 1, Target: 2},
			{Source: 1, Target: 3},
			{Source: 2, Target: 3},
		},
	}

	json.NewEncoder(w).Encode(graphData)
}

func RenderGraph() {
	go func() {
		openBrowser("http://localhost:6969/graph")
	}()
	http.HandleFunc("/graph", graphHandler)
	http.ListenAndServe(":6969", nil)
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = nil
	}

	if err != nil {
		panic(err)
	}
}
