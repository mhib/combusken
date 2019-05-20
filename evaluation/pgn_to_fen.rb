# frozen_string_literal: true

require 'parallel'
require 'thread'
require 'pgn'
require 'pry'

@q = Queue.new
reader = Thread.new do |n|
  acc = +''
  IO.foreach("./games.pgn") do |line|
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
  mate_found = false
  game = PGN.parse(pgn).first
  game.moves.each_with_index do |move, idx|
    if move.comment.include? 'M'
      mate_found = true
      next
    end
    @result.write(game.positions[idx + 1].to_fen.to_s + ";#{game.result}\n")
    @result.flush
  end
  next if mate_found
end
reader.join
