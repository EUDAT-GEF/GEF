import requests
import json
from time import sleep

# Configuration parameters
GEFAddress = "https://localhost:8443" # The GEF host
JobStartEndpoint = "/api/jobs"
VolumesEndpoint = "/api/volumes"
accessToken = "ICAyzW_P8tWWt6-DsCLRlHFtKvqo7dLGqy7Dl-_e6WJ34kSG" # Access token generated in the UI (profile page)
NLTKServiceID = "11edd34a-1096-4759-a08a-f76a3d3ab751" # Service ID of the NLTK demo service (can be found in the UI)
ServiceInput1 = "http://hdl.handle.net/11304/0591b2ed-d5c6-4007-bb99-6b473f3f07fb" # Can be any URL pointing to an English text
ServiceInput2 = "some text to be parsed" # Any text fragment in English

# Starting a job
urlVars = {'access_token': accessToken}
formData = {'serviceID': NLTKServiceID, 'pid_input0': ServiceInput1, 'pid_input1': ServiceInput2}
response = requests.post(GEFAddress + JobStartEndpoint, params = urlVars, data = formData, verify=False) # Certificate verification is OFF, because of the self-signed certificates
jsonResponse = json.loads(response.text)

runningJobID = jsonResponse["jobID"]
if runningJobID:
    print("A job has been started -> " + runningJobID)

# Monitoring a job's status
jobStatusCode = -1
jobOutputVolumeID = ""
while jobStatusCode==-1:
    response = requests.get(GEFAddress + JobStartEndpoint + "/" + runningJobID, verify=False)
    jsonResponse = json.loads(response.text)
    jobStatusCode = jsonResponse["Job"]["State"]["Code"]

    print("Job is running...")
    if jsonResponse["Job"]["OutputVolume"]:
        if len(jsonResponse["Job"]["OutputVolume"])>0:
            jobOutputVolumeID = jsonResponse["Job"]["OutputVolume"][0]["VolumeID"]
    sleep(0.1) # Setting a delay between requests

print("The job has been finished")

# Inspecting and downloading the output
if len(jobOutputVolumeID)>0:
    print("Inspecting the output volume ->" + jobOutputVolumeID) # We need only the first one
    response = requests.get(GEFAddress + VolumesEndpoint + "/" + jobOutputVolumeID + "/", params = urlVars, verify=False)
    jsonResponse = json.loads(response.text)
    if len(jsonResponse["volumeContent"])>0:
        print("Downloading the first file from the output volume")
        outputFileName = jsonResponse["volumeContent"][0]["name"]
        with open(outputFileName, 'wb') as f:
            resp = requests.get(GEFAddress + VolumesEndpoint + "/" + jobOutputVolumeID + "/" + outputFileName + "?content", params = urlVars, verify=False)
            f.write(resp.content)
        print("File has been downloaded -> " + outputFileName)
    else:
        print("No output files were found")
