NAME = "logshipper"

GOPATH = "/tmp/_logshipper_gopath"

all: clean $(NAME)

$(NAME):
	mkdir ${GOPATH}
	GOPATH=${GOPATH} go build -v

package: clean $(NAME)
	equivs-build debian/logshipper.equivs

clean:
	rm -f logshipper
	rm -rf ${GOPATH}
