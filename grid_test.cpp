#include "grid.hpp"
#include <random>
#include <chrono>
#include <cassert>

void test_clone() {
	Grid g;
	g.grid[0][0] = 2;
    g.score = 10;

    Grid g2(g);

	assert(g == g2);

	g2.grid[0][1] = 20;
    assert(g != g2);
}

void test_movement() {
    Grid g;
	g.grid[0][0] = 2;
	g.grid[0][1] = 2;

    Grid comp = from_sparse_grid({0, 0, 2, 0, 1, 2}, 0);
    assert(g == comp);

	g.move(RIGHT);
    assert(g == from_sparse_grid({0, 3, 4}, 4));

	g.grid[1][2] = 4;
	g.move(UP);
    assert(g == from_sparse_grid({0, 2, 4, 0, 3, 4}, 4));

	g.move(LEFT);
    assert(g == from_sparse_grid({0, 0, 8}, 12));

	g.grid[1][0] = 2;
	g.grid[2][2] = 16;
	g.grid[3][2] = 16;
	g.move(DOWN);
    assert(g == from_sparse_grid({3, 0, 2, 2, 0, 8, 3, 2, 32}, 44));
}

void test_place_random() {
    Grid g;
    assert(g.place_random());

	gridnum total = 0;
	for (size_t r = 0; r < rows; r++) {
		for (size_t c = 0; c < cols; c++) {
			total += g.grid[r][c];
		}
	}
    assert(total == 2 || total == 4);


	for (size_t r = 0; r < rows; r++) {
		for (size_t c = 0; c < cols; c++) {
            g.grid[r][c] = 4;
		}
	}
    
    assert(!g.place_random());
}

int main() {
    // Seeds the random number generator
    gen.seed(std::chrono::system_clock::now().time_since_epoch().count());
    test_clone();
    test_movement();
    test_place_random();
}
