package main

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	ps "github.com/shirou/gopsutil/process"
)

var (
	Processes []string // list of all processes by name
)

// entry point
func main() {
	// tview application
	app := tview.NewApplication()

	// tview table
	pTable := tview.NewTable().Select(0, 0).SetSelectable(true, false).SetFixed(2, 1)

	// get processes
	processes, err := ps.Processes()
	if err != nil {
		log.Fatal(err)
	}

	// set table headers
	pTable.SetCell(0, 0, tview.NewTableCell(" Process ID "))
	pTable.SetCell(1, 0, tview.NewTableCell("============"))
	pTable.SetCell(0, 1, tview.NewTableCell(" Parent PID "))
	pTable.SetCell(1, 1, tview.NewTableCell("============"))
	pTable.SetCell(0, 2, tview.NewTableCell(" Process "))
	pTable.SetCell(1, 2, tview.NewTableCell("========="))
	pTable.SetCell(0, 3, tview.NewTableCell(" User "))
	pTable.SetCell(1, 3, tview.NewTableCell("======"))
	pTable.SetCell(0, 4, tview.NewTableCell(" CPU "))
	pTable.SetCell(1, 4, tview.NewTableCell("====="))
	pTable.SetCell(0, 5, tview.NewTableCell(" Memory "))
	pTable.SetCell(1, 5, tview.NewTableCell("========"))

	// add processes to table
	for i, p := range processes {
		j := i + 2
		pid := p.Pid
		if pid == 0 {
			continue
		}

		parentProc, err := p.Parent()
		var ppid int32
		if err == nil {
			ppid = parentProc.Pid
		} else {
			ppid = 0
		}
		user, _ := p.Username()
		exe, _ := p.Name()
		cpu, _ := p.CPUPercent()
		mem, _ := p.MemoryPercent()

		pTable.SetCell(j, 0, tview.NewTableCell(fmt.Sprint(pid)))
		pTable.SetCell(j, 1, tview.NewTableCell(fmt.Sprint(ppid)))
		pTable.SetCell(j, 2, tview.NewTableCell(exe))
		pTable.SetCell(j, 3, tview.NewTableCell(user))
		pTable.SetCell(j, 4, tview.NewTableCell(fmt.Sprintf("%.2f%%", cpu)))
		pTable.SetCell(j, 5, tview.NewTableCell(fmt.Sprintf("%.2f%%", mem)))

		Processes = append(Processes, exe)
	}

	// give process list border
	pTable.SetBorder(true).SetTitle("Processes")

	// custom key captures
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlQ {
			app.Stop()
			return nil
		}

		return event
	})

	// run app
	if err := app.SetRoot(pTable, true).SetFocus(pTable).Run(); err != nil {
		log.Fatal(err)
	}
}
