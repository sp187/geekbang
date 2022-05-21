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
	for i := 0; i < 100000; i++ {
		SetKey(fmt.Sprintf("%05d", i), randomStr(10))
	}
	preMem = curMem
	curMem, _ = GetMemInfo()
	fmt.Printf("10 bytes: %f byte\n", float64(curMem-preMem)/100000)

	FlushAll()
	curMem, _ = GetMemInfo()
	for i := 0; i < 100000; i++ {
		SetKey(fmt.Sprintf("%05d", i), randomStr(20))
	}
	preMem = curMem
	curMem, _ = GetMemInfo()
	fmt.Printf("20 bytes: %f byte\n", float64(curMem-preMem)/100000)

	FlushAll()
	curMem, _ = GetMemInfo()
	for i := 0; i < 100000; i++ {
		SetKey(fmt.Sprintf("%05d", i), randomStr(50))
	}
	preMem = curMem
	curMem, _ = GetMemInfo()
	fmt.Printf("50 bytes: %f byte\n", float64(curMem-preMem)/100000)

	FlushAll()
	curMem, _ = GetMemInfo()
	for i := 0; i < 100000; i++ {
		SetKey(fmt.Sprintf("%05d", i), randomStr(100))
	}
	preMem = curMem
	curMem, _ = GetMemInfo()
	fmt.Printf("100 bytes: %f byte\n", float64(curMem-preMem)/100000)

	FlushAll()
	curMem, _ = GetMemInfo()
	for i := 0; i < 100000; i++ {
		SetKey(fmt.Sprintf("%05d", i), randomStr(200))
	}
	preMem = curMem
	curMem, _ = GetMemInfo()
	fmt.Printf("200 bytes: %f byte\n", float64(curMem-preMem)/100000)

	FlushAll()
	curMem, _ = GetMemInfo()
	for i := 0; i < 50000; i++ {
		SetKey(fmt.Sprintf("%05d", i), randomStr(1024))
	}
	preMem = curMem
	curMem, _ = GetMemInfo()
	fmt.Printf("1024 bytes: %f byte\n", float64(curMem-preMem)/50000)

	FlushAll()
	curMem, _ = GetMemInfo()
	for i := 0; i < 10000; i++ {
		SetKey(fmt.Sprintf("%05d", i), randomStr(5120))
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
