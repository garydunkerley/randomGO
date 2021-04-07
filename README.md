# randomGO

![random_goban](https://user-images.githubusercontent.com/58280110/113188912-3d4d0a80-9220-11eb-9a73-12762b3b4a99.png)

Play go/igo/baduk on randomly generated boards! Gameplay is nearly identical to traditional go with the Japanese ruleset. The game is equipped with a scoring algorithm that should be able to detect simple instances of seki. This project also allows play on rectangular boards and is equipped with a .sgf parser. We are working on making an enriched text format so that randomGO games can be saved and studied. 

Boards are created by constructing a regular equilateral lattice (no position has more than six liberties) and then performing a random but controlled removal of edges. Board visualization is done using the go-graphviz package.


# TODO
 - Fix capturing error (captured stones aren't treated as liberties?)
 - Fix score estimator
 - Tweak random board generation
 - Implement Zobrist hash for super ko
 - Reorganize code into event / command scripts

