import sys, os, nltk, glob

tmpBody = "Sentence; POS tags\n"
textFiles = glob.glob("/root/input?/*.txt")

print("Starting to read files")
for fi in range(len(textFiles)):
    print(textFiles[fi])
    inputFile = open(textFiles[fi], "r")
    inputData = inputFile.read()
    inputData = inputData.strip()

    sentences = nltk.sent_tokenize(inputData)

    for si in range(len(sentences)):
        print("Sentence #"+str(si)+": "+sentences[si])
        sentenceTokens = nltk.word_tokenize(sentences[si])
        sentencePOSTags = nltk.pos_tag(sentenceTokens)
        tagsList = []
        for ti in range(len(sentencePOSTags)):
            tagsList.append("-".join(sentencePOSTags[ti]))

        tmpBody = tmpBody + sentences[si] + ";" + (" ".join(tagsList)) +"\n"

outputFile = open("/root/output/results.csv", "w+")
print("Writing parsing results into a file")
outputFile.write(tmpBody)
outputFile.close()
print("Done")

