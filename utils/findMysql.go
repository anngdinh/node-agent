package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

func FindMysql() {
	// Command to get the process ID of MySQL
	cmd := exec.Command("pgrep", "mysqld", "-u", "mysql")

	// Execute the command and capture the output
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Convert the output to a string and remove any leading/trailing spaces or newlines
	pid := strings.TrimSpace(string(output))
	fmt.Println(pid)

	if pid == "" {
		fmt.Println("MySQL is not running.")
	} else {
		fmt.Printf("MySQL is running with process ID: %s\n", pid)
	}

	process, _, err := FindByPID(1367)
	if err != nil {
		fmt.Printf("explorer.exe PID: %v", process)
	}

	process, _, err = FindByName("mysqld")
	if err != nil {
		fmt.Printf("explorer.exe PID: %v", process)
	}

	p, _ := Find(1367)
	fmt.Println(p)
}

// ErrNotFound - process not found error
var ErrNotFound = errors.New("process not found")

// ProcessEx - os.P
type ProcessEx struct {
	*os.Process
	Name      string
	PID       int
	ParentPID int
}

type processes []*ProcessEx

func (p processes) make() {
	if p == nil {
		p = make([]*ProcessEx, 0, 100)
	}
}

func (p processes) find(name string, pid int) (found []*os.Process, foundEx []*ProcessEx, err error) {
	name = strings.ToLower(name)
	for _, process := range p {
		fmt.Println(process)
		if (name != "" && process.Name == name) || (pid != 0 && process.PID == pid) {
			procTmp, err := os.FindProcess(process.PID)
			if err != nil {
				return nil, nil, err
			}
			process.Process = procTmp
			found = append(found, procTmp)
			foundEx = append(foundEx, process)
		}
	}

	if len(found) > 0 {
		return found, foundEx, nil
	}
	return nil, nil, ErrNotFound
}

// mysql       1674  0.9  2.4 2373280 393256 ?      Ssl  18:26   1:50 /usr/sbin/mysqld
// telegraf    2958  0.9  2.4 2369072 390788 ?      Ssl  18:26   1:55 mysqld
// annd2      31127  0.0  0.0  17864  2392 pts/4    S+   21:44   0:00 grep --color=auto mysql

type linuxProcesses struct {
	processes
	updatedAt int64 //atomic time
}

func (p *linuxProcesses) fetchPID(path string) (int, error) {
	indx := strings.LastIndex(path, "/")
	if indx < 0 || len(path) < 7 {
		return -1, fmt.Errorf("fetch pid error, path: %v", path)
	}
	return strconv.Atoi(path[6:indx])
}

// ------------------------------------------------------------------

func (p *linuxProcesses) fetchName(path string) (string, error) {
	// The status file contains the name of the process in its first line.
	// The line looks like "Name: theProcess".
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	// Extract the process name from within the first line in the buffer
	name := string(f[6:bytes.IndexByte(f, '\n')])
	if len(name) < 2 {
		return "", errors.New("fetch name error")
	}
	return name, nil
}

// ------------------------------------------------------------------

func (p *linuxProcesses) walk(path string, info os.FileInfo, err error) error {
	// We just return in case of errors, as they are likely due to insufficient
	// privileges. We shouldn't get any errors for accessing the information we
	// are interested in. Run as root (sudo) and log the error, in case you want
	// this information.
	if err != nil {
		return nil
	}

	// We are only interested in files with a path looking like /proc/<pid>/status.
	if strings.Count(path, "/") != 3 || !strings.Contains(path, "/status") {
		return nil
	}

	// Let's extract the middle part of the path with the <pid> and
	// convert the <pid> into an integer. Log an error if it fails.
	pid, err := p.fetchPID(path)
	if err != nil {
		return err
	}

	// Extract the process name from within the first line in the buffer
	name, err := p.fetchName(path)
	if err != nil {
		return err
	}

	p.processes = append(p.processes, newProcessEx(name, pid, 0))
	return nil
}
func newProcessEx(name string, pid, parentPID int) *ProcessEx {
	return &ProcessEx{
		Name:      strings.ToLower(name),
		PID:       pid,
		ParentPID: parentPID,
	}
}

// ------------------------------------------------------------------

func (p *linuxProcesses) make() {
	if p.processes == nil {
		p.processes = make([]*ProcessEx, 0, 100)
	} else {
		p.processes = p.processes[:0]
	}
}

// ------------------------------------------------------------------

func (p *linuxProcesses) getUpdatedAt() time.Time {
	return time.Unix(atomic.LoadInt64(&p.updatedAt), 0)
}

func (p *linuxProcesses) setUpdatedAt() {
	atomic.StoreInt64(&p.updatedAt, time.Now().Unix())
}

func (p *linuxProcesses) getProcesses() error {
	if time.Now().Sub(p.getUpdatedAt()) < time.Second*3 {
		return nil
	}
	p.setUpdatedAt()
	p.make()
	return filepath.Walk("/proc", p.walk)
}

// ------------------------------------------------------------------

// FindByName - FindByName
func (p *linuxProcesses) FindByName(name string) ([]*os.Process, []*ProcessEx, error) {
	if name == "" {
		return nil, nil, fmt.Errorf("%w, name is empty", ErrNotFound)
	}
	err := p.getProcesses()
	if err != nil {
		return nil, nil, err
	}
	return p.find(name, 0)
}

// FindByPID - FindByPID
func (p *linuxProcesses) FindByPID(pid int) ([]*os.Process, []*ProcessEx, error) {
	if pid == 0 {
		return nil, nil, fmt.Errorf("%w, pid == 0", ErrNotFound)
	}
	err := p.getProcesses()
	if err != nil {
		return nil, nil, err
	}
	return p.find("", pid)
}

// ------------------------------------------------------------------

func Find(pid int) (*os.Process, error) {
	return os.FindProcess(pid)
}

// ------------------------------------------------------------------

// FindByName looks for a running process by its name.
//
// The Process it returns can be used to obtain information
// about the underlying operating system process.
func FindByName(name string) ([]*os.Process, []*ProcessEx, error) {
	return NewFinder().FindByName(name)
}

// ------------------------------------------------------------------

// FindByPID looks for a running process by its PID.
//
// The Process it returns can be used to obtain information
// about the underlying operating system process.
func FindByPID(pid int) ([]*os.Process, []*ProcessEx, error) {
	return NewFinder().FindByPID(pid)
}

// ------------------------------------------------------------------

// Start starts a new process with the program, arguments and attributes
// specified by name, argv and attr. The argv slice will become os.Args in the
// new process, so it normally starts with the program name.
func Start(name string, argv []string, attr *os.ProcAttr) (*os.Process, error) {
	return os.StartProcess(name, argv, attr)
}

// Finder - system processes Finder
type Finder interface {
	FindByName(name string) ([]*os.Process, []*ProcessEx, error)
	FindByPID(pid int) ([]*os.Process, []*ProcessEx, error)
}

// ------------------------------------------------------------------

// NewFinder - NewFinder
func NewFinder() Finder {
	switch {
	default:
		return &linuxProcesses{}
	}
}
