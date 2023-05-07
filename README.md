# curconv

A very simple currency converter CLI.

I heard somewhere that Google asked this for their phone screen; actually neat application of graph traversal. You could think of this as that currency
converter on Google that lets you choose any pair and get the rate for it.

Written with Go because why not; it's good at computation.

### Usage

```
NAME:
   curconv - A currency converter CLI tool that can find the rate between any two pairs in your given currency list

USAGE:
   ./curconv [global options] <currency from> <currency to> [args]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --verbose, -v           Do it verbosely (default: false)
   --file value, -f value  Path to CSV file containing comma-delimited currency pairs
   --help, -h              show help

EXAMPLE:
$ cat currencies.csv
USD,CAD,1.35
CHF,CAD,1.53

$ ./curconv -f currencies.csv -v CHF USD
2023/05/06 23:42:22 Created adjacency list of currency pairs: map[CAD:[{USD 0.7407407407407407} {CHF 0.6535947712418301}] CHF:[{CAD 1.53}] USD:[{CAD 1.35}]]
2023/05/06 23:42:22 Done traversing currency pairs; successor table: map[CAD:{USD 0.7407407407407407} CHF:{CAD 1.53}]
Steps:
1 CHF = 1.530000 CAD
1 CAD = 0.740741 USD
--------------
1 CHF = 1.133333 USD


```

### TODO

- [ ] benchmark it against big dataset (us_fxrate.csv but need to put it in pair style)
