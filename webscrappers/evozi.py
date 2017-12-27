#!/usr/bin/env python3
import selenium
from selenium import webdriver
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.support.ui import WebDriverWait
import time
import sys
import urllib
import hashlib
import os
import subprocess
import requests

if len(sys.argv) < 2:
	print("evozi.py <file of package names>")
	sys.exit(1)

package_file = sys.argv[1]
package_list = []

file = open(package_file, "r")
for line in file:
	package_list.append(line)

url = "https://apps.evozi.com/apk-downloader/"

profile = webdriver.FirefoxProfile()
profile.set_preference('browser.download.folderList', 2) # custom location
profile.set_preference('browser.download.manager.showWhenStarting', False)
profile.set_preference('browser.download.dir', '/tmp/')
profile.set_preference('browser.helperApps.neverAsk.saveToDisk', '*')

driver = webdriver.Firefox(profile)
for package_name in package_list:
	driver.get(url)
	time.sleep(5)
	package_name_box = driver.find_element_by_class_name("form-control")
	package_name_box.send_keys(package_name)
	submit_button = driver.find_element_by_class_name("btn").click()
	time.sleep(5)
	dl_button = driver.find_element_by_class_name("btn-success")
	link = dl_button.get_attribute("href")
	#link = "http" + link[5:]
	print(link)
	response = requests.get(link, stream=True) #, verify=False)
	response.raise_for_status()
	with open('temp.apk', 'wb') as handle:
	    for block in response.iter_content(1024):
	        handle.write(block)
	filehash = hashlib.sha256(open("temp.apk", 'rb').read()).hexdigest()
	os.rename("temp.apk", filehash + ".apk")
	print("Downloaded: " + package_name)
driver.quit()