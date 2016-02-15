NAME = "logshipper"

all: clean $(NAME)

$(NAME):
	go build -v

package:
	gbp buildpackage --git-pbuilder

clean:
	rm -f logshipper
