.PHONY: clean

all: reflex-scheduler reflex-executor

reflex-scheduler: reflex/
	go build -o reflex-scheduler ./cmd/reflex/

reflex-executor: executor/
	go build -o reflex-executor ./executor/

clean:
	rm reflex-scheduler reflex-executor
