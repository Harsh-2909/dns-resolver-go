# DNS Resolver

![GitHub](https://img.shields.io/github/license/Harsh-2909/dns-resolver-go)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Harsh-2909/dns-resolver-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/Harsh-2909/dns-resolver-go)](https://goreportcard.com/report/github.com/Harsh-2909/dns-resolver-go)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/Harsh-2909/dns-resolver-go/go.yml?branch=main)

![DNS Resolver](assets/dns-resolver.png)

DNS Resolver is a lightweight DNS resolver implemented in Go, designed for simplicity and efficiency.

## Features

-   **DNS Query Resolution:** Resolves DNS queries using UDP.
-   **IPv4 Support:** Capable of resolving IPv4 addresses.
-   **Timeout Handling:** Includes timeout handling for queries to prevent blocking.
-   **Caching:** Implements a caching mechanism to improve query response times.

## Getting Started

### Prerequisites

-   Go 1.18 or higher installed.

### Building

Clone the repository and navigate to the project directory:

```bash
git clone https://github.com/Harsh-2909/dns-resolver-go
cd dns-resolver-go
```

Build the project using Go:

```bash
go build
```

### Running

To run the DNS resolver, use the following command:

```bash
./dns-resolver <domain>
```

If you want to disable caching, use the `--no-cache` flag:

```bash
./dns-resolver <domain> --no-cache
```

### Testing

Unit tests are included to verify the functionality of the resolver. Run the tests with:

```bash
go test ./...
```

## Features to be added

-   **IPv6 Support:** Add support for IPv6 addresses.
-   **DNS Query Resolution:** Add support for resolving DNS queries using TCP or HTTPS.

## Usage

The DNS resolver takes in a domain name as a command-line argument. It then resolves the domain name using UDP and prints the IP addresses associated with the domain name.

The resolver sends a DNS query to the root servers at the start and then sends subsequent queries to the closest parent domain servers. The resolver stops sending queries when it receives an answer response from the nameservers.

To know more about how DNS works, I have written a blog post on the topic. You can find it [here](https://harshagarwal29.hashnode.dev/unveiling-the-magic-of-dns-how-the-internets-directory-works).

## Contributing

If you'd like to contribute to DNS Resolver, please fork the repository and submit a pull request. Feel free to open issues for bug reports, feature requests, or general feedback.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
