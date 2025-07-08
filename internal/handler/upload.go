package handler

import (
	"fmt"
	"github.com/duke-git/lancet/v2/random"
	"github.com/gin-gonic/gin"
	v1 "go-dianping/api/v1"
	"go.uber.org/zap"
	"hash/fnv"
	"net/http"
	"os"
	"path/filepath"
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
	image, err := ctx.FormFile("file")
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
	var suffix string
	if dot := strings.LastIndex(originalFilename, "."); dot != -1 {
		suffix = originalFilename[dot+1:]
	}

	name, err := random.UUIdV4()
	if err != nil {
		return "", err
	}
	v := fnv.New32a()
	_, err = v.Write([]byte(name))
	if err != nil {
		return "", err
	}
	hash := v.Sum32()
	d1 := int(hash & 0xF)
	d2 := int((hash >> 4) & 0xF)
	dir := filepath.Join(imageUploadPath, fmt.Sprintf("%06d", d1), fmt.Sprintf("%06d", d2))
	if err := os.MkdirAll(dir, 0777); err != nil {
		return "", err
	}
	return dir + suffix, nil
}
