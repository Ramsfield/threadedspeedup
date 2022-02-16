package main

import (
  "fmt"
  "sync"
  "math/rand"
  "time"
  "sort"
  "runtime"
  "strconv"
  "os"
)

type ThreadInfo struct {
  sync.WaitGroup
  lock sync.Mutex
  array []int
  done bool
  shuffles int
}

func bogosort (tinfo *ThreadInfo) {
  defer tinfo.Done()
  localArray := make([]int, len(tinfo.array))
  copy(localArray, tinfo.array)
  numSorts := 0
  for (!sort.SliceIsSorted(localArray, func(p, q int) bool { return localArray[p] < localArray[q] }) && !tinfo.done) {
    rand.Shuffle(len(localArray), func(i, j int) { localArray[i], localArray[j] = localArray[j], localArray[i] })
    numSorts++
  }
  tinfo.lock.Lock()
  tinfo.shuffles += numSorts
  if !tinfo.done {
    copy(tinfo.array, localArray)
    tinfo.done = true
  }
  tinfo.lock.Unlock()
}

func usage() {
  fmt.Printf("Usage: %v [num threads] [size of array] [max number]", os.Args[0])
}

//Returns: num threads, size of slice, max number
func parseArgs() (int, int, int) {
  defaultThreads := runtime.NumCPU()
  defaultSliceSize := 10
  defaultMaxNum := 10000
  threads, size, max := defaultThreads, defaultSliceSize, defaultMaxNum
  var err error
  if len(os.Args) > 1 {
    threads, err = strconv.Atoi(os.Args[1])
    if err != nil {
      fmt.Printf("%v invalid, defaulting number of threads to %v\n", os.Args[1], defaultThreads)
      threads = defaultThreads
    }
  }
  if len(os.Args) > 2 {
    size, err = strconv.Atoi(os.Args[2])
    if err != nil {
      fmt.Printf("%v invalid, defaulting size of array to %v\n", os.Args[2], defaultSliceSize)
      size = defaultSliceSize
    }
  }
  if len(os.Args) > 3 {
    max, err = strconv.Atoi(os.Args[3])
    if err != nil {
      fmt.Printf("%v invalid, defaulting max number to %v\n", os.Args[3], defaultMaxNum)
      max = defaultMaxNum
    }
  }
  return threads, size, max
}

func main() {
  rand.Seed(time.Now().UnixNano())
  numthreads, ArraySize, maxnum := parseArgs()
  fmt.Printf("Running %v threads over an array size of %v. Max number: %v\n", numthreads, ArraySize, maxnum)
  tinfo := ThreadInfo { 
    array: make([]int, ArraySize),
    done: false,
  }
  for i:=0; i < ArraySize; i++ {
    tinfo.array[i] = rand.Int() % maxnum
  }
  tinfo.Add(numthreads)
  start := time.Now()
  for i:=0; i < numthreads; i++ {
    go bogosort(&tinfo)
  }
  tinfo.Wait()
  duration := time.Since(start)

  //Print Information
  fmt.Print(tinfo.shuffles)
  fmt.Print(" shuffles in ")
  fmt.Print(duration)
  fmt.Print(" ")
  fmt.Print(float64(tinfo.shuffles) / duration.Seconds())
  fmt.Print(" shuffles per second")
  fmt.Printf(" over %d threads\n", numthreads)
  fmt.Printf("Array: %v\n", tinfo.array)
}
