EXTENSION=gh-subrepo

remove_ext:
	gh extension remove $(EXTENSION)

build: remove_ext
	go build -o $(EXTENSION)
	gh extension install .
