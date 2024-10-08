package routes

import (
	"gin-mall/api"
	"gin-mall/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	//store := cookie.NewStore([]byte("secret"))
	//r.Use(sessions.Sessions("mysession", store))
	//r.StaticFS("/static", http.Dir("./static"))
	v1 := r.Group("api/v1")
	{
		//用户操作
		v1.POST("user/register", api.UserRegisterHandler())
		v1.POST("use/login", api.UserLoginHandler())

		//商品操作
		v1.GET("product/list", api.ListProductHandler())
		v1.GET("product/show", api.ShowProductHandler())
		v1.POST("product/search", api.SearchProductsHandler())
		v1.GET("product/imgs/list", api.ListProductImgHandler()) // 商品图片
		v1.GET("category/list", api.ListCategoryHandler())       // 商品分类
		v1.GET("carousels", api.ListCarouselsHandler())          // 轮播图

		authed := v1.Group("/")
		authed.Use(middleware.AuthMiddleware())
		{
			authed.POST("user/update", api.UserUpdateHandler())
			authed.POST("user/showInfo", api.ShowUserInfoHandler())
			authed.POST("user/send_email", api.SendEmailHandler())
			authed.POST("user/valid_email", api.ValidEmailHandler())
			authed.POST("user/follwing", api.UserFollowingHnadler())
			authed.POST("user/unfollowing", api.UserUnFollowingHnadler())

			//商品操作
			authed.POST("product/create", api.CreateProductHandler())
			authed.POST("product/update", api.UpdateProductHandler())
			authed.POST("product/delete", api.DeleteProductHandler())

			//收藏夹
			authed.POST("favorites/create", api.CreateFavoriteHandler())
			authed.GET("favorites/list", api.ListFavoritesHandler())
			authed.POST("favorites/delete", api.DeleteFavoriteHandler())

			//订单操作
			authed.POST("orders/create", api.CreateOrderHandler())
			authed.GET("orders/list", api.ListOrdersHandler())
			authed.GET("orders/show", api.ShowOrderHandler())
			authed.POST("orders/delete", api.DeleteOrderHandler())

			//购物车操作
			// 购物车
			authed.POST("carts/create", api.CreateCartHandler())
			authed.GET("carts/list", api.ListCartHandler())
			authed.POST("carts/update", api.UpdateCartHandler()) // 购物车id
			authed.POST("carts/delete", api.DeleteCartHandler())

			//收货地址
			// 收获地址操作
			authed.POST("addresses/create", api.CreateAddressHandler())
			authed.GET("addresses/show", api.ShowAddressHandler())
			authed.GET("addresses/list", api.ListAddressHandler())
			authed.POST("addresses/update", api.UpdateAddressHandler())
			authed.POST("addresses/delete", api.DeleteAddressHandler())

			// 支付功能
			authed.POST("paydown", api.OrderPaymentHandler())

			//显示金额
			authed.POST("money", api.ShowMoneyHandler())

			//秒杀
			// 秒杀专场
			authed.POST("skill_product/init", api.InitSkillProductHandler())
			authed.GET("skill_product/list", api.ListSkillProductHandler())
			authed.GET("skill_product/show", api.GetSkillProductHandler())
			authed.POST("skill_product/skill", api.SkillProductHandler())
		}
	}
	return r
}
