package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func UploadProductImage(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get image from request",
		})
	}

	// validasi gambar

	// validasi ekstensi
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowedExts[ext] {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid file type. Only JPG, JPEG, PNG, and WEBP are allowed",
		})
	}

	//validasi content type
	contentType := file.Header.Get("Content-Type")
	allowedContentTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}
	if !allowedContentTypes[contentType] {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid content type. Only image/jpeg, image/png, and image/webp are allowed",
		})
	}

	timestamp := time.Now().UnixNano()
	newFileName := fmt.Sprintf("product_%d%s", timestamp, filepath.Ext(file.Filename))
	err = c.SaveFile(file, "./storage/product/"+newFileName)

	if err != nil {
		fmt.Println("Error saving file:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Internal server error",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Image uploaded successfully",
		"file_name": newFileName,
	})
}