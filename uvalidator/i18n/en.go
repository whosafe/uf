package i18n

// en 英文错误消息
// 注意：此 map 为只读，不可在运行时修改，以保证并发安全
var en = map[string]string{
	// Required 规则
	"required": "{field} is required",

	// Min 规则
	"min":        "{field} must be at least {param}",
	"min_length": "{field} must be at least {param} characters",

	// Max 规则
	"max":        "{field} must be at most {param}",
	"max_length": "{field} must be at most {param} characters",

	// Len 规则
	"len": "{field} must be exactly {param} characters",

	// 比较规则
	"gt":  "{field} must be greater than {param}",
	"gte": "{field} must be greater than or equal to {param}",
	"lt":  "{field} must be less than {param}",
	"lte": "{field} must be less than or equal to {param}",

	// Email 规则
	"email": "{field} must be a valid email address",

	// URL 规则
	"url": "{field} must be a valid URL",

	// Phone 规则
	"phone": "{field} must be a valid phone number",

	// Alpha 规则
	"alpha":    "{field} must contain only letters",
	"alphanum": "{field} must contain only letters and numbers",
	"numeric":  "{field} must contain only numbers",

	// OneOf 规则
	"oneof": "{field} must be one of: {param}",

	// 字符串匹配规则
	"contains":    "{field} must contain {param}",
	"starts_with": "{field} must start with {param}",
	"ends_with":   "{field} must end with {param}",

	// Regex 规则
	"regex": "{field} must match pattern {param}",

	// 数字验证规则
	"between":  "{field} must be between {min} and {max}",
	"positive": "{field} must be a positive number",
	"negative": "{field} must be a negative number",
	"integer":  "{field} must be an integer",
	"decimal":  "{field} must be a decimal number with at most {param} decimal places",

	// 字符串增强规则
	"uuid":      "{field} must be a valid UUID",
	"json":      "{field} must be valid JSON",
	"base64":    "{field} must be valid Base64 encoding",
	"lowercase": "{field} must be lowercase",
	"uppercase": "{field} must be uppercase",
	"ascii":     "{field} must contain only ASCII characters",
	"not_blank": "{field} must not be blank",

	// 网络相关规则
	"ip":     "{field} must be a valid IP address",
	"ipv4":   "{field} must be a valid IPv4 address",
	"ipv6":   "{field} must be a valid IPv6 address",
	"mac":    "{field} must be a valid MAC address",
	"domain": "{field} must be a valid domain name",
	"port":   "{field} must be a valid port number (1-65535)",

	// 日期时间规则
	"date":         "{field} must be a valid date in format {param}",
	"datetime":     "{field} must be a valid datetime in format {param}",
	"date_before":  "{field} must be before {param}",
	"date_after":   "{field} must be after {param}",
	"date_between": "{field} must be between {min} and {max}",

	// 中国特色规则
	"idcard":        "{field} must be a valid ID card number",
	"bankcard":      "{field} must be a valid bank card number",
	"social_credit": "{field} must be a valid unified social credit code",
	"postalcode":    "{field} must be a valid postal code",
	"chinese_name":  "{field} must be a valid Chinese name",

	// 文件相关规则
	"file_extension": "{field} must have one of the following extensions: {param}",
	"mime_type":      "{field} must be one of the following MIME types: {param}",
	"file_size":      "{field} size must be between {min} and {max} bytes",

	// 集合/数组规则
	"unique":         "{field} must contain unique values",
	"array_min":      "{field} must contain at least {param} items",
	"array_max":      "{field} must contain at most {param} items",
	"array_contains": "{field} must contain {param}",

	// 安全相关规则
	"strong_password": "{field} must be a strong password (at least 8 characters with uppercase, lowercase, digit, and special character)",
	"no_html":         "{field} must not contain HTML tags",
	"no_sql":          "{field} must not contain SQL injection characters",
	"no_xss":          "{field} must not contain XSS attack characters",

	// 其他实用规则
	"confirmed": "{field} confirmation does not match",
	"distinct":  "{field} must be different from {param}",
	"not_in":    "{field} must not be one of: {param}",
	"nullable":  "{field} can be null",
}
