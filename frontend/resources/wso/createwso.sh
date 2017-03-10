LOCAL_WSO = "./w.wso"
IRODS_WSO = "/ironZone/home/rods/wso/wso.mss"
IRODS_WORKFLOW_COLLECTION = "/ironZone/home/rods/wso/instance"

# 1. put the iRODS workflow file in iRODS
iput -D  "msso file" $LOCAL_WSO $IRODS_WSO

# 2. link the wso file with a staging collection
imkdir $IRODS_WORKFLOW_COLLECTION
imcoll -m msso $IRODS_WSO $IRODS_WORKFLOW_COLLECTION

# Using the workflow: 
# put the parameters file and any other data
#   iput w.mpf /zone/path/workflow/w.mpf
# execute the workflow (by "getting" the run file)
#   iget /zone/path/workflow/w.run -
# get the result of the execution
#   ils /zone/path/workflow/
