package constant

const (
	TokenExpireCode = 441
	TokenValid      = 241
	TokenNotValid   = 451

	/**设置token起效时间,单位:秒*/
	TokenBefore = 1000
	/**设置token有效时间, 单位:秒*/
	TokenExpire       = 1800
	EmptyString       = ""
	DefaultHarborAddr = "reg.image.com"
)

var (
	WhiteUrlList = []string{"/login", "dad/addsa"}
)
