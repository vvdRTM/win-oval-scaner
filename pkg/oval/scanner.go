package oval

import (
	"context"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/sys/windows/registry"
)

// OVALScanner performs OVAL security assessments
type OVALScanner struct {
}

// NewOVALScanner creates a new OVAL scanner instance
func NewOVALScanner() *OVALScanner {
	return &OVALScanner{}
}

// ScanFromXML parses and executes OVAL definitions from XML
func (s *OVALScanner) ScanFromXML(ctx context.Context, xmlData []byte) ([]TestResult, error) {
	var def OVALDefinition
	if err := xml.Unmarshal(xmlData, &def); err != nil {
		return nil, fmt.Errorf("error parsing XML: %w", err)
	}

	var results []TestResult

	// Execute registry tests
	for _, test := range extractRegistryTests(def) {
		result, err := s.executeRegistryTest(ctx, test)
		if err != nil {
			result = TestResult{
				TestID:  test.ID,
				Status:  "error",
				Message: fmt.Sprintf("Error: %v", err),
			}
		}
		result.Timestamp = time.Now().Format(time.RFC3339)
		results = append(results, result)
	}

	// Execute file tests
	for _, test := range extractFileTests(def) {
		result, err := s.executeFileTest(ctx, test)
		if err != nil {
			result = TestResult{
				TestID:  test.ID,
				Status:  "error",
				Message: fmt.Sprintf("Error: %v", err),
			}
		}
		result.Timestamp = time.Now().Format(time.RFC3339)
		results = append(results, result)
	}

	return results, nil
}

// executeRegistryTest runs a registry test
func (s *OVALScanner) executeRegistryTest(ctx context.Context, test RegistryTest) (TestResult, error) {
	result := TestResult{
		TestID:  test.ID,
		Details: make(map[string]interface{}),
	}

	// Open registry key
	hive := registryStringToHive(test.Hive)
	key, err := registry.OpenKey(hive, test.Key, registry.QUERY_VALUE)
	if err != nil {
		result.Status = "fail"
		result.Message = fmt.Sprintf("Registry key not found: %s\\%s", test.Hive, test.Key)
		return result, nil
	}
	defer key.Close()

	// Read value
	value, _, err := key.GetStringValue(test.ValueName)
	if err != nil {
		result.Status = "fail"
		result.Message = fmt.Sprintf("Registry value not found: %s", test.ValueName)
		return result, nil
	}

	// Compare values
	matches := compareValues(value, test.ExpectedValue, test.Operation)
	if matches {
		result.Status = "pass"
		result.Message = fmt.Sprintf("Registry check passed: %s = %s", test.ValueName, value)
	} else {
		result.Status = "fail"
		result.Message = fmt.Sprintf("Expected '%s', got '%s'", test.ExpectedValue, value)
	}

	result.Details["hive"] = test.Hive
	result.Details["key"] = test.Key
	result.Details["value_name"] = test.ValueName
	result.Details["expected"] = test.ExpectedValue
	result.Details["actual"] = value

	return result, nil
}

// executeFileTest runs a file test
func (s *OVALScanner) executeFileTest(ctx context.Context, test FileTest) (TestResult, error) {
	result := TestResult{
		TestID:  test.ID,
		Details: make(map[string]interface{}),
	}

	result.Details["file_path"] = test.FilePath

	// Check if file exists
	info, err := os.Stat(test.FilePath)
	if err != nil {
		result.Status = "fail"
		result.Message = fmt.Sprintf("File not found: %s", test.FilePath)
		return result, nil
	}

	switch test.TestType {
	case "exists":
		result.Status = "pass"
		result.Message = fmt.Sprintf("File exists: %s (%d bytes)", test.FilePath, info.Size())
		result.Details["size"] = info.Size()
		result.Details["modified"] = info.ModTime().String()

	case "writable":
		// Try to open file for writing (without actually writing)
		if info.Mode()&0200 != 0 {
			result.Status = "pass"
			result.Message = fmt.Sprintf("File is writable: %s", test.FilePath)
		} else {
			result.Status = "fail"
			result.Message = fmt.Sprintf("File is not writable: %s", test.FilePath)
		}
		result.Details["permissions"] = fmt.Sprintf("%o", info.Mode().Perm())

	default:
		result.Status = "unknown"
		result.Message = fmt.Sprintf("Unknown file test type: %s", test.TestType)
	}

	return result, nil
}

// Helper functions
func registryStringToHive(hive string) registry.Key {
	switch strings.ToUpper(hive) {
	case "HKEY_LOCAL_MACHINE":
		return registry.LOCAL_MACHINE
	case "HKEY_CURRENT_USER":
		return registry.CURRENT_USER
	case "HKEY_CLASSES_ROOT":
		return registry.CLASSES_ROOT
	default:
		return registry.LOCAL_MACHINE
	}
}

func compareValues(actual, expected, operation string) bool {
	switch operation {
	case "equals":
		return actual == expected
	case "contains":
		return strings.Contains(actual, expected)
	case "pattern_match":
		// Simple pattern matching (can be enhanced with regex)
		return strings.Contains(actual, strings.TrimSuffix(strings.TrimPrefix(expected, ".*"), ".*"))
	default:
		return false
	}
}

func extractRegistryTests(def OVALDefinition) []RegistryTest {
	var tests []RegistryTest
	// TODO: Implement extraction logic from XML
	return tests
}

func extractFileTests(def OVALDefinition) []FileTest {
	var tests []FileTest
	// TODO: Implement extraction logic from XML
	return tests
}