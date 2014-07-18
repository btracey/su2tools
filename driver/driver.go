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
	Error Status = iota
	ComputedSuccessfully
	ComputedWithError
	NoConfigFile
	ErrorParsingConfig
	UnequalOptions
	NoSolutionFile
	HasErrorFile
)

var statusMap = map[Status]string{
	ComputedSuccessfully: "computed",
	NoConfigFile:         "no existing config file",
	ErrorParsingConfig:   "error parsing the existing config file in the working directory",
	UnequalOptions:       "options in the config file are not the same",
	NoSolutionFile:       "no solution file found",
	HasErrorFile:         "computed with error",
	Error:                "has both solution and error file",
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

// IsComputed returns true if the status indicates that the solution has been computed
func (d *Driver) IsComputed(stat Status) bool {
	return stat == ComputedSuccessfully || stat == ComputedWithError
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

	var hasErrorFile bool
	var hasSolutionFile bool

	// The config file may have been run, but SU^2 errored. See if the error file
	// exists
	if _, err = os.Stat(d.fullpath(d.errorFilename())); err == nil {
		hasErrorFile = true
		//return ComputedWithError
	}

	if _, err = os.Stat(d.fullpath(d.Options.SolutionFlowFilename)); err == nil {
		hasSolutionFile = true
	}

	if hasSolutionFile && hasErrorFile {
		return Error
	}

	if hasErrorFile {
		return ComputedWithError
	}
	if !hasSolutionFile {
		return NoSolutionFile
	}
	// The cases are the same
	return ComputedSuccessfully
}

// Run executes SU2 with the given Syscaller. Run writes the config file (as
// specified by driver.Config and driver.Wd), uses the Syscaller to get
// the arguments to exec.Command, and calls exec.Cmd.Run with the provided
// working directory, stdout, and stderr. If run has an error, an error file
// will be written.
func (d *Driver) Run(su2call Syscaller) error {
	if d.Wd != "" {
		err := os.MkdirAll(d.Wd, 0700)
		if err != nil {
			return errors.New("driver: error creating working directory: " + err.Error())
		}
	}
	// Erase the error file if is exists (we are going to re-run SU^2 so the
	// old error file is not relevant)
	errFilename := d.fullpath(d.errorFilename())
	if _, err := os.Stat(errFilename); err == nil {
		if err := os.Remove(errFilename); err != nil {
			return err
		}
	}

	// Write the config file
	f, err := os.Create(d.fullpath(d.Config))
	if err != nil {
		return err
	}
	d.Options.WriteTo(f, d.OptionList)
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

	fmt.Println("executing execname", "name = ", name, " args = ", args, "wd = ", cmd.Dir)
	runerr := cmd.Run()
	if runerr == nil {
		return err
	}
	// There was an error running SU^2. Create it with the error.
	errFile, err := os.Create(d.fullpath(d.errorFilename()))
	if err != nil {
		return err
	}
	defer errFile.Close()
	errFile.WriteString(runerr.Error())
	return runerr
}

// MoveSolutionToOld moves the solution file to d.Options.SolutionFlowFilename +
// This can be useful for automated tools. (In case of error, the old solution file)
// is not lost
func (d *Driver) MoveSolutionToOld() error {
	// Get the solution flow filename and split it up into a beginning and an end
	// Append _old to it.
	// Copy the old file to the new file

	currFullPath := d.fullpath(d.Options.SolutionFlowFilename)
	oldFullPath := d.fullpath(oldSolutionFilename(d.Options.SolutionFlowFilename))

	return os.Rename(currFullPath, oldFullPath)
}

// oldSolutionFilename translaties the current solutionFlowFilename into the old
// solution flow filename
func oldSolutionFilename(sol string) string {
	ext := filepath.Ext(sol)
	prefix := sol[:len(sol)-len(ext)]
	return prefix + "_old" + ext
}

// EraseDoneFiles deletes the solution file and crash file
func (d *Driver) EraseDoneFiles() error {
	solFullPath := d.fullpath(d.Options.SolutionFlowFilename)
	if err := os.Remove(solFullPath); err != nil {
		return err
	}
	errorFullPath := d.fullpath(d.errorFilename())
	return os.Remove(errorFullPath)
}

// crashFilename is the name of the file that is written if SU^2 exited with error
func (d *Driver) errorFilename() string {
	// Error filename should live in the pwd and should have name su2_error.txt
	return "su2_err.txt"
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
