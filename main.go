package main

import (
	"fmt"
	"log"
	"os"
	"time"

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

	for {
		usedC, totalC := getCpuUsage()
		cpuPercentage := float64(usedC) / float64(totalC) * 100.0
		cpo, _ := bw.CreateMsgPackPayloadObject(bw.FromDotForm(bw.PODFMsgPack), map[string]interface{}{"cpu": cpuPercentage})
		err := iface.PublishSignal("cpu", cpo)
		if err != nil {
			log.Fatal(err)
		}

		usedM, totalM := getMemUsage()
		mpo, _ := bw.CreateMsgPackPayloadObject(bw.FromDotForm(bw.PODFMsgPack), map[string]interface{}{"used": usedM / 1024 / 1024, "total": totalM / 1024 / 1024})
		err = iface.PublishSignal("mem", mpo)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("cpu %v mem %v\n", cpuPercentage, usedM/1024/1024)

		time.Sleep(3 * time.Second)
	}
}
