package controllers

import (
	"go-simple-app/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VersionController struct {
	versionService *services.VersionService
}

func NewVersionController(versionService *services.VersionService) *VersionController {
	return &VersionController{
		versionService: versionService,
	}
}

// GetVersion 獲取版本資訊
func (vc *VersionController) GetVersion(c *gin.Context) {
	version := vc.versionService.GetVersionInfo()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    version,
	})
}

// GetShortVersion 獲取簡短版本號
func (vc *VersionController) GetShortVersion(c *gin.Context) {
	shortVersion := vc.versionService.GetShortVersion()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"version": shortVersion,
		},
	})
}
