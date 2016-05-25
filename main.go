package main

import (
	"fmt"
	"os"
	"time"

	"github.com/immesys/bw2bind"
)

func main() {
	//Connect
	cl := bw2bind.ConnectOrExit("")
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
	iface := svc.RegisterInterface("all", "i.top")
	hostname, _ := os.Hostname()
	iface.PublishSignal("hostname", []bw2bind.PayloadObject{
		bw2bind.CreateTextPayloadObject(bw2bind.PONumText,
			hostname),
	})

	for {
		usedC, totalC := getCpuUsage()
		cpuPercentage := float64(usedC) / float64(totalC) * 100.0
		cpo := bw2bind.CreateTextPayloadObject(bw2bind.PONumText,
			fmt.Sprintf("%.2f %%", cpuPercentage))
		iface.PublishSignal("cpu", []bw2bind.PayloadObject{cpo})

		usedM, totalM := getMemUsage()
		mpo := bw2bind.CreateTextPayloadObject(bw2bind.PONumText,
			fmt.Sprintf("%d / %d MB", usedM/1024/1024, totalM/1024/1024))
		iface.PublishSignal("mem", []bw2bind.PayloadObject{mpo})

		time.Sleep(3 * time.Second)
	}
}
