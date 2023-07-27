# Datax Go Client

This project provides an executable program to replace `datax.py`. suitable for launching Datax in environments lacking a Python setup. you can find binarys in [release page](https://github.com/rainbowatcher/datax-go-client/release)

一个可执行的单文件，用来替代`datax.py`，适合在没有 python 环境的情况下启动 `Datax`，release页面提供[文件下载地址](https://github.com/rainbowatcher/datax-go-client/release)

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

## License

This project is licensed under the [MIT license](LICENSE).
