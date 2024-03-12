
# gptunzip 📦

`gptunzip` is a command-line utility designed to transform a zipped Git repository into a structured text format, ideal for feeding into ChatGPT 🤖 and other NLP models. This tool simplifies the conversion process, enabling efficient processing and analysis of repository contents.

## Installation 🛠️

Before installing `gptunzip`, ensure you have the Go programming language 🐹 installed on your system. If not, you can download it from [the official Go website](https://golang.org/dl/).

Install `gptunzip` by running:

```bash
go install github.com/ginkida/gptunzip@latest
```

This command installs the `gptunzip` binary into your `$GOPATH/bin` directory. Ensure `$GOPATH/bin` is in your `$PATH` for easy access to the `gptunzip` command.

## Usage 🚀

To use `gptunzip`, execute the following command:

```bash
gptunzip [flags] /path/to/zipped/git/repository.zip
```

### Flags 🚩

`gptunzip` supports the creation of smaller, manageable parts for ChatGPT with a single flag:

* `-p`, `--parts`: When enabled, the utility creates smaller parts of the chatGPT text, making it easier to manage and process. This is particularly useful for large repositories.

### Simplified File Handling 🗂️

`gptunzip` streamlines the process by focusing on the essential task of converting and structuring repository data. There's no need for a `.gptignore` file or handling `.gitignore` files, as the tool automatically processes the entire zipped repository.

## Contributing 🤝

Your contributions are welcome! Feel free to submit pull requests or open issues on the GitHub repository to improve `gptunzip` or suggest new features.

## License 📄

`gptunzip` is made available under the MIT License, offering flexibility and freedom for personal and commercial use. Check out the [LICENSE](LICENSE) file for more information.
