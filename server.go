package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

type Commander struct {
	Name        string
	Description string
	Command     string
	Args        []string
	RawOutput   []byte
	Error       error
}

func NewCommander() *Commander {
	return &Commander{}
}

func (c *Commander) Run() *Commander {
	c.RawOutput, c.Error = exec.Command(c.Command, c.Args...).Output()
	return c
}

func (c *Commander) Parse(s string) *Commander {
	arr := strings.Split(s, " ")
	c.Command = arr[0]
	c.Args = arr[1:]
	return c
}

func (c *Commander) isError() bool {
	return c.Error != nil
}

func (c *Commander) Output() string {
	return string(c.RawOutput)
}

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

type Storage struct {
	All  string `json:"all"`
	Used string `json:"used"`
	Free string `json:"free"`
}

type SystemInfo struct {
	Hostname string  `json:"hostname"`
	IP       string  `json:"ip"`
	OS       string  `json:"os"`
	CPUs     int     `json:"cpus"`
	Storage  Storage `json:"storage"`
	Uptime   string  `json:"uptime"`
}

func (s *SystemInfo) init() {
	s.Hostname = s.getHostname()
	s.IP = s.getLocalIP()
	s.OS = s.getOS()
	s.CPUs = s.getCPUs()
	s.Storage = s.getStorage("/")
	s.Uptime = s.getUptime()
}

func (s *SystemInfo) getHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}

func (s *SystemInfo) getOS() string {
	return strings.Title(runtime.GOOS)
}

func (s *SystemInfo) getCPUs() int {
	return runtime.NumCPU()
}

func (s *SystemInfo) getUptime() string {
	return NewCommander().Parse("uptime -p").Run().Output()
}

func (s *SystemInfo) getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func (s *SystemInfo) getStorage(path string) Storage {
	storage := Storage{}

	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return storage
	}

	diskAll := fs.Blocks * uint64(fs.Bsize)
	diskFree := fs.Bfree * uint64(fs.Bsize)
	diskUsed := diskAll - diskFree

	storage.All = fmt.Sprintf("%.2f GB", float64(diskAll)/float64(GB))
	storage.Free = fmt.Sprintf("%.2f GB", float64(diskUsed)/float64(GB))
	storage.Used = fmt.Sprintf("%.2f GB", float64(diskFree)/float64(GB))

	return storage
}

func NewSystemInfo() *SystemInfo {
	info := SystemInfo{}
	info.init()
	return &info
}

func main() {

	fmt.Println("Server start")
	si := NewSystemInfo()

	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))

	m.Get("/api/system", func(r render.Render) {
		si.Storage = si.getStorage("/")
		si.Uptime = si.getUptime()
		r.JSON(200, si)
	})

	m.Get("/api/cmd/ps", func(r render.Render) {
		r.Text(200, NewCommander().Parse("ps -aux").Run().Output())
	})

	m.Get("/api/cmd/df", func(r render.Render) {
		r.Text(200, NewCommander().Parse("df -h").Run().Output())
	})

	m.Get("/api/cmd/ifconfig", func(r render.Render) {
		r.Text(200, NewCommander().Parse("ifconfig").Run().Output())
	})

	m.Post("/api/cmd/reboot", func(r render.Render) {
		r.Text(200, NewCommander().Parse("reboot").Run().Output())
	})

	m.Post("/api/cmd/shutdown", func(r render.Render) {
		r.Text(200, NewCommander().Parse("shutdown -h now").Run().Output())
	})

	m.Run()
}
