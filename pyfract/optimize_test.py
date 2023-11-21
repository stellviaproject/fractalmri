import optimize
import os
import cfg
from os import path

def filter(files:list[str])->list:
    nomasks = []
    for fileName in files:
        if not fileName.endswith("-mask.png"):
            nomasks.append(fileName)
    return nomasks

def join_all(parent:str, files:list[str]):
    i = 0
    while i < len(files):
        files[i] = path.join(parent, files[i])
        i += 1

tumors = filter(os.listdir("../gofract/lib/tumors"))
notumors = filter(os.listdir("../gofract/lib/notumors"))
join_all("../gofract/lib/tumors", tumors)
join_all("../gofract/lib/notumors", notumors)

s = optimize.Sample(0.8, tumors, notumors)
s = s.optimize(10, cfg.GetCFG())
s.save("./sample.json")
