package repositories

import (
	"backend/domain"
	"backend/models"
	"backend/utils/log"
	"database/sql"
	"errors"
	"strings"
	"time"

	"backend/pkg/v1/mysql"

	"github.com/gin-gonic/gin"
)

type AdminRepositories struct {
	sql *sql.DB
}

func NewTestRepoAdmin(sql *sql.DB) domain.AdminRepositories {
	return &AdminRepositories{
		sql: sql,
	}
}

func (t *AdminRepositories) CheckIdAdmin(c *gin.Context, username string) (string, error) {
	// Get connectionDB
	db, err := mysql.GetConnectionAdmin()
	if err != nil {
		log.Log.Errorln("Error GetConnectionAdmin")
		return "", err
	}

	var id string

	err = db.QueryRow("SELECT id FROM tbl_admin WHERE username = ?", username).Scan(&id)
	if err != nil {
		log.Log.Errorln("Error scanning query CheckIdAdmin")
		return "", err
	}

	return id, nil
}

func (t *AdminRepositories) CheckLoginAdmin(c *gin.Context, req models.ReqLoginAdmin, uid string) (exist bool, err error) {
	// Get connection DB
	db, err := mysql.GetConnectionAdmin()
	if err != nil {
		log.Log.Errorln("Error GetConnectionAdmin")
		return false, err
	}

	err = db.QueryRow("SELECT username, password FROM tbl_admin WHERE id = ?", uid).Scan(&req.Username, &req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			log.Log.Errorln("Error scanning query CheckLoginAdmin")
			return false, errors.New("invalid")
		}
		return false, err
	}

	return true, nil
}

func (t *AdminRepositories) GetPassword(c *gin.Context, uid string) (string, error) {
	// Get connection DB
	db, err := mysql.GetConnectionAdmin()
	if err != nil {
		log.Log.Errorln("Error GetConnectionAdmin")
		return "", err
	}

	var pass string
	if err := db.QueryRow("SELECT password FROM tbl_admin where id = ?", uid).Scan(&pass); err != nil {
		log.Log.Errorln("Error scanning query GetPassword")
		return "", err
	}

	return pass, nil
}

func (t *AdminRepositories) GenerateSessionID(c *gin.Context, processName string) (string, error) {
	// Get connection DB
	db, err := mysql.GetConnectionAdmin()
	if err != nil {
		log.Log.Errorln("Error GetConnectionAdmin")
	}

	var sessionId string
	if err := db.QueryRow("SELECT generateSession(?)", processName).Scan(&sessionId); err != nil {
		log.Log.Errorln("Error GenerateSessionID : ", err.Error())
		return "", err
	}

	return sessionId, nil
}

func (t *AdminRepositories) GetAdminProp(c *gin.Context, uid string) (models.AdminProp, error) {
	var adminProp models.AdminProp
	db, err := mysql.GetConnectionAdmin()
	if err != nil {
		log.Log.Errorln("Error GetConnectionAdmin", err.Error())
		return adminProp, err
	}

	if err := db.QueryRow("SELECT name, level FROM tbl_admin WHERE id = ?", uid).Scan(&adminProp.Name, &adminProp.Level); err != nil {
		log.Log.Errorln("Error running query GetAdminProp", err.Error())
		return adminProp, err
	}

	return adminProp, nil
}

func (t *AdminRepositories) UpdateLastLogin(c *gin.Context, uid string) error {
	// Get connection DB
	db, err := mysql.GetConnectionAdmin()
	if err != nil {
		log.Log.Errorln("Error GetConnectionAdmin", err.Error())
		return err
	}

	if _, err = db.Exec("UPDATE tbl_admin SET latest_login = now() WHERE id = ?", uid); err != nil {
		log.Log.Errorln("Error running query UpdateLastLogin : ", err.Error())
		return err
	}

	return nil
}

func (t *AdminRepositories) InsertSession(c *gin.Context, session, uid, tipe, token string, exptime time.Time) error {
	// Get connection DB
	db, err := mysql.GetConnectionAdmin()
	if err != nil {
		log.Log.Errorln("Error GetConnectionAdmin", err.Error())
		return err
	}

	var (
		sessionId string
		crtTime   string
	)

	// Convert time
	dt := exptime.Format("2006-01-02 15:04:05")
	if err := db.QueryRow("CALL sp_i_session(?, ?, ?, ?, ?)", session, uid, tipe, token, dt).Scan(&sessionId, &crtTime); err != nil {
		// SP masih error, padahal udah jalan. Jadi gini aja dulu
		if err == sql.ErrNoRows {
			return nil
		} else {
			log.Log.Errorln("Error running store procedure InsertSession : ", err.Error())
			return err
		}
	}

	return nil
}

func (t *AdminRepositories) GetItems(c *gin.Context) ([]models.RespGetList, error) {
	// Get connection DB
	db, err := mysql.GetConnectionItem()
	if err != nil {
		log.Log.Errorln("Error GetConnectionItem", err.Error())
		return []models.RespGetList{}, err
	}

	rows, err := db.Query("SELECT a.id, a.name, a.price, a.stock_quantity, b.category, a.description, a.image_path FROM tbl_item AS a LEFT JOIN tbl_category AS b ON b.id = a.id_category;")
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

func (t *AdminRepositories) CreateItem(c *gin.Context, req models.ReqCreateItem, imagePath string, uid int) error {
	// Get connection DB
	db, err := mysql.GetConnectionItem()
	if err != nil {
		log.Log.Errorln("Error GetConnectionItem", err.Error())
		return err
	}

	if _, err = db.Exec("INSERT INTO tbl_item (name, price, stock_quantity, id_category, description, image_path, user_id) VALUES (?, ?, ?, ?, ?, ?, ?)", req.Name, req.Price, req.Stock, req.IdCategory, req.Description, imagePath, uid); err != nil {
		log.Log.Errorln("Error running query CreateItem : ", err.Error())
		return err
	}

	return nil
}

func (t *AdminRepositories) GetFilePath(c *gin.Context, req *models.ReqUpdateItem) (string, error) {
	db, err := mysql.GetConnectionItem()
	if err != nil {
		log.Log.Errorln("Error GetConnectionItem", err.Error())
		return "", err
	}

	var path string
	if err := db.QueryRow("SELECT image_path FROM tbl_item where id = ?", req.Id).Scan(&path); err != nil {
		log.Log.Errorln("Error scanning query GetFilePath")
		return "", err
	}

	return path, nil
}

func (t *AdminRepositories) UpdateItem(c *gin.Context, req models.ReqUpdateItem, imagePath string, uid int) error {
	// Get connection DB
	db, err := mysql.GetConnectionItem()
	if err != nil {
		log.Log.Errorln("Error GetConnectionItem", err.Error())
		return err
	}

	if _, err = db.Exec("UPDATE tbl_item SET name = ?, price = ?, stock_quantity = ?, id_category = ?, description = ?, image_path = ?, user_id = ? WHERE id = ?", req.Name, req.Price, req.Stock, req.IdCategory, req.Description, imagePath, uid, req.Id); err != nil {
		log.Log.Errorln("Error running query UpdateItem : ", err.Error())
		return err
	}

	return nil
}

func (t *AdminRepositories) DeleteItem(c *gin.Context, req models.ReqDeleteItem) error {
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

func (t *AdminRepositories) LogoutAdmin(c *gin.Context, uid int) error {
	// Get connection DB
	db, err := mysql.GetConnectionAdmin()
	if err != nil {
		log.Log.Errorln("Error GetConnectionAdmin", err.Error())
		return err
	}

	var id int
	// Get the latest session that inserted into tbl_session
	if err := db.QueryRow("SELECT id FROM tbl_session ORDER BY expired_time DESC LIMIT 1").Scan(&id); err != nil {
		log.Log.Errorln("Error scanning query LogoutAdmin : ")
		return err
	}

	if _, err := db.Exec("UPDATE tbl_session SET status = 0 WHERE id = ?", id); err != nil {
		log.Log.Errorln("Error running query LogoutAdmin : ")
		return err
	}

	return nil
}
