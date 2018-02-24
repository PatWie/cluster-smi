package proc

// #include "proc.h"
// #include <stdlib.h>
import "C"
import "unsafe"
import "errors"
import "strconv"
import "io/ioutil"

// PIDInfo is a collection of properties for a process
type PIDInfo struct {
	PID       int
	Command   string
	UsedTime  int64
	StartTime int64
	UpTime    int64
}

// ClockTicks returns ticks per clock cycle
func ClockTicks() int64 {
	hz := C.long(0)
	C.clock_ticks(&hz)
	return int64(hz)
}

// TimeOfDay returns current time-stamp
func TimeOfDay() int64 {
	val := C.float(0)
	C.time_of_day(&val)
	return int64(val)
}

// BootTime returns timestamps when machine was started
func BootTime() (int64, error) {
	uptime := C.float(0)
	idletime := C.float(0)
	success := C.int(0)
	success = C.boot_time(&uptime, &idletime)

	if int(success) == 0 {
		return 0, errors.New("cannot get bootime")
	}

	return TimeOfDay() - int64(uptime), nil
}

func mallocCStringBuffer(size uint) *C.char {
	buf := make([]byte, size)
	return C.CString(string(buf))
}

// CPUTick returns the total number of Jiffies
func CPUTick() (t int64) {
	return int64(C.read_cpu_tick())
}

// InfoFromPid returns used time (int) and command (string)
func InfoFromPid(pid int) PIDInfo {

	info := PIDInfo{}

	usedTime := C.ulong(0)
	startTime := C.ulonglong(0)

	var cDst *C.char = mallocCStringBuffer(128 + 1)
	defer C.free(unsafe.Pointer(cDst))

	C.read_pid_info(C.ulong(pid), &usedTime, &startTime, cDst)

	info.PID = pid
	info.Command = C.GoString(cDst)
	info.UsedTime = int64(usedTime)
	info.StartTime = int64(startTime)

	return info
}

// NumberCPUCores returns number of CPU cores
func NumberCPUCores() (n int) {
	return int(C.num_cores())
}

// UIDFromPID returns user id (UID) for a given process id (PID)
func UIDFromPID(pid int) (uid int) {
	cUID := C.ulong(0)
	C.get_uid_from_pid(C.ulong(pid), &cUID)
	return int(cUID)
}

// CmdFromPID returns cmd which initiated the process
func CmdFromPID(pid int) string {
	fn := "/proc/" + strconv.Itoa(pid) + "/cmdline"
	b, err := ioutil.ReadFile(fn) // just pass the file name
	if err != nil {
		return "unkown"
	}

	return string(b)
	// var cCMD *C.char = mallocCStringBuffer(128 + 1)
	// defer C.free(unsafe.Pointer(cCMD))

	// C.get_cmd(C.ulong(pid), cCMD)

	// return C.GoString(cCMD)
}

// GetRAMMemoryInfo returns memory information of RAM
func GetRAMMemoryInfo() (total int64, free int64, available int64) {
	cTotal := C.ulong(0)
	cFree := C.ulong(0)
	cAvailable := C.ulong(0)
	C.get_mem(&cTotal, &cFree, &cAvailable)

	return int64(cTotal), int64(cFree), int64(cAvailable)
}
