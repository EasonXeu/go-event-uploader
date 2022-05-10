package main

import (
	"fmt"
	"github.com/go-echarts/statsview"
	"github.com/go-echarts/statsview/viewer"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	// start a memory monitor
	// Visit your browser at http://localhost:18066/debug/statsview
	// Or debug as always via http://localhost:18066/debug/pprof, http://localhost:18066/debug/pprof/heap, ...
	viewer.SetConfiguration(viewer.WithTheme(viewer.ThemeWesteros), viewer.WithAddr("localhost:18066"))
	mgr := statsview.New()
	go mgr.Start()
	time.Sleep(2 * time.Second)
	// the business code
	produce := InitProduce()
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Kill, os.Interrupt)
	produce.Start()
	var m sync.WaitGroup
	for i := 0; i < 20; i++ {
		m.Add(1)
		go func() {
			defer m.Done()
			for i := 0; i < 1349; i++ {
				log := generateLog()
				var err error
				err = produce.SendLog(log)
				if err != nil {
					fmt.Println(err)
				}
			}
		}()
	}
	m.Wait()
	fmt.Println("Send completion")

	if _, ok := <-ch; ok {
		fmt.Println("Get the shutdown signal and start to shut down")
		produce.Stop()
		fmt.Println(SendCount.Load())
		fmt.Println(ErrorCount.Load())
	}
	time.Sleep(5 * time.Second)
	mgr.Stop() // stop the analyser tool

}
