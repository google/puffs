// Copyright 2017 The Wuffs Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package token

// MaxIntBits is the largest size (in bits) of the i8, u8, i16, u16, etc.
// integer types.
const MaxIntBits = 64

// ID is a token type. Every identifier (in the programming language sense),
// keyword, operator and literal has its own ID.
//
// Some IDs are built-in: the "func" keyword always has the same numerical ID
// value. Others are mapped at runtime. For example, the ID value for the
// "foobar" identifier (e.g. a variable name) is looked up in a Map.
type ID uint32

// Str returns a string form of x.
func (x ID) Str(m *Map) string { return m.ByID(x) }

func (x ID) AmbiguousForm() ID {
	if x >= ID(len(ambiguousForms)) {
		return 0
	}
	return ambiguousForms[x]
}

func (x ID) UnaryForm() ID {
	if x >= ID(len(unaryForms)) {
		return 0
	}
	return unaryForms[x]
}

func (x ID) BinaryForm() ID {
	if x >= ID(len(binaryForms)) {
		return 0
	}
	return binaryForms[x]
}

func (x ID) AssociativeForm() ID {
	if x >= ID(len(associativeForms)) {
		return 0
	}
	return associativeForms[x]
}

func (x ID) IsBuiltIn() bool { return x < nBuiltInIDs }

func (x ID) IsUnaryOp() bool       { return minOp <= x && x <= maxOp && unaryForms[x] != 0 }
func (x ID) IsBinaryOp() bool      { return minOp <= x && x <= maxOp && binaryForms[x] != 0 }
func (x ID) IsAssociativeOp() bool { return minOp <= x && x <= maxOp && associativeForms[x] != 0 }

func (x ID) IsLiteral(m *Map) bool {
	if x < nBuiltInIDs {
		return minBuiltInLiteral <= x && x <= maxBuiltInLiteral
	} else if s := m.ByID(x); s != "" {
		return !alpha(s[0])
	}
	return false
}

func (x ID) IsNumLiteral(m *Map) bool {
	if x < nBuiltInIDs {
		return minBuiltInNumLiteral <= x && x <= maxBuiltInNumLiteral
	} else if s := m.ByID(x); s != "" {
		return numeric(s[0])
	}
	return false
}

func (x ID) IsStrLiteral(m *Map) bool {
	if x < nBuiltInIDs {
		return false
	} else if s := m.ByID(x); s != "" {
		return s[0] == '"'
	}
	return false
}

func (x ID) IsIdent(m *Map) bool {
	if x < nBuiltInIDs {
		return minBuiltInIdent <= x && x <= maxBuiltInIdent
	} else if s := m.ByID(x); s != "" {
		return alpha(s[0])
	}
	return false
}

func (x ID) IsOpen() bool       { return x < ID(len(isOpen)) && isOpen[x] }
func (x ID) IsClose() bool      { return x < ID(len(isClose)) && isClose[x] }
func (x ID) IsTightLeft() bool  { return x < ID(len(isTightLeft)) && isTightLeft[x] }
func (x ID) IsTightRight() bool { return x < ID(len(isTightRight)) && isTightRight[x] }

func (x ID) IsAssign() bool         { return minAssign <= x && x <= maxAssign }
func (x ID) IsNumType() bool        { return minNumType <= x && x <= maxNumType }
func (x ID) IsNumTypeOrIdeal() bool { return minNumTypeOrIdeal <= x && x <= maxNumTypeOrIdeal }

func (x ID) IsImplicitSemicolon(m *Map) bool {
	return x.IsLiteral(m) || x.IsIdent(m) ||
		(x < ID(len(isImplicitSemicolon)) && isImplicitSemicolon[x])
}

func (x ID) IsXOp() bool            { return minXOp <= x && x <= maxXOp }
func (x ID) IsXUnaryOp() bool       { return minXOp <= x && x <= maxXOp && unaryForms[x] != 0 }
func (x ID) IsXBinaryOp() bool      { return minXOp <= x && x <= maxXOp && binaryForms[x] != 0 }
func (x ID) IsXAssociativeOp() bool { return minXOp <= x && x <= maxXOp && associativeForms[x] != 0 }

func (x ID) SmallPowerOf2Value() int {
	switch x {
	case ID1:
		return 1
	case ID2:
		return 2
	case ID4:
		return 4
	case ID8:
		return 8
	case ID16:
		return 16
	case ID32:
		return 32
	case ID64:
		return 64
	case ID128:
		return 128
	case ID256:
		return 256
	}
	return 0
}

// QID is a qualified ID, such as "foo.bar". QID[0] is "foo"'s ID and QID[1] is
// "bar"'s. QID[0] may be 0 for a plain "bar".
type QID [2]ID

func (x QID) IsZero() bool { return x == QID{} }

// Str returns a string form of x.
func (x QID) Str(m *Map) string {
	if x[0] != 0 {
		return m.ByID(x[0]) + "." + m.ByID(x[1])
	}
	return m.ByID(x[1])
}

// QQID is a double-qualified ID, such as "receiverPkg.receiverName.funcName".
type QQID [3]ID

func (x QQID) IsZero() bool { return x == QQID{} }

// Str returns a string form of x.
func (x QQID) Str(m *Map) string {
	if x[0] != 0 {
		return m.ByID(x[0]) + "." + m.ByID(x[1]) + "." + m.ByID(x[2])
	}
	if x[1] != 0 {
		return m.ByID(x[1]) + "." + m.ByID(x[2])
	}
	return m.ByID(x[2])
}

// Token combines an ID and the line number it was seen.
type Token struct {
	ID   ID
	Line uint32
}

// nBuiltInIDs is the number of built-in IDs. The packing is:
//  - Zero is invalid.
//  - [ 0x01,  0x0F] are squiggly punctuation, such as "(", ")" and ";".
//  - [ 0x10,  0x1F] are squiggly assignments, such as "=" and "+=".
//  - [ 0x20,  0x3F] are operators, such as "+", "==" and "not".
//  - [ 0x40,  0x6F] are x-ops (disambiguation forms): unary vs binary "+".
//  - [ 0x70,  0x8F] are keywords, such as "if" and "return".
//  - [ 0x90,  0x93] are type modifiers, such as "ptr" and "slice".
//  - [ 0x94,  0x9F] are literals, such as "false" and "true".
//  - [ 0xA0,  0xFF] are reserved.
//  - [0x100, 0x3FF] are identifiers, such as "bool", "u32" and "read_u8".
//
// "Squiggly" means a sequence of non-alpha-numeric characters, such as "+" and
// "&=". Roughly speaking, their IDs range in [0x01, 0x3F], or disambiguation
// forms range in [0x40, 0x6F], but vice versa does not necessarily hold. For
// example, the "and" operator is not "squiggly" but it is within [0x01, 0x3F].
const (
	nBuiltInSymbolicIDs = ID(0x70)  // 112
	nBuiltInIDs         = ID(0x400) // 1024
)

const (
	IDInvalid = ID(0)

	IDOpenParen    = ID(0x02)
	IDCloseParen   = ID(0x03)
	IDOpenBracket  = ID(0x04)
	IDCloseBracket = ID(0x05)
	IDOpenCurly    = ID(0x06)
	IDCloseCurly   = ID(0x07)

	IDDot       = ID(0x08)
	IDDotDot    = ID(0x09)
	IDComma     = ID(0x0A)
	IDExclam    = ID(0x0B)
	IDQuestion  = ID(0x0C)
	IDColon     = ID(0x0D)
	IDSemicolon = ID(0x0E)
	IDDollar    = ID(0x0F)
)

const (
	minAssign = 0x10
	maxAssign = 0x1F

	IDEq              = ID(0x10)
	IDPlusEq          = ID(0x11)
	IDMinusEq         = ID(0x12)
	IDStarEq          = ID(0x13)
	IDSlashEq         = ID(0x14)
	IDShiftLEq        = ID(0x15)
	IDShiftREq        = ID(0x16)
	IDAmpEq           = ID(0x17)
	IDPipeEq          = ID(0x18)
	IDHatEq           = ID(0x19)
	IDPercentEq       = ID(0x1A)
	IDTildeModPlusEq  = ID(0x1B)
	IDTildeModMinusEq = ID(0x1C)
	IDTildeSatPlusEq  = ID(0x1D)
	IDTildeSatMinusEq = ID(0x1E)
	IDEqColon         = ID(0x1F)
)

const (
	minOp          = 0x20
	minAmbiguousOp = 0x20
	maxAmbiguousOp = 0x3F
	minXOp         = 0x40
	maxXOp         = 0x6F
	maxOp          = 0x6F

	IDPlus          = ID(0x21)
	IDMinus         = ID(0x22)
	IDStar          = ID(0x23)
	IDSlash         = ID(0x24)
	IDShiftL        = ID(0x25)
	IDShiftR        = ID(0x26)
	IDAmp           = ID(0x27)
	IDPipe          = ID(0x28)
	IDHat           = ID(0x29)
	IDPercent       = ID(0x2A)
	IDTildeModPlus  = ID(0x2B)
	IDTildeModMinus = ID(0x2C)
	IDTildeSatPlus  = ID(0x2D)
	IDTildeSatMinus = ID(0x2E)

	IDNotEq       = ID(0x30)
	IDLessThan    = ID(0x31)
	IDLessEq      = ID(0x32)
	IDEqEq        = ID(0x33)
	IDGreaterEq   = ID(0x34)
	IDGreaterThan = ID(0x35)

	IDAnd   = ID(0x38)
	IDOr    = ID(0x39)
	IDNot   = ID(0x3A)
	IDAs    = ID(0x3B)
	IDRef   = ID(0x3C)
	IDDeref = ID(0x3D)

	// The IDXFoo IDs are not returned by the tokenizer. They are used by the
	// ast.Node ID-typed fields to disambiguate e.g. unary vs binary plus.

	IDXUnaryPlus  = ID(0x40)
	IDXUnaryMinus = ID(0x41)
	IDXUnaryNot   = ID(0x42)
	IDXUnaryRef   = ID(0x43)
	IDXUnaryDeref = ID(0x44)

	IDXBinaryPlus          = ID(0x48)
	IDXBinaryMinus         = ID(0x49)
	IDXBinaryStar          = ID(0x4A)
	IDXBinarySlash         = ID(0x4B)
	IDXBinaryShiftL        = ID(0x4C)
	IDXBinaryShiftR        = ID(0x4D)
	IDXBinaryAmp           = ID(0x4E)
	IDXBinaryPipe          = ID(0x4F)
	IDXBinaryHat           = ID(0x50)
	IDXBinaryPercent       = ID(0x51)
	IDXBinaryTildeModPlus  = ID(0x52)
	IDXBinaryTildeModMinus = ID(0x53)
	IDXBinaryTildeSatPlus  = ID(0x54)
	IDXBinaryTildeSatMinus = ID(0x55)
	IDXBinaryNotEq         = ID(0x56)
	IDXBinaryLessThan      = ID(0x57)
	IDXBinaryLessEq        = ID(0x58)
	IDXBinaryEqEq          = ID(0x59)
	IDXBinaryGreaterEq     = ID(0x5A)
	IDXBinaryGreaterThan   = ID(0x5B)
	IDXBinaryAnd           = ID(0x5C)
	IDXBinaryOr            = ID(0x5D)
	IDXBinaryAs            = ID(0x5E)

	IDXAssociativePlus = ID(0x60)
	IDXAssociativeStar = ID(0x61)
	IDXAssociativeAmp  = ID(0x62)
	IDXAssociativePipe = ID(0x63)
	IDXAssociativeHat  = ID(0x64)
	IDXAssociativeAnd  = ID(0x65)
	IDXAssociativeOr   = ID(0x66)
)

const (
	minKeyword = 0x70
	maxKeyword = 0x8F

	// TODO: sort these by name, when the list has stabilized.
	IDFunc       = ID(0x70)
	IDAssert     = ID(0x71)
	IDWhile      = ID(0x72)
	IDIf         = ID(0x73)
	IDElse       = ID(0x74)
	IDReturn     = ID(0x75)
	IDBreak      = ID(0x76)
	IDContinue   = ID(0x77)
	IDStruct     = ID(0x78)
	IDUse        = ID(0x79)
	IDVar        = ID(0x7A)
	IDPre        = ID(0x7B)
	IDInv        = ID(0x7C)
	IDPost       = ID(0x7D)
	IDVia        = ID(0x7E)
	IDPub        = ID(0x7F)
	IDPri        = ID(0x80)
	IDError      = ID(0x81)
	IDSuspension = ID(0x82)
	IDPackageID  = ID(0x83)
	IDConst      = ID(0x84)
	IDTry        = ID(0x85)
	IDIterate    = ID(0x86)
	IDYield      = ID(0x87)
	IDIOBind     = ID(0x88)
)

const (
	minTypeModifier = 0x90
	maxTypeModifier = 0x93

	IDArray = ID(0x90)
	IDNptr  = ID(0x91)
	IDPtr   = ID(0x92)
	IDSlice = ID(0x93)
)

const (
	minBuiltInLiteral    = 0x94
	minBuiltInNumLiteral = 0x96
	maxBuiltInNumLiteral = 0x9F
	maxBuiltInLiteral    = 0x9F

	IDFalse = ID(0x94)
	IDTrue  = ID(0x95)
	ID0     = ID(0x96)
	ID1     = ID(0x97)
	ID2     = ID(0x98)
	ID4     = ID(0x99)
	ID8     = ID(0x9A)
	ID16    = ID(0x9B)
	ID32    = ID(0x9C)
	ID64    = ID(0x9D)
	ID128   = ID(0x9E)
	ID256   = ID(0x9F)
)

const (
	minBuiltInIdent   = 0x100
	minNumTypeOrIdeal = 0x10F
	minNumType        = 0x110
	maxNumType        = 0x117
	maxNumTypeOrIdeal = 0x117
	maxBuiltInIdent   = 0x3FF

	// -------- 0x100 block.

	IDEmptyStruct = ID(0x100)
	IDBool        = ID(0x101)
	IDStatus      = ID(0x102)
	IDIOReader    = ID(0x103)
	IDIOWriter    = ID(0x104)
	IDCapitalT    = ID(0x105)
	IDImageConfig = ID(0x106)

	IDDiamond = ID(0x10E)

	// It is important that IDDoubleZ is right next to the IDI8..IDU64 block.
	// See the ID.IsNumTypeOrIdeal method.
	IDDoubleZ = ID(0x10F)

	IDI8  = ID(0x110)
	IDI16 = ID(0x111)
	IDI32 = ID(0x112)
	IDI64 = ID(0x113)
	IDU8  = ID(0x114)
	IDU16 = ID(0x115)
	IDU32 = ID(0x116)
	IDU64 = ID(0x117)

	// TODO Read/Write 24 bits? It might be useful for RGB triples.

	IDUnreadU8  = ID(0x120)
	IDReadU8    = ID(0x121)
	IDReadU16BE = ID(0x122)
	IDReadU16LE = ID(0x123)
	IDReadU32BE = ID(0x124)
	IDReadU32LE = ID(0x125)
	IDReadU64BE = ID(0x126)
	IDReadU64LE = ID(0x127)

	// TODO: IDUnwriteU8?
	IDWriteU8    = ID(0x131)
	IDWriteU16BE = ID(0x132)
	IDWriteU16LE = ID(0x133)
	IDWriteU32BE = ID(0x134)
	IDWriteU32LE = ID(0x135)
	IDWriteU64BE = ID(0x136)
	IDWriteU64LE = ID(0x137)

	IDUnderscore = ID(0x140)
	IDThis       = ID(0x141)
	IDIn         = ID(0x142)
	IDOut        = ID(0x143)
	IDBase       = ID(0x144)
	IDReset      = ID(0x145)

	// -------- 0x200 block.

	IDIsError           = ID(0x200)
	IDIsOK              = ID(0x201)
	IDIsSuspension      = ID(0x202)
	IDCopyFromHistory32 = ID(0x203)
	IDCopyFromReader32  = ID(0x204)
	IDCopyFromSlice     = ID(0x205)
	IDCopyFromSlice32   = ID(0x206)
	IDSkip32            = ID(0x207)
	IDSkip64            = ID(0x208)
	IDLength            = ID(0x209)
	IDAvailable         = ID(0x20A)
	IDPrefix            = ID(0x20B)
	IDSuffix            = ID(0x20C)
	IDSetLimit          = ID(0x20D)
	IDLowBits           = ID(0x20E)
	IDHighBits          = ID(0x20F)
	IDUnroll            = ID(0x210)
	IDMark              = ID(0x211)
	IDSinceMark         = ID(0x212)
)

var builtInsByID = [nBuiltInIDs]string{
	IDOpenParen:    "(",
	IDCloseParen:   ")",
	IDOpenBracket:  "[",
	IDCloseBracket: "]",
	IDOpenCurly:    "{",
	IDCloseCurly:   "}",

	IDDot:       ".",
	IDDotDot:    "..",
	IDComma:     ",",
	IDExclam:    "!",
	IDQuestion:  "?",
	IDColon:     ":",
	IDSemicolon: ";",
	IDDollar:    "$",

	IDEq:              "=",
	IDPlusEq:          "+=",
	IDMinusEq:         "-=",
	IDStarEq:          "*=",
	IDSlashEq:         "/=",
	IDShiftLEq:        "<<=",
	IDShiftREq:        ">>=",
	IDAmpEq:           "&=",
	IDPipeEq:          "|=",
	IDHatEq:           "^=",
	IDPercentEq:       "%=",
	IDTildeModPlusEq:  "~mod+=",
	IDTildeModMinusEq: "~mod-=",
	IDTildeSatPlusEq:  "~sat+=",
	IDTildeSatMinusEq: "~sat-=",
	IDEqColon:         "=:",

	IDPlus:          "+",
	IDMinus:         "-",
	IDStar:          "*",
	IDSlash:         "/",
	IDShiftL:        "<<",
	IDShiftR:        ">>",
	IDAmp:           "&",
	IDPipe:          "|",
	IDHat:           "^",
	IDPercent:       "%",
	IDTildeModPlus:  "~mod+",
	IDTildeModMinus: "~mod-",
	IDTildeSatPlus:  "~sat+",
	IDTildeSatMinus: "~sat-",

	IDNotEq:       "!=",
	IDLessThan:    "<",
	IDLessEq:      "<=",
	IDEqEq:        "==",
	IDGreaterEq:   ">=",
	IDGreaterThan: ">",

	IDAnd:   "and",
	IDOr:    "or",
	IDNot:   "not",
	IDAs:    "as",
	IDRef:   "ref",
	IDDeref: "deref",

	IDFunc:       "func",
	IDAssert:     "assert",
	IDWhile:      "while",
	IDIf:         "if",
	IDElse:       "else",
	IDReturn:     "return",
	IDBreak:      "break",
	IDContinue:   "continue",
	IDStruct:     "struct",
	IDUse:        "use",
	IDVar:        "var",
	IDPre:        "pre",
	IDInv:        "inv",
	IDPost:       "post",
	IDVia:        "via",
	IDPub:        "pub",
	IDPri:        "pri",
	IDError:      "error",
	IDSuspension: "suspension",
	IDPackageID:  "packageid",
	IDConst:      "const",
	IDTry:        "try",
	IDIterate:    "iterate",
	IDYield:      "yield",
	IDIOBind:     "io_bind",

	IDArray: "array",
	IDNptr:  "nptr",
	IDPtr:   "ptr",
	IDSlice: "slice",

	IDFalse: "false",
	IDTrue:  "true",
	ID0:     "0",
	ID1:     "1",
	ID2:     "2",
	ID4:     "4",
	ID8:     "8",
	ID16:    "16",
	ID32:    "32",
	ID64:    "64",
	ID128:   "128",
	ID256:   "256",

	// -------- 0x100 block.

	IDEmptyStruct: "empty_struct",
	IDBool:        "bool",
	IDStatus:      "status",
	IDIOReader:    "io_reader",
	IDIOWriter:    "io_writer",
	IDCapitalT:    "T",
	IDImageConfig: "image_config",

	// IDDiamond and IDDoubleZ are never returned by the tokenizer, as the
	// tokenizer rejects non-ASCII input.
	//
	// The string representations ""◊ and "ℤ" are specifically non-ASCII so
	// that no user-defined (non built-in) identifier will conflict with them.

	// IDDiamond is used by the type checker as a dummy-valued built-in ID to
	// represent a generic type.
	IDDiamond: "◊", // U+25C7 WHITE DIAMOND

	// IDDoubleZ is used by the type checker as a dummy-valued built-in ID to
	// represent an ideal integer type (in mathematical terms, the integer ring
	// ℤ), as opposed to a realized integer type whose range is restricted. For
	// example, the base.u16 type is restricted to [0x0000, 0xFFFF].
	IDDoubleZ: "ℤ", // U+2124 DOUBLE-STRUCK CAPITAL Z

	// Change MaxIntBits if a future update adds an i128 or u128 type.
	IDI8:  "i8",
	IDI16: "i16",
	IDI32: "i32",
	IDI64: "i64",
	IDU8:  "u8",
	IDU16: "u16",
	IDU32: "u32",
	IDU64: "u64",

	IDUnreadU8:  "unread_u8",
	IDReadU8:    "read_u8",
	IDReadU16BE: "read_u16be",
	IDReadU16LE: "read_u16le",
	IDReadU32BE: "read_u32be",
	IDReadU32LE: "read_u32le",
	IDReadU64BE: "read_u64be",
	IDReadU64LE: "read_u64le",

	IDWriteU8:    "write_u8",
	IDWriteU16BE: "write_u16be",
	IDWriteU16LE: "write_u16le",
	IDWriteU32BE: "write_u32be",
	IDWriteU32LE: "write_u32le",
	IDWriteU64BE: "write_u64be",
	IDWriteU64LE: "write_u64le",

	IDUnderscore: "_",
	IDThis:       "this",
	IDIn:         "in",
	IDOut:        "out",
	IDBase:       "base",
	IDReset:      "reset",

	// -------- 0x200 block.

	IDIsError:           "is_error",
	IDIsOK:              "is_ok",
	IDIsSuspension:      "is_suspension",
	IDCopyFromHistory32: "copy_from_history32",
	IDCopyFromReader32:  "copy_from_reader32",
	IDCopyFromSlice:     "copy_from_slice",
	IDCopyFromSlice32:   "copy_from_slice32",
	IDSkip32:            "skip32",
	IDSkip64:            "skip64",
	IDLength:            "length",
	IDAvailable:         "available",
	IDPrefix:            "prefix",
	IDSuffix:            "suffix",
	IDSetLimit:          "set_limit",
	IDLowBits:           "low_bits",
	IDHighBits:          "high_bits",
	IDUnroll:            "unroll",
	IDMark:              "mark",
	IDSinceMark:         "since_mark",
}

var builtInsByName = map[string]ID{}

func init() {
	for i, name := range builtInsByID {
		if name != "" {
			builtInsByName[name] = ID(i)
		}
	}
}

// squiggles are built-in IDs that aren't alpha-numeric.
var squiggles = [256]ID{
	'(': IDOpenParen,
	')': IDCloseParen,
	'[': IDOpenBracket,
	']': IDCloseBracket,
	'{': IDOpenCurly,
	'}': IDCloseCurly,

	',': IDComma,
	'?': IDQuestion,
	':': IDColon,
	';': IDSemicolon,
	'$': IDDollar,
}

type suffixLexer struct {
	suffix string
	id     ID
}

// lexers lex ambiguous 1-byte squiggles. For example, "&" might be the start
// of "&^" or "&=".
//
// The order of the []suffixLexer elements matters. The first match wins. Since
// we want to lex greedily, longer suffixes should be earlier in the slice.
var lexers = [256][]suffixLexer{
	'.': {
		{".", IDDotDot},
		{"", IDDot},
	},
	'!': {
		{"=", IDNotEq},
		{"", IDExclam},
	},
	'&': {
		{"=", IDAmpEq},
		{"", IDAmp},
	},
	'|': {
		{"=", IDPipeEq},
		{"", IDPipe},
	},
	'^': {
		{"=", IDHatEq},
		{"", IDHat},
	},
	'+': {
		{"=", IDPlusEq},
		{"", IDPlus},
	},
	'-': {
		{"=", IDMinusEq},
		{"", IDMinus},
	},
	'*': {
		{"=", IDStarEq},
		{"", IDStar},
	},
	'/': {
		{"=", IDSlashEq},
		{"", IDSlash},
	},
	'%': {
		{"=", IDPercentEq},
		{"", IDPercent},
	},
	'=': {
		{"=", IDEqEq},
		{":", IDEqColon},
		{"", IDEq},
	},
	'<': {
		{"<=", IDShiftLEq},
		{"<", IDShiftL},
		{"=", IDLessEq},
		{"", IDLessThan},
	},
	'>': {
		{">=", IDShiftREq},
		{">", IDShiftR},
		{"=", IDGreaterEq},
		{"", IDGreaterThan},
	},
	'~': {
		{"mod+=", IDTildeModPlusEq},
		{"mod+", IDTildeModPlus},
		{"mod-=", IDTildeModMinusEq},
		{"mod-", IDTildeModMinus},
		{"sat+=", IDTildeSatPlusEq},
		{"sat+", IDTildeSatPlus},
		{"sat-=", IDTildeSatMinusEq},
		{"sat-", IDTildeSatMinus},
	},
}

var ambiguousForms = [nBuiltInSymbolicIDs]ID{
	IDXUnaryPlus:  IDPlus,
	IDXUnaryMinus: IDMinus,
	IDXUnaryNot:   IDNot,
	IDXUnaryRef:   IDRef,
	IDXUnaryDeref: IDDeref,

	IDXBinaryPlus:          IDPlus,
	IDXBinaryMinus:         IDMinus,
	IDXBinaryStar:          IDStar,
	IDXBinarySlash:         IDSlash,
	IDXBinaryShiftL:        IDShiftL,
	IDXBinaryShiftR:        IDShiftR,
	IDXBinaryAmp:           IDAmp,
	IDXBinaryPipe:          IDPipe,
	IDXBinaryHat:           IDHat,
	IDXBinaryPercent:       IDPercent,
	IDXBinaryTildeModPlus:  IDTildeModPlus,
	IDXBinaryTildeModMinus: IDTildeModMinus,
	IDXBinaryTildeSatPlus:  IDTildeSatPlus,
	IDXBinaryTildeSatMinus: IDTildeSatMinus,
	IDXBinaryNotEq:         IDNotEq,
	IDXBinaryLessThan:      IDLessThan,
	IDXBinaryLessEq:        IDLessEq,
	IDXBinaryEqEq:          IDEqEq,
	IDXBinaryGreaterEq:     IDGreaterEq,
	IDXBinaryGreaterThan:   IDGreaterThan,
	IDXBinaryAnd:           IDAnd,
	IDXBinaryOr:            IDOr,
	IDXBinaryAs:            IDAs,

	IDXAssociativePlus: IDPlus,
	IDXAssociativeStar: IDStar,
	IDXAssociativeAmp:  IDAmp,
	IDXAssociativePipe: IDPipe,
	IDXAssociativeHat:  IDHat,
	IDXAssociativeAnd:  IDAnd,
	IDXAssociativeOr:   IDOr,
}

func init() {
	addXForms(&unaryForms)
	addXForms(&binaryForms)
	addXForms(&associativeForms)
}

// addXForms modifies table so that, if table[x] == y, then table[y] = y.
//
// For example, for the unaryForms table, the explicit entries are like:
//  IDPlus:        IDXUnaryPlus,
// and this function implicitly addes entries like:
//  IDXUnaryPlus:  IDXUnaryPlus,
func addXForms(table *[nBuiltInSymbolicIDs]ID) {
	implicitEntries := [nBuiltInSymbolicIDs]bool{}
	for _, y := range table {
		if y != 0 {
			implicitEntries[y] = true
		}
	}
	for y, implicit := range implicitEntries {
		if implicit {
			table[y] = ID(y)
		}
	}
}

var unaryForms = [nBuiltInSymbolicIDs]ID{
	IDPlus:  IDXUnaryPlus,
	IDMinus: IDXUnaryMinus,
	IDNot:   IDXUnaryNot,
	IDRef:   IDXUnaryRef,
	IDDeref: IDXUnaryDeref,
}

var binaryForms = [nBuiltInSymbolicIDs]ID{
	IDPlusEq:          IDXBinaryPlus,
	IDMinusEq:         IDXBinaryMinus,
	IDStarEq:          IDXBinaryStar,
	IDSlashEq:         IDXBinarySlash,
	IDShiftLEq:        IDXBinaryShiftL,
	IDShiftREq:        IDXBinaryShiftR,
	IDAmpEq:           IDXBinaryAmp,
	IDPipeEq:          IDXBinaryPipe,
	IDHatEq:           IDXBinaryHat,
	IDPercentEq:       IDXBinaryPercent,
	IDTildeModPlusEq:  IDXBinaryTildeModPlus,
	IDTildeModMinusEq: IDXBinaryTildeModMinus,
	IDTildeSatPlusEq:  IDXBinaryTildeSatPlus,
	IDTildeSatMinusEq: IDXBinaryTildeSatMinus,

	IDPlus:          IDXBinaryPlus,
	IDMinus:         IDXBinaryMinus,
	IDStar:          IDXBinaryStar,
	IDSlash:         IDXBinarySlash,
	IDShiftL:        IDXBinaryShiftL,
	IDShiftR:        IDXBinaryShiftR,
	IDAmp:           IDXBinaryAmp,
	IDPipe:          IDXBinaryPipe,
	IDHat:           IDXBinaryHat,
	IDPercent:       IDXBinaryPercent,
	IDTildeModPlus:  IDXBinaryTildeModPlus,
	IDTildeModMinus: IDXBinaryTildeModMinus,
	IDTildeSatPlus:  IDXBinaryTildeSatPlus,
	IDTildeSatMinus: IDXBinaryTildeSatMinus,

	IDNotEq:       IDXBinaryNotEq,
	IDLessThan:    IDXBinaryLessThan,
	IDLessEq:      IDXBinaryLessEq,
	IDEqEq:        IDXBinaryEqEq,
	IDGreaterEq:   IDXBinaryGreaterEq,
	IDGreaterThan: IDXBinaryGreaterThan,
	IDAnd:         IDXBinaryAnd,
	IDOr:          IDXBinaryOr,
	IDAs:          IDXBinaryAs,
}

var associativeForms = [nBuiltInSymbolicIDs]ID{
	IDPlus: IDXAssociativePlus,
	IDStar: IDXAssociativeStar,
	IDAmp:  IDXAssociativeAmp,
	IDPipe: IDXAssociativePipe,
	IDHat:  IDXAssociativeHat,
	// TODO: IDTildeModPlus, IDTildeSatPlus?
	IDAnd: IDXAssociativeAnd,
	IDOr:  IDXAssociativeOr,
}

var isOpen = [...]bool{
	IDOpenParen:   true,
	IDOpenBracket: true,
	IDOpenCurly:   true,
}

var isClose = [...]bool{
	IDCloseParen:   true,
	IDCloseBracket: true,
	IDCloseCurly:   true,
}

var isTightLeft = [...]bool{
	IDCloseParen:   true,
	IDOpenBracket:  true,
	IDCloseBracket: true,

	IDDot:       true,
	IDDotDot:    true,
	IDComma:     true,
	IDExclam:    true,
	IDQuestion:  true,
	IDColon:     true,
	IDSemicolon: true,
}

var isTightRight = [...]bool{
	IDOpenParen:   true,
	IDOpenBracket: true,

	IDDot:      true,
	IDDotDot:   true,
	IDExclam:   true,
	IDQuestion: true,
	IDColon:    true,
	IDDollar:   true,
}

var isImplicitSemicolon = [...]bool{
	IDCloseParen:   true,
	IDCloseBracket: true,
	IDCloseCurly:   true,

	IDReturn:   true,
	IDBreak:    true,
	IDContinue: true,
}
