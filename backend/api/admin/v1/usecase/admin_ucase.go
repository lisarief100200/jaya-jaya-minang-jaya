package usecase

import (
	"backend/domain"
	"backend/models"
	"backend/utils/helpers"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminUsecase struct {
	ARepositories domain.AdminRepositories
}

func NewAdminUsecase(adminRepo domain.AdminRepositories) domain.AdminUsecase {
	return &AdminUsecase{
		ARepositories: adminRepo,
	}
}

func (r *AdminUsecase) CheckIdAdmin(c *gin.Context, req models.ReqLoginAdmin) (string, error) {
	uid, err := r.ARepositories.CheckIdAdmin(c, strings.TrimSpace(req.Username))
	if err != nil {
		return "", err
	}
	return uid, nil
}

func (r *AdminUsecase) CheckLoginAdmin(c *gin.Context, req models.ReqLoginAdmin, uid string) error {
	validAdmin, err := r.ARepositories.CheckLoginAdmin(c, req, uid)
	if err != nil {
		return err
	}

	if !validAdmin {
		return err
	}
	return nil
}

func (r *AdminUsecase) GetPassword(c *gin.Context, uid string, req models.ReqLoginAdmin) error {
	pass, err := r.ARepositories.GetPassword(c, uid)
	if err != nil {
		return err
	}

	ps := helpers.Hashed(req.Password)

	if ps != pass {
		return err
	}
	return nil
}

func (r *AdminUsecase) GenerateSessionID(c *gin.Context, typeLogin string) (string, error) {
	sessionId, err := r.ARepositories.GenerateSessionID(c, typeLogin)
	if err != nil {
		return "", err
	}
	return sessionId, nil
}

func (r *AdminUsecase) GetAdminProp(c *gin.Context, uid string) (models.AdminProp, error) {
	prop, err := r.ARepositories.GetAdminProp(c, uid)
	if err != nil {
		return models.AdminProp{}, err
	}
	return prop, nil
}

func (r *AdminUsecase) UpdateLastLogin(c *gin.Context, uid string) error {
	err := r.ARepositories.UpdateLastLogin(c, uid)
	if err != nil {
		return err
	}
	return nil
}

func (r *AdminUsecase) InsertSession(c *gin.Context, sessionId string, uid string, typeLogin string, token string, exp time.Time) error {
	err := r.ARepositories.InsertSession(c, sessionId, uid, typeLogin, token, exp)
	if err != nil {
		return err
	}
	return nil
}

func (r *AdminUsecase) GetItems(c *gin.Context) ([]models.RespGetList, error) {
	list, err := r.ARepositories.GetItems(c)
	if err != nil {
		return []models.RespGetList{}, err
	}

	return list, nil
}

func (r *AdminUsecase) CreateItem(c *gin.Context, req *models.ReqCreateItem, uid int) error {
	// Periksa apakah ada form yang diunggah
	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB (ukuran maksimum file yang diizinkan)
	if err != nil {
		return err
	}

	// Dapatkan file dari form
	file, handler, err := c.Request.FormFile("image")
	if err != nil {
		return err
	}
	defer file.Close()

	// Simpan file ke folder yang ditentukan
	imageFolderPath := "uploads"
	uniqueFilename := helpers.GenerateUniqueFilename(filepath.Ext(handler.Filename))
	imagePath := filepath.Join(imageFolderPath, uniqueFilename)
	imagePath = filepath.ToSlash(imagePath)
	outFile, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Salin isi gambar yang diunggah ke file baru
	_, err = io.Copy(outFile, file)
	if err != nil {
		return err
	}

	// Simpan informasi gambar ke dalam basis data
	err = r.ARepositories.CreateItem(c, *req, imagePath, uid)
	if err != nil {
		return err
	}

	return nil
}

func (r *AdminUsecase) UpdateItem(c *gin.Context, req *models.ReqUpdateItem, uid int) error {
	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		return err
	}

	file, handler, err := c.Request.FormFile("image")
	if err != nil {
		return err
	}
	defer file.Close()

	// Get File Path
	path, err := r.ARepositories.GetFilePath(c, req)
	if err != nil {
		return err
	}

	// Hapus file
	if path != "" {
		err := os.Remove(path)
		if err != nil {
			log.Println("Failed to delete old image", err)
			return err
		}
	}

	// Simpan ke folder yang ditentukan
	imageFolderPath := "uploads"
	uniqueFilename := helpers.GenerateUniqueFilename(filepath.Ext(handler.Filename))
	imagePath := filepath.Join(imageFolderPath, uniqueFilename)
	imagePath = filepath.ToSlash(imagePath)
	outFile, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Salin isi gambar yang diunggah ke file baru
	_, err = io.Copy(outFile, file)
	if err != nil {
		return err
	}

	// Simpan informasi gambar ke dalam basis data
	err = r.ARepositories.UpdateItem(c, *req, imagePath, uid)
	if err != nil {
		return err
	}

	return nil
}

func (r *AdminUsecase) DeleteItem(c *gin.Context, req models.ReqDeleteItem) error {
	err := r.ARepositories.DeleteItem(c, req)
	if err != nil {
		return err
	}
	return nil
}

func (r *AdminUsecase) LogoutAdmin(c *gin.Context, uid int) error {
	err := r.ARepositories.LogoutAdmin(c, uid)
	if err != nil {
		return err
	}
	return nil
}
