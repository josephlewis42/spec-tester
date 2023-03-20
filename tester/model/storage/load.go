package storage

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	v1 "github.com/josephlewis42/scheme-compliance/tester/model/v1"
	"sigs.k8s.io/yaml"
)

// LoadSuite loads a suite at the given path.
func LoadSuite(path string) (*Suite, error) {
	out := &Suite{
		RootPath: path,
	}
	var err error
	out.Implementations, err = decode[v1.Implementation](path, "implementations")
	if err != nil {
		return nil, err
	}

	out.Tests, err = decode[v1.Test](path, "tests")
	if err != nil {
		return nil, err
	}

	out.Specifications, err = decode[v1.Specification](path, "specifications")
	if err != nil {
		return nil, err
	}

	return out, nil
}

func decode[T any](root, subdir string) ([]File[T], error) {
	var results []File[T]

	err := filepath.Walk(filepath.Join(root, subdir), func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		switch filepath.Ext(info.Name()) {
		case ".yaml", ".yml", ".json":
			var tmp T
			bytes, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("couldn't read %s: %e", path, err)
			}
			if err := yaml.UnmarshalStrict(bytes, &tmp); err != nil {
				return fmt.Errorf("couldn't decode %s: %e", path, err)
			}

			results = append(results, File[T]{
				Path:  strings.TrimPrefix(path, root),
				Value: tmp,
			})
		}

		return nil
	})

	// Ensure inputs are deterministic.
	sort.Slice(results, func(i, j int) bool {
		return results[i].Path < results[i].Path
	})

	return results, err
}
