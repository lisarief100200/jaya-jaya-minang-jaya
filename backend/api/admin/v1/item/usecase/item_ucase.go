package usecase

import (
	"backend/domain"
	"backend/models"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type ItemUsecase struct {
	IRepositories domain.ItemRepositories
}

func NewItemUsecase(itemRepo domain.ItemRepositories) domain.ItemUsecase {
	return &ItemUsecase{
		IRepositories: itemRepo,
	}
}

func (r *ItemUsecase) GetItems(c *gin.Context) ([]models.RespGetList, error) {
	list, err := r.IRepositories.GetItems(c)
	if err != nil {
		return []models.RespGetList{}, err
	}

	return list, nil
}

func (r *ItemUsecase) CreateItem(c *gin.Context, req *models.ReqCreateItem, uid int) error {
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
	err = r.IRepositories.CreateItem(c, *req, imageBytes, uid)
	if err != nil {
		return err
	}

	return nil
}

func (r *ItemUsecase) UpdateItem(c *gin.Context, req *models.ReqUpdateItem, uid int) error {
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
	err = r.IRepositories.UpdateItem(c, *req, imageBytes, uid)
	if err != nil {
		return err
	}

	return nil
}

func (r *ItemUsecase) DeleteItem(c *gin.Context, req models.ReqDeleteItem) error {
	err := r.IRepositories.DeleteItem(c, req)
	if err != nil {
		return err
	}
	return nil
}
