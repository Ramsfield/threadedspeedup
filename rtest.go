package main

import (
  "fmt"
  "sync"
  "time"
)

type ThreadInfo struct {
  sync.WaitGroup
  sync.Mutex
  actionsComplete int64
  actionsPerThread int
}

func actionCaller (tinfo *ThreadInfo) {
  defer tinfo.Done()
  num := 0
  numActions := 0
  totalActions := tinfo.actionsPerThread
  for i := 0; i < totalActions; i++ {
    num += i
    numActions++
  }

  tinfo.Lock()
  defer tinfo.Unlock()
  tinfo.actionsComplete += int64(numActions)
}

func main() {
  numthreads := 2
  tinfo := ThreadInfo { 
    actionsComplete: 0,
    actionsPerThread: 100000000000,
  }
  tinfo.Add(numthreads)
  start := time.Now()
  for i:=0; i < numthreads; i++ {
    go actionCaller(&tinfo)
  }
  tinfo.Wait()
  duration := time.Since(start)

  //Print Information
  fmt.Printf("%v actions in %v: %v actions per second over %v threads\n", tinfo.actionsComplete, duration, float64(tinfo.actionsComplete) / duration.Seconds(), numthreads)
}
