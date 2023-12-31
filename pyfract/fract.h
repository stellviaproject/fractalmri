/* Code generated by cmd/cgo; DO NOT EDIT. */

/* package fractalmri/export */


#line 1 "cgo-builtin-export-prolog"

#include <stddef.h>

#ifndef GO_CGO_EXPORT_PROLOGUE_H
#define GO_CGO_EXPORT_PROLOGUE_H

#ifndef GO_CGO_GOSTRING_TYPEDEF
typedef struct { const char *p; ptrdiff_t n; } _GoString_;
#endif

#endif

/* Start of preamble from import "C" comments.  */


#line 3 "export.go"
#include "export.h"

#line 1 "cgo-generated-wrapper"


/* End of preamble from import "C" comments.  */


/* Start of boilerplate cgo prologue.  */
#line 1 "cgo-gcc-export-header-prolog"

#ifndef GO_CGO_PROLOGUE_H
#define GO_CGO_PROLOGUE_H

typedef signed char GoInt8;
typedef unsigned char GoUint8;
typedef short GoInt16;
typedef unsigned short GoUint16;
typedef int GoInt32;
typedef unsigned int GoUint32;
typedef long long GoInt64;
typedef unsigned long long GoUint64;
typedef GoInt64 GoInt;
typedef GoUint64 GoUint;
typedef size_t GoUintptr;
typedef float GoFloat32;
typedef double GoFloat64;
#ifdef _MSC_VER
#include <complex.h>
typedef _Fcomplex GoComplex64;
typedef _Dcomplex GoComplex128;
#else
typedef float _Complex GoComplex64;
typedef double _Complex GoComplex128;
#endif

/*
  static assertion to make sure the file is being used on architecture
  at least with matching size of GoInt.
*/
typedef char _check_for_64_bit_pointer_matching_GoInt[sizeof(void*)==64/8 ? 1:-1];

#ifndef GO_CGO_GOSTRING_TYPEDEF
typedef _GoString_ GoString;
#endif
typedef void *GoMap;
typedef void *GoChan;
typedef struct { void *t; void *v; } GoInterface;
typedef struct { void *data; GoInt len; GoInt cap; } GoSlice;

#endif

/* End of boilerplate cgo prologue.  */

#ifdef __cplusplus
extern "C" {
#endif

extern __declspec(dllexport) char* NewImage(int width, int height, int len, int stride, uint8* data, struct RGBA* rgba);
extern __declspec(dllexport) char* ReadImage(char* fileName, struct RGBA* rgba);
extern __declspec(dllexport) void FreeImage(struct RGBA rgba);
extern __declspec(dllexport) char* WriteImage(char* fileName, struct RGBA img);
extern __declspec(dllexport) struct CFG GetCFG();
extern __declspec(dllexport) void FreeCFG(struct CFG cfg);
extern __declspec(dllexport) void Free(void* ptr);
extern __declspec(dllexport) void* Malloc(size_t size);
extern __declspec(dllexport) void* MemCopy(void* dst, void* src, size_t size);
extern __declspec(dllexport) char* SaveCFG(char* fileName, struct CFG in);
extern __declspec(dllexport) char* LoadCFG(char* fileName, struct CFG* out);
extern __declspec(dllexport) int NewPoints();
extern __declspec(dllexport) char* LoadPoints(char* fileName, int* listID);
extern __declspec(dllexport) char* SavePoints(char* fileName, int listID);
extern __declspec(dllexport) void FreePoints(int listID);
extern __declspec(dllexport) int NewImageSet();
extern __declspec(dllexport) void ImageSetAppend(int listID, struct RGBA img, char* mriLabel);
extern __declspec(dllexport) void FreeImageSet(int listID);
extern __declspec(dllexport) int NewFileSet();
extern __declspec(dllexport) void FileSetAppend(int listID, char* fileName, char* mriLabel);
extern __declspec(dllexport) void FreeFileSet(int listID);
extern __declspec(dllexport) int NewKNNFractal(struct CFG cfg, int listID);
extern __declspec(dllexport) void FreeKNNFractal(int knnID);
extern __declspec(dllexport) char* TrainWithImages(int knnID, int listID);
extern __declspec(dllexport) char* TrainWithFiles(int knnID, int listID);
extern __declspec(dllexport) char* Fit(int knnID, struct RGBA img, char** labelOut);
extern __declspec(dllexport) int GetPoints(int knnID);
extern __declspec(dllexport) int NewModel(struct CFG cfg);
extern __declspec(dllexport) void FreeModel(int modelID);
extern __declspec(dllexport) char* Eval(int modelID, struct RGBA img, int* evalID);
extern __declspec(dllexport) void FreeEval(int evalID);
extern __declspec(dllexport) int EvalLen(int evalID);
extern __declspec(dllexport) void FreeFractalDim(struct FractalDim fd);
extern __declspec(dllexport) void EvalFDAt(int evalID, int index, struct FractalDim* fdptr);
extern __declspec(dllexport) void EvalUmbralAt(GoInt evalID, int index, struct Umbral* umbralptr);
extern __declspec(dllexport) void EvalLogLogAt(GoInt evalID, int index, struct RGBA* rgba);
extern __declspec(dllexport) void EvalMFS(GoInt evalID, struct MFS* mfs);
extern __declspec(dllexport) void FreeMFS(struct MFS mfs);
extern __declspec(dllexport) double* EvalGetPoints(int evalID, int* n);
extern __declspec(dllexport) void FreeEvalPoints(double* ptr);
extern __declspec(dllexport) int NewSample(char** ctumors, int lenTumors, char** cnotumors, int lenNoTumors, double cpart);
extern __declspec(dllexport) void FreeSample(int sampleID);
extern __declspec(dllexport) char* LoadSample(char* fileName, int* sampleID);
extern __declspec(dllexport) char* SaveSample(char* fileName, int sampleID);
extern __declspec(dllexport) char* Optimize(int sampleID, int n, struct CFG cfg, int* new_sample_id);
extern __declspec(dllexport) int Points(int sampleID);

#ifdef __cplusplus
}
#endif
