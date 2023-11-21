import ctypes
import os
from os import path
from ctypes import cdll
import alloc
from alloc import lib
import numpy as np

#Representa una imagen en Golang
class RGBA(ctypes.Structure):
    #son los campos de la imagen
    _fields_ = [
        ("__pix__", ctypes.POINTER(ctypes.c_uint8)),
        ("__len__", ctypes.c_int),
        ("__stride__", ctypes.c_int),
        ("Width", ctypes.c_int),
        ("Height", ctypes.c_int),
    ]
    @property
    def buffer(self)->list[int]:
        sliceType = ctypes.c_uint8 * int(self.__len__)
        pix_array = sliceType.from_address(ctypes.cast(self.__pix__, ctypes.c_void_p).value)
        return list[int](pix_array)
    @property
    def as_np(self)->np.ndarray:
        sliceType = ctypes.c_uint8 * int(self.__len__)
        buffer = sliceType.from_address(ctypes.cast(self.__pix__, ctypes.c_void_p).value)
        np_array = np.frombuffer(buffer, dtype=np.uint8, count=self.len)
        np_array = np_array.reshape((self.width, self.height, 4))
        return np_array
    @property
    def width(self)->int:
        return self.Width
    @property
    def height(self)->int:
        return self.Height
    @property
    def stride(self)->int:
        return self.__stride__
    @property
    def len(self)->int:
        return self.__len__
    #libera la memoria de la imagen
    def __del__(self):
        free_image(self)

#funcion para liberar la memoria de la imagen    
free_image = lib.FreeImage
free_image.argtypes = [RGBA]

#funcion para leer una imagen Golang desde un archivo
read_image = lib.ReadImage
#func ReadImage(*char,*RGBA)*char
read_image.argtypes = [ctypes.c_char_p, ctypes.POINTER(RGBA)]
read_image.restype = ctypes.c_char_p #retorna un error

#funcion para leer una imagen desde un archivo
def ReadImage(fileName:str)->RGBA:
    image = RGBA() #crear la imagen
    #pasar el nombre del archivo y el puntero a la imagen
    err = read_image(fileName.encode(), ctypes.byref(image))
    #si no hay error
    if err is None:
        return image #retornar la imagen
    #lanzar el error
    raise Exception(alloc.tostr(err))

#funcion para escribir una imagen en un archivo desde Golang
write_image = lib.WriteImage
#func WriteImage(*char,RGBA)*char
write_image.argtypes = [ctypes.c_char_p, RGBA]
write_image.restype = ctypes.c_char_p #retorna un error

#funcion para escribir una imagen en un archivo
def WriteImage(fileName:str, image:RGBA):
    #pasar el nombre del archivo y la imagen
    err = lib.WriteImage(fileName.encode(), image)
    #si no hubo error retornar
    if err is None:
        return
    #lanzar el error
    raise Exception(alloc.tostr(err))

new_image = lib.NewImage
new_image.argtypes = [ctypes.c_int32, ctypes.c_int32, ctypes.c_int32, ctypes.c_int32,ctypes.POINTER(ctypes.c_uint8)]
new_image.restype = ctypes.c_char_p

def NewImage(img:np.ndarray) -> RGBA:
    img_uint8 = img.astype(np.uint8).ctypes.data_as(ctypes.POINTER(ctypes.c_uint8))
    stride = img.strides[0]
    width, height, _ = img.shape
    image = RGBA()
    err = lib.NewImage(width, height, width*height*4, stride, img_uint8, ctypes.byref(image))
    if err is None:
        return image
    raise Exception(alloc.tostr(err))