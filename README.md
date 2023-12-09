# SimRawPrinter
Simulation of a printer receiving PDF print jobs using the RAW protocol. Written in Golang

If you want to find out what is sent to the printer when you are sending pdf print jobs to it this
program will help you to find out. It implements a RAW protocol listener (typically on port 9100). 
You can change the port in the config file config.json.

This program stores individual print jobs in local files named printjob-x.pdf,
where x is asequentially growing integer starting with 1.

Optionally, the default system pdf viewer is launched automatically.
You can enable/disable this option in the config file.
