#!/usr/bin/env bash
cutechess-cli \
        -srand $RANDOM \
        -pgnout games.pgn \
        -repeat \
        -recover \
        -tournament gauntlet \
        -rounds 500000 \
        -concurrency 3 \
        -ratinginterval 50 \
        -draw movenumber=50 movecount=5 score=20 \
        -openings file=./2moves_v1.pgn format=pgn order=random \
        -engine cmd=../combusken name=combusken1 tc=40/2+0.05 \
        -engine cmd=../combusken name=combusken2 tc=40/2+0.05 \
        -each timemargin=60000 option.Hash=128 proto=uci
