# sim-printer
Simulation of printer receiving print jobs using LPD and RAW protocols. Written in Golang

If you want to find out what is sent to the printer when you are sending print job to it this
program will help you to find out. It implements LPD server (usually listening on port 515)
as well as RAW protocol listener (in printers usually on port 9100). You can change those ports
in file printerconfig.json.

Program stores individual printer jobs in local files named lprjob#number.prn or rawjob#number.prn,
where number is sequentially growing integer starting with 1.

Program was tested on all combinations of desktop operatinmg systems and processor architectures
on the market - Windows, Linux, MacOS - on Intel as well as arm64 processors.
