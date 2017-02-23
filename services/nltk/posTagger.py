import sys, os, nltk, glob

tmpBody = "Sentence; POS tags\n"
textFiles = glob.glob("/root/input/*.txt")

for fi in range(len(textFiles)):
    inputFile = open(textFiles[fi], "r")
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

