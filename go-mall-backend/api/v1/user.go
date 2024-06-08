package v1

import (
	"mall/consts"
	util "mall/pkg/utils"
	"mall/service"

	"github.com/gin-gonic/gin"
)

func UserRegister(c *gin.Context) {
	// Create an instance of UserService
	var userRegisterService service.UserService
	// Bind the request data to the UserService instance
	if err := c.ShouldBind(&userRegisterService); err == nil {
		res := userRegisterService.Register(c.Request.Context())
		c.JSON(consts.StatusOK, res)
	} else {
		c.JSON(consts.IlleageRequest, ErrorResponse(err))
		util.LogrusObj.Infoln(err)
	}
}

func UserLogin(c *gin.Context) {
	var userLoginService service.UserService
	if err := c.ShouldBind(&userLoginService); err == nil {
		res := userLoginService.Login(c.Request.Context())
		c.JSON(consts.StatusOK, res)
	} else {
		c.JSON(consts.IlleageRequest, ErrorResponse(err))
		util.LogrusObj.Infoln(err)
	}
}

func UserUpdate(c *gin.Context) {
	var userUpdateService service.UserService
	claims, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&userUpdateService); err == nil {
		res := userUpdateService.Update(c.Request.Context(), claims.ID)
		c.JSON(consts.StatusOK, res)
	} else {
		c.JSON(consts.IlleageRequest, ErrorResponse(err))
		util.LogrusObj.Infoln(err)
	}
}

func UploadAvatar(c *gin.Context) {
	file, fileHeader, _ := c.Request.FormFile("file")
	fileSize := fileHeader.Size
	uploadAvatarService := service.UserService{}
	chaim, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&uploadAvatarService); err == nil {
		res := uploadAvatarService.Post(c.Request.Context(), chaim.ID, file, fileSize)
		c.JSON(consts.StatusOK, res)
	} else {
		c.JSON(consts.IlleageRequest, ErrorResponse(err))
		util.LogrusObj.Infoln(err)
	}
}

func SendEmail(c *gin.Context) {
	var sendEmailService service.SendEmailService
	chaim, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&sendEmailService); err == nil {
		res := sendEmailService.Send(c.Request.Context(), chaim.ID)
		c.JSON(consts.StatusOK, res)
	} else {
		c.JSON(consts.IlleageRequest, ErrorResponse(err))
		util.LogrusObj.Infoln(err)
	}
}

func ValidEmail(c *gin.Context) {
	var vaildEmailService service.ValidEmailService
	if err := c.ShouldBind(vaildEmailService); err == nil {
		res := vaildEmailService.Valid(c.Request.Context(), c.GetHeader("Authorization"))
		c.JSON(consts.StatusOK, res)
	} else {
		c.JSON(consts.IlleageRequest, ErrorResponse(err))
		util.LogrusObj.Infoln(err)
	}
}
