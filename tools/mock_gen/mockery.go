package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run mockgen.go <package_directory>")
		os.Exit(1)
	}

	packageDir := os.Args[1]
	packageName := filepath.Base(packageDir)
	testPackageName := packageName + "_test"
	outputPath := packageDir + "mock_test.go"

	fmt.Println("Generating mocks...")
	cmd := exec.Command("mockery", "--name", ".*", "--testonly", "--print", "--with-expecter=true", "--dir", packageDir, "--outpkg", testPackageName)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to generate mocks: %v\n", err)
		os.Exit(1)
	}
	if _, err := os.Stat(outputPath); err == nil {
		err := os.Remove(outputPath)
		if err != nil {
			fmt.Println("Error deleting file:", err)
		} else {
			fmt.Println("File successfully deleted.")
		}
	}
	f, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Failed to create output file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	f.Write([]byte(processContent(string(outb.Bytes()))))
}

func processContent(content string) string {
	lines := strings.Split(content, "\n")
	var newLines []string
	var imports []string
	packageLine := ""
	importBlock := false

	for _, line := range lines {
		if strings.HasPrefix(line, "package") {
			packageLine = line // Save the package line
		} else if strings.HasPrefix(line, "import (") {
			importBlock = true
		} else if importBlock {
			if line == ")" {
				importBlock = false
			} else if line != "" {
				imports = append(imports, line)
			}
		} else if packageLine == "" || importBlock {
		} else {
			newLines = append(newLines, line)
		}
	}

	var result []string
	if packageLine != "" {
		result = append(result, packageLine)
		result = append(result, "")
	}
	if len(imports) > 0 {
		result = append(result, "import (")
		result = append(result, imports...)
		result = append(result, ")")
		result = append(result, "")
	}
	result = append(result, newLines...)

	return strings.Join(result, "\n")
}
