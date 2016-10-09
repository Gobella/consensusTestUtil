package statistics

import (
	"testing"
	"fmt"
)

//func TestLogger(t *testing.T){
//	ls:=NewLogStatic("test for test!")
//
//	go send(ls,"tag111")
//
//	go send(ls,"tag222")
//
//	time.Sleep(3*time.Second)
//
//	ls.PrintRecord()
//}
//
//func send(ls *LogStatic,s string){
//	for i:=0;i<1000;i++{
//		ls.RecordCount("test",s,"A","B","C")
//		ls.RecordCount("test",s,"summary")
//	}
//}

//func TestFun(t *testing.T){
//	c:=&NodeBasicInfo {
//		Id :"1",
//		//Addr:nil,
//	}
//	ls:=NewLogStatic(c)
//
//	testMutiArg(ls,"aaaa","nnnn")
//	ls.PrintRecord()
//}
//func testMutiArg(ls *LogStatic,tags ...string){
//	//ls.RecordCount("cc",tags)
//
//}

func TestInfo(t *testing.T){
	c:=&NodeBasicInfo {
		Id :"1",
		//Addr:nil,
	}
	ls:=NewLogStatic(c)
	ls.RecordList("hello","")
	ls.RecordInfo("qwww","test")
	ls.PrintRecord(0)
	ls.RecordInfo("qwfdsfsdfww","test")

	ls.RecordList("hello1","")
	ls.PrintRecord(0)
}

func TestSs(t *testing.T){
	n:=uint64(30)
	K:=uint64(10)
	h := n / K * K
	fmt.Println(h)
}
