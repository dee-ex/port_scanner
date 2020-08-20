package main

import "fmt"
import "net"
import "strings"
import "strconv"
import "time"
import "runtime"
import "sort"
import "flag"

type scan_result struct {
    Port int
    State string
}

func check_ip_address(ip string) (bool, string) {
    if net.ParseIP(ip) == nil {
        return false, "x"
    }
    if strings.Contains(ip, ".") {
        return true, "v4"
    }
    return true, "v6"
}

func check_hostname(hostname string) bool {
    if _, err := net.LookupIP(hostname); err != nil {
        return false
    } else {
        return true
    }
}

func get_range(r []string) (int, int) {
    var start, stop int64 = 1, 65000

    if r[0] != "" {
        start, _ = strconv.ParseInt(r[0], 10, 32)
    }
    if len(r) == 2 {
        stop, _ = strconv.ParseInt(r[1], 10, 32)
    }
    return int(start), int(stop)
}

func scan_port(protocol, hostname string, port int, lock chan bool) scan_result {
    result := scan_result {Port: port}
    address := hostname + ":" + strconv.Itoa(port)
    lock <- true
    conn, err := net.DialTimeout(protocol, address, 1*time.Second)
    <-lock

    if err != nil {
        result.State = "closed"
        return result
    }

    defer conn.Close()
    result.State = "open"
    return result
}

func disp_info(hostname, protocol string, start, stop , num_port int, open_ports []int) {
    var num_closed_port = num_port - len(open_ports)
    fmt.Printf("Scan report for %s\n", hostname)
    fmt.Printf("From port %d -> %d\n", start, stop)
    fmt.Printf("Not shown: %d closed ports\n", num_closed_port)
    fmt.Println("PORT       STATE")
    for _, p := range open_ports {
        fmt.Printf("%-9s  open\n", strconv.Itoa(p) + "/" + protocol)
    }
}

func main() {
    var hostname, r, protocol string
    flag.StringVar(&hostname, "host", "127.0.0.1", "hostname that long to scan")
    flag.StringVar(&r, "r", "1-65000", "range that long to scan")
    flag.StringVar(&protocol, "p", "tcp", "protocal that long to scan")
    flag.Parse()

    if !check_hostname(hostname) {
        fmt.Println("Invalid hostname")
        return
    }

    start, stop := get_range(strings.Split(r, "-"))

    runtime.GOMAXPROCS(100)
    var limit_chan = 4000

    var res scan_result
    var open_ports []int

    lock_chan := make(chan bool, limit_chan)
    results_chan := make(chan scan_result, limit_chan)

    fmt.Println("Starting")
    start_time := time.Now()

    for p := start; p <= stop; p++ {
        go func(p int) {
            results_chan <- scan_port(protocol, hostname, p, lock_chan)
        }(p)
    }

    num_port := stop - start + 1
    for x := 0; x < num_port; x++ {
        res = <-results_chan
        if res.State == "open" {
            open_ports = append(open_ports, res.Port)
        }
    }

    elapsed := time.Since(start_time)

    close(results_chan)
    close(lock_chan)

    sort.Ints(open_ports)
    disp_info(hostname, protocol, start, stop, num_port, open_ports)

    fmt.Printf("Done: scanned in %s seconds", elapsed)
}
