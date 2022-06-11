package ioutils

import (
	"fmt"
	"sync"

	"github.com/muesli/termenv"
)

var Cout *Print = NewPrint()

type Print struct {
	m            *sync.Mutex
	Verbose      bool
	DisableColor bool
}

// Header : Print as a header
func (p *Print) Header(format string, a ...any) {

	if p.DisableColor {
		p.Printf(format, a...)
		return
	}

	s := termenv.String(fmt.Sprintf(format, a...))

	s = s.Bold().Foreground(Orange)

	p.Printf("%v", s)

}

func (p *Print) Value(format string, a ...any) {
	if p.DisableColor {
		p.Printf(format, a...)
		return
	}

	s := termenv.String(fmt.Sprintf(format, a...))

	s = s.Foreground(AquaMarine)

	p.Printf("%v", s)
}

func (p *Print) Seperator(len int) {

	tmp := ""

	for i := 0; i < len; i++ {
		tmp += "-"
	}

	if p.DisableColor {
		p.m.Lock()
		fmt.Printf("\n%v\n%v\n", tmp, tmp)
		p.m.Unlock()
		return
	}

	s := termenv.String(fmt.Sprintf("\n%v\n%v\n", tmp, tmp))

	s = s.Bold().Foreground(Azure)

	p.m.Lock()
	fmt.Printf("%v\n", s)
	p.m.Unlock()
}

func (p *Print) GetColor(z termenv.Color, format string, a ...any) termenv.Style {
	s := termenv.String(fmt.Sprintf(format, a...))
	if p.DisableColor {
		return s
	}

	s = s.Foreground(z)

	return s
}

func (p *Print) PrintColor(z termenv.Color, format string, a ...any) {
	if p.DisableColor {
		p.Printf(format, a...)
		return
	}

	s := termenv.String(fmt.Sprintf(format, a...))

	s = s.Bold().Foreground(z)

	p.Printf("%v", s)
}

func (p *Print) PrintInfo(format string, a ...any) {
	if !p.Verbose {
		return
	}

	if p.DisableColor {
		p.m.Lock()
		fmt.Printf("[Info] "+format+"\n", a...)
		p.m.Unlock()
	} else {
		z := fmt.Sprintf("%v %v\n", p.GetColor(Orange, "[Warn]"), p.GetColor(Azure, format, a...))
		p.m.Lock()
		fmt.Print(z)
		p.m.Unlock()
	}

}

func (p *Print) ErrColor(er error) termenv.Style {
	s := termenv.String(er.Error())
	if p.DisableColor {
		return s
	}

	s = s.Foreground(Red)

	return s
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
	if p.DisableColor {
		p.m.Lock()
		fmt.Printf("[Warn] "+format+"\n", a...)
		p.m.Unlock()
	} else {
		z := fmt.Sprintf("%v %v\n", p.GetColor(Orange, "[Warn]"), p.GetColor(LightGreen, format, a...))
		p.m.Lock()
		fmt.Print(z)
		p.m.Unlock()
	}
}

func NewPrint() *Print {
	x := Print{
		m: &sync.Mutex{},
	}

	return &x
}
