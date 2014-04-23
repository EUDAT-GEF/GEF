#Input parameters: *StageDir, *File1, *Filter
gefWorkflow {
	msiExecCmd("geffilter.sh", "*StageDir *File1 *Filter", "null","null","null", *Result1);
	msiGetFormattedSystemTime(*myTime,"human","%d-%d-%d %ldh:%ldm:%lds");
	msiGetStdoutInExecCmdOut(*Result1,*Out);
	writeLine("stdout", *Out);
	writeLine("stdout", "Workflow Executed Successfully at *myTime");
}
