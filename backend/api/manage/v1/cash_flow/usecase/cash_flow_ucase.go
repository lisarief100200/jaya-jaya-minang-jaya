package usecase

import (
	"backend/domain"
	"backend/models"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CashFlowUsecase struct {
	CFRepositories domain.CashFlowRepositories
}

func NewCashFlowUsecase(cashFlowRepo domain.CashFlowRepositories) domain.CashFlowUsecase {
	return &CashFlowUsecase{
		CFRepositories: cashFlowRepo,
	}
}

func (r *CashFlowUsecase) GetCashFlow(c *gin.Context, uid, level string) (models.RespGetCashFlow, error) {
	var list models.RespGetCashFlow
	uidInt, _ := strconv.Atoi(uid)

	list.IncomeInfo = []models.RespIncomeDetail{}
	list.ExpensesInfo = []models.RespExpensesDetail{}

	if level == "admin" {
		incomes, err := r.CFRepositories.GetIncome(c)
		if err != nil {
			return models.RespGetCashFlow{}, err
		}

		expenses, err := r.CFRepositories.GetExpenses(c)
		if err != nil {
			return models.RespGetCashFlow{}, err
		}

		// Menggabungkan data incomes ke dalam IncomeInfo
		for _, income := range incomes {
			list.IncomeInfo = append(list.IncomeInfo, income)
		}

		// Menggabungkan data expenses ke dalam ExpensesInfo
		for _, expense := range expenses {
			list.ExpensesInfo = append(list.ExpensesInfo, expense)
		}
	} else {
		incomes, err := r.CFRepositories.GetIncomeById(c, uidInt)
		if err != nil {
			return models.RespGetCashFlow{}, err
		}

		expenses, err := r.CFRepositories.GetExpensesById(c, uidInt)
		if err != nil {
			return models.RespGetCashFlow{}, err
		}

		fmt.Println(incomes, expenses)
	}

	return list, nil
}
