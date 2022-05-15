package ioutils

import (
	"fmt"
	"sync"
)

var Cout *Print = NewPrint()

type Print struct {
	m       *sync.Mutex
	Verbose bool
}

func (p *Print) PrintInfo(format string, a ...any) {
	if !p.Verbose {
		return
	}
	p.m.Lock()
	fmt.Printf("[Info] "+format+"\n", a...)
	p.m.Unlock()
}

func (p *Print) Printf(format string, a ...any) {
	p.m.Lock()
	fmt.Printf(format+"\n", a...)
	p.m.Unlock()
}

func (p *Print) PrintWarning(format string, a ...any) {
	if !p.Verbose {
		return
	}
	p.m.Lock()
	fmt.Printf("[Warn] "+format+"\n", a...)
	p.m.Unlock()
}

func (p *Print) DrawLine(len int) {
	tmp := ""

	for i := 0; i < len; i++ {
		tmp += "-"
	}

	p.m.Lock()
	fmt.Printf("\n%v\n", tmp)
	p.m.Unlock()
}

func NewPrint() *Print {
	x := Print{
		m: &sync.Mutex{},
	}

	return &x
}
