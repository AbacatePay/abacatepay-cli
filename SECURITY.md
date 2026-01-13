<h1 align="center">Security Policy — AbacatePay CLI</h1>

<p align="center">At AbacatePay, security is a top priority.  
This document describes how to report security vulnerabilities related to the <em>AbacatePay CLI</em>, as well as best practices for handling credentials and sensitive data when using the tool.

The AbacatePay CLI interacts directly with user environments, local servers, and AbacatePay APIs. Responsible disclosure and secure usage are essential to keep the ecosystem safe.</p>

---

<h2 align="center">Reporting Security Vulnerabilities</h2>

<p align="center">If you discover a security vulnerability in the <em>AbacatePay CLI</em>, please report it <em>privately</em>.</p>

📧 Email: security@abacatepay.com  
🔐 Alternative: Use <em>GitHub Security Advisories</em> in the official repository.

<p align="center">When reporting, include as much detail as possible:</p>

- Clear description of the vulnerability
- Steps to reproduce
- Affected versions (if known)
- Potential impact (e.g. token exposure, RCE, privilege escalation)
- Suggested mitigation or fix (optional, but appreciated)

> **Do not open public issues for security vulnerabilities.**

---

<h2 align="center">What to Expect From Us</h2>

- Acknowledgement of your report within **48 business hours**
- Triage and severity assessment
- Fix development based on the criticality of the issue
- Coordinated and responsible disclosure once a fix is available

<p align="center">We aim to act quickly and transparently while protecting users.</p>

---

<h2 align="center">Responsible Disclosure</h2>

<p align="center">We ask that you <em>do not publicly disclose</em> vulnerabilities before AbacatePay has had the opportunity to investigate and release a fix.</br>

We strongly support <em>responsible disclosure</em> and value collaboration with the security community.</p>

---

<h2 align="center">Authentication & Token Security</h2>

<p align="center">The AbacatePay CLI uses <em>OAuth2 Device Flow</em> for authentication.</p>

<h3 align="center">Token Storage</h3>

<p align="center">Authentication tokens are stored securely using the operating system’s native keyring:</p>

- **macOS:** Keychain
- **Linux:** gnome-keyring or kwallet
- **Windows:** Credential Manager

<p align="center">Tokens are <em>never stored in plain text files</em> by default.</p>

<h3 align="center">Recommendations</h3>

- Never share screenshots or logs containing tokens
- Avoid running the CLI on shared or untrusted machines
- Always log out (`abacatepay logout`) on compromised environments
- Keep your system keyring properly configured and locked

---

<h2 align="center">Logs & Sensitive Data</h2>

<p align="center">The CLI generates logs in `~/.abacatepay/logs/`.</br>Logs may include:</p>

- Request metadata
- Event identifiers
- Timing and status information

<h3 align="center">Important Notes</h3>

- Tokens and secrets are <em>not intentionally logged</em>
- Webhook payloads may contain sensitive business data
- Treat log files as sensitive information

<p align="center">Do not commit logs to repositories or share them publicly.</p>

---

<h2 align="center">Webhook Forwarding Security</h2>

<p align="center">When using webhook forwarding:</p>

- Ensure your local server is trusted
- Avoid exposing forwarded endpoints to the public internet
- Use firewalls or local-only bindings when possible
- Validate incoming webhook payloads on your server

<p align="center">The CLI acts as a transport layer — <strong>your application is responsible for payload validation</strong>.</p>

---

<h2 align="center">Binary Integrity & Installation</h2>

<p align="center">We recommend installing the CLI using official channels only:</p>

- `go install github.com/AbacatePay/abacatepay-cli@latest`
- Official Homebrew tap (when available)

<p align="center">Avoid running binaries from unknown sources or forks.</p>

---

<h2 align="center">Scope</h2>

<p align="center">This security policy applies to:</p>

- The AbacatePay CLI source code
- Distributed binaries
- Authentication, token handling, and webhook forwarding logic

<p align="center">For API-level or platform vulnerabilities, refer to the main AbacatePay security policies.</p>

---

<h2 align="center">Acknowledgements</h2>

<p align="center">We appreciate and recognize all responsible disclosures that help improve the security of the AbacatePay ecosystem.</br>Your contributions help keep our users and developers safe.</p>

---

<h2 align="center">References</h2>

- AbacatePay Documentation: https://docs.abacatepay.com
- CLI Documentation: https://docs.abacatepay.com/pages/cli
- Main Security Contact: security@abacatepay.com
