
# http://ymsir.com/papers/pmds-iciss.pdf
# Pmds: Permission-based malware detection system

import json
from os import listdir
from os.path import isfile, join
from constants import PERMISSIONS
from sklearn import svm
from sklearn import tree
from sklearn import linear_model
from sklearn.naive_bayes import GaussianNB
from sklearn.cluster import KMeans
from sklearn.model_selection import KFold
from sklearn.tree import DecisionTreeClassifier
from sklearn.ensemble import RandomForestClassifier

# dl: 50 minutes for perms, filesize, and malicious for 94,000 apks, 12.9KB/s
def get_apk_json(filepath):
	d = {}
	with open(filepath) as json_data:
		d = json.load(json_data)
	return d

def get_permissions(apk):
	perms = []
	for permission in PERMISSIONS:
		status = 1 if permission in apk['permissions'] else 0
		perms.append(status)
	return perms

def main():
	models = {}
	models["SimpleLogistic"] = linear_model.LogisticRegression() # simple logistic
	models["SMO"] = svm.SVC() # SMO
	models["NaiveBayes"] = GaussianNB() # naive bayes
	models["RandomTree"] = DecisionTreeClassifier(random_state=0) # random tree
	models["RandomForest10"] = RandomForestClassifier(max_depth=10, random_state=0) # random forest, max depth: 10
	models["RandomForest50"] = RandomForestClassifier(max_depth=50, random_state=0) # random forest, max depth: 50
	models["RandomForest100"] = RandomForestClassifier(max_depth=100, random_state=0) # random forest, max depth: 100

	apks = get_apk_json("permsfilesize.json")
	feature_vector = []
	target_vector = []
	apks = apks['data']

	for model_name, model in models:
		kf = KFold(10, True, None)
		kf.get_n_splits(apks)
		train_apks = []
		predict_apks = []
		for train, test in kf.split(apks):
			train_apks = train
			predict_apks = test
			break
		print("training...")
		for idx in train_apks:
			feature_vector.append(get_permissions(apks[idx]))
			target_type = 1 if apks[idx]['Malicious'] == 'true' else 0
			target_vector.append(target_type)
		clf = model
		clf.fit(feature_vector, target_vector)

		total_test = len(predict_apks)
		number_correct = 0
		print("predicting...")
		for idx in predict_apks:
				test_feature_vector = get_permissions(apks[idx])
				result = clf.predict([test_feature_vector])
				mal_status = 1 if apks[idx]['Malicious'] == 'true' else 0
				if result == mal_status:
						number_correct += 1
		percent = (float(number_correct)/float(total_test))*float(100)
		print("Model: " + model_name)
		print("Accuracy: " + str(percent) + "%")

main()
