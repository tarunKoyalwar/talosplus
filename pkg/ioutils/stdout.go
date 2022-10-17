package ioutils

import (
	"fmt"
	"os"
	"sync"

	"github.com/muesli/termenv"
)

// Cout : Global Print Instance
var Cout *Print = NewPrint()

// Print : Concurrency Safe Stdout Printing
type Print struct {
	m            *sync.Mutex
	Verbose      bool // Verbose (includes warnings)
	VeryVerbose  bool // VeryVerbose (includes warnings+info)
	DisableColor bool // Disable Color Output
}

// Header : Stdout Header Printing Style
func (p *Print) Header(format string, a ...any) {

	if p.DisableColor {
		p.Printf(format, a...)
		return
	}

	s := termenv.String(fmt.Sprintf(format, a...))

	s = s.Bold().Foreground(Orange)

	p.Printf("%v", s)

}

// Value : Stdout Header Value Printing Style
func (p *Print) Value(format string, a ...any) {
	if p.DisableColor {
		p.Printf(format, a...)
		return
	}

	s := termenv.String(fmt.Sprintf(format, a...))

	s = s.Foreground(AquaMarine)

	p.Printf("%v", s)
}

// Seperator : Dash(-) Seperator of given length
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

// GetColor : Get Colorized text
func (p *Print) GetColor(z termenv.Color, format string, a ...any) termenv.Style {
	s := termenv.String(fmt.Sprintf(format, a...))
	if p.DisableColor {
		return s
	}

	s = s.Foreground(z)

	return s
}

// PrintColor : Print Colored Output
func (p *Print) PrintColor(z termenv.Color, format string, a ...any) {
	if p.DisableColor {
		p.Printf(format, a...)
		return
	}

	s := termenv.String(fmt.Sprintf(format, a...))

	s = s.Bold().Foreground(z)

	p.Printf("%v", s)
}

// PrintInfo : Print Info
func (p *Print) PrintInfo(format string, a ...any) {
	if !p.VeryVerbose {
		return
	}

	if p.DisableColor {
		p.m.Lock()
		fmt.Printf("[Info] "+format+"\n", a...)
		p.m.Unlock()
	} else {
		z := fmt.Sprintf("%v %v\n", p.GetColor(Orange, "[Info]"), p.GetColor(Azure, format, a...))
		p.m.Lock()
		fmt.Print(z)
		p.m.Unlock()
	}

}

// Fatalf : Output Followed by panic
func (p *Print) Fatalf(er error, format string, a ...any) {
	if p.DisableColor {
		p.m.Lock()
		fmt.Printf("[Fatal] "+format+"\n", a...)
		p.m.Unlock()
	} else {
		z := fmt.Sprintf("%v %v\n", p.GetColor(Red, "[Fatal]"), fmt.Sprintf(format, a...))
		p.m.Lock()
		fmt.Print(z)
		p.m.Unlock()
	}

	panic(er)

}

// ErrExit : Error Followed by exit
func (p *Print) ErrExit(format string, a ...any) {
	if p.DisableColor {
		p.m.Lock()
		fmt.Printf("[Fatal] "+format+"\n", a...)
		p.m.Unlock()
	} else {
		z := fmt.Sprintf("%v %v\n", p.GetColor(Red, "[Fatal]"), fmt.Sprintf(format, a...))
		p.m.Lock()
		fmt.Print(z)
		p.m.Unlock()
	}

	os.Exit(1)
}

// ErrColor : Colorize Error
func (p *Print) ErrColor(er error) termenv.Style {
	s := termenv.String(er.Error())
	if p.DisableColor {
		return s
	}

	s = s.Foreground(Red)

	return s
}

// Printf : Normal Print
func (p *Print) Printf(format string, a ...any) {
	p.m.Lock()
	fmt.Printf(format+"\n", a...)
	p.m.Unlock()
}

// PrintWarning : Print Warnings
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

// NewPrint : New Print Instance
func NewPrint() *Print {
	x := Print{
		m: &sync.Mutex{},
	}

	return &x
}
