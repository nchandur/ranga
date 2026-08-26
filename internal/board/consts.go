package board

// file exclusion masks used to prevent wrapped across board edges
const (
	NotAFile  BitBoard = 0xfefefefefefefefe // masks out file A
	NotHFile  BitBoard = 0x7f7f7f7f7f7f7f7f // masks out file H
	NotABFile BitBoard = 0xfcfcfcfcfcfcfcfc // masks out files A and B
	NotGHFile BitBoard = 0x3f3f3f3f3f3f3f3f // masks out files G and H

)

// precomputed non-sliding piece attack tables
var (
	PawnAttacks   = [2][64]BitBoard{} // lookup pawn attack masks indexed by [color][square]
	KnightAttacks = [64]BitBoard{}    // lookup knight attack masks indexed by [square]
	KingAttacks   = [64]BitBoard{}    // lookup king attack masks indexed by [square]
)

// precomputed sliding piece occupancy masks
var (
	BishopMasks = [64]BitBoard{} // stores relevant occupancy masks for bishops per square
	RookMasks   = [64]BitBoard{} // stores relevant occupancy masks for rooks per square

)

// magic bitboard attack lookup tables indexed by [square][magic_index]
var (
	BishopAttacks = [64][512]BitBoard{}  // attack tables for bishops
	RookAttacks   = [64][4096]BitBoard{} // attack tables for rooks
)

// defines number of set bits in relevant occupancy mask for bishop on each square
var BishopRelevantOccupancy = [64]int{
	6, 5, 5, 5, 5, 5, 5, 6,
	5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	6, 5, 5, 5, 5, 5, 5, 6,
}

// defines number of set bits in relevant occupancy mask for rook on each square
var RookRelevantOccupancy = [64]int{
	12, 11, 11, 11, 11, 11, 11, 12,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	12, 11, 11, 11, 11, 11, 11, 12,
}

// castle rights for each square
var CastleRights = []Castle{
	7, 15, 15, 15, 3, 15, 15, 11,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	13, 15, 15, 15, 12, 15, 15, 14,
}

// magic numbers
var RookMagicNumbers = [64]BitBoard{
	0x80001080204004,
	0x2840100020004001,
	0x8100084011002001,
	0x200082042001004,
	0x880080080040002,
	0x7000a1801004400,
	0x2080420001000080,
	0x180002480104100,
	0x1044800080400020,
	0x50400050002006,
	0x101004020001101,
	0x802000810204600,
	0x4101001008010004,
	0x128808004000200,
	0x2011000200040100,
	0x80200110200805c,
	0x8000c000406000,
	0x1090004020004000,
	0x8210008010802009,
	0x520420008120020,
	0x50008010010,
	0x4008002000480,
	0x200040030410288,
	0x105020020804401,
	0x4100802080004000,
	0x2010004040002000,
	0x20010100201040,
	0x10008080080012,
	0x88b040080080080,
	0x1000300080400,
	0x4000021400300851,
	0x84630282000acb0c,
	0x80004000c02001,
	0x2080210081004000,
	0x4cd0802000801002,
	0x4441080084801000,
	0x40082800800,
	0x2018802000410,
	0x890002c804000130,
	0x6000169106000064,
	0x4400480248000,
	0x11000402010c004,
	0x609002000410010,
	0x1082042020a0010,
	0x1020040008008080,
	0x10020004008080,
	0x4000420801040010,
	0x400110050820004,
	0xa090150a0800300,
	0x1a200080400580,
	0xa021801000a00180,
	0xa0100008008480,
	0x4901001004080100,
	0x100800200040080,
	0x810d000a00040500,
	0x480084400910a00,
	0x4404110026044082,
	0x128225004001,
	0x884411a008804202,
	0x880100021000409,
	0x301002800500205,
	0x802000c08011016,
	0x1004020090010804,
	0x1080010020804c02,
}

var BishopMagicNumbers = [64]BitBoard{
	0x820040400541020,
	0x101000c0808000,
	0x2021120408400528,
	0x49040700024018,
	0x12220210d8102092,
	0x2221040840010,
	0x2002441c04401404,
	0x200108801109000,
	0x90e04494084040,
	0x820218120084,
	0x1040800930460,
	0x4240580c000,
	0x40308000008,
	0x4020803080000,
	0x7000004202202146,
	0x9020212602104408,
	0x8064108080080,
	0x8002010810210205,
	0x800400280044040a,
	0x900802004004,
	0x4002110401040064,
	0x2001012210120102,
	0x101030084012000,
	0x2090300042021000,
	0x44040440090804,
	0x1401084020820c28,
	0x1000280014081021,
	0x802080054004088,
	0xc0848024002000,
	0x280810020806000,
	0x8020860004124608,
	0x508005040100,
	0x4104040840408,
	0x8822100080808,
	0x120804104100400,
	0x122110800040040,
	0x10010040300404,
	0x1802108200010040,
	0x1c2220040920840,
	0x2040110434148,
	0x80040c2242100800,
	0x101140144002000,
	0x802020202008110,
	0x1050020d1002800,
	0x300242200a000102,
	0x10151008a10100,
	0x110100600b004c8,
	0x10241c6182021020,
	0x8820443008880001,
	0x2004404040820,
	0x100011088040502,
	0x42000210440000,
	0x400810240001,
	0xc40d022089020,
	0x2060880a08114a08,
	0x8c5100082008000,
	0x250082080a0a1212,
	0x1008010046126049,
	0x9000400e01108802,
	0x1100081405421203,
	0xc220002008030400,
	0x8000104002042448,
	0x51a8081004614a,
	0x144504883040080,
}

// NNUE hidden layer size
const NNUEHiddenSize int = 256

// print pieces
const PieceChar = "PNBRQKpnbrqk "

// piece color lookup
var PieceColor = [13]Color{
	WP:    White,
	WN:    White,
	WB:    White,
	WR:    White,
	WQ:    White,
	WK:    White,
	BP:    Black,
	BN:    Black,
	BB:    Black,
	BR:    Black,
	BQ:    Black,
	BK:    Black,
	Empty: Both,
}

// static piece value
var PieceValue = [13]int{
	WP: 100,
	WN: 325,
	WB: 325,
	WR: 500,
	WQ: 800,
	WK: 10000,
	BP: -100,
	BN: -325,
	BB: -325,
	BR: -500,
	BQ: -800,
	BK: -10000,
}
