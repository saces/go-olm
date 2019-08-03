Fork from https://github.com/Dhole/go-olm

build:
issue all commands from this directory

get a decent version of olm
git clone https://gitlab.matrix.org/matrix-org/olm.git olm

build static version

cmake olm/ -Bolm/build -DBUILD_SHARED_LIBS=NO

cmake --build olm/build


run 'go test -v'

you may need  'go get -u github.com/fatih/structs'


--
# Go olm/megolm bindings [![GoDoc](https://godoc.org/github.com/Dhole/go-olm?status.svg)](https://godoc.org/github.com/Dhole/go-olm)



Work In Progress.

All original library functions exposed and documented.

Unstable API for now, as I'm figuring out the best way to use the functions
from the point of view of a matrix chat client.
