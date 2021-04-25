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

var PawnValue = S(83, 150)
var KnightValue = S(421, 536)
var BishopValue = S(395, 527)
var RookValue = S(572, 862)
var QueenValue = S(1424, 1552)

// Pawns Square scores
// Bishop Flag, Rank, Row
var PawnScores = [16][7][8]Score{
	{
		{},
		{S(-24, 17), S(-17, 6), S(-24, 15), S(-9, 17), S(-23, 25), S(18, 14), S(41, 4), S(-6, 11)},
		{S(-22, 6), S(-22, 7), S(-21, 9), S(-12, 8), S(-3, 2), S(1, 7), S(12, 5), S(-12, 10)},
		{S(-21, 17), S(-26, 18), S(-2, 6), S(3, 0), S(16, -7), S(18, 1), S(8, 14), S(-9, 17)},
		{S(-16, 30), S(-9, 20), S(-8, 10), S(0, -1), S(31, -10), S(36, 2), S(40, 10), S(8, 22)},
		{S(-10, 64), S(10, 48), S(20, 38), S(36, 9), S(65, -4), S(75, 15), S(58, 40), S(3, 63)},
		{S(-8, 9), S(-6, 28), S(5, 29), S(27, 8), S(31, -4), S(31, 14), S(-89, 52), S(-128, 61)},
	},
	{
		{},
		{S(-16, 23), S(-8, 1), S(-21, 9), S(1, 15), S(-18, 16), S(24, 17), S(34, -7), S(-8, 10)},
		{S(-12, 6), S(-18, 0), S(-3, 6), S(-5, -1), S(10, 14), S(-9, 2), S(20, 10), S(-22, 10)},
		{S(-9, 17), S(-12, 11), S(6, -2), S(28, 1), S(25, -12), S(26, 6), S(7, 11), S(-10, 11)},
		{S(0, 30), S(14, -1), S(12, 15), S(19, -20), S(46, -8), S(42, -7), S(37, 12), S(-7, 26)},
		{S(-4, 60), S(23, 46), S(29, 12), S(51, 2), S(66, -27), S(86, 6), S(51, 10), S(5, 54)},
		{S(-4, 26), S(-9, 0), S(4, 6), S(23, -21), S(35, -15), S(28, -13), S(-90, 30), S(-130, 35)},
	},
	{
		{},
		{S(-16, 13), S(-8, 10), S(-9, 19), S(-15, 5), S(8, 30), S(14, 5), S(61, 1), S(-1, 17)},
		{S(-18, 6), S(-15, 6), S(-8, 0), S(4, 7), S(6, -8), S(8, 16), S(6, -3), S(-9, 17)},
		{S(-9, 2), S(-5, 10), S(14, 6), S(22, -7), S(44, -5), S(28, -3), S(17, 12), S(-4, 18)},
		{S(-4, 17), S(10, 11), S(6, 1), S(28, -1), S(46, -15), S(55, 4), S(34, 6), S(3, 30)},
		{S(-3, 59), S(16, 20), S(38, 24), S(39, -24), S(76, -9), S(84, 2), S(53, 24), S(11, 54)},
		{S(-11, -3), S(-6, 17), S(-1, -14), S(26, -9), S(32, -26), S(30, -9), S(-91, 26), S(-128, 46)},
	},
	{
		{},
		{S(-17, 7), S(-6, -9), S(-4, 4), S(-4, 9), S(-6, 19), S(13, 7), S(43, 5), S(-1, 0)},
		{S(-13, -1), S(-23, -7), S(-1, -3), S(-1, -9), S(16, 2), S(-7, 6), S(10, -1), S(-12, -2)},
		{S(-8, 2), S(-7, 4), S(10, -16), S(38, -6), S(34, -12), S(24, 3), S(11, 4), S(-12, 6)},
		{S(-6, 20), S(10, 2), S(8, 1), S(22, -18), S(42, -18), S(39, -8), S(32, 8), S(-6, 18)},
		{S(-2, 62), S(17, 34), S(31, 19), S(42, -9), S(67, -19), S(81, 4), S(54, 21), S(3, 51)},
		{S(-6, 14), S(-8, 10), S(3, 5), S(26, -14), S(33, -18), S(30, -9), S(-90, 30), S(-129, 44)},
	},
	{
		{},
		{S(-16, 14), S(-5, 13), S(-2, 8), S(0, 16), S(-16, 19), S(11, 15), S(41, 6), S(-3, 11)},
		{S(-11, 2), S(-18, 2), S(-5, -12), S(5, -4), S(-5, -7), S(-2, 12), S(1, 4), S(-12, 14)},
		{S(-10, 18), S(2, 9), S(16, -5), S(17, -23), S(36, -21), S(13, -4), S(10, 19), S(-10, 24)},
		{S(2, 41), S(17, 12), S(5, -4), S(29, -13), S(38, -33), S(50, 2), S(36, 9), S(-1, 42)},
		{S(-8, 75), S(18, 41), S(29, 21), S(47, -13), S(68, -20), S(84, -2), S(53, 24), S(6, 65)},
		{S(2, 35), S(-8, 10), S(5, 7), S(29, -10), S(33, -22), S(32, -9), S(-89, 32), S(-126, 49)},
	},
	{
		{},
		{S(-4, 20), S(1, 7), S(-8, 12), S(-5, 14), S(-11, 19), S(14, 24), S(26, -4), S(-1, 10)},
		{S(-7, 10), S(-13, -3), S(-1, 5), S(3, -8), S(5, 8), S(-2, 4), S(4, 13), S(-13, 9)},
		{S(-4, 23), S(4, 12), S(17, -6), S(29, -8), S(35, -13), S(23, 6), S(7, 17), S(-10, 18)},
		{S(12, 35), S(16, 6), S(7, 13), S(30, -18), S(49, -20), S(45, -6), S(41, 18), S(0, 29)},
		{S(-1, 71), S(21, 43), S(33, 19), S(47, -5), S(69, -30), S(82, 5), S(51, 19), S(6, 55)},
		{S(-2, 23), S(-11, 3), S(5, 8), S(28, -11), S(33, -18), S(30, -9), S(-89, 30), S(-130, 42)},
	},
	{
		{},
		{S(-8, -2), S(1, 1), S(-1, 6), S(-6, 5), S(1, 21), S(1, 1), S(50, 12), S(-2, 4)},
		{S(-9, -4), S(-10, -1), S(-13, -12), S(0, -3), S(-7, -14), S(-5, 18), S(-10, -9), S(-8, 9)},
		{S(-5, -1), S(1, -8), S(19, -7), S(14, -32), S(31, -8), S(7, -3), S(13, 14), S(-8, 14)},
		{S(0, 15), S(10, 8), S(7, -16), S(21, -6), S(37, -39), S(57, 2), S(27, 2), S(4, 26)},
		{S(-13, 59), S(18, 25), S(27, 19), S(32, -21), S(72, -18), S(82, 0), S(53, 21), S(10, 53)},
		{S(-7, 8), S(-8, 12), S(-4, -10), S(27, -10), S(32, -26), S(30, -9), S(-93, 20), S(-128, 45)},
	},
	{
		{},
		{S(-9, 4), S(-3, 5), S(-2, 4), S(-7, 9), S(2, 23), S(5, 13), S(37, 6), S(-3, 9)},
		{S(-10, -2), S(-18, -4), S(-6, 1), S(2, -3), S(4, 1), S(-5, 15), S(0, -1), S(-11, -3)},
		{S(-6, 6), S(-6, 9), S(19, -10), S(27, -10), S(34, -7), S(13, 7), S(9, 7), S(-13, 16)},
		{S(3, 26), S(15, 9), S(7, 7), S(27, -12), S(41, -20), S(48, -3), S(24, 7), S(-12, 22)},
		{S(-8, 60), S(16, 36), S(29, 20), S(45, -9), S(69, -20), S(85, 6), S(50, 19), S(8, 54)},
		{S(-6, 11), S(-8, 10), S(1, 3), S(26, -12), S(31, -20), S(30, -9), S(-90, 30), S(-129, 44)},
	},
	{
		{},
		{S(-15, 17), S(-9, 5), S(-15, 11), S(-8, 7), S(-12, 14), S(16, 8), S(28, 5), S(-3, 13)},
		{S(-14, 7), S(-18, -1), S(-10, -1), S(-5, -11), S(3, -5), S(-14, 4), S(16, 7), S(-11, 10)},
		{S(-9, 22), S(-7, 18), S(1, -11), S(19, -19), S(22, -27), S(15, 1), S(1, 15), S(-10, 20)},
		{S(0, 48), S(10, 8), S(-2, 5), S(8, -28), S(40, -9), S(26, -15), S(34, 16), S(-6, 33)},
		{S(9, 73), S(17, 43), S(41, 18), S(39, -12), S(68, -21), S(86, 12), S(63, 27), S(3, 55)},
		{S(-3, 28), S(-8, 11), S(7, 12), S(25, -22), S(35, -14), S(30, -12), S(-90, 30), S(-129, 45)},
	},
	{
		{},
		{S(-10, 13), S(-11, -4), S(-22, 8), S(4, 14), S(-20, 12), S(23, 16), S(21, -8), S(-4, 2)},
		{S(-13, 3), S(-19, -10), S(2, -1), S(-10, -10), S(10, 11), S(-18, -1), S(18, 7), S(-16, -7)},
		{S(-7, 20), S(-13, 6), S(3, -22), S(27, -5), S(19, -29), S(28, 5), S(1, 3), S(-14, 12)},
		{S(5, 31), S(11, -8), S(5, 4), S(3, -41), S(45, -14), S(25, -16), S(34, 11), S(-8, 21)},
		{S(6, 64), S(17, 37), S(30, 12), S(46, -1), S(66, -22), S(88, 5), S(51, 17), S(4, 52)},
		{S(-4, 18), S(-11, 2), S(5, 9), S(22, -23), S(33, -18), S(28, -15), S(-90, 30), S(-132, 38)},
	},
	{
		{},
		{S(-7, 16), S(-2, 9), S(-9, 22), S(-11, 12), S(-12, 27), S(14, 8), S(38, 10), S(2, 6)},
		{S(-11, 12), S(-7, 6), S(-8, 8), S(2, 5), S(3, -11), S(-2, 13), S(16, 2), S(-7, 6)},
		{S(-2, 16), S(3, 16), S(16, 2), S(20, -4), S(32, -5), S(25, -1), S(20, 10), S(-7, 14)},
		{S(0, 35), S(23, 17), S(9, 6), S(22, -7), S(44, -23), S(56, -2), S(38, 4), S(2, 29)},
		{S(5, 73), S(20, 38), S(35, 24), S(34, -24), S(76, -14), S(85, 8), S(63, 23), S(7, 51)},
		{S(-7, 11), S(-5, 16), S(-1, -4), S(26, -9), S(35, -19), S(30, -12), S(-91, 27), S(-130, 44)},
	},
	{
		{},
		{S(-12, 14), S(-4, 1), S(-10, 10), S(1, 13), S(-15, 21), S(17, 14), S(31, -1), S(-5, 1)},
		{S(-11, -1), S(-14, 1), S(-4, 7), S(-1, 1), S(6, 6), S(-9, 6), S(15, 3), S(-16, 3)},
		{S(-2, 15), S(1, 8), S(12, -7), S(28, 0), S(28, -10), S(29, 4), S(9, 6), S(-20, 11)},
		{S(8, 23), S(20, 14), S(13, 9), S(23, -28), S(39, -19), S(47, -5), S(26, 9), S(-16, 26)},
		{S(0, 62), S(21, 38), S(31, 17), S(40, -12), S(67, -20), S(81, 4), S(55, 24), S(2, 52)},
		{S(-8, 12), S(-8, 10), S(3, 5), S(26, -12), S(33, -18), S(29, -10), S(-90, 30), S(-129, 43)},
	},
	{
		{},
		{S(-18, 20), S(-7, 10), S(-9, 12), S(-2, 17), S(-28, 13), S(4, 12), S(19, 4), S(-4, 2)},
		{S(-11, 15), S(-21, 3), S(-16, 3), S(-3, -8), S(-3, -3), S(-7, 1), S(2, 4), S(-11, 1)},
		{S(-6, 29), S(-5, 11), S(10, -23), S(12, -16), S(23, -37), S(20, -5), S(2, 9), S(-6, 8)},
		{S(7, 30), S(14, 3), S(5, 6), S(20, -32), S(38, -16), S(30, -18), S(34, 14), S(0, 27)},
		{S(5, 71), S(21, 42), S(30, 19), S(46, -5), S(60, -23), S(84, 4), S(53, 21), S(6, 58)},
		{S(-5, 18), S(-9, 8), S(3, 7), S(22, -20), S(32, -18), S(30, -11), S(-90, 29), S(-129, 41)},
	},
	{
		{},
		{S(-10, 18), S(-3, 7), S(-14, 10), S(-5, 9), S(-20, 16), S(10, 28), S(11, 0), S(-4, 4)},
		{S(-10, 11), S(-24, -1), S(-1, 5), S(-2, -15), S(3, 8), S(-5, 0), S(3, 6), S(-16, -2)},
		{S(-2, 20), S(-5, 16), S(11, -13), S(27, -8), S(26, -25), S(23, 4), S(2, 5), S(-12, 10)},
		{S(14, 35), S(12, 7), S(10, 9), S(22, -33), S(46, -16), S(27, -15), S(40, 13), S(-7, 28)},
		{S(10, 65), S(22, 40), S(31, 21), S(44, -7), S(55, -19), S(82, 4), S(52, 21), S(7, 54)},
		{S(-4, 16), S(-8, 9), S(3, 5), S(26, -14), S(33, -18), S(28, -11), S(-90, 30), S(-129, 44)},
	},
	{
		{},
		{S(-15, 11), S(-1, 10), S(-6, 20), S(-5, 11), S(-19, 19), S(-2, 6), S(28, 16), S(0, 8)},
		{S(-14, 0), S(-9, 9), S(-11, -6), S(1, 1), S(-11, -10), S(-7, 18), S(3, -6), S(-8, 9)},
		{S(-8, 13), S(1, 14), S(16, 1), S(15, -21), S(35, -16), S(13, -7), S(13, 7), S(-11, 9)},
		{S(1, 33), S(31, 16), S(1, 5), S(28, -13), S(33, -34), S(54, -9), S(30, 4), S(-3, 32)},
		{S(9, 74), S(17, 34), S(36, 24), S(32, -22), S(68, -17), S(79, 1), S(49, 22), S(4, 51)},
		{S(-9, 11), S(-9, 11), S(0, 2), S(25, -12), S(32, -20), S(31, -8), S(-89, 29), S(-130, 44)},
	},
	{
		{},
		{S(-9, 8), S(1, 11), S(-12, 13), S(-8, 9), S(-17, 9), S(3, 17), S(17, 4), S(-5, 4)},
		{S(-8, 3), S(-14, 1), S(-5, 0), S(0, -2), S(0, -1), S(-11, 20), S(-2, 1), S(-16, 1)},
		{S(-1, 13), S(2, 12), S(16, -3), S(24, -9), S(30, -18), S(17, 1), S(4, 11), S(-19, 13)},
		{S(10, 34), S(27, 12), S(6, 15), S(27, -25), S(38, -22), S(45, -9), S(25, 8), S(-12, 24)},
		{S(8, 66), S(21, 39), S(30, 22), S(44, -10), S(67, -20), S(76, 2), S(47, 20), S(1, 53)},
		{S(-11, 14), S(-11, 8), S(1, 5), S(25, -13), S(33, -18), S(30, -9), S(-89, 30), S(-129, 44)},
	},
}

// Piece Square Values
var PieceScores = [King + 1][8][4]Score{
	{},
	{ // knight
		{S(-93, -40), S(-25, -44), S(-30, -34), S(-22, -11)},
		{S(-23, -23), S(-24, -20), S(-17, -36), S(-19, -16)},
		{S(-17, -49), S(-4, -26), S(-12, -17), S(1, 2)},
		{S(-10, -6), S(15, -5), S(5, 17), S(4, 24)},
		{S(7, -5), S(4, -1), S(22, 20), S(14, 35)},
		{S(-26, -14), S(-17, 1), S(7, 26), S(18, 24)},
		{S(6, -32), S(-12, -13), S(32, -30), S(41, 4)},
		{S(-231, -36), S(-92, -11), S(-150, 24), S(-11, -13)},
	},
	{ // Bishop
		{S(24, -37), S(29, -3), S(9, -4), S(6, 6)},
		{S(30, -49), S(21, -38), S(31, -12), S(14, 0)},
		{S(21, -10), S(40, -3), S(10, -13), S(28, 14)},
		{S(13, -12), S(21, 3), S(23, 16), S(35, 19)},
		{S(-6, 8), S(23, 11), S(14, 19), S(37, 26)},
		{S(-3, 9), S(-5, 24), S(7, 3), S(23, 21)},
		{S(-38, 27), S(-20, 7), S(14, 20), S(-1, 28)},
		{S(-44, 4), S(-39, 29), S(-113, 42), S(-97, 47)},
	},
	{ // Rook
		{S(-17, -26), S(-9, -17), S(-2, -14), S(3, -24)},
		{S(-58, -11), S(-13, -31), S(-5, -29), S(-15, -26)},
		{S(-33, -19), S(-8, -11), S(-18, -11), S(-15, -15)},
		{S(-26, 1), S(-15, 12), S(-13, 13), S(-10, 5)},
		{S(-12, 14), S(13, 12), S(24, 13), S(33, 9)},
		{S(-13, 24), S(48, 7), S(45, 19), S(57, 5)},
		{S(0, 30), S(6, 34), S(40, 25), S(43, 34)},
		{S(43, 32), S(65, 31), S(40, 38), S(37, 33)},
	},
	{ // Queen
		{S(10, -105), S(10, -90), S(11, -95), S(17, -61)},
		{S(7, -86), S(17, -79), S(20, -90), S(10, -50)},
		{S(3, -65), S(16, -31), S(7, 2), S(1, -11)},
		{S(1, -23), S(9, 5), S(2, 31), S(-13, 60)},
		{S(14, -12), S(7, 38), S(-8, 51), S(-20, 84)},
		{S(7, 4), S(16, 0), S(-7, 55), S(4, 50)},
		{S(-27, 10), S(-73, 54), S(-13, 51), S(-37, 84)},
		{S(-6, 10), S(28, 4), S(25, 35), S(33, 38)},
	},
	{ // King
		{S(167, -25), S(149, 19), S(55, 60), S(62, 80)},
		{S(165, 33), S(131, 40), S(56, 76), S(24, 136)},
		{S(107, 32), S(154, 30), S(113, 63), S(80, 136)},
		{S(87, 24), S(237, 25), S(156, 70), S(96, 145)},
		{S(89, 44), S(270, 37), S(181, 82), S(108, 144)},
		{S(117, 43), S(298, 58), S(248, 80), S(184, 121)},
		{S(115, 0), S(211, 66), S(189, 73), S(149, 113)},
		{S(157, -98), S(290, -17), S(173, 18), S(161, 75)},
	},
}

var PawnsConnected = [7][4]Score{
	{S(0, 0), S(0, 0), S(0, 0), S(0, 0)},
	{S(-2, -13), S(8, 4), S(3, 0), S(17, 18)},
	{S(11, 6), S(22, 4), S(24, 11), S(27, 16)},
	{S(10, 5), S(21, 7), S(11, 9), S(16, 19)},
	{S(5, 17), S(15, 25), S(31, 27), S(28, 20)},
	{S(38, 27), S(30, 60), S(76, 61), S(95, 72)},
	{S(176, 30), S(300, 24), S(294, 40), S(346, 44)},
}

var MobilityBonus = [...][32]Score{
	{S(-60, -134), S(-43, -82), S(-27, -36), S(-19, -13), S(-12, 2), S(-8, 17), // Knights
		S(1, 20), S(9, 15), S(20, 0)},
	{S(-1, -138), S(6, -65), S(14, -26), S(22, -2), S(29, 14), S(33, 30), // Bishops
		S(36, 39), S(36, 42), S(35, 46), S(39, 46), S(42, 43), S(53, 35),
		S(83, 35), S(92, 14)},
	{S(-127, -146), S(-15, -37), S(-3, 14), S(-4, 42), S(0, 54), S(2, 65), // Rooks
		S(4, 75), S(8, 78), S(13, 83), S(17, 87), S(20, 91), S(22, 94),
		S(28, 93), S(43, 82), S(98, 46)},
	{S(-413, -159), S(-122, -138), S(-31, -174), S(-13, -119), S(-1, -84), S(-1, -14), // Queens
		S(4, 11), S(7, 29), S(11, 44), S(14, 56), S(17, 63), S(21, 63),
		S(23, 64), S(23, 68), S(26, 61), S(23, 64), S(22, 59), S(20, 57),
		S(24, 44), S(31, 29), S(41, 10), S(38, -2), S(39, -21), S(49, -50),
		S(14, -38), S(-88, -10), S(127, -127), S(49, -81)},
}

var PassedFriendlyDistance = [8]Score{
	S(0, 0), S(-11, 41), S(-13, 24), S(-13, 8),
	S(-6, -6), S(-3, -15), S(16, -27), S(-1, -38),
}

var PassedEnemyDistance = [8]Score{
	S(0, 0), S(-129, -57), S(-15, -8), S(2, 11),
	S(14, 22), S(14, 29), S(11, 37), S(19, 41),
}
var PawnPsqt [16][2][64]Score   // BishopFlag, colour, Square
var Psqt [2][King + 1][64]Score // One row for every colour purposefelly left empty

var PawnsConnectedSquare [2][64]Score
var pawnsConnectedMask [2][64]uint64

// PassedRank[Rank] contains a bonus according to the rank of a passed pawn
var PassedRank = [7]Score{S(0, 0), S(-13, -31), S(-20, -13), S(-9, 25), S(29, 67), S(52, 145), S(179, 237)}

// PassedFile[File] contains a bonus according to the file of a passed pawn
var PassedFile = [8]Score{S(-11, 34), S(-12, 34), S(-7, 12), S(-4, -4),
	S(-10, 1), S(-17, 10), S(-16, 23), S(12, 10),
}

var PassedStacked = [8]Score{S(0, 0), S(-14, -53), S(-20, -34), S(-31, -56), S(-3, -85), S(28, -207), S(0, 0), S(0, 0)}

var Isolated = S(-7, -16)
var Doubled = S(-9, -33)
var Backward = S(7, -1)
var BackwardOpen = S(-4, -15)

var BishopPair = S(20, 75)
var BishopRammedPawns = S(-8, -23)

var BishopOutpostUndefendedBonus = S(20, -2)
var BishopOutpostDefendedBonus = S(51, 13)

var LongDiagonalBishop = S(23, 24)

var KnightOutpostUndefendedBonus = S(20, -19)
var KnightOutpostDefendedBonus = S(45, 13)

var DistantKnight = [4]Score{S(-12, -2), S(-18, -11), S(-27, -16), S(-52, -20)}

var MinorBehindPawn = S(8, 32)

var Tempo int16 = 15

// Rook on semiopen, open file
var RookOnFile = [2]Score{S(9, 18), S(36, 11)}
var RookOnQueenFile = S(6, 30)

var KingDefenders = [12]Score{
	S(-5, -9), S(2, -9), S(5, -8), S(7, -6),
	S(8, -3), S(10, 5), S(11, 11), S(13, 10),
	S(11, 8), S(15, -70), S(-3, -10), S(11, 0),
}

var KingShelter = [2][8][8]Score{
	{{S(-27, 4), S(8, -17), S(16, -5), S(19, -1),
		S(23, -15), S(17, -10), S(12, -40), S(-95, 34)},
		{S(9, -2), S(24, -16), S(-2, -2), S(-10, -4),
			S(-2, -10), S(24, -10), S(40, -25), S(-51, 17)},
		{S(12, 11), S(-1, 8), S(-32, 11), S(-28, 9),
			S(-9, -7), S(-7, 12), S(14, -1), S(-37, 10)},
		{S(-28, 29), S(10, 7), S(-5, -2), S(0, -3),
			S(7, -20), S(5, -5), S(27, 8), S(-15, 4)},
		{S(-28, 19), S(-15, 11), S(-33, 11), S(-22, 9),
			S(-1, -5), S(-28, 5), S(3, -7), S(-28, 6)},
		{S(39, -13), S(21, -9), S(4, -7), S(3, -17),
			S(12, -24), S(25, -15), S(31, -18), S(-13, 1)},
		{S(17, -11), S(-8, -16), S(-28, -6), S(-17, -6),
			S(-9, -24), S(9, -12), S(13, -24), S(-23, 12)},
		{S(-38, 7), S(-26, -4), S(-18, 22), S(-20, 20),
			S(-23, 12), S(-11, 8), S(-12, -17), S(-75, 46)}},
	{{S(48, 42), S(-45, -10), S(-28, 2), S(-39, -1),
		S(-51, -12), S(-22, 17), S(-67, 5), S(-92, 38)},
		{S(149, 9), S(-1, -9), S(-6, 0), S(-28, 13),
			S(-15, -14), S(15, -7), S(7, -7), S(-86, 29)},
		{S(5, 22), S(42, -13), S(14, -14), S(23, -20),
			S(30, -19), S(16, -12), S(46, -12), S(-27, 1)},
		{S(6, 49), S(-11, 34), S(-17, 28), S(-24, 32),
			S(-2, 18), S(-17, 17), S(-12, 21), S(-43, 8)},
		{S(-10, 44), S(-3, 25), S(-6, 18), S(-12, 15),
			S(-4, 14), S(4, -8), S(15, -6), S(-31, 10)},
		{S(34, 3), S(-3, -5), S(-8, -2), S(-3, -4),
			S(2, -11), S(-44, -2), S(-8, -13), S(-40, 2)},
		{S(43, -1), S(-1, -17), S(-5, -15), S(-23, -9),
			S(-6, -23), S(3, -20), S(3, -22), S(-70, 18)},
		{S(-5, -3), S(-15, -31), S(-8, -19), S(-4, -23),
			S(-8, -34), S(-1, -25), S(-6, -50), S(-69, 29)}},
}

var KingStorm = [2][4][8]Score{
	{{S(12, 4), S(3, 11), S(9, 4), S(4, 6),
		S(2, 3), S(10, 0), S(7, 6), S(-6, -10)},
		{S(21, -7), S(15, -3), S(17, -4), S(3, 4),
			S(17, -6), S(24, -15), S(16, -14), S(-9, -13)},
		{S(24, 11), S(4, 7), S(6, 12), S(2, 14),
			S(5, 9), S(12, 1), S(9, -12), S(0, 2)},
		{S(8, 18), S(6, 14), S(10, 8), S(0, 3),
			S(-2, 6), S(15, 3), S(9, -2), S(-10, 9)}},
	{{S(0, 0), S(19, -5), S(-18, 7), S(10, -10),
		S(-2, 16), S(-15, 31), S(4, 51), S(-2, -8)},
		{S(0, 0), S(-41, -20), S(2, -10), S(44, -12),
			S(1, -8), S(-11, -10), S(13, 56), S(-9, -8)},
		{S(0, 0), S(-17, -9), S(-7, -3), S(17, -1),
			S(7, -7), S(-7, -18), S(47, -47), S(-6, 7)},
		{S(0, 0), S(-13, -13), S(-17, -10), S(-8, 2),
			S(2, -15), S(1, -49), S(6, 3), S(-12, 12)}},
}
var KingOnPawnlessFlank = S(19, -89)

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

var KingSafetyAttacksWeights = [King + 1]int16{0, -3, -7, -4, 7, 0}
var KingSafetyAttackValue int16 = 122
var KingSafetyWeakSquares int16 = 43
var KingSafetyFriendlyPawns int16 = -33
var KingSafetyNoEnemyQueens int16 = -160
var KingSafetySafeQueenCheck int16 = 92
var KingSafetySafeRookCheck int16 = 76
var KingSafetySafeBishopCheck int16 = 60
var KingSafetySafeKnightCheck int16 = 113
var KingSafetyAdjustment int16 = -14

var Hanging = S(34, 13)
var ThreatByKing = S(-9, 33)
var ThreatByMinor = [King + 1]Score{S(0, 0), S(20, 39), S(17, 34), S(33, 26), S(30, 32), S(10, 23)}
var ThreatByRook = [King + 1]Score{S(0, 0), S(-3, 11), S(-2, 16), S(-5, -9), S(34, 7), S(26, -1)}

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

			Psqt[White][King][y*8+x] = PieceScores[King][y][x]
			Psqt[White][King][y*8+(7-x)] = PieceScores[King][y][x]
			Psqt[Black][King][(7-y)*8+x] = PieceScores[King][y][x]
			Psqt[Black][King][(7-y)*8+(7-x)] = PieceScores[King][y][x]

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
	score += Psqt[White][King][whiteKingLocation]
	score += KingDefenders[whiteKingDefenders]
	if tuning {
		T.PieceScores[King][Rank(whiteKingLocation)][FileMirror[File(whiteKingLocation)]]++
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
	score -= Psqt[Black][King][blackKingLocation]
	score -= KingDefenders[blackKingDefenders]
	if tuning {
		T.PieceScores[King][7-Rank(blackKingLocation)][FileMirror[File(blackKingLocation)]]--
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
