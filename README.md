# goscan

![image](https://user-images.githubusercontent.com/1621058/32154543-63c4e560-bcff-11e7-8a92-5281e18f221e.png)

**Features:**
* Scan the entire intranet IPv4 space
* Send ARP packets to the entire intranet
* Display IP/MAC address/host name/device vendor name
* Use SMB (Windows) and mDNS (Mac OS) to sniff intranet hostname (hostname)
* Use MAC address to calculate device manufacturer information


**Modifications from [timest/goscan](https://github.com/timest/goscan)**
* Complete English translation
* Removal of main() method and introduction of API methods
  - Intended for use as an exported C library
  - Returns an array of structs with results

### Usage: ###

```
$ go build test/main.go
$ sudo ./main
```


More details can be viewed [here](https://github.com/timest/goscan/issues/1).