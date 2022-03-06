package main

import (
	"fmt"
	"time"

	"net/http"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

var cpu_percents = make([]opts.LineData, 0)
var mem_percents = make([]opts.LineData, 0)

var period = 2 * 60 * 12
var interval = 5 * time.Second

func collectCPUStat() {

	counts, err := cpu.Counts(false)
	if err != nil {
		fmt.Printf("Get cpu counts failed, err: %v\n", err)
	}

	fmt.Printf("CPU count: %v\n", counts)

	for {
		percent, err := cpu.Percent(1, false)
		if err != nil {
			fmt.Printf("Get cpu percent failed, err: %v\n", err)
		}
		cpu_percents = append(cpu_percents, opts.LineData{Value: percent[0]})

		for len(cpu_percents) > period {
			cpu_percents = cpu_percents[1:]
		}

		time.Sleep(interval)
	}
}

func collectMemStat() {

	for {
		memstat, err := mem.VirtualMemory()
		if err != nil {
			fmt.Printf("Get mem status failed, err: %v\n", err)
		}
		mem_percents = append(mem_percents, opts.LineData{Value: memstat.UsedPercent})

		for len(mem_percents) > period {
			mem_percents = mem_percents[1:]
		}

		time.Sleep(interval)
	}
}

func httpserver(w http.ResponseWriter, _ *http.Request) {
	line := charts.NewLine()

	s := make([]int, period)
	start := 0
	for i := range s {
		s[i] = start
		start += 1
	}

	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "CPU/Memory Status",
			Subtitle: "CPU/Memory Usage Percentage",
		}))

	// Put data into instance
	line.SetXAxis(s).
		AddSeries("CPU Usage", cpu_percents).
		AddSeries("Memory Usage", mem_percents).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	line.Render(w)
}

func main() {
	go collectCPUStat()
	go collectMemStat()

	http.HandleFunc("/", httpserver)
	http.ListenAndServe(":8081", nil)
}
