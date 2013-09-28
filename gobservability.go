package gobservability

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	//"syscall"
	"time"
)

type J map[string]interface{}

func (r J) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

func Run(listen_port string, interval time.Duration) {

	start_time := time.Now()

	var output J

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, output)
		return
	}

	ServeMux := http.NewServeMux()
	ServeMux.HandleFunc("/", handler)

	HttpServer := &http.Server{Addr: listen_port, Handler: ServeMux}
	go HttpServer.ListenAndServe()

	for {
		memStats := runtime.MemStats{}
		runtime.ReadMemStats(&memStats)

		/*
			//do something with this soon
			rLimitStats := &syscall.Rlimit{}
			syscall.Getrlimit(syscall.RLIMIT_NOFILE, rLimitStats)

			rUsageStats := &syscall.Rusage{}
			syscall.Getrusage(syscall.RUSAGE_SELF, rUsageStats)

			statsFs := &syscall.Statfs_t{}
			syscall.Statfs("/", statsFs)
		*/

		hostname, _ := os.Hostname()

		host_details := J{
			//"statFs":  statsFs,
			"name":    hostname,
			"num_cpu": runtime.NumCPU(),
		}

		memory_details := J{
			"allocated":        memStats.Alloc,
			"mallocs":          memStats.Mallocs,
			"frees":            memStats.Frees,
			"total_pause_time": int64(memStats.PauseTotalNs) / int64(1000000), //ms
			"heap":             memStats.HeapAlloc,
			"stack":            memStats.StackInuse,
		}

		proc_details := J{
			"uptime":     int(time.Since(start_time).Seconds()),
			"memory":     memory_details,
			"goroutines": runtime.NumGoroutine(),
			"gomaxprocs": runtime.GOMAXPROCS(0),
			//"rlimit_nofiles": rLimitStats,
			//"rusage":         rUsageStats,
			"args": os.Args,
		}

		go_details := J{
			"version": runtime.Version(),
		}

		output = J{"host": host_details, "proc": proc_details, "go": go_details}
		time.Sleep(interval)
	}
}
