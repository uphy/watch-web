package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type (
	configDirectory struct {
		parent *configDirectory
		// baseDirectory which is the relative path from the parent
		baseDirectory string
		// paths are the additional relative paths from the baseDirectory
		paths []string
	}
)

/*
 * /(baseDirectory)
 *   config.yml
 *   sources/           : re-usable sources
 *     siteA.yml
 *     siteB.yml
 *   constants/         : constant files for 'constant' source
 *     siteC.html
 *   test/source/data   : test data
 *     001_string.yml
 *     constants/       : constant files referenced by test data
 *       siteA.html
 */

func newConfigDirectory(baseDirectory string) *configDirectory {
	return &configDirectory{nil, baseDirectory, make([]string, 0)}
}

func (d *configDirectory) resolve(file string) (string, error) {
	if filepath.IsAbs(file) {
		return file, nil
	}
	f := filepath.Join(d.baseDirectory, file)
	if d.exist(f) {
		return f, nil
	}
	for _, path := range d.paths {
		f := filepath.Join(d.baseDirectory, path, file)
		if d.exist(f) {
			return f, nil
		}
	}
	if d.parent != nil {
		return d.parent.resolve(file)
	}
	return "", fmt.Errorf("file not found: file=%s, %v", file, d)
}

func (d *configDirectory) exist(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

func (d *configDirectory) addPath(path string) {
	d.paths = append(d.paths, path)
}

func (d *configDirectory) child(baseDirectory string) *configDirectory {
	return &configDirectory{d, baseDirectory, make([]string, 0)}
}

func (d *configDirectory) childRelative(relativeBaseDirectory string) *configDirectory {
	base := filepath.Join(d.baseDirectory, relativeBaseDirectory)
	return &configDirectory{d, base, make([]string, 0)}
}

func (d *configDirectory) String() string {
	return fmt.Sprintf("configDirectory{baseDirectory=%s, paths=%v}", d.baseDirectory, d.paths)
}
