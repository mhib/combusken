# Combusken
Combusken is a UCI-compliant open source chess engine using [Alpha-beta algorithm](https://en.wikipedia.org/wiki/Alpha%E2%80%93beta_pruning). You can play with on [lichess](https://lichess.org/@/combuskengine).


## UCI options
### Hash
Size of transposition table in megabytes. Usually the more the better.
### Threads
Number of threads used in search. Usually the more the better.
### PawnHash
Size of Pawn Hash Table. Default value should always work ok, as hit-ratio in Pawn Hash Table is usually pretty high.
### Move Overhead
Time buffer in ms. Should be increased when you notice time-losses.

## CLI options
### `combusken bench`
Runs benchmark

### `combusken tune`
Runs tuning that is a combination of coordinate descent and gradient descent where gradient is calculated with symmetric derivative.

### `combusken trace-tune`
Runs tuning based on gradient descent where gradient is calculated with a vectors that stores how much each evaluation-constant was used in a given position.
In order to work it requires compilation with `tuning` constant set to `true` in `evaluation/eval.go` file.

Games for tuning must be put in `games.fen` file.

## Thanks
+ [Counter](https://github.com/ChizhovVadim/CounterGo) by Vadim Chizhov

UCI protocol implementation and search cancelation pattern is based on CounterGO's. 
Also some miscellaneous things like LMP weights, or EPD parsing.
 
+ [Ethereal](https://github.com/AndyGrant/Ethereal) by Andrew Grant, Alayan & Laldon 

Combusken's search procedure is heavily inspired by Ethereal's, some parts of evaluation(for example king safety) were taken directly from it.
Andrew Grant's [OpenBench](https://github.com/AndyGrant/OpenBench) is used for testing.

+ [Laser](https://github.com/jeffreyan11/laser-chess-engine/) by Jeffrey An

SEE, Internal Iterative Deepening constants.

+ [Stockfish](https://github.com/official-stockfish/Stockfish/) by Tord Romstad, Marco Costalba, Joona Kiiski & Gary Linscott

Some evaluation heuristics.

+ [Zurichess](https://bitbucket.org/zurichess/zurichess/src) by Alexandru Moșoi

It's [tuning positions set](http://www.zurichess.xyz/blog/texels-tuning-method/) was used in Combusken's tuning

## License
Combusken is distributed under the GNU General Public License version 3 (GPL v3). 