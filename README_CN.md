# Datax Go 客户端

该项目提供了一个可执行程序，用于替代 `datax.py`。

## 要求

- [GoLang](https://github.com/golang/go) - Go 编程语言
- [Just](https://github.com/casey/just) - 命令运行器，用于您的工作流

## 用法

### 命令

以下是该项目中可用的命令列表：

| 命令         | 描述                                     |
| ------------ | ---------------------------------------- |
| `just`       | 列出所有可用的命令                       |
| `just fmt`   | 格式化项目中的所有 Go 文件               |
| `just build` | 构建项目并将可执行文件输出到 `/out` 目录 |
| `just start` | 构建并运行 `/out` 目录下的可执行文件     |
| `just test`  | 运行项目中的所有测试                     |

#### 别名

以下别名可作为某些较长命令的快捷方式：

- `s` = `start`
- `b` = `build`
- `t` = `test`
- `f` = `fmt`

### 示例

要构建该项目，请键入以下命令：

```
just build
```

这将编译代码并在 `/out` 目录中输出名为 `datax` 的可执行文件。

要启动项目，简单地键入以下命令：

```
just start
```

这将构建并执行 `datax` 可执行文件。

要运行项目中的所有测试，请键入以下命令：

```
just test
```

这将运行 Go 语言编写的所有测试。

要格式化代码，请键入以下命令：

```
just fmt
```

这将对项目中的所有 Go 文件应用格式化规则。

### 构建 Windows 可执行文件

如果您想构建 Windows 可执行文件，请使用以下命令：

```
just build-windows
```

这将编译代码并在 `/out` 目录中输出名为 `datax.exe` 的可执行文件。

## 许可证

该项目使用 [MIT 许可证](LICENSE)。
