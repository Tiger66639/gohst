# gohst
> A webhosting platform powered by golang

Created for the sake of learning and memes

## Installation & Running
>Prerequisite: go must be installed and nginx or apache must listen for fcgi
on port 8000

To install, simply run the following

    $ go get github.com/cosban/gohst
    $ cd $GOPATH/src/github.com/cosban/gohst && go install

Once this has been done, cd into the directory containing your static and
template directories then run the following:

    $ gohst

This is excellent for running inside screen sessions and whatnot. The preferred
method to run is in the following manner:

    #!/bin/bash
    if [ "$(pidof gohst)" ]
    then
        killall gohst
    fi
    gohst > log 2> err &

This script can also be found at [gohst-example](https://github.com/cosban/gohst-example/master/start.sh)
You can also view a website running gohst at [cosban.net](https://cosban.net)
