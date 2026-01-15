# Text File Inputs Example

This example demonstrates how to send text file content to models using the OpenRouter API.

## Features

- **Single file input**: Send code, configuration, or documentation files
- **Multiple files**: Compare or analyze multiple files in one request
- **ContentBuilder**: Build complex messages with text files and other content
- **Format support**: .txt, .md, .json, .csv, .js, .py, .go, .java, .rs, .ts, .yaml, .toml, .xml, .ini, and more

## Usage

```bash
export OPENROUTER_API_KEY="your-api-key"
go run main.go
```

## Supported File Formats

### Common Text Formats
- `.txt` - Plain text
- `.md` - Markdown
- `.json` - JSON files
- `.csv` - CSV files

### Code Files
- `.js`, `.jsx` - JavaScript
- `.ts`, `.tsx` - TypeScript
- `.py` - Python
- `.go` - Go
- `.java` - Java
- `.rs` - Rust
- `.c`, `.cpp`, `.h` - C/C++
- `.rb` - Ruby
- `.php` - PHP
- `.swift` - Swift
- `.kt` - Kotlin

### Configuration Files
- `.yaml`, `.yml` - YAML
- `.toml` - TOML
- `.xml` - XML
- `.ini` - INI
- `.env` - Environment files

## Notes

- Text files are sent as **inline text content**, not base64-encoded
- Files must contain valid UTF-8 encoded text
- File content is included with a header showing the filename for context
- Maximum file size depends on the model's context window
