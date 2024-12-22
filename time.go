package main

import (
	"strconv"
	"time"
)

func parseHexToUnix(hexStr string) (time.Time, error) {
	unixTime, err := strconv.ParseInt(hexStr, 16, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(unixTime, 0), nil
}
