package week08

import (
	"fmt"
	"testing"
)

/*


 */

func TestMem(t *testing.T) {
	InitClient()
	defer FlushAll()
	var (
		curMem, preMem int64
	)

	FlushAll()
	curMem, _ = GetMemInfo()
	string10 := randomStr(10)
	for i := 0; i < 100000; i++ {
		SetKey(fmt.Sprintf("%05d", i), string10)
	}
	preMem = curMem
	curMem, _ = GetMemInfo()
	fmt.Printf("10 bytes: %f byte\n", float64(curMem-preMem)/100000)

	FlushAll()
	curMem, _ = GetMemInfo()
	string20 := randomStr(20)
	for i := 0; i < 100000; i++ {
		SetKey(fmt.Sprintf("%05d", i), string20)
	}
	preMem = curMem
	curMem, _ = GetMemInfo()
	fmt.Printf("20 bytes: %f byte\n", float64(curMem-preMem)/100000)

	FlushAll()
	curMem, _ = GetMemInfo()
	string50 := randomStr(50)
	for i := 0; i < 100000; i++ {
		SetKey(fmt.Sprintf("%05d", i), string50)
	}
	preMem = curMem
	curMem, _ = GetMemInfo()
	fmt.Printf("50 bytes: %f byte\n", float64(curMem-preMem)/100000)

	FlushAll()
	curMem, _ = GetMemInfo()
	string100 := randomStr(100)
	for i := 0; i < 100000; i++ {
		SetKey(fmt.Sprintf("%05d", i), string100)
	}
	preMem = curMem
	curMem, _ = GetMemInfo()
	fmt.Printf("100 bytes: %f byte\n", float64(curMem-preMem)/100000)

	FlushAll()
	curMem, _ = GetMemInfo()
	string200 := randomStr(200)
	for i := 0; i < 100000; i++ {
		SetKey(fmt.Sprintf("%05d", i), string200)
	}
	preMem = curMem
	curMem, _ = GetMemInfo()
	fmt.Printf("200 bytes: %f byte\n", float64(curMem-preMem)/100000)

	FlushAll()
	curMem, _ = GetMemInfo()
	string1k := randomStr(1024)
	for i := 0; i < 50000; i++ {
		SetKey(fmt.Sprintf("%05d", i), string1k)
	}
	preMem = curMem
	curMem, _ = GetMemInfo()
	fmt.Printf("1024 bytes: %f byte\n", float64(curMem-preMem)/50000)

	FlushAll()
	curMem, _ = GetMemInfo()
	string5k := randomStr(5120)
	for i := 0; i < 10000; i++ {
		SetKey(fmt.Sprintf("%05d", i), string5k)
	}
	preMem = curMem
	curMem, _ = GetMemInfo()
	fmt.Printf("5120 bytes: %f byte\n", float64(curMem-preMem)/10000)
}

func TestInfo(t *testing.T) {
	mem, _ := GetMemInfo()
	fmt.Println(mem)
	fmt.Printf("%05d\n", 3)
}
