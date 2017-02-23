import sys, os, nltk

tmpBody = "Sentence; POS tags\n"

inputFile = open("/root/input/text.txt", "r")
inputData = inputFile.read()
inputData = inputData.strip()

sentences = nltk.sent_tokenize(inputData)

for si in range(len(sentences)):
	sentenceTokens = nltk.word_tokenize(sentences[si])
	sentencePOSTags = nltk.pos_tag(sentenceTokens)
	tagsList = []
	for ti in range(len(sentencePOSTags)):
	    tagsList.append("-".join(sentencePOSTags[ti]))

	tmpBody = tmpBody + sentences[si] + ";" + (" ".join(tagsList)) +"\n"

outputFile = open("/root/output/results.csv", "w+")
outputFile.write(tmpBody)
outputFile.close()

