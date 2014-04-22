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

	"github.com/btracey/su2tools/config"
)

// Status is a type that represents if the computation specified by the driver
// has been computed, and if not, why not. Note that this is not the status
// during a call to Driver.Run, but rather used for checking if this config file
// has already been successfully been completed
type Status int

func (s Status) String() string {
	return statusMap[s]
}

const (
	Computed Status = iota + 1
	NoConfigFile
	ErrorParsingConfig
	UnequalOptions
	NoSolutionFile
)

var statusMap = map[Status]string{
	Computed:           "computed",
	NoConfigFile:       "no existing config file",
	ErrorParsingConfig: "error parsing the existing config file in the working directory",
	UnequalOptions:     "options in the config file are not the same",
	NoSolutionFile:     "no solution file found",
}

// Syscaller is a provides the system call to execute a Driver. The returned
// arguments will be passed directly to exec.Cmd
type Syscaller interface {
	Syscall(d *Driver) (string, []string) // Returns the exec name and the arguments to be called by exec.Command
	NumCores() int                        // How many cores does it want
}

/*
// A FileWriter is a type that needs to print a file to disk before executing the driver
type FileWriter interface {
	WriteFile(d *Driver) error
}
*/

// Serial runs SU2 on one processor.
type Serial struct{}

func (s Serial) Syscall(d *Driver) (execname string, args []string) {
	return "SU2_CFD", []string{d.Config}
}

func (s Serial) NumCores() int {
	return 1
}

// Parallel runs SU2 in parallel with the specified number of cores. SU2 must
// be compiled with MPI
type Parallel struct {
	Cores int
}

func (p Parallel) Syscall(d *Driver) (execname string, args []string) {
	execname = "parallel_computation.py"
	args = []string{"-f", d.Config, "-p", strconv.Itoa(p.Cores)}
	return
}

func (p Parallel) NumCores() int {
	return p.Cores
}

// Slurm runs SU2 via a call to Slurm. This does not work at present.
type Slurm struct {
	Cores int
}

func (c Slurm) Syscall(d *Driver) (execname string, args []string) {
	fmt.Println("In cluster syscall")
	//execname = "sbatch"
	//args = []string{"slurm_script.run"}

	execname = "srun"
	coresStr := strconv.Itoa(c.Cores)
	args = []string{"--job-name", d.Name, "-n", coresStr, "--output", d.Name + "slurm.out", "parallel_computation.py", "-f", d.Config, "-p", coresStr}

	return
}

func (c Slurm) NumCores() int {
	return c.Cores
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

// Driver specifies a case for running SU2.
type Driver struct {
	Name       string                 // Identifier for the case
	Options    *config.Options        // OptionList for the case
	Config     string                 // Name of the config filename (relative to working directory)
	Wd         string                 // Working directory of SU2
	Stdout     string                 // Where to redict Stdout (relative to working directory). Will be set to Stdout if ==""
	Stderr     string                 // Where to redirct Stderr (relative to working directory). Will be set to Stderr if ==""
	OptionList map[config.Option]bool // Which options to print to the config file
	FancyName  string                 // Longer name (can be used for plot legends or something)
}

// IsComputed returns true if driver.Status() == Computed
func (d *Driver) IsComputed() bool {
	return d.Status() == Computed
}

// Status returns the status of the computation. If this particular set of options
// has been computed before (defined by the config file and the working directory)
// it returns Computed, otherwise it returns one of the other status constants.
func (d *Driver) Status() Status {
	// First, check if the file exists
	f, err := os.Open(d.fullpath(d.Config))
	if err != nil {
		return NoConfigFile
	}

	// Next, check if the options file is the same
	oldOptions, _, err := config.Read(f)
	if err != nil {
		return ErrorParsingConfig
	}
	if !reflect.DeepEqual(d.Options, oldOptions) {
		return UnequalOptions
	}
	// Same config file already exists and is the same, see if the solution file exists
	_, err = os.Open(d.fullpath(d.Options.SolutionFlowFilename))
	if err != nil {
		return NoSolutionFile
	}
	// The cases are the same
	return Computed
}

// Run executes SU2 with the given Syscaller. Run writes the config file (as
// specified by driver.Config and driver.Wd), uses the Syscaller to get
// the arguments to exec.Command, and calls exec.Cmd.Run with the provided
// working directory, stdout, and stderr
func (d *Driver) Run(su2call Syscaller) error {

	err := os.MkdirAll(d.Wd, 0700)
	if err != nil {
		return errors.New("driver: error creating working directory: " + err.Error())
	}
	// Write the config file
	f, err := os.Create(d.fullpath(d.Config))
	if err != nil {
		return err
	}
	d.Options.WriteConfigTo(f, d.OptionList)
	f.Close()

	var stdout io.Writer

	if d.Stdout == "" {
		stdout = os.Stdout
	} else {
		// Open the standard out file
		f, err := os.Create(d.fullpath(d.Stdout))
		if err != nil {
			return errors.New("driver: error creating stdout: " + err.Error())
		}
		defer f.Close()
		stdout = f
	}
	var stderr io.Writer
	if d.Stderr == "" {
		stderr = os.Stderr
	} else {
		// Open the standard error file
		f, err := os.Create(d.fullpath(d.Stderr))
		if err != nil {
			return errors.New("driver: error creating stderr: " + err.Error())
		}
		defer f.Close()
		stderr = f
	}

	// Create the command
	name, args := su2call.Syscall(d)
	cmd := exec.Command(name, args...)
	cmd.Stdout = stdout
	cmd.Dir = d.Wd
	cmd.Stderr = stderr

	fmt.Println("executing execname", "name = ", name, " args = ", args)
	return cmd.Run()
}

// CopyRestartToSolution copies the restart file to the solution file
func (d *Driver) CopyRestartToSolution() error {
	s, err := os.Create(d.fullpath(d.Options.SolutionFlowFilename))
	if err != nil {
		return fmt.Errorf("copyrestart: error opening SolutionFlowFilename: %v", err)
	}
	defer s.Close()

	r, err := os.Open(d.fullpath(d.Options.RestartFlowFilename))
	if err != nil {
		return fmt.Errorf("copyrestart: error opening RestartFlowFilename: %v", err)
	}
	defer r.Close()
	_, err = io.Copy(s, r)
	return err
}

// Fullpath returns the full path
func (d *Driver) fullpath(relpath string) string {
	return filepath.Join(d.Wd, relpath)
}

// Load reads in from the input and sets the Options and OptionList fields.
func (d *Driver) Load(reader io.Reader) error {
	// TODO: Add example
	options, optionList, err := config.Read(reader)
	if err != nil {
		return err
	}
	d.Options = options
	d.OptionList = optionList
	return nil
}

/*

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
*/

func (d *Driver) fileFromFilename(filename string, isrel bool, create bool) (*os.File, error) {
	if isrel {
		if d.Wd == "" {
			return nil, errors.New("driver: cannot use relative path; working directory not set")
		}
		filename = d.fullpath(filename)
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

/*
// RunCases runs a list of Drivers with the input SU2Syscaller if they have not been computed.
// Depending on the settings of the Syscaller, the cases will either be
// run concurrently or in parallel. The return is a list of errors, one for
// each element in drivers. If there were no errors, nil is returned.
// If redo is true, recompute the cases even if IsComputed() returns true
func RunCases(drivers []*Driver, su2call SU2Syscaller, redo bool) []error {

	Errors := make([]error, len(drivers))
	if !su2call.Concurrently() {
		// Run all of the cases in serial
		for i, driver := range drivers {
			Errors[i] = runcase(redo, driver, su2call)
		}
	} else {
		// Run all of the training cases concurrently
		w := &sync.WaitGroup{}
		for i, driver := range drivers {
			w.Add(1)
			go func(i int, driver *Driver) {
				Errors[i] = runcase(redo, driver, su2call)
				w.Done()
			}(i, driver)
		}
		w.Wait()
	}
	// If there were any errors, return the list, otherwise return nil
	for _, err := range Errors {
		if err != nil {
			return Errors
		}
	}
	return nil
}

func runcase(redo bool, driver *Driver, su2call SU2Syscaller) error {
	// See if the case needs to be run
	if !redo && driver.IsComputed() {
		fmt.Printf("\t%s: already valid solution file %s\n", driver.Name, driver.Options.SolutionFlowFilename)
		return nil
	}

	// Need to rerun the cases, so do so
	fmt.Printf("\t%s: starting: wd: %s, config: %s\n", driver.Name, driver.Wd, driver.Config)

	// See if the caller needs to write a file, and if so, write it
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
}
*/
