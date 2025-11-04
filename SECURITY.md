# Security Policy

## Supported Versions

We currently support the following versions with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 1.0.0   | :white_check_mark: |
| < 1.0.0 | :x:                |

## Reporting a Vulnerability

We take the security of openrouter-go seriously. If you discover a security vulnerability, please follow these steps:

### How to Report

1. **Do NOT** open a public GitHub issue for security vulnerabilities
2. Report security issues through [GitHub Security Advisories](https://github.com/hra42/openrouter-go/security/advisories/new)
3. Alternatively, you can email the maintainer directly through GitHub

### What to Include

When reporting a vulnerability, please include:

- Description of the vulnerability
- Steps to reproduce the issue
- Potential impact
- Suggested fix (if any)
- Your contact information for follow-up

### Response Timeline

- **Initial Response**: We aim to acknowledge your report within 48 hours
- **Status Updates**: You can expect updates on the progress every 5-7 days
- **Resolution**: We will work to resolve confirmed vulnerabilities as quickly as possible, typically within 30 days for moderate to high severity issues

### What to Expect

- **Accepted Vulnerabilities**: If your report is accepted, we will:
  - Work on a fix in a private branch
  - Credit you in the security advisory (unless you prefer to remain anonymous)
  - Notify you when the fix is released
  - Publish a security advisory with details after the fix is deployed

- **Declined Reports**: If we determine the report is not a security vulnerability, we will:
  - Explain our reasoning
  - Suggest alternative channels if the issue is still valid but not security-related

## Security Best Practices

When using this library:

1. **API Key Security**: Never commit your OpenRouter API keys to version control
2. **Environment Variables**: Store API keys in environment variables or secure secret management systems
3. **Dependencies**: Keep the library updated to the latest version for security patches
4. **Input Validation**: Always validate and sanitize user input before passing to API calls
5. **Error Handling**: Properly handle errors to avoid leaking sensitive information

## Disclosure Policy

We follow responsible disclosure practices:

- Security issues are fixed privately before public disclosure
- We coordinate with reporters on disclosure timing
- Security advisories are published after fixes are available
- We credit security researchers who report issues (with permission)

Thank you for helping keep openrouter-go secure!
