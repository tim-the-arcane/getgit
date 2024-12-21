# getgit

`getgit` is a command-line tool written in Go that allows you to download specific files or entire folders from GitHub repositories directly into your current working directory. Whether you need a single file or a complete folder with its contents, `getgit` simplifies the process without the need to clone the entire repository.

## Features

- **Download Single Files:** Easily download individual files from any public GitHub repository.
- **Download Folders:** Recursively download entire folders along with all their subdirectories and files.
- **Simple Command-Line Interface:** Intuitive commands to get you started quickly.
- **Lightweight:** Minimal dependencies and efficient performance.

## Installation

### Prerequisites

- **Go:** Ensure you have Go installed on your system. You can download it from [here](https://golang.org/dl/).

### Clone the Repository

```bash
git clone https://github.com/tim-the-arcane/getgit.git
cd getgit
```

### Build the Executable

```bash
go build -o getgit main.go
```

This will generate an executable named `getgit` (or `getgit.exe` on Windows) in the current directory.

## Usage

### Download a Single File

To download a specific file from a GitHub repository:

```bash
./getgit https://github.com/username/repo/blob/branch/path/to/file.ext
```

**Example:**

```bash
./getgit https://github.com/facebook/lexical/blob/main/packages/lexical-website/docs/index.md
```

This command will download `index.md` directly into your current working directory.

### Download a Folder

To download an entire folder along with its contents:

```bash
./getgit https://github.com/username/repo/tree/branch/path/to/folder
```

**Example:**

```bash
./getgit https://github.com/facebook/lexical/tree/main/packages/lexical-website/docs
```

This will create a `docs` folder in your current directory and populate it with all files and subdirectories from the specified GitHub folder.

## Compilation as a Binary (Optional)

If you prefer to use `getgit` without rebuilding it every time, you can compile it into a binary:

```bash
go build -o getgit main.go
```

You can then move the binary to a directory in your `PATH` for easy access:

```bash
mv getgit /usr/local/bin/
```

Now, you can use `getgit` from anywhere in your terminal:

```bash
getgit <GitHub-URL>
```

## Contributing

Contributions are welcome! If you'd like to contribute to `getgit`, please follow these steps:

1. **Fork the Repository:** Click the "Fork" button at the top of the repository page.
2. **Clone Your Fork:**

   ```bash
   git clone https://github.com/your-username/getgit.git
   cd getgit
   ```

3. **Create a New Branch:**

   ```bash
   git checkout -b feature/YourFeatureName
   ```

4. **Make Your Changes:** Implement your feature or bug fix.
5. **Commit Your Changes:**

   ```bash
   git commit -m "Add Your Feature Description"
   ```

6. **Push to Your Fork:**

   ```bash
   git push origin feature/YourFeatureName
   ```

7. **Create a Pull Request:** Go to the original repository and create a pull request from your fork.

## License

This project is licensed under the [MIT License](LICENSE.md).

## Acknowledgements

- [GitHub API](https://docs.github.com/en/rest)
- [Go Programming Language](https://golang.org/)

## Contact

For any questions or suggestions, feel free to open an issue or contact me at [your-email@example.com](mailto:your-email@example.com).
