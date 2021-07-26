package baseline

import (
	"context"
	"github.com/metrumresearchgroup/command"
	. "github.com/metrumresearchgroup/pkgr/testhelper"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

const (
	baselinePlanTest1="BSLNPLN-E2E-001"
	baselinePlanTest2="BSLNPLN-E2E-002"
	baselinePlanTest3="BSLNPLN-E2E-003"
)

// Golden Test File Names
const (
	baselinePlanGolden="baseline-plan-golden"
)

func TestPlan(t *testing.T) {



	t.Run(MakeTestName(baselinePlanTest1, "plan indicates packages to be installed, as well as version, source repo, and whether pkg is user-defined or a dependency"), func(t *testing.T) {
		DeleteTestLibrary("test-library", t)
		ctx := context.TODO()
		planCmd := command.New()

		capture, err := planCmd.Run(ctx, "pkgr", "plan", "--loglevel=debug", "--logjson")
		if err != nil {
			t.Fatalf("error running pkgr plan: %s", err)
		}
		//output := string(capture.Output)


		pkgRepoSetLogs := CollectPkgRepoSetLogs(t, capture)
		//assert.True(t, pkgRepoSetLogs.Contains("R6", "2.5.0", "CRAN", "user_defined"), "failed to find expected log message")
		//assert.True(t, pkgRepoSetLogs.Contains("pillar", "1.6.1", "CRAN", "user_defined"), "failed to find expected log message")
		//assert.True(t, pkgRepoSetLogs.Contains("glue", "1.4.2", "CRAN", "dependency"), "failed to find expected log message")
		//assert.True(t, pkgRepoSetLogs.Contains("fansi", "0.5.0", "CRAN", "dependency"), "failed to find expected log message")
		//assert.True(t, pkgRepoSetLogs.Contains("rlang", "0.4.11", "CRAN", "dependency"), "failed to find expected log message")
		//assert.True(t, pkgRepoSetLogs.Contains("utf8", "1.2.1", "CRAN", "dependency"), "failed to find expected log message")
		//assert.True(t, pkgRepoSetLogs.Contains("crayon", "1.4.1", "CRAN", "dependency"), "failed to find expected log message")
		//assert.True(t, pkgRepoSetLogs.Contains("lifecycle", "1.0.0", "CRAN", "dependency"), "failed to find expected log message")
		//assert.True(t, pkgRepoSetLogs.Contains("vctrs", "0.3.8", "CRAN", "dependency"), "failed to find expected log message")
		//assert.True(t, pkgRepoSetLogs.Contains("ellipsis", "0.3.2", "CRAN", "dependency"), "failed to find expected log message")
		//assert.True(t, pkgRepoSetLogs.Contains("cli", "2.5.0", "CRAN", "dependency"), "failed to find expected log message")

		//alternative
		g := goldie.New(t)

		g.Assert(t, baselinePlanGolden, (pkgRepoSetLogs.ToBytes()))
	})

	t.Run(MakeTestName(baselinePlanTest2, "number of workers (threads) can be set"), func(t *testing.T) {
		DeleteTestLibrary("test-library", t)
		ctx := context.TODO()
		planCmd := command.New()

		outputBytes, err := planCmd.Run(ctx, "pkgr", "plan", "--threads=5", "--logjson")
		if err != nil {
			t.Fatalf("error running pkgr plan: %s", err)
		}
		output := string(outputBytes.Output)

		jsonRegex := `\{"level":"info","msg":"Installation would launch 5 workers.*\}`
		assert.Regexp(t, jsonRegex, output)
	})

	// We can't determine the specific global cache location up front because it varies by OS and might be a tempdir.
	// So this test is only going to verify that it is printed out and used.
	t.Run(MakeTestName(baselinePlanTest3, "pkgr will indicate the location of the global cache and use it for installation"), func(t *testing.T) {
		DeleteTestLibrary("test-library", t)
		ctx := context.TODO()
		installCmd := command.New()

		outputBytes, err := installCmd.Run(ctx, "pkgr", "plan", "--config=pkgr-global-cache.yml", "--logjson")
		if err != nil {
			t.Fatalf("error running pkgr plan: %s", err)
		}
		output := string(outputBytes.Output)

		pkgDbRegex := `\{"level":"info","msg":"Database cache directory: .*\}` // We can't determine this up front because it varies by OS and might be a temp dir.
		pkgCacheRegex:= `\{"level":"info","msg":"Package installation cache directory: .*\}` // We can't determine this up front because it varies by OS and might be a temp dir.
		assert.Regexp(t, pkgDbRegex, output)
		assert.Regexp(t, pkgCacheRegex, output)

		// If that passed, we can assume this will work:
		rPkgDbs := regexp.MustCompile(`\{"level":"info","msg":"Database cache directory:  (.*)"\}`)
		pkgDbDir := rPkgDbs.FindStringSubmatch(output)[1]
		rPkgCache := regexp.MustCompile(`\{"level":"info","msg":"Package installation cache directory:  (.*)"\}`)
		pkgCacheDir := rPkgCache.FindStringSubmatch(output)[1]

		//t.Log(pkgDbDir)
		//t.Log(pkgCacheDir)
		pkgDbDirContents, err := os.ReadDir(filepath.Clean(pkgDbDir))
		if err != nil {
			t.Fatalf("error attempting to read global pkgDb dir: %s", err)
		}
		assert.NotEmpty(t, pkgDbDirContents)


		pkgCacheDirContents, err := os.ReadDir(pkgCacheDir)
		if err != nil {
			t.Fatalf("error attempting to read global package cache dir: %s", err)
		}

		assert.NotEmpty(t, pkgCacheDirContents)
	})


}
