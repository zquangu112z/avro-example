package main

import "strings"

// Return 2 lists of header: one for general and one for specific data (section)
func getHeaderLists(mapping [][]string) ([]string, []string) {
	var generalHeaders []string
	var specificHeaders []string
	for _, tuple := range mapping {
		if strings.HasPrefix(tuple[0], ">") {
			generalHeaders = appendIfMissing(generalHeaders, tuple[1])
		} else {
			specificHeaders = appendIfMissing(specificHeaders, tuple[1])
		}
	}
	return generalHeaders, specificHeaders
}

// Return 2 lists of segments' name: one for general and one for specific data (section)
func getNameLists(mapping [][]string) ([]string, string) {
	var generalNames []string
	var specificName string
	for _, tuple := range mapping {
		if strings.HasPrefix(tuple[0], ">") {
			newItem := strings.TrimPrefix(strings.Split(tuple[0], ".")[0], ">")
			generalNames = appendIfMissing(generalNames, newItem)
		} else {
			if len(specificName) == 0 {
				newItem := strings.Split(tuple[0], ".")[0]
				specificName = newItem
			}
		}
	}
	return generalNames, specificName
}

// Return the mapping for section only (not include general information)
func getMappingSection(mapping [][]string, specificName string) [][]string {
	var mappingSpecific [][]string
	for _, tuple := range mapping {
		if strings.HasPrefix(tuple[0], specificName) {
			mappingSpecific = append(mappingSpecific, tuple)
		}
	}
	return mappingSpecific
}

// Return a slice with unique items
func appendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}
