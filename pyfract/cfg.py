import ctypes
import os
from os import path
from ctypes import cdll
import rgba
from enum import Enum
import alloc
from alloc import lib

DIST_EUCLIDEAN = "euclidean"
DIST_MANHATTAN = "manhattan"
DIST_MINKOWSKI = "minkowski"
DIST_CHEBYSHEV = "chebyshev"
DIST_PEARSON_CORRELATION = "pearson-correlation"
DIST_HAMMING = "hamming"
    
SELECTOR_BINARY = "binary"
SELECTOR_MULTICLASS = "multi-class"
SELECTOR_SMOOTH_INVERSE = "smooth-inverse"

#representa la configuracion para el KNN
class KNN(ctypes.Structure):
    _fields_ = [
        ("__k__", ctypes.c_int),
        ("__distance__", ctypes.c_char_p),
        ("__selector__", ctypes.c_char_p),
        ("__minkowski_ratio__", ctypes.c_double),
        ("__weight_param__", ctypes.c_double),
        ("__smoothing_param__", ctypes.c_double),
    ]
    @property
    def k(self)->int:
        return int(self.__k__)
    @k.setter
    def k(self, value:int):
        self.__k__ = ctypes.c_int(value)
    @property
    def distance(self)->str:
        return alloc.tostr(self.__distance__)
    @distance.setter
    def distance(self,value:str):
        self.__distance__ = ctypes.c_char_p(value.encode())
    @property
    def selector(self)->str:
        return alloc.tostr(self.__selector__)
    @selector.setter
    def selector(self,value:str):
        self.__selector__ = ctypes.c_char_p(value.encode())
    @property
    def minkowski_ratio(self)->float:
        return float(self.__minkowski_ratio__)
    @minkowski_ratio.setter
    def minkowski_ratio(self,value:float):
        self.__minkowski_ratio__ = ctypes.c_double(value)
    @property
    def weight_param(self)->float:
        return float(self.__weight_param__)
    @weight_param.setter
    def weight_param(self,value:float):
        self.__weight_param__ = ctypes.c_double(value)
    @property
    def smoothing_param(self)->float:
        return float(self.__smoothing_param__)
    @smoothing_param.setter
    def smoothing_param(self,value:float):
        self.__smoothing_param__ = ctypes.c_double(value)

#representa la configuracion del algoritmo fractal y el KNN
class CFG(ctypes.Structure):
    _fields_ = [
        ("__knn__", KNN),
        ("__buffer__", ctypes.c_int),
        ("__parallel__", ctypes.c_int),
        ("__window_ratio__", ctypes.c_int),
        ("__denoiser_sigma_color__", ctypes.c_double),
        ("__denoiser_sigma_space__", ctypes.c_double),
        ("__denoiser_diameter__", ctypes.c_int),
        ("__denoiser_umbral_color__", ctypes.c_double),
        ("__min_umbral__", ctypes.c_double),
        ("__max_umbral__", ctypes.c_double),
        ("__min_area__", ctypes.c_int),
        ("__max_area__", ctypes.c_int),
        ("__ratio__", ctypes.c_int),
        ("__box_sizes__", ctypes.POINTER(ctypes.c_int)),
        ("__len_box_sizes__", ctypes.c_int),
        ("__log_log_width__", ctypes.c_int),
        ("__log_log_height__", ctypes.c_int),
        ("__umbral__", ctypes.POINTER(ctypes.c_double)),
        ("__len_umbral__", ctypes.c_int),
    ]
    @property
    def knn(self)->KNN:
        return self.__knn__
    @property
    def buffer(self)->int:
        return int(self.__buffer__)
    @buffer.setter
    def buffer(self,value:int):
        if value <= 0:
            raise Exception("buffer can not be lesser than one")
        self.__buffer__ = ctypes.c_int(value)
    @property
    def parallel(self)->int:
        return int(self.__parallel__)
    @parallel.setter
    def parallel(self,value:int):
        if value <= 0:
            raise Exception("number of thread can not be lesser than one")
        self.__parallel__ = ctypes.c_int(value)
    @property
    def window_ratio(self)->int:
        return int(self.__window_ratio__)
    @window_ratio.setter
    def window_ratio(self,value:int):
        if value <= 1:
            raise Exception("window ratio can not be lesser than one")
        self.__window_ratio__ = ctypes.c_int(value)
    @property
    def denoiser_sigma_color(self)->float:
        return float(self.__denoiser_sigma_color__)
    @denoiser_sigma_color.setter
    def denoiser_sigma_color(self,value:float):
        if value < 0.0:
            raise Exception("denoiser color can not be lesser than zero")
        self.__denoiser_sigma_color__ = ctypes.c_double(value)
    @property
    def denoiser_sigma_space(self)->float:
        return float(self.__denoiser_sigma_space__)
    @denoiser_sigma_space.setter
    def denoiser_sigma_space(self,value:float):
        if value < 0.0:
            raise Exception("denoiser sigma space can not be lesser than zero")
        self.__denoiser_sigma_space__ = ctypes.c_double(value)
    @property
    def denoiser_diameter(self)->int:
        return int(self.__denoiser_diameter__)
    @denoiser_diameter.setter
    def denoiser_diameter(self,value:int):
        if value < 1:
            raise Exception("denoiser diameter can not be lesser than zero")
        self.__denoiser_diameter__ = ctypes.c_int(value)
    
    @property
    def denoiser_umbral_color(self)->float:
        return float(self.__denoiser_umbral_color__)
    @denoiser_umbral_color.setter
    def denoiser_umbral_color(self,value:float):
        if value < 0.0:
            raise Exception("denoiser umbral color can not be lesser than zero")
        self.__denoiser_umbral_color__ = ctypes.c_double(value)
    @property
    def min_umbral(self)->float:
        return float(self.__min_umbral__)
    @min_umbral.setter
    def min_umbral(self,value:float):
        if value < 0.0:
            raise Exception("umbral can not be lesser than zero")
        self.__min_umbral__ = ctypes.c_double(value)
    @property
    def max_umbral(self)->float:
        return float(self.__max_umbral__)
    @max_umbral.setter
    def max_umbral(self,value:float):
        if value < 0.0:
            raise Exception("umbral can not be lesser than zero")
        self.__max_umbral__ = ctypes.c_double(value)
    @property
    def min_area(self)->int:
        return int(self.__min_area__)
    @min_area.setter
    def min_area(self,value:int):
        if value < 0:
            raise Exception("area can not be lesser than zero")
        self.__min_area__ = ctypes.c_int(value)
    @property
    def max_area(self)->int:
        return int(self.__max_area__)
    @max_area.setter
    def max_area(self,value:int):
        if value < 0:
            raise Exception("area can not be lesser than zero")
        self.__max_area__ = ctypes.c_int(value)
    @property
    def ratio(self)->int:
        return int(self.__ratio__)
    @ratio.setter
    def ratio(self,value:int):
        self.__ratio__ = ctypes.c_int(value)
    @property
    def box_sizes(self)->list[int]:
        sliceType = ctypes.c_int * int(self.__len_box_sizes__)
        box_sizes_array = sliceType.from_address(ctypes.cast(self.__box_sizes__, ctypes.c_void_p).value)
        return list[int](box_sizes_array)
    @box_sizes.setter
    def box_sizes(self, value: list[int]):
        if len(value) == 0:
            raise Exception("box sizes list needs to have at least one element")
        for size in value:
            if size < 1:
                raise Exception("box sizes can not be lesser than zero")
        # free previous pointer
        alloc.free(ctypes.cast(self.__box_sizes__, ctypes.c_void_p))
        # get python data as a pointer
        sliceType = ctypes.c_int * len(value)
        box_sizes_array = sliceType(*value)
        src_ptr = ctypes.cast(box_sizes_array, ctypes.c_void_p)
        # create a pointer to set go data
        dst_ptr = alloc.malloc(len(value)*4)
        # copy python to go
        alloc.memcpy(dst_ptr, src_ptr, len(value)*4)
        self.__box_sizes__ = ctypes.cast(dst_ptr, ctypes.POINTER(ctypes.c_int))
        self.__len_box_sizes__ = len(value)
    @property
    def loglog_width(self)->int:
        return int(self.__log_log_width__)
    @loglog_width.setter
    def loglog_width(self,value:int):
        if value < 0:
            raise Exception("the width of graphic can not be lesser than zero")
        self.__log_log_width__ = ctypes.c_int(value)
    @property
    def loglog_height(self)->int:
        return int(self.__log_log_height__)
    @loglog_height.setter
    def loglog_height(self,value:int):
        if value < 0:
            raise Exception("the height of graphic can not be lesser than zero")
        self.__log_log_height__ = ctypes.c_int(value)
    @property
    def umbral(self)->list[float]:
        sliceType = ctypes.c_double * int(self.__len_umbral__)
        umbral_array = sliceType.from_address(ctypes.cast(self.__umbral__, ctypes.c_void_p).value)
        return list[int](umbral_array)
    @umbral.setter
    def umbral(self, value: list[float]):
        if len(value) == 0:
            raise Exception("umbral list needs to have at least one element")
        for size in value:
            if size < 0.0:
                raise Exception("umbrals can not be lesser than zero")
        # free previous pointer
        alloc.free(ctypes.cast(self.__umbral__, ctypes.c_void_p))
        # prepare python data
        sliceType = ctypes.c_double * len(value)
        umbral_array = sliceType(*value)
        src_ptr = ctypes.cast(umbral_array, ctypes.c_void_p)
        # alloc new pointer
        dst_ptr = alloc.malloc(len(value)*8)
        # copy python to go
        alloc.memcpy(dst_ptr, src_ptr, len(value)*8)
        self.__umbral__ = ctypes.cast(dst_ptr, ctypes.POINTER(ctypes.c_double))
        self.__len_umbral__ = len(value)
    #libera la memoria de la configuracion
    def __del__(self):
        alloc.free(ctypes.cast(self.__box_sizes__, ctypes.c_void_p))
        alloc.free(ctypes.cast(self.__umbral__, ctypes.c_void_p))

#funcion para liberar la memoria de la configuracion
#func FreeCFG(CFG)
free_cfg = lib.FreeCFG
free_cfg.argtypes = [CFG]

#funcion para obtener la configuracion por defecto
#func GetCFG()Config
get_cfg = lib.GetCFG
get_cfg.restype = CFG #retorna la configuracion

#funcion para obtener la configuracion por defecto
def GetCFG()->CFG:
    return get_cfg() #obtiene la configuracion por defecto

#funcion para guardar la configuracion en un archivo json
#func SaveCFG(*char,CFG)*char
save_cfg = lib.SaveCFG
save_cfg.argtypes = [ctypes.c_char_p, CFG] #nombre del archivo, configuracion
save_cfg.restype = ctypes.c_char_p #algun error

#funcion para guardar la configuracion en un archivo json
def SaveCFG(fileName:str, cfg:CFG):
    #guardar la configuracion en un archivo
    err = save_cfg(fileName.encode(), cfg)
    #si hubo algun error lanzarlo
    if err is not None:
        raise Exception(alloc.tostr(err))

#funcion para cargar la configuracion desde un archivo json
#func LoadCFG(*char,*CFG)*char
load_cfg = lib.LoadCFG
load_cfg.argtypes = [ctypes.c_char_p, ctypes.POINTER(CFG)] #el nombre del archivo, un puntero a la configuracion
load_cfg.restype = ctypes.c_char_p #algun error

#funcion para cargar la configuracion desde un archivo json
def LoadCFG(fileName:str)->CFG:
    #crear la estructura para la configuracion
    cfg = CFG()
    #pasar el nombre del archivo y el puntero a la estructura
    err = load_cfg(fileName.encode(), ctypes.byref(cfg))
    #si hubo algun error
    if err is not None:
        raise Exception(alloc.tostr(err)) #convertir error a str y lanzar
    return cfg