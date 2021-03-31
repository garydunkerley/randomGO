# randomGO

![board](https://user-images.githubusercontent.com/58280110/113186815-bb5be200-921d-11eb-9c7a-14403291b5ec.png)

Play go/igo/baduk on randomly generated boards! Gameplay is nearly identical to traditional go with the Japanese ruleset. The game is equipped with a scoring algorithm that should be able to detect simple instances of seki. This project also allows play on rectangular boards and is equipped with a .sgf parser. We are working on making an enriched text format so that randomGO games can be saved and studied. 

Boards are created by constructing a regular equilateral lattice (no position has more than six liberties) and then performing a random but controlled removal of edges. Board visualization is done using the go-graphviz package.


# TODO
 - Tweak random board generation
 - Implement Zobrist hashing for super ko rule
 - Create format for random games
 - Create UI
 - create algorithms to play random games
 - Develop network play

