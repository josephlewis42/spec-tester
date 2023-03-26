package storage

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"sigs.k8s.io/yaml"
)

// NewYamlFile creates a new YamlFile with the given type.
func NewYamlFile[T any](path string) *YamlFile[T] {
	return &YamlFile[T]{
		Path: path,
	}
}

// OpenYamlFile opens an existing YAML file with the given type.
func OpenYamlFile[T any](path string) (*YamlFile[T], error) {
	var tmp T
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("couldn't read %s: %e", path, err)
	}
	if err := yaml.UnmarshalStrict(bytes, &tmp); err != nil {
		return nil, fmt.Errorf("couldn't decode %s: %e", path, err)
	}

	return &YamlFile[T]{
		originalPath: path,
		originalData: bytes,
		Path:         path,
		Value:        tmp,
	}, nil
}

// YamlFile maps an on-disk file with a specific parsed value from that file.
type YamlFile[T any] struct {
	originalPath string
	originalData []byte

	Path  string
	Value T
}

func (file *YamlFile[T]) marshal() ([]byte, error) {
	return yaml.Marshal(file.Value)
}

// Diff returns a human-readable diff of the file compared to its stored version.
// The output isn't stable.
func (file *YamlFile[T]) Diff() string {
	// Build intro message.
	pathDiff := ""
	switch {
	case file.originalPath == "":
		pathDiff = fmt.Sprintf("new file: %s", file.Path)
	case file.Path == "":
		pathDiff = fmt.Sprintf("deleted: %s", file.originalPath)
	case file.originalPath != file.Path:
		pathDiff = fmt.Sprintf("renamed: %q -> %q", file.originalPath, file.Path)
	}

	// Coalesce data into the same type for diffing
	orig := make(map[string]any)
	curr := make(map[string]any)

	if err := yaml.Unmarshal(file.originalData, &orig); err != nil {
		return fmt.Sprintf("%s\nERROR: Invalid original data: %s\n", file.Path, err.Error())
	}

	currBytes, err := file.marshal()
	if err != nil {
		return fmt.Sprintf("%s\nERROR: Couldn't marshal current data: %s\n", file.Path, err.Error())
	}

	if err := yaml.Unmarshal(currBytes, &curr); err != nil {
		return fmt.Sprintf("%s\nERROR: Invalid current data: %s\n", file.Path, err.Error())
	}

	// Generate and return diff.
	delta := cmp.Diff(orig, curr, cmpopts.EquateEmpty())

	switch {
	case delta == "" && pathDiff == "": // No change
		return ""

	case delta == "" && pathDiff != "": // Path change
		return fmt.Sprintf("%s\n [file contents match]\n", pathDiff)

	case pathDiff == "": // Content only change
		return fmt.Sprintf("modified: %s\n%s\n", file.Path, delta)

	default: // Content and path change
		return fmt.Sprintf("%s\n%s\n", pathDiff, delta)
	}
}

// Save updates the file on disk.
func (file *YamlFile[T]) Save() error {
	// Write new file.
	if file.Path != "" {
		newContents, err := file.marshal()
		if err != nil {
			return fmt.Errorf("couldn't marshal file: %e", err)
		}

		ioutil.WriteFile(file.Path, newContents, 0600)
	}

	// File was deleted or moved.
	if file.Path == "" || file.Path != file.originalPath {
		return os.Remove(file.originalPath)
	}

	return nil
}
