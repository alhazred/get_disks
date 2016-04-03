package main

import (
	"C"
	"fmt"
	"io/ioutil"
	"regexp"
	"syscall"
	"unsafe"
)

const sysIoctl = 54
const DKIOCREMOVABLE = (0x04 << 8) | 16
const DKIOCGMEDIAINFO = (0x04 << 8) | 42
const GB = 1 << (10 * 3)

type dk_minfo struct {
	dki_media_type C.uint
	dki_lbsize     C.uint
	dki_capacity   C.longlong
}

func main() {
	dsk, _ := ioutil.ReadDir("/dev/dsk/")
	for _, d := range dsk {
		re := regexp.MustCompile(`c.*d0p0$`)
		if re.MatchString(string(d.Name())) {
			rmdsk := uint(0)
			rdsk := fmt.Sprintf("/dev/rdsk/%s", string(d.Name()))
			fd, e := syscall.Open(rdsk, syscall.O_RDONLY, 0600)
			if e == nil {
				_, _, err := syscall.Syscall(sysIoctl, uintptr(fd),
					DKIOCREMOVABLE, uintptr(unsafe.Pointer(&rmdsk)))
				if err != 0 {
					fmt.Println(err.Error())
					return
				}
				if rmdsk == 0 {
					var media dk_minfo
					_, _, err = syscall.Syscall(sysIoctl, uintptr(fd),
						DKIOCGMEDIAINFO, uintptr(unsafe.Pointer(&media)))
					if err != 0 {
						fmt.Println(err.Error())
						return
					}
					value := media.dki_capacity * C.longlong(media.dki_lbsize) / GB
					dname := d.Name()[:len(d.Name())-2]
					disk := fmt.Sprintf("%s %.2fGB", dname, float64(value))
					fmt.Println(disk)
				}
			}
			syscall.Close(fd)
		}
	}
}
