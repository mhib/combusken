package evaluation

import (
	. "github.com/mhib/combusken/backend"
	. "github.com/mhib/combusken/utils"
)

const tuning = false

var T Trace

const PawnPhase = 0
const KnightPhase = 1
const BishopPhase = 1
const RookPhase = 2
const QueenPhase = 4
const TotalPhase = PawnPhase*16 + KnightPhase*4 + BishopPhase*4 + RookPhase*4 + QueenPhase*2

var PawnValue = S(82, 150)
var KnightValue = S(420, 537)
var BishopValue = S(393, 528)
var RookValue = S(570, 864)
var QueenValue = S(1416, 1558)

// Pawns Square scores
// Bishop Flag, Rank, Row
var PawnScores = [16][7][8]Score{
	{
		{},
		{S(-24, 17), S(-17, 6), S(-24, 15), S(-10, 17), S(-23, 25), S(17, 13), S(41, 4), S(-6, 11)},
		{S(-22, 6), S(-22, 7), S(-21, 9), S(-12, 8), S(-3, 3), S(1, 7), S(12, 5), S(-12, 10)},
		{S(-21, 17), S(-26, 18), S(-2, 6), S(3, 0), S(15, -7), S(19, 1), S(8, 14), S(-9, 17)},
		{S(-16, 30), S(-9, 21), S(-8, 11), S(0, -1), S(31, -11), S(36, 2), S(40, 10), S(8, 23)},
		{S(-11, 64), S(9, 48), S(19, 38), S(35, 9), S(64, -4), S(74, 15), S(58, 40), S(3, 63)},
		{S(-8, 9), S(-6, 29), S(5, 30), S(27, 9), S(31, -3), S(31, 15), S(-89, 54), S(-128, 62)},
	},
	{
		{},
		{S(-15, 24), S(-9, 2), S(-21, 9), S(2, 16), S(-18, 16), S(25, 17), S(35, -6), S(-8, 11)},
		{S(-12, 7), S(-18, 0), S(-3, 7), S(-5, -1), S(11, 14), S(-8, 1), S(20, 10), S(-21, 10)},
		{S(-8, 16), S(-11, 11), S(6, -3), S(28, 1), S(25, -12), S(27, 6), S(7, 11), S(-10, 10)},
		{S(0, 30), S(14, -1), S(12, 15), S(18, -20), S(46, -8), S(43, -8), S(37, 12), S(-7, 26)},
		{S(-4, 60), S(23, 46), S(29, 11), S(52, 3), S(66, -28), S(86, 6), S(51, 9), S(5, 54)},
		{S(-4, 28), S(-9, -1), S(4, 6), S(23, -22), S(35, -15), S(28, -14), S(-90, 30), S(-130, 34)},
	},
	{
		{},
		{S(-16, 13), S(-9, 10), S(-10, 19), S(-16, 5), S(8, 30), S(14, 5), S(60, 1), S(-1, 17)},
		{S(-18, 6), S(-16, 6), S(-9, -1), S(4, 8), S(6, -8), S(8, 16), S(6, -2), S(-10, 17)},
		{S(-9, 2), S(-6, 10), S(14, 6), S(21, -8), S(44, -5), S(27, -4), S(17, 12), S(-4, 18)},
		{S(-5, 17), S(9, 11), S(6, 0), S(28, 0), S(45, -16), S(56, 4), S(34, 5), S(4, 31)},
		{S(-3, 59), S(16, 19), S(38, 24), S(39, -25), S(77, -8), S(84, 2), S(53, 24), S(11, 54)},
		{S(-11, -4), S(-6, 18), S(-1, -15), S(26, -9), S(32, -27), S(30, -9), S(-91, 26), S(-128, 46)},
	},
	{
		{},
		{S(-17, 7), S(-6, -9), S(-4, 3), S(-4, 9), S(-5, 19), S(13, 7), S(44, 5), S(-2, 0)},
		{S(-13, -1), S(-23, -7), S(-1, -3), S(-1, -9), S(15, 2), S(-7, 6), S(10, -1), S(-12, -3)},
		{S(-7, 1), S(-7, 3), S(11, -16), S(38, -6), S(35, -11), S(24, 3), S(11, 3), S(-11, 6)},
		{S(-8, 19), S(9, 2), S(9, 0), S(22, -18), S(42, -19), S(39, -8), S(32, 8), S(-6, 17)},
		{S(-3, 62), S(17, 34), S(31, 19), S(42, -9), S(67, -19), S(81, 4), S(54, 21), S(3, 51)},
		{S(-6, 14), S(-8, 10), S(3, 5), S(26, -14), S(33, -18), S(30, -9), S(-90, 30), S(-129, 44)},
	},
	{
		{},
		{S(-17, 14), S(-6, 13), S(-2, 7), S(0, 16), S(-16, 19), S(11, 16), S(42, 5), S(-3, 12)},
		{S(-12, 2), S(-18, 2), S(-5, -13), S(6, -4), S(-6, -7), S(-2, 12), S(2, 3), S(-12, 15)},
		{S(-11, 18), S(1, 9), S(16, -5), S(16, -24), S(36, -22), S(13, -5), S(10, 19), S(-10, 25)},
		{S(1, 41), S(17, 13), S(5, -5), S(29, -12), S(37, -34), S(50, 3), S(36, 9), S(-1, 43)},
		{S(-8, 76), S(18, 41), S(29, 22), S(47, -14), S(69, -19), S(84, -3), S(54, 25), S(6, 65)},
		{S(2, 35), S(-8, 11), S(5, 7), S(29, -9), S(33, -22), S(32, -9), S(-89, 32), S(-126, 50)},
	},
	{
		{},
		{S(-4, 21), S(1, 8), S(-8, 13), S(-5, 14), S(-10, 19), S(16, 26), S(29, -1), S(1, 11)},
		{S(-6, 10), S(-12, -2), S(0, 6), S(4, -8), S(5, 9), S(-2, 5), S(5, 14), S(-11, 9)},
		{S(-3, 23), S(5, 13), S(18, -6), S(29, -8), S(35, -13), S(24, 6), S(8, 18), S(-9, 18)},
		{S(13, 36), S(17, 6), S(8, 14), S(31, -18), S(50, -20), S(46, -6), S(41, 19), S(-1, 30)},
		{S(-1, 72), S(22, 44), S(33, 19), S(48, -5), S(69, -31), S(82, 5), S(50, 19), S(6, 56)},
		{S(-1, 24), S(-11, 2), S(5, 8), S(28, -11), S(33, -18), S(30, -9), S(-89, 30), S(-130, 42)},
	},
	{
		{},
		{S(-9, -3), S(0, 0), S(-2, 5), S(-6, 5), S(1, 21), S(0, 0), S(50, 11), S(-2, 3)},
		{S(-9, -5), S(-10, -1), S(-14, -13), S(0, -3), S(-9, -15), S(-4, 17), S(-11, -10), S(-8, 9)},
		{S(-6, -2), S(1, -10), S(19, -7), S(12, -34), S(32, -8), S(6, -5), S(13, 14), S(-8, 13)},
		{S(0, 14), S(10, 8), S(6, -18), S(21, -5), S(35, -41), S(58, 2), S(26, 1), S(5, 26)},
		{S(-13, 59), S(18, 24), S(27, 19), S(32, -22), S(72, -18), S(82, -1), S(53, 21), S(11, 53)},
		{S(-7, 7), S(-8, 12), S(-5, -11), S(27, -10), S(32, -27), S(30, -9), S(-93, 19), S(-128, 45)},
	},
	{
		{},
		{S(-8, 4), S(-3, 6), S(-3, 4), S(-7, 9), S(3, 23), S(6, 14), S(38, 7), S(-2, 10)},
		{S(-9, -2), S(-17, -3), S(-6, 1), S(2, -2), S(3, 1), S(-6, 16), S(1, -1), S(-10, -3)},
		{S(-5, 6), S(-5, 9), S(19, -10), S(27, -11), S(34, -6), S(12, 8), S(9, 7), S(-13, 16)},
		{S(4, 26), S(15, 9), S(7, 7), S(28, -11), S(41, -21), S(49, -3), S(23, 7), S(-12, 22)},
		{S(-8, 60), S(16, 36), S(29, 20), S(45, -9), S(69, -20), S(85, 6), S(49, 18), S(8, 54)},
		{S(-6, 11), S(-8, 10), S(1, 3), S(26, -12), S(31, -20), S(30, -9), S(-90, 30), S(-129, 44)},
	},
	{
		{},
		{S(-15, 18), S(-9, 5), S(-15, 11), S(-8, 7), S(-12, 14), S(15, 8), S(27, 5), S(-4, 12)},
		{S(-14, 7), S(-18, -1), S(-10, -1), S(-6, -11), S(3, -5), S(-14, 4), S(15, 7), S(-12, 10)},
		{S(-10, 23), S(-7, 18), S(1, -11), S(19, -19), S(22, -27), S(16, 1), S(1, 15), S(-9, 20)},
		{S(-1, 49), S(10, 8), S(-2, 6), S(7, -28), S(40, -9), S(26, -16), S(34, 17), S(-6, 33)},
		{S(9, 74), S(17, 43), S(42, 17), S(39, -12), S(68, -22), S(86, 12), S(64, 27), S(3, 55)},
		{S(-3, 29), S(-8, 11), S(7, 13), S(25, -23), S(35, -13), S(30, -12), S(-90, 30), S(-129, 45)},
	},
	{
		{},
		{S(-10, 13), S(-10, -4), S(-22, 8), S(5, 14), S(-20, 12), S(24, 15), S(21, -9), S(-4, 2)},
		{S(-13, 3), S(-19, -11), S(3, -1), S(-10, -10), S(11, 11), S(-18, -1), S(18, 7), S(-15, -8)},
		{S(-7, 20), S(-13, 5), S(3, -23), S(28, -5), S(19, -30), S(29, 4), S(1, 2), S(-13, 11)},
		{S(5, 31), S(11, -9), S(5, 3), S(2, -43), S(45, -14), S(24, -17), S(34, 11), S(-9, 20)},
		{S(7, 64), S(17, 37), S(30, 11), S(46, -1), S(65, -23), S(88, 5), S(51, 16), S(4, 52)},
		{S(-4, 19), S(-11, 1), S(5, 10), S(21, -25), S(33, -18), S(28, -16), S(-90, 30), S(-132, 37)},
	},
	{
		{},
		{S(-7, 17), S(-2, 10), S(-8, 23), S(-10, 12), S(-12, 28), S(13, 8), S(37, 10), S(0, 6)},
		{S(-10, 13), S(-7, 6), S(-8, 9), S(2, 5), S(3, -10), S(-2, 14), S(14, 2), S(-9, 6)},
		{S(-2, 16), S(3, 17), S(16, 3), S(20, -4), S(33, -5), S(25, 0), S(20, 11), S(-7, 15)},
		{S(0, 36), S(23, 18), S(9, 7), S(22, -7), S(44, -23), S(56, -2), S(38, 4), S(2, 29)},
		{S(5, 74), S(20, 38), S(35, 24), S(34, -24), S(77, -14), S(85, 8), S(64, 23), S(8, 51)},
		{S(-7, 11), S(-5, 17), S(-1, -5), S(26, -9), S(35, -19), S(30, -12), S(-91, 27), S(-130, 44)},
	},
	{
		{},
		{S(-11, 14), S(-3, 2), S(-9, 10), S(2, 13), S(-14, 22), S(16, 14), S(31, -2), S(-5, 1)},
		{S(-9, -1), S(-14, 1), S(-3, 8), S(-2, 2), S(6, 7), S(-9, 6), S(14, 3), S(-17, 2)},
		{S(-1, 15), S(2, 8), S(13, -7), S(29, 1), S(29, -11), S(29, 4), S(10, 6), S(-19, 11)},
		{S(9, 22), S(21, 14), S(14, 9), S(24, -29), S(40, -19), S(47, -6), S(26, 9), S(-17, 26)},
		{S(0, 62), S(21, 38), S(31, 17), S(40, -13), S(67, -20), S(81, 4), S(55, 24), S(2, 52)},
		{S(-8, 12), S(-8, 10), S(3, 5), S(26, -12), S(33, -18), S(29, -10), S(-90, 30), S(-129, 43)},
	},
	{
		{},
		{S(-19, 20), S(-7, 11), S(-8, 12), S(-2, 17), S(-28, 13), S(3, 11), S(19, 4), S(-3, 1)},
		{S(-12, 15), S(-21, 3), S(-18, 3), S(-3, -8), S(-4, -4), S(-7, 1), S(2, 3), S(-10, 1)},
		{S(-7, 29), S(-5, 11), S(10, -23), S(11, -17), S(22, -38), S(19, -6), S(2, 10), S(-6, 8)},
		{S(6, 30), S(15, 4), S(5, 6), S(21, -31), S(37, -16), S(30, -18), S(33, 14), S(1, 28)},
		{S(5, 72), S(21, 42), S(30, 19), S(46, -5), S(60, -23), S(84, 4), S(53, 21), S(6, 58)},
		{S(-5, 18), S(-9, 8), S(3, 7), S(22, -20), S(32, -18), S(30, -11), S(-90, 29), S(-129, 41)},
	},
	{
		{},
		{S(-9, 19), S(-2, 8), S(-13, 10), S(-5, 9), S(-20, 16), S(10, 28), S(12, 0), S(-4, 4)},
		{S(-9, 12), S(-23, -1), S(-1, 6), S(-1, -16), S(4, 8), S(-5, 0), S(4, 7), S(-15, -2)},
		{S(-1, 21), S(-4, 17), S(12, -14), S(28, -8), S(27, -26), S(23, 5), S(3, 5), S(-11, 10)},
		{S(15, 36), S(13, 7), S(11, 10), S(23, -34), S(47, -15), S(27, -16), S(40, 14), S(-7, 28)},
		{S(11, 65), S(23, 41), S(31, 21), S(44, -7), S(54, -19), S(82, 4), S(52, 21), S(7, 54)},
		{S(-4, 16), S(-8, 9), S(3, 5), S(26, -14), S(33, -18), S(28, -11), S(-90, 30), S(-129, 44)},
	},
	{
		{},
		{S(-15, 12), S(-1, 10), S(-6, 21), S(-5, 11), S(-19, 19), S(-2, 5), S(27, 17), S(0, 8)},
		{S(-14, 0), S(-9, 9), S(-11, -8), S(1, 1), S(-11, -11), S(-7, 17), S(3, -7), S(-8, 9)},
		{S(-9, 13), S(1, 14), S(16, 1), S(15, -22), S(35, -17), S(12, -8), S(13, 7), S(-11, 9)},
		{S(0, 33), S(31, 16), S(0, 4), S(28, -12), S(33, -35), S(55, -9), S(30, 3), S(-2, 33)},
		{S(9, 75), S(17, 34), S(36, 24), S(31, -23), S(68, -17), S(79, 1), S(49, 22), S(5, 51)},
		{S(-9, 11), S(-9, 11), S(0, 2), S(25, -12), S(32, -20), S(31, -8), S(-89, 29), S(-130, 44)},
	},
	{
		{},
		{S(-8, 8), S(0, 12), S(-11, 14), S(-8, 9), S(-18, 8), S(4, 18), S(17, 5), S(-5, 4)},
		{S(-8, 3), S(-14, 2), S(-5, 1), S(-1, -3), S(0, 0), S(-12, 21), S(-2, 2), S(-16, 1)},
		{S(-1, 13), S(1, 13), S(15, -2), S(23, -8), S(29, -18), S(16, 2), S(3, 11), S(-20, 13)},
		{S(10, 35), S(26, 12), S(5, 16), S(27, -26), S(37, -21), S(44, -10), S(24, 8), S(-13, 24)},
		{S(8, 67), S(21, 39), S(30, 22), S(44, -10), S(67, -20), S(75, 2), S(46, 20), S(0, 53)},
		{S(-12, 14), S(-12, 8), S(1, 5), S(25, -13), S(33, -18), S(30, -9), S(-89, 30), S(-129, 44)},
	},
}

// Piece Square Values
var PieceScores = [Queen + 1][8][4]Score{
	{},
	{ // knight
		{S(-93, -40), S(-25, -44), S(-30, -34), S(-22, -11)},
		{S(-22, -23), S(-24, -20), S(-17, -36), S(-18, -16)},
		{S(-17, -49), S(-4, -26), S(-12, -17), S(1, 3)},
		{S(-10, -6), S(14, -5), S(4, 17), S(3, 24)},
		{S(7, -4), S(4, -1), S(21, 20), S(14, 35)},
		{S(-26, -14), S(-17, 1), S(6, 27), S(18, 24)},
		{S(6, -32), S(-12, -13), S(32, -30), S(41, 4)},
		{S(-230, -36), S(-92, -11), S(-150, 24), S(-11, -13)},
	},
	{ // Bishop
		{S(24, -37), S(29, -3), S(9, -5), S(6, 6)},
		{S(30, -49), S(21, -39), S(31, -11), S(14, 0)},
		{S(22, -10), S(40, -2), S(10, -13), S(27, 15)},
		{S(14, -12), S(21, 3), S(23, 17), S(35, 20)},
		{S(-7, 8), S(23, 11), S(14, 19), S(37, 27)},
		{S(-3, 9), S(-5, 25), S(7, 3), S(23, 22)},
		{S(-37, 27), S(-20, 7), S(14, 20), S(-1, 28)},
		{S(-44, 4), S(-39, 29), S(-113, 42), S(-97, 47)},
	},
	{ // Rook
		{S(-17, -26), S(-8, -17), S(-2, -15), S(2, -25)},
		{S(-58, -11), S(-13, -31), S(-6, -29), S(-15, -26)},
		{S(-33, -19), S(-8, -11), S(-18, -11), S(-15, -15)},
		{S(-26, 1), S(-15, 12), S(-13, 13), S(-9, 6)},
		{S(-12, 14), S(14, 12), S(24, 13), S(33, 9)},
		{S(-12, 24), S(48, 7), S(45, 20), S(57, 5)},
		{S(0, 31), S(5, 34), S(40, 25), S(43, 34)},
		{S(43, 32), S(64, 31), S(39, 38), S(37, 34)},
	},
	{ // Queen
		{S(9, -105), S(10, -90), S(11, -95), S(16, -61)},
		{S(7, -86), S(16, -79), S(21, -89), S(10, -49)},
		{S(3, -65), S(15, -31), S(6, 2), S(1, -10)},
		{S(1, -23), S(9, 5), S(2, 31), S(-13, 60)},
		{S(14, -12), S(6, 38), S(-9, 51), S(-20, 84)},
		{S(7, 4), S(16, 0), S(-8, 55), S(4, 50)},
		{S(-27, 10), S(-72, 54), S(-14, 51), S(-37, 84)},
		{S(-6, 10), S(28, 4), S(25, 35), S(32, 38)},
	},
}

var KingScores = [16][8][8]Score{
	{
		{S(168, -23), S(154, 20), S(56, 61), S(62, 84), S(58, 85), S(54, 60), S(153, 20), S(166, -22)},
		{S(166, 34), S(132, 40), S(56, 76), S(24, 131), S(21, 133), S(55, 75), S(136, 45), S(163, 39)},
		{S(107, 31), S(155, 31), S(113, 63), S(79, 131), S(79, 134), S(112, 63), S(153, 33), S(106, 34)},
		{S(87, 23), S(237, 24), S(156, 70), S(96, 142), S(95, 143), S(155, 68), S(236, 25), S(87, 26)},
		{S(89, 45), S(270, 37), S(181, 81), S(108, 143), S(107, 142), S(181, 80), S(269, 37), S(89, 44)},
		{S(117, 43), S(298, 58), S(248, 80), S(184, 121), S(184, 120), S(248, 78), S(298, 59), S(117, 43)},
		{S(115, 0), S(211, 65), S(189, 73), S(149, 113), S(149, 113), S(189, 73), S(211, 66), S(115, 0)},
		{S(157, -98), S(290, -17), S(173, 18), S(161, 75), S(161, 75), S(173, 18), S(290, -17), S(157, -98)},
	},
	{
		{S(167, -25), S(149, 17), S(53, 58), S(62, 81), S(59, 79), S(54, 60), S(149, 14), S(167, -25)},
		{S(165, 33), S(130, 39), S(56, 76), S(24, 136), S(25, 139), S(55, 74), S(136, 44), S(164, 30)},
		{S(107, 32), S(154, 31), S(113, 63), S(81, 138), S(80, 135), S(114, 66), S(154, 28), S(108, 33)},
		{S(87, 24), S(237, 25), S(156, 71), S(96, 145), S(97, 147), S(156, 70), S(238, 26), S(87, 24)},
		{S(89, 44), S(270, 37), S(181, 82), S(108, 145), S(108, 144), S(181, 83), S(270, 37), S(89, 44)},
		{S(117, 43), S(298, 58), S(248, 80), S(184, 121), S(184, 121), S(248, 80), S(298, 57), S(117, 43)},
		{S(115, 0), S(211, 66), S(189, 73), S(149, 113), S(149, 113), S(189, 73), S(211, 66), S(115, 0)},
		{S(157, -98), S(290, -17), S(173, 18), S(161, 75), S(161, 75), S(173, 18), S(290, -17), S(157, -98)},
	},
	{
		{S(167, -25), S(151, 19), S(56, 60), S(61, 80), S(59, 80), S(54, 59), S(151, 18), S(167, -25)},
		{S(164, 33), S(131, 41), S(55, 75), S(25, 136), S(24, 136), S(56, 76), S(129, 37), S(166, 35)},
		{S(107, 32), S(154, 30), S(113, 62), S(80, 136), S(80, 137), S(113, 63), S(155, 32), S(107, 31)},
		{S(87, 24), S(237, 25), S(156, 71), S(96, 146), S(96, 145), S(156, 69), S(237, 24), S(87, 24)},
		{S(89, 44), S(270, 37), S(181, 83), S(108, 144), S(108, 144), S(181, 82), S(270, 37), S(89, 43)},
		{S(117, 43), S(298, 58), S(248, 80), S(184, 121), S(184, 121), S(248, 80), S(298, 58), S(117, 43)},
		{S(115, 0), S(211, 66), S(189, 73), S(149, 113), S(149, 113), S(189, 73), S(211, 66), S(115, 0)},
		{S(157, -98), S(290, -17), S(173, 18), S(161, 75), S(161, 75), S(173, 18), S(290, -17), S(157, -98)},
	},
	{
		{S(167, -25), S(150, 19), S(55, 60), S(62, 80), S(61, 80), S(55, 60), S(148, 18), S(168, -25)},
		{S(165, 33), S(131, 40), S(56, 76), S(24, 136), S(24, 136), S(56, 76), S(132, 40), S(166, 33)},
		{S(107, 32), S(154, 30), S(113, 63), S(80, 136), S(80, 136), S(113, 64), S(154, 30), S(107, 32)},
		{S(87, 24), S(237, 25), S(156, 70), S(96, 145), S(96, 145), S(156, 70), S(237, 25), S(87, 24)},
		{S(89, 44), S(270, 37), S(181, 82), S(108, 144), S(108, 144), S(181, 82), S(270, 37), S(89, 44)},
		{S(117, 43), S(298, 58), S(248, 80), S(184, 121), S(184, 121), S(248, 80), S(298, 58), S(117, 43)},
		{S(115, 0), S(211, 66), S(189, 73), S(149, 113), S(149, 113), S(189, 73), S(211, 66), S(115, 0)},
		{S(157, -98), S(290, -17), S(173, 18), S(161, 75), S(161, 75), S(173, 18), S(290, -17), S(157, -98)},
	},
	{
		{S(167, -25), S(154, 22), S(57, 61), S(62, 81), S(61, 80), S(54, 60), S(148, 16), S(166, -26)},
		{S(165, 32), S(131, 40), S(55, 75), S(24, 135), S(24, 137), S(56, 77), S(130, 37), S(162, 32)},
		{S(107, 32), S(154, 30), S(113, 63), S(80, 136), S(80, 136), S(113, 64), S(153, 29), S(107, 32)},
		{S(87, 24), S(237, 25), S(156, 70), S(96, 146), S(96, 145), S(156, 70), S(237, 25), S(87, 23)},
		{S(89, 44), S(270, 37), S(181, 82), S(108, 144), S(108, 144), S(181, 83), S(270, 37), S(89, 44)},
		{S(117, 43), S(298, 58), S(248, 80), S(184, 121), S(184, 121), S(248, 80), S(298, 59), S(117, 43)},
		{S(115, 0), S(211, 66), S(189, 73), S(149, 113), S(149, 113), S(189, 73), S(211, 66), S(115, 0)},
		{S(157, -98), S(290, -17), S(173, 18), S(161, 75), S(161, 75), S(173, 18), S(290, -17), S(157, -98)},
	},
	{
		{S(167, -25), S(149, 19), S(53, 58), S(62, 81), S(60, 79), S(56, 61), S(142, 11), S(165, -27)},
		{S(166, 33), S(130, 39), S(56, 76), S(24, 135), S(26, 139), S(56, 76), S(131, 40), S(164, 31)},
		{S(107, 32), S(154, 30), S(113, 63), S(81, 139), S(79, 135), S(114, 65), S(154, 29), S(108, 33)},
		{S(87, 24), S(237, 25), S(156, 71), S(96, 145), S(97, 148), S(156, 70), S(237, 26), S(87, 24)},
		{S(89, 44), S(270, 38), S(181, 82), S(108, 145), S(108, 144), S(181, 83), S(270, 37), S(89, 44)},
		{S(117, 43), S(298, 58), S(248, 80), S(184, 121), S(184, 121), S(248, 80), S(298, 58), S(117, 43)},
		{S(115, 0), S(211, 66), S(189, 73), S(149, 113), S(149, 113), S(189, 73), S(211, 66), S(115, 0)},
		{S(157, -98), S(290, -17), S(173, 18), S(161, 75), S(161, 75), S(173, 18), S(290, -17), S(157, -98)},
	},
	{
		{S(167, -25), S(151, 20), S(57, 61), S(62, 80), S(61, 80), S(54, 59), S(151, 21), S(166, -25)},
		{S(165, 33), S(132, 41), S(56, 75), S(24, 136), S(24, 136), S(56, 75), S(129, 38), S(164, 32)},
		{S(107, 32), S(154, 30), S(113, 63), S(80, 136), S(80, 137), S(113, 64), S(154, 30), S(107, 32)},
		{S(87, 24), S(237, 25), S(156, 70), S(96, 145), S(96, 145), S(156, 70), S(237, 25), S(87, 24)},
		{S(89, 44), S(270, 37), S(181, 82), S(108, 144), S(108, 144), S(181, 82), S(270, 37), S(89, 44)},
		{S(117, 43), S(298, 58), S(248, 80), S(184, 121), S(184, 121), S(248, 80), S(298, 58), S(117, 43)},
		{S(115, 0), S(211, 66), S(189, 73), S(149, 113), S(149, 113), S(189, 73), S(211, 66), S(115, 0)},
		{S(157, -98), S(290, -17), S(173, 18), S(161, 75), S(161, 75), S(173, 18), S(290, -17), S(157, -98)},
	},
	{
		{S(167, -25), S(150, 19), S(56, 60), S(61, 80), S(65, 81), S(55, 60), S(144, 17), S(165, -25)},
		{S(165, 33), S(130, 40), S(56, 77), S(24, 137), S(26, 139), S(56, 77), S(128, 38), S(165, 32)},
		{S(107, 32), S(154, 30), S(113, 63), S(80, 136), S(80, 136), S(113, 63), S(154, 30), S(107, 32)},
		{S(87, 24), S(237, 25), S(156, 70), S(96, 145), S(96, 145), S(156, 70), S(237, 25), S(87, 24)},
		{S(89, 44), S(270, 37), S(181, 82), S(108, 144), S(108, 144), S(181, 82), S(270, 37), S(89, 44)},
		{S(117, 43), S(298, 58), S(248, 80), S(184, 121), S(184, 121), S(248, 80), S(298, 58), S(117, 43)},
		{S(115, 0), S(211, 66), S(189, 73), S(149, 113), S(149, 113), S(189, 73), S(211, 66), S(115, 0)},
		{S(157, -98), S(290, -17), S(173, 18), S(161, 75), S(161, 75), S(173, 18), S(290, -17), S(157, -98)},
	},
	{
		{S(167, -25), S(151, 19), S(56, 60), S(62, 80), S(58, 80), S(53, 60), S(155, 21), S(168, -23)},
		{S(165, 33), S(132, 41), S(56, 76), S(23, 134), S(23, 135), S(56, 77), S(132, 41), S(162, 32)},
		{S(107, 32), S(154, 31), S(113, 63), S(80, 136), S(80, 136), S(113, 63), S(154, 29), S(106, 31)},
		{S(87, 24), S(237, 25), S(156, 70), S(96, 146), S(96, 144), S(156, 70), S(237, 25), S(87, 24)},
		{S(89, 44), S(270, 38), S(181, 82), S(108, 144), S(108, 144), S(181, 82), S(270, 36), S(89, 44)},
		{S(117, 43), S(298, 58), S(248, 80), S(184, 121), S(184, 121), S(248, 80), S(298, 58), S(117, 43)},
		{S(115, 0), S(211, 66), S(189, 73), S(149, 113), S(149, 113), S(189, 73), S(211, 66), S(115, 0)},
		{S(157, -98), S(290, -17), S(173, 18), S(161, 75), S(161, 75), S(173, 18), S(290, -17), S(157, -98)},
	},
	{
		{S(167, -25), S(149, 19), S(54, 59), S(61, 79), S(60, 79), S(54, 60), S(152, 21), S(169, -23)},
		{S(165, 33), S(130, 39), S(56, 76), S(24, 136), S(24, 136), S(56, 76), S(134, 43), S(164, 31)},
		{S(107, 32), S(154, 30), S(113, 63), S(80, 136), S(80, 136), S(113, 65), S(154, 29), S(107, 33)},
		{S(87, 24), S(237, 25), S(156, 70), S(96, 145), S(96, 144), S(156, 69), S(237, 25), S(87, 24)},
		{S(89, 44), S(270, 37), S(181, 82), S(108, 144), S(108, 144), S(181, 82), S(270, 37), S(89, 44)},
		{S(117, 43), S(298, 58), S(248, 80), S(184, 121), S(184, 121), S(248, 80), S(298, 58), S(117, 43)},
		{S(115, 0), S(211, 66), S(189, 73), S(149, 113), S(149, 113), S(189, 73), S(211, 66), S(115, 0)},
		{S(157, -98), S(290, -17), S(173, 18), S(161, 75), S(161, 75), S(173, 18), S(290, -17), S(157, -98)},
	},
	{
		{S(167, -25), S(147, 18), S(58, 61), S(61, 80), S(62, 79), S(54, 59), S(154, 22), S(165, -26)},
		{S(165, 33), S(130, 39), S(55, 74), S(25, 138), S(24, 136), S(57, 78), S(130, 37), S(168, 34)},
		{S(107, 32), S(154, 30), S(113, 64), S(80, 135), S(80, 138), S(113, 63), S(154, 30), S(107, 32)},
		{S(87, 24), S(237, 25), S(156, 70), S(96, 147), S(96, 145), S(156, 72), S(237, 25), S(87, 24)},
		{S(89, 44), S(270, 37), S(181, 83), S(108, 144), S(108, 145), S(181, 82), S(270, 38), S(89, 44)},
		{S(117, 43), S(298, 58), S(248, 80), S(184, 121), S(184, 121), S(248, 80), S(298, 58), S(117, 43)},
		{S(115, 0), S(211, 66), S(189, 73), S(149, 113), S(149, 113), S(189, 73), S(211, 66), S(115, 0)},
		{S(157, -98), S(290, -17), S(173, 18), S(161, 75), S(161, 75), S(173, 18), S(290, -17), S(157, -98)},
	},
	{
		{S(167, -25), S(144, 17), S(53, 59), S(62, 80), S(62, 78), S(55, 60), S(153, 22), S(168, -25)},
		{S(165, 33), S(130, 39), S(56, 76), S(25, 136), S(24, 136), S(57, 76), S(132, 39), S(168, 34)},
		{S(107, 32), S(154, 30), S(113, 63), S(80, 136), S(80, 136), S(113, 64), S(154, 30), S(107, 32)},
		{S(87, 24), S(237, 25), S(156, 70), S(96, 145), S(96, 145), S(156, 70), S(237, 25), S(87, 24)},
		{S(89, 44), S(270, 37), S(181, 82), S(108, 144), S(108, 144), S(181, 82), S(270, 37), S(89, 44)},
		{S(117, 43), S(298, 58), S(248, 80), S(184, 121), S(184, 121), S(248, 80), S(298, 58), S(117, 43)},
		{S(115, 0), S(211, 66), S(189, 73), S(149, 113), S(149, 113), S(189, 73), S(211, 66), S(115, 0)},
		{S(157, -98), S(290, -17), S(173, 18), S(161, 75), S(161, 75), S(173, 18), S(290, -17), S(157, -98)},
	},
	{
		{S(167, -25), S(150, 20), S(56, 60), S(62, 80), S(60, 79), S(55, 60), S(151, 19), S(167, -25)},
		{S(165, 33), S(131, 40), S(56, 76), S(24, 136), S(24, 136), S(56, 76), S(130, 40), S(164, 33)},
		{S(107, 32), S(154, 30), S(113, 63), S(80, 136), S(80, 136), S(113, 64), S(154, 30), S(107, 32)},
		{S(87, 24), S(237, 25), S(156, 70), S(96, 145), S(96, 145), S(156, 70), S(237, 25), S(87, 24)},
		{S(89, 44), S(270, 37), S(181, 82), S(108, 144), S(108, 144), S(181, 82), S(270, 37), S(89, 44)},
		{S(117, 43), S(298, 58), S(248, 80), S(184, 121), S(184, 121), S(248, 80), S(298, 58), S(117, 43)},
		{S(115, 0), S(211, 66), S(189, 73), S(149, 113), S(149, 113), S(189, 73), S(211, 66), S(115, 0)},
		{S(157, -98), S(290, -17), S(173, 18), S(161, 75), S(161, 75), S(173, 18), S(290, -17), S(157, -98)},
	},
	{
		{S(167, -25), S(150, 19), S(55, 60), S(62, 80), S(63, 80), S(56, 61), S(150, 17), S(169, -25)},
		{S(165, 33), S(131, 40), S(56, 77), S(24, 136), S(23, 137), S(56, 77), S(130, 40), S(161, 31)},
		{S(107, 32), S(154, 30), S(113, 63), S(80, 136), S(80, 136), S(113, 64), S(154, 30), S(107, 32)},
		{S(87, 24), S(237, 25), S(156, 70), S(96, 145), S(96, 145), S(156, 70), S(237, 25), S(87, 24)},
		{S(89, 44), S(270, 37), S(181, 82), S(108, 144), S(108, 144), S(181, 82), S(270, 37), S(89, 44)},
		{S(117, 43), S(298, 58), S(248, 80), S(184, 121), S(184, 121), S(248, 80), S(298, 58), S(117, 43)},
		{S(115, 0), S(211, 66), S(189, 73), S(149, 113), S(149, 113), S(189, 73), S(211, 66), S(115, 0)},
		{S(157, -98), S(290, -17), S(173, 18), S(161, 75), S(161, 75), S(173, 18), S(290, -17), S(157, -98)},
	},
	{
		{S(167, -25), S(150, 19), S(58, 60), S(62, 80), S(65, 79), S(55, 60), S(148, 22), S(165, -26)},
		{S(165, 33), S(131, 40), S(55, 76), S(24, 136), S(23, 136), S(57, 77), S(128, 38), S(167, 33)},
		{S(107, 32), S(154, 30), S(113, 63), S(80, 136), S(80, 136), S(113, 64), S(154, 30), S(107, 32)},
		{S(87, 24), S(237, 25), S(156, 70), S(96, 145), S(96, 145), S(156, 70), S(237, 25), S(87, 24)},
		{S(89, 44), S(270, 37), S(181, 82), S(108, 144), S(108, 144), S(181, 82), S(270, 37), S(89, 44)},
		{S(117, 43), S(298, 58), S(248, 80), S(184, 121), S(184, 121), S(248, 80), S(298, 58), S(117, 43)},
		{S(115, 0), S(211, 66), S(189, 73), S(149, 113), S(149, 113), S(189, 73), S(211, 66), S(115, 0)},
		{S(157, -98), S(290, -17), S(173, 18), S(161, 75), S(161, 75), S(173, 18), S(290, -17), S(157, -98)},
	},
	{
		{S(166, -25), S(138, 18), S(53, 59), S(63, 81), S(74, 77), S(57, 61), S(143, 24), S(168, -26)},
		{S(165, 33), S(129, 39), S(56, 77), S(25, 136), S(27, 138), S(60, 79), S(126, 38), S(167, 31)},
		{S(107, 32), S(154, 30), S(113, 63), S(80, 136), S(80, 136), S(112, 63), S(154, 30), S(108, 33)},
		{S(87, 24), S(237, 25), S(156, 70), S(96, 145), S(96, 145), S(156, 70), S(237, 25), S(87, 24)},
		{S(89, 44), S(270, 37), S(181, 82), S(108, 144), S(108, 144), S(181, 82), S(270, 37), S(89, 44)},
		{S(117, 43), S(298, 58), S(248, 80), S(184, 121), S(184, 121), S(248, 80), S(298, 58), S(117, 43)},
		{S(115, 0), S(211, 66), S(189, 73), S(149, 113), S(149, 113), S(189, 73), S(211, 66), S(115, 0)},
		{S(157, -98), S(290, -17), S(173, 18), S(161, 75), S(161, 75), S(173, 18), S(290, -17), S(157, -98)},
	},
}

var PawnsConnected = [7][4]Score{
	{S(0, 0), S(0, 0), S(0, 0), S(0, 0)},
	{S(-2, -13), S(8, 4), S(3, 0), S(17, 18)},
	{S(11, 6), S(22, 4), S(24, 11), S(27, 16)},
	{S(10, 5), S(21, 7), S(11, 8), S(17, 20)},
	{S(5, 17), S(15, 26), S(31, 27), S(28, 20)},
	{S(38, 27), S(30, 60), S(76, 61), S(95, 73)},
	{S(176, 30), S(300, 24), S(294, 40), S(346, 44)},
}

var MobilityBonus = [...][32]Score{
	{S(-60, -134), S(-43, -82), S(-28, -36), S(-19, -12), S(-12, 2), S(-8, 17), // Knights
		S(1, 20), S(9, 15), S(19, 1)},
	{S(-1, -138), S(6, -64), S(14, -25), S(22, -1), S(29, 14), S(33, 31), // Bishops
		S(36, 39), S(36, 43), S(35, 46), S(39, 46), S(42, 42), S(53, 34),
		S(82, 35), S(91, 13)},
	{S(-127, -146), S(-15, -37), S(-3, 14), S(-4, 42), S(0, 54), S(2, 66), // Rooks
		S(3, 75), S(8, 78), S(13, 83), S(17, 87), S(20, 91), S(22, 94),
		S(28, 93), S(42, 83), S(98, 47)},
	{S(-413, -159), S(-122, -138), S(-31, -174), S(-13, -119), S(-2, -84), S(-1, -14), // Queens
		S(3, 11), S(7, 29), S(10, 45), S(13, 57), S(16, 63), S(21, 63),
		S(22, 65), S(23, 69), S(26, 62), S(23, 64), S(21, 59), S(19, 57),
		S(24, 44), S(31, 29), S(41, 10), S(38, -2), S(39, -21), S(49, -50),
		S(14, -38), S(-88, -10), S(127, -127), S(49, -81)},
}

var PassedFriendlyDistance = [8]Score{
	S(0, 0), S(-11, 41), S(-14, 24), S(-13, 9),
	S(-6, -6), S(-2, -15), S(16, -27), S(0, -38),
}

var PassedEnemyDistance = [8]Score{
	S(0, 0), S(-129, -57), S(-15, -8), S(2, 12),
	S(13, 22), S(14, 29), S(11, 36), S(18, 41),
}
var PawnPsqt [16][2][64]Score    // BishopFlag, colour, Square
var Psqt [2][Queen + 1][64]Score // One row for every colour purposefelly left empty
var KingPsqt [16][2][64]Score    // BishopFlag, colour, Square

var PawnsConnectedSquare [2][64]Score
var pawnsConnectedMask [2][64]uint64

// PassedRank[Rank] contains a bonus according to the rank of a passed pawn
var PassedRank = [7]Score{S(0, 0), S(-13, -31), S(-20, -13), S(-9, 25), S(29, 67), S(52, 145), S(178, 237)}

// PassedFile[File] contains a bonus according to the file of a passed pawn
var PassedFile = [8]Score{S(-11, 34), S(-13, 34), S(-7, 12), S(-4, -4),
	S(-10, 1), S(-17, 10), S(-16, 23), S(12, 10),
}

var PassedStacked = [8]Score{S(0, 0), S(-14, -53), S(-20, -34), S(-31, -56), S(-3, -85), S(28, -207), S(0, 0), S(0, 0)}

var Isolated = S(-8, -16)
var Doubled = S(-9, -33)
var Backward = S(7, -1)
var BackwardOpen = S(-5, -15)

var BishopPair = S(20, 75)
var BishopRammedPawns = S(-8, -23)

var BishopOutpostUndefendedBonus = S(20, -2)
var BishopOutpostDefendedBonus = S(50, 13)

var LongDiagonalBishop = S(23, 23)

var KnightOutpostUndefendedBonus = S(19, -20)
var KnightOutpostDefendedBonus = S(44, 13)

var DistantKnight = [4]Score{S(-12, -3), S(-18, -11), S(-27, -16), S(-52, -20)}

var MinorBehindPawn = S(8, 32)

var Tempo int16 = 15

// Rook on semiopen, open file
var RookOnFile = [2]Score{S(9, 18), S(35, 11)}
var RookOnQueenFile = S(6, 30)

var KingDefenders = [12]Score{
	S(-4, -9), S(3, -8), S(5, -8), S(6, -6),
	S(8, -3), S(9, 5), S(11, 11), S(13, 10),
	S(11, 8), S(15, -70), S(-3, -10), S(11, 0),
}

var KingShelter = [2][8][8]Score{
	{{S(-27, 4), S(8, -17), S(16, -5), S(19, -1),
		S(23, -15), S(18, -10), S(13, -40), S(-95, 34)},
		{S(9, -2), S(24, -16), S(-2, -2), S(-10, -4),
			S(-2, -10), S(24, -10), S(41, -25), S(-51, 17)},
		{S(12, 11), S(0, 8), S(-32, 11), S(-28, 9),
			S(-8, -7), S(-6, 12), S(15, -1), S(-37, 9)},
		{S(-28, 29), S(10, 7), S(-5, -2), S(-1, -2),
			S(7, -19), S(4, -5), S(26, 8), S(-15, 4)},
		{S(-28, 19), S(-15, 12), S(-33, 11), S(-22, 9),
			S(-1, -5), S(-28, 5), S(3, -7), S(-28, 6)},
		{S(38, -14), S(21, -9), S(3, -7), S(3, -17),
			S(12, -25), S(24, -16), S(31, -17), S(-13, 1)},
		{S(16, -11), S(-9, -15), S(-28, -6), S(-18, -6),
			S(-8, -23), S(8, -13), S(13, -24), S(-24, 12)},
		{S(-39, 6), S(-25, -4), S(-17, 21), S(-20, 19),
			S(-22, 11), S(-10, 8), S(-11, -18), S(-75, 46)}},
	{{S(48, 42), S(-45, -10), S(-28, 2), S(-39, -1),
		S(-51, -12), S(-22, 17), S(-67, 5), S(-92, 38)},
		{S(149, 9), S(0, -9), S(-5, 0), S(-28, 13),
			S(-15, -14), S(15, -7), S(8, -7), S(-87, 29)},
		{S(5, 22), S(42, -13), S(14, -14), S(23, -20),
			S(30, -19), S(16, -12), S(46, -12), S(-27, 1)},
		{S(6, 49), S(-11, 34), S(-17, 28), S(-24, 32),
			S(-2, 18), S(-17, 17), S(-12, 21), S(-43, 7)},
		{S(-10, 44), S(-3, 25), S(-7, 19), S(-13, 16),
			S(-5, 15), S(3, -7), S(14, -5), S(-31, 9)},
		{S(34, 3), S(-3, -5), S(-8, -2), S(-3, -3),
			S(2, -11), S(-44, -2), S(-8, -13), S(-40, 2)},
		{S(43, -1), S(-1, -18), S(-4, -15), S(-22, -10),
			S(-6, -24), S(4, -21), S(4, -23), S(-71, 18)},
		{S(-5, -3), S(-16, -31), S(-8, -18), S(-4, -24),
			S(-8, -34), S(-1, -25), S(-6, -49), S(-70, 28)}},
}

var KingStorm = [2][4][8]Score{
	{{S(11, 5), S(4, 11), S(9, 4), S(5, 5),
		S(2, 2), S(10, -1), S(8, 6), S(-6, -11)},
		{S(21, -7), S(15, -3), S(17, -4), S(3, 4),
			S(17, -6), S(25, -17), S(17, -15), S(-10, -13)},
		{S(23, 10), S(4, 7), S(6, 12), S(1, 13),
			S(5, 8), S(12, 1), S(10, -12), S(1, 2)},
		{S(8, 18), S(6, 14), S(9, 9), S(-1, 4),
			S(-3, 7), S(15, 4), S(8, -2), S(-10, 10)}},
	{{S(0, 0), S(19, -5), S(-18, 7), S(10, -10),
		S(-2, 15), S(-14, 31), S(4, 51), S(-2, -8)},
		{S(0, 0), S(-41, -20), S(2, -10), S(44, -12),
			S(1, -8), S(-11, -10), S(13, 56), S(-10, -8)},
		{S(0, 0), S(-17, -9), S(-7, -3), S(17, -1),
			S(8, -7), S(-7, -18), S(47, -47), S(-6, 7)},
		{S(0, 0), S(-13, -13), S(-17, -10), S(-8, 3),
			S(2, -14), S(0, -49), S(6, 3), S(-12, 12)}},
}
var KingOnPawnlessFlank = S(15, -87)

var passedMask [2][64]uint64

var outpustMask [2][64]uint64

var distanceBetween [64][64]int16

var adjacentFilesMask [8]uint64

var whiteKingAreaMask [64]uint64
var blackKingAreaMask [64]uint64

var forwardRanksMask [2][8]uint64

var forwardFileMask [2][64]uint64

// Outpost bitboards
const whiteOutpustRanks = RANK_4_BB | RANK_5_BB | RANK_6_BB
const blackOutpustRanks = RANK_5_BB | RANK_4_BB | RANK_3_BB

var KingSafetyAttacksWeights = [King + 1]int16{0, -3, -7, -4, 5, 0}
var KingSafetyAttackValue int16 = 122
var KingSafetyWeakSquares int16 = 44
var KingSafetyFriendlyPawns int16 = -34
var KingSafetyNoEnemyQueens int16 = -176
var KingSafetySafeQueenCheck int16 = 89
var KingSafetySafeRookCheck int16 = 74
var KingSafetySafeBishopCheck int16 = 53
var KingSafetySafeKnightCheck int16 = 111
var KingSafetyAdjustment int16 = -10

var Hanging = S(34, 13)
var ThreatByKing = S(-9, 33)
var ThreatByMinor = [King + 1]Score{S(0, 0), S(20, 39), S(17, 35), S(32, 26), S(29, 32), S(10, 23)}
var ThreatByRook = [King + 1]Score{S(0, 0), S(-3, 11), S(-1, 16), S(-5, -9), S(34, 7), S(27, -1)}

func LoadScoresToPieceSquares() {
	for x := 0; x < 4; x++ {
		for y := 0; y < 8; y++ {
			Psqt[White][Knight][y*8+x] = PieceScores[Knight][y][x] + KnightValue
			Psqt[White][Knight][y*8+(7-x)] = PieceScores[Knight][y][x] + KnightValue
			Psqt[Black][Knight][(7-y)*8+x] = PieceScores[Knight][y][x] + KnightValue
			Psqt[Black][Knight][(7-y)*8+(7-x)] = PieceScores[Knight][y][x] + KnightValue

			Psqt[White][Bishop][y*8+x] = PieceScores[Bishop][y][x] + BishopValue
			Psqt[White][Bishop][y*8+(7-x)] = PieceScores[Bishop][y][x] + BishopValue
			Psqt[Black][Bishop][(7-y)*8+x] = PieceScores[Bishop][y][x] + BishopValue
			Psqt[Black][Bishop][(7-y)*8+(7-x)] = PieceScores[Bishop][y][x] + BishopValue

			Psqt[White][Rook][y*8+x] = PieceScores[Rook][y][x] + RookValue
			Psqt[White][Rook][y*8+(7-x)] = PieceScores[Rook][y][x] + RookValue
			Psqt[Black][Rook][(7-y)*8+x] = PieceScores[Rook][y][x] + RookValue
			Psqt[Black][Rook][(7-y)*8+(7-x)] = PieceScores[Rook][y][x] + RookValue

			Psqt[White][Queen][y*8+x] = PieceScores[Queen][y][x] + QueenValue
			Psqt[White][Queen][y*8+(7-x)] = PieceScores[Queen][y][x] + QueenValue
			Psqt[Black][Queen][(7-y)*8+x] = PieceScores[Queen][y][x] + QueenValue
			Psqt[Black][Queen][(7-y)*8+(7-x)] = PieceScores[Queen][y][x] + QueenValue

			if y != 7 {
				PawnsConnectedSquare[White][y*8+x] = PawnsConnected[y][x]
				PawnsConnectedSquare[White][y*8+(7-x)] = PawnsConnected[y][x]
				PawnsConnectedSquare[Black][(7-y)*8+x] = PawnsConnected[y][x]
				PawnsConnectedSquare[Black][(7-y)*8+(7-x)] = PawnsConnected[y][x]
			}
		}
	}

	for bishopFlag := 0; bishopFlag <= 15; bishopFlag++ {
		for y := 1; y < 7; y++ {
			for x := 0; x < 8; x++ {
				PawnPsqt[bishopFlag][White][y*8+x] = PawnScores[bishopFlag][y][x] + PawnValue
				PawnPsqt[bishopFlag][Black][(7-y)*8+x] = PawnScores[BishopFlag(bishopFlag).BlackPerspective()][y][x] + PawnValue
			}
		}
	}

	for bishopFlag := 0; bishopFlag <= 15; bishopFlag++ {
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				KingPsqt[bishopFlag][White][y*8+x] = KingScores[bishopFlag][y][x]
				KingPsqt[bishopFlag][Black][(7-y)*8+x] = KingScores[BishopFlag(bishopFlag).BlackPerspective()][y][x]
			}
		}
	}
}

func init() {
	LoadScoresToPieceSquares()

	// Pawn is passed if no pawn of opposite color can stop it from promoting
	for i := 8; i <= 55; i++ {
		passedMask[White][i] = 0
		for file := File(i) - 1; file <= File(i)+1; file++ {
			if file < FILE_A || file > FILE_H {
				continue
			}
			for rank := Rank(i) + 1; rank < RANK_8; rank++ {
				passedMask[White][i] |= 1 << uint(rank*8+file)
			}
		}
	}
	// Outpust is similar to passed bitboard bot we do not care about pawns in same file
	for i := 8; i <= 55; i++ {
		outpustMask[White][i] = passedMask[White][i] & ^FILES[File(i)]
	}

	for i := 55; i >= 8; i-- {
		passedMask[Black][i] = 0
		for file := File(i) - 1; file <= File(i)+1; file++ {
			if file < FILE_A || file > FILE_H {
				continue
			}
			for rank := Rank(i) - 1; rank > RANK_1; rank-- {
				passedMask[Black][i] |= 1 << uint(rank*8+file)
			}
		}
	}
	for i := 55; i >= 8; i-- {
		outpustMask[Black][i] = passedMask[Black][i] & ^FILES[File(i)]
	}

	for i := 8; i <= 55; i++ {
		pawnsConnectedMask[White][i] = PawnAttacks[Black][i] | PawnAttacks[Black][i+8]
		pawnsConnectedMask[Black][i] = PawnAttacks[White][i] | PawnAttacks[White][i-8]
	}

	for i := range FILES {
		adjacentFilesMask[i] = 0
		if i != 0 {
			adjacentFilesMask[i] |= FILES[i-1]
		}
		if i != 7 {
			adjacentFilesMask[i] |= FILES[i+1]
		}
	}

	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			distanceBetween[y][x] = int16(Max(Abs(Rank(y)-Rank(x)), Abs(File(y)-File(x))))
		}
	}

	for y := 0; y < 64; y++ {
		whiteKingAreaMask[y] = KingAttacks[y] | SquareMask[y] | North(KingAttacks[y])
		blackKingAreaMask[y] = KingAttacks[y] | SquareMask[y] | South(KingAttacks[y])
		if File(y) > FILE_A {
			whiteKingAreaMask[y] |= West(whiteKingAreaMask[y])
			blackKingAreaMask[y] |= West(blackKingAreaMask[y])
		}
		if File(y) < FILE_H {
			whiteKingAreaMask[y] |= East(whiteKingAreaMask[y])
			blackKingAreaMask[y] |= East(blackKingAreaMask[y])
		}
	}

	for rank := RANK_1; rank <= RANK_8; rank++ {
		for y := rank; y <= RANK_8; y++ {
			forwardRanksMask[White][rank] |= RANKS[y]
		}
		forwardRanksMask[Black][rank] = (^forwardRanksMask[White][rank]) | RANKS[rank]
	}

	for y := 0; y < 64; y++ {
		forwardFileMask[White][y] = forwardRanksMask[White][Rank(y)] & FILES[File(y)] & ^SquareMask[y]
		forwardFileMask[Black][y] = forwardRanksMask[Black][Rank(y)] & FILES[File(y)] & ^SquareMask[y]
	}
}

func IsLateEndGame(pos *Position) bool {
	return ((pos.Pieces[Rook] | pos.Pieces[Queen] | pos.Pieces[Bishop] | pos.Pieces[Knight]) & pos.Colours[pos.SideToMove]) == 0
}

func evaluateKingPawns(pos *Position) Score {
	if !tuning {
		if ok, score := GlobalPawnKingTable.Get(pos.PawnKey); ok {
			return score
		}
	}
	var fromBB uint64
	var fromId int
	whiteKingLocation := BitScan(pos.Pieces[King] & pos.Colours[White])
	blackKingLocation := BitScan(pos.Pieces[King] & pos.Colours[Black])
	score := SCORE_ZERO

	// white pawns
	for fromBB = pos.Pieces[Pawn] & pos.Colours[White]; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)

		// Passed bonus
		if passedMask[White][fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 {
			// Bonus is calculated based on rank, file, distance from friendly and enemy king
			score +=
				PassedRank[Rank(fromId)] +
					PassedFile[File(fromId)] +
					PassedFriendlyDistance[distanceBetween[whiteKingLocation][fromId]] +
					PassedEnemyDistance[distanceBetween[blackKingLocation][fromId]]

			if tuning {
				T.PassedRank[Rank(fromId)]++
				T.PassedFile[File(fromId)]++
				T.PassedFriendlyDistance[distanceBetween[whiteKingLocation][fromId]]++
				T.PassedEnemyDistance[distanceBetween[blackKingLocation][fromId]]++
			}

			if pos.Pieces[Pawn]&pos.Colours[White]&forwardFileMask[White][fromId] != 0 {
				score += PassedStacked[Rank(fromId)]
				if tuning {
					T.PassedStacked[Rank(fromId)]++
				}
			}
		}

		// Isolated pawn penalty
		if adjacentFilesMask[File(fromId)]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 {
			score += Isolated
			if tuning {
				T.Isolated++
			}
		}

		// Pawn is backward if there are no pawns behind it and cannot increase rank without being attacked by enemy pawn
		if passedMask[Black][fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 &&
			PawnAttacks[White][fromId+8]&(pos.Pieces[Pawn]&pos.Colours[Black]) != 0 {
			if FILES[File(fromId)]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 {
				score += BackwardOpen
				if tuning {
					T.BackwardOpen++
				}
			} else {
				score += Backward
				if tuning {
					T.Backward++
				}
			}
		} else if pawnsConnectedMask[White][fromId]&(pos.Colours[White]&pos.Pieces[Pawn]) != 0 {
			score += PawnsConnectedSquare[White][fromId]
			if tuning {
				T.PawnsConnected[Rank(fromId)][FileMirror[File(fromId)]]++
			}
		}
	}

	// white doubled pawns
	score += Score(PopCount(pos.Pieces[Pawn]&pos.Colours[White]&South(pos.Pieces[Pawn]&pos.Colours[White]))) * Doubled
	if tuning {
		T.Doubled += PopCount(pos.Pieces[Pawn] & pos.Colours[White] & South(pos.Pieces[Pawn]&pos.Colours[White]))
	}

	// black pawns
	for fromBB = pos.Pieces[Pawn] & pos.Colours[Black]; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)

		if passedMask[Black][fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 {
			score -=
				PassedRank[7-Rank(fromId)] +
					PassedFile[File(fromId)] +
					PassedFriendlyDistance[distanceBetween[blackKingLocation][fromId]] +
					PassedEnemyDistance[distanceBetween[whiteKingLocation][fromId]]
			if tuning {
				T.PassedRank[7-Rank(fromId)]--
				T.PassedFile[File(fromId)]--
				T.PassedFriendlyDistance[distanceBetween[blackKingLocation][fromId]]--
				T.PassedEnemyDistance[distanceBetween[whiteKingLocation][fromId]]--
			}

			if pos.Pieces[Pawn]&pos.Colours[Black]&forwardFileMask[Black][fromId] != 0 {
				score -= PassedStacked[7-Rank(fromId)]
				if tuning {
					T.PassedStacked[7-Rank(fromId)]--
				}
			}
		}
		if adjacentFilesMask[File(fromId)]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 {
			score -= Isolated
			if tuning {
				T.Isolated--
			}
		}
		if passedMask[White][fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 &&
			PawnAttacks[Black][fromId-8]&(pos.Pieces[Pawn]&pos.Colours[White]) != 0 {
			if FILES[File(fromId)]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 {
				score -= BackwardOpen
				if tuning {
					T.BackwardOpen--
				}
			} else {
				score -= Backward
				if tuning {
					T.Backward--
				}
			}
		} else if pawnsConnectedMask[Black][fromId]&(pos.Colours[Black]&pos.Pieces[Pawn]) != 0 {
			score -= PawnsConnectedSquare[Black][fromId]
			if tuning {
				T.PawnsConnected[7-Rank(fromId)][FileMirror[File(fromId)]]--
			}
		}
	}

	// black doubled pawns
	score -= Score(PopCount(pos.Pieces[Pawn]&pos.Colours[Black]&North(pos.Pieces[Pawn]&pos.Colours[Black]))) * Doubled
	if tuning {
		T.Doubled -= PopCount(pos.Pieces[Pawn] & pos.Colours[Black] & North(pos.Pieces[Pawn]&pos.Colours[Black]))
	}

	// White king storm shelter
	for file := Max(File(whiteKingLocation)-1, FILE_A); file <= Min(File(whiteKingLocation)+1, FILE_H); file++ {
		ours := pos.Pieces[Pawn] & FILES[file] & pos.Colours[White] & forwardRanksMask[White][Rank(whiteKingLocation)]
		var ourDist int
		if ours == 0 {
			ourDist = 7
		} else {
			ourDist = Abs(Rank(whiteKingLocation) - Rank(BitScan(ours)))
		}
		theirs := pos.Pieces[Pawn] & FILES[file] & pos.Colours[Black] & forwardRanksMask[White][Rank(whiteKingLocation)]
		var theirDist int
		if theirs == 0 {
			theirDist = 7
		} else {
			theirDist = Abs(Rank(whiteKingLocation) - Rank(BitScan(theirs)))
		}
		sameFile := BoolToInt(file == File(whiteKingLocation))
		score += KingShelter[sameFile][file][ourDist]
		if tuning {
			T.KingShelter[sameFile][file][ourDist]++
		}

		blocked := BoolToInt(ourDist != 7 && ourDist == theirDist-1)
		score += KingStorm[blocked][FileMirror[file]][theirDist]

		if tuning {
			T.KingStorm[blocked][FileMirror[file]][theirDist]++
		}
	}
	if KING_FLANK_BB[File(whiteKingLocation)]&pos.Pieces[Pawn] == 0 {
		score += KingOnPawnlessFlank
		if tuning {
			T.KingOnPawnlessFlank++
		}
	}

	// Black king storm / shelter
	for file := Max(File(blackKingLocation)-1, FILE_A); file <= Min(File(blackKingLocation)+1, FILE_H); file++ {
		ours := pos.Pieces[Pawn] & FILES[file] & pos.Colours[Black] & forwardRanksMask[Black][Rank(blackKingLocation)]
		var ourDist int
		if ours == 0 {
			ourDist = 7
		} else {
			ourDist = Abs(Rank(blackKingLocation) - Rank(MostSignificantBit(ours)))
		}
		theirs := pos.Pieces[Pawn] & FILES[file] & pos.Colours[White] & forwardRanksMask[Black][Rank(blackKingLocation)]
		var theirDist int
		if theirs == 0 {
			theirDist = 7
		} else {
			theirDist = Abs(Rank(blackKingLocation) - Rank(MostSignificantBit(theirs)))
		}
		sameFile := BoolToInt(file == File(blackKingLocation))
		score -= KingShelter[sameFile][file][ourDist]
		if tuning {
			T.KingShelter[sameFile][file][ourDist]--
		}

		blocked := BoolToInt(ourDist != 7 && ourDist == theirDist-1)
		score -= KingStorm[blocked][FileMirror[file]][theirDist]
		if tuning {
			T.KingStorm[blocked][FileMirror[file]][theirDist]--
		}
	}
	if KING_FLANK_BB[File(blackKingLocation)]&pos.Pieces[Pawn] == 0 {
		score -= KingOnPawnlessFlank
		if tuning {
			T.KingOnPawnlessFlank--
		}
	}
	if !tuning {
		GlobalPawnKingTable.Set(pos.PawnKey, score)
	}
	return score
}

func Evaluate(pos *Position) int {
	var fromId int
	var fromBB uint64
	var attacks uint64

	var whiteAttacked uint64
	var whiteAttackedBy [King + 1]uint64
	var whiteAttackedByTwo uint64
	var blackAttacked uint64
	var whiteKingAttacksCount int16
	var whiteKingAttackersCount int16
	var whiteKingAttackersWeight int16
	var blackAttackedBy [King + 1]uint64
	var blackAttackedByTwo uint64
	var blackKingAttacksCount int16
	var blackKingAttackersCount int16
	var blackKingAttackersWeight int16

	phase := TotalPhase
	whiteMobilityArea := ^((pos.Pieces[Pawn] & pos.Colours[White]) | (BlackPawnsAttacks(pos.Pieces[Pawn] & pos.Colours[Black])))
	blackMobilityArea := ^((pos.Pieces[Pawn] & pos.Colours[Black]) | (WhitePawnsAttacks(pos.Pieces[Pawn] & pos.Colours[White])))
	allOccupation := pos.Colours[White] | pos.Colours[Black]

	whiteKingLocation := BitScan(pos.Pieces[King] & pos.Colours[White])
	attacks = KingAttacks[whiteKingLocation]
	whiteAttacked |= attacks
	whiteAttackedBy[King] |= attacks
	whiteKingArea := whiteKingAreaMask[whiteKingLocation]

	blackKingLocation := BitScan(pos.Pieces[King] & pos.Colours[Black])
	attacks = KingAttacks[blackKingLocation]
	blackAttacked |= attacks
	blackAttackedBy[King] |= attacks
	blackKingArea := blackKingAreaMask[blackKingLocation]

	// white pawns
	attacks = WhitePawnsAttacks(pos.Pieces[Pawn] & pos.Colours[White])
	whiteAttackedByTwo |= whiteAttacked & attacks
	whiteAttacked |= attacks
	whiteAttackedBy[Pawn] |= attacks
	whiteKingAttacksCount += int16(PopCount(attacks & blackKingArea))

	// black pawns
	attacks = BlackPawnsAttacks(pos.Pieces[Pawn] & pos.Colours[Black])
	blackAttackedByTwo |= blackAttacked & attacks
	blackAttacked |= attacks
	blackAttackedBy[Pawn] |= attacks
	blackKingAttacksCount += int16(PopCount(attacks & whiteKingArea))

	score := evaluateKingPawns(pos)

	// white pawns
	for fromBB = pos.Pieces[Pawn] & pos.Colours[White]; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)

		score += PawnPsqt[pos.BishopFlag][White][fromId]
		if tuning {
			T.PawnValue++
			T.PawnScores[pos.BishopFlag][Rank(fromId)][File(fromId)]++
		}
	}

	// black pawns
	for fromBB = pos.Pieces[Pawn] & pos.Colours[Black]; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)

		score -= PawnPsqt[pos.BishopFlag][Black][fromId]

		if tuning {
			T.PawnValue--
			T.PawnScores[pos.BishopFlag.BlackPerspective()][7-Rank(fromId)][File(fromId)]--
		}
	}

	// white knights
	for fromBB = pos.Pieces[Knight] & pos.Colours[White]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= KnightPhase
		fromId = BitScan(fromBB)

		attacks = KnightAttacks[fromId]
		mobility := PopCount(whiteMobilityArea & attacks)
		score += Psqt[White][Knight][fromId]
		score += MobilityBonus[0][mobility]
		if tuning {
			T.KnightValue++
			T.PieceScores[Knight][Rank(fromId)][FileMirror[File(fromId)]]++
			T.MobilityBonus[0][mobility]++
		}

		whiteAttackedByTwo |= whiteAttacked & attacks
		whiteAttacked |= attacks
		whiteAttackedBy[Knight] |= attacks

		if (pos.Pieces[Pawn]>>8)&SquareMask[fromId] != 0 {
			score += MinorBehindPawn
			if tuning {
				T.MinorBehindPawn++
			}
		}
		if SquareMask[fromId]&whiteOutpustRanks != 0 && outpustMask[White][fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 {
			if PawnAttacks[Black][fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) != 0 {
				score += KnightOutpostDefendedBonus
				if tuning {
					T.KnightOutpostDefendedBonus++
				}
			} else {
				score += KnightOutpostUndefendedBonus
				if tuning {
					T.KnightOutpostUndefendedBonus++
				}
			}
		}

		kingDistance := Min(int(distanceBetween[fromId][whiteKingLocation]), int(distanceBetween[fromId][blackKingLocation]))
		if kingDistance >= 4 {
			score += DistantKnight[kingDistance-4]
			if tuning {
				T.DistantKnight[kingDistance-4]++
			}
		}
		if attacks&blackKingArea != 0 {
			whiteKingAttacksCount += int16(PopCount(attacks & blackKingArea))
			whiteKingAttackersCount++
			whiteKingAttackersWeight += KingSafetyAttacksWeights[Knight]
		}
	}

	// black knights
	for fromBB = pos.Pieces[Knight] & pos.Colours[Black]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= KnightPhase
		fromId = BitScan(fromBB)

		attacks = KnightAttacks[fromId]
		mobility := PopCount(blackMobilityArea & attacks)
		score -= Psqt[Black][Knight][fromId]
		score -= MobilityBonus[0][mobility]
		if tuning {
			T.KnightValue--
			T.PieceScores[Knight][7-Rank(fromId)][FileMirror[File(fromId)]]--
			T.MobilityBonus[0][mobility]--
		}

		blackAttackedByTwo |= blackAttacked & attacks
		blackAttacked |= attacks
		blackAttackedBy[Knight] |= attacks

		if (pos.Pieces[Pawn]<<8)&SquareMask[fromId] != 0 {
			score -= MinorBehindPawn
			if tuning {
				T.MinorBehindPawn--
			}
		}
		if SquareMask[fromId]&blackOutpustRanks != 0 && outpustMask[Black][fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 {
			if PawnAttacks[White][fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) != 0 {
				score -= KnightOutpostDefendedBonus
				if tuning {
					T.KnightOutpostDefendedBonus--
				}
			} else {
				score -= KnightOutpostUndefendedBonus
				if tuning {
					T.KnightOutpostUndefendedBonus--
				}
			}
		}
		kingDistance := Min(int(distanceBetween[fromId][whiteKingLocation]), int(distanceBetween[fromId][blackKingLocation]))
		if kingDistance >= 4 {
			score -= DistantKnight[kingDistance-4]
			if tuning {
				T.DistantKnight[kingDistance-4]--
			}
		}
		if attacks&whiteKingArea != 0 {
			blackKingAttacksCount += int16(PopCount(attacks & whiteKingArea))
			blackKingAttackersCount++
			blackKingAttackersWeight += KingSafetyAttacksWeights[Knight]
		}
	}

	// white bishops
	whiteRammedPawns := South(pos.Pieces[Pawn]&pos.Colours[Black]) & (pos.Pieces[Pawn] & pos.Colours[White])
	for fromBB = pos.Pieces[Bishop] & pos.Colours[White]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= BishopPhase
		fromId = BitScan(fromBB)

		attacks = BishopAttacks(fromId, allOccupation)
		mobility := PopCount(whiteMobilityArea & attacks)
		score += MobilityBonus[1][mobility]
		score += Psqt[White][Bishop][fromId]
		if tuning {
			T.BishopValue++
			T.PieceScores[Bishop][Rank(fromId)][FileMirror[File(fromId)]]++
			T.MobilityBonus[1][mobility]++
		}

		whiteAttackedByTwo |= whiteAttacked & attacks
		whiteAttacked |= attacks
		whiteAttackedBy[Bishop] |= attacks

		if (pos.Pieces[Pawn]>>8)&SquareMask[fromId] != 0 {
			score += MinorBehindPawn
			if tuning {
				T.MinorBehindPawn++
			}
		}
		if (LONG_DIAGONALS&SquareMask[fromId]) != 0 && (MoreThanOne(BishopAttacks(fromId, pos.Pieces[Pawn]) & CENTER)) {
			score += LongDiagonalBishop
			if tuning {
				T.LongDiagonalBishop++
			}
		}
		if SquareMask[fromId]&whiteOutpustRanks != 0 && outpustMask[White][fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 {
			if PawnAttacks[Black][fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) != 0 {
				score += BishopOutpostDefendedBonus
				if tuning {
					T.BishopOutpostDefendedBonus++
				}
			} else {
				score += BishopOutpostUndefendedBonus
				if tuning {
					T.BishopOutpostUndefendedBonus++
				}
			}
		}

		// Bishop is worth less if there are friendly rammed pawns of its color
		var rammedCount Score
		if SquareMask[fromId]&WHITE_SQUARES != 0 {
			rammedCount = Score(PopCount(whiteRammedPawns & WHITE_SQUARES))
		} else {
			rammedCount = Score(PopCount(whiteRammedPawns & BLACK_SQUARES))
		}
		score += BishopRammedPawns * rammedCount
		if tuning {
			T.BishopRammedPawns += int(rammedCount)
		}
		if attacks&blackKingArea != 0 {
			whiteKingAttacksCount += int16(PopCount(attacks & blackKingArea))
			whiteKingAttackersCount++
			whiteKingAttackersWeight += KingSafetyAttacksWeights[Bishop]
		}
	}

	// Bishop pair bonus
	// It is not checked if bishops have opposite colors, but that is almost always the case
	if MoreThanOne(pos.Pieces[Bishop] & pos.Colours[White]) {
		score += BishopPair
		if tuning {
			T.BishopPair++
		}
	}

	// black bishops
	blackRammedPawns := North(pos.Pieces[Pawn]&pos.Colours[White]) & (pos.Pieces[Pawn] & pos.Colours[Black])
	for fromBB = pos.Pieces[Bishop] & pos.Colours[Black]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= BishopPhase
		fromId = BitScan(fromBB)

		attacks = BishopAttacks(fromId, allOccupation)
		mobility := PopCount(blackMobilityArea & attacks)
		score -= MobilityBonus[1][mobility]
		score -= Psqt[Black][Bishop][fromId]
		if tuning {
			T.BishopValue--
			T.PieceScores[Bishop][7-Rank(fromId)][FileMirror[File(fromId)]]--
			T.MobilityBonus[1][mobility]--
		}

		blackAttackedByTwo |= blackAttacked & attacks
		blackAttacked |= attacks
		blackAttackedBy[Bishop] |= attacks

		if (pos.Pieces[Pawn]<<8)&SquareMask[fromId] != 0 {
			score -= MinorBehindPawn
			if tuning {
				T.MinorBehindPawn--
			}
		}
		if (LONG_DIAGONALS&SquareMask[fromId]) != 0 && (MoreThanOne(BishopAttacks(fromId, pos.Pieces[Pawn]) & CENTER)) {
			score -= LongDiagonalBishop
			if tuning {
				T.LongDiagonalBishop--
			}
		}
		if SquareMask[fromId]&blackOutpustRanks != 0 && outpustMask[Black][fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 {
			if PawnAttacks[White][fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) != 0 {
				score -= BishopOutpostDefendedBonus
				if tuning {
					T.BishopOutpostDefendedBonus--
				}
			} else {
				score -= BishopOutpostUndefendedBonus
				if tuning {
					T.BishopOutpostUndefendedBonus--
				}
			}
		}
		var rammedCount Score
		if SquareMask[fromId]&WHITE_SQUARES != 0 {
			rammedCount = Score(PopCount(blackRammedPawns & WHITE_SQUARES))
		} else {
			rammedCount = Score(PopCount(blackRammedPawns & BLACK_SQUARES))
		}
		score -= BishopRammedPawns * rammedCount
		if tuning {
			T.BishopRammedPawns -= int(rammedCount)
		}
		if attacks&whiteKingArea != 0 {
			blackKingAttacksCount += int16(PopCount(attacks & whiteKingArea))
			blackKingAttackersCount++
			blackKingAttackersWeight += KingSafetyAttacksWeights[Bishop]
		}
	}

	if MoreThanOne(pos.Pieces[Bishop] & pos.Colours[Black]) {
		score -= BishopPair

		if tuning {
			T.BishopPair--
		}
	}

	// white rooks
	for fromBB = pos.Pieces[Rook] & pos.Colours[White]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= RookPhase
		fromId = BitScan(fromBB)

		attacks = RookAttacks(fromId, allOccupation)
		mobility := PopCount(whiteMobilityArea & attacks)
		score += MobilityBonus[2][mobility]
		score += Psqt[White][Rook][fromId]

		if tuning {
			T.RookValue++
			T.PieceScores[Rook][Rank(fromId)][FileMirror[File(fromId)]]++
			T.MobilityBonus[2][mobility]++
		}

		whiteAttackedByTwo |= whiteAttacked & attacks
		whiteAttacked |= attacks
		whiteAttackedBy[Rook] |= attacks

		if pos.Pieces[Pawn]&FILES[File(fromId)] == 0 {
			score += RookOnFile[1]
			if tuning {
				T.RookOnFile[1]++
			}
		} else if (pos.Pieces[Pawn]&pos.Colours[White])&FILES[File(fromId)] == 0 {
			score += RookOnFile[0]
			if tuning {
				T.RookOnFile[0]++
			}
		}

		if FileBB(fromId)&pos.Pieces[Queen] != 0 {
			score += RookOnQueenFile
			if tuning {
				T.RookOnQueenFile++
			}
		}

		if attacks&blackKingArea != 0 {
			whiteKingAttacksCount += int16(PopCount(attacks & blackKingArea))
			whiteKingAttackersCount++
			whiteKingAttackersWeight += KingSafetyAttacksWeights[Rook]
		}
	}

	// black rooks
	for fromBB = pos.Pieces[Rook] & pos.Colours[Black]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= RookPhase
		fromId = BitScan(fromBB)

		attacks = RookAttacks(fromId, allOccupation)
		mobility := PopCount(blackMobilityArea & attacks)
		score -= MobilityBonus[2][mobility]
		score -= Psqt[Black][Rook][fromId]

		if tuning {
			T.RookValue--
			T.PieceScores[Rook][7-Rank(fromId)][FileMirror[File(fromId)]]--
			T.MobilityBonus[2][mobility]--
		}

		blackAttackedByTwo |= blackAttacked & attacks
		blackAttacked |= attacks
		blackAttackedBy[Rook] |= attacks

		if pos.Pieces[Pawn]&FILES[File(fromId)] == 0 {
			score -= RookOnFile[1]
			if tuning {
				T.RookOnFile[1]--
			}
		} else if (pos.Pieces[Pawn]&pos.Colours[Black])&FILES[File(fromId)] == 0 {
			score -= RookOnFile[0]
			if tuning {
				T.RookOnFile[0]--
			}
		}

		if FileBB(fromId)&pos.Pieces[Queen] != 0 {
			score -= RookOnQueenFile
			if tuning {
				T.RookOnQueenFile--
			}
		}

		if attacks&whiteKingArea != 0 {
			blackKingAttacksCount += int16(PopCount(attacks & whiteKingArea))
			blackKingAttackersCount++
			blackKingAttackersWeight += KingSafetyAttacksWeights[Rook]
		}
	}

	//white queens
	for fromBB = pos.Pieces[Queen] & pos.Colours[White]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= QueenPhase
		fromId = BitScan(fromBB)

		attacks = QueenAttacks(fromId, allOccupation)
		mobility := PopCount(whiteMobilityArea & attacks)
		score += MobilityBonus[3][mobility]
		score += Psqt[White][Queen][fromId]

		if tuning {
			T.QueenValue++
			T.PieceScores[Queen][Rank(fromId)][FileMirror[File(fromId)]]++
			T.MobilityBonus[3][mobility]++
		}

		whiteAttackedByTwo |= whiteAttacked & attacks
		whiteAttacked |= attacks
		whiteAttackedBy[Queen] |= attacks

		if attacks&blackKingArea != 0 {
			whiteKingAttacksCount += int16(PopCount(attacks & blackKingArea))
			whiteKingAttackersCount++
			whiteKingAttackersWeight += KingSafetyAttacksWeights[Queen]
		}
	}

	// black queens
	for fromBB = pos.Pieces[Queen] & pos.Colours[Black]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= QueenPhase
		fromId = BitScan(fromBB)

		attacks = QueenAttacks(fromId, allOccupation)
		mobility := PopCount(blackMobilityArea & attacks)
		score -= MobilityBonus[3][mobility]
		score -= Psqt[Black][Queen][fromId]

		if tuning {
			T.QueenValue--
			T.PieceScores[Queen][7-Rank(fromId)][FileMirror[File(fromId)]]--
			T.MobilityBonus[3][mobility]--
		}

		blackAttackedByTwo |= blackAttacked & attacks
		blackAttacked |= attacks
		blackAttackedBy[Queen] |= attacks
		if attacks&whiteKingArea != 0 {
			blackKingAttacksCount += int16(PopCount(attacks & whiteKingArea))
			blackKingAttackersCount++
			blackKingAttackersWeight += KingSafetyAttacksWeights[Queen]
		}
	}

	if phase < 0 {
		phase = 0
	}

	// white king
	whiteKingDefenders := PopCount(
		(pos.Pieces[Pawn] | pos.Pieces[Bishop] | pos.Pieces[Knight]) & pos.Colours[White] & whiteKingAreaMask[whiteKingLocation],
	)
	score += KingPsqt[pos.BishopFlag][White][whiteKingLocation]
	score += KingDefenders[whiteKingDefenders]
	if tuning {
		T.KingScores[pos.BishopFlag][Rank(whiteKingLocation)][File(whiteKingLocation)]++
		T.KingDefenders[whiteKingDefenders]++
	}

	// Weak squares are attacked by the enemy, defended no more
	// than once and only defended by our Queens or our King
	weakForWhite := blackAttacked & ^whiteAttackedByTwo & (^whiteAttacked | whiteAttackedBy[Queen] | whiteAttackedBy[King])
	if int(blackKingAttackersCount) > 1-PopCount(pos.Colours[Black]&pos.Pieces[Queen]) {
		safe := ^pos.Colours[Black] & (^whiteAttacked | (weakForWhite & blackAttackedByTwo))

		knightThreats := KnightAttacks[whiteKingLocation]
		bishopThreats := BishopAttacks(whiteKingLocation, allOccupation)
		rookThreats := RookAttacks(whiteKingLocation, allOccupation)
		queenThreats := bishopThreats | rookThreats

		knightChecks := knightThreats & safe & blackAttackedBy[Knight]
		bishopChecks := bishopThreats & safe & blackAttackedBy[Bishop]
		rookChecks := rookThreats & safe & blackAttackedBy[Rook]
		queenChecks := queenThreats & safe & blackAttackedBy[Queen]

		count := int(blackKingAttackersCount) * int(blackKingAttackersWeight)
		count += int(KingSafetyAttackValue) * 9 * int(blackKingAttackersCount) / PopCount(whiteKingArea)
		count += int(KingSafetyWeakSquares) * PopCount(whiteKingArea&weakForWhite)
		count += int(KingSafetyFriendlyPawns) * PopCount(pos.Colours[White]&pos.Pieces[Pawn]&whiteKingArea & ^weakForWhite)
		count += int(KingSafetyNoEnemyQueens) * BoolToInt(pos.Colours[Black]&pos.Pieces[Queen] == 0)
		count += int(KingSafetySafeQueenCheck) * PopCount(queenChecks)
		count += int(KingSafetySafeRookCheck) * PopCount(rookChecks)
		count += int(KingSafetySafeBishopCheck) * PopCount(bishopChecks)
		count += int(KingSafetySafeKnightCheck) * PopCount(knightChecks)
		count += int(KingSafetyAdjustment)
		if count > 0 {
			score -= S(int16(count*count/720), int16(count/20))
		}
	}

	// black king
	blackKingDefenders := PopCount(
		(pos.Pieces[Pawn] | pos.Pieces[Bishop] | pos.Pieces[Knight]) & pos.Colours[Black] & blackKingAreaMask[blackKingLocation],
	)
	score -= KingPsqt[pos.BishopFlag][Black][blackKingLocation]
	score -= KingDefenders[blackKingDefenders]
	if tuning {
		T.KingScores[pos.BishopFlag.BlackPerspective()][7-Rank(blackKingLocation)][File(blackKingLocation)]--
		T.KingDefenders[blackKingDefenders]--
	}

	// Weak squares are attacked by the enemy, defended no more
	// than once and only defended by our Queens or our King
	weakForBlack := whiteAttacked & ^blackAttackedByTwo & (^blackAttacked | blackAttackedBy[Queen] | blackAttackedBy[King])
	if int(whiteKingAttackersCount) > 1-PopCount(pos.Colours[White]&pos.Pieces[Queen]) {
		safe := ^pos.Colours[White] & (^blackAttacked | (weakForBlack & whiteAttackedByTwo))

		knightThreats := KnightAttacks[blackKingLocation]
		bishopThreats := BishopAttacks(blackKingLocation, allOccupation)
		rookThreats := RookAttacks(blackKingLocation, allOccupation)
		queenThreats := bishopThreats | rookThreats

		knightChecks := knightThreats & safe & whiteAttackedBy[Knight]
		bishopChecks := bishopThreats & safe & whiteAttackedBy[Bishop]
		rookChecks := rookThreats & safe & whiteAttackedBy[Rook]
		queenChecks := queenThreats & safe & whiteAttackedBy[Queen]

		count := int(whiteKingAttackersCount) * int(whiteKingAttackersWeight)
		count += int(KingSafetyAttackValue) * int(whiteKingAttackersCount) * 9 / PopCount(blackKingArea) // Scale value to king area size
		count += int(KingSafetyWeakSquares) * PopCount(blackKingArea&weakForBlack)
		count += int(KingSafetyFriendlyPawns) * PopCount(pos.Colours[Black]&pos.Pieces[Pawn]&blackKingArea & ^weakForBlack)
		count += int(KingSafetyNoEnemyQueens) * BoolToInt(pos.Colours[White]&pos.Pieces[Queen] == 0)
		count += int(KingSafetySafeQueenCheck) * PopCount(queenChecks)
		count += int(KingSafetySafeRookCheck) * PopCount(rookChecks)
		count += int(KingSafetySafeBishopCheck) * PopCount(bishopChecks)
		count += int(KingSafetySafeKnightCheck) * PopCount(knightChecks)
		count += int(KingSafetyAdjustment)
		if count > 0 {
			score += S(int16(count*count/720), int16(count/20))
		}
	}

	// White threats
	blackStronglyProtected := blackAttackedBy[Pawn] | (blackAttackedByTwo & ^whiteAttackedByTwo)
	blackDefended := pos.Colours[Black] & ^pos.Pieces[Pawn] & blackStronglyProtected
	if ((pos.Colours[Black] & weakForBlack) | blackDefended) != 0 {
		for fromBB = pos.Colours[Black] & (blackDefended | weakForBlack) & (whiteAttackedBy[Knight] | whiteAttackedBy[Bishop]) & ^pos.Pieces[Pawn]; fromBB != 0; fromBB &= (fromBB - 1) {
			fromId = BitScan(fromBB)
			threatenedPiece := pos.TypeOnSquare(SquareMask[fromId])
			score += ThreatByMinor[threatenedPiece]
			if tuning {
				T.ThreatByMinor[threatenedPiece]++
			}
		}

		for fromBB = pos.Colours[Black] & (blackDefended | weakForBlack) & whiteAttackedBy[Rook] & ^pos.Pieces[Pawn]; fromBB != 0; fromBB &= (fromBB - 1) & ^pos.Pieces[Pawn] {
			fromId = BitScan(fromBB)
			threatenedPiece := pos.TypeOnSquare(SquareMask[fromId])
			score += ThreatByRook[threatenedPiece]
			if tuning {
				T.ThreatByRook[threatenedPiece]++
			}
		}

		if weakForBlack&pos.Colours[Black]&whiteAttackedBy[King] != 0 {
			score += ThreatByKing
			if tuning {
				T.ThreatByKing++
			}
		}

		// Bonus if enemy has a hanging piece
		score += Hanging *
			Score(PopCount((pos.Colours[Black] & ^pos.Pieces[Pawn] & whiteAttackedByTwo)&weakForBlack))

		if tuning {
			T.Hanging += PopCount((pos.Colours[Black] & ^pos.Pieces[Pawn] & whiteAttackedByTwo) & weakForBlack)
		}

	}

	// Black threats
	whiteStronglyProtected := whiteAttackedBy[Pawn] | (whiteAttackedByTwo & ^blackAttackedByTwo)
	whiteDefended := pos.Colours[White] & ^pos.Pieces[Pawn] & whiteStronglyProtected
	if ((pos.Colours[White] & weakForWhite) | whiteDefended) != 0 {
		for fromBB = pos.Colours[White] & (whiteDefended | weakForWhite) & (blackAttackedBy[Knight] | blackAttackedBy[Bishop]) & ^pos.Pieces[Pawn]; fromBB != 0; fromBB &= (fromBB - 1) {
			fromId = BitScan(fromBB)
			threatenedPiece := pos.TypeOnSquare(SquareMask[fromId])
			score -= ThreatByMinor[threatenedPiece]
			if tuning {
				T.ThreatByMinor[threatenedPiece]--
			}
		}

		for fromBB = pos.Colours[White] & (whiteDefended | weakForWhite) & blackAttackedBy[Rook] & ^pos.Pieces[Pawn]; fromBB != 0; fromBB &= (fromBB - 1) {
			fromId = BitScan(fromBB)
			threatenedPiece := pos.TypeOnSquare(SquareMask[fromId])
			score -= ThreatByRook[threatenedPiece]
			if tuning {
				T.ThreatByRook[threatenedPiece]--
			}
		}

		if weakForWhite&pos.Colours[White]&blackAttackedBy[King] != 0 {
			score -= ThreatByKing
			if tuning {
				T.ThreatByKing--
			}
		}

		// Bonus if enemy has a hanging piece
		score -= Hanging *
			Score(PopCount(pos.Colours[White] & ^pos.Pieces[Pawn] & blackAttackedByTwo & weakForWhite))

		if tuning {
			T.Hanging -= PopCount(pos.Colours[White] & ^pos.Pieces[Pawn] & blackAttackedByTwo & weakForWhite)
		}
	}

	// Scale Factor inlined
	scale := SCALE_NORMAL
	if OnlyOne(pos.Colours[Black]&pos.Pieces[Bishop]) &&
		OnlyOne(pos.Colours[White]&pos.Pieces[Bishop]) &&
		OnlyOne(pos.Pieces[Bishop]&WHITE_SQUARES) &&
		(pos.Pieces[Knight]|pos.Pieces[Rook]|pos.Pieces[Queen]) == 0 {
		scale = SCALE_HARD
	} else if (score.End() > 0 && PopCount(pos.Colours[White]) == 2 && (pos.Colours[White]&(pos.Pieces[Bishop]|pos.Pieces[Knight])) != 0) ||
		(score.End() < 0 && PopCount(pos.Colours[Black]) == 2 && (pos.Colours[Black]&(pos.Pieces[Bishop]|pos.Pieces[Knight])) != 0) {
		scale = SCALE_DRAW
	}

	if tuning {
		T.Scale = scale
	}

	// tapering eval
	phase = (phase*256 + (TotalPhase / 2)) / TotalPhase
	result := (int(score.Middle())*(256-phase) + (int(score.End()) * phase * scale / SCALE_NORMAL)) / 256

	if pos.SideToMove == White {
		return result + int(Tempo)
	}
	return -result + int(Tempo)
}

const SCALE_NORMAL = 2
const SCALE_HARD = 1
const SCALE_DRAW = 0

func ScaleFactor(pos *Position, endResult int16) int {
	// OCB without other pieces endgame
	if OnlyOne(pos.Colours[Black]&pos.Pieces[Bishop]) &&
		OnlyOne(pos.Colours[White]&pos.Pieces[Bishop]) &&
		OnlyOne(pos.Pieces[Bishop]&WHITE_SQUARES) &&
		(pos.Pieces[Knight]|pos.Pieces[Rook]|pos.Pieces[Queen]) == 0 {
		return SCALE_HARD
	}
	if (endResult > 0 && PopCount(pos.Colours[White]) == 2 && (pos.Colours[White]&(pos.Pieces[Bishop]|pos.Pieces[Knight])) != 0) ||
		(endResult < 0 && PopCount(pos.Colours[Black]) == 2 && (pos.Colours[Black]&(pos.Pieces[Bishop]|pos.Pieces[Knight])) != 0) {
		return SCALE_DRAW
	}
	return SCALE_NORMAL
}
