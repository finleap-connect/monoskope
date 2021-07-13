package acceptance

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	messages "github.com/cucumber/messages-go/v10"
	flag "github.com/spf13/pflag"
)

var opts = godog.Options{Output: colors.Colored(os.Stdout)}

// Table represents the Table argument made to a step definition
type Table = messages.PickleStepArgument_PickleTable

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func TestMain(m *testing.M) {
	flag.Parse()
	opts.Paths = flag.Args()

	status := godog.TestSuite{
		Name:                 "cluster",
		TestSuiteInitializer: ClusterInitializeTestSuite,
		ScenarioInitializer:  ClusterInitializeScenario,
		Options:              &opts,
	}.Run()

	os.Exit(status)
}

func Fail(message string, callerSkip ...int) {
	// skip := 0
	// if len(callerSkip) > 0 {
	// 	skip = callerSkip[0]
	// }

	// global.Failer.Fail(message, codelocation.New(skip+1))
	panic(message)

}
