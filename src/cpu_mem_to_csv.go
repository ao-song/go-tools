package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func main() {
	file, err := os.Create("cpu_mem_stat.csv")
	defer file.Close()
	if err != nil {
		fmt.Printf("Failed to open csv file, %v\n", err)
	}

	w := csv.NewWriter(file)
	defer w.Flush()

	h, _ := strconv.Atoi(os.Args[1])
	hours := h * 60 * 60
	count := 0
	flushCount := 0

	for {
		cpuPercent, err := cpu.Percent(time.Second, false)
		if err != nil {
			fmt.Printf("Get cpu percent failed, err: %v\n", err)
		}

		memStat, err := mem.VirtualMemory()
		if err != nil {
			fmt.Printf("Get mem status failed, err: %v\n", err)
		}

		row := []string{time.Now().Format(time.UnixDate), fmt.Sprintf("%f", cpuPercent[0]), fmt.Sprintf("%f", memStat.UsedPercent)}
		w.Write(row)

		flushCount += 1
		if flushCount > 10 {
			w.Flush()
			flushCount = 0
		}

		time.Sleep(time.Second)

		count += 1
		if count > hours {
			break
		}
	}
}
