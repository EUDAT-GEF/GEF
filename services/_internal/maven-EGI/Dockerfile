FROM zenika/alpine-maven

# Install bash to allow manual inspection i.e docker run -it $IMAGE_NAME /bin/bash
RUN apk add --update bash && rm -rf /var/cache/apk/* \
	apk upgrade && \
	apk add --update autoconf gcc

# Curl set up
RUN apk add --update curl && \
	rm -rf /var/cache/apk/* 

# Shell script to install the certificates allowed by EGI
COPY fedcloud-certificates.sh fedcloud-certificates.sh

# Make the sh file executable
RUN chmod +x fedcloud-certificates.sh

# Download the certificates
RUN ./fedcloud-certificates.sh

# Copy the VOMS project - use to generate the proxy
COPY VOMS-proxy /home/VOMS-proxy

RUN mvn -f /home/VOMS-proxy/ compile & mvn -f /home/VOMS-proxy/ package

# Copy the maven project inside the image
COPY workflowEGI /home/workflowEGI

# Compile and package the maven project
RUN mvn -f /home/workflowEGI/ compile & mvn -f /home/workflowEGI/ package

# Shell script which wrap the 2 command line
COPY run_java_workflow.sh /home/.

# Make the sh file executable
RUN chmod +x /home/run_java_workflow.sh

# Execute the command
CMD sh /home/run_java_workflow.sh