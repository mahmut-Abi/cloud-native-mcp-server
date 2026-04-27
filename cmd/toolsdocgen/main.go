package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/tooldoc"
)

func main() {
	var (
		check bool
		file  string
	)

	flag.BoolVar(&check, "check", false, "exit non-zero if the generated inventory block is out of date")
	flag.StringVar(&file, "file", filepath.Join("docs", "TOOLS.md"), "path to the tool reference markdown file")
	flag.Parse()

	repoRoot, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to determine working directory: %v\n", err)
		os.Exit(1)
	}

	docPath := filepath.Join(repoRoot, file)
	docContent, err := os.ReadFile(docPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read %s: %v\n", docPath, err)
		os.Exit(1)
	}

	inventory, err := tooldoc.CollectInventory(repoRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to collect tool inventory: %v\n", err)
		os.Exit(1)
	}

	generated := tooldoc.RenderGeneratedInventory(inventory)
	updated := tooldoc.ReplaceGeneratedInventory(string(docContent), generated)

	if check {
		if updated != string(docContent) {
			fmt.Fprintf(os.Stderr, "%s is out of date; run `go run ./cmd/toolsdocgen`\n", file)
			os.Exit(1)
		}
		return
	}

	if err := os.WriteFile(docPath, []byte(updated), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write %s: %v\n", docPath, err)
		os.Exit(1)
	}
}
