package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"shop-aggregator/internal/auth"
)

type AuthStorer interface {
	EnsureValidToken(context.Context, string) (uuid.UUID, error)
}

type AuthHandler interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
}

type UserHandler interface {
	CreateUser(c *gin.Context)
	UpdatePassword(c *gin.Context)
	GetUser(c *gin.Context)
	UpdateEmail(c *gin.Context)
}

type BrandHandler interface {
	Create(c *gin.Context)
	GetByPartialName(c *gin.Context)
}

type CompanyHandler interface {
	GetByPartialName(c *gin.Context)
}

type BillHandler interface {
	Start(c *gin.Context)
	Close(c *gin.Context)
	GetBillsByUserID(c *gin.Context)
	GetLastBill(c *gin.Context)
	Cancel(c *gin.Context)
}

type StoreHandler interface {
	CreateStore(c *gin.Context)
	GetStoreByZipCodeOrName(c *gin.Context)
}

type ProductHandler interface {
	Create(c *gin.Context)
	GetProductByEAN(c *gin.Context)
}

type UserProductHandler interface {
	Create(c *gin.Context)
	SelectProductsByBillID(c *gin.Context)
	UpdateQuantity(c *gin.Context)
	Delete(c *gin.Context)
}

type InitialisationHandler interface {
	AppInitialisation(c *gin.Context)
}

// NewRouter creates and configures a Gin engine instance
func NewRouter(
	router *gin.Engine,
	as AuthStorer,
	ah AuthHandler,
	uh UserHandler,
	bh BrandHandler,
	ch CompanyHandler,
	bih BillHandler,
	sh StoreHandler,
	ph ProductHandler,
	uph UserProductHandler,
	ih InitialisationHandler,
) *gin.Engine {
	router.GET("/init", ih.AppInitialisation)

	router.POST("/create-user", uh.CreateUser)
	router.POST("/login", ah.Login)

	protected := router.Group("/")
	protected.Use(auth.Middleware(as))

	user := protected.Group("/user")
	{
		user.GET("/get", uh.GetUser)
		user.POST("/logout", ah.Logout)
		user.POST("/reset-password", uh.UpdatePassword)
		user.POST("/update-email", uh.UpdateEmail)
	}

	brand := protected.Group("/brand")
	{
		brand.GET("/get/:name", bh.GetByPartialName)
		brand.POST("/create", bh.Create)
	}

	company := protected.Group("/company")
	{
		company.GET("/get/:name", ch.GetByPartialName)
	}

	bill := protected.Group("/bill")
	{
		bill.GET("/get-all", bih.GetBillsByUserID)
		bill.GET("/get-last", bih.GetLastBill)
		bill.POST("/start", bih.Start)
		bill.POST("/stop", bih.Close)
		bill.POST("/cancel", bih.Cancel)
	}

	store := protected.Group("/store")
	{
		store.GET("/get/:store_type/:search", sh.GetStoreByZipCodeOrName)
		store.POST("/create-store", sh.CreateStore)
	}

	product := protected.Group("/product")
	{
		product.GET("/get/:ean", ph.GetProductByEAN)
		product.POST("/create-product", ph.Create)
	}

	userProduct := protected.Group("/user-product")
	{
		userProduct.GET("/get-bill-id/:bill_id", uph.SelectProductsByBillID)
		userProduct.POST("/create-user-product", uph.Create)
		userProduct.PUT("/quantity", uph.UpdateQuantity)
		userProduct.DELETE("/delete/:user_product_id", uph.Delete)
	}

	return router
}
