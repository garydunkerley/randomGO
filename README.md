# randomGO

![random_goban](https://user-images.githubusercontent.com/58280110/113188912-3d4d0a80-9220-11eb-9a73-12762b3b4a99.png)

# ** hiatus, feel free to fork! ** 

Play go/igo/baduk on randomly generated boards! Gameplay is nearly identical to traditional go with the Japanese ruleset. Future iterations of this project will implement a scoring algorithm that should be able to account for simple cases of seki. This project also allows play on rectangular boards.

Boards are created by constructing a regular equilateral lattice (no position has more than six liberties) and then performing a random but controlled removal of edges. Board visualization is done using the go-graphviz package.


# TODO
 - Finish .sgf parser and custom hypertext format 
 - Fix capturing error (captured stones aren't treated as liberties?)
 - Fix score estimator
 - Implement Zobrist hash for super ko
 - Reorganize code into event / command scripts

