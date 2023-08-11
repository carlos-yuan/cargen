package captcha

import (
	"comm/aes"
	"comm/md5"
	rd "comm/random"
	"comm/timeUtil"
	"errors"
	"image/color"
	"strconv"
	"strings"

	"github.com/coocood/freecache"

	"github.com/google/uuid"
	"github.com/mojocn/base64Captcha"
)

// SetStore 设置store
func SetStore(s base64Captcha.Store) {
	base64Captcha.DefaultMemStore = s
}

// configJsonBody json request body.
type configJsonBody struct {
	Id            string
	CaptchaType   string
	VerifyValue   string
	DriverAudio   *base64Captcha.DriverAudio
	DriverString  *base64Captcha.DriverString
	DriverChinese *base64Captcha.DriverChinese
	DriverMath    *base64Captcha.DriverMath
	DriverDigit   *base64Captcha.DriverDigit
}

func DriverStringFunc() (id, b64s string, err error) {
	e := configJsonBody{}
	e.Id = uuid.New().String()
	e.DriverString = base64Captcha.NewDriverString(46, 140, 2, 2, 4, "234567890abcdefghjkmnpqrstuvwxyz", &color.RGBA{240, 240, 246, 246}, nil, []string{"wqy-microhei.ttc"})
	driver := e.DriverString.ConvertFonts()
	cap := base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore)
	return cap.Generate()
}

func DriverDigitFunc() (id, b64s string, err error) {
	e := configJsonBody{}
	e.Id = uuid.New().String()
	e.DriverDigit = base64Captcha.NewDriverDigit(80, 240, 4, 0.7, 80)
	driver := e.DriverDigit
	cap := base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore)
	return cap.Generate()
}

var cache = freecache.NewCache(1 * 1024 * 1024)

func AesDigit(key string) (id, b64s string, err error) {
	if len(key) != 32 {
		key = md5.Encode(key)
	}
	driver := base64Captcha.NewDriverDigit(80, 240, 4, 0.7, 80)
	number := rd.GetNumber(4)
	dr, err := driver.DrawCaptcha(number)
	if err != nil {
		return
	}
	b64s = dr.EncodeB64string()
	id, err = aes.EncryptCBC5([]byte(number+"-"+strconv.Itoa(int(timeUtil.Milli()+expireTime))), key[:16], key[16:])
	return
}

const expireTime = timeUtil.Minute3

var CodeExpire = errors.New("code expire")
var CodeErr = errors.New("code error")

func AesVerify(key, id, code string) error {
	if len(key) != 32 {
		key = md5.Encode(key)
	}
	id, err := aes.DecryptCBC5(id, key[:16], key[16:])
	if err != nil {
		return CodeErr
	}
	numberExpire := strings.Split(id, "-")
	if len(numberExpire) != 2 {
		return CodeErr
	}
	if numberExpire[0] != code {
		return CodeErr
	}
	expire, err := strconv.ParseInt(numberExpire[1], 10, 64)
	if err != nil {
		return CodeErr
	}
	if expire > timeUtil.Milli() {
		//缓存判断是否存在，单机中如果已存在则返回失败
		_, err := cache.Get([]byte(id))
		if err != nil && err.Error() != "Entry not found" {
			return CodeErr
		}
		if err != nil && err.Error() == "Entry not found" {
			err = cache.Set([]byte(id), []byte(code), int(expireTime/timeUtil.Second))
			if err != nil {
				return CodeErr
			}
			return nil
		} else {
			return CodeExpire
		}
	} else {
		return CodeExpire
	}
}

// Verify 校验验证码
func Verify(id, code string, clear bool) bool {
	return base64Captcha.DefaultMemStore.Verify(id, code, clear)
}
