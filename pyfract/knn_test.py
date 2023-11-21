def filter(files:list[str])->list:
    nomasks = []
    for fileName in files:
        if not fileName.endswith("-mask.png"):
            nomasks.append(fileName)
    return nomasks

tumors = filter(os.listdir("./gofract/lib/tumors"))
notumors = filter(os.listdir("./gofract/lib/notumors"))

fileSet = FileSet()
tumorsPath = "./gofract/lib/tumors"
notumorsPath = "./gofract/lib/notumors"
for tumor in tumors:
    fileSet.append(path.join(tumorsPath, tumor), "tumor")
for notumor in notumors:
    fileSet.append(path.join(notumorsPath, notumor), "notumor")

knn = KNNFractal(GetCFG(), PointList())
knn.train_files(fileSet)