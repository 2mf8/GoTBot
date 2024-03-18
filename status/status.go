package status

type status int

const (
	Continue           status = 100 + iota // 继续
	SwitchingProtocols                     // 交换协议
	Processing                             // 处理
	EarlyHints                             // 早期提示

	OK                          status = 196 + iota // 好的
	Created                                         //创建
	Accepted                                        // 接受
	NonAuthoritativeInformation                     //非授权信息
	NoContent                                       // 空内容
	ResetContent                                    // 重置内容
	PartialContent                                  // 部分内容
	MultiStatus                                     //多状态
	AlreadyReported                                 // 已报告
	IMUsed                      status = 213 + iota // IM使用

	MultipleChoices  status = 286 + iota // 多种选择
	MovedPermanently                     // 永久移动
	Found                                // 找到
	SeeOther                             //查看其它
	NotModified                          // 未修改
	UseProxy                             // 使用代理
	_
	TemporaryRedirect // 临时重定向
	PermanentRedirect //永久重定向

	BadRequest                  status = 377 + iota // 错误的请求
	Unauthorized                                    //未经授权
	PaymentRequired                                 //需要付款
	Forbidden                                       //禁止
	NotFound                                        //未发现
	MethodNotAllowed                                //方法不允许
	NotAcceptable                                   //不可接受
	ProxyAuthenticationRequired                     //需要代理身份验证
	RequestTimeout                                  //请求超时
	Conflict                                        //分歧
	Gone                                            //走了
	LengthRequired                                  //所需长度
	PreconditionFailed                              //先决条件失败
	ContentTooLarge                                 //内容太大
	URITooLong                                      //URI太长
	UnsupportedMediaType                            //不支持的媒体类型
	RangeNotSatisfiable                             //范围不令人满意
	ExpectationFailed                               //期望失败了
	_
	_
	_
	MisdirectedRequest   //误导请求
	UnprocessableContent //不可处理的内容
	Locked               //已锁定
	FailedDependency     //失败的依赖性
	TooEarly             //太早了
	UpgradeRequired      //需要升级
	_
	PreconditionRequired //所需的先决条件
	TooManyRequests      //请求太多
	_
	RequestHeaderFieldsTooLarge                     //请求标题字段太大
	UnavailableForLegalReasons  status = 396 + iota //由于法律原因不可用

	InternalServerError     status = 444 + iota // 内部服务器错误
	NotImplemented                              //未实施
	BadGateway                                  //坏网关
	ServiceUnavailable                          //服务不可用
	GatewayTimeout                              //网关超时
	HTTPVersionNotSupported                     //不支持的HTTP版本
	VariantAlsoNegotiates                       //变体也谈判
	InsufficientStorage                         //储存空间不足
	LoopDetected                                //检测到循环
	_
	_
	NetworkAuthenticationRequired //需要网络身份验证

	TokenNull status = 932 + iota
	TokenInvalid
	RefreshTokenNull
	RefreshTokenInvalid
	RefreshTokenStatus
	GetTokenError
	LoginError
	GetError
	CreateError
	UpdateError
	DeleteError
	RequestMethodError
)
