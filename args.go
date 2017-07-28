package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gesquive/cli"
	"github.com/pkg/errors"
)

// Package is a combination of OS/arch/archive that can be packaged.
type Package struct {
	OS          string
	Arch        string
	Archive     string
	ExePath     string
	ArchivePath string
	FileList    []string
	Dir         string
}

func (p *Package) String() string {
	return fmt.Sprintf("%s/%s/%s", p.OS, p.Arch, p.Archive)
}

func ParsePackage(pkgString string) (Package, error) {
	pkg := Package{}
	parts := strings.SplitN(pkgString, "/", 3)
	if len(parts) != 3 {
		return pkg, errors.Errorf("could not parse package '%s'", pkgString)
	}
	pkg.OS = parts[0]
	pkg.Arch = parts[1]
	pkg.Archive = parts[2]
	return pkg, nil
}

var (
	// OSList is the full list of golang OSs
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

	// ArchList is the full list of golang architectures
	ArchList = []string{
		"386",
		"amd64",
		"amd64p32",
		"arm",
		"arm64",
		"ppc64",
		"ppc64le",
	}

	// DefaultArchiveList is the list of default archives
	DefaultArchiveList = []string{
		"zip",
		"tar.gz",
		"tar.xz",
	}

	// ArchiveList is the full list of supported archives
	ArchiveList = []string{
		"zip",
		"tar",
		"tar.gz",
		"tar.bz2",
		"tar.xz",
		"tar.lz4",
		"tar.sz",
	}
)

// GetUserArchs generates a list of architectures from the user defined list
func GetUserArchs(userArch []string) ([]string, error) {
	cleanList := splitListItems(userArch)
	pList, nList := splitNegatedItems(cleanList)
	if len(pList) == 0 {
		pList = ArchList
	}
	cleanList = negateList(pList, nList)
	return cleanList, nil
}

// GetUserOSs generates a list of OSs from the user defined list
func GetUserOSs(userOS []string) ([]string, error) {
	cleanList := splitListItems(userOS)
	pList, nList := splitNegatedItems(cleanList)
	if len(pList) == 0 {
		pList = OSList
	}
	cleanList = negateList(pList, nList)
	return cleanList, nil
}

// GetUserArchives generates a list of valid archive types from the user defined list
func GetUserArchives(userArchive []string) ([]string, error) {
	cleanList := splitListItems(userArchive)
	pList, nList := splitNegatedItems(cleanList)
	if len(pList) == 0 {
		pList = ArchiveList
	}
	cleanList = negateList(pList, nList)

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

func GetUserPackages(userPkgs []string) ([]Package, error) {
	pkgs := []Package{}
	userPkgs = splitListItems(userPkgs)
	for _, userPkg := range userPkgs {
		pkg, err := ParsePackage(userPkg)
		if err != nil {
			continue
		}
		pkgs = append(pkgs, pkg)
	}
	return pkgs, nil
}

// AssemblePackageInfo generates a list of packages from the user defined arguments
func AssemblePackageInfo(userArch []string, userOS []string,
	userArchive []string, userPackages []string) ([]Package, error) {
	archList, _ := GetUserArchs(userArch)
	osList, _ := GetUserOSs(userOS)
	archiveList, _ := GetUserArchives(userArchive)
	specificList, _ := GetUserPackages(userPackages)

	packageList := []Package{}
	for _, arch := range archList {
		for _, os := range osList {
			for _, archive := range archiveList {
				pkg := Package{Arch: arch, OS: os, Archive: archive}
				packageList = appendIfMissing(packageList, pkg)
			}
		}
	}

	for _, userPkg := range specificList {
		if strings.HasPrefix(userPkg.String(), "!") {
			userPkg.OS = userPkg.OS[1:]
			packageList = removeIfPresent(packageList, userPkg)
		} else {
			packageList = appendIfMissing(packageList, userPkg)
		}
	}

	return packageList, nil
}

// GetPackagePaths generates info about the archives
func GetPackagePaths(packages []Package, dirs []string, inputTemplate string,
	outputTemplate string) ([]Package, error) {
	filledPackages := []Package{}
	for _, pkg := range packages {
		for _, path := range dirs {
			filledPkg := Package{
				Dir:     filepath.Base(path),
				OS:      pkg.OS,
				Arch:    pkg.Arch,
				Archive: pkg.Archive,
			}

			inputTpl, err := template.New("input").Parse(inputTemplate)
			if err != nil {
				return nil, errors.Wrap(err, "input template error")
			}
			var inputPath bytes.Buffer
			if err = inputTpl.Execute(&inputPath, &filledPkg); err != nil {
				return nil, errors.Wrap(err, "error generating input path")
			}
			filledPkg.ExePath = inputPath.String()

			outputTpl, err := template.New("output").Parse(outputTemplate)
			if err != nil {
				return nil, errors.Wrap(err, "output template error")
			}

			var outputPath bytes.Buffer
			if err := outputTpl.Execute(&outputPath, &filledPkg); err != nil {
				return nil, errors.Wrap(err, "error generating output path")
			}
			filledPkg.ArchivePath = outputPath.String()

			filledPackages = append(filledPackages, filledPkg)
		}
	}
	return filledPackages, nil
}

func GetPackageFiles(packages []Package, fileList []string) ([]Package, error) {
	pkgs := []Package{}
	fileList = splitListItems(fileList)
	for _, pkg := range packages {
		files := append([]string{pkg.ExePath}, fileList...)
		pkg.FileList = files
		pkgs = append(pkgs, pkg)
	}
	return pkgs, nil
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

func splitNegatedItems(list []string) (p []string, n []string) {
	for _, item := range list {
		if strings.HasPrefix(item, "!") {
			n = append(n, item[1:])
		} else {
			p = append(p, item)
		}
	}
	return
}

func negateList(pList []string, nList []string) []string {
	finalList := []string{}

	for _, item := range pList {
		lowerItem := strings.ToLower(item)
		found := false
		for _, negation := range nList {
			if lowerItem == strings.ToLower(negation) {
				found = true
				break
			}
		}
		if !found {
			finalList = append(finalList, item)
		}
	}

	return finalList
}

func removeIfPresent(pkgs []Package, pkg Package) []Package {
	match := strings.ToLower(pkg.String())
	result := []Package{}
	for _, existing := range pkgs {
		if strings.ToLower(existing.String()) != match {
			result = append(result, existing)
		}
	}
	return result
}

func appendIfMissing(pkgs []Package, pkg Package) []Package {
	match := strings.ToLower(pkg.String())
	missing := true
	for _, existing := range pkgs {
		if strings.ToLower(existing.String()) == match {
			missing = false
		}
	}

	if missing {
		pkgs = append(pkgs, pkg)
	}
	return pkgs
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

// IsEmpty checks to see if a directory is empty
// src: https://stackoverflow.com/a/30708914/613218
func IsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}
