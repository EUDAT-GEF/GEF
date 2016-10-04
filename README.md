GEF-Docker
==========

The GEF backend for Docker


# How to compile
You should add the current directory to your GOPATH, then simply run
```
make build
```
The GEF binary executable should be generated in ``` /bin```

To generate a binary for linux, run
```
make build-linux
```

Then you can use the script ```onchanged-gef-docker.sh``` to monitor the changes, and copy the binary to target VM automatically.
