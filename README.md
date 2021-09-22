uMagNUS
======

GPU accelerated micromagnetic simulator based on OpenCL.


Downloads and documentation
---------------------------

The frontend is based on MuMax3 and accepts simulation files written for MuMax3.

Refer to MuMax3 documentation at:
http://mumax.github.io


Paper
-----

- To be updated.


Building from source (for linux)
--------------------

Consider downloading a pre-compiled binary. If you want to compile nevertheless:

  * install the OpenCL driver, if not yet present.
   - if unsure, it's probably already there
   - requires OpenCL 1.2 support
  * install Go 
    - https://golang.org/dl/
    - set $GOPATH
  * if you have git installed: 
    - `go install -v github.com/seeder-research/uMagNUS/cmd/`
  * if you don't have git:
    - seriously, no git?
    - get the source from https://github.com/seeder-research/uMagNUS/releases
    - unzip the source into $GOPATH/src/github.com/seeder-research/uMagNUS
    - `cd $GOPATH/src/github.com/seeder-research/uMagNUS/cmd/uMagNUS`
    - `go install`
  * optional: install gnuplot if you want pretty graphs
    - Ubuntu: `sudo apt-get install gnuplot`

Your binary is now at `$GOPATH/bin/umagnus`

To do all at once on Ubuntu:
```
sudo apt-get install git golang-go gcc gnuplot
export GOPATH=$HOME go install -u -v github.com/seeder-research/uMagNUS/cmd/uMagNUS
```

Contributing
------------

Contributions are gratefully accepted. To contribute code, fork the repo on github and send a pull request.
