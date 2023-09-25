package repositories

import (
	"backend/domain"
	"backend/models"
	"backend/pkg/v1/mysql"
	"backend/utils/log"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthRepositories struct {
	sql *sql.DB
}

func NewTestRepoAuth(sql *sql.DB) domain.AuthRepositories {
	return &AuthRepositories{
		sql: sql,
	}
}

func (t *AuthRepositories) CheckIdAdmin(c *gin.Context, username string) (string, error) {
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

func (t *AuthRepositories) CheckLoginAdmin(c *gin.Context, req models.ReqLoginAdmin, uid string) (exist bool, err error) {
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

func (t *AuthRepositories) GetPassword(c *gin.Context, uid string) (string, error) {
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

func (t *AuthRepositories) GenerateSessionID(c *gin.Context, processName string) (string, error) {
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

func (t *AuthRepositories) GetAdminProp(c *gin.Context, uid string) (models.AdminProp, error) {
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

func (t *AuthRepositories) UpdateLastLogin(c *gin.Context, uid string) error {
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

func (t *AuthRepositories) InsertSession(c *gin.Context, session, uid, tipe, token string, exptime time.Time) error {
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

func (t *AuthRepositories) LogoutAdmin(c *gin.Context, uid int) error {
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
