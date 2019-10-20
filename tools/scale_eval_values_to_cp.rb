#!/usr/local/bin/ruby

INITAL_PAWN_VALUE = 173
SCALE = 100.0 / INITAL_PAWN_VALUE
REGEXP = /Score{(-?\d+, -?\d+)}/

EVAL_PATH = File.join(__dir__, '..', 'evaluation', 'eval.go')

eval_content = IO.read(EVAL_PATH)

new_eval = eval_content.gsub(REGEXP) do |match|
  middle_score, end_score = match.scan(/-?\d+/).map(&:to_i)
  "Score{#{(middle_score * SCALE).ceil}, #{(end_score * SCALE).ceil}}"
end

IO.write(EVAL_PATH, new_eval, mode: 'w')
