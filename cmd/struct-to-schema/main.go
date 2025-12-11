// Command struct-to-schema generates OpenAPI schemas from Go struct definitions.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"text/template"

	"go.yaml.in/yaml/v4"
)

var (
	typeName  = flag.String("type", "", "Type name to generate schema for (required)")
	pkgPath   = flag.String("pkg", "", "Package import path (required)")
	output    = flag.String("output", "", "Output file path (required)")
	internal  = flag.String("internal", "", "Value for x-internal marker (optional)")
	simpleStr = regexp.MustCompile(`^[a-zA-Z0-9_./:-]+$`)
)

const generatorTemplate = `package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/archesai/archesai/internal/jsonschema"
	target "{{.PkgPath}}"
)

func main() {
	os.Setenv("JSONSCHEMAGODEBUG", "typeschemasnull=1")

	var v target.{{.TypeName}}
	t := reflect.TypeOf(v)
	schema, err := jsonschema.ForType(t, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error generating schema: %v\n", err)
		os.Exit(1)
	}

	jsonBytes, err := json.Marshal(schema)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshaling schema: %v\n", err)
		os.Exit(1)
	}

	os.Stdout.Write(jsonBytes)
}
`

func needsQuoting(s string) bool {
	if s == "" {
		return true
	}
	if simpleStr.MatchString(s) {
		return false
	}
	return true
}

func setYAMLStyle(node *yaml.Node, isKey bool) {
	switch node.Kind {
	case yaml.DocumentNode, yaml.SequenceNode:
		node.Style = 0
		for _, child := range node.Content {
			setYAMLStyle(child, false)
		}
	case yaml.MappingNode:
		node.Style = 0
		for i, child := range node.Content {
			setYAMLStyle(child, i%2 == 0)
		}
	case yaml.ScalarNode:
		if isKey {
			node.Style = 0
		} else if node.Tag == "!!str" {
			if needsQuoting(node.Value) {
				node.Style = yaml.SingleQuotedStyle
			} else {
				node.Style = 0
			}
		}
	}
}

func main() {
	flag.Parse()

	if *typeName == "" || *pkgPath == "" || *output == "" {
		fmt.Fprintln(
			os.Stderr,
			"Usage: struct-to-schema -type <TypeName> -pkg <package/path> -output <output.yaml> [-internal <marker>]",
		)
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Find module root
	modRoot := findModuleRoot()
	if modRoot == "." {
		fmt.Fprintln(os.Stderr, "error: could not find go.mod in parent directories")
		os.Exit(1)
	}

	// Create temp file within module to access internal packages
	tmpDir := filepath.Join(modRoot, ".tmp-schema-gen")
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "error creating temp dir: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Generate the generator program
	tmpl, err := template.New("generator").Parse(generatorTemplate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing template: %v\n", err)
		os.Exit(1)
	}

	var progBuf bytes.Buffer
	err = tmpl.Execute(&progBuf, map[string]string{
		"PkgPath":  *pkgPath,
		"TypeName": *typeName,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error executing template: %v\n", err)
		os.Exit(1)
	}

	// Write generator program
	genFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(genFile, progBuf.Bytes(), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing generator: %v\n", err)
		os.Exit(1)
	}

	// Run the generator from module root
	runCmd := exec.Command("go", "run", genFile)
	runCmd.Dir = modRoot
	var stdout, stderr bytes.Buffer
	runCmd.Stdout = &stdout
	runCmd.Stderr = &stderr
	if err := runCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error running generator: %v\n%s\n", err, stderr.String())
		os.Exit(1)
	}

	// Convert JSON to YAML with proper styling
	var node yaml.Node
	if err := yaml.Unmarshal(stdout.Bytes(), &node); err != nil {
		fmt.Fprintf(os.Stderr, "error unmarshaling to YAML: %v\n", err)
		os.Exit(1)
	}

	setYAMLStyle(&node, false)

	// Add x-internal marker if specified
	if *internal != "" && node.Kind == yaml.DocumentNode && len(node.Content) > 0 {
		rootMap := node.Content[0]
		if rootMap.Kind == yaml.MappingNode {
			rootMap.Content = append(rootMap.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "x-internal"},
				&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: *internal},
			)
		}
	}

	// Marshal to YAML
	var yamlBuf bytes.Buffer
	enc := yaml.NewEncoder(&yamlBuf)
	enc.SetIndent(2)
	if err := enc.Encode(&node); err != nil {
		fmt.Fprintf(os.Stderr, "error encoding YAML: %v\n", err)
		os.Exit(1)
	}
	if err := enc.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "error closing encoder: %v\n", err)
		os.Exit(1)
	}

	// Ensure output directory exists
	outDir := filepath.Dir(*output)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "error creating output dir: %v\n", err)
		os.Exit(1)
	}

	// Write output file
	if err := os.WriteFile(*output, yamlBuf.Bytes(), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing output: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated schema for %s.%s -> %s\n", *pkgPath, *typeName, *output)
}

// findModuleRoot finds the root of the current Go module
func findModuleRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "."
		}
		dir = parent
	}
}
