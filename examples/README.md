# Poodle Go SDK Examples

This directory contains example applications demonstrating how to use the Poodle Go SDK.

## Structure

Each example is contained in its own subdirectory with its own `go.mod` file. This structure prevents "main redeclared" errors when running tests or builds from the parent directory.

## Available Examples

### basic_usage/

Shows basic email sending functionality including:

- Sending HTML emails
- Sending plain text emails
- Sending emails with both HTML and text content
- Using the Email model

### error_handling/

Demonstrates comprehensive error handling including:

- Validation errors
- Authentication errors
- Rate limiting errors
- Subscription errors
- Account suspension errors
- Network errors

## Running Examples

To run an example, navigate to its directory and run:

```bash
cd basic_usage
go run .
```

Or build and run:

```bash
cd basic_usage
go build .
./basic_usage
```

## Environment Variables

All examples require the `POODLE_API_KEY` environment variable to be set:

```bash
export POODLE_API_KEY="your_api_key_here"
go run .
```
