# -*- coding: utf-8 -*-

# This script is used by B2SHARE_access_image_builder.py. Please keep it in the same directory.

# set_access_token.py is a script that is part of a Docker image solution for accessing the B2SHARE HTTP API.
# It sets the access token used to authenticate with the selected B2SHARE instance.

import sys

try:
    import os
except ImportError:
    print('Failed to import os library. Exiting.')
    sys.exit()

try:
    import argparse
except ImportError:
    print('Failed to import argparse library. Exiting.')
    sys.exit()

# Creates the parser for the argument string and sets the first required positional argument to be 'access_token'. 

arg_parser = argparse.ArgumentParser()
arg_parser.add_argument('access_token', help = 'Specify the access token string.')

# Parses the argument string and retrieves the access token specified by the user.

args = arg_parser.parse_args()

access_token = args.access_token

# Writes the access_token into a file.

filename_access_token = 'access_token.txt'

if os.path.isfile(filename_access_token):
    print('An access token has already been set.')
    file = open(filename_access_token, 'r')
    access_token = file.read()
    file.close()
    print('It is: %s' % access_token)

else:
    file = open(filename_access_token, 'w')
    file.write(access_token)
    file.close()
    print('The current access token has been set to: %s' % access_token)