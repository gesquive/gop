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

func TestGetDefaultPackages(t *testing.T) {
	results, err := GetPackages([]string{}, []string{}, []string{})
	assert.NoError(t, err, "unexpected error")

	assert.Equal(t, 189, len(results), "package results do not match")
}

func TestGetPackages(t *testing.T) {
	pkg := Package{Arch: "amd64", OS: "linux", Archive: "tar.xz"}
	results, err := GetPackages([]string{pkg.Arch}, []string{pkg.OS},
		[]string{pkg.Archive})
	assert.NoError(t, err, "unexpected error")

	expected := []Package{pkg}
	assert.Equal(t, expected, results, "package results do not match")
}
