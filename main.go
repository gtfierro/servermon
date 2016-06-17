package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	bw "gopkg.in/immesys/bw2bind.v5"
)

func main() {
	//Connect
	cl := bw.ConnectOrExit("")
	cl.SetEntityFromEnvironOrExit()

	uri := os.Getenv("BASEURI")
	if uri == "" {
		fmt.Println("You need to set $BASEURI")
		os.Exit(1)
	}
	svc := cl.RegisterService(uri, "s.top")
	if os.Getenv("DESC") != "" {
		svc.SetMetadata("description", os.Getenv("DESC"))
	}
	iface := svc.RegisterInterface("pantry", "i.top")
	fmt.Println(iface.SignalURI("cpu"))
	cl.SetMetadata(iface.SignalURI("cpu"), "SourceName", "top")
	cl.SetMetadata(iface.SignalURI("mem"), "SourceName", "top")

	cl.SetMetadata(iface.SignalURI("cpu"), "UnitofMeasure", "CPU")
	cl.SetMetadata(iface.SignalURI("mem"), "UnitofMeasure", "MEM")

	for {
		cpus, err := cpu.Percent(3*time.Second, true)
		cpuPercentage := float64(0)
		for _, num := range cpus {
			cpuPercentage += num
		}
		cpo, _ := bw.CreateMsgPackPayloadObject(bw.FromDotForm(bw.PODFMsgPack), map[string]interface{}{"cpu": cpuPercentage})
		err = iface.PublishSignal("cpu", cpo)
		if err != nil {
			log.Fatal(err)
		}

		memstat, err := mem.VirtualMemory()
		mpo, _ := bw.CreateMsgPackPayloadObject(bw.FromDotForm(bw.PODFMsgPack), map[string]interface{}{"used": memstat.UsedPercent})
		err = iface.PublishSignal("mem", mpo)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("cpu %v mem %v\n", cpuPercentage, memstat.UsedPercent)

		time.Sleep(3 * time.Second)
	}
}
