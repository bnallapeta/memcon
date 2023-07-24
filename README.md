# MemCon: Memory Consumption Monitor

Welcome to MemCon, a fun and useful tool to monitor the memory usage of your selected applications over time. This project is intended to help you understand how much memory your applications are consuming in real-time. 

*Note: At the moment, this application works only on MacBook*

## Key Features

- **Monitor your applications:** Simply list the applications you want to monitor, and MemCon will keep track of their memory usage.
- **Local storage:** All data is stored in a SQLite database on your local system. This means that your data remains private and secure - no data is sent externally.
- **In-depth Memory Usage Insights:** MemCon provides a means to examine your apps' memory consumption across different time periods, supporting informed decisions about system resource usage.

## Getting Started

### Prerequisites

- Golang
- Root access (some processes require `sudo` permissions to monitor)

1. Download the binary from the [releases page](https://github.com/bnallapeta/memcon/releases) on the GitHub project page.
2. In the terminal, navigate to the directory where you downloaded the binary and start the program:

```bash
cd <path-to-binary>
sudo ./memcon -apps <path-to-apps.env> -interval <fetchInterval>
```

`<path-to-apps.env>` is the path to your environment file. If you don't specify -apps, it will take the default file at config/apps.env. 

`<fetchInterval>` is the time in seconds between successive fetches for memory usage data. If you don't specify -interval, it will take 360 seconds as the default.

### For Developers
1. Clone the repository:

```bash
git clone https://github.com/bnallapeta/memcon.git
cd memcon
```

2. Ensure Golang is installed in your system.
3. Compile the program into a binary:

```bash
GOOS=darwin GOARCH=amd64 go build -o memcon main.go
```

This creates a binary named memcon in your current directory.

4. To start the program, run:

```bash
sudo ./memcon -apps <path-to-apps.env> -interval <fetchInterval>
```
`<path-to-apps.env>` is the path to your environment file. If you don't specify -apps, it will take the default file at config/apps.env. 

`<fetchInterval>` is the time in seconds between successive fetches for memory usage data. If you don't specify -interval, it will take 360 seconds as the default.

### How it Works

MemCon periodically polls the system to fetch memory usage data of the specified apps. This data is stored in a local SQLite database (`process_memory.db`). A web server is started which serves a page at `http://localhost:8088/visualize` where you can visualize the data.

The `-apps` flag is used to specify an environment file that contains a list of the applications to be monitored. The file contains a comma separated list of values which are the app names. See [apps.env](./config/apps.env) for example. The `-interval` flag specifies the time interval in seconds between successive memory usage data fetches.

## Data Privacy

MemCon stores all its data in a local SQLite database on your system. No data is transmitted externally, ensuring that your usage data remains private and secure.

## Contribution

Contributions are welcome! Feel free to fork the project, open issues, and submit PRs. For major changes, please open an issue first to discuss the change.

## Disclaimer

This tool requires root permissions to monitor certain processes. Please ensure that you understand the implications of this and use it responsibly.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.