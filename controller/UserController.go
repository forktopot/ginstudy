package controller

import (
	"bytes"
	"ginstudy/common"
	"ginstudy/dto"
	"ginstudy/model"
	"ginstudy/response"
	"ginstudy/util"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dchest/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

//注册用户
func Register(c *gin.Context) {
	DB := common.GetDB()
	//获取参数
	name := c.PostForm("name")
	telephone := c.PostForm("telephone")
	password := c.PostForm("password")

	//数据验证
	if len(telephone) != 11 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码不能少于六位")
		return
	}
	if len(name) == 0 {
		name = util.RandomString(10)
	}

	log.Println(name, telephone, password)

	//判断手机号是否存在
	if isTelephoneExist(DB, telephone) {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "用户已经存在")
		return
	}

	//创建用户并对密码加密
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "加密错误")
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
	}

	DB.Create(&newUser)

	//返回结果
	c.JSON(200, gin.H{"code": 200, "msg": "注册成功"})

}

//登录
func Login(c *gin.Context) {
	DB := common.GetDB()

	//获取参数
	telephone := c.PostForm("telephone")
	password := c.PostForm("password")
	captcha := c.PostForm("captcha")

	//数据验证
	if len(telephone) != 11 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
		return
	}
	if len(password) < 6 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于六位"})
		return
	}

	//判断验证码是否正确
	if !CaptchaVerify(c, captcha) {
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": "验证码错误"})
		return
	}

	//判断手机号是否存在
	var user model.User
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户不存在"})
	}

	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "密码错误"})
		return
	}

	//发送token
	token, err := common.ReleaseToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "系统异常"})
		log.Printf("token generate error: %v", err)
		return
	}

	//返回结果
	c.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{"token": token},
		"msg":  "登录成功",
	})

}

//获取用户信息
func Info(c *gin.Context) {
	user, _ := c.Get("user")

	//user.(model.User) 接口断言，将接口实际类型值返回，参考：https://studygolang.com/articles/20238
	c.JSON(http.StatusOK, gin.H{"user": dto.ToUserDto(user.(model.User))})
}

//判读是否存在该手机号的用户
func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}

	return false
}

//文件上传
func Upload(c *gin.Context) {
	//获取文件头
	file, err := c.FormFile("file")
	if err != nil {
		c.String(500, "上传图片出错")
	}
	//获取文件后缀
	// fileSuffix := path.Ext(file.Filename)
	// if fileSuffix != ".jpg" {
	// 	c.JSON(http.StatusOK, gin.H{"code": "200", "msg": "上传失败，只允许上传.jpg后缀的文件"})
	// 	return
	// }
	// c.JSON(200, gin.H{"message": file.Header.Context})
	c.SaveUploadedFile(file, "./uploadfile/"+file.Filename)
	c.JSON(http.StatusOK, gin.H{"code": "200", "msg": "上传成功"})
}

//文件下载
func HandleDownloadFile(c *gin.Context) {
	filename := c.Query("file")
	file, err := os.Open("./uploadfile/" + filename)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "文件未找到"})
		return
	}
	defer file.Close()
	bytes, _ := ioutil.ReadAll(file)
	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	// c.Header("Accept-Length", fmt.Sprintf("%d", len(file)))
	c.Writer.Write([]byte(bytes))
}

//验证码
func Captcha(c *gin.Context) {
	length := make([]int, 1)
	length[0] = 4
	l := captcha.DefaultLen
	w, h := 107, 36
	if len(length) == 1 {
		l = length[0]
	}
	if len(length) == 2 {
		w = length[1]
	}
	if len(length) == 3 {
		h = length[2]
	}
	captchaId := captcha.NewLen(l)
	session := sessions.Default(c)
	session.Set("captcha", captchaId)
	_ = session.Save()
	_ = Serve(c.Writer, c.Request, captchaId, ".png", "zh", false, w, h)
}

//验证码验证
func CaptchaVerify(c *gin.Context, code string) bool {
	session := sessions.Default(c)
	captchaId := session.Get("captcha")
	if captchaId != nil {
		session.Delete("captcha")
		_ = session.Save()
		if captcha.VerifyString(captchaId.(string), code) {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func Serve(w http.ResponseWriter, r *http.Request, id, ext, lang string, download bool, width, height int) error {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	var content bytes.Buffer
	switch ext {
	case ".png":
		w.Header().Set("Content-Type", "image/png")
		_ = captcha.WriteImage(&content, id, width, height)
	case ".wav":
		w.Header().Set("Content-Type", "audio/x-wav")
		_ = captcha.WriteAudio(&content, id, lang)
	default:
		return captcha.ErrNotFound
	}

	if download {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	http.ServeContent(w, r, id+ext, time.Time{}, bytes.NewReader(content.Bytes()))
	return nil
}
