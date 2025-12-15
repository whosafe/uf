package cli

import (
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/whosafe/uf/cmd/uf/template"
	"github.com/whosafe/uf/cmd/uf/util"
	"github.com/whosafe/uf/uerror"
)

// InitCommand 初始化命令
type InitCommand struct {
	protocol   string
	modulePath string
	clean      bool
	example    bool
}

// HandleInit 处理 init 命令
func HandleInit(args []string) error {
	cmd := &InitCommand{}

	// 解析命令行参数
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	fs.StringVar(&cmd.protocol, "protocol", DefaultProtocol, "协议类型 (http/tcp/quic)")
	fs.StringVar(&cmd.modulePath, "module", "", "Go 模块路径 (默认: example.com/<project-name>)")
	fs.BoolVar(&cmd.clean, "clean", false, "创建纯净版项目")
	fs.BoolVar(&cmd.example, "example", false, "创建带示例版项目")

	fs.Usage = func() {
		fmt.Println(`用法: uf init [选项] <项目名称>

选项:
  --protocol string    协议类型 (http/tcp/quic) (默认: http)
  --module string      Go 模块路径 (默认: example.com/<项目名称>)
  --clean              创建纯净版项目
  --example            创建带示例版项目 (默认)

示例:
  uf init my-project                           # 创建带示例的 HTTP 项目
  uf init my-project --clean                   # 创建纯净版 HTTP 项目
  uf init my-project --module github.com/user/my-project
`)
	}

	fs.Parse(args)

	// 获取项目名称
	if fs.NArg() == 0 {
		fs.Usage()
		return uerror.New("请指定项目名称")
	}

	projectName := fs.Arg(0)

	// 确保 clean 和 example 互斥，默认创建带示例版
	if cmd.clean {
		cmd.example = false
	} else if !cmd.example {
		cmd.example = true
	}

	// 默认模块路径
	if cmd.modulePath == "" {
		cmd.modulePath = fmt.Sprintf("%s/%s", DefaultModulePrefix, projectName)
	}

	return cmd.Run(projectName)
}

// Run 执行初始化
func (cmd *InitCommand) Run(projectName string) error {
	fmt.Println("🚀 UF 项目初始化")
	fmt.Println()

	// 验证协议
	if cmd.protocol != DefaultProtocol {
		util.Warning(fmt.Sprintf("%s 协议模板开发中，将使用 HTTP 模板", strings.ToUpper(cmd.protocol)))
		cmd.protocol = DefaultProtocol
	}

	// 显示配置信息
	cmd.printConfig(projectName)

	// 创建项目
	return cmd.createProject(projectName)
}

// printConfig 打印配置信息
func (cmd *InitCommand) printConfig(projectName string) {
	util.Info(fmt.Sprintf("项目名称: %s", projectName))
	util.Info(fmt.Sprintf("协议类型: %s", strings.ToUpper(cmd.protocol)))
	if cmd.clean {
		util.Info("项目类型: 纯净版")
	} else {
		util.Info("项目类型: 带示例")
	}
	util.Info(fmt.Sprintf("模块路径: %s", cmd.modulePath))
	fmt.Println()
}

// createProject 创建项目文件和目录结构
// 执行步骤：
// 1. 验证项目目录不存在
// 2. 创建目录结构
// 3. 生成项目文件（main.go, config.yaml等）
// 4. 为空目录添加 .gitkeep
// 5. 初始化 Go 模块
func (cmd *InitCommand) createProject(projectName string) error {
	// 检查目录是否存在
	if util.FileExists(projectName) {
		return uerror.New(fmt.Sprintf("目录 %s 已存在，请使用其他名称或删除现有目录", projectName))
	}

	// 步骤1: 创建目录结构
	util.PrintStep(1, TotalSteps, "创建项目目录...")
	if err := cmd.createDirectories(projectName); err != nil {
		return err
	}

	// 步骤2: 生成项目文件
	util.PrintStep(2, TotalSteps, "生成项目文件...")
	if err := cmd.generateFiles(projectName); err != nil {
		return err
	}

	// 步骤3: 初始化 Go 模块
	util.PrintStep(3, TotalSteps, "初始化 Go 模块...")
	if err := util.RunCommandInDir(projectName, "go", "mod", "tidy"); err != nil {
		util.Warning("go mod tidy 失败，请手动执行")
	}

	// 步骤4: 格式化代码
	util.PrintStep(4, TotalSteps, "格式化代码...")
	util.RunCommandInDir(projectName, "go", "fmt", "./...")

	// 步骤5: 完成
	util.PrintStep(TotalSteps, TotalSteps, "项目创建完成!")
	fmt.Println()

	cmd.printSuccess(projectName)
	return nil
}

// createDirectories 创建项目目录结构
func (cmd *InitCommand) createDirectories(projectName string) error {
	dirs := []string{
		projectName,
		filepath.Join(projectName, DirConfig),
		filepath.Join(projectName, DirHack),
		filepath.Join(projectName, DirInternal),
		filepath.Join(projectName, DirInternal, DirConsts),
		filepath.Join(projectName, DirInternal, DirHandler),
		filepath.Join(projectName, DirInternal, DirModel),
		filepath.Join(projectName, DirInternal, DirDAO),
		filepath.Join(projectName, DirInternal, DirRouter),
		filepath.Join(projectName, DirInternal, DirService),
		filepath.Join(projectName, DirInternal, DirLogic),
		filepath.Join(projectName, DirInternal, DirMiddleware),
		filepath.Join(projectName, DirUtility),
	}

	for _, dir := range dirs {
		if err := util.CreateDir(dir); err != nil {
			return uerror.Wrap(err, fmt.Sprintf("创建目录 %s 失败", dir))
		}
	}

	return nil
}

// generateFiles 生成所有项目文件
func (cmd *InitCommand) generateFiles(projectName string) error {
	// 准备模板变量
	vars := cmd.prepareTemplateVars(projectName)

	// 生成 main.go
	if err := cmd.generateMainFile(projectName, vars); err != nil {
		return err
	}

	// 生成示例文件（仅带示例版）
	if cmd.example {
		if err := cmd.generateExampleFiles(projectName, vars); err != nil {
			return err
		}
	}

	// 生成配置和其他文件
	if err := cmd.generateCommonFiles(projectName, vars); err != nil {
		return err
	}

	// 为空目录添加 .gitkeep
	if err := cmd.addGitkeepFiles(projectName); err != nil {
		return err
	}

	return nil
}

// prepareTemplateVars 准备模板变量
func (cmd *InitCommand) prepareTemplateVars(projectName string) map[string]string {
	// 获取 UF 框架路径（相对路径）
	currentDir, _ := util.GetCurrentDir()
	ufPath, _ := filepath.Rel(filepath.Join(currentDir, projectName), filepath.Join(currentDir, "..", ".."))
	if ufPath == "" {
		ufPath = "../.."
	}

	return map[string]string{
		"ProjectName": projectName,
		"ModulePath":  cmd.modulePath,
		"Protocol":    strings.ToUpper(cmd.protocol),
		"UFPath":      ufPath,
	}
}

// generateMainFile 生成 main.go 文件
func (cmd *InitCommand) generateMainFile(projectName string, vars map[string]string) error {
	var templateContent string
	if cmd.clean {
		templateContent = template.HTTPCleanMain
	} else {
		templateContent = template.HTTPExampleMain
	}

	return cmd.generateFile(projectName, FileMainGo, templateContent, vars)
}

// generateExampleFiles 生成示例文件（router、handler 和 middleware）
func (cmd *InitCommand) generateExampleFiles(projectName string, vars map[string]string) error {
	// internal/router/router.go
	routerPath := filepath.Join(DirInternal, DirRouter, FileRouterGo)
	if err := cmd.generateFile(projectName, routerPath, template.RouterTemplate, vars); err != nil {
		return err
	}

	// internal/handler/index.go
	handlerPath := filepath.Join(DirInternal, DirHandler, FileIndexGo)
	if err := cmd.generateFile(projectName, handlerPath, template.HandlerIndexTemplate, vars); err != nil {
		return err
	}

	// internal/middleware/session.go
	middlewarePath := filepath.Join(DirInternal, DirMiddleware, "session.go")
	if err := cmd.generateFile(projectName, middlewarePath, template.MiddlewareSessionTemplate, vars); err != nil {
		return err
	}

	return nil
}

// generateCommonFiles 生成通用文件（config, go.mod, .gitignore, README）
func (cmd *InitCommand) generateCommonFiles(projectName string, vars map[string]string) error {
	files := map[string]string{
		filepath.Join(DirConfig, FileConfigYAML): template.HTTPCleanConfig,
		FileGoMod:                                template.GoModTemplate,
		FileGitignore:                            template.GitignoreTemplate,
		FileReadme:                               template.ReadmeTemplate,
	}

	for relativePath, templateContent := range files {
		if err := cmd.generateFile(projectName, relativePath, templateContent, vars); err != nil {
			return err
		}
	}

	return nil
}

// generateFile 生成单个文件（辅助函数）
func (cmd *InitCommand) generateFile(projectName, relativePath, templateContent string, vars map[string]string) error {
	content := util.RenderTemplate(templateContent, vars)
	fullPath := filepath.Join(projectName, relativePath)

	if err := util.WriteFile(fullPath, content); err != nil {
		return uerror.Wrap(err, fmt.Sprintf("创建 %s 失败", relativePath))
	}

	return nil
}

// addGitkeepFiles 为空目录添加 .gitkeep 文件
func (cmd *InitCommand) addGitkeepFiles(projectName string) error {
	emptyDirs := []string{
		filepath.Join(projectName, DirHack),
		filepath.Join(projectName, DirUtility),
		filepath.Join(projectName, DirInternal, DirConsts),
		filepath.Join(projectName, DirInternal, DirHandler),
		filepath.Join(projectName, DirInternal, DirModel),
		filepath.Join(projectName, DirInternal, DirDAO),
		filepath.Join(projectName, DirInternal, DirRouter),
		filepath.Join(projectName, DirInternal, DirService),
		filepath.Join(projectName, DirInternal, DirLogic),
		filepath.Join(projectName, DirInternal, DirMiddleware),
	}

	for _, dir := range emptyDirs {
		gitkeepPath := filepath.Join(dir, FileGitkeep)
		if err := util.WriteFile(gitkeepPath, ""); err != nil {
			util.Warning(fmt.Sprintf("创建 %s 失败", gitkeepPath))
		}
	}

	return nil
}

// printSuccess 打印成功信息
func (cmd *InitCommand) printSuccess(projectName string) {
	util.Success(fmt.Sprintf("项目 %s 创建成功!", projectName))
	fmt.Println()
	fmt.Println("下一步:")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Println("  go run main.go")
	fmt.Println()
}
