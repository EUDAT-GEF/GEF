# -*- coding: utf-8 -*-
import os.path
import sys
import re
import requests
from subprocess import call
from requests.packages.urllib3.exceptions import InsecureRequestWarning
requests.packages.urllib3.disable_warnings(InsecureRequestWarning)

download_dir = '/volume'
hdl_prefix = 'http://hdl.handle.net/'
b2share_prefixes = ['https://b2share.eudat.eu/',
                    'https://trng-b2share.eudat.eu/']
pid_pattern = re.compile(r'\d+/[^/]+')


def main():
    if len(sys.argv) != 2:
        print("Usage: {} PID-or-URL".format(sys.argv[0]))
        exit(1)

    kind, url = analyze_and_resolve(sys.argv[1])
    if kind == 'b2share_record':
        for f in list_b2share_record_files(url):
            download_file(f['url'], local_filename=f['key'])
    else:
        download_file(url)



def analyze_and_resolve(url):
    is_pid = False
    if pid_pattern.match(url):
        is_pid = True
        url = '{}{}'.format(hdl_prefix, url)
    elif url.startswith(hdl_prefix):
        is_pid = True
    elif not (url.startswith('http://') or url.startswith('https://')):
        error("Argument is neither PID nor URL\n{}".format(url))

    if is_pid:
        res = requests.get(url, allow_redirects=False)
        if res.status_code < 200 or res.status_code >= 400:
            error("global handle system returned code {} for url {}".format(
                res.status_code, url))
        url = res.headers['Location']

    for prefix in b2share_prefixes:
        if url.startswith(prefix):
            if '/records/' in url:
                return 'b2share_record', url
            elif '/files/' in url:
                return 'b2share_file', url
            else:
                error("Cannot classify kind of b2share resource:\n{}".format(url))

    return 'url', url


def list_b2share_record_files(url):
    if '/api/records/' not in url:
        url = url.replace('/records/', '/api/records/')

    r = requests.get(url, headers={'Accept':'application/json'}, verify=False)
    if r.status_code != 200:
        error("Getting b2share record failed, status code: {}\n{}".format(
            r.status_code, r.text))
    links = r.json().get('links', {})
    file_bucket_url = links['files']
    r = requests.get(file_bucket_url, headers={'Accept':'application/json'}, verify=False)
    if r.status_code != 200:
        error("Getting b2share file bucket failed, status code: {}\n{}".format(
            r.status_code, r.text))

    for f in r.json()['contents']:
        f['url'] = f['links']['self']
        yield f


def download_file(url, local_filename=None):
    r = requests.get(url, stream=True, verify=False)
    disp = r.headers['content-disposition']
    fname = re.findall("filename=(.+)", disp, flags=re.IGNORECASE)
    if fname and fname[0]:
        fname = fname[0]
        fname = fname.strip('\'"')
    else:
        if local_filename:
            fname = local_filename
        else:
            fname = url.split('/')[-1]
    print('downloading', fname)
    with open(os.path.join(download_dir, fname), 'wb') as f:
        for chunk in r.iter_content(chunk_size=1024):
            if chunk: # filter out keep-alive new chunks
                f.write(chunk)
    return local_filename


def error(msg):
    print(msg)
    exit(1)


if __name__ == "__main__":
    main()
