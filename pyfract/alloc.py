import ctypes
from ctypes import cdll
from os import path
import platform

dll_file = path.dirname(path.abspath(__file__))

if platform.system() == 'Linux':
    dll_file = path.join(dll_file, 'fract.so')
elif platform.system() == 'Windows':
    dll_file = path.join(dll_file, 'fract.dll')
else:
    raise Exception('unsuported platform ' + platform.system())

#cargar la biblioteca gofract
lib = cdll.LoadLibrary(dll_file)

__free__ = lib.Free
__free__.argtypes = [ctypes.c_void_p]

__malloc__ = lib.Malloc
__malloc__.argtypes = [ctypes.c_ulonglong]
__malloc__.restype = ctypes.c_void_p

__memcpy__ = lib.MemCopy
__memcpy__.argtypes = [ctypes.c_void_p, ctypes.c_void_p, ctypes.c_ulonglong]
__memcpy__.restype = ctypes.c_void_p

def free(ptr:ctypes.c_void_p):
    __free__(ptr)
    
def malloc(len:int)->ctypes.c_void_p:
    return __malloc__(ctypes.c_ulonglong(len))

def memcpy(dst:ctypes.c_void_p, src:ctypes.c_void_p, len:int)->ctypes.c_void_p:
    return __memcpy__(dst, src, ctypes.c_ulonglong(len))

#obtiene un str de python a partir de un puntero char de C
def tostr(pstr:ctypes.c_char_p)->str:
    msg_str = ctypes.cast(pstr, ctypes.c_char_p).value.decode()
    #free cstring memory
    #free(ctypes.cast(pstr, ctypes.c_void_p))
    return msg_str