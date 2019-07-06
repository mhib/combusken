#!/usr/local/bin/ruby

RESULT = '[Score{186, 216} Score{-17, -2} Score{29, -15} Score{11, -3} Score{33, -17} Score{35, -22} Score{16, -9} Score{34, -12} Score{-18, 0} Score{-11, -29} Score{-23, -28} Score{-3, -35} Score{-6, -40} Score{-7, -39} Score{9, -41} Score{-17, -35} Score{-6, -33} Score{-24, -10} Score{-24, -22} Score{7, -46} Score{36, -53} Score{33, -51} Score{11, -43} Score{-21, -21} Score{-21, -13} Score{-1, 40} Score{40, 16} Score{29, -12} Score{76, -42} Score{70, -50} Score{19, 0} Score{49, 10} Score{-1, 37} Score{46, 203} Score{74, 198} Score{136, 114} Score{134, 101} Score{146, 88} Score{158, 124} Score{44, 205} Score{39, 211} Score{153, 510} Score{227, 473} Score{156, 400} Score{191, 363} Score{281, 368} Score{129, 424} Score{189, 460} Score{30, 568} Score{0, 0} Score{0, 0} Score{0, 0} Score{0, 0} Score{19, -46} Score{17, 0} Score{24, -11} Score{7, 16} Score{21, 7} Score{69, 2} Score{23, 20} Score{51, 44} Score{22, 15} Score{48, 11} Score{31, 27} Score{57, 22} Score{8, 21} Score{8, 26} Score{41, 39} Score{57, 36} Score{-42, 66} Score{35, 41} Score{72, 69} Score{97, 78} Score{0, 263} Score{127, 51} Score{147, 0} Score{0, 128} Score{0, 0} Score{0, 0} Score{0, 0} Score{0, 0} Score{0, 0} Score{-49, -7} Score{-124, 7} Score{-82, -11} Score{-69, -21} Score{-29, -37} Score{21, -65} Score{-37, -45} Score{0, 0} Score{19, 1} Score{38, -2} Score{57, 1} Score{49, 2} Score{41, 6} Score{48, 3} Score{18, 22} Score{0, 0} Score{63, 100} Score{28, 122} Score{50, 124} Score{16, 137} Score{39, 142} Score{-4, 163} Score{-91, -4} Score{-79, -1} Score{-42, -24} Score{-28, -24} Score{-50, -22} Score{-57, -14} Score{-61, -3} Score{-62, -14} Score{-21, -16} Score{-26, -41} Score{21, -8} Score{-27, -2} Score{11, 0} Score{-11, 0}]'
REGEXP = /Score{(-?\d+, -?\d+)}/

EVAL_PATH = File.join(__dir__, '..', 'evaluation', 'eval.go')

values = RESULT.scan(REGEXP).flatten

eval_content = IO.read(EVAL_PATH)

idx = -1

new_eval = eval_content.gsub(REGEXP) do |match|
  idx += 1
  "Score{#{values[idx]}}"
end

IO.write(EVAL_PATH, new_eval, mode: 'w')
