// A tool that takes a JSON or YAML on standard input and produces tj tree from it
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ridge/must"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

func args() formatConfig {
	var format, packageName, variableName, tjPrefix, tjPackage string

	pflag.StringVar(&format, "format", "json", "Input format to use: [json | yaml]")
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
		format:       inputFormat(format),
		packageName:  packageName,
		variableName: variableName,
		tjPrefix:     tjPrefix,
		tjPackage:    tjPackage,
	}
}

func main() {
	cfg := args()

	var val interface{}

	switch cfg.format {
	case inputJSON:
		must.OK(json.NewDecoder(os.Stdin).Decode(&val))
	case inputYAML:
		must.OK(yaml.NewDecoder(os.Stdin).Decode(&val))
		val = keysToStrings(val)
	default:
		panic(fmt.Sprintf("unsupported format %s", cfg.format))
	}

	formatHeader(cfg)
	format(cfg, val)
}

func keysToStrings(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = keysToStrings(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = keysToStrings(v)
		}
	}
	return i
}
