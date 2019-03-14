package processor

import (
	"os"
	"strings"
)

// useful for using kube service environment variables

// SubEnvVariables substitutes environment variables into config
func SubEnvVariables(strConf *string) {
	subCount := strings.Count(*strConf, "$$")
	replaceCount := 0
	if subCount > 0 {
		for _, e := range os.Environ() {
			pair := strings.Split(e, "=")
			if len(pair) == 2 {
				if strings.Contains(*strConf, "$$"+pair[0]) {
					*strConf = strings.Replace(*strConf, "$$"+pair[0], pair[1], -1)
					replaceCount++
				}
			}
			if replaceCount >= subCount {
				break
			}
		}
	}
}
