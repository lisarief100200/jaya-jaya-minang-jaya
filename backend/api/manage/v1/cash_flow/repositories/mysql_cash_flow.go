package repositories

import (
	"backend/domain"
	"backend/models"
	"backend/pkg/v1/mysql"
	"backend/utils/log"
	"database/sql"

	"github.com/gin-gonic/gin"
)

type CashFlowRepositories struct {
	sql *sql.DB
}

func NewTestRepoCashFlow(sql *sql.DB) domain.CashFlowRepositories {
	return &CashFlowRepositories{
		sql: sql,
	}
}

func (t *CashFlowRepositories) GetIncome(c *gin.Context) ([]models.RespIncomeDetail, error) {
	// Get connection DB
	db, err := mysql.GetConnectionCashFlow()
	if err != nil {
		log.Log.Errorln("Error GetConnectionCashFlow", err.Error())
		return []models.RespIncomeDetail{}, err
	}

	rows, err := db.Query("SELECT id, name_income, type_income, date, total_income, desc_income FROM tbl_income")
	if err != nil {
		if err == sql.ErrNoRows {
			return []models.RespIncomeDetail{}, err
		}
		return []models.RespIncomeDetail{}, err
	}
	defer rows.Close()

	// Mapping to struct rows
	var getIncomes []models.RespIncomeDetail

	for rows.Next() {
		getIncome := models.RespIncomeDetail{}
		if errScan := rows.Scan(&getIncome.Id, &getIncome.NameIncome, &getIncome.TypeIncome, &getIncome.Date, &getIncome.TotalIncome, &getIncome.DescIncome); errScan != nil {
			return []models.RespIncomeDetail{}, err
		}
		getIncomes = append(getIncomes, getIncome)
	}

	if errRows := rows.Err(); errRows != nil {
		return []models.RespIncomeDetail{}, err
	}

	return getIncomes, nil
}

func (t *CashFlowRepositories) GetIncomeById(c *gin.Context, uid int) ([]models.RespIncomeDetail, error) {
	// Get connection DB
	db, err := mysql.GetConnectionCashFlow()
	if err != nil {
		log.Log.Errorln("Error GetConnectionCashFlow", err.Error())
		return []models.RespIncomeDetail{}, err
	}

	rows, err := db.Query("SELECT id, name_income, type_income, date, total_income, desc_income FROM tbl_income WHERE user_id = ?", uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return []models.RespIncomeDetail{}, err
		}
		return []models.RespIncomeDetail{}, err
	}
	defer rows.Close()

	// Mapping to struct rows
	var getIncomes []models.RespIncomeDetail

	for rows.Next() {
		getIncome := models.RespIncomeDetail{}
		if errScan := rows.Scan(&getIncome.Id, &getIncome.NameIncome, &getIncome.TypeIncome, &getIncome.Date, &getIncome.TotalIncome, &getIncome.DescIncome); errScan != nil {
			return []models.RespIncomeDetail{}, err
		}
		getIncomes = append(getIncomes, getIncome)
	}

	if errRows := rows.Err(); errRows != nil {
		return []models.RespIncomeDetail{}, err
	}

	return getIncomes, nil
}

func (t *CashFlowRepositories) GetExpenses(c *gin.Context) ([]models.RespExpensesDetail, error) {
	// Get connection DB
	db, err := mysql.GetConnectionCashFlow()
	if err != nil {
		log.Log.Errorln("Error GetConnectionCashFlow", err.Error())
		return []models.RespExpensesDetail{}, err
	}

	rows, err := db.Query("SELECT id, name_expenses, type_expenses, date, total_expenses, desc_expenses FROM tbl_expenses")
	if err != nil {
		if err == sql.ErrNoRows {
			return []models.RespExpensesDetail{}, err
		}
		return []models.RespExpensesDetail{}, err
	}
	defer rows.Close()

	// Mappung to struct rows
	var getExpenses []models.RespExpensesDetail

	for rows.Next() {
		getExpense := models.RespExpensesDetail{}
		if errScan := rows.Scan(&getExpense.Id, &getExpense.NameExpenses, &getExpense.TypeExpenses, &getExpense.Date, &getExpense.TotalExpenses, &getExpense.DescExpenses); errScan != nil {
			return []models.RespExpensesDetail{}, err
		}
		getExpenses = append(getExpenses, getExpense)
	}

	if errRows := rows.Err(); errRows != nil {
		return []models.RespExpensesDetail{}, err
	}

	return getExpenses, nil
}

func (t *CashFlowRepositories) GetExpensesById(c *gin.Context, uid int) ([]models.RespExpensesDetail, error) {
	// Get connection DB
	db, err := mysql.GetConnectionCashFlow()
	if err != nil {
		log.Log.Errorln("Error GetConnectionCashFlow", err.Error())
		return []models.RespExpensesDetail{}, err
	}

	rows, err := db.Query("SELECT id, name_expenses, type_expenses, date, total_expenses, desc_expenses FROM tbl_expenses WHERE user_id = ?", uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return []models.RespExpensesDetail{}, err
		}
		return []models.RespExpensesDetail{}, err
	}
	defer rows.Close()

	// Mapping to struct rows
	var getExpenses []models.RespExpensesDetail

	for rows.Next() {
		getExpense := models.RespExpensesDetail{}
		if errScan := rows.Scan(&getExpense.Id, &getExpense.NameExpenses, &getExpense.TypeExpenses, &getExpense.Date, &getExpense.TotalExpenses, &getExpense.DescExpenses); errScan != nil {
			return []models.RespExpensesDetail{}, err
		}
		getExpenses = append(getExpenses, getExpense)
	}

	if errRows := rows.Err(); errRows != nil {
		return []models.RespExpensesDetail{}, err
	}

	return getExpenses, nil
}
