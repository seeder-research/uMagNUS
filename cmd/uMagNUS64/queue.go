package main

// File que for distributing multiple input files over GPUs.

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	engine "github.com/seeder-research/uMagNUS/engine64"
	opencl "github.com/seeder-research/uMagNUS/opencl64"
)

var (
	exitStatus       atom = 0
	numOK, numFailed atom = 0, 0
)

func RunQueue(files []string) {
	s := NewStateTab(files)
	s.PrintTo(os.Stdout)
	go s.ListenAndServe(*engine.Flag_port)
	s.Run()
	fmt.Println(numOK.get(), "OK, ", numFailed.get(), "failed")
	os.Exit(int(exitStatus))
}

// StateTab holds the queue state (list of jobs + statuses).
// All operations are atomic.
type stateTab struct {
	lock sync.Mutex
	jobs []job
	next int
}

// Job info.
type job struct {
	inFile  string // input file to run
	webAddr string // http address for gui of running process
	uid     int
}

// NewStateTab constructs a queue for the given input files.
// After construction, it is accessed atomically.
func NewStateTab(inFiles []string) *stateTab {
	s := new(stateTab)
	s.jobs = make([]job, len(inFiles))
	for i, f := range inFiles {
		s.jobs[i] = job{inFile: f, uid: i}
	}
	return s
}

// StartNext advances the next job and marks it running, setting its webAddr to indicate the GUI url.
// A copy of the job info is returned, the original remains unmodified.
// ok is false if there is no next job.
func (s *stateTab) StartNext(webAddr string) (next job, ok bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.next >= len(s.jobs) {
		return job{}, false
	}
	s.jobs[s.next].webAddr = webAddr
	jobCopy := s.jobs[s.next]
	s.next++
	return jobCopy, true
}

// Finish marks the job with j's uid as finished.
func (s *stateTab) Finish(j job) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.jobs[j.uid].webAddr = ""
}

// Runs all the jobs in stateTab.
func (s *stateTab) Run() {
	var gpu_arr []int
	numGPUsAvailable := len(opencl.ClDevices)
	useAllGPUs := false
	if len(*engine.Flag_gpulist) > 0 {
		gpuList := strings.Split(*engine.Flag_gpulist, ",")
		if len(gpuList) > 0 {
			for _, item := range gpuList {
				if id, err := strconv.Atoi(item); err == nil {
					if id < 0 {
						log.Println("Invalid GPU id detected. Must be an integer >= 0!")
					} else {
						if id < numGPUsAvailable {
							gpu_arr = append(gpu_arr, id)
						} else {
							log.Printf("Invalid GPU id detected. Must be an integer < %v!\n", numGPUsAvailable)
						}
					}
				}
			}
			if len(gpu_arr) == 0 {
				useAllGPUs = true
			}
		} else {
			log.Println("Empty GPU list given. Will use all detected GPUs instead.")
			useAllGPUs = true
		}
	} else {
		useAllGPUs = true
	}
	if useAllGPUs {
		for idx := 0; idx < numGPUsAvailable; idx++ {
			gpu_arr = append(gpu_arr, idx)
		}
	}
	nGPU := len(gpu_arr)
	idle := initGPUs(nGPU)
	for {
		var gpuIdx int
		var gpu int
		var addr string
		if *engine.Flag_host {
			gpu = -5
			addr = fmt.Sprint(":", 35368)
		} else {
			gpuIdx = <-idle
			gpu = gpu_arr[gpuIdx]
			addr = fmt.Sprint(":", 35368+gpuIdx)
		}
		j, ok := s.StartNext(addr)
		if !ok {
			break
		}
		go func() {
			run(j.inFile, gpu, j.webAddr)
			s.Finish(j)
			if !(*engine.Flag_host) {
				idle <- gpuIdx
			}
		}()
	}
	// drain remaining tasks (one already done)
	if !(*engine.Flag_host) {
		for i := 1; i < nGPU; i++ {
			<-idle
		}
	}
}

type atom int32

func (a *atom) set(v int) { atomic.StoreInt32((*int32)(a), int32(v)) }
func (a *atom) get() int  { return int(atomic.LoadInt32((*int32)(a))) }
func (a *atom) inc()      { atomic.AddInt32((*int32)(a), 1) }

func run(inFile string, gpu int, webAddr string) {
	// overridden flags
	gpuFlag := fmt.Sprint(`-gpu=`, gpu)
	httpFlag := fmt.Sprint(`-http=`, webAddr)

	// pass through flags
	flags := []string{gpuFlag, httpFlag}
	flag.Visit(func(f *flag.Flag) {
		if f.Name != "gpu" && f.Name != "http" && f.Name != "failfast" {
			flags = append(flags, fmt.Sprintf("-%v=%v", f.Name, f.Value))
		}
	})
	flags = append(flags, inFile)

	cmd := exec.Command(os.Args[0], flags...)
	log.Println(os.Args[0], flags)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(inFile, err)
		log.Printf("%s\n", output)
		exitStatus.set(1)
		numFailed.inc()
		if *engine.Flag_failfast {
			os.Exit(1)
		}
	} else {
		numOK.inc()
	}
}

func initGPUs(nGpu int) chan int {
	if nGpu == 0 {
		log.Fatal("no GPUs available")
		panic(0)
	}
	idle := make(chan int, nGpu)
	for i := 0; i < nGpu; i++ {
		idle <- i
	}
	return idle
}

func (s *stateTab) PrintTo(w io.Writer) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for i, j := range s.jobs {
		fmt.Fprintf(w, "%3d %v %v\n", i, j.inFile, j.webAddr)
	}
}

func (s *stateTab) RenderHTML(w io.Writer) {
	s.lock.Lock()
	defer s.lock.Unlock()
	fmt.Fprintln(w, ` 
<!DOCTYPE html> <html> <head> 
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	<meta http-equiv="refresh" content="1">
`+engine.CSS+`
	</head><body>
	<span style="color:gray; font-weight:bold; font-size:1.5em"> uMagNUS queue status </span><br/>
	<hr/>
	<pre>
`)

	hostname := "localhost"
	hostname, _ = os.Hostname()
	for _, j := range s.jobs {
		if j.webAddr != "" {
			fmt.Fprint(w, `<b>`, j.uid, ` <a href="`, "http://", hostname+j.webAddr, `">`, j.inFile, " ", j.webAddr, "</a></b>\n")
		} else {
			fmt.Fprint(w, j.uid, " ", j.inFile, "\n")
		}
	}

	fmt.Fprintln(w, `</pre><hr/></body></html>`)
}

func (s *stateTab) ListenAndServe(addr string) {
	http.Handle("/", s)
	go http.ListenAndServe(addr, nil)
}

func (s *stateTab) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.RenderHTML(w)
}
