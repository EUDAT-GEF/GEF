# -*- coding: utf-8 -*-

try:
    import os
except ImportError:
    print('Failed to import os library.')

try:
    import requests
except ImportError:
    print('Failed to import requests library.')

try:
    from subprocess import call
except ImportError:
    print('Failed to import call from subprocess library.')   

try:
    import string
except ImportError:
    print('Failed to import string library.')   

try:
    import random
except ImportError:
    print('Failed to import random library.') 

# The content string for the Python script that sets the target B2SHARE URL.

set_target_url_script_content = """# -*- coding: utf-8 -*-

try:
    import os
except ImportError:
    print('Failed to import os library.')

try:
    import argparse
except ImportError:
    print('Failed to import argparse library.')

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
    print('The current B2SHARE instance has been set to: %s' % target_url)"""

# The content string for the Python script that sets the access token.

set_access_token_script_content = """# -*- coding: utf-8 -*-

try:
    import os
except ImportError:
    print('Failed to import os library.')

try:
    import argparse
except ImportError:
    print('Failed to import argparse library.')

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
    print('The current access token has been set to: %s' % access_token)"""

# The content string for the Python script that lists all records in B2SHARE.

list_all_records_script_content = """# -*- coding: utf-8 -*-

try:
    import sys
except ImportError:
    print('Failed to import sys library')

try:
    import os
except ImportError:
    print('Failed to import os library.')

try:
    import urlparse
except ImportError:
    print('Failed to import urlparse library.')

try:
    import requests
except ImportError:
    print('Failed to import requests library.')

try:
    import json
except ImportError:
    print('Failed to import json library.')

try:
    import argparse
except ImportError:
    print('Failed to import argparse library.')

# Creates the parser for the argument string and sets the first optional positional argument to 'access_token' and
# the second optional argument to be 'target_url'.

arg_parser = argparse.ArgumentParser()
arg_parser.add_argument('--access_token', help = 'Specify the required token for accessing B2SHARE.')
arg_parser.add_argument('--target_url', help = 'Specify the URL of the B2SHARE instance to be addressed.')

# Parses the argument string and retrieves the access token specified by the user.

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

response = requests.get('%s/api/records' % target_url, params={'access_token':token, 'page_size':5, 'page_offset':2}, verify=False)

# Pretty-printing the response.

text_response = json.loads(response.text)

print json.dumps(text_response, indent=4)"""


# The content string for the Python script that lists all the records of a specific community from B2SHARE. 

list_community_records_script_content = """# -*- coding: utf-8 -*-

try:
    import sys
except ImportError:
    print('Failed to import sys library')

try:
    import os
except ImportError:
    print('Failed to import os library.')

try:
    import urlparse
except ImportError:
    print('Failed to import urlparse library.')

try:
    import requests
except ImportError:
    print('Failed to import requests library.')

try:
    import json
except ImportError:
    print('Failed to import json library.')

try:
    import argparse
except ImportError:
    print('Failed to import argparse library.')

# Creates the parser for the argument string and sets the first required positional argument to be 'community_name', 
# the second optional argument to be 'access_token', and the third optional argument to be 'target_url'.

arg_parser = argparse.ArgumentParser()
arg_parser.add_argument('community_name', help = 'Specify the name of the scientific community whose records are to be listed.')
arg_parser.add_argument('--access_token', help = 'Specify the required token for accessing B2SHARE.')
arg_parser.add_argument('--target_url', help = 'Specify the URL of the B2SHARE instance to be adressed.')

# Parses the argument string and retrieves the record id and the access token specified by the user.

args = arg_parser.parse_args()

community_name = args.community_name

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

response = requests.get('%s/api/records/%s' % (target_url, community_name), params={'access_token':token, 'page_size':5, 'page_offset':2}, verify=False)

# Pretty-printing the response.

text_response = json.loads(response.text)

print json.dumps(text_response, indent=4)""" 


# The content string for the Python script that lists a specific a record from B2SHARE.

list_specific_record_script_content = """# -*- coding: utf-8 -*-

try:
    import sys
except ImportError:
    print('Failed to import sys library')

try:
    import os
except ImportError:
    print('Failed to import os library.')

try:
    import urlparse
except ImportError:
    print('Failed to import urlparse library.')

try:
    import requests
except ImportError:
    print('Failed to import requests library.')

try:
    import json
except ImportError:
    print('Failed to import json library.')

try:
    import argparse
except ImportError:
    print('Failed to import argparse library.')

# Creates the parser for the argument string and sets the first required positional argument to be 'record_id' and 
# the second optional argument to be 'access_token' and the third optional argument to be 'target_url'.

arg_parser = argparse.ArgumentParser()
arg_parser.add_argument('record_id', help = 'Specify the id of the record to be listed.')
arg_parser.add_argument('--access_token', help = 'Specify the required token for accessing B2SHARE.')
arg_parser.add_argument('--target_url', help = 'Specify the URL of the B2SHARE instance to be adressed.')

# Parses the argument string and retrieves the record id and the access token specified by the user.

args = arg_parser.parse_args()

record_id = args.record_id

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

response = requests.get('%s/api/record/%s' % (target_url, record_id), params={'access_token': token}, verify=False)

# Pretty-printing the response.

text_response = json.loads(response.text)

print json.dumps(text_response, indent=4)"""


# The content string for the Python script that creates a new deposition in B2SHARE.

create_new_deposition_script_content = """# -*- coding: utf-8 -*-

try:
    import sys
except ImportError:
    print('Failed to import sys library')

try:
    import os
except ImportError:
    print('Failed to import os library.')

try:
    import urlparse
except ImportError:
    print('Failed to import urlparse library.')

try:
    import requests
except ImportError:
    print('Failed to import requests library.')

try:
    import json
except ImportError:
    print('Failed to import json library.')

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

response = requests.post('%s/api/depositions % target_url, params={'access_token': token}, verify=False)

# Pretty-printing the response.

text_response = json.loads(response.text)

print json.dumps(text_response, indent=4)"""


# The content string for the Python script that uploads a new file into a deposition object.

load_file_into_deposition_script_content = """# -*- coding: utf-8 -*-

try:
    import sys
except ImportError:
    print('Failed to import sys library')

try:
    import os
except ImportError:
    print('Failed to import os library.')

try:
    import urlparse
except ImportError:
    print('Failed to import urlparse library.')

try:
    import requests
except ImportError:
    print('Failed to import requests library.')

try:
    import json
except ImportError:
    print('Failed to import json library.')

try:
    import argparse
except ImportError:
    print('Failed to import argparse library.')

# Creates the parser for the argument string and sets the first required positional argument to be 'deposition_id',
# the second required argument to be 'filename', the name of the file to be uploaded, the third optional argument to be 'access_ token',
# and the fourth optional argument to be 'target_url'.

arg_parser = argparse.ArgumentParser()
arg_parser.add_argument('deposition_id', help = 'Specify the id of the deposition to load into.')
arg_parser.add_argument('filename', help = 'Specify the name of the file to be uploaded into the deposition.')
arg_parser.add_argument('--access_token', help = 'Specify the required token for accessing B2SHARE. ')
arg_parser.add_argument('--target_url', help = 'Specify the URL of the B2SHARE instance to be adressed.')

# Parses the argument string and retrieves the deposition id and the filename specified by the user.

args = arg_parser.parse_args()

deposition_id = args.deposition_id

filename = args.filename

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


# In order to hand the file to be uploaded to the requests library this dict is created.

file_to_be_uploaded = {'file' : open(filename, 'rb')}

# Accessing the B2SHARE instance through its API.

response = requests.post('%s/api/%s' % (target_url, deposition_id), files=file_to_be_uploaded, params={'access_token': token}, verify=False)

# Pretty-printing the response.

text_response = json.loads(response.text)

print json.dumps(text_response, indent=4)"""


# The content string for the Python script that lists the files uploaded into a deposition object.

list_files_uploaded_into_deposition_script_content = """# -*- coding: utf-8 -*-

try:
    import sys
except ImportError:
    print('Failed to import sys library')

try:
    import os
except ImportError:
    print('Failed to import os library.')

try:
    import urlparse
except ImportError:
    print('Failed to import urlparse library.')

try:
    import requests
except ImportError:
    print('Failed to import requests library.')

try:
    import json
except ImportError:
    print('Failed to import json library.')

try:
    import argparse
except ImportError:
    print('Failed to import argparse library.')

# Creates the parser for the argument string and sets the first required positional argument to be 'deposition_id',
# the second optional argument to be 'access_token', and the third optional argument to be 'target_url'.

arg_parser = argparse.ArgumentParser()
arg_parser.add_argument('deposition_id', help = 'Specify the id of the deposition to be listed.')
arg_parser.add_argument('--access_token', help = 'Specify the required token for accessing B2SHARE. ')
arg_parser.add_argument('--target_url', help = 'Specify the URL of the B2SHARE instance to be adressed.')

# Parses the argument string and retrieves the deposition id specified by the user.

args = arg_parser.parse_args()

deposition_id = args.deposition_id

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

response = requests.get('%s/api/%s/files' % (target_id, deposition_id), params={'access_token': token}, verify=False)

# Pretty-printing the response.

text_response = json.loads(response.text)

print json.dumps(text_response, indent=4)"""


# The content string for the Python script that commits a deposition, transforming it into an immutable object.

commit_deposition_script_content = """# -*- coding: utf-8 -*-

try:
    import sys
except ImportError:
    print('Failed to import sys library')

try:
    import os
except ImportError:
    print('Failed to import os library.')

try:
    import urlparse
except ImportError:
    print('Failed to import urlparse library.')

try:
    import requests
except ImportError:
    print('Failed to import requests library.')

try:
    import json
except ImportError:
    print('Failed to import json library.')

try:
    import argparse
except ImportError:
    print('Failed to import argparse library.')

try:
    import ast
except ImportError:
    print('Failed to import ast library.')

# Creates the parser for the argument string and sets the first required positional argument to be 'deposition_id',
# the second required argument to be the string defining the metadata dict 'metadata_dict_string', 
# the third optional argument to be 'access_token', and the fourth optional argument to be 'target_url'.

arg_parser = argparse.ArgumentParser()
arg_parser.add_argument('deposition_id', help = 'Specify the id of the deposition to be listed.')
arg_parser.add_argument('metadata_dict_string', help = 'Specify the metadata for the deposition as string defining a Python dict.')
arg_parser.add_argument('--access_token', help = 'Specify the required token for accessing B2SHARE. ')
arg_parser.add_argument('--target_url', help = 'Specify the URL of the B2SHARE instance to be adressed.')

# Parses the argument string and retrieves the deposition id, the metadata dict string and the access token specified by the user.

args = arg_parser.parse_args()

deposition_id = args.deposition_id

metadata_dict_string = args.metadata_dict_string()

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

# Creates the dict that describes the header for the requests post call. It's always the same.

metadata_header = {'content-type': 'application/json'}

# Creates the dict that describes the metadata for the requests post call.

record_metadata = ast.literal_eval(metadata_dict_string)

# Accessing the B2SHARE instance through its API.

response = requests.post('%s/api/%s/commit' % (target_url, deposition_id), data = json.dumps(record_metadata), headers = metadata_heade, params={'access_token': token}, verify=False)

# Pretty-printing the response.

text_response = json.loads(response.text)

print json.dumps(text_response, indent=4)"""


# Creates the Python script that sets a target B2SHARE instance

set_target_url_script_filename = 'set_target_url.py'

if not os.path.isabs(set_target_url_script_filename):
    abs_path_set_target_url_script = os.path.abspath(set_target_url_script_filename)
else:
    abs_path_set_target_url_script = set_target_url_script_filename

if not os.path.isfile(abs_path_set_target_url_script):
    file = open(abs_path_set_target_url_script, 'w')
    file.write(set_target_url_script_content)
    file.close()
else:
    print("A file with the name \'%s\' already exists!" % set_target_url_script_filename)

# Creates the Python script that sets the access token for B2SHARE

set_access_token_script_filename = 'set_access_token.py'

if not os.path.isabs(set_access_token_script_filename):
    abs_path_set_access_token_script = os.path.abspath(set_access_token_script_filename)
else:
    abs_path_set_access_token_script = set_access_token_script_filename

if not os.path.isfile(abs_path_set_access_token_script):
    file = open(abs_path_set_access_token_script, 'w')
    file.write(set_access_token_script_content)
    file.close()
else:
    print("A file with the name \'%s\' already exists!" % set_access_token_script_filename)

# Creates the Python script that lists all records from a B2SHARE instance

list_all_records_script_filename = 'list_all_records.py'

if not os.path.isabs(list_all_records_script_filename):
    abs_path_list_all_records_script = os.path.abspath(list_all_records_script_filename)
else:
    abs_path_list_all_records_script = list_all_records_script_filename

if not os.path.isfile(abs_path_list_all_records_script):
    file = open(abs_path_list_all_records_script, 'w')
    file.write(list_all_records_script_content)
    file.close()
else:
    print("A file with the name \'%s\' already exists!" % list_all_records_script_filename)

# Creates the Python script that lists all records of a specific community from a B2SHARE instance

list_community_records_script_filename = 'list_community_records.py'

if not os.path.isabs(list_community_records_script_filename):
    abs_path_list_community_records_script = os.path.abspath(list_community_records_script_filename)
else:
    abs_path_list_community_records_script = list_community_records_script_filename

if not os.path.isfile(abs_path_list_community_records_script):
    file = open(abs_path_list_community_records_script, 'w')
    file.write(list_community_records_script_content)
    file.close()
else:
    print("A file with the name \'%s\' already exists!" % list_community_records_script_filename)

# Creates the Python script that reads a record from B2SHARE

list_specific_record_script_filename = 'list_specific_record.py'

if not os.path.isabs(list_specific_record_script_filename):
    abs_path_list_specific_record_script = os.path.abspath(list_specific_record_script_filename)
else:
    abs_path_list_specific_record_script = list_specific_record_script_filename

if not os.path.isfile(abs_path_list_specific_record_script):
    file = open(abs_path_list_specific_record_script, 'w')
    file.write(list_specific_record_script_content)
    file.close()
else:
    print("A file with the name \'%s\' already exists!" % list_specific_record_script_filename)

# Creates the Python script that creates a new deposition in B2SHARE

create_new_deposition_script_filename = 'create_new_deposition.py'

if not os.path.isabs(create_new_deposition_script_filename):
    abs_path_create_new_deposition_script = os.path.abspath(create_new_deposition_script_filename)
else:
    abs_path_create_new_deposition_script = create_new_deposition_script_filename

if not os.path.isfile(abs_path_create_new_deposition_script):
    file = open(abs_path_create_new_deposition_script, 'w')
    file.write(create_new_deposition_script_content)
    file.close()
else:
    print("A file with the name \'%s\' already exists!" % create_new_deposition_script_filename)

# Creates the Python script that loads a new filo into a deposition.

load_file_into_deposition_script_filename = 'load_file_into_deposition.py'

if not os.path.isabs(load_file_into_deposition_script_filename):
    abs_path_load_file_into_deposition_script = os.path.abspath(load_file_into_deposition_script_filename)
else:
    abs_path_load_file_into_deposition_script = load_file_into_deposition_script_filename

if not os.path.isfile(abs_path_load_file_into_deposition_script):
    file = open(abs_path_load_file_into_deposition_script, 'w')
    file.write(load_file_into_deposition_script_content)
    file.close()
else:
    print("A file with the name \'%s\' already exists!" % load_file_into_deposition_script_filename)

# Creates the Python script that lists the files that were uploaded into a deposition.

list_files_uploaded_into_deposition_script_filename = 'list_files_uploaded_into_deposition.py'

if not os.path.isabs(list_files_uploaded_into_deposition_script_filename):
    abs_path_list_files_uploaded_into_deposition_script = os.path.abspath(list_files_uploaded_into_deposition_script_filename)
else:
    abs_path_list_files_uploaded_into_deposition_script = list_files_uploaded_into_deposition_script_filename

if not os.path.isfile(abs_path_list_files_uploaded_into_deposition_script):
    file = open(abs_path_list_files_uploaded_into_deposition_script, 'w')
    file.write(list_files_uploaded_into_deposition_script_content)
    file.close()
else:
    print("A file with the name \'%s\' already exists!" % list_files_uploaded_into_deposition_script_filename)

# Creates the Python script that commits the deposition to B2SHARE.

commit_deposition_script_filename = 'commit_deposition.py'

if not os.path.isabs(commit_deposition_script_filename):
    abs_path_commit_deposition_script = os.path.abspath(commit_deposition_script_filename)
else:
    abs_path_commit_deposition_script = commit_deposition_script_filename

if not os.path.isfile(abs_path_commit_deposition_script):
    file = open(abs_path_commit_deposition_script, 'w')
    file.write(commit_deposition_script_content)
    file.close()
else:
    print("A file with the name \'%s\' already exists!" % commit_deposition_script_filename)

# The content string for the Dockerfile that creates the image to access B2SHARE

dockerfile_content = """FROM ubuntu:latest
MAINTAINER Asela Rajapakse <asela.rajapakse@mpimet.mpg.de>
LABEL 'eudat.gef.service.name'='Image for accessing a B2SHARE instance though its HTTP API.'
LABEL 'eudat.gef.service.input_directory.path'='/input_directory'
LABEL 'eudat.gef.service.output_directory.path'='/output_directory'
VOLUME /input_directory /output_directory
RUN apt-get update
RUN apt-get --yes --allow-unauthenticated install apt-utils
RUN apt-get --yes --allow-unauthenticated install python
RUN apt-get --yes --allow-unauthenticated install python-pip
RUN pip install --upgrade pip && pip install requests
RUN apt-get --yes --allow-unauthenticated install python-lxml
COPY %s /scripts_for_user_invocation/
RUN chmod +x /scripts_for_user_invocation/%s
COPY %s /scripts_for_user_invocation/
RUN chmod +x /scripts_for_user_invocation/%s
COPY %s /scripts_for_user_invocation/
RUN chmod +x /scripts_for_user_invocation/%s
COPY %s /scripts_for_user_invocation/
RUN chmod +x /scripts_for_user_invocation/%s
COPY %s /scripts_for_user_invocation/
RUN chmod +x /scripts_for_user_invocation/%s
COPY %s /scripts_for_user_invocation/
RUN chmod +x /scripts_for_user_invocation/%s
COPY %s /scripts_for_user_invocation/
RUN chmod +x /scripts_for_user_invocation/%s
COPY %s /scripts_for_user_invocation/
RUN chmod +x /scripts_for_user_invocation/%s
COPY %s /scripts_for_user_invocation/
RUN chmod +x /scripts_for_user_invocation/%s\n""" % (set_target_url_script_filename, set_target_url_script_filename, 
    set_access_token_script_filename, set_access_token_script_filename, 
    list_all_records_script_filename, list_all_records_script_filename, 
    list_community_records_script_filename, list_community_records_script_filename, 
    list_specific_record_script_filename, list_specific_record_script_filename,
    create_new_deposition_script_filename, create_new_deposition_script_filename,
    load_file_into_deposition_script_filename, load_file_into_deposition_script_filename,
    list_files_uploaded_into_deposition_script_filename, list_files_uploaded_into_deposition_script_filename,
    commit_deposition_script_filename, commit_deposition_script_filename)



# Creating a Dockerfile that builds the image to access B2SHARE

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

# Removing temporary files.

os.remove(abs_path_set_access_token_script)
os.remove(abs_path_set_target_url_script)
os.remove(abs_path_list_all_records_script)
os.remove(abs_path_list_community_records_script)
os.remove(abs_path_list_specific_record_script)
os.remove(abs_path_create_new_deposition_script)
os.remove(abs_path_load_file_into_deposition_script)
os.remove(abs_path_list_files_uploaded_into_deposition_script)
os.remove(abs_path_commit_deposition_script)
os.remove(abs_path_dockerfile)
