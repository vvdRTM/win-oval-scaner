package oval

// RegistryTest represents a Windows registry test
type RegistryTest struct {
	ID            string
	Hive          string
	Key           string
	ValueName     string
	ExpectedValue string
	Operation     string
}

// FileTest represents a file system test
type FileTest struct {
	ID       string
	FilePath string
	TestType string
}