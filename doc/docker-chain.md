Proposal for docker chain in GEFx
=======

# 1. Requirements
The GEF should be able to run a chain of docker containers.
The chain is supplied by users.
The GEF can run the chains in sequence. Any failure executing a container will abort the execution of the whole chain.

# 2. Notes
A chain of docker containers can use named docker volume to exchange data.
A example of data exchange is demostrated below

```
docker volume create
66b45e625b566574dab8e8b574ce393e54a0bfbe4976ea328ca84335b5770374

docker volume ls
DRIVER              VOLUME NAME
local               66b45e625b566574dab8e8b574ce393e54a0bfbe4976ea328ca84335b5770374

docker volume inspect  66b45e625b566574dab8e8b574ce393e54a0bfbe4976ea328ca84335b5770374
[
    {
        "Name": "66b45e625b566574dab8e8b574ce393e54a0bfbe4976ea328ca84335b5770374",
        "Driver": "local",
        "Mountpoint": "/var/lib/docker/volumes/66b45e625b566574dab8e8b574ce393e54a0bfbe4976ea328ca84335b5770374/_data",
        "Labels": {},
        "Scope": "local"
    }
]

docker run  --rm -v 66b45e625b566574dab8e8b574ce393e54a0bfbe4976ea328ca84335b5770374:/data alpine /bin/sh -c "echo 'this is a test' > /data/test.txt"

docker run  --rm -v 66b45e625b566574dab8e8b574ce393e54a0bfbe4976ea328ca84335b5770374:/data1 alpine /bin/sh -c "cat /data1/test.txt"
this is a test
```

