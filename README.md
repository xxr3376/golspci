# golspci

`golspci` is a Go library that parses the output of the `lspci` command, make it easier to collect hardware info by program.  
This library will call `lspci -vmm -D` to get canonical output for easier parsing.

## Note

`lspci` provides `-n` flag, which show PCI vendor and device codes as numbers instead of looking them up in the [PCI ID list](https://pci-ids.ucw.cz/v2.2/pci.ids). The code never change between different version of `lspci`, but name does. Pass vendorInNumber=true when calling `lspci.New` to get codes instead of name.

## Usage

```go
package main

import (
    "fmt"
    "github.com/xxr3376/golspci/lspci"
)

func main() {
    // vendorInNumber=false to get text version of vendor
    // vendorInNumber=true to get codes version of vendor
    l := lspci.New(false)

    if err := l.Run(); err != nil {
        panic(err)
    }

    fmt.Println(l.Data)
    // You will get something like:
    /* map[
        0000:7f:14.1:
            map[
                SVendor:Intel Corporation
                SDevice:Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 0 Channel 1 Thermal Control
                Rev:02
                Slot:0000:7f:14.1
                Class:System peripheral
                Vendor:Intel Corporation
                Device:Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 0 Channel 1 Thermal Control
            ]
        0000:ff:14.0:
            map[
        ...
     */
}
```
