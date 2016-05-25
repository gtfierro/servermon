package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// Plagiarised
// from http://stackoverflow.com/questions/11356330/getting-cpu-usage-with-golang
func getCpuUsage() (used, total uint64) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	var idle uint64
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val // tally up all the numbers to get total ticks
				if i == 4 {  // idle is the 5th field in the cpu line
					idle = val
				}
			}
			return total - idle, total
		}
	}
	return
}

func getMemUsage() (used, total uint64) {
	contents, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return
	}
	var free uint64
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 3 {
			if fields[0] == "MemTotal:" {
				total, _ = strconv.ParseUint(fields[1], 10, 64)
				total *= 1024
			}
			if fields[0] == "MemFree:" {
				free, _ = strconv.ParseUint(fields[1], 10, 64)
				free *= 1024
			}
		}
	}
	return total - free, free
}
