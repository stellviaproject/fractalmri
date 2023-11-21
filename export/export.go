package main

//#include "export.h"
import "C"
import (
	"fractalmri/gofract/cfg"
	"fractalmri/gofract/lib"
	"fractalmri/gofract/pipes"
	"image"
	"unsafe"
)

func main() {}

//export NewImage
func NewImage(width, height, len, stride C.int, data *C.uint8, rgba *C.struct_RGBA) (err_msg *C.char) {
	defer func() {
		if e := recover(); e != nil {
			err := e.(error)
			err_msg = C.CString(err.Error())
		}
	}()
	rgba.Pix = (*C.uint8)(C.malloc(C.size_t(len)))
	rgba.Width = width
	rgba.Height = height
	rgba.Len = len
	rgba.Stride = stride
	C.memcpy(unsafe.Pointer(rgba.Pix), unsafe.Pointer(data), C.size_t(len))
	return
}

//export ReadImage
func ReadImage(fileName *C.char, rgba *C.struct_RGBA) *C.char {
	img, err := lib.ReadImage(C.GoString(fileName))
	if err != nil {
		return C.CString(err.Error())
	}
	*rgba = ImageToC(img)
	return nil
}

func ImageToC(img *image.RGBA) C.struct_RGBA {
	var rgba C.struct_RGBA
	// Crear un búfer de memoria para Pix
	size := C.size_t(len(img.Pix))
	rgba.Pix = (*C.uint8)(C.malloc(size))

	// Copiar los datos del puntero de Go al búfer de memoria
	C.memcpy(unsafe.Pointer(rgba.Pix), unsafe.Pointer(&img.Pix[0]), size)

	rgba.Stride = C.int(img.Stride)
	rgba.Len = C.int(len(img.Pix))
	rgba.Width = C.int(img.Rect.Dx())
	rgba.Height = C.int(img.Rect.Dy())
	return rgba
}

//export FreeImage
func FreeImage(rgba C.struct_RGBA) {
	C.free(unsafe.Pointer(rgba.Pix))
}

//export WriteImage
func WriteImage(fileName *C.char, img C.struct_RGBA) *C.char {
	rgba := ToGoImage(img)
	// Llama a la función WriteImage en el paquete lib para escribir la imagen
	err := lib.WriteImage(C.GoString(fileName), rgba)
	if err != nil {
		return C.CString(err.Error())
	}
	return nil
}

func ToGoImage(img C.struct_RGBA) *image.RGBA {
	rgba := &image.RGBA{
		Pix:    (*[1 << 30]uint8)(unsafe.Pointer(img.Pix))[:img.Len:img.Len],
		Stride: int(img.Stride),
		Rect:   image.Rect(0, 0, int(img.Width), int(img.Height)), // Ajusta las dimensiones según corresponda
	}
	return rgba
}

//export GetCFG
func GetCFG() C.struct_CFG {
	return CFGToC(cfg.GetCFG())
}

func CFGToC(cfg cfg.Configuration) C.struct_CFG {
	var out C.struct_CFG
	//BoxSizes
	out.BoxSizes = (*C.int)(C.malloc(C.size_t(4 * len(cfg.BoxSizes))))
	C.memcpy(unsafe.Pointer(out.BoxSizes), unsafe.Pointer(&toInt32(cfg.BoxSizes)[0]), C.size_t(len(cfg.BoxSizes)*4))
	out.LenBoxSizes = C.int(len(cfg.BoxSizes))
	//Buffer
	out.Buffer = C.int(cfg.Buffer)
	//parallel
	out.Parallel = C.int(cfg.Parallel)
	//denoiser
	out.DenoiserDiameter = C.int(cfg.DenoiserDiameter)
	out.DenoiserSigmaColor = C.double(cfg.DenoiserSigmaColor)
	out.DenoiserSigmaSpace = C.double(cfg.DenoiserSigmaSpace)
	out.DenoiserUmbralColor = C.double(cfg.DenoiserUmbralColor)
	//knn
	out.KNN.K = C.int(cfg.KNN.K)
	out.KNN.Distance = C.CString(string(cfg.KNN.Distance))
	out.KNN.Selector = C.CString(string(cfg.KNN.Selector))
	out.KNN.MinkowskiRatio = C.double(cfg.KNN.MinkowskiRatio)
	out.KNN.SmoothingParam = C.double(cfg.KNN.SmoothingParam)
	out.KNN.WeightParam = C.double(cfg.KNN.WeightParam)
	//Umbral
	out.MinUmbral = C.double(cfg.MinUmbral)
	out.MaxUmbral = C.double(cfg.MaxUmbral)
	out.Ratio = C.int(cfg.Ratio)
	//Umbral
	out.Umbral = (*C.double)(C.malloc(C.size_t(8 * len(cfg.Umbral))))
	C.memcpy(unsafe.Pointer(out.Umbral), unsafe.Pointer(&cfg.Umbral[0]), C.size_t(8*len(cfg.Umbral)))
	out.LenUmbral = C.int(len(cfg.Umbral))
	//Area
	out.MinArea = C.int(cfg.MinArea)
	out.MaxArea = C.int(cfg.MaxArea)
	//LogLog
	out.LogLogHeight = C.int(cfg.LogLogHeight)
	out.LogLogWidth = C.int(cfg.LogLogWidth)
	//WindowRatio
	out.WindowRatio = C.int(cfg.WindowRatio)
	return out
}

func toInt32(slice []int) []int32 {
	out := make([]int32, len(slice))
	for i := 0; i < len(slice); i++ {
		out[i] = int32(slice[i])
	}
	return out
}

func CfgToGo(in C.struct_CFG) cfg.Configuration {
	boxSizes := make([]int, in.LenBoxSizes)
	boxSzN32 := make([]int32, in.LenBoxSizes)
	C.memcpy(unsafe.Pointer(&boxSzN32[0]), unsafe.Pointer(in.BoxSizes), C.size_t(4*in.LenBoxSizes))
	for i := 0; i < len(boxSzN32); i++ {
		boxSizes[i] = int(boxSzN32[i])
	}

	umbral := make([]float64, in.LenUmbral)
	C.memcpy(unsafe.Pointer(&umbral[0]), unsafe.Pointer(in.Umbral), C.size_t(8*in.LenUmbral))

	return cfg.Configuration{
		KNN: cfg.KNN{
			K:              int(in.KNN.K),
			Distance:       cfg.Distance(C.GoString(in.KNN.Distance)),
			Selector:       cfg.Selector(C.GoString(in.KNN.Selector)),
			MinkowskiRatio: float64(in.KNN.MinkowskiRatio),
			SmoothingParam: float64(in.KNN.SmoothingParam),
			WeightParam:    float64(in.KNN.WeightParam),
		},
		Buffer:              int(in.Buffer),
		Parallel:            int(in.Parallel),
		WindowRatio:         int(in.WindowRatio),
		DenoiserSigmaColor:  float64(in.DenoiserSigmaColor),
		DenoiserSigmaSpace:  float64(in.DenoiserSigmaSpace),
		DenoiserDiameter:    int(in.DenoiserDiameter),
		DenoiserUmbralColor: float64(in.DenoiserUmbralColor),
		MinUmbral:           float64(in.MinUmbral),
		MaxUmbral:           float64(in.MaxUmbral),
		MinArea:             int(in.MinArea),
		MaxArea:             int(in.MaxArea),
		Ratio:               int(in.Ratio),
		BoxSizes:            boxSizes,
		LogLogWidth:         int(in.LogLogWidth),
		LogLogHeight:        int(in.LogLogHeight),
		Umbral:              umbral,
	}
}

//export FreeCFG
func FreeCFG(cfg C.struct_CFG) {
	C.free(unsafe.Pointer(cfg.KNN.Distance))
	C.free(unsafe.Pointer(cfg.KNN.Selector))
	C.free(unsafe.Pointer(cfg.BoxSizes))
}

//export Free
func Free(ptr *C.void) {
	C.free(unsafe.Pointer(ptr))
}

//export Malloc
func Malloc(size C.size_t) *C.void {
	return (*C.void)(C.malloc(size))
}

//export MemCopy
func MemCopy(dst, src *C.void, size C.size_t) *C.void {
	return (*C.void)(C.memcpy(unsafe.Pointer(dst), unsafe.Pointer(src), size))
}

//export SaveCFG
func SaveCFG(fileName *C.char, in C.struct_CFG) *C.char {
	gocfg := CfgToGo(in)
	err := cfg.SaveCFG(C.GoString(fileName), gocfg)
	if err != nil {
		return C.CString(err.Error())
	}
	return nil
}

//export LoadCFG
func LoadCFG(fileName *C.char, out *C.struct_CFG) *C.char {
	ldcfg, err := cfg.LoadCFG(C.GoString(fileName))
	if err != nil {
		return C.CString(err.Error())
	}
	*out = CFGToC(ldcfg)
	return nil
}

var pointsMp = map[int][]*lib.DataPoint{}

//export NewPoints
func NewPoints() C.int {
	listID := len(pointsMp)
	pointsMp[listID] = []*lib.DataPoint{}
	return C.int(listID)
}

//export LoadPoints
func LoadPoints(fileName *C.char, listID *C.int) *C.char {
	ldPoints, err := lib.LoadPoints(C.GoString(fileName))
	if err != nil {
		return C.CString(err.Error())
	}
	id := len(pointsMp)
	pointsMp[id] = ldPoints
	*listID = C.int(id)
	return nil
}

//export SavePoints
func SavePoints(fileName *C.char, listID C.int) *C.char {
	err := lib.SavePoints(C.GoString(fileName), pointsMp[int(listID)])
	if err != nil {
		return C.CString(err.Error())
	}
	return nil
}

//export FreePoints
func FreePoints(listID C.int) {
	delete(pointsMp, int(listID))
}

var imgSetMp = map[int][]*lib.ImageSetItem{}

//export NewImageSet
func NewImageSet() C.int {
	list := []*lib.ImageSetItem{}
	id := len(imgSetMp)
	imgSetMp[id] = list
	return C.int(id)
}

//export ImageSetAppend
func ImageSetAppend(listID C.int, img C.struct_RGBA, mriLabel *C.char) {
	list := imgSetMp[int(listID)]
	list = append(list, lib.NewImageSetItem(ToGoImage(img), C.GoString(mriLabel)))
	imgSetMp[int(listID)] = list
}

//export FreeImageSet
func FreeImageSet(listID C.int) {
	delete(imgSetMp, int(listID))
}

var fileSetMp = map[int][]*lib.FileSetItem{}

//export NewFileSet
func NewFileSet() C.int {
	list := make([]*lib.FileSetItem, 0, 10)
	listID := len(fileSetMp)
	fileSetMp[listID] = list
	return C.int(listID)
}

//export FileSetAppend
func FileSetAppend(listID C.int, fileName, mriLabel *C.char) {
	list := fileSetMp[int(listID)]
	list = append(list, lib.NewFileSetItem(C.GoString(fileName), C.GoString(mriLabel)))
	fileSetMp[int(listID)] = list
}

//export FreeFileSet
func FreeFileSet(listID C.int) {
	delete(fileSetMp, int(listID))
}

var knnMp = map[int]*lib.KNNFractal{}

//export NewKNNFractal
func NewKNNFractal(cfg C.struct_CFG, listID C.int) C.int {
	knnCFG := CfgToGo(cfg)
	points := pointsMp[int(listID)]
	knn := lib.NewKNNFractal(knnCFG, points)
	knnID := len(knnMp)
	knnMp[knnID] = knn
	return C.int(knnID)
}

//export FreeKNNFractal
func FreeKNNFractal(knnID C.int) {
	delete(knnMp, int(knnID))
}

//export TrainWithImages
func TrainWithImages(knnID, listID C.int) *C.char {
	knn := knnMp[int(knnID)]
	images := imgSetMp[int(listID)]
	err := knn.TrainWithImages(images)
	if err != nil {
		return C.CString(err.Error())
	}
	return nil
}

//export TrainWithFiles
func TrainWithFiles(knnID C.int, listID C.int) *C.char {
	knn := knnMp[int(knnID)]
	files := fileSetMp[int(listID)]
	err := knn.TrainWithFiles(files)
	if err != nil {
		return C.CString(err.Error())
	}
	return nil
}

//export Fit
func Fit(knnID C.int, img C.struct_RGBA, labelOut **C.char) *C.char {
	knn := knnMp[int(knnID)]
	label, err := knn.Fit(ToGoImage(img))
	if err != nil {
		return C.CString(err.Error())
	}
	*labelOut = C.CString(label)
	return nil
}

//export GetPoints
func GetPoints(knnID C.int) C.int {
	knn := knnMp[int(knnID)]
	points := knn.GetPoints()
	listID := len(pointsMp)
	pointsMp[listID] = points
	return C.int(listID)
}

var modelMp = map[int]*lib.FDModel{}

//export NewModel
func NewModel(cfg C.struct_CFG) C.int {
	model := lib.NewFDModel(CfgToGo(cfg))
	modelID := len(modelMp)
	modelMp[modelID] = model
	return C.int(modelID)
}

//export FreeModel
func FreeModel(modelID C.int) {
	delete(modelMp, int(modelID))
}

var evalMp = map[int]*lib.ModelEvaluation{}

//export Eval
func Eval(modelID C.int, img C.struct_RGBA, evalID *C.int) *C.char {
	model := modelMp[int(modelID)]
	eval, err := model.Eval(ToGoImage(img))
	if err != nil {
		return C.CString(err.Error())
	}
	id := len(evalMp)
	*evalID = C.int(id)
	evalMp[id] = eval
	return nil
}

//export FreeEval
func FreeEval(evalID C.int) {
	delete(evalMp, int(evalID))
}

//export EvalLen
func EvalLen(evalID C.int) C.int {
	return C.int(len(evalMp[int(evalID)].FDs))
}

func FractalDimToC(fd *pipes.FractalDim) C.struct_FractalDim {
	var fdc C.struct_FractalDim
	fdc.LogSizes = (*C.double)(C.malloc(C.size_t(8 * len(fd.LogMeasure))))
	fdc.LogMeasure = (*C.double)(C.malloc(C.size_t(8 * len(fd.LogMeasure))))
	C.memcpy(unsafe.Pointer(fdc.LogSizes), unsafe.Pointer(&fd.LogSizes[0]), C.size_t(8*len(fd.LogSizes)))
	C.memcpy(unsafe.Pointer(fdc.LogMeasure), unsafe.Pointer(&fd.LogMeasure[0]), C.size_t(8*len(fd.LogMeasure)))
	fdc.Len = C.int(len(fd.LogSizes))
	fdc.FD = C.double(fd.FD)
	return fdc
}

//export FreeFractalDim
func FreeFractalDim(fd C.struct_FractalDim) {
	C.free(unsafe.Pointer(fd.LogMeasure))
	C.free(unsafe.Pointer(fd.LogSizes))
}

//export EvalFDAt
func EvalFDAt(evalID C.int, index C.int, fdptr *C.struct_FractalDim) {
	eval := evalMp[int(evalID)]
	fd := eval.FDs[int(index)]
	*fdptr = FractalDimToC(fd)
}

//export EvalUmbralAt
func EvalUmbralAt(evalID int, index C.int, umbralptr *C.struct_Umbral) {
	umbral := evalMp[int(evalID)].Umbrals[int(index)]
	var umc C.struct_Umbral
	umc.Min = C.double(umbral.Min)
	umc.Max = C.double(umbral.Max)
	*umbralptr = umc
}

//export EvalLogLogAt
func EvalLogLogAt(evalID int, index C.int, rgba *C.struct_RGBA) {
	img := evalMp[int(evalID)].LogLogs[int(index)]
	*rgba = ImageToC(img)
}

//export EvalMFS
func EvalMFS(evalID int, mfs *C.struct_MFS) {
	goMfs := evalMp[int(evalID)].MFS
	dataPtr := (*C.double)(C.malloc(C.size_t(int(unsafe.Sizeof(float64(0))) * goMfs.Width() * goMfs.Height())))

	for i := 0; i < goMfs.Height(); i++ {
		for j := 0; j < goMfs.Width(); j++ {
			//dataPtr + (i*width + j)*sizeof(float64)
			*(*C.double)(unsafe.Pointer(
				uintptr(unsafe.Pointer(dataPtr)) +
					(uintptr(
						i*goMfs.Width())+uintptr(j))*
						uintptr(unsafe.Sizeof(float64(0))))) =
				C.double(goMfs[i][j])
		}
	}

	mfs.Width = C.int(goMfs.Width())
	mfs.Height = C.int(goMfs.Height())
	mfs.Data = dataPtr
}

//export FreeMFS
func FreeMFS(mfs C.struct_MFS) {
	C.free(unsafe.Pointer(mfs.Data))
}

//export EvalGetPoints
func EvalGetPoints(evalID C.int, n *C.int) *C.double {
	pts := []float64(evalMp[int(evalID)].GetPoints())
	*n = C.int(len(pts))
	points := (*C.double)(C.malloc(C.size_t(8 * len(pts))))
	C.memcpy(unsafe.Pointer(points), unsafe.Pointer(&pts[0]), C.size_t(8*len(pts)))
	return points
}

//export FreeEvalPoints
func FreeEvalPoints(ptr *C.double) {
	C.free(unsafe.Pointer(ptr))
}

var sampleMp = map[int]*lib.Sample{}

//export NewSample
func NewSample(ctumors **C.char, lenTumors C.int, cnotumors **C.char, lenNoTumors C.int, cpart C.double) C.int {
	tumorCharSlice := unsafe.Slice(ctumors, lenTumors)
	noTumorCharSlice := unsafe.Slice(cnotumors, lenNoTumors)
	tumors := make([]string, len(tumorCharSlice))
	notumors := make([]string, len(noTumorCharSlice))
	for i := 0; i < len(tumors); i++ {
		tumors[i] = C.GoString(tumorCharSlice[i])
	}
	for i := 0; i < len(notumors); i++ {
		notumors[i] = C.GoString(noTumorCharSlice[i])
	}
	partition := float64(cpart)
	sample := lib.NewSample(partition, tumors, notumors)
	sampleID := len(sampleMp)
	sampleMp[sampleID] = sample
	return C.int(sampleID)
}

//export FreeSample
func FreeSample(sampleID C.int) {
	delete(sampleMp, int(sampleID))
}

//export LoadSample
func LoadSample(fileName *C.char, sampleID *C.int) *C.char {
	sample, err := lib.LoadSample(C.GoString(fileName))
	if err != nil {
		return C.CString(err.Error())
	}
	ID := len(sampleMp)
	*sampleID = C.int(ID)
	sampleMp[ID] = sample
	return nil
}

//export SaveSample
func SaveSample(fileName *C.char, sampleID C.int) *C.char {
	sample := sampleMp[int(sampleID)]
	err := lib.SaveSample(C.GoString(fileName), sample)
	return C.CString(err.Error())
}

//export Optimize
func Optimize(sampleID, n C.int, cfg C.struct_CFG, new_sample_id *C.int) *C.char {
	sample := sampleMp[int(sampleID)]
	new_sample, err := sample.Optimize(int(n), CfgToGo(cfg))
	if err != nil {
		return C.CString(err.Error())
	}
	ID := len(sampleMp)
	sampleMp[ID] = new_sample
	*new_sample_id = C.int(ID)
	return nil
}

//export Points
func Points(sampleID C.int) C.int {
	sample := sampleMp[int(sampleID)]
	ID := len(pointsMp)
	pointsMp[ID] = sample.GetPoints()
	return C.int(ID)
}
