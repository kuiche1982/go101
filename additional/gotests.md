gotests: github.com/cweill/gotestsIndex | Files | Directories
package gotests
import "github.com/cweill/gotests"

Package gotests contains the core logic for generating table-driven tests.

Index
type GeneratedTest
func GenerateTests(srcPath string, opt *Options) ([]*GeneratedTest, error)
type Options
Package Files
gotests.go

type GeneratedTest
type GeneratedTest struct {
    Path      string             // The test file's absolute path.
    Functions []*models.Function // The functions with new test methods.
    Output    []byte             // The contents of the test file.
}
A GeneratedTest contains information about a test file with generated tests.

func GenerateTests
func GenerateTests(srcPath string, opt *Options) ([]*GeneratedTest, error)
GenerateTests generates table-driven tests for the function and method signatures defined in the target source path file(s). The source path parameter can be either a Go source file or directory containing Go files.

type Options
type Options struct {
    Only           *regexp.Regexp         // Includes only functions that match.
    Exclude        *regexp.Regexp         // Excludes functions that match.
    Exported       bool                   // Include only exported methods
    PrintInputs    bool                   // Print function parameters in error messages
    Subtests       bool                   // Print tests using Go 1.7 subtests
    Importer       func() types.Importer  // A custom importer.
    Template       string                 // Name of custom template set
    TemplateDir    string                 // Path to custom template set
    TemplateParams map[string]interface{} // Custom external parameters
}
Options provides custom filters and parameters for generating tests.