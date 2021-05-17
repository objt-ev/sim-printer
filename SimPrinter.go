/* Simulated printer, listens on 9100 and 515
 */
package main

import (
	"fmt"
	"net"
	"os"
//	"os/exec"

	//  "strings"
	"encoding/json"
	"io"
	"io/ioutil"
	"strconv"
)

type Configuration struct {
	Lprport string `json:"lprport"`
	Rawport string `json:"rawport"`
}

func main() {
	var configuration Configuration
	configfile, err := os.Open("printerconfig.json")
	if err != nil {
    configuration.Lprport = "515"
    configuration.Rawport = "9100"
	} else {
  	defer configfile.Close()
  	decoder := json.NewDecoder(configfile)
  	err := decoder.Decode(&configuration)
  	if err != nil {
  		fmt.Println("error:", err)
  	}  
  }
  

	fmt.Println("Simulated printer, Jan Ftacnik, 2021 \n")

	fmt.Println("Using LPR port: ", configuration.Lprport)
	fmt.Println("Using RAW port: ", configuration.Rawport)

	go server9100(configuration.Rawport)
	server515(configuration.Lprport)

}

func server515(lprport string) {
	var name string
	service := ":" + lprport
	icount := 0
	name1 := "lprjob"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		// run as a goroutine
		icount = icount + 1
		name = name1 + strconv.Itoa(icount)
		go handleClient515(conn, name)
	}
}

func handleClient515(conn net.Conn, filename string) {
	// close connection on exit
	defer conn.Close()

	//      file, err4 := os.Create(filename)
	//      if err4 != nil {
	//        fmt.Fprintf(os.Stderr, "Fatal error: %s", err4.Error())
	//        return
	//      }
	//    defer file.Close()

	var buf [512]byte
	var code byte = 0x0
	for {
		ack := []byte{code}

		// read up to 512 bytes
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}

		//           queue := string(buf[1:n])
		//           fmt.Println("   queue: ",queue)
		// write the n bytes read
		_, err2 := conn.Write(ack)
		if err2 != nil {
			fmt.Fprintf(os.Stderr, "Could not write message back....")
			return
		}

		// read up to 512 bytes
		_, err = conn.Read(buf[0:])
		if err != nil {
			return
		}

		//           work := string(buf[1:n])
		//           i :=  strings.Index(work, "cfA")
		//           if i>-1 {
		//             comp  := string(buf[7+i:n-1])
		//             jobid := string(buf[4+i:7+i])
		//             fmt.Println("   comp: ",comp,"   jobid: ",jobid)
		//           }
		// write the n bytes read
		_, err2 = conn.Write(ack)
		if err2 != nil {
			fmt.Fprintf(os.Stderr, "Could not write message back....")
			return
		}

		// read up to 512 bytes
		n, err = conn.Read(buf[0:])
		if err != nil {
			return
		}

		errc := ioutil.WriteFile(filename+".cfg", buf[0:n], 0666)
		if errc != nil {
			fmt.Fprintf(os.Stderr, "Fatal error: %s", errc.Error())
		}

		//           work5 := string(buf[0:n])
		//           s := strings.Fields(work5)
		//           for index, element := range s {
		//		         fmt.Println(index, "--", element)
		//           }

		// write the n bytes read
		_, err2 = conn.Write(ack)
		if err2 != nil {
			fmt.Fprintf(os.Stderr, "Could not write message back....")
			return
		}

		// read up to 512 bytes
		_, err = conn.Read(buf[0:])
		if err != nil {
			return
		}

		//         fmt.Println(n," Read: ",string(buf[1:n]))

		// write the n bytes read
		_, err2 = conn.Write(ack)
		if err2 != nil {
			fmt.Fprintf(os.Stderr, "Could not write message back....")
			return
		}

		// read up to 512 bytes
		n, err = conn.Read(buf[0:])
		if err != nil {
			return
		}

		err5 := ioutil.WriteFile(filename+".prn", buf[0:n], 0666)
		if err5 != nil {
			fmt.Fprintf(os.Stderr, "Fatal error: %s", err5.Error())
		}

		//OPEN FILE TO APPEND INTO
		file, err6 := os.OpenFile(filename+".prn", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err6 != nil {
			fmt.Fprintf(os.Stderr, "Fatal error: %s", err6.Error())
			return
		}
		defer file.Close()

		//           fmt.Println(n," Header of data file: ",string(buf[1:n]))

		// write the n bytes read
		_, err2 = conn.Write(ack)
		if err2 != nil {
			fmt.Fprintf(os.Stderr, "Could not write message back....")
			return
		}

		// read up to 512 bytes

		_, err = io.Copy(file, conn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
			return
		}

		//           n, err    = conn.Read(buf[0:])
		//           if err   != nil {
		//              return
		//           }

		//           fmt.Println(" Data file done \n")

		// write the n bytes read
		_, err2 = conn.Write(ack)
		if err2 != nil {
			fmt.Fprintf(os.Stderr, "Could not write message back....")
			return
		}

		file.Close()
		fmt.Printf("Saved file %s \n", filename+".prn")
//		out, errex := exec.Command(`cmd.exe`, `/C`, `E:\Go\src\simprinter\copyfile.bat`, filename+".prn").Output()
//		if errex != nil {
//			fmt.Printf("Problem executing batch file: %s", errex)
//		}
//		fmt.Printf("Output: %s", out)

	}
}

func server9100(rawport string) {
	var name string
	service := ":" + rawport
	jcount := 0
	name1 := "rawjob"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		// run as a goroutine
		jcount = jcount + 1
		name = name1 + strconv.Itoa(jcount) + ".prn"

		go handleClient9100(conn, name)
	}
}

func handleClient9100(conn net.Conn, filename string) {
	// close connection on exit
	defer conn.Close()

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
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
