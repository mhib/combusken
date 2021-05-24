package evaluation

import (
	. "github.com/mhib/combusken/backend"
)

var PawnValue = S(85, 154)
var KnightValue = S(416, 547)
var BishopValue = S(389, 539)
var RookValue = S(561, 908)
var QueenValue = S(1244, 1745)

// Pawns Square scores
// Bishop Flag, Rank, Row
var PawnScores = [16][7][8]Score{
	{
		{},
		{S(-17, 15), S(-12, 4), S(-23, 12), S(-8, 13), S(-23, 21), S(16, 12), S(33, 5), S(-12, 11)},
		{S(-15, 6), S(-17, 5), S(-19, 5), S(-7, -1), S(-3, -5), S(-1, 4), S(7, 4), S(-19, 10)},
		{S(-14, 17), S(-20, 15), S(1, 0), S(7, -9), S(17, -16), S(17, -5), S(1, 12), S(-15, 17)},
		{S(-7, 29), S(-4, 17), S(-3, 3), S(4, -11), S(30, -18), S(32, -4), S(31, 9), S(1, 23)},
		{S(1, 51), S(10, 33), S(20, 22), S(46, -9), S(63, -18), S(64, 11), S(64, 28), S(3, 54)},
		{S(11, -7), S(-23, 24), S(0, 19), S(15, -3), S(-22, -1), S(-19, 20), S(-129, 58), S(-146, 50)},
	},
	{
		{},
		{S(-12, 27), S(-7, 6), S(-21, 10), S(9, 17), S(-19, 17), S(28, 15), S(28, -2), S(-12, 12)},
		{S(-8, 11), S(-17, 2), S(-2, 7), S(-2, -4), S(10, 10), S(-4, -3), S(15, 11), S(-25, 12)},
		{S(-4, 21), S(-6, 12), S(6, -4), S(33, -4), S(26, -17), S(28, 1), S(0, 11), S(-13, 12)},
		{S(6, 34), S(16, 3), S(20, 7), S(19, -21), S(45, -13), S(47, -13), S(28, 14), S(-11, 32)},
		{S(-12, 82), S(36, 26), S(24, 27), S(69, -17), S(60, -13), S(104, -12), S(38, 24), S(-4, 52)},
		{S(-21, 69), S(10, 5), S(28, 7), S(17, -8), S(57, -10), S(16, -22), S(-83, 36), S(-108, 28)},
	},
	{
		{},
		{S(-18, 17), S(-14, 15), S(-16, 23), S(-19, 3), S(5, 29), S(9, 5), S(50, 4), S(-9, 20)},
		{S(-21, 12), S(-22, 12), S(-13, 0), S(0, 10), S(4, -12), S(3, 16), S(-1, -3), S(-17, 19)},
		{S(-12, 7), S(-8, 12), S(9, 7), S(20, -12), S(42, -9), S(23, -7), S(7, 12), S(-13, 22)},
		{S(-8, 25), S(3, 16), S(2, 2), S(26, -5), S(37, -15), S(60, -2), S(21, 9), S(-4, 34)},
		{S(1, 45), S(9, 39), S(65, -6), S(44, -10), S(98, -26), S(69, 18), S(25, 26), S(6, 67)},
		{S(-32, 21), S(-18, 42), S(12, -19), S(21, 6), S(27, -29), S(30, -22), S(-91, 27), S(-101, 29)},
	},
	{
		{},
		{S(-19, 6), S(-13, -5), S(-4, -8), S(2, -18), S(7, -1), S(10, -1), S(44, -9), S(-2, -17)},
		{S(-12, -4), S(-29, -7), S(-2, -8), S(-2, -12), S(10, -5), S(-8, -3), S(6, -11), S(-10, -22)},
		{S(-2, -10), S(-7, -12), S(13, -24), S(37, -19), S(33, -14), S(22, -8), S(22, -23), S(1, -24)},
		{S(-17, 18), S(-1, -5), S(34, -44), S(19, -16), S(42, -42), S(45, -25), S(17, 5), S(-5, -12)},
		{S(-27, 73), S(13, 16), S(40, -25), S(24, 1), S(83, -32), S(40, 7), S(38, 13), S(-16, 31)},
		{S(-15, 32), S(-28, 33), S(-27, 9), S(17, -36), S(52, 0), S(44, -8), S(-91, 32), S(-120, 71)},
	},
	{
		{},
		{S(-23, 20), S(-11, 19), S(-7, 8), S(-4, 14), S(-23, 20), S(6, 17), S(34, 8), S(-11, 16)},
		{S(-18, 9), S(-21, 6), S(-9, -14), S(3, -7), S(-11, -10), S(-5, 11), S(-1, 2), S(-19, 19)},
		{S(-18, 27), S(-2, 11), S(12, -5), S(14, -29), S(33, -27), S(10, -9), S(5, 18), S(-17, 29)},
		{S(-5, 48), S(11, 18), S(5, -9), S(24, -14), S(34, -40), S(36, 8), S(28, 9), S(-8, 47)},
		{S(-23, 88), S(5, 42), S(8, 32), S(64, -24), S(72, -15), S(83, -7), S(57, 38), S(4, 73)},
		{S(37, 42), S(-27, 43), S(5, 25), S(42, 18), S(73, -34), S(52, 17), S(-78, 60), S(-125, 98)},
	},
	{
		{},
		{S(-2, 31), S(1, 19), S(-9, 21), S(-3, 20), S(-10, 24), S(15, 30), S(24, 9), S(-2, 16)},
		{S(-3, 20), S(-10, 5), S(1, 11), S(5, -4), S(6, 11), S(-1, 9), S(5, 19), S(-12, 15)},
		{S(-1, 34), S(9, 18), S(19, -1), S(31, -6), S(38, -13), S(24, 8), S(6, 23), S(-9, 24)},
		{S(17, 47), S(20, 15), S(11, 18), S(34, -14), S(50, -20), S(49, -6), S(38, 28), S(-6, 42)},
		{S(-5, 99), S(33, 45), S(38, 35), S(69, -17), S(78, -18), S(85, 7), S(29, 45), S(-3, 71)},
		{S(19, 72), S(-36, 21), S(14, 41), S(52, 19), S(22, -2), S(51, 31), S(-52, 21), S(-132, 47)},
	},
	{
		{},
		{S(-7, -17), S(1, -10), S(-4, -7), S(3, -26), S(2, 4), S(4, -19), S(47, -1), S(-1, -9)},
		{S(-5, -18), S(-8, -10), S(-13, -30), S(1, -11), S(-9, -34), S(5, -7), S(-11, -23), S(-8, -1)},
		{S(-6, -9), S(11, -32), S(19, -18), S(14, -56), S(39, -30), S(9, -30), S(14, -5), S(-9, 6)},
		{S(8, -5), S(10, 1), S(10, -40), S(22, -16), S(35, -63), S(75, -26), S(28, -26), S(18, 10)},
		{S(-7, 23), S(45, -6), S(23, -8), S(45, -33), S(80, -54), S(94, -50), S(85, -13), S(47, 9)},
		{S(86, -82), S(5, 23), S(19, -83), S(27, -11), S(76, -106), S(21, 3), S(-38, -119), S(-111, 21)},
	},
	{
		{},
		{S(-6, 3), S(-5, 15), S(-5, 8), S(-4, 4), S(3, 20), S(5, 16), S(31, 16), S(-3, 11)},
		{S(-7, 0), S(-19, 7), S(-5, 0), S(0, 6), S(2, -2), S(-11, 29), S(-2, 1), S(-10, -3)},
		{S(-3, 6), S(-2, 6), S(19, -9), S(30, -22), S(33, -5), S(11, 10), S(4, 12), S(-12, 15)},
		{S(9, 27), S(15, 18), S(8, 7), S(26, -4), S(41, -29), S(53, -4), S(14, 12), S(-13, 21)},
		{S(-6, 56), S(26, 30), S(15, 22), S(44, 14), S(85, -33), S(87, 17), S(21, -10), S(-4, 78)},
		{S(12, -6), S(-41, 20), S(-21, -2), S(58, -8), S(-26, -55), S(33, -39), S(-64, 39), S(-135, 18)},
	},
	{
		{},
		{S(-15, 23), S(-8, 8), S(-17, 13), S(-8, 6), S(-14, 12), S(12, 10), S(23, 7), S(-12, 17)},
		{S(-12, 11), S(-17, 1), S(-10, -1), S(-5, -14), S(0, -7), S(-15, 1), S(6, 10), S(-18, 13)},
		{S(-9, 27), S(-6, 20), S(1, -13), S(20, -24), S(22, -33), S(12, -2), S(-5, 15), S(-16, 23)},
		{S(2, 54), S(10, 11), S(0, 3), S(8, -32), S(37, -14), S(23, -19), S(25, 19), S(-18, 40)},
		{S(17, 72), S(12, 44), S(64, 5), S(45, -18), S(69, -22), S(86, 13), S(99, 21), S(-8, 60)},
		{S(-30, 57), S(-12, 38), S(44, 17), S(45, -19), S(38, 22), S(42, -16), S(-86, 20), S(-123, 57)},
	},
	{
		{},
		{S(5, -1), S(4, -19), S(-9, -8), S(27, -4), S(-8, -7), S(38, -2), S(33, -25), S(6, -13)},
		{S(2, -9), S(-4, -28), S(15, -18), S(5, -32), S(21, -6), S(-4, -20), S(24, -6), S(-4, -23)},
		{S(6, 9), S(3, -13), S(18, -45), S(42, -23), S(34, -56), S(43, -18), S(12, -19), S(1, -7)},
		{S(23, 15), S(29, -32), S(21, -20), S(18, -67), S(58, -38), S(41, -49), S(36, 1), S(5, -4)},
		{S(43, 32), S(21, 6), S(59, -21), S(74, -39), S(83, -51), S(136, -46), S(67, -19), S(28, 8)},
		{S(2, 22), S(36, -71), S(7, 46), S(53, -79), S(17, -31), S(64, -119), S(-100, -9), S(-118, -69)},
	},
	{
		{},
		{S(-7, 28), S(-2, 20), S(-9, 32), S(-8, 21), S(-16, 35), S(11, 14), S(33, 18), S(-3, 13)},
		{S(-10, 23), S(-7, 15), S(-7, 16), S(4, 7), S(1, -6), S(-4, 18), S(7, 9), S(-11, 12)},
		{S(-2, 26), S(4, 26), S(17, 7), S(21, -1), S(33, -4), S(24, 3), S(14, 17), S(-11, 22)},
		{S(2, 50), S(23, 27), S(10, 14), S(25, -7), S(41, -18), S(55, -2), S(30, 12), S(-3, 39)},
		{S(11, 77), S(19, 68), S(47, 11), S(34, 9), S(94, -22), S(78, 29), S(95, 6), S(4, 66)},
		{S(-10, 48), S(-8, 67), S(-13, 14), S(28, 20), S(80, -14), S(66, -50), S(-89, 24), S(-155, 76)},
	},
	{
		{},
		{S(-9, 23), S(-1, 9), S(-7, 14), S(9, 18), S(-15, 29), S(17, 14), S(32, -2), S(-4, 0)},
		{S(-7, 8), S(-14, 8), S(-3, 20), S(2, -2), S(6, 12), S(-9, 5), S(12, 3), S(-14, 1)},
		{S(1, 23), S(7, 8), S(17, -10), S(31, 7), S(35, -21), S(28, 4), S(10, 3), S(-17, 11)},
		{S(15, 23), S(22, 24), S(16, 20), S(30, -45), S(40, -20), S(53, -24), S(17, 22), S(-21, 39)},
		{S(2, 80), S(25, 54), S(34, 8), S(44, -17), S(62, -11), S(90, -13), S(40, 42), S(0, 34)},
		{S(-16, -8), S(3, 27), S(9, 5), S(21, 5), S(34, -7), S(37, 0), S(-88, 35), S(-117, 35)},
	},
	{
		{},
		{S(-28, 26), S(-16, 21), S(-8, 1), S(-5, 4), S(-27, -12), S(-1, 5), S(20, -7), S(0, -10)},
		{S(-23, 25), S(-26, 13), S(-26, 4), S(-6, -9), S(-1, -29), S(-2, -11), S(-1, -5), S(-7, -2)},
		{S(-14, 37), S(-12, 21), S(5, -23), S(8, -34), S(24, -58), S(22, -28), S(-10, 18), S(-5, 8)},
		{S(-7, 41), S(3, 32), S(0, 12), S(21, -30), S(28, -24), S(41, -51), S(20, 17), S(4, 40)},
		{S(-21, 131), S(-2, 74), S(4, 47), S(48, 15), S(63, -2), S(84, 2), S(70, 31), S(-17, 106)},
		{S(-15, 49), S(-13, 39), S(-16, 36), S(-1, -8), S(15, -9), S(30, -45), S(-52, 33), S(-93, 78)},
	},
	{
		{},
		{S(-9, 33), S(-1, 15), S(-12, 13), S(-1, 9), S(-19, 13), S(12, 27), S(14, 2), S(-2, 2)},
		{S(-7, 23), S(-19, 5), S(1, 9), S(4, -28), S(4, 9), S(-1, -7), S(3, 9), S(-13, -3)},
		{S(-1, 37), S(0, 21), S(17, -20), S(28, -6), S(33, -39), S(25, 2), S(5, 3), S(-8, 6)},
		{S(16, 57), S(15, 12), S(16, 14), S(30, -49), S(48, -13), S(33, -29), S(37, 23), S(-7, 35)},
		{S(12, 89), S(29, 75), S(30, 39), S(40, 16), S(32, 33), S(89, -14), S(35, 35), S(5, 66)},
		{S(13, 99), S(-28, 48), S(-3, 30), S(33, -15), S(30, -5), S(-11, -50), S(-120, -9), S(-106, 83)},
	},
	{
		{},
		{S(-21, 24), S(-4, 16), S(-10, 28), S(-3, 4), S(-26, 28), S(-5, 5), S(22, 23), S(-6, 16)},
		{S(-18, 10), S(-10, 12), S(-9, -19), S(1, 0), S(-12, -16), S(-7, 10), S(-2, -6), S(-13, 15)},
		{S(-14, 27), S(0, 18), S(15, 3), S(16, -35), S(33, -18), S(12, -14), S(8, 7), S(-15, 16)},
		{S(-6, 52), S(26, 30), S(2, -1), S(24, 0), S(33, -44), S(57, -20), S(26, -3), S(-5, 44)},
		{S(1, 111), S(17, 44), S(49, 33), S(24, 6), S(92, -17), S(62, 11), S(48, 31), S(9, 61)},
		{S(-13, 14), S(-17, 76), S(-50, 25), S(10, 26), S(8, -24), S(28, 10), S(-76, 21), S(-134, 64)},
	},
	{
		{},
		{S(-12, 26), S(-4, 33), S(-16, 32), S(-8, 16), S(-22, 19), S(0, 31), S(12, 20), S(-9, 15)},
		{S(-11, 21), S(-17, 21), S(-8, 20), S(-4, 9), S(-4, 13), S(-17, 38), S(-9, 18), S(-19, 14)},
		{S(-4, 30), S(0, 31), S(13, 12), S(20, 6), S(28, -9), S(13, 13), S(-2, 25), S(-21, 26)},
		{S(8, 60), S(23, 29), S(4, 40), S(25, -18), S(34, -8), S(43, -8), S(18, 28), S(-16, 41)},
		{S(6, 99), S(20, 73), S(17, 75), S(45, 8), S(60, -15), S(64, -1), S(28, 27), S(-9, 70)},
		{S(-87, 107), S(-74, 21), S(-36, 47), S(5, 1), S(15, -41), S(32, 3), S(-60, 16), S(-162, 25)},
	},
}

// Piece Square Values
var PieceScores = [King + 1][8][8]Score{
	{},
	{ // knight
		{S(-81, -87), S(-23, -46), S(-46, -32), S(-21, -6), S(-23, -13), S(-26, -34), S(-28, -32), S(-89, -44)},
		{S(-31, -31), S(-19, -19), S(-15, -37), S(-15, -13), S(-25, -15), S(-17, -33), S(-26, -18), S(-22, -22)},
		{S(-13, -46), S(-7, -19), S(-7, -12), S(-2, 8), S(0, 6), S(-16, -14), S(-4, -24), S(-23, -46)},
		{S(-11, -4), S(-7, 6), S(9, 21), S(-3, 31), S(11, 28), S(2, 22), S(24, -5), S(-12, 1)},
		{S(-2, 9), S(-3, 7), S(11, 33), S(12, 49), S(11, 42), S(29, 25), S(8, 7), S(9, 5)},
		{S(-53, 8), S(-27, 16), S(-2, 40), S(9, 39), S(29, 28), S(29, 28), S(-13, 5), S(-20, -4)},
		{S(1, -19), S(-12, -5), S(26, -17), S(34, 17), S(41, 9), S(63, -38), S(-21, -8), S(15, -33)},
		{S(-201, -52), S(-98, -9), S(-153, 26), S(-30, -6), S(10, -7), S(-127, 32), S(-75, -10), S(-187, -68)},
	},
	{ // Bishop
		{S(34, -39), S(52, -12), S(7, 0), S(12, 1), S(-1, 10), S(15, -10), S(13, -11), S(22, -40)},
		{S(46, -31), S(20, -35), S(47, -10), S(8, 5), S(17, -2), S(17, -14), S(32, -42), S(23, -64)},
		{S(21, -7), S(52, 1), S(7, -5), S(35, 15), S(17, 20), S(16, -15), S(32, -2), S(34, -8)},
		{S(15, -5), S(21, 9), S(32, 18), S(22, 27), S(44, 22), S(15, 21), S(27, 5), S(14, -8)},
		{S(-13, 19), S(30, 15), S(0, 30), S(44, 33), S(24, 36), S(29, 22), S(19, 18), S(2, 15)},
		{S(-1, 15), S(-14, 33), S(13, 6), S(7, 34), S(28, 30), S(-3, 16), S(10, 29), S(-10, 20)},
		{S(-30, 31), S(-6, 8), S(13, 25), S(-11, 38), S(-6, 35), S(17, 26), S(-19, 4), S(-33, 30)},
		{S(-22, -1), S(-64, 38), S(-127, 48), S(-108, 58), S(-107, 51), S(-100, 45), S(21, 16), S(-40, 0)},
	},
	{ // Rook
		{S(-23, -27), S(-18, -17), S(-12, -15), S(-2, -26), S(-4, -29), S(-5, -17), S(-3, -26), S(-14, -43)},
		{S(-62, -22), S(-26, -31), S(-22, -25), S(-21, -28), S(-22, -30), S(-1, -40), S(2, -44), S(-60, -19)},
		{S(-44, -17), S(-28, -2), S(-35, -5), S(-22, -13), S(-24, -12), S(-17, -13), S(7, -18), S(-28, -23)},
		{S(-40, 8), S(-33, 23), S(-29, 24), S(-15, 12), S(-22, 10), S(-15, 13), S(-7, 9), S(-17, 1)},
		{S(-32, 30), S(-3, 28), S(7, 28), S(36, 17), S(11, 20), S(25, 15), S(22, 9), S(3, 17)},
		{S(-35, 42), S(24, 24), S(20, 33), S(44, 17), S(45, 15), S(77, 19), S(93, -4), S(13, 25)},
		{S(-25, 50), S(-17, 50), S(14, 43), S(38, 42), S(17, 47), S(74, 12), S(22, 28), S(38, 22)},
		{S(12, 45), S(30, 46), S(0, 55), S(8, 48), S(12, 47), S(63, 40), S(86, 34), S(71, 39)},
	},
	{ // Queen
		{S(5, -108), S(2, -81), S(6, -88), S(10, -51), S(10, -87), S(-13, -79), S(3, -122), S(1, -112)},
		{S(-4, -81), S(5, -61), S(13, -84), S(3, -44), S(4, -54), S(18, -98), S(25, -109), S(17, -95)},
		{S(-9, -60), S(7, -24), S(0, 1), S(-6, -7), S(-4, -17), S(-1, 7), S(17, -36), S(7, -56)},
		{S(-12, -47), S(-5, 8), S(-10, 17), S(-20, 70), S(-21, 54), S(5, 41), S(8, 12), S(12, 10)},
		{S(-28, -19), S(-14, 21), S(-31, 24), S(-25, 77), S(-30, 99), S(-4, 94), S(20, 69), S(13, 52)},
		{S(-47, -5), S(-28, -3), S(-30, 22), S(-11, 38), S(6, 74), S(50, 66), S(57, 51), S(26, 78)},
		{S(-49, 14), S(-84, 42), S(-28, 21), S(-61, 70), S(-36, 103), S(22, 82), S(-33, 80), S(15, 27)},
		{S(-46, 18), S(-7, 8), S(-14, 44), S(11, 42), S(14, 56), S(59, 47), S(92, 19), S(64, 8)},
	},
	{ // King
		{S(142, -9), S(129, 27), S(53, 65), S(51, 71), S(79, 51), S(50, 68), S(150, 23), S(168, -29)},
		{S(156, 33), S(101, 53), S(41, 88), S(40, 110), S(30, 120), S(51, 87), S(126, 48), S(161, 32)},
		{S(113, 26), S(154, 43), S(94, 81), S(92, 118), S(76, 126), S(103, 78), S(143, 42), S(104, 35)},
		{S(142, 11), S(270, 33), S(162, 90), S(83, 138), S(89, 140), S(149, 87), S(228, 38), S(83, 32)},
		{S(130, 50), S(326, 47), S(156, 109), S(107, 139), S(86, 145), S(191, 97), S(269, 49), S(91, 49)},
		{S(138, 46), S(324, 68), S(226, 102), S(174, 114), S(177, 118), S(253, 92), S(301, 63), S(131, 44)},
		{S(136, -4), S(222, 61), S(212, 77), S(177, 89), S(127, 99), S(174, 76), S(203, 67), S(124, -3)},
		{S(180, -123), S(312, -43), S(201, -1), S(173, 52), S(154, 52), S(167, 11), S(287, -24), S(145, -100)},
	},
}

var PawnsConnected = [7][4]Score{
	{S(0, 0), S(0, 0), S(0, 0), S(0, 0)},
	{S(-1, -14), S(7, 4), S(3, -1), S(16, 17)},
	{S(11, 5), S(22, 3), S(24, 10), S(27, 16)},
	{S(11, 5), S(21, 8), S(10, 10), S(17, 21)},
	{S(5, 17), S(16, 25), S(29, 29), S(28, 23)},
	{S(37, 34), S(25, 73), S(73, 74), S(79, 94)},
	{S(176, 37), S(280, 25), S(275, 55), S(319, 57)},
}

var MobilityBonus = [...][32]Score{
	{S(-62, -127), S(-44, -74), S(-28, -26), S(-19, -3), S(-12, 12), S(-7, 27), // Knights
		S(1, 30), S(10, 25), S(21, 10)},
	{S(0, -133), S(7, -62), S(15, -20), S(23, 3), S(30, 18), S(35, 35), // Bishops
		S(37, 43), S(38, 45), S(37, 48), S(42, 47), S(46, 43), S(58, 32),
		S(76, 36), S(85, 10)},
	{S(-127, -156), S(-16, -38), S(-5, 15), S(-5, 43), S(-1, 56), S(1, 69), // Rooks
		S(3, 79), S(7, 85), S(12, 90), S(16, 96), S(19, 103), S(20, 108),
		S(26, 109), S(41, 99), S(99, 61)},
	{S(-413, -159), S(-130, -143), S(-48, -178), S(-26, -119), S(-11, -85), S(-11, -16), // Queens
		S(-6, 11), S(-3, 32), S(2, 50), S(5, 64), S(8, 72), S(12, 74),
		S(13, 78), S(15, 83), S(18, 77), S(15, 80), S(15, 76), S(13, 74),
		S(18, 63), S(25, 48), S(37, 28), S(40, 11), S(42, -8), S(56, -39),
		S(24, -29), S(-63, -12), S(140, -118), S(56, -79)},
}

var PassedFriendlyDistance = [8]Score{
	S(0, 0), S(-7, 39), S(-11, 23), S(-8, 9),
	S(-4, -5), S(1, -16), S(18, -27), S(5, -41),
}

var PassedEnemyDistance = [8]Score{
	S(0, 0), S(-101, -42), S(-5, -14), S(6, 5),
	S(13, 18), S(12, 27), S(7, 36), S(19, 42),
}

// PassedRank[Rank] contains a bonus according to the rank of a passed pawn, whether it can be pushed and whether the push would be safe
var PassedRank = [2][2][2][7]Score{
	{
		{
			{S(0, 0), S(-48, -21), S(-40, 9), S(-27, 13), S(29, 14), S(45, 19), S(176, 91)},
			{S(0, 0), S(-32, -50), S(-32, -18), S(-19, 6), S(30, 29), S(60, 40), S(171, 134)},
		},
		{
			{S(0, 0), S(-17, -27), S(-29, 0), S(-18, 19), S(14, 45), S(59, 80), S(255, 140)},
			{S(0, 0), S(-6, -50), S(-25, -6), S(-18, 24), S(14, 59), S(77, 85), S(231, 179)},
		},
	},
	{
		{
			{S(0, 0), S(-21, -9), S(-33, -3), S(-13, 22), S(29, 50), S(67, 98), S(278, 183)},
			{S(0, 0), S(-14, -24), S(-21, -7), S(-11, 35), S(37, 66), S(75, 108), S(239, 233)},
		},
		{
			{S(0, 0), S(-30, -15), S(-29, -11), S(-15, 32), S(19, 86), S(31, 229), S(70, 395)},
			{S(0, 0), S(-13, -35), S(-18, -17), S(-11, 32), S(23, 91), S(57, 216), S(137, 413)},
		},
	},
}

// PassedFile[File] contains a bonus according to the file of a passed pawn
var PassedFile = [8]Score{S(0, 23), S(-4, 27), S(-2, 11), S(-4, 0),
	S(-11, 4), S(-20, 14), S(-17, 25), S(11, 7),
}

var PassedStacked = [7]Score{S(0, 0), S(-17, -45), S(-25, -27), S(-34, -36), S(-16, -42), S(43, -99), S(0, 0)}
var PassedUncontested = [6]Score{S(0, 0), S(-93, 43), S(-85, 39), S(-96, 42), S(-91, 43), S(-74, 29)}
var PassedPushDefended = [6]Score{S(0, 0), S(-3, 21), S(-1, 8), S(6, 2), S(-2, 7), S(-10, 16)}
var PassedPushUncontestedDefended = [6]Score{S(0, 0), S(-57, 28), S(-36, 22), S(-66, 43), S(-73, 53), S(-57, 66)}

var Isolated = S(-7, -17)
var Doubled = S(-11, -31)
var Backward = S(8, -1)
var BackwardOpen = S(-4, -17)

var BishopPair = S(17, 74)
var BishopRammedPawns = S(-8, -23)

var BishopOutpostUndefendedBonus = S(23, -4)
var BishopOutpostDefendedBonus = S(56, 12)

var LongDiagonalBishop = S(23, 25)
var DistantBishop = [4]Score{S(-6, 1), S(-11, -1), S(-14, -4), S(-19, -21)}

var KnightOutpostUndefendedBonus = S(23, -22)
var KnightOutpostDefendedBonus = S(49, 13)

var DistantKnight = [4]Score{S(-16, -3), S(-17, -14), S(-29, -17), S(-55, -22)}

var MinorBehindPawn = S(8, 28)

var Tempo int16 = 15

// Rook on semiopen, open file
var RookOnFile = [2]Score{S(12, 12), S(41, 4)}
var RookOnQueenFile = S(6, 33)

var KingDefenders = [12]Score{
	S(-14, -10), S(-2, -11), S(1, -10), S(3, -7),
	S(6, -3), S(8, 6), S(12, 12), S(15, 11),
	S(14, 14), S(15, -62), S(-16, -12), S(11, 0),
}

var KingShelter = [2][8][8]Score{
	{{S(-25, 8), S(11, -12), S(19, 2), S(24, 6),
		S(30, -13), S(21, -6), S(17, -38), S(-95, 37)},
		{S(17, 6), S(25, -7), S(-6, 11), S(-8, 7),
			S(1, 0), S(24, 4), S(42, -13), S(-47, 22)},
		{S(16, 3), S(5, -2), S(-26, 0), S(-22, -1),
			S(0, -17), S(0, 5), S(20, -8), S(-31, -1)},
		{S(-19, 17), S(7, -1), S(-7, -9), S(-2, -8),
			S(6, -24), S(1, -9), S(18, 20), S(-16, -4)},
		{S(-32, 17), S(-15, 9), S(-32, 9), S(-20, 6),
			S(0, -8), S(-23, 0), S(7, -3), S(-27, 2)},
		{S(40, -18), S(23, -14), S(5, -12), S(4, -21),
			S(13, -27), S(26, -20), S(31, -19), S(-12, -4)},
		{S(17, -2), S(-7, -7), S(-29, 4), S(-18, 3),
			S(-8, -14), S(8, -3), S(15, -17), S(-24, 19)},
		{S(-40, 11), S(-26, -1), S(-19, 27), S(-20, 24),
			S(-22, 17), S(-10, 12), S(-11, -16), S(-71, 47)}},
	{{S(34, 56), S(-36, -11), S(-24, 4), S(-30, -4),
		S(-41, -15), S(-4, 6), S(-42, -16), S(-91, 36)},
		{S(138, 18), S(7, -4), S(-5, 6), S(-24, 19),
			S(-6, -11), S(14, 4), S(13, 0), S(-81, 24)},
		{S(-4, 32), S(35, 4), S(12, 0), S(20, -5),
			S(27, -4), S(11, 2), S(38, 5), S(-25, 4)},
		{S(6, 34), S(-13, 22), S(-15, 16), S(-23, 22),
			S(-1, 10), S(-5, 3), S(-5, 11), S(-42, -6)},
		{S(4, 23), S(-1, 9), S(-7, 5), S(-12, 2),
			S(-5, 1), S(0, -11), S(12, -2), S(-31, -6)},
		{S(17, 18), S(-6, 7), S(-13, 9), S(-8, 8),
			S(-3, 1), S(-39, 10), S(-10, 2), S(-39, 7)},
		{S(32, 4), S(-4, -15), S(-9, -11), S(-26, -6),
			S(-9, -19), S(0, -17), S(3, -20), S(-72, 19)},
		{S(-27, -3), S(-15, -30), S(-8, -17), S(-3, -22),
			S(-5, -34), S(-1, -22), S(-4, -50), S(-65, 24)}},
}

var KingStorm = [2][4][8]Score{
	{{S(12, 8), S(3, 15), S(9, 7), S(4, 9),
		S(3, 6), S(11, 2), S(8, 8), S(-11, -5)},
		{S(23, 1), S(16, 6), S(19, 3), S(5, 12),
			S(20, 2), S(27, -9), S(18, -6), S(-8, -5)},
		{S(23, 10), S(2, 7), S(6, 11), S(1, 12),
			S(4, 8), S(11, 1), S(7, -9), S(-2, 4)},
		{S(12, 2), S(8, -2), S(11, -10), S(-1, -13),
			S(-3, -7), S(15, -15), S(8, -12), S(-11, 2)}},
	{{S(0, 0), S(3, 4), S(-18, 10), S(10, -6),
		S(-2, 20), S(-13, 31), S(6, 50), S(-6, -4)},
		{S(0, 0), S(-53, -9), S(6, -3), S(47, -6),
			S(4, 0), S(-8, -5), S(15, 54), S(-8, 0)},
		{S(0, 0), S(-18, -11), S(-6, -5), S(17, -2),
			S(6, -6), S(-11, -18), S(39, -38), S(-7, 8)},
		{S(0, 0), S(-17, -24), S(-14, -24), S(-7, -11),
			S(1, -24), S(-1, -55), S(4, -2), S(-11, -5)}},
}
var KingOnPawnlessFlank = S(4, -67)

var Hanging = S(38, 16)
var ThreatByKing = S(-7, 35)
var ThreatByMinor = [King + 1]Score{S(0, 0), S(20, 38), S(17, 38), S(33, 28), S(30, 28), S(0, 27)}
var ThreatByRook = [King + 1]Score{S(0, 0), S(-4, 13), S(-2, 17), S(-5, -12), S(33, 9), S(13, 10)}

// This weights are from black piece on black square perspective
var RookBishopExistence = [16]Score{
	S(25, -29), S(5, -13), S(6, -13), S(3, 3), S(-1, -11), S(-3, 38), S(0, -3), S(-3, 21), S(2, -11), S(-1, -4), S(-3, 38), S(-4, 22), S(-9, -19), S(-4, 19), S(-2, 18), S(-18, 57),
}
var QueenBishopExistence = [16]Score{
	S(85, -73), S(-4, -9), S(-3, -23), S(-22, -5), S(-6, -11), S(10, 78), S(-18, 0), S(0, -2), S(-1, -8), S(-13, 26), S(12, 74), S(4, 9), S(-37, -16), S(-8, 37), S(-12, 28), S(-12, 36),
}
var KingBishopExistence = [16]Score{
	S(0, 0), S(1, 4), S(2, 4), S(-17, -11), S(-3, -10), S(-2, -7), S(-1, -3), S(0, 4), S(0, -1), S(1, 3), S(2, 7), S(3, 13), S(17, 11), S(0, -11), S(1, -8), S(0, 0),
}

// King safety
//

var KingSafetyAttacksWeights = [Queen + 1]Score{S(0, 0), S(-3, -3), S(-7, -7), S(-4, -4), S(4, 4)}
var KingSafetyWeakSquares = S(44, 44)
var KingSafetyFriendlyPawns = S(-35, -35)
var KingSafetyNoEnemyQueens = S(-176, -176)
var KingSafetySafeQueenCheck = S(90, 90)
var KingSafetySafeRookCheck = S(73, 73)
var KingSafetySafeBishopCheck = S(51, 51)
var KingSafetySafeKnightCheck = S(112, 112)
var KingSafetyAdjustment = S(-12, -12)

// Attack value is special as it is scaled by a fraction
var KingSafetyAttackValue = S(124, 124)
