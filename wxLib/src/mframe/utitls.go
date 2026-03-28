package mframe

import (
	"os"
	"strconv"
)

func GetServerIdFromRun() int {
	svrid := -1
	if len(os.Args) > 1 {
		svrid, _ = strconv.Atoi(os.Args[1])
	}
	return svrid
}
