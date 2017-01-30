FROM python:3.6.0-alpine

MAINTAINER Emanuel Dima <emanueldima@gmail.com>

RUN pip install requests && mkdir /dataset && mkdir /status

WORKDIR /root
ADD . /root

LABEL "eudat.gef.service.name"="Volume Stage In"
LABEL "eudat.gef.service.description"="Stage In datasets specified by PIDs or URLs"
LABEL "eudat.gef.service.version"="0.1"

# when running this image, the main.py script takes the pid/url arguments
# and downloads the data into /dataset
CMD ["python3", "/root/main.py"]
