## GEF rules: append this file (or move it) to:
## the irods rule file: server/config/reConfigs/core.re
## uses gefcommand binary (the go executor module) from 
acPostProcForPut {
    ON($objPath like "\*.gefcommand") {msiExecCmd("gefcommand", $objPath, "null", "null", "null", *out);}
}
