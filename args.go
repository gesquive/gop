package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/gesquive/cli"
)

// Package is a combination of OS/arch/archive that can be packaged.
type Package struct {
	OS      string
	Arch    string
	Archive string
}

// func (p *Package) String() string {
// 	return fmt.Sprintf("%s-%s-%s", p.OS, p.Arch, p.Archive)
// }

var (
	OSList = []string{
		"darwin",
		"dragonfly",
		"freebsd",
		"linux",
		"netbsd",
		"openbsd",
		"plan9",
		"solaris",
		"windows",
	}

	ArchList = []string{
		"386",
		"amd64",
		"amd64p32",
		"arm",
		"arm64",
		"ppc64",
		"ppc64le",
	}

	DefaultArchiveList = []string{
		"zip",
		"tar.gz",
		"tar.xz",
	}

	ArchiveList = []string{
		"zip",
		"tar",
		"tgz",
		"tar.gz",
		"tbz2",
		"tar.bz2",
		"txz",
		"tar.xz",
		"tlz4",
		"tar.lz4",
		"tsz",
		"tar.sz",
	}
)

func GetPackages(userArch []string, userOS []string, userArchive []string) ([]Package, error) {
	archList, _ := GetArchs(userArch)
	osList, _ := GetOSs(userOS)
	archiveList, _ := GetArchiveTypes(userArchive)

	packageList := []Package{}
	for _, arch := range archList {
		for _, os := range osList {
			for _, archive := range archiveList {
				packageList = append(packageList, Package{Arch: arch, OS: os, Archive: archive})
			}
		}
	}

	return packageList, nil
}

func GetArchs(userArch []string) ([]string, error) {
	if len(userArch) == 0 {
		return ArchList, nil
	}
	if len(userArch) == 1 && strings.ToLower(userArch[0]) == "all" {
		return ArchList, nil
	}

	cleanList := splitListItems(userArch)
	return cleanList, nil
}

func GetOSs(userOS []string) ([]string, error) {
	if len(userOS) == 0 {
		return OSList, nil
	}
	if len(userOS) == 1 && strings.ToLower(userOS[0]) == "all" {
		return OSList, nil
	}

	cleanList := splitListItems(userOS)
	return cleanList, nil
}

func GetArchiveTypes(userArchive []string) ([]string, error) {
	if len(userArchive) == 0 || len(userArchive) == 1 && strings.ToLower(userArchive[0]) == "default" {
		return DefaultArchiveList, nil
	}
	if len(userArchive) == 1 && strings.ToLower(userArchive[0]) == "all" {
		return ArchiveList, nil
	}

	cleanList := splitListItems(userArchive)

	validArchives := []string{}
	for _, archive := range cleanList {
		for _, dArchive := range ArchiveList {
			if strings.ToLower(archive) == dArchive {
				validArchives = append(validArchives, archive)
				break
			}
		}
	}

	return validArchives, nil
}

func splitListItems(list []string) []string {
	cleanList := []string{}
	for _, item := range list {
		if parts := strings.Split(item, " "); len(parts) > 1 {
			cleanList = append(cleanList, parts...)
		} else if parts := strings.Split(item, ","); len(parts) > 1 {
			cleanList = append(cleanList, parts...)
		} else {
			cleanList = append(cleanList, item)
		}
	}
	return cleanList
}

// GetAppDirs returns the file paths to the packages that are "main"
// packages, from the list of packages given. The list of packages can
// include relative paths, the special "..." Go keyword, etc.
func GetAppDirs(packages []string) ([]string, error) {
	if len(packages) < 1 {
		packages = []string{"."}
	}

	// Get the packages that are in the given paths
	args := make([]string, 0, len(packages)+3)
	args = append(args, "list", "-f", "{{.Name}}|{{.ImportPath}}")
	args = append(args, packages...)

	output, err := execGo("go", nil, "", args...)
	if err != nil {
		return nil, err
	}

	results := make([]string, 0, len(output))
	for _, line := range strings.Split(output, "\n") {
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 {
			cli.Warn("Bad line reading packages: %s", line)
			continue
		}

		if parts[0] == "main" {
			results = append(results, parts[1])
		}
	}

	return results, nil
}

// NOTE: The original code can be found at the gox repo
//	https://raw.githubusercontent.com/mitchellh/gox/master/go.go
func execGo(GoCmd string, env []string, dir string, args ...string) (string, error) {
	var stderr, stdout bytes.Buffer
	cmd := exec.Command(GoCmd, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if env != nil {
		cmd.Env = env
	}
	if dir != "" {
		cmd.Dir = dir
	}
	if err := cmd.Run(); err != nil {
		err = fmt.Errorf("%s\nStderr: %s", err, stderr.String())
		return "", err
	}

	return stdout.String(), nil
}

const versionSource = `package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Print(runtime.Version())
}`
