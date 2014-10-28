# GEF rules

#### triggered by a rule in server/config/reConfigs/core.re, like this:
#### uses gefcommand binary (the go executor module) from 
## acPostProcForPut {
##        ON($objPath like "\*.gefcommand") {msiExecCmd("gefcommand", $objPath, "null", "null", "null", *out);}
## }

processGefWorkflowFile(*cmdPath) {
	logDebug("processGefWorkflowFile(*cmdPath)");
	processTavernaWorkflow();
}

processTavernaWorkflow() {
	logDebug("processTavernaWorkflow(*cmdPath)");
    msiExecCmd("/home/irods/taverna-commandline-2.4.0/executeworkflow.sh", "*serverID*path", "null", "null", "null", *out);
    msiGetStdoutInExecCmdOut(*out, *message);
}
