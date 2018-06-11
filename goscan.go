package main

import "C"

import (
    "github.com/Sirupsen/logrus"
    "time"
    "context"
    "sort"
)

// Format the output
// xxx.xxx.xxx.xxx  xx:xx:xx:xx:xx:xx  hostname  manuf
// xxx.xxx.xxx.xxx  xx:xx:xx:xx:xx:xx  hostname  manuf

type ScanResult struct {
    IPAddress string
    MACAddress string
    Hostname string
    Manufacturer string
}

func GetData() []ScanResult {
    var result []ScanResult
    var keys IPSlice
    for k := range data {
        keys = append(keys, ParseIPString(k))
    }
    sort.Sort(keys)
    for _, k := range keys {
        d := data[k.String()]
        mac := ""
        if d.Mac != nil {
            mac = d.Mac.String()
        }
        result = append(result, ScanResult{k.String(), mac, d.Hostname, d.Manuf})
    }

    return result
}

//export GetDefaultScanResults
func GetDefaultScanResults() []ScanResult {
    return GetCustomScanResults("", false)
}

//export GetDefaultScanResultsVerbose
func GetDefaultScanResultsVerbose() []ScanResult {
    return GetCustomScanResults("", true)
}

//export GetCustomScanResults
func GetCustomScanResults(paramIface string, verbose bool) []ScanResult{
    if verbose {
        log.Level = logrus.InfoLevel
    } else {
        log.Level = logrus.WarnLevel
    }
    // initialization data
    data = make(map[string]Info)
    do = make(chan string)
    // Initialize network information
    setupNetInfo(paramIface)

    ctx, cancel := context.WithCancel(context.Background())
    go listenARP(ctx)
    go listenMDNS(ctx)
    go listenNBNS(ctx)
    go sendARP()
    go localHost()

    t = time.NewTicker(4 * time.Second)
    var result []ScanResult
    for {
        select {
        case <-t.C:
            result = GetData()
            cancel()
            goto END
        case d := <-do:
            switch d {
            case START:
                t.Stop()
            case END:
                // Received new data, reset counter for 2 seconds
                t = time.NewTicker(2 * time.Second)
            }
        }
    }
END:
    return result
}

// required, but not used
func main() { }
