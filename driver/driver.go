package driver

// Tools for running SU2 from Go

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"sync"

	"github.com/btracey/su2tools/config"
)

// SU2Syscall is a way of running SU2 from a system call
type SU2Syscaller interface {
	SyscallString(d *Driver) (string, []string)
	Concurrently() bool
}

type FileWriter interface {
	WriteFile(d *Driver) error
}

// Serial runs SU2 on one processor
type Serial struct {
	Concurrent bool
}

func (s Serial) SyscallString(d *Driver) (execname string, args []string) {
	return "SU2_CFD", []string{d.Config}
}

func (s Serial) Concurrently() bool {
	return s.Concurrent
}

type Parallel struct {
	NumCores   int
	Concurrent bool
}

func (p Parallel) SyscallString(d *Driver) (execname string, args []string) {
	execname = "parallel_computation.py"
	args = []string{"-f", d.Config, "-p", strconv.Itoa(p.NumCores)}
	return
}

func (p Parallel) Concurrently() bool {
	return p.Concurrent
}

type Cluster struct {
	Cores      int
	Concurrent bool
}

func (c Cluster) SyscallString(d *Driver) (execname string, args []string) {
	fmt.Println("In cluster syscall")
	//execname = "sbatch"
	//args = []string{"slurm_script.run"}

	execname = "srun"
	coresStr := strconv.Itoa(c.Cores)
	args = []string{"--job-name", d.Name, "-n", coresStr, "--output", d.Name + "slurm.out", "parallel_computation.py", "-f", d.Config, "-p", coresStr}
	return
}

/*
func (c Cluster) WriteFile(d *Driver) error {
	fmt.Println("In write string")
	name := filepath.Join(d.Wd, "slurm_script.run")
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	f.WriteString("#!/bin/bash\n")
	f.WriteString("#SBATCH --job-name=" + d.Name + "\n")
	f.WriteString("#SBATCH -n" + strconv.Itoa(c.Cores) + "\n")
	f.WriteString("#SBATCH --output=slurm.out\n")
	f.WriteString("#SBATCH --error=slurm.err\n")
	f.WriteString("parallel_computation.py -f " + d.Config + " -p " + strconv.Itoa(c.Cores))
	f.Close()
	return nil
	//panic("check bash script")
}
*/

func (c Cluster) Concurrently() bool {
	return c.Concurrent
}

// Driver is an SU2 case to be run
type Driver struct {
	Name       string            // Identifier for the case
	Options    *config.Options   // OptionList for the case
	Config     string            // Name of the config filename (relative to working directory)
	Wd         string            // Working directory of SU2
	Stdout     string            // Where should the run output go (relative to working directory or StIO)
	OptionList config.OptionList // Which options to print
}

// IsComputed checks if the case specified by the driver has already been run
func (d *Driver) IsComputed() bool {
	// First, check if the file exists
	f, err := os.Open(d.Fullpath(d.Config))
	if err != nil {
		fmt.Println("not computed, no config file")
		return false
	}

	// Next, check if the options file is the same
	oldOptions, _, err := config.Read(f)
	if err != nil {
		fmt.Println("not computed, error reading config file")
		return false
	}
	if !reflect.DeepEqual(d.Options, oldOptions) {
		fmt.Println("not computed, config files unequal")
		return false
	}
	// Same config file already exists, see if the solution file exists
	_, err = os.Open(d.Fullpath(d.Options.SolutionFlowFilename))
	if err != nil {
		fmt.Println("not computed, no solution file")
		return false
	}
	return true
}

// Run writes the config file and calls SU2 with the output specified
func (d *Driver) Run(su2call SU2Syscaller) error {

	err := os.MkdirAll(d.Wd, 0700)
	if err != nil {
		return err
	}
	// Write the config file
	f, err := os.Create(d.Fullpath(d.Config))
	if err != nil {
		return err
	}
	d.Options.WriteConfig(f, d.OptionList)
	f.Close()

	// Open the standard out file
	stdout, err := os.Create(d.Fullpath(d.Stdout))
	defer stdout.Close()
	if err != nil {
		return err
	}

	// Create the command
	name, args := su2call.SyscallString(d)
	cmd := exec.Command(name, args...)
	cmd.Stdout = stdout
	cmd.Dir = d.Wd
	cmd.Stderr = os.Stdout
	fmt.Println("Command name: ", name)
	fmt.Println("Command args: ", args)

	return cmd.Run()
}

// CopyRestartToSolution copies the restart file to the solution file
func (d *Driver) CopyRestartToSolution() error {
	s, err := os.Create(d.Fullpath(d.Options.SolutionFlowFilename))
	if err != nil {
		return fmt.Errorf("copyrestart: error opening SolutionFlowFilename: %v", err)
	}
	defer s.Close()

	r, err := os.Open(d.Fullpath(d.Options.RestartFlowFilename))
	if err != nil {
		return fmt.Errorf("copyrestart: error opening RestartFlowFilename: %v", err)
	}
	defer r.Close()
	_, err = io.Copy(s, r)
	return err
}

/*
// SetStdout sets the standard output from the filename.
// isrel specifies if the path is a path relative to the working
// directory (or is an absolute path)
func (d *Driver) SetStdout(filename string, isrel bool) error {
	f, err := d.fileFromFilename(filename, isrel, true)
	d.Stdout = f
	return err
}
*/

// Fullpath returns the full path
func (d *Driver) Fullpath(relpath string) string {
	return filepath.Join(d.Wd, relpath)
}

// SetRelativeOptions sets the options structure relative to the from the config file in filename.
// filename is the base config file
// is rel specifies if the filepath is a relative path or an absolute path
// delta specifies the changes to be made to the config file
// This initializes the Options struct and sets the fieldmap to be all of the fields
// read in from the base config file
func (d *Driver) SetRelativeOptions(filename string, isrel bool, delta config.FieldMap) error {
	f, err := d.fileFromFilename(filename, isrel, false)
	if err != nil {
		return err
	}
	// FromBase generates a new config file with the file in reader as a
	// base, then adding the changes in delta
	options, optionList, err := config.Read(f)
	if err != nil {
		return err
	}
	err = options.SetFields(delta)
	if err != nil {
		return err
	}
	d.Options = options
	d.OptionList = optionList
	return nil
}

func (d *Driver) fileFromFilename(filename string, isrel bool, create bool) (*os.File, error) {
	if isrel {
		if d.Wd == "" {
			return nil, errors.New("driver: cannot use relative path; working directory not set")
		}
		filename = d.Fullpath(filename)
	}

	//fmt.Println("Filename = ", filename)
	if create {
		path := filepath.Dir(filename)
		err := os.MkdirAll(path, 0700)
		if err != nil {
			return nil, err
		}
		return os.Create(filename)
	}
	return os.Open(filename)
}

// RunCases runs the cases concurrently and wait until they're done.
// If redo is true, recompute the cases even if IsComputed() returns true
func RunCases(drivers []*Driver, su2call SU2Syscaller, redo bool) []error {

	Errors := make([]error, len(drivers))
	// TODO: Need to combine parallel and non-parallel code
	if !su2call.Concurrently() {
		for i, driver := range drivers {
			Errors[i] = runcase(redo, driver, su2call)
		}
	} else {
		w := &sync.WaitGroup{}
		// Run all of the training cases

		for i, driver := range drivers {
			w.Add(1)
			go func(i int, driver *Driver) {
				Errors[i] = runcase(redo, driver, su2call)
				w.Done()
			}(i, driver)
		}
		w.Wait()
	}
	for _, err := range Errors {
		if err != nil {
			return Errors
		}
	}
	return nil
}

func runcase(redo bool, driver *Driver, su2call SU2Syscaller) error {
	if redo || !driver.IsComputed() {
		fmt.Printf("\t%s: starting: wd: %s, config: %s\n", driver.Name, driver.Wd, driver.Config)

		// See if the caller needs to write a file
		fw, ok := su2call.(FileWriter)
		if ok {
			err := fw.WriteFile(driver)
			if err != nil {
				return err
			}
		}

		err := driver.Run(su2call)
		if err != nil {
			fmt.Printf("Error case %s\n", driver.Name)
		} else {
			fmt.Printf("Finished running %s\n", driver.Name)
		}
		return err
	} else {
		fmt.Printf("\t%s: already valid solution file %s\n", driver.Name, driver.Options.SolutionFlowFilename)
		return nil
	}
}
