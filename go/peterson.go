package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	NrOfProcesses = 2
	MinSteps      = 50
	MaxSteps      = 100
	MinDelay      = 10 * time.Millisecond
	MaxDelay      = 50 * time.Millisecond
)

type ProcessState int

const (
	LocalSection ProcessState = iota
	EntryProtocol
	CriticalSection
	ExitProtocol
)

var stateNames = []string{"LOCAL_SECTION", "ENTRY_PROTOCOL", "CRITICAL_SECTION", "EXIT_PROTOCOL"}

var C1 int = 0 
var C2 int = 0 
var Last int = 1 

const (
	BoardWidth  = NrOfProcesses
	BoardHeight = int(ExitProtocol) + 1
)

type Position struct {
	X int
	Y int
}

type Trace struct {
	TimeStamp time.Duration
	ID        int
	Position  Position
	Symbol    rune
}

type TracesSequence struct {
	Traces []Trace
}

type Process struct {
	ID       int
	Symbol   rune
	Position Position
}

var startTime time.Time
var printerChan = make(chan TracesSequence, NrOfProcesses)
var wg sync.WaitGroup

func printer() {
	defer wg.Done()
	
	for i := 0; i < NrOfProcesses; i++ {
		traces := <-printerChan
		printTraces(traces)
	}
	
	fmt.Printf("-1 %d %d %d ", NrOfProcesses, BoardWidth, BoardHeight)
	for i := 0; i < len(stateNames); i++ {
		fmt.Printf("%s;", stateNames[i])
	}
	fmt.Println("EXTRA_LABEL;")
}

func printTraces(traces TracesSequence) {
	for _, trace := range traces.Traces {
		printTrace(trace)
	}
}

func printTrace(trace Trace) {
	fmt.Printf("%.9f %d %d %d %c\n",
		trace.TimeStamp.Seconds(),
		trace.ID,
		trace.Position.X,
		trace.Position.Y,
		trace.Symbol)
}

func processTask(id int, symbol rune) {
	defer wg.Done()
	
	rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(id)))
	
	process := Process{
		ID:     id,
		Symbol: symbol,
		Position: Position{
			X: id,
			Y: int(LocalSection),
		},
	}
	
	nrOfSteps := MinSteps + rng.Intn(MaxSteps-MinSteps+1)
	
	var traces TracesSequence
	
	timeStamp := time.Since(startTime)
	traces.Traces = append(traces.Traces, Trace{
		TimeStamp: timeStamp,
		ID:        process.ID,
		Position:  process.Position,
		Symbol:    process.Symbol,
	})
	
	changeState := func(state ProcessState) {
		timeStamp := time.Since(startTime)
		process.Position.Y = int(state)
		traces.Traces = append(traces.Traces, Trace{
			TimeStamp: timeStamp,
			ID:        process.ID,
			Position:  process.Position,
			Symbol:    process.Symbol,
		})
	}
	
	for step := 0; step < nrOfSteps/4; step++ {
		// LOCAL_SECTION - start
		delay := MinDelay + time.Duration(float64(MaxDelay-MinDelay)*rng.Float64())
		time.Sleep(delay)
		// LOCAL_SECTION - end
		
		changeState(EntryProtocol) // starting ENTRY_PROTOCOL
		if process.ID == 0 {
			C1 = 1
			Last = 1
			for C2 == 1 && Last == 1 {
				time.Sleep(1 * time.Millisecond)
			}
		} else {
			C2 = 1
			Last = 2
			for C1 == 1 && Last == 2 {
				time.Sleep(1 * time.Millisecond)
			}
		}

		changeState(CriticalSection) // starting CRITICAL_SECTION

		// CRITICAL_SECTION - start
		delay = MinDelay + time.Duration(float64(MaxDelay-MinDelay)*rng.Float64())
		time.Sleep(delay)
		// CRITICAL_SECTION - end

		changeState(ExitProtocol) // starting EXIT_PROTOCOL
		if process.ID == 0 {
			C1 = 0
		} else {
			C2 = 0
		}

		changeState(LocalSection) // starting LOCAL_SECTION
	}
	
	printerChan <- traces
}

func main() {
	startTime = time.Now()
	
	wg.Add(1)
	go printer()
	
	symbol := 'A'
	for i := 0; i < NrOfProcesses; i++ {
		wg.Add(1)
		go processTask(i, symbol)
		symbol++
	}
	
	wg.Wait()
}