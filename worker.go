/*
This file has been modified.
*/

package main

import (
    "net"
    "github.com/Sirupsen/logrus"
    "time"
    "sync"
    "os"
    "strings"
    "github.com/timest/gomanuf"
)

var log = logrus.New()
// ipNet Store IP address and subnet mask
var ipNet *net.IPNet
//This machine's mac address, need to use Ethernet packet
var localHaddr net.HardwareAddr
var iface string
// Stores the final data, key[string] holds the IP address
var data map[string]Info
// Timer, in a period of time no new data is written to data, exit the program, otherwise reset the timer
var t *time.Ticker
var do chan string

const (
    // 3 second timer
    START = "start"
    END = "end"
)

type Info struct {
    // IP Address
    Mac      net.HardwareAddr
    // CPU name
    Hostname string
    // Vendor information
    Manuf    string
}

// Add the captured data set to data and reset the timer
func pushData(ip string, mac net.HardwareAddr, hostname, manuf string) {
    // Stop timer
    do <- START
    var mu sync.RWMutex
    mu.RLock()
    defer func() {
        // reset the timer
        do <- END
        mu.RUnlock()
    }()
    if _, ok := data[ip]; !ok {
        data[ip] = Info{Mac: mac, Hostname: hostname, Manuf: manuf}
        return
    }
    info := data[ip]
    if len(hostname) > 0 && len(info.Hostname) == 0 {
        info.Hostname = hostname
    }
    if len(manuf) > 0 && len(info.Manuf) == 0 {
        info.Manuf = manuf
    }
    if mac != nil {
        info.Mac = mac
    }
    data[ip] = info
}

func setupNetInfo(f string) {
    var ifs []net.Interface
    var err error
    if f == "" {
        ifs, err = net.Interfaces()
    } else {
        // Iface has been selected
        var it *net.Interface
        it, err = net.InterfaceByName(f)
        if err == nil {
            ifs = append(ifs, *it)
        }
    }
    if err != nil {
        log.Fatal("Unable to get local network information:", err)
        goto END
    }
    for _, it := range ifs {
        addr, _ := it.Addrs()
        for _, a := range addr {
            if ip, ok := a.(*net.IPNet); ok && !ip.IP.IsLoopback() {
                if ip.IP.To4() != nil {
                    ipNet = ip
                    localHaddr = it.HardwareAddr
                    iface = it.Name
                    goto END
                }
            }
        }
    }
    END:
    if ipNet == nil || len(localHaddr) == 0 {
        log.Fatal("Unable to get local network information")
    }
}

func localHost() {
    host, _ := os.Hostname()
    data[ipNet.IP.String()] = Info{Mac: localHaddr, Hostname: strings.TrimSuffix(host, ".local"), Manuf: manuf.Search(localHaddr.String())}
}

func sendARP() {
    // ips is a collection of intranet IP addresses
    ips := Table(ipNet)
    for _, ip := range ips {
        go sendArpPackage(ip)
    }
}