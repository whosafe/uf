package i18n

// zh 中文错误消息
// 注意：此 map 为只读，不可在运行时修改，以保证并发安全
var zh = map[string]string{
	// Required 规则
	"required": "{field}不能为空",

	// Min 规则
	"min":        "{field}不能少于{param}",
	"min_length": "{field}长度不能少于{param}个字符",

	// Max 规则
	"max":        "{field}不能超过{param}",
	"max_length": "{field}长度不能超过{param}个字符",

	// Len 规则
	"len": "{field}长度必须为{param}个字符",

	// 比较规则
	"gt":  "{field}必须大于{param}",
	"gte": "{field}必须大于或等于{param}",
	"lt":  "{field}必须小于{param}",
	"lte": "{field}必须小于或等于{param}",

	// Email 规则
	"email": "{field}必须是有效的邮箱地址",

	// URL 规则
	"url": "{field}必须是有效的URL",

	// Phone 规则
	"phone": "{field}必须是有效的手机号",

	// Alpha 规则
	"alpha":    "{field}只能包含字母",
	"alphanum": "{field}只能包含字母和数字",
	"numeric":  "{field}只能包含数字",

	// OneOf 规则
	"oneof": "{field}必须是以下值之一: {param}",

	// 字符串匹配规则
	"contains":    "{field}必须包含{param}",
	"starts_with": "{field}必须以{param}开头",
	"ends_with":   "{field}必须以{param}结尾",

	// Regex 规则
	"regex": "{field}必须匹配模式{param}",

	// 数字验证规则
	"between":  "{field}必须在{min}和{max}之间",
	"positive": "{field}必须是正数",
	"negative": "{field}必须是负数",
	"integer":  "{field}必须是整数",
	"decimal":  "{field}必须是小数,最多{param}位小数",

	// 字符串增强规则
	"uuid":      "{field}必须是有效的UUID",
	"json":      "{field}必须是有效的JSON格式",
	"base64":    "{field}必须是有效的Base64编码",
	"lowercase": "{field}必须是小写",
	"uppercase": "{field}必须是大写",
	"ascii":     "{field}只能包含ASCII字符",
	"not_blank": "{field}不能为空白",

	// 网络相关规则
	"ip":     "{field}必须是有效的IP地址",
	"ipv4":   "{field}必须是有效的IPv4地址",
	"ipv6":   "{field}必须是有效的IPv6地址",
	"mac":    "{field}必须是有效的MAC地址",
	"domain": "{field}必须是有效的域名",
	"port":   "{field}必须是有效的端口号(1-65535)",

	// 日期时间规则
	"date":         "{field}必须是有效的日期,格式为{param}",
	"datetime":     "{field}必须是有效的日期时间,格式为{param}",
	"date_before":  "{field}必须早于{param}",
	"date_after":   "{field}必须晚于{param}",
	"date_between": "{field}必须在{min}和{max}之间",

	// 中国特色规则
	"idcard":        "{field}必须是有效的身份证号",
	"bankcard":      "{field}必须是有效的银行卡号",
	"social_credit": "{field}必须是有效的统一社会信用代码",
	"postalcode":    "{field}必须是有效的邮政编码",
	"chinese_name":  "{field}必须是有效的中文姓名",

	// 文件相关规则
	"file_extension": "{field}必须是以下扩展名之一: {param}",
	"mime_type":      "{field}必须是以下MIME类型之一: {param}",
	"file_size":      "{field}大小必须在{min}和{max}字节之间",

	// 集合/数组规则
	"unique":         "{field}必须包含唯一值",
	"array_min":      "{field}至少包含{param}个元素",
	"array_max":      "{field}最多包含{param}个元素",
	"array_contains": "{field}必须包含{param}",

	// 安全相关规则
	"strong_password": "{field}必须是强密码(至少8位,包含大写、小写、数字和特殊字符)",
	"no_html":         "{field}不能包含HTML标签",
	"no_sql":          "{field}不能包含SQL注入字符",
	"no_xss":          "{field}不能包含XSS攻击字符",

	// 其他实用规则
	"confirmed": "{field}确认不匹配",
	"distinct":  "{field}必须不同于{param}",
	"not_in":    "{field}不能是以下值之一: {param}",
	"nullable":  "{field}可以为空",
}
