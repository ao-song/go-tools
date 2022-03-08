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
var tstamp = make([]time.Time, 0)

var period = 2 * 60 * 60
var interval = time.Second

func collectTstamp() {
	for {
		tstamp = append(tstamp, time.Now())

		for len(tstamp) > period {
			tstamp = tstamp[1:]
		}

		time.Sleep(interval)
	}
}

func collectCPUStat() {

	counts, err := cpu.Counts(false)
	if err != nil {
		fmt.Printf("Get cpu counts failed, err: %v\n", err)
	}

	fmt.Printf("CPU count: %v\n", counts)

	for {
		percent, err := cpu.Percent(time.Second, false)
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

	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Past 2 hours",
			Type: "category"}),
		charts.WithTitleOpts(opts.Title{
			Title:    "CPU/Memory Status",
			Subtitle: "CPU/Memory Usage Percentage",
		}))

	// Put data into instance
	line.SetXAxis(tstamp).
		AddSeries("CPU Usage", cpu_percents).
		AddSeries("Memory Usage", mem_percents).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	line.Render(w)
}

// func main() {
// 	go collectTstamp()
// 	go collectCPUStat()
// 	go collectMemStat()

// 	http.HandleFunc("/", httpserver)
// 	http.ListenAndServe(":8081", nil)
// }
