# frozen_string_literal: true

require 'parallel'
require 'pgn'

@q = Queue.new
reader = Thread.new do
  acc = +''
  IO.foreach('./games.pgn') do |line|
    if !acc.empty? && line.include?('Event')
      @q << acc
      acc = +line
    else
      acc << line
    end
  end
  @q << acc
  @q << Parallel::Stop
  @q.close
end

@result = File.new('./games.fen', 'a')

Parallel.each(@q, in_processes: 7) do |pgn|
  game = PGN.parse(pgn).first
  game.moves.each_with_index do |move, idx|
    next if idx.zero? # do not include position from book
    break if move.comment.include? 'M' # stop if mate was found in position

    @result.write(game.positions[idx].to_fen.to_s + ";#{game.result}\n")
    @result.flush
  end
end
reader.join
