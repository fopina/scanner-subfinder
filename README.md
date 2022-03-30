# Surface Scanner: subfinder



## Usage

* Create `input/input.txt` (example in [testdata](testdata/inp1.txt))
* Run it
  ```
  docker run --rm -v $(pwd)/input:/input:ro \
                  -v $(pwd)/output:/output \
                  ghcr.io/fopina/scanner-subfinder \
                  /input/input.txt
  ```

Changes made here do not break any standard usage of subfinder, so it can also be used as original subfinder:

```
docker run --rm ghcr.io/fopina/scanner-subfinder \
                -d github.io
```


## Parameters

```
Subfinder is a subdomain discovery tool that discovers subdomains for websites by using passive online sources.

Usage:
  ./subfinder [flags]

Flags:
INPUT:
   -d, -domain string[]  domains to find subdomains for
   -dL, -list string     file containing list of domains for subdomain discovery

SOURCE:
   -s, -sources string[]           specific sources to use for discovery (-s crtsh,github
   -recursive                      use only recursive sources
   -all                            Use all sources (slow) for enumeration
   -es, -exclude-sources string[]  sources to exclude from enumeration (-es archiveis,zoomeye)

RATE-LIMIT:
   -rl, -rate-limit int  maximum number of http requests to send per second
   -t int                number of concurrent goroutines for resolving (-active only) (default 10)

OUTPUT:
   -o, -output string       file to write output to
   -oJ, -json               write output in JSONL(ines) format
   -oD, -output-dir string  directory to write output (-dL only)
   -cs, -collect-sources    include all sources in the output (-json only)
   -oI, -ip                 include host IP in output (-active only)

CONFIGURATION:
   -r string[]         comma separated list of resolvers to use
   -rL, -rlist string  file containing list of resolvers to use
   -nW, -active        display active subdomains only
   -proxy string       http proxy to use with subfinder

DEBUG:
   -silent             show only subdomains in output
   -version            show version of subfinder
   -v                  show verbose output
   -nc, -no-color      disable color in output
   -ls, -list-sources  list all available sources

OPTIMIZATION:
   -timeout int   seconds to wait before timing out (default 30)
   -max-time int  minutes to wait for enumeration results (default 10)
```
