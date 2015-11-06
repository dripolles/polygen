package generator

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

func getDestinationFromTypes(source string, types TypeAssignments) string {
	sourceSection := getSourceSection(source)
	typesSection := getTypesSection(types)

	return fmt.Sprintf("%s_%s.go", sourceSection, typesSection)
}

func getTypesSection(types TypeAssignments) string {
	typeStrings := sort.StringSlice{}
	for alias, typ := range types {
		str := fmt.Sprintf("%s:%s", alias, typ)
		typeStrings = append(typeStrings, str)
	}

	sort.Sort(typeStrings)
	destList := []string{}
	for _, str := range typeStrings {
		parts := strings.Split(str, ":")
		destList = append(destList, parts[1])
	}

	return strings.Join(destList, "")
}

func getSourceSection(source string) string {
	extension := filepath.Ext(source)
	basename := filepath.Base(source)
	noextension := strings.TrimSuffix(basename, extension)

	return noextension
}
