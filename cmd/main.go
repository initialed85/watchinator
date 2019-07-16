package main

import (
	"bytes"
	"compress/zlib"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("usage: watchinator [file] [UDPv4/v6 address]")
	}

	fileName := os.Args[1]
	rawAddr := os.Args[2]

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	err = watcher.Add(fileName)
	if err != nil {
		log.Fatal(err)
	}

	network := "udp4"
	if strings.Count(rawAddr, ":") > 1 {
		network = "udp6"
	}

	addr, err := net.ResolveUDPAddr(network, rawAddr)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP(network, nil, addr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op == fsnotify.Remove {
				log.Fatal("file removed, exiting")
			} else if event.Op != fsnotify.Write {
				continue
			}

			sourceData, err := ioutil.ReadFile(fileName)
			if err != nil {
				log.Fatal(err)
			}

			var b bytes.Buffer

			w := zlib.NewWriter(&b)

			_, err = w.Write(sourceData)
			if err != nil {
				log.Fatal(err)
			}

			err = w.Close()
			if err != nil {
				log.Fatal(err)
			}

			compressedData := b.Bytes()

			log.Printf("file changed, sending %v bytes to %v", len(compressedData), addr)

			_, err = conn.Write(compressedData)
			if err != nil {
				// fire and forget
			}
		case err := <-watcher.Errors:
			log.Fatal(err)
		}
	}

	defer watcher.Close()
}
