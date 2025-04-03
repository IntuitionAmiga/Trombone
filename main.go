package main

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// Evidence represents findings from each detection technique
type Evidence struct {
	Technique string
	Findings  []string
	Positive  bool
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: rust-detector <binary_path>")
		os.Exit(1)
	}

	binaryPath := os.Args[1]
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Check if file exists
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		logger.Error("Binary file does not exist", "path", binaryPath)
		os.Exit(1)
	}

	// Create context with timeout for commands
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Collect evidence from all techniques
	allEvidence := []Evidence{
		checkRustStrings(ctx, binaryPath),
		checkSymbolNames(ctx, binaryPath),
		checkElfInfo(ctx, binaryPath),
		checkMemoryAllocator(ctx, binaryPath),
		checkPanicHandlers(ctx, binaryPath),
	}

	// Determine if it's a Rust binary based on evidence
	isRustBinary, confidence := analyseEvidence(allEvidence)

	// Display results
	if isRustBinary {
		fmt.Printf("✅ The binary '%s' appears to be compiled from Rust (confidence: %s)\n\n", binaryPath, confidence)
	} else {
		fmt.Printf("❌ The binary '%s' does not appear to be compiled from Rust\n\n", binaryPath)
	}

	// Display detailed evidence
	fmt.Println("Detailed analysis:")
	fmt.Println("=================")
	for _, e := range allEvidence {
		resultSymbol := "❌"
		if e.Positive {
			resultSymbol = "✅"
		}

		fmt.Printf("%s %s:\n", resultSymbol, e.Technique)
		if len(e.Findings) > 0 {
			for _, finding := range e.Findings {
				fmt.Printf("  • %s\n", finding)
			}
		} else {
			fmt.Println("  • No evidence found")
		}
		fmt.Println()
	}
}

// checkRustStrings looks for Rust-specific string patterns
func checkRustStrings(ctx context.Context, binaryPath string) Evidence {
	evidence := Evidence{
		Technique: "Rust-specific strings",
		Findings:  []string{},
		Positive:  false,
	}

	cmd := exec.CommandContext(ctx, "strings", binaryPath)
	output, err := cmd.Output()
	if err != nil {
		return evidence
	}

	rustPatterns := []string{
		"rust", "rustc", "cargo",
		"core::result::Result",
		"core::option::Option",
		"std::panic",
		"alloc::",
		"thread 'main' panicked",
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		for _, pattern := range rustPatterns {
			if strings.Contains(strings.ToLower(line), pattern) {
				evidence.Findings = append(evidence.Findings, line)
				evidence.Positive = true
				break
			}
		}
		// Limit the number of findings to avoid overwhelming output
		if len(evidence.Findings) >= 5 {
			evidence.Findings = append(evidence.Findings, "... (more findings omitted)")
			break
		}
	}

	return evidence
}

// checkSymbolNames looks for Rust-specific symbol name patterns
func checkSymbolNames(ctx context.Context, binaryPath string) Evidence {
	evidence := Evidence{
		Technique: "Rust symbol names",
		Findings:  []string{},
		Positive:  false,
	}

	cmd := exec.CommandContext(ctx, "nm", "-a", binaryPath)
	output, err := cmd.Output()
	if err != nil {
		return evidence
	}

	rustSymbolRegex := regexp.MustCompile(`_ZN.*rust.*`)
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if rustSymbolRegex.MatchString(line) {
			evidence.Findings = append(evidence.Findings, line)
			evidence.Positive = true
		}
		// Limit the number of findings
		if len(evidence.Findings) >= 5 {
			evidence.Findings = append(evidence.Findings, "... (more findings omitted)")
			break
		}
	}

	return evidence
}

// checkElfInfo analyses ELF headers for Rust-specific sections and symbols
func checkElfInfo(ctx context.Context, binaryPath string) Evidence {
	evidence := Evidence{
		Technique: "ELF headers and sections",
		Findings:  []string{},
		Positive:  false,
	}

	cmd := exec.CommandContext(ctx, "readelf", "-a", binaryPath)
	output, err := cmd.Output()
	if err != nil {
		return evidence
	}

	rustPatterns := []string{
		"rust", "rustc",
		".rust_",
		"core::fmt",
		"core::panicking",
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		for _, pattern := range rustPatterns {
			if strings.Contains(strings.ToLower(line), pattern) {
				trimmedLine := strings.TrimSpace(line)
				evidence.Findings = append(evidence.Findings, trimmedLine)
				evidence.Positive = true
				break
			}
		}
		// Limit the number of findings
		if len(evidence.Findings) >= 5 {
			evidence.Findings = append(evidence.Findings, "... (more findings omitted)")
			break
		}
	}

	return evidence
}

// checkMemoryAllocator looks for Rust's memory allocator symbols
func checkMemoryAllocator(ctx context.Context, binaryPath string) Evidence {
	evidence := Evidence{
		Technique: "Rust memory allocator",
		Findings:  []string{},
		Positive:  false,
	}

	cmd := exec.CommandContext(ctx, "nm", "-D", binaryPath)
	output, err := cmd.Output()
	if err != nil {
		return evidence
	}

	rustAllocFunctions := []string{
		"__rdl_alloc",
		"__rdl_dealloc",
		"__rdl_realloc",
		"__rdl_alloc_zeroed",
		"__rust_alloc",
		"__rust_dealloc",
		"__rust_realloc",
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		for _, funcName := range rustAllocFunctions {
			if strings.Contains(line, funcName) {
				evidence.Findings = append(evidence.Findings, line)
				evidence.Positive = true
				break
			}
		}
	}

	return evidence
}

// checkPanicHandlers looks for Rust's panic handling code
func checkPanicHandlers(ctx context.Context, binaryPath string) Evidence {
	evidence := Evidence{
		Technique: "Rust panic handlers",
		Findings:  []string{},
		Positive:  false,
	}

	cmd := exec.CommandContext(ctx, "objdump", "-d", binaryPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	output, err := cmd.Output()
	if err != nil {
		// objdump might fail with some binaries, but we can continue
		return evidence
	}

	// Look for panic-related functions in disassembly
	panicPatterns := []string{
		"panic",
		"core::panicking",
		"rust_panic",
		"rust_begin_unwind",
	}

	// Process output to find panic handlers
	lines := strings.Split(string(output), "\n")
	for i, line := range lines {
		for _, pattern := range panicPatterns {
			if strings.Contains(strings.ToLower(line), pattern) {
				// Get a few lines of currentContext for the finding
				startLine := max(0, i-2)
				endLine := min(len(lines), i+3)
				currentContext := strings.Join(lines[startLine:endLine], "\n")
				evidence.Findings = append(evidence.Findings, "Found panic-related code:\n"+currentContext)
				evidence.Positive = true
				break
			}
		}
		// Limit the number of findings
		if len(evidence.Findings) >= 3 {
			evidence.Findings = append(evidence.Findings, "... (more findings omitted)")
			break
		}
	}

	return evidence
}

// analyseEvidence determines if the binary is Rust based on the collected evidence
func analyseEvidence(allEvidence []Evidence) (bool, string) {
	positiveCount := 0
	for _, e := range allEvidence {
		if e.Positive {
			positiveCount++
		}
	}

	// Determine confidence level
	var confidence string
	switch positiveCount {
	case 0:
		return false, "none"
	case 1:
		confidence = "low"
	case 2:
		confidence = "medium"
	case 3:
		confidence = "high"
	default:
		confidence = "very high"
	}

	return positiveCount > 0, confidence
}
