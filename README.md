# gohst
> A webhosting platform powered by golang

Created for the sake of education and memes
##Installation & Running
To install, simply run the following

    $ go get github.com/cosban/gohst
    $ cd $GOPATH/src/github.com/cosban/gohst/gohst && go install

Once this has been done, cd into the directory containing your static and
template directories then run the following:

    $ gohst

This is excellent for running inside screen sessions and whatnot. The preferred
method to run is in the following manner though:

    $ gohst > log 2> err &

Which allows it to run in a detached process while also logging to files

## File Structure
In order for files to be served correctly, files need to be arranged in a
general manner. The structure is pretty leanient, but must fall under the
following pattern

    static/       # static content which is not rendered through templating
        txt/      # .txt files which are displayed inside a <pre> block
    templates/    # content which must be rendered through the go code
        base.html # The base template which all others are rendered through
        home.html # This is the web app's home page

###TODO: Example gohst web app link et al
