package week08

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/shirou/gopsutil/v3/mem"
)

func TestMem(t *testing.T) {
	InitClient()
	FlushAll()
	defer FlushAll()
	var (
		v           *mem.VirtualMemoryStat
		currentFree uint64
	)

	v, _ = mem.VirtualMemory()
	currentFree = v.Available
	string10 := randomStr(10)
	for i := 0; i < 100000; i++ {
		SetKey("10:"+strconv.Itoa(i), string10)
	}
	v, _ = mem.VirtualMemory()
	fmt.Printf("10 bytes value: %d bytes per key\n", (currentFree-v.Available)/100000)

	v, _ = mem.VirtualMemory()
	currentFree = v.Available
	string20 := randomStr(20)
	for i := 0; i < 100000; i++ {
		SetKey("20:"+strconv.Itoa(i), string20)
	}
	v, _ = mem.VirtualMemory()
	fmt.Printf("20 bytes value: %d bytes per key\n", (currentFree-v.Available)/100000)

	v, _ = mem.VirtualMemory()
	currentFree = v.Available
	string50 := randomStr(50)
	for i := 0; i < 100000; i++ {
		SetKey("50:"+strconv.Itoa(i), string50)
	}
	v, _ = mem.VirtualMemory()
	fmt.Printf("50 bytes value: %d bytes per key\n", (currentFree-v.Available)/100000)

	v, _ = mem.VirtualMemory()
	currentFree = v.Available
	string100 := randomStr(100)
	for i := 0; i < 100000; i++ {
		SetKey("100:"+strconv.Itoa(i), string100)
	}
	v, _ = mem.VirtualMemory()
	fmt.Printf("100 bytes value: %d bytes per key\n", (currentFree-v.Available)/100000)

	v, _ = mem.VirtualMemory()
	currentFree = v.Available
	string200 := randomStr(200)
	for i := 0; i < 100000; i++ {
		SetKey("200:"+strconv.Itoa(i), string200)
	}
	v, _ = mem.VirtualMemory()
	fmt.Printf("200 bytes value: %d bytes per key\n", (currentFree-v.Available)/100000)

	v, _ = mem.VirtualMemory()
	currentFree = v.Available
	string1k := randomStr(1024)
	for i := 0; i < 50000; i++ {
		SetKey("1k:"+strconv.Itoa(i), string1k)
	}
	v, _ = mem.VirtualMemory()
	fmt.Printf("1k bytes value: %d bytes per key\n", (currentFree-v.Available)/50000)

	v, _ = mem.VirtualMemory()
	currentFree = v.Available
	string5k := randomStr(5120)
	for i := 0; i < 10000; i++ {
		SetKey("1k:"+strconv.Itoa(i), string5k)
	}
	v, _ = mem.VirtualMemory()
	fmt.Printf("5k bytes value: %d bytes per key\n", (currentFree-v.Available)/10000)

}
