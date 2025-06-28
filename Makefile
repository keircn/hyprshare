BINARY=build/hyprshare
SRC=cmd/hyprshare/main.go

.PHONY: all clean

all: $(BINARY)

$(BINARY): $(SRC)
	mkdir -p $(dir $(BINARY))
	go build -o $(BINARY) $(SRC)

clean:
	rm -rf build