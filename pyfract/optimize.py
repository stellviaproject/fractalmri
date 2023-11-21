import ctypes
from ctypes import cdll
import alloc
import cfg
import knn
from os import path
from alloc import lib

new_sample = lib.NewSample
new_sample.argtypes = [ctypes.POINTER(ctypes.c_char_p), ctypes.c_int, ctypes.POINTER(ctypes.c_char_p), ctypes.c_int, ctypes.c_double]
new_sample.restype = ctypes.c_int

class Sample(ctypes.Structure):
    _fields_ = [
        ("__id___", ctypes.c_int),
    ]
    def __init__():
        self.__id__ = ctypes.c_int(0)
    def __init__(self, fileName:str):
        err = load_sample(fileName.encode(), ctypes.byref(self.__id__))
        if err is not None:
            raise Exception(alloc.tostr(err))
    def __init__(self, part:float, tumors:list[str], notumors:list[str]):
        tumors_array_str = (ctypes.c_char_p * len(tumors))()
        notumors_array_str = (ctypes.c_char_p * len(notumors))()
        for i, ls_str in enumerate(tumors):
            tumors_array_str[i] = ls_str.encode()
        for i, ls_str in enumerate(notumors):
            notumors_array_str[i] = ls_str.encode()
        tumors_ptr = ctypes.POINTER(ctypes.c_char_p)(tumors_array_str)
        notumors_ptr = ctypes.POINTER(ctypes.c_char_p)(notumors_array_str)
        self.__id__ = new_sample(tumors_ptr, ctypes.c_int(len(tumors)), notumors_ptr, ctypes.c_int(len(notumors)), ctypes.c_double(part))
    def save(self, fileName:str):
        err = save_sample(fileName.encode(), self.__id__)
        if err is not None:
            raise Exception(alloc.tostr(err))
    def optimize(self, n:int, cfg:cfg.CFG)->'Sample':
        new_id = ctypes.c_int(0)
        err = optimize(self.__id__, ctypes.c_int(n), cfg, ctypes.byref(new_id))
        if err is not None:
            raise Exception(alloc.tostr(err))
        s = Sample()
        s.__id__ = new_id
        return s
    @property
    def points()->knn.PointList:
        pts_list = knn.PointList()
        pts_list.__listID__ = points(self.__id__)
        return pts_list
    def __del__(self):
        free_sample(self)
        
free_sample = lib.FreeSample
free_sample.argtypes = [Sample]

load_sample = lib.LoadSample
load_sample.argtypes = [ctypes.c_char_p, ctypes.POINTER(ctypes.c_int)]
load_sample.restype = ctypes.c_char_p

save_sample = lib.SaveSample
save_sample.argtypes = [ctypes.c_char_p, ctypes.c_int]
load_sample.restype = ctypes.c_char_p

optimize = lib.Optimize
optimize.argtypes = [ctypes.c_int, ctypes.c_int, cfg.CFG, ctypes.POINTER(ctypes.c_int)]
optimize.restype = ctypes.c_char_p

points = lib.Points
points.argtypes = [ctypes.c_int]
points.restype = ctypes.c_int