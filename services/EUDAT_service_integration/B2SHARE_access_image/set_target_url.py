# -*- coding: utf-8 -*-

# This script is used by B2SHARE_access_image_builder.py. Please keep it in the same directory.

# set_target_url.py is a script that is part of a Docker image solution for accessing the B2SHARE HTTP API.
# It sets the URL of the B2SHARE instance to address.

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

# Creates the parser for the argument string and sets the first required positional argument to be the URL of the target B2SHARE instance. 

arg_parser = argparse.ArgumentParser()
arg_parser.add_argument('target_url', help = 'Specify the URL of the B2SHARE instance to adress. ')

# Parses the argument string and retrieves the URL specified by the user.

args = arg_parser.parse_args()

target_url = args.target_url

# Writes the target URL into a file.

filename_target_url = 'target_url.txt'

if os.path.isfile(filename_target_url):
    print('A target B2SHARE instance has already been set.')
    file = open(filename_target_url, 'r')
    target_url = file.read()
    file.close()
    print('It is at: %s' % target_url)

else:
    file = open(filename_target_url, 'w')
    file.write(target_url)
    file.close()
    print('The current B2SHARE instance has been set to: %s' % target_url)