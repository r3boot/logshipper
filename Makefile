NAME = "logshipper"

all: clean $(NAME)

$(NAME):
	go build -v

clean:
	rm -f logshipper
