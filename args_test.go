package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAppDirs(t *testing.T) {
	results, err := GetAppDirs([]string{})
	assert.NoError(t, err, "error not expected")

	assert.Equal(t, []string{"github.com/gesquive/gop"}, results, "results do not match")
}

func TestGetUserArchs(t *testing.T) {
	testArchs := []string{"386", "amd64", "arm", "x86"}
	results, err := GetUserArchs(testArchs)
	assert.NoError(t, err, "unexpected error")

	assert.Equal(t, testArchs, results, "arch results do not match")
}

func TestGetUserArchs_Default(t *testing.T) {
	testArchs := []string{}
	results, err := GetUserArchs(testArchs)
	assert.NoError(t, err, "unexpected error")

	assert.Equal(t, ArchList, results, "arch results do not match")
}

func TestGetUserArchs_WithDelimiters(t *testing.T) {
	testArchs := []string{"386 amd64", "arm,arm64"}
	results, err := GetUserArchs(testArchs)
	assert.NoError(t, err, "unexpected error")

	expected := []string{"386", "amd64", "arm", "arm64"}
	assert.Equal(t, expected, results, "arch results do not match")
}

func TestGetUserArchs_DefaultNegations(t *testing.T) {
	testArchs := []string{"!amd64p32", "!ppc64le"}
	results, err := GetUserArchs(testArchs)
	assert.NoError(t, err, "unexpected error")

	expected := []string{"386", "amd64", "arm", "arm64", "ppc64"}
	assert.Equal(t, expected, results, "arch results do not match")
}

func TestGetUserArchs_WithNegations(t *testing.T) {
	testArchs := []string{"386 amd64", "!arm,!arm64"}
	results, err := GetUserArchs(testArchs)
	assert.NoError(t, err, "unexpected error")

	expected := []string{"386", "amd64"}
	assert.Equal(t, expected, results, "arch results do not match")
}

func TestGetUserOSs(t *testing.T) {
	testOSs := []string{"darwin", "linux", "windows", "rasbian"}
	results, err := GetUserOSs(testOSs)
	assert.NoError(t, err, "unexpected error")

	assert.Equal(t, testOSs, results, "os results do not match")
}

func TestGetUserOSs_Default(t *testing.T) {
	testOSs := []string{}
	results, err := GetUserOSs(testOSs)
	assert.NoError(t, err, "unexpected error")

	assert.Equal(t, OSList, results, "os results do not match")
}

func TestGetUserOSs_WithDelimiters(t *testing.T) {
	testOSs := []string{"darwin linux", "windows,rasbian"}
	results, err := GetUserOSs(testOSs)
	assert.NoError(t, err, "unexpected error")

	expected := []string{"darwin", "linux", "windows", "rasbian"}
	assert.Equal(t, expected, results, "os results do not match")
}

func TestGetUserOSs_DefaultNegations(t *testing.T) {
	testOSs := []string{"!dragonfly", "!netbsd", "!openbsd", "!plan9", "!solaris"}
	results, err := GetUserOSs(testOSs)
	assert.NoError(t, err, "unexpected error")

	expected := []string{"darwin", "freebsd", "linux", "windows"}
	assert.Equal(t, expected, results, "os results do not match")
}

func TestGetUserOSs_WithNegations(t *testing.T) {
	testOSs := []string{"darwin", "linux", "!windows"}
	results, err := GetUserOSs(testOSs)
	assert.NoError(t, err, "unexpected error")

	expected := []string{"darwin", "linux"}
	assert.Equal(t, expected, results, "os results do not match")
}

func TestGetUserArchives(t *testing.T) {
	testArchives := []string{"zip", "tar.gz", "tar.xz", "rar"}
	results, err := GetUserArchives(testArchives)
	assert.NoError(t, err, "unexpected error")

	expected := []string{"zip", "tar.gz", "tar.xz"}
	assert.Equal(t, expected, results, "archive results do not match")
}

func TestGetUserArchives_Default(t *testing.T) {
	testArchives := []string{}
	results, err := GetUserArchives(testArchives)
	assert.NoError(t, err, "unexpected error")

	assert.Equal(t, ArchiveList, results, "archive results do not match")
}

func TestGetUserArchives_DefaultNegations(t *testing.T) {
	testArchives := []string{"!zip", "!tar", "!tar.bz2", "!rar"}
	results, err := GetUserArchives(testArchives)
	assert.NoError(t, err, "unexpected error")

	expected := []string{"tar.gz", "tar.xz", "tar.lz4", "tar.sz"}
	assert.Equal(t, expected, results, "archive results do not match")
}

func TestGetUserArchives_WithNegations(t *testing.T) {
	testArchives := []string{"zip", "tar.gz", "tar.xz", "rar", "!zip", "!tar.bz2"}
	results, err := GetUserArchives(testArchives)
	assert.NoError(t, err, "unexpected error")

	expected := []string{"tar.gz", "tar.xz"}
	assert.Equal(t, expected, results, "archive results do not match")
}

func TestGetUserPackages(t *testing.T) {
	pkg := Package{Arch: "amd64", OS: "linux", Archive: "tar.xz"}
	pkg2 := Package{Arch: "x86", OS: "linux", Archive: "tar.gz"}
	results, err := GetUserPackages([]string{pkg.String(), pkg2.String()})
	assert.NoError(t, err, "unexpected error")

	assert.Len(t, results, 2, "incorrect number of packages")
	assert.Contains(t, results, pkg, "missing expected package")
	assert.Contains(t, results, pkg2, "missing expected package")
}

func TestGetUserPackages_Default(t *testing.T) {
	results, err := GetUserPackages([]string{})
	assert.NoError(t, err, "unexpected error")

	assert.Len(t, results, 0, "incorrect number of packages")
}

func TestGetUserPackages_WithDelimiters(t *testing.T) {
	pkg := Package{Arch: "amd64", OS: "linux", Archive: "tar.xz"}
	pkg2 := Package{Arch: "x86", OS: "linux", Archive: "tar.gz"}
	results, err := GetUserPackages([]string{"linux/amd64/tar.xz linux/x86/tar.gz"})
	assert.NoError(t, err, "unexpected error")

	assert.Len(t, results, 2, "incorrect number of packages")
	assert.Contains(t, results, pkg, "missing expected package")
	assert.Contains(t, results, pkg2, "missing expected package")
}

func TestGetUserPackages_InvalidPackage(t *testing.T) {
	results, err := GetUserPackages([]string{"linux/amd64"})
	assert.NoError(t, err, "unexpected error")

	assert.Len(t, results, 0, "incorrect number of packages")
}

func TestAssemblePackageInfo_DefaultList(t *testing.T) {
	results, err := AssemblePackageInfo([]string{}, []string{}, []string{}, []string{})
	assert.NoError(t, err, "unexpected error")

	assert.Equal(t, 441, len(results), "package results do not match")
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

func TestAssemblPackageInfo_NegateArch(t *testing.T) {
	results, err := AssemblePackageInfo([]string{"!arm", "amd64", "arm", "x86"},
		[]string{"linux"}, []string{"tar.gz", "tar.xz"}, []string{})
	assert.NoError(t, err, "unexpected error")
	assert.Len(t, results, 4, "unexpected number of results")

}

func TestAssemblPackageInfo_NegatePackage(t *testing.T) {
	results, err := AssemblePackageInfo([]string{"arm", "amd64", "x86"},
		[]string{"linux"}, []string{"tar.gz", "tar.xz"}, []string{"!linux/arm/tar.xz"})
	assert.NoError(t, err, "unexpected error")
	assert.Len(t, results, 5, "unexpected number of results")
	assert.NotContains(t, results, Package{Arch: "arm", OS: "linux", Archive: "tar.xz"},
		"negated package found in results")
}

func TestAssemblePackageInfo_OnlyNegate(t *testing.T) {
	results, err := AssemblePackageInfo([]string{}, []string{}, []string{},
		[]string{"!linux/arm/tar.xz", "!darwin/arm/tar.gz"})
	assert.NoError(t, err, "unexpected error")
	assert.Len(t, results, 439, "unexpected number of results")
	assert.NotContains(t, results, Package{Arch: "arm", OS: "linux", Archive: "tar.xz"},
		"negated package found in results")
	assert.NotContains(t, results, Package{Arch: "arm", OS: "darwin", Archive: "tar.gz"},
		"negated package found in results")
}

func TestAssemblePackageInfo_Precedence(t *testing.T) {
	// if included in packages, it should be built even if negated in user flags
	results, err := AssemblePackageInfo([]string{"arm", "amd64", "!x86"},
		[]string{"linux"}, []string{"tar.gz", "tar.xz"}, []string{"linux/x86/tar.gz"})
	assert.NoError(t, err, "unexpected error")
	assert.Len(t, results, 5, "unexpected number of results")
	assert.Contains(t, results, Package{Arch: "x86", OS: "linux", Archive: "tar.gz"},
		"negated package found in results")
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
