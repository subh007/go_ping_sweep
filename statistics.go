package go_ping_sweep

func average(results []Result) int64 {
	var average int64
	for _, res := range results {
		average += res.TimePing
	}

	return average / int64(len(results))
}
