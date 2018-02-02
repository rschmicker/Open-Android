
# http://ymsir.com/papers/pmds-iciss.pdf

import json
from os import listdir
from os.path import isfile, join
from constants import PERMISSIONS
from sklearn import svm
from sklearn import tree
from sklearn import linear_model
from sklearn.naive_bayes import GaussianNB
from sklearn.cluster import KMeans

def get_apk_json(filepath):
	json_files = [f for f in listdir(filepath) if isfile(join(filepath, f))]
	data = []
	for f in json_files:
		with open(filepath + f) as json_data:
			d = json.load(json_data)
			data.append(d)
	return data

def get_permissions(apk):
	perms = []
	for permission in PERMISSIONS:
		status = 1 if permission in apk['Permissions'] else 0
		perms.append(status)
	return perms

def main():
	model = linear_model.LogisticRegression()
	#model = svm.SVC()
	#model = GaussianNB()
	#model = KMeans(n_clusters = 10)
	filepath = "../data/output/"
	apks = get_apk_json(filepath)
	feature_vector = []
	target_vector = []
	for apk in apks:
		feature_vector.append(get_permissions(apk))
		target_type = 1 if apk['Malicious'] == 'true' else 0
		target_vector.append(target_type)
	clf = model
	clf.fit(feature_vector, target_vector)

	test_data = get_apk_json("../data/test/")[0]
	test_feature_vector = get_permissions(test_data)
	result = clf.predict([test_feature_vector])
	if result == 1:
		print("test data found to be malware")
	else:
		print("test data found to be benign")


main()
