package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mmtaee/ocserv-users-management/common/pkg/config"
	"github.com/mmtaee/ocserv-users-management/common/pkg/database"
	"github.com/mmtaee/ocserv-users-management/log_stream/internal/readers"
	"github.com/mmtaee/ocserv-users-management/log_stream/internal/sse"
	"github.com/mmtaee/ocserv-users-management/log_stream/internal/stats"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	debug   bool
	host    string
	port    int
	systemd bool
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	flag.BoolVar(&debug, "d", false, "debug mode")
	flag.StringVar(&host, "h", "0.0.0.0", "Server Host")
	flag.IntVar(&port, "p", 8080, "Server Port")
	flag.BoolVar(&systemd, "systemd", false, "Systemd Mode")
	flag.Parse()

	if debug {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}

	ctx, cancel := context.WithCancel(context.Background())
	service := "ocserv"

	config.Init(debug, host, port)
	cfg := config.Get()
	database.Connect()

	streamChan := make(chan string, 1000)
	lineLogChan := make(chan string, 1000)
	broadcastChan := make(chan string, 1000)

	if systemd {
		log.Println("Running on host – using systemd logs")
		go func() {
			if err := readers.SystemdStreamLogs(ctx, service, streamChan); err != nil {
				log.Printf("Systemd log error: %v\n", err)
			}
		}()
	} else {
		log.Println("Running in Docker – using Docker logs")
		go func() {
			if err := readers.DockerStreamLogs(ctx, service, streamChan); err != nil {
				log.Println(err)
			}
		}()
	}

	statService := stats.NewStatService(ctx, lineLogChan)
	go func() {
		statService.CalculateUserStats()
	}()

	sseServer := sse.NewSSEServer()
	sseServer.StartBroadcast(broadcastChan)

	go func() {
		server := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
		http.HandleFunc("/logs", sseServer.SSEHandler())
		log.Println("Starting server on ", server)
		if err := http.ListenAndServe(server, nil); err != nil {
			log.Fatalf("ListenAndServe failed: %v", err)
		}
	}()

	go func() {
		start(ctx, streamChan, broadcastChan, lineLogChan)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("\nReceived signal: %s\n", sig)
		cancel()
	}()

	<-ctx.Done()
	log.Println("Service shutting down successfully")
}

func start(ctx context.Context, streamText <-chan string, broadcaster, lineLogChan chan<- string) {
	for {
		select {
		case <-ctx.Done():
			return
		case line, ok := <-streamText:
			if !ok {
				return
			}
			// Send to broadcaster
			go func(l string) {
				select {
				case broadcaster <- l:
				case <-ctx.Done():
					return
				default:
					// skip log, continue
				}
			}(line)

			// Send to lineLogChan
			go func(l string) {
				select {
				case lineLogChan <- l:
				case <-ctx.Done():
				}
			}(line)
		}
	}
}
