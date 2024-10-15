// Package testutilsext is purely for test utilities that may access other packages
// in the codebase that tend to use testutils.
package testutilsext

import (
	"fmt"
	"os"

	"github.com/edaniels/golog"
	"go.uber.org/goleak"

	"go.viam.com/utils"
	"go.viam.com/utils/artifact"
	"go.viam.com/utils/testutils"
)

// VerifyTestMain preforms various runtime checks on code that tests run.
func VerifyTestMain(m goleak.TestingM) {
	cache, err := artifact.GlobalCache()
	if err != nil {
		golog.Global().Fatalw("error opening artifact", "error", err)
	}
	currentGoroutines := goleak.IgnoreCurrent()
	//nolint:ifshort
	exitCode := m.Run()
	testutils.Teardown()
	if err := cache.Close(); err != nil {
		golog.Global().Errorw("error closing artifact", "error", err)
	}
	if exitCode != 0 {
		os.Exit(exitCode)
	}
	if err := utils.FindGoroutineLeaks(currentGoroutines); err != nil {
		fmt.Fprintf(os.Stderr, "goleak: Errors on successful test run: %v\n", err)
		os.Exit(1)
	}
}
