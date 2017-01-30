# -*- coding: utf-8 -*-

# This script builds the Docker image that enables access to a B2SHARE instance. It does so by copying scripts that call
# the respective B2SHARE HTTP API functions into the image and by preparing the image for their invocation by installing the required dependencies.

import sys

try:
    import os
except ImportError:
    print('Failed to import os library. Exiting.')
    sys.exit()

try:
    from subprocess import call
except ImportError:
    print('Failed to import call from subprocess library. Exiting.')
    sys.exit()   

try:
    import string
except ImportError:
    print('Failed to import string library. Exiting.')
    sys.exit()   

try:
    import random
except ImportError:
    print('Failed to import random library. Exiting.')
    sys.exit()

try:
    import hashlib
except ImportError:
    print('Failed to import hashlib library. Exiting.')
    sys.exit()

# This dict contains the names of the scripts to be copied into the image and their respective sha224 hashes.

scripts_and_their_sha224_hashes = {
    "set_target_url.py":"26d5ad61dc5e4a86646e66b74d19a5009fbdc5964f61e01871f69ce9",
    "set_access_token.py":"cbf13fd3ff81a9ff963c1421e35307677f2721ebe5ff7e8d7ca4af96",
    "list_all_records.py":"649e4762a0c2b298a9ec1f0b6d0687c0f87caa43c7add4559930f278",
    "list_community_records.py":"72167bb36fa90cd282bb17dde22c15a80947f336bd54c1d681281943",
    "list_specific_record.py":"b77c3d5b0c55c68a6035657cc46ff04e6fe8a60a705a1af0be7683ef",
    "create_new_deposition.py":"910dd7542320d338c309c8d1610bb2677ef5e59c862f817204123f3f",
    "load_file_into_deposition.py":"0d2c97ee7249757d404af636e5d88dd0c9c805fa4820e6e0a8edc7a2",
    "list_files_uploaded_into_deposition.py":"2c81f260d8a281931ca5e0121e4cfe81d2aadfce7394afed1f09799a",
    "commit_deposition.py":"92625f9692c9dd56b5778c14a9316db3ae9665a9a7b0cea37237aec0"
}

# Checks if scripts with these names are available in the the same directory as the builder script. Does not continue if one is not.

for filename in scripts_and_their_sha224_hashes:
    if not os.path.isfile(filename):
        print('Script %s cannot be found. Aborting build process.' % filename)
        sys.exit()

# Checks if scripts with these names are available in the the same directory as the builder script.
# Also ensures that the files have not been tampered with before they are copied into the image by
# comparing sha24 hashes of the file contents with hashes created for verification. Continues only if this is true.

#for filename in scripts_and_their_sha224_hashes:
#
#    if os.path.isfile(filename):
#        file = open(filename, 'r')
#        file_content = file.read()
#        file.close()

#        # The file actully exists and its content is hashed with sha224.

#        file_hash = hashlib.sha224(file_content).hexdigest()
        
#        if not (scripts_and_their_sha224_hashes[filename] == file_hash):
#            print('Script %s has been altered. Aborting build process.' % filename)
#            sys.exit()
#    else:
#        print('Script %s cannot be found. Aborting build process.' % filename)
#        sys.exit()


# The content string for the Dockerfile that creates the image to access B2SHARE

dockerfile_content = """FROM ubuntu:latest
MAINTAINER Asela Rajapakse <asela.rajapakse@mpimet.mpg.de>
LABEL 'eudat.gef.service.name'='Image for accessing a B2SHARE instance though its HTTP API.'
LABEL 'eudat.gef.service.input_directory.path'='/input_directory'
VOLUME /input_directory
RUN apt-get update
RUN apt-get --yes --allow-unauthenticated install apt-utils
RUN apt-get --yes --allow-unauthenticated install python
RUN apt-get --yes --allow-unauthenticated install python-pip
RUN pip install --upgrade pip && pip install requests
RUN apt-get --yes --allow-unauthenticated install python-lxml\n"""

# Another content string for repeated RUN and COPY commands

more_dockerfile_content = ""

for filename in scripts_and_their_sha224_hashes:
    more_dockerfile_content += """COPY %s /scripts_for_user_invocation/
RUN chmod +x /scripts_for_user_invocation/%s\n""" % (filename, filename)

# These two are concatenated to form the content string for the Dockerfile.

dockerfile_content += more_dockerfile_content

# Creating a Dockerfile with a random name that builds the image to access B2SHARE.

empty_string = ""
random_string = empty_string.join(random.choice('abcdefghijklmnopqrs0123456789') for _ in range(20))

filename_dockerfile = random_string

if not os.path.isabs(filename_dockerfile):
    abs_path_dockerfile = os.path.abspath(filename_dockerfile)
else:
    abs_path_dockerfile = filename_dockerfile

if not os.path.isfile(abs_path_dockerfile):
    file = open(abs_path_dockerfile, 'w')
    file.write(dockerfile_content)
    file.close()
else:
    print("A file with the name \'%s\' already exists!" % filename_dockerfile)

# Building the Docker image with the Python scripts for access to EUDAT services.

call(['docker', 'build', '--no-cache=TRUE', '-t', 'b2share_access_image:latest', '-f', '%s' % filename_dockerfile, '.'])

# Removing Dockerfile.

os.remove(abs_path_dockerfile)
