EXTENSION=gh-subrepo

remove_ext:
	gh extension remove $(EXTENSION)

build:
	go mod tidy && go build -o $(EXTENSION)
	gh extension install .

rebuild: remove_ext build

clean:
	rm -rf gh-subrepo coverage.*
