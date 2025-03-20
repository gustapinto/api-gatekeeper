package yamlutil

import (
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	envExpr        = regexp.MustCompile(`\$\{([^}]+)\}`)
	envKeyReplacer = strings.NewReplacer("${", "", "}", "")
)

// Unmarshal Wraps gopkg.in/yaml.v3.Unmarshal and adds environment variable substitution
// using the ${VARIABLE_NAME} syntax
func Unmarshal(data []byte, target any) error {
	parsedData := envExpr.ReplaceAllFunc(data, func(b []byte) []byte {
		key := envKeyReplacer.Replace(string(b))

		return []byte(os.Getenv(key))
	})

	return yaml.Unmarshal(parsedData, target)
}
