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

func (h *UploadHandler) DeleteBlogImg(ctx *gin.Context) {
	var req struct {
		Name string `form:"name"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	info, err := os.Stat(filepath.Join(imageUploadPath, req.Name))
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	if info.IsDir() {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrIncorrectFilename.Error(), nil)
		return
	}
	if err := os.Remove(filepath.Join(imageUploadPath, req.Name)); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

func (h *UploadHandler) createNewFilename(originalFilename string) (string, error) {
	// 获取后缀
	var suffix string
	if dot := strings.LastIndex(originalFilename, "."); dot != -1 {
		suffix = originalFilename[dot+1:]
	}
	// 生成目录
	hashBuilder := fnv.New32a()
	if _, err := hashBuilder.Write([]byte(originalFilename)); err != nil {
		return "", err
	}
	hash := hashBuilder.Sum32()
	d1, d2 := strconv.Itoa(int(hash&0xF)), strconv.Itoa(int(hash>>4&0xF))
	if err := os.MkdirAll(filepath.Join(imageUploadPath, d1, d2), os.ModePerm); err != nil {
		return "", err
	}
	// 生成文件名
	return filepath.Join(imageUploadPath, d1, d2, strconv.Itoa(int(hash)), suffix), nil
}
