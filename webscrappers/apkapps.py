#!/usr/bin/env python

import os
import sys
import requests
import time
import shutil
import base64
import hashlib

dl_loc = ""

def download(url):
	response = requests.get(url, stream=True)
	response.raise_for_status()
	with open('temp4.apk', 'wb') as handle:
		for block in response.iter_content(1024):
			handle.write(block)
	if(os.path.getsize('temp4.apk') < 1024):
		print("invalid file... continuing...")
		time.sleep(10)
		return
	filehash = hashlib.sha256(open("temp4.apk", 'rb').read()).hexdigest()
	shutil.move("temp4.apk", dl_loc + "/" + filehash + ".apk")
	print("Downloaded: " + str(url))
	time.sleep(10)

def crawl_site(url):
	for i in range(10000):
		e_num = base64.b64encode(str(i))
		download(url + e_num)

def main():
	base_url = "https://apkapps.com/apps/download/"
	if(len(sys.argv) < 2):
		print("Usage: apk-downloaders <download dir>")
		sys.exit(1)
	global dl_loc
	dl_loc = sys.argv[1]
	crawl_site(base_url)	
main()
