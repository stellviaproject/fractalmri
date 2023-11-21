import ctypes
import os
from os import path
from ctypes import cdll
import rgba
import alloc
import cfg
from alloc import lib

class PointList:
    def __init__(self):
        self.__listID__ = new_points()
    def SavePoints(self,fileName:str):
        err = save_points(fileName.encode(), self.__listID__)
        if err is not None:
            raise Exception(alloc.tostr(err))
    def __del__(self):
        free_points(self.__listID__)

new_points = lib.NewPoints
new_points.restype = ctypes.c_int

#funcion para liberar los puntos
#func FreePoints(int)
free_points = lib.FreePoints
free_points.argtypes = [ctypes.c_int] #el id de la lista de puntos

#funcion para cargar los puntos desde un archivo
#func LoadPoints(*char,*int)*char
load_points = lib.LoadPoints
load_points.argtypes = [ctypes.c_char_p, ctypes.POINTER(ctypes.c_int)] #nombre del archivo, puntero a id
load_points.restype = ctypes.c_char_p #algun error

#funcion para cargar los puntos desde un archivo
def LoadPoints(fileName:str)->PointList:
    pointList = PointList()
    err = load_points(fileName.encode(), ctypes.byref(pointList.__listID__))
    if err is None:
        return pointList
    raise Exception(alloc.tostr(err))

save_points = lib.SavePoints
save_points.argtypes = [ctypes.c_char_p, ctypes.c_int]
save_points.restype = ctypes.c_char_p

class ImageSet:
    def __init__(self):
        self.__setID__ = new_image_set()
    def __del__(self):
        free_image_set(self.__setID__)
    def append(self, image:rgba.RGBA, label:str):
        image_set_append(self.__setID__, image, label.encode())
        
new_image_set = lib.NewImageSet
new_image_set.restype = ctypes.c_int

free_image_set = lib.FreeImageSet
free_image_set.argtypes = [ctypes.c_int]

image_set_append = lib.ImageSetAppend
image_set_append.argtypes = [ctypes.c_int, rgba.RGBA, ctypes.c_char_p]

class FileSet:
    def __init__(self):
        self.__setID__ = new_file_set()
    def append(self, fileName:str, label:str):
        file_set_append(self.__setID__, fileName.encode(), label.encode())
    def __del__(self):
        free_file_set(self.__setID__)

new_file_set = lib.NewFileSet
new_file_set.restype = ctypes.c_int

file_set_append = lib.FileSetAppend
file_set_append.argtypes = [ctypes.c_int, ctypes.c_char_p, ctypes.c_char_p]

free_file_set = lib.FreeFileSet
free_file_set.archivo = [ctypes.c_int]

class KNNFractal:
    def __init__(self, cfg:cfg.CFG, points:PointList):
        self.__knnID__ = new_knn_fractal(cfg, points.__listID__)
    def train_images(self, image_set:ImageSet):
        err = train_with_images(self.__knnID__, image_set.__setID__)
        if err is not None:
            raise Exception(alloc.tostr(err))
    def train_files(self, file_set:FileSet):
        err = train_with_files(self.__knnID__, file_set.__setID__)        
        if err is not None:
            raise Exception(alloc.tostr(err))
    def fit(self, image:rgba.RGBA)->str:
        label = ctypes.c_char_p()
        err = fit(self.__knnID__, image, ctypes.byref(label))
        if err is not None:
            raise Exception(alloc.tostr(err))
        return alloc.tostr(label)
    def points(self)->PointList:
        points = PointList()
        points.__listID__ = get_points(self.__knnID__)
        return points
    def __del__(self):
        free_knn_fractal(self.__knnID__)
        
new_knn_fractal = lib.NewKNNFractal
new_knn_fractal.argtypes = [cfg.CFG, ctypes.c_int]
new_knn_fractal.restype = ctypes.c_int

free_knn_fractal = lib.FreeKNNFractal
free_knn_fractal.restype = ctypes.c_int

train_with_images = lib.TrainWithImages
train_with_images.argtypes = [ctypes.c_int, ctypes.c_int]
train_with_images.restype = ctypes.c_char_p

train_with_files = lib.TrainWithFiles
train_with_files.argtypes = [ctypes.c_int, ctypes.c_int]
train_with_files.restype = ctypes.c_char_p

fit = lib.Fit
fit.argtypes = [ctypes.c_int, rgba.RGBA, ctypes.POINTER(ctypes.c_char_p)]
fit.restype = ctypes.c_char_p

get_points = lib.GetPoints
get_points.argtypes = [ctypes.c_int]
get_points.restype = ctypes.c_int


#img = rgba.ReadImage("test/1.png")
#buffer = img.buffer
#rgba.WriteImage("test/2.png", img)

# cfg = GetCFG()
# SaveCFG("./cfg.json", cfg)
# cfg = LoadCFG("./cfg.json")
# SaveCFG("./cfg-cmp.json", cfg)

# points = LoadPoints("gofract/lib/points.json")
# points.SavePoints("gofract/lib/copy.json")

# imgSet = ImageSet()

