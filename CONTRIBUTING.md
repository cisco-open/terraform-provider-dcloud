# How to Contribute

Thank you for your interest in contributing to the dCloud Terraform Provider project! To ensure smooth collaboration and effective development, please review the following guidelines. By following these guidelines, you demonstrate respect for the time and efforts of the project's contributors. In return, they will reciprocate by addressing your issues, evaluating changes, and assisting with your pull requests. We strive to review incoming issues and pull requests within 10 days. Any lingering issues or pull requests inactive for 60 days will be closed.

Please note that all interactions within the project are subject to our [Code of Conduct](/CODE_OF_CONDUCT.md). This applies to creating issues or pull requests, commenting on them, as well as any real-time discussions in platforms like Slack, Discord, etc.

## Table Of Contents

- [Troubleshooting and Debugging](#troubleshooting-and-debugging)
- [Reporting Issues](#reporting-issues)
- [Development](#development)
    - [Setting up the Development Environment](#setting-up-the-development-environment)
    - [Building the dCloud Terraform Provider](#building-the-dcloud-topology-builder-go-client)
    - [Running Tests](#running-tests)
- [Sending Pull Requests](#sending-pull-requests)
- [Other Ways to Contribute](#other-ways-to-contribute)

## Reporting Issues

Before reporting a new issue, please search our [issues list](https://github.com/cisco-open/dcloud-tb-go-client/issues) to ensure that it hasn't already been reported or fixed.

When creating a new issue, include a clear title, a detailed description, relevant information, and if possible, a test case.

**If you discover a security vulnerability, please refrain from reporting it through GitHub. Instead, follow the security procedures outlined in [SECURITY.md](/SECURITY.md).**

## Development

### Setting up the Development Environment

To set up your development environment for the dCloud Terraform Provider, follow the instructions below:

1. Clone the repository: `git clone https://github.com/cisco-open/terraform-provider-dcloud.git`
2. Navigate to the project directory: `terraform-provider-dcloud`
3. Install the necessary dependencies: https://go.dev/doc/install, https://developer.hashicorp.com/terraform

### Building the dCloud Terraform Provider

To build the dCloud Terraform Provider, execute the following steps:

1. Install Go
2. Execute the following from the root of the repository: `go build .`

### Running Examples

To run the example terraform definitions for the dCloud Terraform Provider, do the following:

1. Navigate to each example directory under `/examples`
2. Execute `terraform init`
3. Create an environment variable to hold your Cisco Authentication Token (use a valid token) `export TB_AUTH_TOKEN=abc123`
4. Execute `terraform apply`

## Sending Pull Requests

Before submitting a new pull request, please check existing pull requests and issues to ensure that your proposed changes or fixes haven't been discussed in the past or already implemented but not released.

When submitting a pull request, include tests for any affected behavior. As we follow semantic versioning, breaking changes may be reserved for the next major version release.

## Other Ways to Contribute

We welcome contributions beyond code. Here are some ways you can contribute to the dCloud Terraform Provider project:

- Help triage and respond to open issues, providing troubleshooting assistance and suggesting fixes.
- Review existing pull requests and test patches against real applications that use the dCloud Terraform Provider.
- Write new tests or add missing test cases to existing tests.

Thank you for your interest in contributing to the dCloud Terraform Provider project!

:heart: