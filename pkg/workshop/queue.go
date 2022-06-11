package workshop

import (
	"io/ioutil"
	"sync"

	"github.com/tarunKoyalwar/talosplus/pkg/ioutils"
	"github.com/tarunKoyalwar/talosplus/pkg/shell"
)

// ExecQueue : Execute Given Queue with given concurrency
func ExecQueue(q []*shell.CMDWrap, limit int, showout bool) {

	wg := &sync.WaitGroup{}
	wrec := &sync.WaitGroup{}

	if limit > len(q) {
		limit = len(q)
	}

	var z chan *shell.CMDWrap = make(chan *shell.CMDWrap, len(q))
	var done chan *shell.CMDWrap = make(chan *shell.CMDWrap, limit)

	wrec.Add(1)
	go func(d chan *shell.CMDWrap) {
		defer wrec.Done()
		for {
			instance, ok := <-done
			if !ok {
				break
			}

			ioutils.Cout.Printf("[$] %v Executed Successfully", instance.Comment)

			// WIll show output of commands
			if showout {
				if !instance.Ignore {
					ioutils.Cout.Printf("[+] %v\n", instance.Raw)
					if instance.ExportFromFile == "" {
						ioutils.Cout.Printf("%v", instance.CMD.COutStream.String())
					} else {
						dat, _ := ioutil.ReadFile(instance.ExportFromFile)
						ioutils.Cout.Printf("%v", string(dat))
					}
				}
			}

		}
	}(done)

	worker := func(z chan *shell.CMDWrap) {
		defer wg.Done()

		for {
			d, ok := <-z
			if !ok {
				break
			}

			er := d.Execute()
			if er != nil {
				ioutils.Cout.PrintWarning(er.Error())
			}
			done <- d
		}

	}

	for i := 0; i < limit; i++ {
		wg.Add(1)
		go worker(z)
	}

	//assign work
	for _, v := range q {
		z <- v
	}

	close(z)
	wg.Wait()
	close(done)
	wrec.Wait()

}
