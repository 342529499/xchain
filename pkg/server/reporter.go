package server

import (
	"log"
	"os"
	"time"
)

//运行时打印运行时连接信息
var (
	printer_NetManager_Server printerType = 0
	printer_NetManager_Client printerType = 1
	printer_EndPoint          printerType = 2
	printer_All               printerType = 3

	defaultPrinterTimer = time.Second * 60
	printerHelper       = log.New(os.Stderr, "printer:", log.LstdFlags)
)

type printerType int

type TimerPrinter struct {
	timerCh   <-chan time.Time
	triggerCh chan printerType
	printer   printer
}

func (n *Node) StartPrinter(timer time.Duration) {
	ti := &TimerPrinter{
		timerCh:   time.Tick(timer),
		triggerCh: make(chan printerType, 30),
		printer:   n,
	}
	ti.start()
}

func (p *TimerPrinter) start() {
	for {
		select {
		case <-p.timerCh:
			printerHelper.Printf("printer time: %s\n", time.Now().String())
			p.printer.print(printer_All)
		case printerTYpe := <-p.triggerCh:
			printerHelper.Printf("printer time: %s\n", time.Now().String())
			p.printer.print(printerTYpe)
		}
	}
}

type printer interface {
	print(t printerType)
}

func (n *Node) print(t printerType) {

	switch t {
	case printer_NetManager_Server:
		printerHelper.Printf("server connection total: %d\n", len(n.netManager.serverConsManager.Keys()))
		if len(n.netManager.serverConsManager.Keys()) > 0 {
			printerHelper.Printf("details:[%s]\n", n.netManager.serverConsManager.Keys())
		}
	case printer_NetManager_Client:
		printerHelper.Printf("client connection total: %d\n", len(n.netManager.clientConsManager.Keys()))
		if len(n.netManager.clientConsManager.Keys()) > 0 {
			printerHelper.Printf("details:[%s]\n", n.netManager.clientConsManager.Keys())
		}
	case printer_EndPoint:
		n.epManager.printEP()
	case printer_All:
		n.print(printer_NetManager_Server)
		n.print(printer_NetManager_Client)
		n.print(printer_EndPoint)
	default:
		n.print(printer_All)
	}
}

func (n *EndPointManager) printEP() {
	printerHelper.Printf("node total: %d\n", len(n.NonValidateList)+len(n.ValidatorList))

	if len(n.ValidatorList) > 0 {
		printerHelper.Printf("non validate details:\n")
		for _, id := range n.NonValidateList {
			printerHelper.Printf("{id:%s,ip:%s}\n", id, n.IDToAddress[id])
		}
	}
	if len(n.NonValidateList) > 0 {
		printerHelper.Printf("validate details:\n")
		for _, id := range n.ValidatorList {
			printerHelper.Printf("{id:%s,ip:%s}\n", id, n.IDToAddress[id])
		}
	}
}
