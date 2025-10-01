package performance

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/sirupsen/logrus"
)

func Profiling(logger *logrus.Logger) {
	fmt.Println("Run performance profiling")
	fCpu, err := os.Create("cpu-perf")
	if err != nil {
		logger.Fatal("could not create CPU profile: ", err)
	}
	if err := pprof.StartCPUProfile(fCpu); err != nil {
		logger.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	fMem, err := os.Create("mem-perf")
	if err != nil {
		logger.Fatal("could not create memory profile: ", err)
	}
	defer func(fMem *os.File) {
		err := fMem.Close()
		if err != nil {
			logger.Fatal("Unable to close file: ", err)
		}
	}(fMem)
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(fMem); err != nil {
		logger.Fatal("could not write memory profile: ", err)
	}
}
