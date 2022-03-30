build:
	docker build -t tmp-$(notdir $(CURDIR)) .

test:
	rm -fr testdata/output
	go run main.go -v -output-dir testdata/output/ testdata/inp1.txt
