// Creates a game tree and tries to make decisions that maximize the
// score. It's not a minimax algorithm, since the opponent's decisions
// are random.
#pragma once

#include "grid.hpp"
#include <thread>
#include <vector>
#include <atomic>
#include <algorithm>
#include <cmath>
#include <chrono>
#include <random>

class Tree {
public:
    Grid grid;
    Tree** children;
    const static size_t numchildren = 4;
    size_t best_score;
    int best_direction;

    // NewTree returns an tree with an empty grid
    Tree(): children(nullptr), best_score(0),
            best_direction(0) {
    }

    Tree(const Grid& g): grid(g), children(nullptr),
                         best_score(0), best_direction(0) {}

    // Frees the children
    ~Tree() {
        if (children != nullptr) {
            for (size_t i = 0; i < numchildren; i++) {
                if (children[i] != nullptr) {
                    delete children[i];
                }
            }
            delete[] children;
            children = nullptr;
        }
    }

    // Given the height of the tree, it will fill out the tree to
    // nodes of height 0. If the tree already has children, it won't
    // generate new ones, but it will recursively call fill. This
    // should allow for iterative deepening.
    void fill(size_t height) {
        if (height == 0) {
            return;
        }
        // Generate new children
        children = new Tree*[numchildren];
        for (size_t i = 0; i < 4; i++) {
            Tree* node = new Tree;
            node->grid = grid;
            if (node->grid.move(i)) {
                // We only execute the move if tiles would be moving
                children[i] = node;
                if (node->grid.place_random()) {
                    node->fill(height-1);
                }
            } else {
                children[i] = nullptr;
                delete node;
            }
        }
    }

    // Given a filled tree, it takes the scores of the leaf nodes and
    // fills up the scores and directions of the parent nodes using
    // the minimax algorithm. The root of the tree will contain the
    // best score and direction
    void score() {
        if (children == nullptr) {
            best_score = grid.score;
            best_direction = -1;
        } else {
            best_score = 0;
            best_direction = -1;
            for (size_t i = 0; i < numchildren; i++) {
                if (children[i] != nullptr) {
                    children[i]->score();
                    if (children[i]->best_score > best_score) {
                        best_score = children[i]->best_score;
                        best_direction = i;
                    }
                }
            }
        }
    }
};

// Given a grid and some parameters, it figures out the next best
// move. If it returns -1, that means it couldn't find a move.
int next_move(const Grid& g, const size_t height, const size_t reps,
              const size_t num_threads) {
    std::atomic<size_t> counts[4];
    for (size_t i = 0; i < 4; i++) {
        counts[i].store(0);
    }
	// We round the number of reps to a multiple of thread_num when
	// calculating reps_per_thread
    size_t reps_per_thread = (size_t) ceil((float)reps / (float)num_threads);

    std::vector<std::thread> threads;
    auto next_move_helper = [&counts, reps_per_thread, g, height]() {
        // Seeds the random number generator and sets the real_dist
        gen.seed(std::chrono::system_clock::now().time_since_epoch().count());
        real_dist.param(std::uniform_real_distribution<double>(0.0, 1.0).param());
        for (size_t i = 0; i < reps_per_thread; i++) {
            Tree t(g);
            t.fill(height);
            t.score();
            if (t.best_direction >= 0 && t.best_direction < 4) {
                counts[t.best_direction].fetch_add(1);
            }
        }
    };

    for (size_t i = 0; i < num_threads; i++) {
        threads.emplace_back(next_move_helper);
    }
    for (size_t i = 0; i < num_threads; i++) {
        threads[i].join();
    }

    int best_ind = -1;
    size_t best_occ = 0;
    for (size_t i = 0; i < 4; i++) {
        if (counts[i].load() > best_occ) {
            best_occ = counts[i].load();
            best_ind = i;
        }
    }
    return best_ind;
}
