# srf-fwu

This tool can be used to perform firmware upgrades (or downgrades) on SharkRF
devices which have the SharkRF Bootloader and it's USB serial console.

The USB serial console is available in SharkRF Bootloader v5 and up.

## Installing

```
go install github.com/jacobsa/go-serial/serial
go get github.com/sharkrf/srf-fwu
cd $GOPATH/sharkrf/srf-fwu
make
```

You can print the available command line switches by running srf-fwu with the
-h command line switch.
