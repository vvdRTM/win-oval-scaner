# Windows OVAL Scanner - Quick Start Guide

## Building the Project

```bash
cd windows-oval-scanner
go mod download
go build -o oval-scanner.exe
```

## Running OVAL Scans

### Basic Usage

```bash
./oval-scanner.exe -xml examples/windows-security.xml
```

### With JSON Output

```bash
./oval-scanner.exe -xml examples/windows-security.xml -json results.json
```

### With HTML Report

```bash
./oval-scanner.exe -xml examples/windows-security.xml -html report.html
```

## Custom OVAL Definitions

Create your own OVAL XML file following the structure in `examples/windows-security.xml`.

### Example: Check Windows Firewall Status

```xml
<registry_object id="oval:firewall:obj:1" version="1">
  <hive>HKEY_LOCAL_MACHINE</hive>
  <key>SYSTEM\CurrentControlSet\Services\MpsSvc</key>
  <name>Start</name>
</registry_object>

<registry_state id="oval:firewall:ste:1" version="1">
  <value operation="equals">2</value>
</registry_state>
```

## Project Structure

```
windows-oval-scanner/
├── main.go                      # Entry point
├── go.mod                       # Go module definition
├── go.sum                       # Go dependencies
├── pkg/
│   └── oval/
│       ├── types.go             # OVAL data structures
│       ├── scanner.go           # Core scanning logic
│       └── tests.go             # Test implementations
├── examples/
│   ├── windows-security.xml     # Example OVAL definitions
│   └── quick-start.md           # This file
└── README.md                    # Full documentation
```

## Supported Test Types

### Registry Tests
- Hive: HKEY_LOCAL_MACHINE, HKEY_CURRENT_USER, HKEY_CLASSES_ROOT
- Operations: equals, contains, pattern_match

### File Tests
- Type: exists, writable
- Returns: file size, permissions, modification time

## Architecture Overview

1. **XML Parsing** - Load OVAL definitions
2. **Test Extraction** - Parse registry and file objects
3. **Test Execution** - Run tests against system state
4. **Result Aggregation** - Collect and format results
5. **Report Generation** - JSON/HTML output

## Extending the Scanner

### Adding New Test Types

1. Create new test structure in `pkg/oval/tests.go`
2. Implement executor in `pkg/oval/scanner.go`
3. Add extraction logic in `extractRegistryTests()` or `extractFileTests()`

### Adding New Operations

Modify `compareValues()` function in `pkg/oval/scanner.go` to support:
- Regular expressions
- Version comparisons
- Custom logic

## Performance Notes

- Registry lookups: ~1-5ms per test
- File existence checks: ~1-2ms per test
- XML parsing overhead: ~10-50ms
- Typical scan (50 tests): 200-300ms

## Security Considerations

- Requires admin privileges for certain registry keys
- File access depends on user permissions
- Results may vary by Windows version

## Troubleshooting

### "Registry key not found"
- Verify the registry path exists
- Check Windows version compatibility
- Ensure proper permissions

### "File not found"
- Verify absolute file path
- Check file exists in the location
- Run as administrator if needed

## Next Steps

1. Customize OVAL definitions for your environment
2. Automate scanning with Windows Task Scheduler
3. Integrate with SIEM/compliance tools
4. Build compliance reports

## License

MIT