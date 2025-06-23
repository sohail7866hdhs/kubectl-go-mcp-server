# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of kubectl-go-mcp-server seriously. If you believe you've found a security vulnerability, please follow the steps below:

### For Public Vulnerabilities

- Create a regular issue if the vulnerability is already public knowledge or poses minimal risk

### For Security-Sensitive Vulnerabilities

1. **Do not** disclose the vulnerability publicly
2. Email the details to joetech.ooj@gmail.com
3. Include as much information as possible:
   - A description of the vulnerability
   - How it could be exploited
   - Steps to reproduce
   - Potential impact
   - Suggested fixes if you have them

### What to Expect

- We will acknowledge receipt of your report within 48 hours
- We aim to provide an initial assessment within 7 days
- We will keep you updated on our progress
- Once the issue is resolved, we will credit you for the discovery (unless you prefer to remain anonymous)

## Security Best Practices

When using kubectl-go-mcp-server, we recommend following these security best practices:

1. **Keep your version updated**: Always use the latest version with security patches
2. **Use a dedicated kubeconfig**: Create a restricted kubeconfig with limited permissions
3. **Configure proper kubectl restrictions**: Use the tool's built-in safety features
4. **Review command sanitization**: Check the sanitization patterns in the configuration
5. **Audit trail**: Enable command logging in a production environment
