package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAppDirs(t *testing.T) {
	results, err := GetAppDirs([]string{})
	assert.NoError(t, err, "error not expected")

	assert.Equal(t, []string{"github.com/gesquive/gack"}, results, "results do not match")
}

func TestGetArchs(t *testing.T) {
	testArchs := []string{"386", "amd64", "arm", "x86"}
	results, err := GetArchs(testArchs)
	assert.NoError(t, err, "unexpected error")

	assert.Equal(t, testArchs, results, "arch results do not match")
}

func TestGetDefaultArchs(t *testing.T) {
	testArchs := []string{}
	results, err := GetArchs(testArchs)
	assert.NoError(t, err, "unexpected error")

	assert.Equal(t, ArchList, results, "arch results do not match")
}

func TestGetArchsSpace(t *testing.T) {
	testArchs := []string{"386 amd64", "arm,arm64"}
	results, err := GetArchs(testArchs)
	assert.NoError(t, err, "unexpected error")

	expected := []string{"386", "amd64", "arm", "arm64"}
	assert.Equal(t, expected, results, "arch results do not match")
}

func TestGetOSs(t *testing.T) {
	testOSs := []string{"darwin", "linux", "windows", "rasbian"}
	results, err := GetOSs(testOSs)
	assert.NoError(t, err, "unexpected error")

	assert.Equal(t, testOSs, results, "os results do not match")
}

func TestGetDefaultOSs(t *testing.T) {
	testOSs := []string{}
	results, err := GetOSs(testOSs)
	assert.NoError(t, err, "unexpected error")

	assert.Equal(t, OSList, results, "os results do not match")
}

func TestGetOSsSpace(t *testing.T) {
	testOSs := []string{"darwin linux", "windows,rasbian"}
	results, err := GetOSs(testOSs)
	assert.NoError(t, err, "unexpected error")

	expected := []string{"darwin", "linux", "windows", "rasbian"}
	assert.Equal(t, expected, results, "os results do not match")
}

func TestGetArchives(t *testing.T) {
	testArchives := []string{"zip", "tgz", "tar.xz", "rar"}
	results, err := GetArchiveTypes(testArchives)
	assert.NoError(t, err, "unexpected error")

	expected := []string{"zip", "tgz", "tar.xz"}
	assert.Equal(t, expected, results, "archive results do not match")
}

func TestGetDefaultArchives(t *testing.T) {
	testArchives := []string{}
	results, err := GetArchiveTypes(testArchives)
	assert.NoError(t, err, "unexpected error")

	assert.Equal(t, DefaultArchiveList, results, "archive results do not match")
}

func TestAssemblePackageInfo_DefaultList(t *testing.T) {
	results, err := AssemblePackageInfo([]string{}, []string{}, []string{}, []string{})
	assert.NoError(t, err, "unexpected error")

	assert.Equal(t, 189, len(results), "package results do not match")
}

func TestAssemblePackageInfo_SingleAssembled(t *testing.T) {
	pkg := Package{Arch: "amd64", OS: "linux", Archive: "tar.xz"}
	results, err := AssemblePackageInfo([]string{pkg.Arch}, []string{pkg.OS},
		[]string{pkg.Archive}, []string{})
	assert.NoError(t, err, "unexpected error")

	expected := []Package{pkg}
	assert.Equal(t, expected, results, "package results do not match")
}

func TestAssemblePackageInfo_DefineDuplicate(t *testing.T) {
	pkg := Package{Arch: "amd64", OS: "linux", Archive: "tar.xz"}
	results, err := AssemblePackageInfo([]string{pkg.Arch}, []string{pkg.OS},
		[]string{pkg.Archive}, []string{"linux/amd64/tar.xz"})
	assert.NoError(t, err, "unexpected error")
	assert.Len(t, results, 1, "unexpected number of results")

	expected := []Package{pkg}
	assert.Equal(t, expected, results, "package results do not match")
}

func TestAssemblePackageInfo_DefineAndAssemble(t *testing.T) {
	pkg := Package{Arch: "amd64", OS: "linux", Archive: "tar.xz"}
	pkg2 := Package{Arch: "arm", OS: "linux", Archive: "tar.gz"}
	results, err := AssemblePackageInfo([]string{pkg.Arch}, []string{pkg.OS},
		[]string{pkg.Archive}, []string{pkg2.String()})
	assert.NoError(t, err, "unexpected error")
	assert.Len(t, results, 2, "unexpected number of results")

	assert.Contains(t, results, pkg, "package missing from results")
	assert.Contains(t, results, pkg2, "package missing from results")
}

func TestGetPackagePaths(t *testing.T) {
	pkgs := []Package{Package{OS: "linux", Arch: "x64", Archive: "tgz"}}
	dirs := []string{"/this/is/a/test", "/another/test/exe"}
	inputTemplate := "test/{{.Dir}}-{{.OS}}-{{.Arch}}"
	outputTemplate := "test/{{.Dir}}-{{.OS}}-{{.Arch}}.{{.Archive}}"

	result, err := GetPackagePaths(pkgs, dirs, inputTemplate, outputTemplate)
	assert.NoError(t, err, "unexpected error")

	assert.Len(t, result, 2, "incorrect number of packaged results")

	expected := pkgs[0]
	expected.Dir = "test"
	expected.ExePath = "test/test-linux-x64"
	expected.ArchivePath = "test/test-linux-x64.tgz"
	assert.Equal(t, expected, result[0], "package results do not match")

	expected = pkgs[0]
	expected.Dir = "exe"
	expected.ExePath = "test/exe-linux-x64"
	expected.ArchivePath = "test/exe-linux-x64.tgz"
	assert.Equal(t, expected, result[1], "package results do not match")
}

func TestGetPackageFiles(t *testing.T) {
	pkgs := []Package{Package{OS: "linux", Arch: "x64", Archive: "tgz", ExePath: "bin/exe-linux-x64"}}
	fileList := []string{"readme.md", "license", "test/file"}

	result, err := GetPackageFiles(pkgs, fileList)
	assert.NoError(t, err, "unexpected error")
	assert.Len(t, result, 1, "incorrect number of packaged results")

	expected := pkgs[0]
	expected.FileList = []string{"bin/exe-linux-x64", "readme.md", "license", "test/file"}

	assert.Equal(t, expected, result[0], "package results do not match")
}
