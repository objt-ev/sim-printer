/* Simulated printer, listens on 9100
 */
package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"encoding/json"
	"io"
	"strconv"
	"path/filepath"
)

type Configuration struct {
	Rawport int `yaml:"rawport"`
	OpenViewer bool `yaml:"openviewer"`
}

func main() {
	// read config
	var configuration Configuration
	var err error
	configfile, err := os.Open("config.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	defer configfile.Close()
	decoder := json.NewDecoder(configfile)
	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
    
	// create printjobs folder if it does not exist yet
	err = os.MkdirAll("./printjobs", os.ModePerm)
	checkError(err)


	fmt.Println("=========================")
	fmt.Println("= Simulated Raw Printer =")
	fmt.Println("=========================")
	fmt.Println("Settings: RAW port   ", configuration.Rawport)
	fmt.Println("          OpenViewer ", configuration.OpenViewer)

	rawPrintServer(configuration)
}

func rawPrintServer(configuration Configuration) {
	var filename string
	service := ":" + strconv.Itoa(configuration.Rawport)
	jobcount := 0
	baseName := "printjobs/printjob-"

	tcpAddr, err := net.ResolveTCPAddr("tcp", service)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	fmt.Println("Started Listening to print requests...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		jobcount++
		filename = baseName + strconv.Itoa(jobcount) + ".pdf"

		go handleRawPrintJob(conn, filename, configuration)
	}
}

func handleRawPrintJob(conn net.Conn, filename string, configuration Configuration) {
	// close connection on exit
	defer conn.Close()

	fmt.Printf("Receiving PrintJob %s \n", filename)

	file, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		return
	}
	defer file.Close()

	_, err = io.Copy(file, conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		return
	}
	fmt.Printf("Saved file %s \n", filename)

	if (configuration.OpenViewer) {
		// Start default viewer
	    abs, _ := filepath.Abs(filename)
		fmt.Printf("Start viewer for %s \n", abs)

		cmd := exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", abs)
		cmd.Start()
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

