# PDF inputs

PDFs are sent as base64 data with an optional parsing-engine hint (text vs. OCR). File annotations from a first request can be reused on follow-ups to avoid re-parsing cost.

See [`examples/pdf-inputs/main.go`](../../examples/pdf-inputs/main.go) for:

- Basic PDF upload
- Choosing an engine with `WithPDFEngine`
- Reusing annotations across requests
- Base64 vs. URL inputs
- Engine comparison

Models must support document inputs.
