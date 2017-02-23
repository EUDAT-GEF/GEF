from __future__ import division
import sys, os, nltk, numpy


model_dict = {}
for model_item in open("model.csv", "r"):
	model_item = model_item.strip()
	model_list = model_item.split("\t")
	if model_list[0] not in model_dict:
		model_dict[model_list[0]] = {}
		
	model_dict[model_list[0]][model_list[2]] = model_list[1]
	



#print(model_dict)

tmp_body = "QUESTION;COARSE_CLASS;FINE_CLASS;CLASSIFIED_AS;CORRECT\n"
tags_freq = {}
n_gram = 4
class_dict = {}
question_caunter = 0
for anvil_item in open("test.csv", "r"):
	
	anvil_item = anvil_item.strip()
	
	line_list = anvil_item.split("\t")
	question = line_list[3]
	question_caunter = question_caunter + 1
	label=""
	text = nltk.word_tokenize(question)
	tags_set = nltk.pos_tag(text)
	tagged_question = ""
	tags_list = []
	for t in range(len(tags_set)):
		if len(tags_set[t])>1:
			tagged_question = tagged_question + tags_set[t][0] + "/" + tags_set[t][1] + " "
			tags_list.append(tags_set[t][1])
			
	
	tmp_ngram = []
	ngram_list = []
	for tl in range(len(tags_list)):
		if (tl+(n_gram-1))<=(len(tags_list)-1):
			for ng in range(n_gram):
				tmp_ngram.append(tags_list[ng+tl])
			#ngram_list.append(tmp_ngram)
			ngram_list.append("_".join(tmp_ngram))
			
			# Calcuating frequences for n-gramms for all questions
			if "_".join(tmp_ngram) not in tags_freq:
				tags_freq["_".join(tmp_ngram)] = 1
			else:
				tags_freq["_".join(tmp_ngram)] = tags_freq["_".join(tmp_ngram)] + 1
			
			tmp_ngram = []
		
	guestion_ngrams = ""
	for nl in range(len(ngram_list)):
		ngram_item = "_".join(ngram_list[nl])
		guestion_ngrams = guestion_ngrams + ngram_item + " * "
	question = tagged_question
	
	rel_name = ""
	if len(line_list)>=6:
		if len(line_list[5])>0:
			helper = "eatHumanDescription"
			rel_name = line_list[5].lower()
	
	if len(line_list)>=7:	
		if len(line_list[6])>0:
			helper = "eatHumanRelation"
			rel_name = line_list[6].lower()
			
	if len(line_list)>=8:
		if len(line_list[7])>0:
			helper = "eatHumanGroups"
			rel_name = line_list[7].lower()
	
	if len(line_list)>=9:	
		if len(line_list[8])>0:
			helper = "eatEntities"
			rel_name = line_list[8].lower()
	
	if len(line_list)>=10:	
		if len(line_list[9])>0:
			helper = "eatDescription"
			rel_name = line_list[9].lower()
	
	if len(line_list)>=11:		
		if len(line_list[10])>0:
			helper = "eatLocation"
			rel_name = line_list[10].lower()
	
	if len(line_list)>=12:	
		if len(line_list[11])>0:
			helper = "eatTime"
			rel_name = line_list[11].lower()
	
	if len(rel_name)>0:
		print(question)
		print(ngram_list)
		print()
		rel_dict = {}
		for nl in range(len(ngram_list)):
			if ngram_list[nl] in model_dict:
				
				for r in model_dict[ngram_list[nl]]:
					model_dict[ngram_list[nl]]
				
				if r not in rel_dict:
					rel_dict[r] = float(model_dict[ngram_list[nl]][r])
				else:
					rel_dict[r] = rel_dict[r] + float(model_dict[ngram_list[nl]][r])
		
		if len(rel_dict)>0:
			max_rel = 0
			max_label = ""
			for mr in rel_dict:
				if rel_dict[mr]>max_rel:
					max_rel = rel_dict[mr]
					max_label = mr
			
			if max_label == rel_name:
				is_correct = "yes"
			else:
				is_correct = "no"		
			tmp_body = tmp_body + question + ";" + helper + ";" + rel_name + ";" + max_label + ";" + is_correct +"\n"
			#tmp_body = tmp_body + guestion_ngrams + ";" + "" + ";" + "" + "\n"
		#print(rel_dict)
		
	
		




rf = open("results.csv","w+")
rf.write(tmp_body)
rf.close()	

