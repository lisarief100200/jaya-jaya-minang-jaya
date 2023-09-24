package models

type ReqLoginAdmin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RespLoginAdmin struct {
	Token    string `json:"token"`
	RoleCode string `json:"roleCode"`
}

type AdminProp struct {
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
