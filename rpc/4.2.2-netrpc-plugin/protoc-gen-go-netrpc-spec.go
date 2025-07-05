// 包 main 定义了我们自定义的 protoc 插件。
// 这个可执行文件将被 protoc 编译器调用，以生成基于 .proto 文件的 Go 代码。
// 它的目标是为 Go 标准库的 net/rpc 框架生成完整的服务端和客户端代码。
package main

import (
	"bytes"
	"log"
	"text/template"

	// 注意：这些 'github.com/golang/protobuf' 包路径已被弃用。
	// 在新的项目中，推荐使用 'google.golang.org/protobuf' 模块。
	// 这里为了与原文保持一致，我们仍然使用旧的包。
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
	"google.golang.org/protobuf/proto"
)

// ServiceSpec 定义了用于模板渲染的服务元信息。
type ServiceSpec struct {
	ServiceName string              // 服务的驼峰式命名
	MethodList  []ServiceMethodSpec // 服务包含的方法列表
}

// ServiceMethodSpec 定义了用于模板渲染的服务方法的元信息。
type ServiceMethodSpec struct {
	MethodName     string // 方法的驼峰式命名
	InputTypeName  string // 输入参数的类型名称
	OutputTypeName string // 输出参数的类型名称
}

// tmplService 是用于生成 RPC 服务代码的 Go 模板。
const tmplService = `
{{$root := .}}

// {{.ServiceName}}Interface 是 {{.ServiceName}} 服务的接口。
type {{.ServiceName}}Interface interface {
    {{- range $_, $m := .MethodList}}
    {{$m.MethodName}}(*{{$m.InputTypeName}}, *{{$m.OutputTypeName}}) error
    {{- end}}
}

// Register{{.ServiceName}} 用于将服务注册到 rpc.Server。
func Register{{.ServiceName}}(
    srv *rpc.Server, x {{.ServiceName}}Interface,
) error {
    if err := srv.RegisterName("{{.ServiceName}}", x); err != nil {
        return err
    }
    return nil
}

// {{.ServiceName}}Client 是 {{.ServiceName}} 服务的客户端。
type {{.ServiceName}}Client struct {
    *rpc.Client
}

var _ {{.ServiceName}}Interface = (*{{.ServiceName}}Client)(nil)

// Dial{{.ServiceName}} 用于创建 RPC 客户端。
func Dial{{.ServiceName}}(network, address string) (
    *{{.ServiceName}}Client, error,
) {
    c, err := rpc.Dial(network, address)
    if err != nil {
        return nil, err
    }
    return &{{.ServiceName}}Client{Client: c}, nil
}

{{range $_, $m := .MethodList}}
// {{$m.MethodName}} 是客户端的远程调用方法。
func (p *{{$root.ServiceName}}Client) {{$m.MethodName}}(
    in *{{$m.InputTypeName}}, out *{{$m.OutputTypeName}},
) error {
    return p.Client.Call("{{$root.ServiceName}}.{{$m.MethodName}}", in, out)
}
{{end}}
`

// netrpcPlugin 实现了 generator.Plugin 接口。
type netrpcPlugin struct {
	*generator.Generator
}

// Name 返回插件的名称。
func (p *netrpcPlugin) Name() string {
	return "netrpc-spec"
}

// Init 初始化插件。
func (p *netrpcPlugin) Init(g *generator.Generator) {
	p.Generator = g
}

// GenerateImports 为给定的文件生成导入声明。
func (p *netrpcPlugin) GenerateImports(file *generator.FileDescriptor) {
	if len(file.Service) > 0 {
		p.genImportCode(file)
	}
}

// Generate 为给定的文件生成主要代码。
func (p *netrpcPlugin) Generate(file *generator.FileDescriptor) {
	for _, svc := range file.Service {
		p.genServiceCode(svc)
	}
}

// genImportCode 生成所需的 import 语句。
func (p *netrpcPlugin) genImportCode(file *generator.FileDescriptor) {
	p.P(`import "net/rpc"`)
}

// buildServiceSpec 从服务的描述符中构建 ServiceSpec 元信息。
func (p *netrpcPlugin) buildServiceSpec(
	svc *descriptor.ServiceDescriptorProto,
) *ServiceSpec {
	spec := &ServiceSpec{
		ServiceName: generator.CamelCase(svc.GetName()),
	}

	for _, m := range svc.Method {
		spec.MethodList = append(spec.MethodList, ServiceMethodSpec{
			MethodName:     generator.CamelCase(m.GetName()),
			InputTypeName:  p.TypeName(p.ObjectNamed(m.GetInputType())),
			OutputTypeName: p.TypeName(p.ObjectNamed(m.GetOutputType())),
		})
	}

	return spec
}

// genServiceCode 使用模板为每个服务生成完整的 RPC 代码。
func (p *netrpcPlugin) genServiceCode(svc *descriptor.ServiceDescriptorProto) {
	spec := p.buildServiceSpec(svc)

	var buf bytes.Buffer
	t := template.Must(template.New("").Parse(tmplService))
	err := t.Execute(&buf, spec)
	if err != nil {
		log.Fatal(err)
	}

	p.P(buf.String())
}

// init 函数在 main 函数之前自动执行，用于注册插件。
func init() {
	generator.RegisterPlugin(new(netrpcPlugin))
}

// main 函数是插件的入口点。
func main() {
	g := generator.New()

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		g.Error(err, "reading input")
	}

	if err := proto.Unmarshal(data, g.Request); err != nil {
		g.Error(err, "parsing input proto")
	}

	if len(g.Request.FileToGenerate) == 0 {
		g.Fail("no files to generate")
	}

	g.CommandLineParameters(g.Request.GetParameter())

	g.WrapTypes()
	g.SetPackageNames()
	g.BuildTypeNameMap()

	g.GenerateAllFiles()

	data, err = proto.Marshal(g.Response)
	if err != nil {
		g.Error(err, "failed to marshal output proto")
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		g.Error(err, "failed to write output proto")
	}
}
