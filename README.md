# Datax Go Client

This project provides an executable program to replace `datax.py`.

[中文](./README_CN.md)

## Requirements

- [GoLang](https://github.com/golang/go) - The Go programming language
- [Just](https://github.com/casey/just) - A command runner for your workflow

## Usage

### Commands

Below are a list of commands available in this project:

| Command      | Description                                             |
| ------------ | ------------------------------------------------------- |
| `just`       | List all available commands                             |
| `just fmt`   | Formats all go files in the project                     |
| `just build` | Builds the project and outputs the executable to `/out` |
| `just start` | Build and Starts the built executable from `/out`       |
| `just test`  | Runs all tests in the project                           |

#### Aliases

The following aliases have been provided as shortcuts for some of the longer commands above:

- `s` = `start`
- `b` = `build`
- `t` = `test`
- `f` = `fmt`

### Examples

To build the project, run:

```
just build
```

This will compile the code and output an executable file named `datax` in the `/out` directory.

To start the project, simply run:

```
just start
```

This will build and execute the `datax` executable.

To run all tests in the project, run:

```
just test
```

This will run all tests written in Go language.

To format the code, run:

```
just fmt
```

This will apply formatting rules to all Go files in the project.

### Build Windows Executable

If you want to build a Windows executable, use the following command:

```
just build-windows
```

This will compile the code and output an executable file named `datax.exe` in the `/out` directory.

## License

This project is licensed under the [MIT license](LICENSE).
