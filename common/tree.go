package common

import (
	"strings"
)

type treeConnector string

const (
	SIMPLE_CONNECTOR treeConnector = "│   "
	LAST_CONNECTOR   treeConnector = "└── "
	BRANCH_CONNECTOR treeConnector = "├── "
	SIMPLE_INDENT    treeConnector = "    "
	COUNTER_INDENT   treeConnector = "      " // 6 spaces for counter alignment
)

func ReplaceSuffix(start string, old, new treeConnector) string {
	if strings.HasSuffix(start, string(old)) {
		return strings.TrimSuffix(start, string(old)) + string(new)
	}
	return start
}

func AdvanceSuffix(start string) string {
	if strings.HasSuffix(start, string(BRANCH_CONNECTOR)) {
		return ReplaceSuffix(start, BRANCH_CONNECTOR, SIMPLE_CONNECTOR)
	} else if strings.HasSuffix(start, string(LAST_CONNECTOR)) {
		return ReplaceSuffix(start, LAST_CONNECTOR, SIMPLE_INDENT)
	} else {
		return start
	}
}
