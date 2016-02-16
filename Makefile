NAME = "logshipper"

all: clean $(NAME)

$(NAME):
	go build -v

package: clean $(NAME)
	equivs-build debian/logshipper.equivs

clean:
	rm -f logshipper
