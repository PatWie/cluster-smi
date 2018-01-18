package proc

// #include "proc.h"
// #include <stdlib.h>
import "C"
import "unsafe"
import "errors"

type PIDInfo struct {
	PID       int
	Command   string
	UsedTime  int64
	StartTime int64
	UpTime    int64
}

func ClockTicks() int64 {
	hz := C.long(0)
	C.clock_ticks(&hz)
	return int64(hz)
}

func TimeOfDay() int64 {
	val := C.float(0)
	C.time_of_day(&val)
	return int64(val)
}

// returns timestamps when machine was started
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

// CpuTick returns the total number of Jiffies
func CpuTick() (t int64) {
	return int64(C.read_cpu_tick())
}

// TimeAndNameFromPID returns used time (int) and command (string)
func InfoFromPid(pid int) PIDInfo {

	info := PIDInfo{}

	used_time := C.ulong(0)
	starttime := C.ulonglong(0)

	var c_dst *C.char = mallocCStringBuffer(128 + 1)
	defer C.free(unsafe.Pointer(c_dst))

	C.read_pid_info(C.ulong(pid), &used_time, &starttime, c_dst)

	info.PID = pid
	info.Command = C.GoString(c_dst)
	info.UsedTime = int64(used_time)
	info.StartTime = int64(starttime)

	return info
}

// NumberCPUCores returns number of CPU cores
func NumberCPUCores() (n int) {
	return int(C.num_cores())
}

// Find user id (UID) for a given process id (PID)
func UIDFromPID(pid int) (uid int) {
	c_uid := C.ulong(0)
	C.get_uid_from_pid(C.ulong(pid), &c_uid)
	return int(c_uid)
}

// get memory information of RAM
func GetRAMMemoryInfo() (total int64, free int64, available int64) {
	c_total := C.ulong(0)
	c_free := C.ulong(0)
	c_available := C.ulong(0)

	C.get_mem(&c_total, &c_free, &c_available)

	return int64(c_total), int64(c_free), int64(c_available)
}
