package aptutil

import (
	"fmt"
	"strings"
)

func processLinees(linees []string) (string, map[string]string) {
	vars := map[string]string{}

	currentVarContent := ""
	for _, oneLine := range linees {
		currentVarContent += oneLine + "\n"
		if !strings.HasPrefix(oneLine, " ") {
			label := currentVarContent[0:strings.Index(currentVarContent, ": ")]
			content := currentVarContent[len(label)+2:]
			if vars[label] != "" {
				panic(fmt.Errorf("duplicated %#v label", label))
			}
			vars[label] = strings.Trim(content, "\n ")

			currentVarContent = ""
		}
	}

	packageName := vars["Package"]
	if packageName == "" {
		panic(fmt.Errorf("no package name found %#v", linees))
	}
	if vars["Version"] == "" {
		panic(fmt.Errorf("no package Version found for %#v", packageName))
	}
	return packageName, vars
}
