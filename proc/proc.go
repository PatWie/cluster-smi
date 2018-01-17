package proc

// #include "proc.h"
// #include <stdlib.h>
import "C"
import "unsafe"

func mallocCStringBuffer(size uint) *C.char {
	buf := make([]byte, size)
	return C.CString(string(buf))
}

// CpuTick returns the total number of Jiffies
func CpuTick() (t int64) {
	return int64(C.read_cpu_tick())
}

// TimeAndNameFromPID returns used time (int) and command (string)
func TimeAndNameFromPID(pid int) (int64, string) {
	time := C.ulong(0)

	var c_dst *C.char = mallocCStringBuffer(128 + 1)
	defer C.free(unsafe.Pointer(c_dst))

	C.read_time_and_name_from_pid(C.ulong(pid), &time, c_dst)
	return int64(time), C.GoString(c_dst)
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
