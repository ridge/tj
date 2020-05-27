// A tool that takes a JSON file and produces tj tree from it
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ridge/must"
	"github.com/spf13/pflag"
)

func args() formatConfig {
	var packageName, variableName, tjPrefix, tjPackage string

	pflag.StringVar(&packageName, "package", "", "Generate package declaration (requires --variable)")
	pflag.StringVar(&variableName, "variable", "", "Generate variable declaration")
	pflag.StringVar(&tjPrefix, "tj-prefix", "tj", "Import prefix for to use for tj package")
	pflag.StringVar(&tjPackage, "tj-package", "github.com/ridge/tj", "Import path for tj package")

	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTION]... < INPUT.json > OUTPUT.go\n", os.Args[0])
		pflag.PrintDefaults()
	}

	pflag.Parse()

	if packageName != "" && variableName == "" {
		panic("--package requires --variable")
	}

	return formatConfig{
		w:            os.Stdout,
		packageName:  packageName,
		variableName: variableName,
		tjPrefix:     tjPrefix,
		tjPackage:    tjPackage,
	}
}

func main() {
	cfg := args()

	var val interface{}
	must.OK(json.NewDecoder(os.Stdin).Decode(&val))

	formatHeader(cfg)
	format(cfg, val)
}
