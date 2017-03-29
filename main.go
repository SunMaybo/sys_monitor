// SysMonitor project main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type SysBaseInfo struct {
	Cpu  string    `json:"cpu"`
	Mem  int64     `json:"mem"`
	Host string    `json:"hostname"`
	Due  time.Time `json:"due"`
}
type MemoryInfo struct {
	Total     int64     `json:"total"` //总内存
	Used      int64     `json:"used"`  //使用的
	Free      int64     `json:"free"`  //空闲的
	Shared    int64     `json:"shared"`
	Cache     int64     `json:"cache"`
	Available int64     `json:"available"`
	Due       time.Time `json:"due"`
}
type CpuInfo struct {
	TaskProcessCount int       `json:"taskProcessCount"`
	IOProcessCount   int       `json:"ioProcessCount"`
	Us               int       `json:"us"`
	Sy               int       `json:"sy"`
	Id               int       `json:"id"`
	Due              time.Time `json:"due"`
}
type DeviceInfo struct {
	Rx       string    `json:"rx"`
	Tx       string    `json:"tx"`
	RxData   string    `json:"rxData"`
	TxData   string    `json:"txData"`
	RxErrors string    `json:"rxErrors"`
	TxErrors string    `json:"txErrors"`
	RxOver   string    `json:"rxOver"`
	TxColl   string    `json:"txColl"`
	Due      time.Time `json:"due"`
}
type SysInfo struct {
	MemSysInfo      MemoryInfo `json:"memInfo"`
	CpuSysInfo      CpuInfo    `json:"cpuInfo"`
	DeviceSysInfo   DeviceInfo `json:"deviceInfo"`
	PortAcceptCount int        `json:"portAcceptCount"`
	Due             time.Time  `json:"due"`
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/sys/info", sys_info)
	router.HandleFunc("/sys/info/current", sys_current_info)
	router.HandleFunc("/sys/current/ce/info", sys_cpu_mem_current_info)
	router.HandleFunc("/sys/current/tcpAccept/{port}", sys_current_port_acceptCount)
	router.HandleFunc("/sys/current/device/{device}", sys_device_current_info)
	log.Fatal(http.ListenAndServe(":10329", router))
}
func execCommand(command string) string {
	log.Println(command)
	cmd := exec.Command("/bin/sh", "-c", command)

	bytes, err := cmd.Output()
	if err != nil {
		fmt.Println("cmd.Output: ", err)
		return ""
	}
	return strings.Split(string(bytes), "\n")[0]
}

func init_sys() (string, string) {

	return execCommand("grep MemTotal /proc/meminfo | awk '{print $2}'"), execCommand("cat /proc/cpuinfo | grep name | cut -f2 -d: | uniq -c")
}

func current_sys_cpu_mem() (string, string) {
	memString := "free | grep Mem | awk '{print $2,$3,$4,$5,$6,$7}'"
	cpuString := "vmstat |grep '[0-9]'| awk '{print $1,$2,$13,$14,$15}'"
	return execCommand(memString), execCommand(cpuString)
}
func current_port_concurrency(port string) string {
	commandString := "netstat -nat|grep ESTABLISHED |awk '{print $4}' | grep -i "
	commandString += port + "|wc -l"
	log.Println(commandString)
	return execCommand(commandString)
}

func current_device_concurrency(device string) string {
	deviceString := "ifstat -a | grep " + device + " |awk '{print $2,$3,$4,$5,$6,$7,$8,$9}'"
	return execCommand(deviceString)
}
func hostName() string {
	deviceString := "echo $HOSTNAME"
	return execCommand(deviceString)
}
func sys_port_acceptCount_serice(port string) string {
	countStr := current_port_concurrency(port)
	return strings.TrimSpace(countStr)
}
func sys_cpu_mem_service() SysInfo {
	memInfoStr, cpuInfoStr := current_sys_cpu_mem()
	memInfoArray := strings.Split(memInfoStr, " ")
	cpuInfoArray := strings.Split(cpuInfoStr, " ")
	total, err := strconv.ParseInt(memInfoArray[0], 10, 64)
	if err != nil {

	}
	used, err := strconv.ParseInt(memInfoArray[1], 10, 64)
	if err != nil {

	}
	free, err := strconv.ParseInt(memInfoArray[2], 10, 64)
	if err != nil {

	}
	shared, err := strconv.ParseInt(memInfoArray[3], 10, 64)
	if err != nil {

	}
	cache, err := strconv.ParseInt(memInfoArray[4], 10, 64)
	if err != nil {

	}
	available, err := strconv.ParseInt(memInfoArray[5], 10, 64)
	if err != nil {

	}

	memoryInfo := MemoryInfo{Total: total, Used: used, Free: free, Shared: shared, Cache: cache, Available: available, Due: time.Now()}
	taskProcessCount, err := strconv.Atoi(cpuInfoArray[0])
	if err != nil {

	}
	ioProcessCount, err := strconv.Atoi(cpuInfoArray[1])
	if err != nil {

	}
	us, err := strconv.Atoi(cpuInfoArray[2])
	if err != nil {

	}
	sy, err := strconv.Atoi(cpuInfoArray[3])
	if err != nil {

	}
	id, err := strconv.Atoi(cpuInfoArray[4])
	if err != nil {

	}
	cpuInfo := CpuInfo{TaskProcessCount: taskProcessCount, IOProcessCount: ioProcessCount, Us: us, Sy: sy, Id: id, Due: time.Now()}
	sysInfo := SysInfo{MemSysInfo: memoryInfo, CpuSysInfo: cpuInfo, Due: time.Now()}
	return sysInfo
}
func sys_device_info_service(device string) DeviceInfo {
	baseDeviceInfo := current_device_concurrency(device)
	var deviceInfo DeviceInfo
	if !strings.EqualFold("", baseDeviceInfo) {
		deviceInfoArray := strings.Split(baseDeviceInfo, " ")
		deviceInfo = DeviceInfo{Rx: deviceInfoArray[0], Tx: deviceInfoArray[1], RxData: deviceInfoArray[2], TxData: deviceInfoArray[3], RxErrors: deviceInfoArray[4], TxErrors: deviceInfoArray[5], RxOver: deviceInfoArray[6], TxColl: deviceInfoArray[7], Due: time.Now()}
		return deviceInfo
	}
	return deviceInfo
}
func sys_current_port_acceptCount(w http.ResponseWriter, r *http.Request) {
	// Get the port.
	vars := mux.Vars(r)
	portStr := vars["port"]
	countStr := sys_port_acceptCount_serice(portStr)
	json.NewEncoder(w).Encode(countStr)
}
func sys_cpu_mem_current_info(w http.ResponseWriter, r *http.Request) {
	sysInfo := sys_cpu_mem_service()
	json.NewEncoder(w).Encode(sysInfo)
}
func sys_device_current_info(w http.ResponseWriter, r *http.Request) {
	// Get the device.
	vars := mux.Vars(r)
	device := vars["device"]
	deviceInfo := sys_device_info_service(device)
	json.NewEncoder(w).Encode(deviceInfo)
}
func sys_info(w http.ResponseWriter, r *http.Request) {
	memStr, cpu := init_sys()
	host := hostName()
	fmt.Println(host)
	mem, err := strconv.ParseInt(memStr, 10, 64)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		sysInfo := SysBaseInfo{Cpu: cpu, Mem: mem, Host: host, Due: time.Now()}
		json.NewEncoder(w).Encode(sysInfo)
	}
}
func sys_current_info(w http.ResponseWriter, r *http.Request) {
	sysInfo := sys_cpu_mem_service()
	sysInfo.DeviceSysInfo = sys_device_info_service("eth0")
	if strings.EqualFold(sysInfo.DeviceSysInfo.Rx, "") {
		sysInfo.DeviceSysInfo = sys_device_info_service("enp3s0")
	}
	json.NewEncoder(w).Encode(sysInfo)
}
