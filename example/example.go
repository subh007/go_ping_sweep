// example
package main

import (
	"flag"
	"fmt"
	"github.com/subh007/go_ping_sweep"
	tm "github.com/buger/goterm"
	"os"
	"strconv"
)

func main() {
	host := flag.String("host", "", "ip address/ hostname to analyse the ping.")
	flag.Parse()

	if host == nil {
		fmt.Println("usage: ./example --host <host name/ip>")
		os.Exit(-1)
	}

	if go_ping_sweep.IsAdmin() {
		t, err := go_ping_sweep.PingAnalyse(*host, 10)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}

		//	 create the graph with the output
		fmt.Println("Ping Analysis (Graph)")

		tm.Clear()
		tm.MoveCursor(5, 0)

		chart := tm.NewLineChart(100, 20)
		data := new(tm.DataTable)
		data.AddColumn("Ping")
		data.AddColumn("Delay (ns)")


		pingTime := t.GetColumn("TimePing (ns)")
		for i := 0; i < len(pingTime); i++ {
			pingTimeI, _ := strconv.Atoi(pingTime[i])
			data.AddRow(float64(i), float64(pingTimeI))
		}

		tm.Println(chart.Draw(data))
		tm.Flush()

		fmt.Println("Ping Analysis (Table)")
		t.CreateTable()
	}
}
