import ctypes
import os
from os import path
from ctypes import cdll
import rgba
import cfg
import alloc
from alloc import lib

class FractalDim(ctypes.Structure):
    _fields_ = [
        ("__log_sizes__", ctypes.POINTER(ctypes.c_double)),
        ("__log_measure__", ctypes.POINTER(ctypes.c_double)),
        ("__len__", ctypes.c_int),
        ("__fd__", ctypes.c_double),
    ]
    @property
    def len(self)->int:
        return int(self.__len__)
    @property
    def fd(self)->float:
        return float(self.__fd__)
    @property
    def logsizes(self)->list[float]:
        sliceType = ctypes.c_double * int(self.__len__)
        log_size_array = sliceType.from_address(ctypes.cast(self.__log_sizes__, ctypes.c_void_p).value)
        return list[int](log_size_array)
    @property
    def logmeasures(self)->list[float]:
        sliceType = ctypes.c_double * int(self.__len__)
        log_measure_array = sliceType.from_address(ctypes.cast(self.__log_measure__, ctypes.c_void_p).value)
        return list[float](log_measure_array)
    def __del__(self):
        free_fractal_dim(self)
        
free_fractal_dim = lib.FreeFractalDim
free_fractal_dim.argtypes = [FractalDim]
        
class Umbral(ctypes.Structure):
    _fields_ = [
        ("__min__", ctypes.c_double),
        ("__max__", ctypes.c_double),
    ]
    @property
    def min(self)->float:
        return float(self.__min__)
    @property
    def max(self)->float:
        return float(self.__max__)

class MFS(ctypes.Structure):
    _fields_ = [
        ("__data__", ctypes.POINTER(ctypes.POINTER(ctypes.c_double))),
        ("__width__", ctypes.c_int),
        ("__height__", ctypes.c_int),
    ]
    @property
    def width(self)->int:
        return int(self.__width__)
    @property
    def height(self)->int:
        return int(self.__height__)
    @property
    def array(self)->list[list[float]]:
        doubleSlice = ctypes.c_double * (self.width * self.height)
        mfs_array = doubleSlice.from_address(ctypes.cast(self.__data__, ctypes.c_void_p).value)
        mfs_list = list[float](mfs_array)
        mfs = list[list[float]]()
        width = self.width
        for i in range(0, self.height, 1):
            mfs.append(mfs_array[i*width:(i+1)*width])
        return mfs
    def __del__(self):
        free_mfs(self)
    
free_mfs = lib.FreeMFS
free_mfs.argtypes = [MFS]


class ModelEvaluation:
    def __init__(self):
        self.__evalID__ = ctypes.c_int(0)
    @property
    def len(self)->int:
        return eval_len(self.__evalID__)
    def fdAt(self, index:int)->FractalDim:
        fd = FractalDim()
        eval_fd_at(self.__evalID__, ctypes.c_int(index), ctypes.byref(fd))
        return fd
    def umbralAt(self, index:int)->Umbral:
        umbral = Umbral()
        eval_umbral_at(self.__evalID__, ctypes.c_int(index), ctypes.byref(umbral))
        return umbral
    def loglogAt(self, index:int)->rgba.RGBA:
        image = rgba.RGBA()
        eval_loglog_at(self.__evalID__, ctypes.c_int(index), ctypes.byref(image))
        return image
    def mfs(self)->MFS:
        mfs = MFS()
        eval_mfs(self.__evalID__, ctypes.byref(mfs))
        return mfs
    @property
    def points(self)->list[float]:
        n = ctypes.c_int(0)
        points_ptr = eval_get_points(self.__evalID__, ctypes.byref(n))
        sliceType = ctypes.c_double * int(n.value)
        array = sliceType.from_address(ctypes.cast(points_ptr, ctypes.c_void_p).value)
        free_eval_points(points_ptr)
        return list[float](array)
    def __del__(self):
        free_eval(self.__evalID__)
        
free_eval = lib.FreeEval
free_eval.argtypes = [ctypes.c_int]

eval_len = lib.EvalLen
eval_len.argtypes = [ctypes.c_int]
eval_len.restype = ctypes.c_int

eval_fd_at = lib.EvalFDAt
eval_fd_at.argtypes = [ctypes.c_int, ctypes.c_int, ctypes.POINTER(FractalDim)]
    
eval_umbral_at = lib.EvalUmbralAt
eval_umbral_at.argtypes = [ctypes.c_int, ctypes.c_int, ctypes.POINTER(Umbral)]

eval_loglog_at = lib.EvalLogLogAt
eval_loglog_at.argtypes = [ctypes.c_int, ctypes.c_int, ctypes.POINTER(rgba.RGBA)]

eval_mfs = lib.EvalMFS
eval_mfs.argtypes = [ctypes.c_int, ctypes.POINTER(MFS)]

eval_get_points = lib.EvalGetPoints
eval_get_points.argtypes = [ctypes.c_int, ctypes.POINTER(ctypes.c_int)]
eval_get_points.restype = ctypes.POINTER(ctypes.c_double)

free_eval_points = lib.FreeEvalPoints
free_eval_points.argtypes = [ctypes.POINTER(ctypes.c_double)]

class Model:
    def __init__(self, cfg:cfg.CFG):
        self.__modelID__ = new_model(cfg)
    def eval(self, image:rgba.RGBA)->ModelEvaluation:
        evalID = ctypes.c_int(0)
        err = model_eval(self.__modelID__, image, ctypes.byref(evalID))
        if err is not None:
            raise Exception(alloc.tostr(err))    
        evaluation = ModelEvaluation()
        evaluation.__evalID__ = evalID
        return evaluation
    def __del__(self):
        free_model(self.__modelID__)

new_model = lib.NewModel
new_model.argtypes = [cfg.CFG]
new_model.restype = ctypes.c_int

free_model = lib.FreeModel
free_model.argtypes = [ctypes.c_int]

model_eval = lib.Eval
model_eval.argtypes = [ctypes.c_int, rgba.RGBA, ctypes.POINTER(ctypes.c_int)]
model_eval.restype = ctypes.c_char_p
