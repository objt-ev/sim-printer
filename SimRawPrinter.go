/* Simulated printer, listens on 9100
 */
package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Configuration struct {
	Rawport int 
	OpenViewer bool 
}

func main() {
	// params
	rawPortPtr := flag.Int("rawport", 9100, "tcp port to use for raw printer.")
	openViewerPtr := flag.Bool("openviewer", true, "Open system viewer automatically")
	flag.Parse()

	var configuration Configuration
	configuration.Rawport = *rawPortPtr
	configuration.OpenViewer = *openViewerPtr
    
	// create printjobs folder if it does not exist yet
	err := os.MkdirAll("./printjobs", os.ModePerm)
	checkError(err)

	fmt.Println("=========================")
	fmt.Println("= Simulated Raw Printer =")
	fmt.Println("=========================")
	fmt.Printf("Port [%d], Openviewer [%t]\n\n", configuration.Rawport, configuration.OpenViewer)

	rawPrintServer(configuration)
}

func rawPrintServer(configuration Configuration) {
	
	service := ":" + strconv.Itoa(configuration.Rawport)
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	
	var filename string
	jobcount := 0
	baseName := "printjobs/printjob-"

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
	    abs, _ := filepath.Abs(filename)
		contentType, _ := GetFileContentType(abs)
		if (strings.HasSuffix(contentType, "pdf")) {
			fmt.Printf("Start viewer for %s \n", abs)
			
			cmd := exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", abs)
			cmd.Start()
		} else {			
		    fmt.Println("Viewer not started, content Type of file is: " + contentType )
		}
	}
}

func GetFileContentType(filename string) (string, error) {

   file, err := os.Open(filename)

   if err != nil {
      panic(err)
   }

   defer file.Close()

   // to sniff the content type only the first 512 bytes are used.
   buf := make([]byte, 512)
   _, err = file.Read(buf)
   if err != nil {
      return "", err
   }

   // the function that actually does the trick
   contentType := http.DetectContentType(buf)

   if err != nil {
      panic(err)
   }

   return contentType, nil
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

