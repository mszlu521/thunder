# EPUB Document Parser

EPUB解析器是为Eino框架设计的文档解析器，用于解析EPUB格式的电子书文件。

## 功能特点

- 解析EPUB格式的电子书文件
- 按照阅读顺序提取章节内容
- 支持HTML到纯文本的转换
- 提取章节元数据信息（标题、顺序等）
- 与Eino框架无缝集成

## 安装

EPUB解析器已经集成在项目中，确保你已经安装了所需的依赖：

```bash
go get github.com/jaytaylor/html2text
```

## 使用方法

### 1. 在文件加载器中使用

EPUB解析器需要通过文件路径来解析EPUB文件，因此推荐与文件加载器一起使用：

```go
import (
    "context"
    "github.com/cloudwego/eino-ext/components/document/loader/file"
    "your-project/thunder/einos/components/document/parser/epub"
)

// 创建EPUB解析器（默认配置）
epubParser, err := epub.NewParser(context.Background(), nil)
if err != nil {
    // 处理错误
}

// 创建EPUB解析器（带配置选项）
epubParser, err := epub.NewParser(context.Background(), &epub.Config{
    StripHTML: true, // 将HTML内容转换为纯文本
}, epub.WithStripHTML(true))
if err != nil {
    // 处理错误
}

// 创建文件加载器并指定EPUB解析器
loader, err := file.NewFileLoader(context.Background(), &file.FileLoaderConfig{
    Parser: epubParser,
})
if err != nil {
    // 处理错误
}

// 加载EPUB文件
docs, err := loader.Load(context.Background(), document.Source{
    URI: "/path/to/your/book.epub",
})
```

### 2. 直接使用解析器

你也可以直接使用EPUB解析器来解析EPUB文件：

```go
import (
    "context"
    "your-project/thunder/einos/components/document/parser/epub"
)

// 创建EPUB解析器
parser, err := epub.NewParser(context.Background(), nil)
if err != nil {
    // 处理错误
}

// 解析EPUB文件
docs, err := parser.ParseFromPath(context.Background(), "/path/to/your/book.epub")
if err != nil {
    // 处理错误
}

// 处理解析后的文档
for _, doc := range docs {
    fmt.Println("Content:", doc.Content)
    fmt.Println("Metadata:", doc.MetaData)
}
```

### 3. 配置选项

EPUB解析器支持以下配置选项：

- `StripHTML`: 是否将HTML内容转换为纯文本（默认：false）

配置选项可以通过两种方式设置：

1. 通过Config结构体：
```go
config := &epub.Config{
    StripHTML: true,
}
parser, err := epub.NewParser(context.Background(), config)
```

2. 通过函数选项：
```go
parser, err := epub.NewParser(context.Background(), nil, epub.WithStripHTML(true))
```

3. 通过解析时传递额外元数据：
```go
docs, err := parser.ParseFromPath(context.Background(), "/path/to/your/book.epub", 
    parser.WithExtraMeta(map[string]any{
        "strip_html": true,
    }))
```

## 返回的数据结构

解析器会返回一个文档数组，每个文档代表EPUB中的一个章节，包含：

- `Content`: 章节的内容文本
- `MetaData`: 元数据信息，包括：
  - `order`: 章节在书中的顺序
  - `href`: 章节文件的相对路径
  - `media_type`: 文件的媒体类型
  - `source_uri`: 源文件的URI
  - `epub_opf`: OPF文件的路径
  - `id_ref`: 章节的ID引用
  - `title`: 章节标题（如果能从HTML中提取）

## 工作原理

EPUB解析器按照以下步骤解析EPUB文件：

1. 将EPUB文件作为ZIP归档文件打开
2. 读取`META-INF/container.xml`文件找到OPF文件路径
3. 解析OPF文件获取内容清单(manifest)和阅读顺序(spine)
4. 按照阅读顺序解析各个章节文件
5. 提取章节内容和元数据

## 限制

- 不支持加密的EPUB文件
- 不支持从通用reader解析（需要文件路径）
- 只解析XHTML和HTML格式的章节内容
- 不处理EPUB中的CSS样式和图像

## 扩展

你可以通过实现以下方法来扩展EPUB解析器的功能：

1. 添加对其他内容类型的支持
2. 增加对EPUB 2和EPUB 3的更多特性的支持
3. 添加对加密EPUB文件的支持
4. 增加对元数据的更详细解析