build:
	docker build -t tmp-$(notdir $(CURDIR)) .

cleantest:
	rm -fr testdata/output

test: cleantest
	go run main.go -v -output-dir testdata/output/ testdata/inp1.txt

dockertest: cleantest build
	docker run --rm -v $(PWD)/testdata:/input tmp-$(notdir $(CURDIR)) /input/inp1.txt
