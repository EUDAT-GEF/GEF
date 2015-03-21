// Build with:
// `GOOS=linux GOARCH=amd64 go build gefcommand.go`
// then move it to vagrant
// `mv ./gefcommand ../vagrant`
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

const (
	VERSION          string = "0.2.0"
	DEBUG            bool   = true
	ICOMM_DIR        string = "/var/lib/irods/iRODS/clients/icommands/bin/"
	STAGEIN_SUFFIX   string = "IrodsPath"
	STAGED_SUFFIX    string = "FilePath"
	WORKFLOW_KEY     string = "workflowFilePath"
	TAVERNA_EXECUTOR string = "/opt/taverna-commandline-2.4.0/executeworkflow.sh"
)

var (
	stagedir    string = "."
	stdoutName  string = ".out"
	stdoutFile  *os.File
	stderrName  string = ".err"
	stderrFile  *os.File
	logName     string = ".log"
	logFile     *os.File
	stagedFiles []string = make([]string, 0)
)

func main() {
	if len(os.Args) != 2 {
		log.Println("gefcommand v." + VERSION)
		log.Println("\tUsage: gefcommand irods_path_to_.gefcommand")
		log.Fatal("FATAL: missing argument")
	}

	irodsGefCommandFilePath := os.Args[1]

	gefFile := filepath.Base(irodsGefCommandFilePath)

	var err error
	{
		stdoutName = gefFile + ".out"
		if stdoutFile, err = os.Create(stdoutName); err != nil {
			fatal(err)
		}

		stderrName = gefFile + ".err"
		if stderrFile, err = os.Create(stderrName); err != nil {
			fatal(err)
		}

		logName = gefFile + ".log"
		if logFile, err = os.Create(logName); err != nil {
			fatal(err)
		}
	}

	if stagedir, err = ioutil.TempDir("", "gefcommand-stage-"); err != nil {
		fatal(err)
	}

	debug("stagedir = " + stagedir)

	if err = os.Chdir(stagedir); err != nil {
		fatal(err)
	}

	var args map[string]string
	if args, err = readGefCommandFile(irodsGefCommandFilePath); err != nil {
		fatal(err)
	}

	if args, err = stageInFiles(args); err != nil {
		fatal(err)
	}
	debug(args)

	if err = execute(args); err != nil {
		fatal(err)
	}

	stageOutFiles(irodsGefCommandFilePath, args)
}

func readGefCommandFile(irodsGefCommandFilePath string) (map[string]string, error) {
	if err := stageIn(irodsGefCommandFilePath); err != nil {
		return nil, err
	}
	gefCommandName := path.Base(irodsGefCommandFilePath)
	if gefCommandFile, err := os.Open(gefCommandName); err != nil {
		return nil, err
	} else {
		defer gefCommandFile.Close()

		args := make(map[string]string)
		scanner := bufio.NewScanner(gefCommandFile)
		for scanner.Scan() {
			line := scanner.Text()
			tokens := strings.SplitN(line, "=", 2)
			if len(tokens) != 2 {
				return nil, errors.New("unparseable gefcommand line: " + line)
			}
			args[tokens[0]] = tokens[1]
		}
		return args, scanner.Err()
	}
}

func stageInFiles(args map[string]string) (map[string]string, error) {
	newargs := make(map[string]string)
	for k, v := range args {
		if strings.HasSuffix(k, STAGEIN_SUFFIX) {
			irodsPath := v
			if err := stageIn(irodsPath); err != nil {
				return nil, err
			}
			l := len(k) - len(STAGEIN_SUFFIX)
			k = k[:l] + STAGED_SUFFIX
			v = stagedir + "/" + path.Base(irodsPath)
		}
		newargs[k] = v
	}
	return newargs, nil
}

func execute(args map[string]string) error {
	if value, ok := args[WORKFLOW_KEY]; ok {
		delete(args, WORKFLOW_KEY)
		if strings.HasSuffix(value, ".t2flow") {
			return runTavernaWorkflow(value, args)
		} else if strings.HasSuffix(value, ".test") {
			return runTest(value, args)
		}
		fatal("unknown workflow type: " + value)
	}
	fatal("function unspecified, missing workflow argument")
	return nil
}

func runTavernaWorkflow(workflowFile string, args map[string]string) error {
	str := ""
	for k, v := range args {
		if k == "datasetFilePath" {
			str = str + "-inputfile datasetFile " + v + " "
		} else {
			str = str + "-inputvalue " + k + " " + v + " "
		}
	}
	command := []string{workflowFile, str}
	debug(TAVERNA_EXECUTOR, command)

	cmd := exec.Command(TAVERNA_EXECUTOR, command...)
	cmd.Stdout = stdoutFile
	cmd.Stderr = stderrFile
	return cmd.Run()
}

func runTest(workflowFile string, args map[string]string) error {
	debug("TEST EXECUTOR", workflowFile, args)

	cmd := exec.Command("/bin/cat", workflowFile)
	cmd.Stdout = stdoutFile
	cmd.Stderr = stderrFile
	return cmd.Run()
}

func stageOutFiles(irodsGefCommandFilePath string, args map[string]string) error {
	irodsTarget, _ := filepath.Split(irodsGefCommandFilePath)

	if list, err := ioutil.ReadDir("."); err != nil {
		return err
	} else {
		for _, entry := range list {
			name := entry.Name()
			if !entry.IsDir() {
				if !contains(stagedFiles, name) &&
					name != stdoutName && name != stderrName && name != logName {
					stageOut(name, irodsTarget)
				}
			}
			// Taverna default output directory
			if name == "Workflow1_output" && entry.IsDir() {
				if list2, err := ioutil.ReadDir(name); err != nil {
					return err
				} else {
					for _, fi := range list2 {
						stageOut(name+"/"+fi.Name(), irodsTarget)
					}
				}
			}
		}
	}

	stdoutFile.Close()
	stageOut(stdoutName, irodsTarget)
	stderrFile.Close()
	stageOut(stderrName, irodsTarget)
	logFile.Close()
	logFile = nil
	stageOut(logName, irodsTarget)
	return nil
}

func stageIn(irodsPath string) error {
	stagedFiles = append(stagedFiles, filepath.Base(irodsPath))
	debug(ICOMM_DIR + "iget " + irodsPath)
	return exec.Command(ICOMM_DIR+"iget", irodsPath).Run()
}

func stageOut(filePath string, irodsTarget string) error {
	debug(ICOMM_DIR + "iput " + filePath + " " + irodsTarget)
	return exec.Command(ICOMM_DIR+"iput", filePath, irodsTarget).Run()
}

func contains(list []string, x string) bool {
	for _, a := range list {
		if a == x {
			return true
		}
	}
	return false
}

func debug(message ...interface{}) {
	if logFile != nil {
		logFile.WriteString(fmt.Sprintln(message...))
	}
}

func fatal(message ...interface{}) {
	if logFile != nil {
		logFile.WriteString("FATAL: ")
		logFile.WriteString(fmt.Sprintln(message...))
	}
	log.Fatal(message)
}
