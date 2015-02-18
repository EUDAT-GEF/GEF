The GEF specification
=====================


# 1. Requirements

The GEF is required to:

1.1. Expose a fixed set of predefined data processing services

- e.g. data transfers, filtering; see §3.2
- if the input of a service is specified by a PID, the GEF must resolve that PID
  * the PID could be an irods URN; it should be resolved correctly
- it is expected that the indirectly referred data sets are located near the service; the GEF can refuse the processing request otherwise, returning an appropriate error message

1.2. Accept a Docker container image and install it as a service

- the container image must follows certain conventions (e.g. execution paths, file data paths), see §4

1.3. Be secure

- Mangage user authentication and authorization using the EUDAT Unity service



# 2. Architecture Overview

The GEF infrastructure should be based on:
- a set of GEF endpoints, accessible at various URLs, 
- and a Request dispatcher, to route general requests to the most efficient endpoint

## 2.1. A typical endpoint 

2.1.0. Frontend: browser UI / user HTTP calls

2.1.1. Frontend web service: implements the APIs, see §3

- should also resolve PIDs and stage or map the datasets to the backends
- should route calls to the appropriate backend

2.1.2. Backends:

  - 2.1.2.1. simple Docker containers
    - good as prototype but ultimately unscalable

  - 2.1.2.2. map/reduce: Hadoop cluster
    - using PigLatin

  - 2.1.2.3. using a cluster manager: Kubernetes on Mesos?
    - will solve the scheduling problem by delegating it to the cluster manager

  - 2.1.2.4. streaming data? 

TODO: all the backends must connect to iRODS, as it's expected that most data will be stored in iRODS. How do we do that? Do we stage in/out data? Mounting would be a much better choice, but there are rumors that it's unstable.


## 2.2. Request dispatcher

Provide a single service with unique URL (e.g. `http://eudat.eu/gef`) to route a request (by returning HTTP 307 Temporary Redirect) to the best endpoint from a data locality perspective.

## 2.3. Services for other data transfers

See the Lightweight Replication Service / EUDAT HTTP API



# 3. APIs

## 3.1. API Protocols

3.1.1. The GEF will support the WPS protocol for discovery and execution of its functions. 

3.1.2. The GEF will also support a simple, RESTful API w/ json, primarily for the use of the web based UI. It will be less stable than WPS, changing rapidly during UI construction.

- TODO: should the REST UI API be part of the spec?


## 3.2. API Functions

These functions will be exposed as WPS processes/REST resources:

### 3.2.0. Data transfer functions

- based on URLs and HTTP GET method

### 3.2.1. Fixed functions

3.2.1.1. Filter

- input parameters: 
  - `dataPID` or `dataURI`
  - `queryType`: depends on endpoint/community
  - `query`: depends on `queryType`

3.2.1.2. Map-Reduce script

- input params:
  - `dataPID` or `dataURI`
  - `scriptType`: only `PigLatin` is acceptable for now
  - `script`: a string representing a PigLatin script


### 3.2.2. Management of dynamic functions

3.2.2.1. New service

- create a new service by uploading a Docker image
- the new service becomes a new WPS process/REST resource
- TODO: details

3.2.2.1. Delete service

- delete an existing service
- TODO: details



# 4. Docker image specification

As described in §3.2.2.1, a Docker image can be uploaded to the GEF and become a service. When invoking this service the following should happen:

4.1. The input data must be mapped to the container file system:

- 4.1.1. input data specified by PID will be mapped to `/data/PID/ACTUAL_PID`, read-only
- 4.1.2. output data is expected to be written to the `/data/output/` folder

The container should not rely on external resources (databases, web services) as any such request will probably fail.



# 5. Implementation

5.1. module §2.1.1: the frontend web service:

- is currently implemented in Java
- resolves PID calls to dataset urls
- mounts datasets, stage in needed data
- runs executioner/container
- returns output data

5.2. backend §2.1.2.1: 

- now in development

