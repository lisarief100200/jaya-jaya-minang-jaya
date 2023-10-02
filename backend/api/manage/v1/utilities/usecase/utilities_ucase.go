package usecase

import (
	"backend/domain"
	"backend/models"
	"io/ioutil"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UtilitiesUsecase struct {
	URepositories domain.UtilitiesRepositories
}

func NewUtilitiesUsecase(utilitiesRepo domain.UtilitiesRepositories) domain.UtilitiesUsecase {
	return &UtilitiesUsecase{
		URepositories: utilitiesRepo,
	}
}

func (r *UtilitiesUsecase) GetUtilities(c *gin.Context, uid, level string) ([]models.RespGetList, error) {
	var list []models.RespGetList
	var err error
	uidInt, _ := strconv.Atoi(uid)

	if level == "admin" {
		list, err = r.URepositories.GetUtilities(c)
		if err != nil {
			return []models.RespGetList{}, err
		}
	} else {
		list, err = r.URepositories.GetUtilitiesById(c, uidInt)
		if err != nil {
			return []models.RespGetList{}, err
		}
	}

	return list, nil
}

func (r *UtilitiesUsecase) CreateUtilities(c *gin.Context, req *models.ReqCreateUtilities, uid int) error {
	// Check file size
	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB (ukuran maksimum file yang diizinkan)
	if err != nil {
		return err
	}

	// Dapatkan gambar dalam bentuk byte
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		return err
	}
	defer file.Close()

	imageBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	// Save data to database
	err = r.URepositories.CreateUtilities(c, *req, imageBytes, uid)
	if err != nil {
		return err
	}

	return nil
}

func (r *UtilitiesUsecase) UpdateUtilities(c *gin.Context, req *models.ReqUpdateUtilities, uid int) error {
	// Check file size
	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		return err
	}

	file, _, err := c.Request.FormFile("image")
	if err != nil {
		return err
	}
	defer file.Close()

	imageBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	// Save data to database
	err = r.URepositories.UpdateUtilities(c, *req, imageBytes, uid)
	if err != nil {
		return err
	}

	return nil
}

func (r *UtilitiesUsecase) DeleteUtilities(c *gin.Context, req models.ReqDeleteUtilities) error {
	err := r.URepositories.DeleteUtilities(c, req)
	if err != nil {
		return err
	}
	return nil
}
