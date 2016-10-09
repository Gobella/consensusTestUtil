package systemStat

import (
	"github.com/shirou/gopsutil/mem"
	"fmt"
)

func test(){
	c,err:=mem.SwapMemory()
	if err!=nil{
		return
	}
	fmt.Println("swapMemory:",c)
}