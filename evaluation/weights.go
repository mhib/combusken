package evaluation

import (
	. "github.com/mhib/combusken/backend"
)

var PawnValue = S(82, 155)
var KnightValue = S(414, 542)
var BishopValue = S(391, 536)
var RookValue = S(558, 893)
var QueenValue = S(1296, 1671)

// Pawns Square scores
// Bishop Flag, Rank, Row
var PawnScores = [16][7][8]Score{
	{
		{},
		{S(-28, 20), S(-21, 7), S(-28, 16), S(-15, 19), S(-27, 25), S(15, 15), S(39, 5), S(-8, 12)},
		{S(-25, 8), S(-25, 8), S(-25, 9), S(-14, 7), S(-6, 1), S(-1, 7), S(10, 7), S(-15, 11)},
		{S(-24, 19), S(-30, 19), S(-5, 6), S(-1, -2), S(13, -9), S(17, 1), S(5, 15), S(-11, 18)},
		{S(-20, 34), S(-14, 22), S(-11, 11), S(-3, -3), S(29, -13), S(35, 2), S(37, 11), S(5, 25)},
		{S(-20, 71), S(-6, 55), S(2, 45), S(25, 12), S(54, -2), S(60, 20), S(57, 42), S(-6, 69)},
		{S(-10, 23), S(-15, 46), S(-2, 47), S(14, 25), S(15, 14), S(20, 42), S(-101, 81), S(-135, 81)},
	},
	{
		{},
		{S(-16, 25), S(-11, 4), S(-22, 9), S(6, 17), S(-20, 19), S(27, 16), S(35, -5), S(-8, 11)},
		{S(-13, 9), S(-20, 1), S(-5, 8), S(-5, -3), S(10, 13), S(-5, -1), S(19, 11), S(-23, 11)},
		{S(-8, 17), S(-12, 11), S(5, -4), S(29, -1), S(25, -14), S(29, 4), S(4, 12), S(-10, 10)},
		{S(-2, 34), S(13, 0), S(14, 13), S(18, -23), S(46, -9), S(51, -13), S(35, 14), S(-8, 27)},
		{S(-10, 65), S(19, 51), S(27, 8), S(53, 11), S(64, -36), S(93, 4), S(46, -2), S(-2, 62)},
		{S(0, 70), S(-1, -14), S(13, 22), S(23, -29), S(42, 1), S(25, -27), S(-87, 36), S(-127, 20)},
	},
	{
		{},
		{S(-18, 13), S(-14, 12), S(-14, 20), S(-21, 3), S(8, 29), S(11, 5), S(60, 0), S(-1, 17)},
		{S(-21, 7), S(-22, 9), S(-13, -1), S(-1, 10), S(5, -11), S(7, 17), S(4, -3), S(-11, 18)},
		{S(-12, 2), S(-10, 9), S(10, 7), S(19, -10), S(43, -6), S(27, -6), S(16, 11), S(-7, 19)},
		{S(-9, 18), S(1, 16), S(3, -1), S(24, 0), S(42, -17), S(66, 1), S(31, 4), S(4, 31)},
		{S(-10, 63), S(15, 15), S(45, 25), S(42, -32), S(86, -2), S(83, 0), S(46, 29), S(16, 53)},
		{S(-15, -13), S(-2, 40), S(-1, -38), S(29, 3), S(31, -39), S(30, -8), S(-92, 18), S(-122, 56)},
	},
	{
		{},
		{S(-17, 5), S(-12, -6), S(-4, -4), S(-6, 0), S(2, 14), S(10, 6), S(47, 1), S(-4, -5)},
		{S(-10, -7), S(-26, -9), S(-3, -4), S(-2, -11), S(10, 1), S(-8, 5), S(8, -3), S(-12, -12)},
		{S(-2, -11), S(-8, -10), S(12, -21), S(35, -13), S(33, -9), S(24, -1), S(20, -11), S(-4, -11)},
		{S(-19, 14), S(0, -3), S(21, -21), S(20, -18), S(43, -33), S(44, -15), S(27, 6), S(-11, 2)},
		{S(-16, 62), S(11, 24), S(28, 7), S(36, -12), S(70, -22), S(73, 1), S(49, 15), S(-4, 43)},
		{S(-8, 15), S(-9, 11), S(-1, 2), S(24, -19), S(35, -16), S(31, -9), S(-90, 30), S(-127, 47)},
	},
	{
		{},
		{S(-23, 17), S(-10, 15), S(-6, 7), S(-3, 14), S(-22, 22), S(8, 16), S(44, 5), S(-5, 14)},
		{S(-18, 5), S(-23, 4), S(-8, -15), S(3, -5), S(-9, -8), S(-2, 12), S(3, 2), S(-14, 17)},
		{S(-19, 23), S(-3, 10), S(11, -4), S(14, -27), S(34, -25), S(12, -6), S(12, 17), S(-13, 28)},
		{S(-7, 47), S(10, 17), S(4, -8), S(24, -11), S(36, -38), S(43, 9), S(35, 9), S(-2, 46)},
		{S(-23, 92), S(8, 45), S(16, 34), S(52, -27), S(71, -10), S(85, -7), S(61, 41), S(9, 73)},
		{S(12, 52), S(-5, 35), S(7, 13), S(36, 9), S(40, -28), S(37, 3), S(-85, 41), S(-119, 79)},
	},
	{
		{},
		{S(-3, 26), S(2, 13), S(-8, 19), S(-4, 18), S(-10, 23), S(15, 29), S(28, 6), S(0, 13)},
		{S(-5, 14), S(-12, 1), S(0, 9), S(4, -6), S(6, 11), S(-2, 9), S(4, 19), S(-12, 12)},
		{S(-3, 27), S(6, 16), S(18, -4), S(30, -6), S(37, -13), S(25, 9), S(7, 22), S(-9, 21)},
		{S(13, 43), S(18, 9), S(8, 20), S(34, -20), S(52, -19), S(51, -8), S(39, 27), S(-4, 36)},
		{S(-6, 87), S(23, 60), S(39, 19), S(54, 0), S(79, -45), S(85, 11), S(39, 18), S(3, 70)},
		{S(15, 61), S(-15, -8), S(12, 22), S(35, -4), S(33, -13), S(34, -4), S(-84, 30), S(-131, 39)},
	},
	{
		{},
		{S(-11, -14), S(-5, -8), S(-8, -3), S(-10, -13), S(-3, 13), S(-1, -12), S(49, 1), S(-2, -4)},
		{S(-10, -18), S(-15, -7), S(-18, -26), S(-6, -4), S(-14, -26), S(-1, 2), S(-15, -16), S(-11, 5)},
		{S(-12, -7), S(3, -27), S(12, -11), S(8, -48), S(31, -20), S(5, -20), S(13, 1), S(-12, 11)},
		{S(2, -5), S(0, 10), S(5, -37), S(13, -3), S(32, -57), S(68, -11), S(27, -20), S(15, 17)},
		{S(-25, 54), S(22, -13), S(21, 22), S(31, -52), S(71, -18), S(79, -45), S(60, 17), S(27, 12)},
		{S(-8, -55), S(-2, 22), S(-21, -72), S(30, -3), S(26, -75), S(31, -5), S(-102, -43), S(-125, 48)},
	},
	{
		{},
		{S(-6, 0), S(-4, 10), S(-4, 5), S(-7, 5), S(3, 21), S(5, 15), S(37, 11), S(-3, 11)},
		{S(-8, -4), S(-19, 1), S(-6, -2), S(0, 2), S(3, 0), S(-11, 29), S(-1, 0), S(-11, -3)},
		{S(-4, 0), S(-4, 4), S(19, -12), S(29, -19), S(32, -4), S(12, 12), S(7, 9), S(-13, 14)},
		{S(6, 22), S(14, 13), S(7, 4), S(26, -8), S(42, -29), S(54, -4), S(19, 8), S(-13, 18)},
		{S(-11, 59), S(17, 27), S(19, 18), S(43, -10), S(77, -27), S(90, 8), S(31, 3), S(6, 58)},
		{S(-6, 1), S(-11, 9), S(-5, -5), S(30, -11), S(23, -28), S(30, -12), S(-86, 30), S(-129, 41)},
	},
	{
		{},
		{S(-17, 21), S(-11, 7), S(-17, 12), S(-11, 7), S(-12, 12), S(15, 9), S(29, 5), S(-3, 14)},
		{S(-16, 9), S(-18, 0), S(-12, 0), S(-7, -13), S(1, -5), S(-13, 3), S(15, 8), S(-11, 11)},
		{S(-11, 25), S(-9, 20), S(0, -11), S(18, -21), S(22, -29), S(16, 0), S(1, 16), S(-9, 21)},
		{S(-5, 56), S(8, 11), S(-5, 8), S(7, -30), S(40, -11), S(27, -16), S(34, 19), S(-10, 38)},
		{S(11, 80), S(10, 51), S(55, 11), S(37, -13), S(67, -25), S(89, 17), S(83, 28), S(-2, 63)},
		{S(-2, 64), S(-3, 23), S(21, 29), S(33, -26), S(41, 4), S(32, -14), S(-90, 31), S(-126, 55)},
	},
	{
		{},
		{S(-3, 2), S(-6, -14), S(-16, -3), S(15, 10), S(-15, 2), S(31, 5), S(29, -19), S(3, -9)},
		{S(-8, -6), S(-11, -25), S(6, -10), S(-4, -24), S(14, 3), S(-12, -10), S(22, 1), S(-9, -17)},
		{S(-2, 11), S(-8, -6), S(11, -39), S(32, -13), S(26, -46), S(38, -9), S(7, -10), S(-3, -3)},
		{S(8, 24), S(21, -29), S(6, -3), S(11, -62), S(50, -24), S(36, -41), S(35, 8), S(-2, 1)},
		{S(31, 26), S(6, 37), S(42, -27), S(51, 2), S(65, -60), S(100, -4), S(44, -25), S(8, 43)},
		{S(1, 38), S(-17, -62), S(12, 28), S(11, -77), S(33, -17), S(17, -75), S(-91, 28), S(-148, -27)},
	},
	{
		{},
		{S(-6, 22), S(-1, 15), S(-8, 27), S(-11, 19), S(-12, 32), S(14, 11), S(39, 14), S(3, 9)},
		{S(-9, 17), S(-6, 10), S(-8, 13), S(3, 6), S(3, -7), S(-2, 18), S(15, 7), S(-7, 10)},
		{S(-1, 21), S(4, 21), S(17, 6), S(21, -2), S(33, -2), S(27, 1), S(20, 16), S(-6, 18)},
		{S(1, 41), S(24, 24), S(11, 10), S(24, -5), S(45, -21), S(61, 0), S(38, 8), S(3, 36)},
		{S(5, 88), S(22, 45), S(38, 30), S(37, -26), S(88, -7), S(88, 14), S(85, 21), S(11, 54)},
		{S(-4, 17), S(3, 38), S(-3, -16), S(30, 5), S(46, -19), S(35, -18), S(-92, 21), S(-130, 52)},
	},
	{
		{},
		{S(-8, 16), S(-2, 4), S(-7, 11), S(5, 17), S(-14, 27), S(16, 14), S(34, -3), S(-4, 1)},
		{S(-7, 1), S(-12, 4), S(-5, 19), S(0, -1), S(6, 12), S(-9, 5), S(16, 2), S(-15, 3)},
		{S(2, 16), S(5, 5), S(16, -10), S(29, 6), S(32, -18), S(31, 4), S(12, 3), S(-18, 10)},
		{S(13, 17), S(22, 18), S(14, 16), S(29, -45), S(41, -21), S(54, -18), S(22, 15), S(-21, 35)},
		{S(2, 66), S(25, 40), S(31, 6), S(38, -25), S(64, -23), S(86, 0), S(49, 21), S(1, 41)},
		{S(-13, 6), S(-2, 11), S(4, 2), S(24, -12), S(34, -17), S(31, -10), S(-89, 30), S(-126, 40)},
	},
	{
		{},
		{S(-25, 19), S(-14, 18), S(-7, 4), S(-8, 12), S(-31, 3), S(0, 7), S(22, -4), S(3, -9)},
		{S(-21, 19), S(-25, 7), S(-26, 5), S(-5, -10), S(-2, -20), S(-5, -4), S(4, -4), S(-7, -1)},
		{S(-13, 32), S(-10, 18), S(7, -23), S(7, -26), S(22, -50), S(21, -21), S(-5, 17), S(-4, 10)},
		{S(-4, 34), S(10, 20), S(1, 10), S(21, -26), S(29, -18), S(36, -33), S(27, 16), S(8, 38)},
		{S(3, 89), S(19, 48), S(25, 23), S(47, -2), S(63, -17), S(85, 5), S(58, 24), S(5, 67)},
		{S(-5, 23), S(-9, 11), S(1, 10), S(20, -19), S(30, -17), S(30, -14), S(-87, 29), S(-126, 46)},
	},
	{
		{},
		{S(-8, 26), S(-1, 11), S(-11, 10), S(-2, 9), S(-17, 13), S(12, 30), S(14, 4), S(-1, 4)},
		{S(-7, 17), S(-20, 0), S(0, 9), S(3, -27), S(5, 10), S(-2, -4), S(5, 11), S(-14, -1)},
		{S(0, 30), S(-2, 19), S(16, -19), S(29, -4), S(31, -36), S(25, 5), S(5, 6), S(-8, 7)},
		{S(14, 52), S(16, 8), S(13, 16), S(30, -50), S(49, -10), S(33, -25), S(40, 25), S(-7, 35)},
		{S(15, 70), S(31, 56), S(32, 20), S(43, -1), S(41, -14), S(84, 2), S(41, 17), S(11, 59)},
		{S(4, 31), S(-8, 11), S(4, 8), S(26, -18), S(33, -16), S(22, -19), S(-94, 25), S(-124, 47)},
	},
	{
		{},
		{S(-18, 18), S(-2, 11), S(-8, 24), S(-4, 6), S(-21, 25), S(-3, 4), S(26, 21), S(0, 13)},
		{S(-16, 5), S(-9, 9), S(-9, -20), S(1, -1), S(-10, -16), S(-5, 13), S(3, -6), S(-8, 14)},
		{S(-11, 21), S(1, 15), S(15, 2), S(17, -33), S(34, -16), S(14, -12), S(13, 7), S(-11, 14)},
		{S(-5, 44), S(29, 27), S(2, -2), S(25, 1), S(36, -44), S(61, -12), S(31, -2), S(1, 41)},
		{S(8, 95), S(16, 26), S(47, 31), S(26, -33), S(80, -13), S(71, -1), S(53, 26), S(15, 51)},
		{S(-12, 6), S(-7, 19), S(-9, -4), S(24, -7), S(29, -24), S(31, -6), S(-87, 26), S(-130, 46)},
	},
	{
		{},
		{S(-9, 17), S(-1, 23), S(-13, 23), S(-9, 10), S(-18, 12), S(1, 27), S(15, 13), S(-7, 10)},
		{S(-9, 12), S(-15, 12), S(-7, 13), S(-2, 1), S(-1, 8), S(-15, 34), S(-5, 12), S(-18, 9)},
		{S(-2, 20), S(1, 24), S(15, 6), S(22, 0), S(28, -12), S(16, 10), S(1, 19), S(-21, 20)},
		{S(8, 51), S(25, 20), S(4, 32), S(27, -28), S(37, -14), S(46, -10), S(22, 21), S(-15, 34)},
		{S(7, 85), S(22, 46), S(26, 37), S(45, -20), S(64, -30), S(67, -3), S(34, 15), S(-4, 55)},
		{S(-33, 20), S(-28, 4), S(-5, 5), S(21, -16), S(31, -22), S(31, -9), S(-81, 28), S(-134, 40)},
	},
}

// Piece Square Values
var PieceScores = [Queen + 1][8][4]Score{
	{},
	{ // knight
		{S(-94, -35), S(-26, -40), S(-32, -30), S(-24, -7)},
		{S(-24, -19), S(-25, -16), S(-18, -32), S(-20, -12)},
		{S(-19, -45), S(-5, -22), S(-13, -14), S(0, 7)},
		{S(-12, -1), S(14, 0), S(4, 22), S(2, 29)},
		{S(6, 1), S(3, 3), S(21, 25), S(13, 40)},
		{S(-28, -9), S(-18, 6), S(5, 31), S(17, 29)},
		{S(5, -28), S(-12, -10), S(32, -27), S(40, 9)},
		{S(-212, -45), S(-91, -11), S(-145, 26), S(-8, -11)},
	},
	{ // Bishop
		{S(26, -38), S(31, -3), S(10, -4), S(7, 8)},
		{S(32, -48), S(23, -38), S(32, -10), S(15, 2)},
		{S(23, -8), S(41, 1), S(12, -10), S(29, 18)},
		{S(15, -8), S(22, 7), S(24, 21), S(36, 25)},
		{S(-5, 12), S(24, 15), S(15, 24), S(38, 33)},
		{S(-2, 11), S(-3, 28), S(7, 8), S(23, 29)},
		{S(-34, 28), S(-19, 9), S(16, 25), S(-1, 34)},
		{S(-42, 2), S(-38, 31), S(-111, 45), S(-101, 52)},
	},
	{ // Rook
		{S(-19, -26), S(-11, -15), S(-5, -13), S(-1, -24)},
		{S(-62, -9), S(-16, -29), S(-9, -28), S(-19, -24)},
		{S(-36, -17), S(-11, -9), S(-22, -9), S(-20, -13)},
		{S(-29, 3), S(-19, 15), S(-17, 16), S(-14, 8)},
		{S(-15, 16), S(10, 14), S(21, 16), S(30, 12)},
		{S(-16, 27), S(45, 10), S(42, 22), S(53, 8)},
		{S(-3, 33), S(2, 36), S(36, 28), S(39, 38)},
		{S(35, 37), S(55, 36), S(28, 43), S(26, 39)},
	},
	{ // Queen
		{S(7, -103), S(7, -87), S(8, -91), S(14, -57)},
		{S(4, -84), S(14, -75), S(18, -87), S(7, -47)},
		{S(1, -62), S(13, -28), S(4, 5), S(-2, -8)},
		{S(-1, -20), S(6, 9), S(-1, 35), S(-17, 65)},
		{S(12, -9), S(4, 41), S(-12, 57), S(-24, 90)},
		{S(5, 5), S(13, 4), S(-11, 59), S(0, 56)},
		{S(-29, 10), S(-75, 55), S(-17, 56), S(-40, 88)},
		{S(-10, 13), S(25, 8), S(16, 44), S(22, 47)},
	},
}

var KingScores = [8][8]Score{
	{S(151, -3), S(138, 31), S(56, 73), S(47, 81), S(79, 53), S(42, 76), S(148, 29), S(163, -21)},
	{S(162, 37), S(110, 54), S(43, 93), S(40, 113), S(30, 124), S(55, 89), S(129, 50), S(162, 37)},
	{S(112, 25), S(161, 40), S(98, 80), S(94, 117), S(76, 126), S(113, 75), S(153, 39), S(109, 34)},
	{S(100, 14), S(257, 28), S(168, 84), S(92, 131), S(89, 136), S(158, 80), S(237, 32), S(84, 28)},
	{S(101, 44), S(291, 42), S(180, 96), S(110, 130), S(98, 135), S(187, 90), S(266, 41), S(87, 43)},
	{S(120, 40), S(310, 61), S(248, 91), S(183, 108), S(180, 113), S(250, 84), S(296, 59), S(117, 42)},
	{S(116, -6), S(214, 61), S(195, 79), S(151, 104), S(144, 104), S(187, 75), S(210, 70), S(115, -2)},
	{S(158, -105), S(291, -27), S(175, 11), S(162, 69), S(159, 67), S(172, 17), S(290, -15), S(156, -97)},
}

var PawnsConnected = [7][4]Score{
	{S(0, 0), S(0, 0), S(0, 0), S(0, 0)},
	{S(-2, -13), S(8, 3), S(3, 0), S(17, 17)},
	{S(12, 6), S(22, 4), S(24, 10), S(27, 16)},
	{S(10, 5), S(20, 8), S(11, 9), S(17, 20)},
	{S(5, 17), S(16, 25), S(31, 29), S(27, 22)},
	{S(38, 26), S(30, 64), S(78, 64), S(89, 81)},
	{S(176, 32), S(297, 22), S(291, 40), S(344, 47)},
}

var MobilityBonus = [...][32]Score{
	{S(-61, -132), S(-43, -78), S(-28, -30), S(-20, -5), S(-13, 9), S(-8, 25), // Knights
		S(1, 28), S(9, 23), S(20, 9)},
	{S(0, -136), S(7, -64), S(15, -23), S(23, 0), S(30, 16), S(35, 33), // Bishops
		S(37, 41), S(37, 44), S(37, 47), S(41, 47), S(44, 43), S(55, 34),
		S(75, 37), S(81, 15)},
	{S(-127, -147), S(-16, -33), S(-4, 18), S(-5, 47), S(-1, 59), S(1, 71), // Rooks
		S(2, 80), S(7, 84), S(11, 88), S(15, 93), S(19, 98), S(21, 101),
		S(28, 101), S(41, 90), S(97, 54)},
	{S(-413, -159), S(-122, -138), S(-38, -175), S(-20, -117), S(-7, -79), S(-7, -5), // Queens
		S(-3, 22), S(1, 40), S(5, 54), S(7, 66), S(11, 72), S(15, 73),
		S(16, 74), S(17, 77), S(21, 70), S(18, 72), S(17, 66), S(15, 63),
		S(20, 50), S(27, 34), S(38, 13), S(40, -5), S(40, -22), S(51, -51),
		S(14, -38), S(-86, -13), S(127, -127), S(49, -81)},
}

var PassedFriendlyDistance = [8]Score{
	S(0, 0), S(-11, 41), S(-15, 24), S(-14, 9),
	S(-7, -7), S(-3, -16), S(16, -28), S(0, -40),
}

var PassedEnemyDistance = [8]Score{
	S(0, 0), S(-125, -62), S(-14, -9), S(1, 11),
	S(13, 22), S(13, 31), S(10, 38), S(17, 44),
}

// PassedRank[Rank] contains a bonus according to the rank of a passed pawn
var PassedRank = [7]Score{S(0, 0), S(-12, -30), S(-20, -12), S(-9, 28), S(29, 70), S(54, 148), S(183, 233)}

// PassedFile[File] contains a bonus according to the file of a passed pawn
var PassedFile = [8]Score{S(-11, 31), S(-13, 32), S(-7, 12), S(-4, -5),
	S(-10, 0), S(-19, 10), S(-18, 23), S(11, 9),
}

var PassedStacked = [8]Score{S(0, 0), S(-13, -55), S(-19, -36), S(-30, -58), S(-6, -84), S(28, -207), S(0, 0), S(0, 0)}

var Isolated = S(-7, -16)
var Doubled = S(-9, -35)
var Backward = S(8, -1)
var BackwardOpen = S(-5, -16)

var BishopPair = S(18, 75)
var BishopRammedPawns = S(-8, -24)

var BishopOutpostUndefendedBonus = S(20, -2)
var BishopOutpostDefendedBonus = S(51, 13)

var LongDiagonalBishop = S(23, 25)

var KnightOutpostUndefendedBonus = S(20, -21)
var KnightOutpostDefendedBonus = S(45, 13)

var DistantKnight = [4]Score{S(-13, -3), S(-18, -12), S(-27, -18), S(-53, -23)}

var MinorBehindPawn = S(8, 33)

var Tempo int16 = 15

// Rook on semiopen, open file
var RookOnFile = [2]Score{S(10, 18), S(36, 11)}
var RookOnQueenFile = S(5, 31)

var KingDefenders = [12]Score{
	S(-7, -7), S(2, -9), S(5, -9), S(6, -7),
	S(7, -4), S(8, 5), S(10, 11), S(11, 11),
	S(9, 12), S(11, -69), S(-4, -10), S(11, 0),
}

var KingShelter = [2][8][8]Score{
	{{S(-27, 10), S(9, -12), S(18, 0), S(21, 4),
		S(27, -12), S(20, -7), S(16, -40), S(-94, 38)},
		{S(11, 7), S(21, -6), S(-5, 9), S(-9, 6),
			S(2, -2), S(26, -1), S(40, -16), S(-51, 25)},
		{S(16, 3), S(3, -1), S(-29, 2), S(-26, 1),
			S(-3, -16), S(-1, 4), S(19, -9), S(-36, 3)},
		{S(-22, 19), S(11, -1), S(-5, -10), S(-1, -10),
			S(7, -26), S(3, -9), S(24, 14), S(-15, -2)},
		{S(-33, 18), S(-15, 9), S(-32, 8), S(-20, 6),
			S(1, -9), S(-23, 0), S(8, -7), S(-27, 4)},
		{S(41, -19), S(23, -14), S(4, -12), S(4, -21),
			S(13, -28), S(26, -20), S(33, -19), S(-12, -2)},
		{S(16, -3), S(-8, -8), S(-27, 2), S(-17, 1),
			S(-7, -17), S(9, -5), S(14, -18), S(-24, 19)},
		{S(-41, 11), S(-26, 0), S(-17, 26), S(-20, 23),
			S(-22, 16), S(-9, 11), S(-10, -15), S(-75, 50)}},
	{{S(48, 45), S(-42, -11), S(-28, 3), S(-36, -5),
		S(-50, -15), S(-16, 13), S(-59, -2), S(-96, 41)},
		{S(150, 12), S(2, -4), S(-3, 5), S(-26, 19),
			S(-7, -12), S(17, -4), S(9, -2), S(-86, 33)},
		{S(5, 26), S(39, 0), S(13, -1), S(21, -7),
			S(29, -6), S(14, -4), S(42, 0), S(-27, 10)},
		{S(4, 38), S(-11, 20), S(-13, 15), S(-24, 21),
			S(-1, 7), S(-10, 7), S(-9, 19), S(-42, -4)},
		{S(-10, 29), S(-1, 8), S(-5, 3), S(-11, 1),
			S(-4, 0), S(2, -14), S(13, -5), S(-31, -2)},
		{S(32, 12), S(-4, 4), S(-13, 9), S(-6, 7),
			S(-1, -1), S(-37, 5), S(-3, -5), S(-43, 11)},
		{S(40, 3), S(-5, -14), S(-7, -12), S(-25, -7),
			S(-8, -20), S(2, -20), S(2, -21), S(-74, 22)},
		{S(-9, -3), S(-16, -30), S(-8, -18), S(-3, -24),
			S(-7, -35), S(1, -27), S(-4, -52), S(-70, 29)}},
}

var KingStorm = [2][4][8]Score{
	{{S(10, 10), S(3, 17), S(8, 9), S(4, 10),
		S(2, 7), S(10, 2), S(7, 9), S(-6, -8)},
		{S(23, 2), S(17, 6), S(19, 5), S(5, 12),
			S(19, 1), S(27, -10), S(19, -9), S(-7, -7)},
		{S(22, 10), S(3, 7), S(5, 12), S(1, 13),
			S(4, 8), S(11, 1), S(8, -11), S(-1, 1)},
		{S(12, 3), S(9, -1), S(12, -6), S(1, -11),
			S(-2, -7), S(15, -8), S(8, -5), S(-10, -1)}},
	{{S(0, 0), S(16, -1), S(-19, 11), S(9, -6),
		S(-3, 19), S(-13, 31), S(6, 48), S(-2, -4)},
		{S(0, 0), S(-42, -13), S(5, -3), S(46, -5),
			S(4, -2), S(-7, -8), S(16, 57), S(-7, -1)},
		{S(0, 0), S(-17, -8), S(-9, -4), S(16, -1),
			S(6, -7), S(-10, -19), S(46, -45), S(-7, 7)},
		{S(0, 0), S(-15, -24), S(-16, -23), S(-7, -11),
			S(2, -26), S(-1, -52), S(6, 4), S(-11, -2)}},
}
var KingOnPawnlessFlank = S(21, -73)

var KingSafetyAttacksWeights = [King + 1]int16{0, -3, -7, -4, 4, 0}
var KingSafetyAttackValue int16 = 124
var KingSafetyWeakSquares int16 = 44
var KingSafetyFriendlyPawns int16 = -35
var KingSafetyNoEnemyQueens int16 = -176
var KingSafetySafeQueenCheck int16 = 90
var KingSafetySafeRookCheck int16 = 73
var KingSafetySafeBishopCheck int16 = 51
var KingSafetySafeKnightCheck int16 = 112
var KingSafetyAdjustment int16 = -12

var Hanging = S(37, 13)
var ThreatByKing = S(-8, 33)
var ThreatByMinor = [King + 1]Score{S(0, 0), S(21, 39), S(17, 36), S(32, 27), S(29, 30), S(7, 22)}
var ThreatByRook = [King + 1]Score{S(0, 0), S(-3, 11), S(-1, 16), S(-6, -10), S(34, 8), S(30, -1)}

// This weights are from black piece on black square perspective
var RookBishopExistence = [16]Score{
	S(19, -22), S(-1, -8), S(0, -9), S(-1, 6), S(-5, -7), S(-3, 26), S(-11, 9), S(0, 7), S(-2, -7), S(-12, 11), S(-2, 27), S(0, 10), S(-4, -12), S(-1, 10), S(2, 9), S(-15, 38),
}
var QueenBishopExistence = [16]Score{
	S(10, -6), S(-2, 5), S(-1, -9), S(-11, -5), S(-2, 4), S(2, 31), S(-8, -7), S(-6, -12), S(3, 9), S(-4, 23), S(4, 27), S(-2, -3), S(-7, -7), S(1, 12), S(-4, 3), S(-24, -10),
}
var KingBishopExistence = [16]Score{
	S(0, 0), S(-3, 2), S(-1, 2), S(-5, -7), S(0, -7), S(-2, -6), S(-1, -4), S(-1, -2), S(3, 1), S(1, 4), S(2, 6), S(1, 8), S(5, 7), S(-1, -4), S(3, -2), S(0, 0),
}
