package routes

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	api "mall/api/v1"
	"mall/middleware"
)

// router
func NewRouter() *gin.Engine {
	r := gin.Default()                                        // Create a new Gin router with default middleware
	store := cookie.NewStore([]byte("something-very-secret")) // Create a new cookie store for sessions
	r.Use(middleware.Cors())                                  // Apply CORS middleware to all routes
	r.Use(sessions.Sessions("mysession", store))              // Apply session middleware to all routes
	r.StaticFS("/static", http.Dir("./static"))
	v1 := r.Group("api/v1")
	{

		v1.GET("ping", func(c *gin.Context) {
			c.JSON(200, "success")
		})

		// 用户操作
		v1.POST("user/register", api.UserRegister)
		v1.POST("user/login", api.UserLogin)

		// 商品操作
		v1.GET("products", api.ListProducts)
		v1.GET("product/:id", api.ShowProduct)
		v1.POST("products", api.SearchProducts)
		v1.GET("imgs/:id", api.ListProductImg)   // 商品图片
		v1.GET("categories", api.ListCategories) // 商品分类
		v1.GET("carousels", api.ListCarousels)   // 轮播图

		authed := v1.Group("/")      // Create a subgroup for routes that require authentication
		authed.Use(middleware.JWT()) // Apply JWT middleware to all routes in the subgroup
		{

			// 用户操作
			authed.PUT("user", api.UserUpdate)
			authed.POST("user/sending-email", api.SendEmail)
			authed.POST("user/valid-email", api.ValidEmail)
			authed.POST("avatar", api.UploadAvatar) // 上传头像

			// 商品操作
			authed.POST("product", api.CreateProduct)
			authed.PUT("product/:id", api.UpdateProduct)
			authed.DELETE("product/:id", api.DeleteProduct)
			// 收藏夹
			authed.GET("favorites", api.ShowFavorites)
			authed.POST("favorites", api.CreateFavorite)
			authed.DELETE("favorites/:id", api.DeleteFavorite)

			// 订单操作
			authed.POST("orders", api.CreateOrder)
			authed.GET("orders", api.ListOrders)
			authed.GET("orders/:id", api.ShowOrder)
			authed.DELETE("orders/:id", api.DeleteOrder)

			// 购物车
			authed.POST("carts", api.CreateCart)
			authed.GET("carts", api.ShowCarts)
			authed.PUT("carts/:id", api.UpdateCart) // 购物车id
			authed.DELETE("carts/:id", api.DeleteCart)

			// 收获地址操作
			authed.POST("addresses", api.CreateAddress)
			authed.GET("addresses/:id", api.GetAddress)
			authed.GET("addresses", api.ListAddress)
			authed.PUT("addresses/:id", api.UpdateAddress)
			authed.DELETE("addresses/:id", api.DeleteAddress)

			// 支付功能
			authed.POST("paydown", api.OrderPay)

			// 显示金额
			authed.POST("money", api.ShowMoney)

			// 秒杀专场
			authed.POST("import_skill_goods", api.ImportSkillGoods)
			authed.POST("init_skill_goods", api.InitSkillGoods)
			authed.POST("skill_goods", api.SkillGoods)
		}
	}
	return r
}
