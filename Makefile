BASIC_COMMAND = $(CXX) -std=c++11 -stdlib=libc++ -Wall -Wextra -Werror
COMMAND = $(BASIC_COMMAND) -O3 -o simulation
DEBUG_COMMAND = $(BASIC_COMMAND) -g -o simulation

all: grid.hpp ai.hpp simulation.cpp
	$(COMMAND) simulation.cpp

debug: grid.hpp ai.hpp simulation.cpp
	$(DEBUG_COMMAND) simulation.cpp

check: grid.hpp grid_test.cpp
	$(BASIC_COMMAND) -g -o grid_test grid_test.cpp
	./grid_test

clean:
	rm -f simulation grid_test
	rm -rf *.dSYM
