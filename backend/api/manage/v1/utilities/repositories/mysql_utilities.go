package repositories

import (
	"backend/domain"
	"backend/models"

	"backend/pkg/v1/mysql"
	"backend/utils/log"
	"database/sql"

	"github.com/gin-gonic/gin"
)

type UtilitiesRepositories struct {
	sql *sql.DB
}

func NewTestRepoUtilities(sql *sql.DB) domain.UtilitiesRepositories {
	return &UtilitiesRepositories{
		sql: sql,
	}
}

func (t *UtilitiesRepositories) GetUtilities(c *gin.Context) ([]models.RespGetList, error) {
	// Get connetion DB
	db, err := mysql.GetConnectionItem()
	if err != nil {
		log.Log.Errorln("Error GetConnectionUtilities", err.Error())
		return []models.RespGetList{}, err
	}

	rows, err := db.Query("SELECT a.id, a.name, a.price, a.stock_quantity, b.category, a.description, a.image FROM tbl_utilities AS a LEFT JOIN tbl_category AS b ON b.id = a.id_category;")
	if err != nil {
		if err == sql.ErrNoRows {
			return []models.RespGetList{}, err
		}
		return []models.RespGetList{}, err
	}
	defer rows.Close()

	// Mapping to struct rows
	var getUtilities []models.RespGetList
	for rows.Next() {
		getUtility := models.RespGetList{}
		if errScan := rows.Scan(&getUtility.Id, &getUtility.Name, &getUtility.Price, &getUtility.Stock, &getUtility.Category, &getUtility.Description, &getUtility.Image); errScan != nil {
			return []models.RespGetList{}, err
		}
		getUtilities = append(getUtilities, getUtility)
	}
	if errRows := rows.Err(); errRows != nil {
		return []models.RespGetList{}, err
	}

	return getUtilities, nil
}

func (t *UtilitiesRepositories) GetUtilitiesById(c *gin.Context, uid int) ([]models.RespGetList, error) {
	// Get connection DB
	db, err := mysql.GetConnectionItem()
	if err != nil {
		log.Log.Errorln("Error GetConnectionUtilities", err.Error())
		return []models.RespGetList{}, err
	}

	rows, err := db.Query("SELECT a.id, a.name, a.price, a.stock_quantity, b.category, a.description, a.image FROM tbl_utilities AS a LEFT JOIN tbl_category AS b ON b.id = a.id_category WHERE a.user_id = ?;", uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return []models.RespGetList{}, err
		}
		return []models.RespGetList{}, err
	}
	defer rows.Close()

	// Mapping to struct rows
	var getUtilities []models.RespGetList
	for rows.Next() {
		getUtility := models.RespGetList{}
		if errScan := rows.Scan(&getUtility.Id, &getUtility.Name, &getUtility.Price, &getUtility.Stock, &getUtility.Category, &getUtility.Description, &getUtility.Image); errScan != nil {
			return []models.RespGetList{}, err
		}
		getUtilities = append(getUtilities, getUtility)
	}
	if errRows := rows.Err(); errRows != nil {
		return []models.RespGetList{}, err
	}

	return getUtilities, nil
}

func (t *UtilitiesRepositories) CreateUtilities(c *gin.Context, req models.ReqCreateUtilities, image []byte, uid int) error {
	// Get connection DB
	db, err := mysql.GetConnectionItem()
	if err != nil {
		log.Log.Errorln("Error GetConnectionItem", err.Error())
		return err
	}

	if _, err = db.Exec("INSERT INTO tbl_utilities (name, price, stock_quantity, id_category, description, image, user_id) VALUES (?, ?, ?, ?, ?, ?, ?)", req.Name, req.Price, req.Stock, req.IdCategory, req.Description, image, uid); err != nil {
		log.Log.Errorln("Error running query CreateUtilities : ", err.Error())
		return err
	}

	return nil
}

func (t *UtilitiesRepositories) UpdateUtilities(c *gin.Context, req models.ReqUpdateUtilities, image []byte, uid int) error {
	// Get connection DB
	db, err := mysql.GetConnectionItem()
	if err != nil {
		log.Log.Errorln("Error GetConnectionUtilities", err.Error())
		return err
	}

	if _, err = db.Exec("UPDATE tbl_utilities SET name = ?, price = ?, stock_quantity = ?, id_category = ?, description = ?, image = ?, user_id = ? WHERE id = ?", req.Name, req.Price, req.Stock, req.IdCategory, req.Description, image, uid, req.Id); err != nil {
		log.Log.Errorln("Error running query UpdateUtilities : ", err.Error())
		return err
	}

	return nil
}

func (t *UtilitiesRepositories) DeleteUtilities(c *gin.Context, req models.ReqDeleteUtilities) error {
	// Get connection DB
	db, err := mysql.GetConnectionItem()
	if err != nil {
		log.Log.Errorln("Error GetConnectionItem", err.Error())
		return err
	}

	if _, err = db.Exec("DELETE FROM tbl_utilities WHERE id = ?", req.Id); err != nil {
		log.Log.Errorln("Error running query DeleteUtilities : ", err.Error())
		return err
	}

	return nil
}
