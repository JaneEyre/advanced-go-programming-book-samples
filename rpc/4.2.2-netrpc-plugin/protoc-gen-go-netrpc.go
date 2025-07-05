// 包 main 定义了我们自定义的 protoc 插件。
// 这个可执行文件将被 protoc 编译器调用，以生成基于 .proto 文件的 Go 代码。
// 我们将其命名为 protoc-gen-go-netrpc，以区别于官方的 protoc-gen-go。
package main

import (
	// 注意：这些 'github.com/golang/protobuf' 包路径已被弃用。
	// 在新的项目中，推荐使用 'google.golang.org/protobuf' 模块。
	// 这里为了与原文保持一致，我们仍然使用旧的包。
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
	"google.golang.org/protobuf/proto"
)

// netrpcPlugin 实现了 generator.Plugin 接口，用于生成适用于标准库 net/rpc 的代码。
type netrpcPlugin struct {
	*generator.Generator // 通过内嵌 *generator.Generator，继承了代码生成所需的辅助方法，如 P(), In(), Out()。
}

// Name 返回插件的名称。这个名称用于在 protoc 命令的参数中指定插件，例如：plugins=netrpc。
func (p *netrpcPlugin) Name() string {
	//return "netrpc"
	return "abc"
}

// Init 在代码生成开始前被调用，用于初始化插件。
// 参数 g 是一个代码生成器实例，包含了所有解析后的 .proto 文件信息。
// 我们将 g 赋值给内嵌的 Generator，这样插件实例就可以直接调用 g 的方法。
func (p *netrpcPlugin) Init(g *generator.Generator) {
	p.Generator = g
}

// GenerateImports 为给定的文件生成导入声明。
// 它在 Generate 方法之后被调用。
// 我们检查文件中是否定义了任何 service，如果有，则调用 genImportCode 生成导入代码。
func (p *netrpcPlugin) GenerateImports(file *generator.FileDescriptor) {
	if len(file.Service) > 0 {
		p.genImportCode(file)
	}
}

// Generate 为给定的文件生成主要的 Go 代码。
// 它会遍历文件中的所有 service 定义，并为每个 service 调用 genServiceCode 生成相应的代码。
func (p *netrpcPlugin) Generate(file *generator.FileDescriptor) {
	for _, svc := range file.Service {
		p.genServiceCode(svc)
	}
}

// genImportCode 是一个辅助方法，用于生成导入语句。
// 在这个示例中，我们只打印一个占位符注释。
// 在实际应用中，这里会根据需要生成如 "net/rpc" 等包的导入语句。
func (p *netrpcPlugin) genImportCode(file *generator.FileDescriptor) {
	p.P("// TODO: import code for net/rpc")
}

// genServiceCode 是一个辅助方法，用于为具体的服务生成代码。
// svc 参数包含了服务的完整描述信息，如服务名、方法等。
// 在这个示例中，我们只打印一个包含服务名的占位符注释。
func (p *netrpcPlugin) genServiceCode(svc *descriptor.ServiceDescriptorProto) {
	p.P("// TODO: service code, Name = " + svc.GetName())
}

// init 函数在 main 函数之前自动执行。
// 我们在这里通过 generator.RegisterPlugin 将我们的 netrpcPlugin 注册到 protoc-gen-go 的插件系统中。
// 这样，当 main 函数中的 generator 运行时，它就能识别并调用我们的插件。
func init() {
	generator.RegisterPlugin(new(netrpcPlugin))
}

// main 函数是 protoc-gen-go-netrpc 可执行文件的入口点。
// 它的逻辑基本克隆自官方的 protoc-gen-go，负责处理与 protoc 编译器的通信。
func main() {
	// 1. 创建一个新的代码生成器实例。
	g := generator.New()

	// 2. 从标准输入（stdin）读取 protoc 编译器传递过来的 CodeGeneratorRequest 数据。
	// 这是 protoc 与插件通信的方式。
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		g.Error(err, "reading input")
	}

	// 3. 使用 proto.Unmarshal 将读取的数据反序列化到生成器的 Request 字段中。
	if err := proto.Unmarshal(data, g.Request); err != nil {
		g.Error(err, "parsing input proto")
	}

	// 4. 检查是否有文件需要生成。
	if len(g.Request.FileToGenerate) == 0 {
		g.Fail("no files to generate")
	}

	// 5. 解析命令行参数（例如 "plugins=netrpc"）。
	g.CommandLineParameters(g.Request.GetParameter())

	// 6. 构建内部所需的数据结构，如类型映射等。
	g.WrapTypes()
	g.SetPackageNames()
	g.BuildTypeNameMap()

	// 7. 调用 GenerateAllFiles，这个方法会遍历所有注册的插件（包括我们自定义的 netrpcPlugin），
	// 并调用它们的 Generate 和 GenerateImports 方法来生成代码。
	g.GenerateAllFiles()

	// 8. 将生成的结果（g.Response）序列化为 protobuf 格式。
	data, err = proto.Marshal(g.Response)
	if err != nil {
		g.Error(err, "failed to marshal output proto")
	}

	// 9. 将序列化后的数据写入标准输出（stdout），返回给 protoc 编译器。
	// protoc 会负责将这些内容写入最终的 .pb.go 文件。
	_, err = os.Stdout.Write(data)
	if err != nil {
		g.Error(err, "failed to write output proto")
	}
}
