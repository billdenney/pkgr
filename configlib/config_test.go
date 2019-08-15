package configlib

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/metrumresearchgroup/pkgr/cran"
	"github.com/metrumresearchgroup/pkgr/rcmd"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestAddRemovePackage(t *testing.T) {
	tests := []struct {
		fileName    string
		packageName string
	}{
		{
			fileName:    "../integration_tests/simple/pkgr.yml",
			packageName: "packageTestName",
		},
		{
			fileName:    "../integration_tests/simple-suggests/pkgr.yml",
			packageName: "packageTestName",
		},
		{
			fileName:    "../integration_tests/mixed-source/pkgr.yml",
			packageName: "packageTestName",
		},
		{
			fileName:    "../integration_tests/outdated-pkgs/pkgr.yml",
			packageName: "packageTestName",
		},
		{
			fileName:    "../integration_tests/outdated-pkgs-no-update/pkgr.yml",
			packageName: "packageTestName",
		},
		{
			fileName:    "../integration_tests/repo-order/pkgr.yml",
			packageName: "packageTestName",
		},
	}

	appFS := afero.NewOsFs()
	for _, tt := range tests {
		b, _ := afero.Exists(appFS, tt.fileName)
		assert.Equal(t, true, b, fmt.Sprintf("yml file not found:%s", tt.fileName))

		ymlStart, _ := afero.ReadFile(appFS, tt.fileName)

		add(tt.fileName, tt.packageName)
		b, _ = afero.FileContainsBytes(appFS, tt.fileName, []byte(tt.packageName))
		assert.Equal(t, true, b, fmt.Sprintf("Package not added:%s", tt.fileName))

		remove(tt.fileName, tt.packageName)
		b, _ = afero.FileContainsBytes(appFS, tt.fileName, []byte(tt.packageName))
		assert.Equal(t, false, b, fmt.Sprintf("Package not removed:%s", tt.fileName))

		ymlEnd, _ := afero.ReadFile(appFS, tt.fileName)
		b = equal(ymlStart, ymlEnd, false)
		assert.Equal(t, true, b, fmt.Sprintf("Start and End yml files differ:%s", tt.fileName))

		// put file back for Git
		fi, _ := os.Stat(tt.fileName)
		err := afero.WriteFile(appFS, tt.fileName, ymlStart, fi.Mode())
		assert.Equal(t, nil, err, fmt.Sprintf("Error writing file back to original state:%s", tt.fileName))
	}
}

func TestRemoveWhitespace(t *testing.T) {

	tests := []struct {
		in       string
		expected string
		message  string
	}{
		{
			in:       "hello world\n",
			expected: "helloworld",
			message:  "newline",
		},
		{
			in:       "hello world\t",
			expected: "helloworld",
			message:  "h tab",
		},
		{
			in:       "hello world\v",
			expected: "helloworld",
			message:  "v tab",
		},
		{
			in:       "hello world\f",
			expected: "helloworld",
			message:  "feed",
		},
		{
			in:       "hello world\r",
			expected: "helloworld",
			message:  "return",
		},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.expected, string(removeWhitespace([]byte(tt.in))), fmt.Sprintf("fail to remove:%s", tt.message))
	}
}

func removeWhitespace(b []byte) []byte {
	var ws = []byte{'\t', '\n', '\v', '\f', '\r', ' '}
	for _, r := range ws {
		b = bytes.ReplaceAll(b, []byte(string(r)), []byte(""))
	}
	return b
}

func equal(a []byte, b []byte, compareWs bool) bool {
	if compareWs == false {
		a = removeWhitespace(a)
		b = removeWhitespace(b)
	}
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestNewConfigPackrat(t *testing.T) {
	tests := []struct {
		folder   string
		expected string
		message  string
	}{
		{
			folder:   "../integration_tests/packrat-library",
			expected: "packrat",
			message:  "packrat exists",
		},
	}
	for _, tt := range tests {
		var cfg PkgrConfig
		_ = os.Chdir(tt.folder)
		_ = LoadConfigFromPath(viper.GetString("config"))
		NewConfig(&cfg)
		assert.Equal(t, tt.expected, cfg.Lockfile.Type, fmt.Sprintf("Fail:%s", tt.message))
	}
}

func TestNewConfigNoPackrat(t *testing.T) {
	tests := []struct {
		folder   string
		expected string
		message  string
	}{
		{
			folder:   "../integration_tests/simple",
			expected: "",
			message:  "packrat does not exist",
		},
	}
	for _, tt := range tests {
		var cfg PkgrConfig
		_ = os.Chdir(tt.folder)
		_ = LoadConfigFromPath(viper.GetString("config"))
		NewConfig(&cfg)
		assert.Equal(t, tt.expected, cfg.Lockfile.Type, fmt.Sprintf("Fail:%s", tt.message))
	}
}

func TestGetLibraryPath(t *testing.T) {
	tests := []struct {
		lftype   string
		expected string
		message  string
	}{
		{
			lftype:   "renv",
			expected: "renv/library/R-1.2/apple",
		},
		{
			lftype:   "packrat",
			expected: "packrat/lib/apple/1.2.3",
		},
		{
			lftype:   "pkgr",
			expected: "original",
		},
	}
	for _, tt := range tests {
		var rv = cran.RVersion{
			Major: 1,
			Minor: 2,
			Patch: 3,
		}
		library := getLibraryPath(tt.lftype, "myRpath", rv, "apple", "original")
		assert.Equal(t, tt.expected, library, fmt.Sprintf("Fail:%s", tt.expected))
	}
}

func TestSetCustomizations(t *testing.T) {
	tests := []struct {
		pkg   string
		name  string
		value string
	}{
		{
			pkg:   "data.table",
			name:  "R_MAKEVARS_USER",
			value: "~/.R/Makevars_data.table",
		},
		{
			pkg:   "boo",
			name:  "foo",
			value: "soo",
		},
	}
	for _, tt := range tests {
		var cfg PkgrConfig
		NewConfig(&cfg)
		cfg.Customizations.Packages = map[string]PkgConfig{
			tt.pkg: PkgConfig{
				Env: map[string]string{
					tt.name: tt.value,
				},
			},
		}
		var rs rcmd.RSettings
		rs.PkgEnvVars = make(map[string]map[string]string)
		rs2 := SetCustomizations(cfg, rs)
		assert.Equal(t, tt.value, rs2.PkgEnvVars[tt.pkg][tt.name], fmt.Sprintf("Fail to get: %s", tt.value))
	}
}
