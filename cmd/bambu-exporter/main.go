package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/evanofslack/bambulab-client/monitor"
	"github.com/evanofslack/bambulab-client/mqtt"

	"github.com/evanofslack/bambulab-exporter/exporter"
	"github.com/evanofslack/bambulab-exporter/config"
)

const (
    defaultPort = "8080"
)

func main() {
	ctx := context.Background()
    config, err := config.New()
    if err != nil {
        panic(err)
    }

	client, err := mqtt.NewCloudClient(config.Auth.Endpoint, config.Auth.DeviceId, config.Auth.Username, config.Auth.Password)
	if err != nil {
		panic(err)
	}

	if err := client.Connect(); err != nil {
		panic(err)
	}

	// Chan to recieve messages on
	msgs := make(chan mqtt.Message)

	// Subscribe for info updates
	fmt.Println("subscribing")
	client.Subscribe(msgs)

	monitor := monitor.New()
	fmt.Println("starting monitor")
	go monitor.Start(msgs)

	exporter, err := exporter.New(monitor, config.Auth.DeviceId)
	if err != nil {
		panic(err)
	}
	fmt.Println("starting exporter")
	go exporter.Start(ctx)

    port := config.HTTP.Port
    if port == "" {
        port = defaultPort
    }
	go func (){
        if err := exporter.Serve(port); err != nil {
            panic(err)
        }
    }()


    // Pause for 1 second before sending first request for all data
    select {
    case <- ctx.Done():
        break
    case <- time.After(time.Second):
        go requestAll(ctx, client)
    }

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	<-sig
	fmt.Println("")
	client.Disconnect()
	fmt.Println("shutdown complete")
}

// Start thread requesting all every 5 mins
func requestAll(ctx context.Context, client *mqtt.Client) {
	// Make an initial request for all info
	if err := client.PublishPushAll(ctx); err != nil {
		fmt.Println(err)
	}

	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			if err := client.PublishPushAll(ctx); err != nil {
				fmt.Println(err)
			}
		}
	}
}

func printMsgs(ctx context.Context, msgs chan mqtt.Message) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-msgs:
			fmt.Printf("got raw msg:\n")
			b, err := json.MarshalIndent(msg, "", "  ")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(b))
		}
	}
}

func printMonitor(ctx context.Context, mon *monitor.Monitor) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-mon.Update:
			fmt.Println("got state update:")
			b, err := json.MarshalIndent(mon.State, "", "  ")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(b))
		}
	}
}
