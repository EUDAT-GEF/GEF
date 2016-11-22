# -*- coding: utf-8 -*-

# This script is used by B2SHARE_access_image_builder.py. Please keep it in the same directory.

# create_new_deposition.py is a script that is part of a Docker image solution for accessing the B2SHARE HTTP API.
# It creates a new B2SHARE deposition.

import sys

try:
    import os
except ImportError:
    print('Failed to import os library. Exiting.')
    sys.exit()

try:
    import urlparse
except ImportError:
    print('Failed to import urlparse library. Exiting')
    sys.exit()

try:
    import requests

    # Importing urlib3 used in request package and disabling InsecureRequestWarnings that request throws when accessing B2SHARE.

    from requests.packages.urllib3.exceptions import InsecureRequestWarning
    
    requests.packages.urllib3.disable_warnings(InsecureRequestWarning)

except ImportError:
    print('Failed to import requests library. Exiting')
    sys.exit()

try:
    import json
except ImportError:
    print('Failed to import json library. Exiting')
    sys.exit()

try:
    import argparse
except ImportError:
    print('Failed to import argparse library. Exiting')
    sys.exit()

# Creates the parser for the argument string and sets the first optional positional argument to be 'access_token' and the second to be "target_url.

arg_parser = argparse.ArgumentParser()
arg_parser.add_argument('--access_token', help = 'Specify the required token for accessing B2SHARE.')
arg_parser.add_argument('--target_url', help = 'Specify the URL of the B2SHARE instance to be adressed.')

# Parses the argument string.

args = arg_parser.parse_args()

# The access token and the URL of the target B2SHARE instance are either given as arguments or retreived from files.

if args.access_token:   
    token = args.access_token
else:
    filename_access_token = 'access_token.txt'

    if os.path.isfile(filename_access_token):
        file = open(filename_access_token, 'r')
        token = file.read()
        file.close()
        print('Access token %s has been read from file.' % token)
    else:
        print('Calling the B2SHARE instance requires an access token.')
        sys.exit()

if args.target_url:
    target_url = args.target_url
else:
    filename_target_url = 'target_url.txt'

    if os.path.isfile(filename_target_url):
        file = open(filename_target_url, 'r')
        target_url = file.read()
        file.close()
        print('B2SHARE target URL %s has been read from file.' % target_url)
    else:
        print('Calling a B2SHARE instance requires and URL.')
        sys.exit()

# Accessing the B2SHARE instance through its API.

try:
    response = requests.post('%s/api/depositions' % target_url, params={'access_token': token}, verify=False)
except requests.exceptions.RequestException:
    print('Connection to B2SHARE host % could not be established.' % target_url)

# Pretty-printing the response. HTTP status code 201 is returned by B2SHARE if the new deposition has been created.

if (response.status_code == 201):
    print('Content type: '+response.headers['content-type'])
    text_response = json.loads(response.text)
    print json.dumps(text_response, indent=4)
else:
    print('That did not work as expected! Server returned HTTP status code %s.' % response.status_code)