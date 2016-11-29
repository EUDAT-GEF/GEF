FROM busybox:latest
MAINTAINER Emanuel Dima <emanueldima@gmail.com>

LABEL "eudat.gef.service.name"="Clone"
LABEL "eudat.gef.service.description"="Copy input to output"
LABEL "eudat.gef.service.version"="0.1"
LABEL "eudat.gef.service.input.1.name"="Input Directory"
LABEL "eudat.gef.service.input.1.path"="/mydata/input"
LABEL "eudat.gef.service.output.1.name"="Output Directory"
LABEL "eudat.gef.service.output.1.path"="/mydata/output"

COPY clone /clone
RUN chmod +x /clone
CMD ["/clone"]

# COPY taverna-commandline-core-2.5.0-linux_amd64.deb /taverna-commandline-core-2.5.0-linux_amd64.deb
# RUN dpkg -i taverna-commandline-core-2.5.0-linux_amd64.deb
# CMD ["/../taverna-command-line", "-i", "something.tflow "]