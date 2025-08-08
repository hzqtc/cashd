# cashd

`cashd` is a fast and cozy interactive TUI for personal finance management.
It allows you to effortlessly track, analyze, and gain insights into your financial transactions directly from your terminal.
`cashd` currently supports ledger/hledger and CSV as data sources

<p float="left">
  <img src="https://raw.github.com/hzqtc/cashd/master/screenshots/transaction_view.png" width="400" />
  <img src="https://raw.github.com/hzqtc/cashd/master/screenshots/account_view.png" width="400" />
</p>

## Features

- **Interactive TUI:** Navigate through your financial data with an intuitive and responsive terminal interface.
- **Multiple Views:**
  - **Transactions:** View a detailed list of all your financial transactions, with sorting and searching capabilities.
  - **Accounts:** Get an overview of your financial accounts, including balances and transaction insights.
  - **Categories:** Analyze your spending and income by category, helping you understand where your money goes.
- **Flexible Data Loading:** Supports loading financial data from various sources.
  - **Configurable CSV Parsing:** Customize how `cashd` interprets your CSV files to match your data's format.
- **Date Range Filtering:** Filter transactions by custom date ranges (weekly, monthly, quarterly, annually) to focus on specific periods.
- **Search Functionality:** Quickly find specific transactions using keywords.
- **Financial Insights:** Visualize your financial trends with time-series charts for accounts and categories.

## Limitations

The following limitations are known:

- Only supports `Income` and `Expense` transaction types
- Only supports `Cash`, `Bank Account` and `Credit Card` as account types
- Only supports `$` as the currency
- Specficially for `ledger` transactions
  - Only supports 2 postings per transaction

### Supported Data Sources

`cashd` is designed to be flexible with your financial data. Currently, it supports:

- **CSV Files:** Load transactions from a standard CSV file. `cashd` provides extensive configuration options to correctly parse your CSV data.
- **Ledger/Hledger:** Integrate seamlessly with popular plain-text accounting tools like `ledger` and `hledger` by parsing their journal files.
  - Note: `cashd` invokes `ledger print` or `hledger print` and does not read journal files directly

## Installation

### Prerequsites

- A nerd font enabled terminal

### Prebuilt binary

Prebuilt binaries can be found in Github release page.

### Build from source

To build `cashd`, ensure you have Go installed (version 1.18 or higher).

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/your-username/cashd.git
    cd cashd
    ```
2.  **Build the application:**
    ```bash
    make
    ```
3.  **Installs binary to `~/.local/bin`:**
    ```bash
    make install
    ```

## Usage

### Loading Data from Ledger/Hledger (default)

To load transactions from a Ledger or Hledger journal file, use the `--ledger` flag:

```bash
cashd --ledger path/to/your/journal.dat
```

Alternatively, you can set the `LEDGER_FILE` or `HLEDGER_FILE` environment variables:

```bash
export LEDGER_FILE=/path/to/your/journal.dat
cashd
```

### Loading Data from a CSV File

To load transactions from a CSV file, use the `--csv` flag:

```bash
cashd --csv path/to/your/transactions.csv
```

If your CSV file has a custom format, you can provide a JSON configuration file using the `--csv-config` flag:

```bash
cashd --csv path/to/your/transactions.csv --csv-config path/to/your/config.json
```

### Generating a Sample CSV File

The `sample` directory contains a Go program (`gencsv.go`) that can generate a sample CSV file for testing purposes.

```bash
go run sample/gencsv.go -lines 10000
cashd --csv sample.csv --csv-config sample-csv-config.json
```

This will create a `sample.csv` file in the current directory with 10000 random transactions.

### Command Line Flags

- `-h`, `--help`: Show help message.
- `--csv <file_path>`: Specify the path to your CSV transaction file.
- `--csv-config <file_path>`: Specify the path to your CSV configuration JSON file.
- `--ledger <file_path>`: Specify the path to your Ledger/Hledger journal file.
- `--debug`: Enable debug mode (useful for development and troubleshooting).

## CSV Configuration File Format

The CSV configuration file is a JSON file that defines how `cashd` should parse your CSV data.
This is particularly useful if your CSV columns or data formats differ from the default expectations.

Here's an example of the structure:

```json
{
  "columns": {
    "Period": "Date",
    "Accounts": "Account",
    "Category": "Category",
    "Note": "Description",
    "USD": "Amount",
    "Income/Expense": "Type"
  },
  "date_formats": [
    "2006-01-02",
    "2006-01-02 15:04:05",
    "01/02/2006",
    "01/02/2006 15:04:05"
  ],
  "transaction_types": {
    "income": "Income",
    "inc.": "Income",
    "expense": "Expense",
    "exp.": "Expense",
    "exps.": "Expense"
  },
  "account_types": {
    "cash": "Cash",
    "bank": "Bank Account",
    "credit card": "Credit Card"
  },
  "account_type_from_name": {
    "^cash$": "Cash",
    "checking$": "Bank Account",
    "saving(s)?$": "Bank Account",
    "card$": "Credit Card"
  }
}
```

### Config fields:

- `columns`: A map where keys are the actual column headers in your CSV file, and values are the corresponding internal `TransactionField` names (`Date`, `Type`, `AccountType`, `Account`, `Category`, `Amount`, `Description`).
- `column_indexes` (Optional): A map where keys are `TransactionField` names and values are the 0-based index of the column in your CSV. If not provided, `cashd` will attempt to infer column indexes from the `columns` mapping and the CSV header.
- `date_formats`: An array of Go time format strings that `cashd` will attempt to use when parsing the `Date` column. The first format that successfully parses the date will be used.
- `transaction_types`: A map where keys are string values found in your CSV's "Type" column, and values are the internal `TransactionType` (`Income` or `Expense`). This allows `cashd` to understand various representations of income and expense in your data.
- `account_types`: A map where keys are string values found in your CSV's "AccountType" column, and values are the internal `AccountType` (`Cash`, `Bank Account`, `Credit Card`).
- `account_type_from_name`: A map where keys are regular expressions that will be matched against the `Account` name (case-insensitive), and values are the `AccountType` to assign if a match is found. This is useful for inferring account types when they are not explicitly provided in your CSV. If no match is found, it defaults to `Credit Card`.

## Credit

This project is built using:

- [bubbletea](https://github.com/charmbracelet/bubbletea)
- [bubbles](https://github.com/charmbracelet/bubbles)
- [lipgloss](https://github.com/charmbracelet/lipgloss)
- [ntcharts](https://github.com/NimbleMarkets/ntcharts)
- [pflag](https://github.com/spf13/pflag)
