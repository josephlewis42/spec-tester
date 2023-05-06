package storage

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"sort"

	v1 "github.com/josephlewis42/scheme-compliance/tester/model/v1"
	"sigs.k8s.io/yaml"
)

// LoadSuite loads a suite at the given path.
func LoadSuite(path string) (*Suite, error) {
	out := &Suite{
		RootPath: path,
	}
	var err error
	var testSuites []YamlFile[v1.TestSuite]

	testSuites, err = decode[v1.TestSuite](path, false)
	switch {
	case err != nil:
		return nil, err

	case len(testSuites) > 1:
		return nil, errors.New("conflicting TestSuite definitions")

	case len(testSuites) == 0:
		return nil, errors.New("missing TestSuite defintion")

	default:
		out.TestSuite = testSuites[0]
	}

	out.Implementations, err = decode[v1.Implementation](filepath.Join(path, "implementations"), true)
	if err != nil {
		return nil, err
	}

	out.Tests, err = decode[v1.Test](filepath.Join(path, "tests"), true)
	if err != nil {
		return nil, err
	}

	out.Specifications, err = decode[v1.Specification](filepath.Join(path, "specifications"), true)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func decode[T any](path string, recurse bool) ([]YamlFile[T], error) {
	var results []YamlFile[T]

	processedRoot := false

	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			switch {
			case recurse:
				return nil

			case !processedRoot:
				processedRoot = true
				return nil

			default:
				return filepath.SkipDir
			}
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

			results = append(results, YamlFile[T]{
				originalPath: path,
				originalData: bytes,
				Path:         path,
				Value:        tmp,
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
