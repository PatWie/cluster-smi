package nvml

/*
#cgo CFLAGS: -I/usr/local/cuda/include
#cgo LDFLAGS: -lnvidia-ml -L/usr/local/cuda-8.0/targets/x86_64-linux/lib/stubs/

// #cgo CFLAGS: -I/graphics/opt/opt_Ubuntu16.04/cuda/toolkit_8.0/cuda/include
// #cgo LDFLAGS: -lnvidia-ml -L/graphics/opt/opt_Ubuntu16.04/cuda/toolkit_8.0/cuda/lib64/stubs

#include "bridge.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

var (
	errNoErrorString = errors.New("nvml: expected an error from driver but got nothing")
	errNoError       = errors.New("nvml: getGoError called on a successful API call")
)

func getGoError(result C.nvmlReturn_t) error {
	if result == C.NVML_SUCCESS {
		return errNoError
	}

	errString := C.nvmlErrorString(result)
	if errString == nil {
		return errNoErrorString
	}

	return fmt.Errorf("nvml: %s", C.GoString(errString))
}

// Device describes the NVIDIA GPU device attached to the host
type Device struct {
	DeviceName string
	DeviceUUID string
	d          C.nvmlDevice_t
	i          int
}

// NVMLMemory contains information about the memory allocation of a device
type NVMLMemory struct {
	// Unallocated FB memory (in bytes).
	Free int64
	// Total installed FB memory (in bytes).
	Total int64
	// Allocated FB memory (in bytes). Note that the driver/GPU always sets
	// aside a small amount of memory for bookkeeping.
	Used int64
}

type NVMLProcess struct {
	Pid           int
	UsedGpuMemory int64
	// Procs         []C.nvmlProcessInfo_t
}

func newDevice(nvmlDevice C.nvmlDevice_t, idx int) (dev Device, err error) {
	dev = Device{
		d: nvmlDevice,
		i: idx,
	}

	if dev.DeviceUUID, err = dev.UUID(); err != nil {
		return
	}

	if dev.DeviceName, err = dev.Name(); err != nil {
		return
	}
	return
}

func (s *Device) callGetTextFunc(f C.getNvmlCharProperty, sz C.uint) (string, error) {
	buf := make([]byte, sz)
	cs := C.CString(string(buf))
	defer C.free(unsafe.Pointer(cs))

	if result := C.bridge_get_text_property(f, s.d, cs, sz); result != C.NVML_SUCCESS {
		return "", getGoError(result)
	}

	return C.GoString(cs), nil
}

func (s *Device) callGetIntFunc(f C.getNvmlIntProperty) (int, error) {
	var valC C.uint

	if result := C.bridge_get_int_property(f, s.d, &valC); result != C.NVML_SUCCESS {
		return 0, getGoError(result)
	}

	return int(valC), nil
}

// UUID returns the Device's Unique ID
func (s *Device) UUID() (uuid string, err error) {
	return s.callGetTextFunc(C.getNvmlCharProperty(C.nvmlDeviceGetUUID), C.NVML_DEVICE_UUID_BUFFER_SIZE)
}

// Name returns the Device's Name and is not guaranteed to exceed 64 characters in length
func (s *Device) Name() (name string, err error) {
	return s.callGetTextFunc(C.getNvmlCharProperty(C.nvmlDeviceGetName), C.NVML_DEVICE_NAME_BUFFER_SIZE)
}

// GetUtilization returns the GPU and memory usage returned as a percentage used of a given GPU device
func (s *Device) GetUtilization() (gpu, memory int, err error) {
	var utilRates C.nvmlUtilization_t
	if result := C.nvmlDeviceGetUtilizationRates(s.d, &utilRates); result != C.NVML_SUCCESS {
		err = getGoError(result)
		return
	}
	gpu = int(utilRates.gpu)
	memory = int(utilRates.memory)
	return
}

// GetPowerUsage returns the power consumption of the GPU in watts
func (s *Device) GetPowerUsage() (int, error) {
	usage, err := s.callGetIntFunc(C.getNvmlIntProperty(C.nvmlDeviceGetPowerUsage))
	if err != nil {
		return 0, err
	}
	// nvmlDeviceGetPowerUsage returns milliwatts.. convert to watts
	return usage / 1000, nil
}

// GetFanSpeed returns the fan speed in percent
func (s *Device) GetFanSpeed() (int, error) {
	speed, err := s.callGetIntFunc(C.getNvmlIntProperty(C.nvmlDeviceGetFanSpeed))
	if err != nil {
		return 0, err
	}
	return speed, nil
}

// GetTemperature returns the Device's temperature in Farenheit and celsius
func (s *Device) GetTemperature() (int, int, error) {
	var tempc C.uint
	if result := C.nvmlDeviceGetTemperature(s.d, C.NVML_TEMPERATURE_GPU, &tempc); result != C.NVML_SUCCESS {
		return -1, -1, getGoError(result)
	}

	return int(tempc), int(tempc*9/5 + 32), nil
}

// GetMemoryInfo retrieves the amount of used, free and total memory available on the device, in bytes.
func (s *Device) GetMemoryInfo() (memInfo *NVMLMemory, err error) {
	var res C.nvmlMemory_t

	if result := C.nvmlDeviceGetMemoryInfo(s.d, &res); result != C.NVML_SUCCESS {
		return nil, getGoError(result)
	}

	return &NVMLMemory{
		Free:  int64(res.free),
		Total: int64(res.total),
		Used:  int64(res.used),
	}, nil
}

// GetProcessInfo retrieves the active proccesses (pid, used gpu memory) running on the device
func (s *Device) GetProcessInfo() (procInfo []NVMLProcess, err error) {

	var res []C.nvmlProcessInfo_t
	var cnt C.uint = 0

	// first query the number of procs (cnt := 0)
	result := C.nvmlDeviceGetComputeRunningProcesses(s.d, &cnt, nil)
	if result == C.NVML_SUCCESS {
		// no processes
		return nil, nil
	}

	if result == C.NVML_ERROR_INSUFFICIENT_SIZE {
		// we now know the number of procs (see cnt)
		num_procs := int(cnt)
		res = make([]C.nvmlProcessInfo_t, num_procs)

		if result := C.nvmlDeviceGetComputeRunningProcesses(s.d, &cnt, &res[0]); result != C.NVML_SUCCESS {
			return nil, getGoError(result)
		}

		procsInfo := make([]NVMLProcess, num_procs)
		for i := 0; i < num_procs; i++ {
			procsInfo = append(procsInfo, NVMLProcess{
				Pid:           int(res[i].pid),
				UsedGpuMemory: int64(res[i].usedGpuMemory),
			})
		}

		return procsInfo, nil
	}

	return nil, errNoError

}

// InitNVML initializes NVML
func InitNVML() error {
	if result := C.nvmlInit(); result != C.NVML_SUCCESS {
		return getGoError(result)
	}
	return nil
}

// ShutdownNVML all resources that were created when we initialized
func ShutdownNVML() error {
	if result := C.nvmlShutdown(); result != C.NVML_SUCCESS {
		return getGoError(result)
	}
	return nil
}

// GetDeviceCount returns the # of CUDA devices present on the host
func GetDeviceCount() (int, error) {
	var cnt C.uint

	if result := C.nvmlDeviceGetCount(&cnt); result != C.NVML_SUCCESS {
		return 0, getGoError(result)
	}
	return int(cnt), nil
}

// DeviceGetHandleByIndex acquires the handle for a particular device, based on its index.
func DeviceGetHandleByIndex(idx int) (*C.nvmlDevice_t, error) {
	var device *C.nvmlDevice_t

	if result := C.nvmlDeviceGetHandleByIndex(C.uint(idx), device); result != C.NVML_SUCCESS {
		return nil, getGoError(result)
	}
	return device, nil
}

// GetDevices returns a list of all installed CUDA devices
func GetDevices() ([]Device, error) {
	var nvdev C.nvmlDevice_t

	devCount, err := GetDeviceCount()
	if err != nil {
		return nil, err
	}

	devices := make([]Device, devCount)
	for i := 0; i <= devCount-1; i++ {
		if result := C.nvmlDeviceGetHandleByIndex(C.uint(i), &nvdev); result != C.NVML_SUCCESS {
			return nil, getGoError(result)
		}

		if devices[i], err = newDevice(nvdev, i); err != nil {
			return nil, err
		}
	}

	return devices, nil
}
