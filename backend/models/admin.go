package models

import "time"

type ReqLoginUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RespLoginUser struct {
	Token    string `json:"token"`
	RoleCode string `json:"roleCode"`
}

type UserProp struct {
	Name  string `json:"name"`
	Level string `json:"level"`
}

type ReqCreateItem struct {
	Name        string `json:"name" form:"name"`
	Price       int64  `json:"price" form:"price"`
	Stock       int64  `json:"stock" form:"stock"`
	IdCategory  int64  `json:"idCategory" form:"idCategory"`
	Description string `json:"description" form:"description"`
	//Image       string `json:"image" form:"image"`
}

type RespGetList struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Price       int64  `json:"price"`
	Stock       int64  `json:"stock"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type ReqUpdateItem struct {
	Id          int64  `json:"id" form:"id"`
	Name        string `json:"name" form:"name"`
	Price       int64  `json:"price" form:"price"`
	Stock       int64  `json:"stock" form:"stock"`
	IdCategory  int64  `json:"idCategory" form:"idCategory"`
	Description string `json:"description" form:"description"`
	//Image       string `json:"image" form:"image"`
}

type ReqDeleteItem struct {
	Id int64 `json:"id"`
}

type ReqCreateUtilities struct {
	Name        string `json:"name" form:"name"`
	Price       int64  `json:"price" form:"price"`
	Stock       int64  `json:"stock" form:"stock"`
	IdCategory  int64  `json:"idCategory" form:"idCategory"`
	Description string `json:"description" form:"description"`
	//Image       string `json:"image" form:"image"`
}

type ReqUpdateUtilities struct {
	Id          int64  `json:"id" form:"id"`
	Name        string `json:"name" form:"name"`
	Price       int64  `json:"price" form:"price"`
	Stock       int64  `json:"stock" form:"stock"`
	IdCategory  int64  `json:"idCategory" form:"idCategory"`
	Description string `json:"description" form:"description"`
	//Image       string `json:"image" form:"image"`
}

type ReqDeleteUtilities struct {
	Id int64 `json:"id"`
}

type RespGetCashFlow struct {
	IncomeInfo   []RespIncomeDetail   `json:"incomeDetail"`
	ExpensesInfo []RespExpensesDetail `json:"expensesDetail"`
}

type RespIncomeDetail struct {
	Id          int64     `json:"id"`
	NameIncome  string    `json:"nameIncome"`
	TypeIncome  string    `json:"typeIncome"`
	Date        time.Time `json:"date"`
	TotalIncome int64     `json:"totalIncome"`
	DescIncome  string    `json:"descIncome"`
}

type RespExpensesDetail struct {
	Id            int64     `json:"id"`
	NameExpenses  string    `json:"nameExpenses"`
	TypeExpenses  string    `json:"typeExpenses"`
	Date          time.Time `json:"date"`
	TotalExpenses int64     `json:"totalExpenses"`
	DescExpenses  string    `json:"descExpenses"`
}
