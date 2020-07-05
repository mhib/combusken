from skopt import gp_minimize
from tempfile import NamedTemporaryFile
from os import chdir, system, pipe, remove, fdopen, close
from shutil import copy
from typing import List
import subprocess
import re


NUMBER_OF_THREADS = 8
NUMBER_OF_GAMES = 20

class SearchConstant:
    def __init__(self, name, range):
        self.name = name
        self.range = range

def cute_chess_command(engine_1, engine_2):
    return "cutechess-cli -repeat -recover -wait 10 -resign movecount=3 score=400 -draw movenumber=40 movecount=8 score=10 -concurrency {} -games {} -engine cmd={} option.Hash=64 -engine cmd={} option.Hash=64 -each proto=uci tc=60+0.6 -openings file=./2moves_v1.pgn format=pgn order=random plies=16".format(
        NUMBER_OF_THREADS,
        NUMBER_OF_GAMES,
        engine_1,
        engine_2
    )

def const_pattern(parameter: str):
    return re.compile(r'const ' + parameter + r' = (-?\d+)$')

def replace_value(parameter: str, value: int):
    pattern = const_pattern(parameter)
    with NamedTemporaryFile(mode='w') as new_file:
        for _i, line in enumerate(open('../engine/search.go')):
            if pattern.match(line):
                new_file.write('const ' + parameter + ' = ' + str(value) + '\n')
            else:
                new_file.write(line)
        new_file.flush()
        copy(new_file.name, '../engine/search.go')

def get_current_value(pattern):
    for _i, line in enumerate(open('../engine/search.go')):
        for match in re.finditer(pattern, line):
            return int(match.group(1))


elo_regex = re.compile(r'Elo difference: (-?\d+\.?\d*)')
inf_regex = re.compile(r'Elo difference: (-?inf)')
def calculate_elo_diff(parameters: List[SearchConstant], values: List[int]):
    calculate_elo_diff.call_count += 1
    print("{} call {}".format(calculate_elo_diff.call_count, values))
    binary_name = '-'.join(["combusken"] + [p.name for p in parameters] + list(map(str, values)))
    for parameter, value in zip(parameters, values):
        replace_value(parameter.name, value)
    build_engine(binary_name)
    (pipe_read, pipe_write) = pipe()
    process = subprocess.Popen(
        cute_chess_command('../combusken', '../' + binary_name).split(),
        stdout=pipe_write
    )
    close(pipe_write)
    with fdopen(pipe_read) as fd:
        while True:
            line = fd.readline().strip()
            print(line)
            if line.startswith('Elo difference'):
                process.kill()
                remove('../' + binary_name)

                if inf_regex.match(line):
                    match = inf_regex.match(line)
                    return float(match.group(1))

                if elo_regex.match(line):
                    match = elo_regex.match(line)
                    return float(match.group(1))

def build_engine(exe_name):
    chdir('..')
    system('go build -o {} combusken.go'.format(exe_name))
    chdir('./tools')


def optimize_parameter(parameters: List[SearchConstant]):
    calculate_elo_diff.call_count = 0
    build_engine('combusken')
    current_values = [get_current_value(const_pattern(x.name)) for x in parameters]
    res = gp_minimize(
        lambda xs: calculate_elo_diff(parameters, list(map(round, xs))),
        [(value - parameter.range, value + parameter.range) for value, parameter in zip(current_values, parameters)],
        n_calls=100,
        n_random_starts=5,
        acq_func="EI",
    )
    print(res)

if __name__ == '__main__':
    optimize_parameter([SearchConstant('seePruningDepth', 2), SearchConstant('seeQuietMargin', 20), SearchConstant('seeNoisyMargin', 20)])
