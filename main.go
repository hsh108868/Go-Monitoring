package main

import (
	"log"
	"time"
)

func cpuMonitor() {
	procStatRaws, err := cpu.ParseProcStatFirstTime()

	cpuError := make(chan error)
	cpuStats := make(chan []model.cpuStat)
	present := make(chan []model.procStatRaw)

	if err != nil {
		log.Println("Cant Collect CPU Data", err)

	}

	for {
		time.Sleep(10 * time.Second)
		log.Println(procStatRaws[0].TotalTime)
		go cpu.GetCPUStatAsync(procStatRaws, present, cpuStats, cpuError)
		for i := 0; i < 2; i++ {
			select {
			case CPUStats := <-cpuStats:
				log.Println("CPU Used Percent(", CPUStats[0].Timestamp, ")", CPUStats[0].Used, "%")
			case PRESENT := <-present:
				procStatRaws = PRESENT
			case cpuError := <-cpuError:
				log.Println("Cant Collect CPU Data", cpuError)
				break
			}
		}
	}
}

func main() {

	err := config.GetEnvironmentVariable()

	influx.InfluxInit()
	wg.Add(2)
	go memoryMonitor()
	go cpuMonitor()
	wg.wait()
}
