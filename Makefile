EXTENSION=gh-subrepo

remove_ext:
	gh extension remove $(EXTENSION)

build: remove_ext
	go mod tidy && go build -o $(EXTENSION)
	gh extension install .

clean:
	rm -rf gh-subrepo coverage.*
