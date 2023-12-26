# ProxyChecker

ProxyChecker is a blazing fast proxy checker tool developed in Golang. This tool is designed to efficiently validate a list of proxies. It supports both SOCKS4 and SOCKS5 proxies.

## Features

- Concurrent Processing: Utilizes goroutines for fast and efficient proxy checking.
- Bare-bones: Messes with raw TCP connections for maximum efficiency.
- SOCKS4 and SOCKS5 Support: Checks both SOCKS4 and SOCKS5 proxies.

## Installation

To install ProxyChecker, ensure you have Golang installed on your system. You can then clone this repository and build the binary:

```bash
git clone https://github.com/moohr/proxyChecker.git
cd proxyChecker
go build
```
## Usage

Run ProxyChecker with the following command-line options:

```bash
./proxyChecker [options]
```
### Options
- -dt: Dial timeout duration (default 5s).
- -if: File with list of proxies (default "list.txt").
- -max: Maximum number of goroutines (default 5000).
- -of: Output file to write results to (defaults to stdout).
- -rt: Read timeout duration (default 10s).
- -s4: Check for SOCKS4 proxies (default true).
- -s5: Check for SOCKS5 proxies (default true).

### Example
Check proxies from proxies.txt, write results to results.txt, with a dial timeout of 3s and a maximum of 3000 goroutines:

```bash
./proxyChecker -if proxies.txt -of results.txt -dt 3s -max 3000
```

## Contributing

Contributions to ProxyChecker are welcome! Please submit your pull requests or issues to the repository.

## License

This project is licensed under the MIT License.

## Disclaimer

ProxyChecker is developed for educational and legitimate testing purposes only. The developers are not responsible for any misuse or damage caused by this tool.

This tool is designed as a companion tool to [ipBlacklist](https://github.com/moohr/ipBlacklist) to identify legitimate, alive proxy servers that poses a threat rather than ones that are randomly added to the list of harmful proxies.

Enjoy using ProxyChecker for your proxy checking needs! For any questions or feedback, please open an issue in the repository.