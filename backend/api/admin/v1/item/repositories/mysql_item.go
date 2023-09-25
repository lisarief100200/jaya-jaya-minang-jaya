package repositories

import (
	"backend/domain"
	"backend/models"
	"backend/pkg/v1/mysql"
	"backend/utils/log"
	"database/sql"

	"github.com/gin-gonic/gin"
)

type ItemRepositories struct {
	sql *sql.DB
}

func NewTestRepoItem(sql *sql.DB) domain.ItemRepositories {
	return &ItemRepositories{
		sql: sql,
	}
}

func (t *ItemRepositories) GetItems(c *gin.Context) ([]models.RespGetList, error) {
	// Get connection DB
	db, err := mysql.GetConnectionItem()
	if err != nil {
		log.Log.Errorln("Error GetConnectionItem", err.Error())
		return []models.RespGetList{}, err
	}

	rows, err := db.Query("SELECT a.id, a.name, a.price, a.stock_quantity, b.category, a.description, a.image FROM tbl_item AS a LEFT JOIN tbl_category AS b ON b.id = a.id_category;")
	if err != nil {
		if err == sql.ErrNoRows {
			return []models.RespGetList{}, err
		}
		return []models.RespGetList{}, err
	}
	defer rows.Close()

	// Mapping to struct rows
	var getItems []models.RespGetList
	for rows.Next() {
		getItem := models.RespGetList{}
		if errScan := rows.Scan(&getItem.Id, &getItem.Name, &getItem.Price, &getItem.Stock, &getItem.Category, &getItem.Description, &getItem.Image); errScan != nil {
			return []models.RespGetList{}, err
		}
		getItems = append(getItems, getItem)
	}
	if errRows := rows.Err(); errRows != nil {
		return []models.RespGetList{}, err
	}

	return getItems, nil
}

func (t *ItemRepositories) CreateItem(c *gin.Context, req models.ReqCreateItem, image []byte, uid int) error {
	// Get connection DB
	db, err := mysql.GetConnectionItem()
	if err != nil {
		log.Log.Errorln("Error GetConnectionItem", err.Error())
		return err
	}

	if _, err = db.Exec("INSERT INTO tbl_item (name, price, stock_quantity, id_category, description, image, user_id) VALUES (?, ?, ?, ?, ?, ?, ?)", req.Name, req.Price, req.Stock, req.IdCategory, req.Description, image, uid); err != nil {
		log.Log.Errorln("Error running query CreateItem : ", err.Error())
		return err
	}

	return nil
}

func (t *ItemRepositories) GetFilePath(c *gin.Context, req *models.ReqUpdateItem) (string, error) {
	db, err := mysql.GetConnectionItem()
	if err != nil {
		log.Log.Errorln("Error GetConnectionItem", err.Error())
		return "", err
	}

	var path string
	if err := db.QueryRow("SELECT image FROM tbl_item where id = ?", req.Id).Scan(&path); err != nil {
		log.Log.Errorln("Error scanning query GetFilePath")
		return "", err
	}

	return path, nil
}

func (t *ItemRepositories) UpdateItem(c *gin.Context, req models.ReqUpdateItem, image []byte, uid int) error {
	// Get connection DB
	db, err := mysql.GetConnectionItem()
	if err != nil {
		log.Log.Errorln("Error GetConnectionItem", err.Error())
		return err
	}

	if _, err = db.Exec("UPDATE tbl_item SET name = ?, price = ?, stock_quantity = ?, id_category = ?, description = ?, image = ?, user_id = ? WHERE id = ?", req.Name, req.Price, req.Stock, req.IdCategory, req.Description, image, uid, req.Id); err != nil {
		log.Log.Errorln("Error running query UpdateItem : ", err.Error())
		return err
	}

	return nil
}

func (t *ItemRepositories) DeleteItem(c *gin.Context, req models.ReqDeleteItem) error {
	// Get connection DB
	db, err := mysql.GetConnectionItem()
	if err != nil {
		log.Log.Errorln("Error GetConnectionItem", err.Error())
		return err
	}

	if _, err = db.Exec("DELETE FROM tbl_item WHERE id = ?", req.Id); err != nil {
		log.Log.Errorln("Error running query DeleteItem : ", err.Error())
		return err
	}

	return nil
}
