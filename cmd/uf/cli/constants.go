package cli

// 目录常量
const (
	DirConfig     = "config"
	DirHack       = "hack"
	DirInternal   = "internal"
	DirUtility    = "utility"
	DirConsts     = "consts"
	DirHandler    = "handler"
	DirModel      = "model"
	DirDAO        = "dao"
	DirRouter     = "router"
	DirService    = "service"
	DirLogic      = "logic"
	DirMiddleware = "middleware"
)

// 文件常量
const (
	FileMainGo     = "main.go"
	FileConfigYAML = "config.yaml"
	FileGoMod      = "go.mod"
	FileGitignore  = ".gitignore"
	FileReadme     = "README.md"
	FileRouterGo   = "router.go"
	FileIndexGo    = "index.go"
	FileGitkeep    = ".gitkeep"
)

// 默认值常量
const (
	DefaultProtocol     = "http"
	DefaultModulePrefix = "example.com"
	TotalSteps          = 5
)
