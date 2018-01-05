#!/usr/bin/env python

import sys
from bs4 import BeautifulSoup
import requests
import hashlib
import time
import shutil

dl_loc = ""

def get_soup(url):
	res = requests.get(url)
	return BeautifulSoup(res.content, "html.parser")

def download(url):
	response = requests.get(url, stream=True)
	response.raise_for_status()
	with open('temp.apk', 'wb') as handle:
		for block in response.iter_content(1024):
			handle.write(block)
	filehash = hashlib.sha256(open("temp.apk", 'rb').read()).hexdigest()
	shutil.move("temp.apk", dl_loc + "/" + filehash + ".apk")
	print("Downloaded: " + str(url))
	time.sleep(10)

def crawl_site(url):
	soup = get_soup(url)
	a_tags = soup.findAll("a", href=True)
	for link in a_tags:
		if 'Parent' in str(link):
			continue
		elif '.apk' in str(link):
			to_append = link.get('href')
			download(url + to_append)
		else:
			next_link = link.get('href')
			print("Moving to: " + url + next_link)
			crawl_site(url + next_link)

def main():
	base_url = "http://apk-downloaders.com/apps/2018/"
	if(len(sys.argv) < 2):
		print("Usage: apk-downloaders <download dir>")
		sys.exit(1)
	global dl_loc
	dl_loc = sys.argv[1]
	crawl_site(base_url)	
main()
