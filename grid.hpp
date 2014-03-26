// The grid representation

#pragma once

#include <iostream>
#include <cstdint>
#include <cstdlib>
#include <random>
#include <limits>
#include <algorithm>
#include <initializer_list>

const size_t rows = 4;
const size_t cols = 4;

// A thread_local random number generator
__thread std::mt19937_64 gen;
__thread std::uniform_real_distribution<double> real_dist;

// Directions of movement
const size_t LEFT = 0;
const size_t RIGHT = 1;
const size_t UP = 2;
const size_t DOWN = 3;

// Converts direction to string
const char* direction_to_string(size_t d) {
    switch (d) {
    case LEFT:
        return "LEFT";
    case RIGHT:
        return "RIGHT";
    case UP:
        return "UP";
    case DOWN:
        return "DOWN";
    }
    return "Invalid direction";
}

using gridnum = uint16_t;

class Grid {
public:
    gridnum grid[rows][cols];
    size_t score;

    Grid(): score(0) {
        memset(grid, 0, sizeof(Grid));
    }

    // Clones the given grid
    Grid(const Grid& g) {
        memcpy(this, &g, sizeof(Grid));
    }

    void operator=(const Grid& g) {
        memcpy(this, &g, sizeof(Grid));
    }

    // Places a 2 or 4 tile (90% chance it's a 2) at a random place in
    // the board. Assumes the random number generator is already
    // seeded. If there is no place for a tile, it returns false.
    bool place_random() {
        // The random number is a position to start searching at, which
        // wraps around.
        const size_t total = rows*cols;
        gridnum tileval = 2;
        if (real_dist(gen) < 0.1) {
            tileval = 4;
        }
        size_t startpos = gen() % total;
        for (size_t i = startpos+1 % total; i != startpos; i = (i+1) % total) {
            size_t row = i / cols;
            size_t col = i % cols;
            if (grid[row][col] == 0) {
                grid[row][col] = tileval;
                return true;
            }
        }
        return false;
    }

    bool operator==(const Grid& g) {
        return memcmp(this, &g, sizeof(Grid)) == 0;
    }

    bool operator!=(const Grid& g) {
        return !(*this == g);
    }

    // Move moves the tiles in the grid in the given direction. It
    // returns false if no tiles moved, and true otherwise.
    bool move(size_t d) {
        bool ret = false;
        switch (d) {
        case LEFT:
        case RIGHT:
            // Goes row by row
            for (size_t r = 0; r < rows; r++) {
                // The column to merge with, start column, end column,
                // and column increment
                int moveCol, start, end, inc;
                if (d == LEFT) {
                    moveCol = 0;
                    start = 1;
                    end = cols;
                    inc = 1;
                } else {
                    moveCol = cols-1;
                    start = cols-2;
                    end = -1;
                    inc = -1;
                }
                for (int c = start; c != end; c += inc) {
                    if (grid[r][c] == 0) {
                        continue;
                    } else if (grid[r][moveCol] == 0) {
                        // Move grid[r][c] all the way down
                        std::swap(grid[r][c], grid[r][moveCol]);
                        ret = true;
                    } else if (grid[r][moveCol] == grid[r][c]) {
                        // Merge grid[r][c] with grid[r][moveCol]
                        grid[r][c] = 0;
                        grid[r][moveCol] *= 2;
                        score += grid[r][moveCol];
                        ret = true;
                    } else {
                        // Increment moveCol and move grid[r][c]
                        // there, if it isn't already
                        moveCol += inc;
                        if (moveCol != c) {
                            std::swap(grid[r][c], grid[r][moveCol]);
                            ret = true;
                        }
                    }
                }
            }
            break;
        case UP:
        case DOWN:
            // Goes column by column
            for (size_t c = 0; c < cols; c++) {
                // The row to merge with, start row, end row, and row
                // increment
                int moveRow, start, end, inc;
                if (d == UP) {
                    moveRow = 0;
                    start = 1;
                    end = rows;
                    inc = 1;
                } else {
                    moveRow = rows-1;
                    start = rows-2;
                    end = -1;
                    inc = -1;
                }
                for (int r = start; r != end; r += inc) {
                    if (grid[r][c] == 0) {
                        continue;
                    } else if (grid[moveRow][c] == 0) {
                        // Move grid[r][c] all the way down
                        std::swap(grid[r][c], grid[moveRow][c]);
                        ret = true;
                    } else if (grid[moveRow][c] == grid[r][c]) {
                        // Merge grid[r][c] with grid[r][moveCol]
                        grid[r][c] = 0;
                        grid[moveRow][c] *= 2;
                        score += grid[moveRow][c];
                        ret = true;
                    } else { 
                        // Increment moveRow and move grid[r][c] there, if
                        // it isn't already
                        moveRow += inc;
                        if (moveRow != r) {
                            std::swap(grid[r][c], grid[moveRow][c]);
                            ret = true;
                        }
                    }
                }
            }
            break;
        default:
            std::cerr << "Invalid direction" << std::endl;
            exit(1);
        }
        return ret;
    }
};

// from_sparse_grid takes a slice consisting of triples of numbers,
// where the first is the row, the second the column, and the third
// the value, and creates a grid.
Grid from_sparse_grid(std::initializer_list<size_t> vals, size_t score) {
    Grid ret;
    ret.score = score;
    for (auto it = vals.begin(); it < vals.end(); it += 3) {
        ret.grid[*it][*(it+1)] = *(it+2);
    }
    return ret;
}

// Overload for printing
std::ostream& operator<< (std::ostream& stream, const Grid& g) {
    for (size_t i = 0; i < rows; i++) {
        for (size_t j = 0; j < cols; j++) {
            stream << g.grid[i][j] << '\t';
        }
        stream << std::endl;
    }
    stream << "Score: " << g.score;
    return stream;
}
