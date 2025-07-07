package handler

import (
	"github.com/gin-gonic/gin"
	v1 "go-dianping/api/v1"
	"go.uber.org/zap"
	"hash/fnv"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var imageUploadPath = filepath.Join("storage", "images")

type UploadHandler struct {
	*Handler
}

func NewUploadHandler(
	handler *Handler,
) *UploadHandler {
	return &UploadHandler{
		Handler: handler,
	}
}

func (h *UploadHandler) UploadImage(ctx *gin.Context) {
	image, err := ctx.FormFile("image")
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// 获取原始文件名称
	originalFilename := image.Filename
	// 生成新文件名
	filename, err := h.createNewFilename(originalFilename)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	// 保存文件
	if err := ctx.SaveUploadedFile(image, filepath.Join(imageUploadPath, filename)); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	// 返回结果
	h.logger.Debug("文件上传成功", zap.String("filename", filename))
	v1.HandleSuccess(ctx, filename)
}

func (h *UploadHandler) createNewFilename(originalFilename string) (string, error) {
	// 获取后缀
	var suffix string
	if dot := strings.LastIndex(originalFilename, "."); dot != -1 {
		suffix = originalFilename[dot+1:]
	}

	hashBuilder := fnv.New32a()
	if _, err := hashBuilder.Write([]byte(originalFilename)); err != nil {
		return "", err
	}
	hash := hashBuilder.Sum32()
	d1, d2 := strconv.Itoa(int(hash&0xF)), strconv.Itoa(int(hash>>4&0xF))
	if err := os.MkdirAll(filepath.Join(imageUploadPath, d1, d2), os.ModePerm); err != nil {
		return "", err
	}
	return filepath.Join(imageUploadPath, d1, d2, strconv.Itoa(int(hash)), suffix), nil
}
