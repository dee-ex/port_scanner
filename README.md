# port_scanner
A Simple Port Scanner in Golang
## Usage
Run command
```
go run main.go hostname
```
where hostname can be a domain name or ip address. :D  
This tool helps you scan from port No.1 to port No.65000.  
If you wanna scan in you desirable range (e.g `[a, b]`), run command  
```
go run main.go hostname a-b
```
In case you forget `a`, like
```
go run main.go hostname -b
```
The tool will scan from 1 to b. Or, if you drop `-b`
```
go run main.go hostname a
```
The range which the tool scan is from `a` to `65000`.
