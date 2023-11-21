#include <stdlib.h>
#include <string.h>

typedef unsigned char uint8;

struct RGBA {
	uint8 *Pix;
	int Len;
	int Stride;
	int Width;
	int Height;
};

typedef struct {
	int K;
	char *Distance;
	char *Selector;
	double MinkowskiRatio;
	double WeightParam;
	double SmoothingParam;
} KNN;

struct CFG {
	KNN KNN;
	int Buffer;
	int Parallel;
	int WindowRatio;
	double DenoiserSigmaColor;
	double DenoiserSigmaSpace;
	int DenoiserDiameter;
	double DenoiserUmbralColor;
	double MinUmbral;
	double MaxUmbral;
	int MinArea;
	int MaxArea;
	int Ratio;
	int *BoxSizes;
	int LenBoxSizes;
	int LogLogWidth;
	int LogLogHeight;
	double *Umbral;
	int LenUmbral;
};

struct DataPoint {
	int id;
	char *MRILabel;
};

struct FractalDim {
	double *LogSizes;
	double *LogMeasure;
	int Len;
	double FD;
};

// FractalDim* NewFractalDim(int len) {
// 	FractalDim *fd = (FractalDim*)malloc();
// 	fd->LogMeasure = (double*)malloc(sizeof(double)*len);
// 	fd->LogSizes = (int*)malloc(sizeof(int)*len);
// 	fd->Len = len;
// 	return fd;
// }

// void FreeFractalDim(FractalDim *fd) {
// 	free(fd->LogMeasure);
// 	free(fd->LogSizes);
// 	free(fd);
// }

struct Umbral {
	double Min;
	double Max;
};

struct MFS {
	double* Data;
	int Width;
	int Height;
};

// MFS* NewMFS(int width, int height) {
// 	double** data = (double**)malloc(height);
// 	for (int i = 0; i < height; i++) {
// 		*(data+i) = (double*)malloc(width);
// 	}
// 	MFS* mfs =(MFS*)malloc(sizeof(MFS));
// 	mfs->data = data;
// 	mfs->height = height;
// 	mfs->width = width;
// 	return mfs;
// }

// void FreeMFS(MFS *mfs) {
// 	int height = mfs->height;
// 	double** data = mfs->data;
// 	for (int i = 0; i < height; i++) {
// 		free(*(data+i));
// 	}
// 	free(mfs);
// }