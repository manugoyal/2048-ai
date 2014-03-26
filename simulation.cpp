#include "grid.hpp"
#include "ai.hpp"

#include <iostream>

const size_t height = 6;
const size_t thread_num = std::thread::hardware_concurrency();
const size_t reps = thread_num*2;

int main() {
    Grid g;
    size_t moves = 0;
    g.place_random();
    g.place_random();

    std::cout << g << std::endl;

    while (true) {
        int best_move = next_move(g, height, reps, thread_num);
        if (best_move == -1) {
            std::cout << "Couldn't find good direction. Game over" << std::endl;
			break;
        }

        std::cout << "Moving " << direction_to_string(best_move) << std::endl;
        g.move(best_move);
        moves++;

        if (!g.place_random()) {
            std::cout << "Couldn't place piece. Game over" << std::endl;
			break;
        }

        std::cout << g << std::endl;
    }

    std::cout << "After " << moves << " moves:" << std::endl;
    std::cout << g << std::endl;
}
