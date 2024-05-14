package model

import "errors"

var (
	ErrUserError            = errors.New("an error occurred on user")
	ErrUserNotFound         = errors.New("user not found")
	ErrPasswordError        = errors.New("invalid password")
	ErrOldPasswordError     = errors.New("invalid old password")
	ErrInsertCompanyError   = errors.New("error on insert company")
	ErrSelectCompaniesError = errors.New("error on select companies")
	ErrBrandExists          = errors.New("error brand exists")
	ErrBrandError           = errors.New("brand error")
	ErrCompanyExists        = errors.New("error company exists")
	ErrCompanyError         = errors.New("error company error")
	ErrBillError            = errors.New("bill error")
	ErrStoreError           = errors.New("store error")
	ErrProductError         = errors.New("product error")
	ErrNotExistsError       = errors.New("product not exists")
	ErrUserProductError     = errors.New("user product error")
)
