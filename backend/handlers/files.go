package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/ankuragarwal/fernfs/backend/storage"
)

type FileHandler struct {
	adapter storage.Adapter
}

func NewFileHandler(adapter storage.Adapter) *FileHandler {
	return &FileHandler{
		adapter: adapter,
	}
}

// RegisterRoutes registers file operation routes
func (h *FileHandler) RegisterRoutes(router *gin.Engine) {
	files := router.Group("/files")
	{
		files.POST("/mkdir", h.Mkdir)
		files.POST("/rmdir", h.Rmdir)
		files.GET("/readdir/*path", h.Readdir)
		files.GET("/stat/*path", h.Stat)
		files.GET("/read/*path", h.ReadFile)
		files.POST("/write/*path", h.WriteFile)
		files.POST("/unlink", h.Unlink)
		files.POST("/rename", h.Rename)
		files.POST("/chmod", h.Chmod)
	}
}

// Mkdir handles directory creation
func (h *FileHandler) Mkdir(c *gin.Context) {
	var req struct {
		Path string `json:"path" binding:"required"`
		Mode uint32 `json:"mode"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Mode == 0 {
		req.Mode = 0755
	}

	if err := h.adapter.Mkdir(c.Request.Context(), req.Path, os.FileMode(req.Mode)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// Rmdir handles directory removal
func (h *FileHandler) Rmdir(c *gin.Context) {
	var req struct {
		Path string `json:"path" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.adapter.Rmdir(c.Request.Context(), req.Path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// Readdir handles directory listing
func (h *FileHandler) Readdir(c *gin.Context) {
	path := filepath.Clean(c.Param("path"))

	entries, err := h.adapter.Readdir(c.Request.Context(), path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entries)
}

// Stat handles file/directory info retrieval
func (h *FileHandler) Stat(c *gin.Context) {
	path := filepath.Clean(c.Param("path"))

	info, err := h.adapter.Stat(c.Request.Context(), path)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}

// ReadFile handles file reading
func (h *FileHandler) ReadFile(c *gin.Context) {
	path := filepath.Clean(c.Param("path"))

	// Parse range header if present
	var opts storage.ReadOptions
	if rangeHeader := c.GetHeader("Range"); rangeHeader != "" {
		var start, end int64
		if _, err := fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end); err == nil {
			opts.Offset = start
			if end > start {
				opts.Length = end - start + 1
			}
		}
	}

	reader, err := h.adapter.ReadFile(c.Request.Context(), path, opts)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer reader.Close()

	info, err := h.adapter.Stat(c.Request.Context(), path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.DataFromReader(http.StatusOK, info.Size, "application/octet-stream", reader, nil)
}

// WriteFile handles file writing
func (h *FileHandler) WriteFile(c *gin.Context) {
	path := filepath.Clean(c.Param("path"))

	var opts storage.WriteOptions
	if mode := c.GetHeader("X-File-Mode"); mode != "" {
		var modeInt uint32
		if _, err := fmt.Sscanf(mode, "%o", &modeInt); err == nil {
			opts.Mode = os.FileMode(modeInt)
		}
	}
	if opts.Mode == 0 {
		opts.Mode = 0644
	}

	if err := h.adapter.WriteFile(c.Request.Context(), path, c.Request.Body, opts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// Unlink handles file deletion
func (h *FileHandler) Unlink(c *gin.Context) {
	var req struct {
		Path string `json:"path" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.adapter.Unlink(c.Request.Context(), req.Path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// Rename handles file/directory renaming
func (h *FileHandler) Rename(c *gin.Context) {
	var req struct {
		OldPath string `json:"old_path" binding:"required"`
		NewPath string `json:"new_path" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.adapter.Rename(c.Request.Context(), req.OldPath, req.NewPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// Chmod handles file permission changes
func (h *FileHandler) Chmod(c *gin.Context) {
	var req struct {
		Path string `json:"path" binding:"required"`
		Mode uint32 `json:"mode" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.adapter.Chmod(c.Request.Context(), req.Path, os.FileMode(req.Mode)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}