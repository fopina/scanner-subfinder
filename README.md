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

## Parameters

```
Usage of scan:
  -b, --bin string      Path to scanner binary (default "subfinder")
  -o, --output string   Scanner results directory (default "/output")
  -H, --scanner-help    Show help for the scanner extra flags
```

