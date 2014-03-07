#!/usr/bin/python
import os
import os.path
import sys
import string
import tempfile

DEBUG = 1
ICOMM_DIR = "/opt/iRODS/clients/icommands/bin/"
STAGEIN_SUFFIX = "IrodsPath"
STAGED_SUFFIX = "FilePath"
WORKFLOW_KEY = "workflowFilePath"
TAVERNA_EXECUTOR = "/opt/taverna-commandline-2.4.0/executeworkflow.sh"

stagedir = "."
stdoutFile = ".out"
stderrFile = ".err"
stagedFiles = []

class Logger(object):
    def __init__(self, terminal, filename):
        self.terminal = terminal
        self.log = open(filename, "a")

    def write(self, message):
        self.terminal.write(message)
        self.log.write(message)

    def flush(self):
        self.log.flush()

def raiseError(message):
    sys.stderr.write(message+"\n") 
    raise Exception(message)
    
def stageIn(irodsPath):
    stagedFiles.append(os.path.basename(irodsPath))
    command = ICOMM_DIR + "iget " + irodsPath
    if DEBUG: print "DEBUG: ", command
    os.system(command)

def stageOut(filePath, irodsTarget):
    command = ICOMM_DIR + "iput " + filePath + " " + irodsTarget
    if DEBUG: print "DEBUG: ", command
    os.system(command)

def runTavernaWorkflow(workflowFile, args):
    strlist = []
    for k, v in args.iteritems():
        if k == "datasetFilePath":
            strlist.append("-inputfile datasetFile " + v)
        else:
            strlist.append("-inputvalue " + k + " " + v)
    str = " ".join(strlist)
    command = TAVERNA_EXECUTOR + " " + workflowFile + " " + str + " >> " + stdoutFile + " 2>> " + stderrFile
    print "DEBUG: ", command
    os.system(command)


def setup(irodsGefCommandFilePath):
    global stagedir
    stagedir = tempfile.mkdtemp()
    if DEBUG: print "stagedir = " + stagedir
    gefFile = os.path.basename(irodsGefCommandFilePath)
    global stdoutFile
    global stderrFile
    stdoutFile = gefFile + ".out"
    stderrFile = gefFile + ".err"
    sys.stdout = Logger(sys.stdout, stagedir + "/" + stdoutFile)
    sys.stderr = Logger(sys.stderr, stagedir + "/" + stderrFile)
    os.chdir(stagedir);


def readGefCommandFile(irodsGefCommandFilePath): 
    stageIn(irodsGefCommandFilePath)
    gefCommandFileName = os.path.basename(irodsGefCommandFilePath)
    gefCommandFile = open(gefCommandFileName)
    args = []
    for line in gefCommandFile.readlines():
        tokens = string.split(line[:-1], "=", 1)
        if len(tokens) != 2:
            raiseError("unparseable gefcommand line: "+line)
        args.append((tokens[0], tokens[1]))
    return args


def stageInFiles(args):
    newargs = {}
    for var, value in args:
        if var.endswith(STAGEIN_SUFFIX):
            irodsPath = value
            stageIn(irodsPath)
            var = var[:-len(STAGEIN_SUFFIX)] + STAGED_SUFFIX
            value = stagedir + "/" + os.path.basename(irodsPath)
        newargs[var] = value;
    return newargs


def execute(args):
    if WORKFLOW_KEY in args:
        value = args[WORKFLOW_KEY]
        del args[WORKFLOW_KEY]
        if value.endswith("t2flow"):
            return runTavernaWorkflow(value, args)
        raiseError("unknown workflow type: " + value)
    raiseError("function unspecified, missing workflow argument")


def stageOutFiles(irodsGefCommandFilePath, args):
    irodsTarget = os.path.dirname(irodsGefCommandFilePath)
    paths = [f for f in os.listdir('.')]
    for p in paths:
        if os.path.isfile(p):
            if (not p in stagedFiles) and p != stdoutFile and p != stderrFile:
                stageOut(p, irodsTarget)
        if p == "Workflow1_output" and os.path.isdir(p):
            for f in os.listdir(p):
                stageOut(p+"/"+f, irodsTarget)
    sys.stdout.flush();
    stageOut(stdoutFile, irodsTarget)
    sys.stderr.flush();
    stageOut(stderrFile, irodsTarget)


def main(argv):
    if len(argv) != 2:
	raiseError("need one argument")
    setup(argv[1])
    try:
        args = readGefCommandFile(argv[1])
        args = stageInFiles(args)
        if DEBUG: print "DEBUG: ", args
        execute(args)
    finally:
        stageOutFiles(argv[1], args)

if __name__ == "__main__":
    main(sys.argv)

