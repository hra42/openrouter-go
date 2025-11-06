# PDF Testing Documentation

This document describes the PDF testing capabilities added to the OpenRouter Go client test suite.

## Test PDF File

**Location**: `cmd/openrouter-test/test-pdf.pdf`
**Size**: 65 KB
**Public URL**: https://hra42.com/test-pdf.pdf

The same PDF is available both as a local file and via URL, allowing comprehensive testing of both PDF input methods.

## PDF Test Cases

### 1. PDF Input (URL) - `RunPDFURLTest`
Tests basic PDF input using a public URL.

```bash
go run cmd/openrouter-test/main.go -test pdf -model anthropic/claude-sonnet-4
```

**What it tests:**
- Sending PDF via public URL
- Basic PDF processing
- Response generation from PDF content

### 2. PDF with Parser Engine - `RunPDFWithEngineTest`
Tests PDF parsing with specific engine configuration.

```bash
go run cmd/openrouter-test/main.go -test pdfengine -model google/gemma-3-27b-it
```

**What it tests:**
- Plugin configuration for PDF parsing
- Using the `pdf-text` engine (free tier)
- Engine-specific processing

**Available engines:**
- `pdf-text` - Free, for digital PDFs with text
- `mistral-ocr` - $0.0004/1K pages, for scanned PDFs
- `native` - Uses model's native file support

### 3. PDF with Annotations - `RunPDFWithAnnotationsTest`
Tests reusing file annotations to avoid re-parsing costs.

```bash
go run cmd/openrouter-test/main.go -test pdfannotations -model google/gemma-3-27b-it
```

**What it tests:**
- Initial PDF processing with annotations
- Follow-up requests reusing annotations
- Cost optimization through annotation reuse

### 4. Multiple Files - `RunMultipleFilesTest`
Tests sending multiple files in a single request.

```bash
go run cmd/openrouter-test/main.go -test multiplefiles -model anthropic/claude-sonnet-4
```

**What it tests:**
- Multiple file attachments (PDF + image)
- Mixed content types in one request
- Combined analysis of different file types

### 5. ContentBuilder with PDF - `RunPDFContentBuilderTest`
Tests using ContentBuilder for complex messages with PDFs.

```bash
go run cmd/openrouter-test/main.go -test pdfbuilder -model anthropic/claude-sonnet-4
```

**What it tests:**
- ContentBuilder API with PDFs
- Interleaving text and files
- Complex message structure

### 6. Base64 Encoded PDF - `RunBase64PDFTest` ⭐ NEW
Tests PDF input using base64-encoded local file.

```bash
go run cmd/openrouter-test/main.go -test base64pdf -model anthropic/claude-sonnet-4
```

**What it tests:**
- Reading local PDF file
- Base64 encoding
- Sending encoded PDF to API
- Same functionality as URL method but for local files

**Implementation details:**
- Uses `test-pdf.pdf` from the test directory
- Automatically encodes file to base64 data URL
- Validates response from encoded PDF

### 7. PDF URL vs Base64 Comparison - `RunPDFComparisonTest` ⭐ NEW
Tests comparing the same PDF sent via URL and base64 encoding.

```bash
go run cmd/openrouter-test/main.go -test pdfcomparison -model anthropic/claude-sonnet-4
```

**What it tests:**
- PDF sent via URL (https://hra42.com/test-pdf.pdf)
- Same PDF sent via base64 encoding (local file)
- Response consistency between methods
- Performance comparison (timing)
- Token usage comparison

**What to expect:**
- Both methods should produce similar responses
- Base64 may have slightly higher prompt tokens (due to encoding overhead)
- URL method may be faster (no encoding step)
- Functionally equivalent results

## Running All PDF Tests

Run all PDF tests at once:

```bash
go run cmd/openrouter-test/main.go -test all -model anthropic/claude-sonnet-4
```

This will execute all PDF tests as part of the complete test suite.

## Test Recommendations

### For URL Testing
Best for:
- Testing with public documents
- Faster execution (no encoding needed)
- Lower bandwidth for large files

**Recommended models:**
- `anthropic/claude-sonnet-4` - Native file support
- `google/gemini-2.0-flash-thinking-exp` - Good vision/document support
- `openai/gpt-4o` - Excellent document understanding

### For Base64 Testing
Best for:
- Testing with local/private files
- Validating encoding logic
- Testing without internet access

**Recommended models:**
- Same as URL testing (functionality is identical)

### Cost Considerations

1. **Free Tier** (`pdf-text` engine):
   - Use for well-structured digital PDFs
   - No additional cost beyond model tokens
   - Good for most testing scenarios

2. **OCR Engine** (`mistral-ocr`):
   - $0.0004 per 1,000 pages
   - Better for scanned documents
   - May not be necessary for test PDF

3. **Native Support**:
   - Charged as input tokens
   - Cost varies by model
   - Often highest quality

## Unit Tests

The library also includes comprehensive unit tests for PDF functionality:

```bash
# Run PDF-specific unit tests
go test -v -run TestEncodePDF
go test -v -run TestCreateUserMessageWithPDF
go test -v -run TestCreateFileParserPlugin
go test -v -run TestContentBuilderWithPDF

# Run all tests
go test ./...
```

## Troubleshooting

### Test Fails with "Model not available"
Some models may not support file input. Try:
- `anthropic/claude-sonnet-4`
- `google/gemini-2.0-flash-thinking-exp`
- `openai/gpt-4o`

### Base64 Test Fails with "Failed to encode PDF"
Check that `test-pdf.pdf` exists in `cmd/openrouter-test/` directory.

### Comparison Test Shows Different Results
This is normal - LLM responses can vary even with identical input. Look for:
- Similar topics/themes in responses
- Comparable token counts
- Both tests completing successfully

## Example Output

```
🔄 Test: Base64 Encoded PDF
✅ Success! (3.45s)
   Response: This document discusses...
   Model: anthropic/claude-3.5-sonnet
   PDF: Base64-encoded local file
   Tokens: 1245 prompt, 98 completion, 1343 total

🔄 Test: PDF URL vs Base64 Comparison
   Testing URL version...
   Testing Base64 version...
✅ Success! URL: 2.87s, Base64: 3.12s

   URL Response: The document covers...
   URL Tokens: 1234 prompt, 95 completion

   Base64 Response: The document details...
   Base64 Tokens: 1245 prompt, 97 completion

   ✓ Both methods successfully processed the same PDF
```

## Integration with CI/CD

These tests can be integrated into CI/CD pipelines:

```yaml
# Example GitHub Actions workflow
- name: Run PDF Tests
  run: |
    export OPENROUTER_API_KEY=${{ secrets.OPENROUTER_API_KEY }}
    go run cmd/openrouter-test/main.go -test pdf
    go run cmd/openrouter-test/main.go -test base64pdf
    go run cmd/openrouter-test/main.go -test pdfcomparison
```

## Future Enhancements

Potential additions to the PDF test suite:
- Multi-page PDF testing
- Large PDF handling (>10MB)
- PDF with complex formatting
- Scanned document testing with OCR engine
- Performance benchmarking
