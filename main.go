package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var (
		xmlFile = flag.String("xml", "", "Path to OVAL XML definition file")
		jsonOut = flag.String("json", "", "Output results as JSON")
		htmlOut = flag.String("html", "", "Output results as HTML report")
	)
	flag.Parse()

	if *xmlFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	log.Println("🛡️  Windows OVAL Scanner")
	log.Printf("Loading OVAL definitions from: %s\n", *xmlFile)

	xmlData, err := os.ReadFile(*xmlFile)
	if err != nil {
		log.Fatalf("Error reading XML file: %v", err)
	}

	ctx := context.Background()
	scanner := NewOVALScanner()

	results, err := scanner.ScanFromXML(ctx, xmlData)
	if err != nil {
		log.Fatalf("Error scanning: %v", err)
	}

	log.Printf("✓ Executed %d tests\n", len(results))
	
	// Display results
	displayResults(results)

	// Output JSON if requested
	if *jsonOut != "" {
		if err := outputJSON(results, *jsonOut); err != nil {
			log.Fatalf("Error writing JSON: %v", err)
		}
		log.Printf("✓ JSON report saved to: %s\n", *jsonOut)
	}

	// Output HTML if requested
	if *htmlOut != "" {
		if err := outputHTML(results, *htmlOut); err != nil {
			log.Fatalf("Error writing HTML: %v", err)
		}
		log.Printf("✓ HTML report saved to: %s\n", *htmlOut)
	}
}

func displayResults(results []TestResult) {
	fmt.Println("\n" + "="*60)
	fmt.Println("TEST RESULTS")
	fmt.Println("="*60)

	passed := 0
	failed := 0

	for _, result := range results {
		status := "✓"
		if result.Status != "pass" {
			status = "✗"
			failed++
		} else {
			passed++
		}

		fmt.Printf("%s [%s] %s\n", status, result.TestID, result.Message)
	}

	fmt.Println("="*60)
	fmt.Printf("Passed: %d | Failed: %d\n", passed, failed)
	fmt.Println("="*60)
}

func outputJSON(results []TestResult, filename string) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func outputHTML(results []TestResult, filename string) error {
	// TODO: Implement HTML report generation
	return nil
}