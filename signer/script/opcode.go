package script

// An Opcode defines the information related to a txscript opcode.  opfunc, if
// present, is the function to call to perform the opcode on the script.  The
// current script is passed in as a slice with the first member being the opcode
// itself.
type Opcode struct {
	Value  byte
	Length int
	//opfunc func(*parsedOpcode, *Engine) error
}

// These constants are the values of the official opcodes used on the btc wiki,
// in bitcoin core and in most if not all other references and software related
// to handling BTC scripts.
const (
	Op0                   = 0x00 // 0
	OpFalse               = 0x00 // 0 - AKA Op0
	OpData1               = 0x01 // 1
	OpData2               = 0x02 // 2
	OpData3               = 0x03 // 3
	OpData4               = 0x04 // 4
	OpData5               = 0x05 // 5
	OpData6               = 0x06 // 6
	OpData7               = 0x07 // 7
	OpData8               = 0x08 // 8
	OpData9               = 0x09 // 9
	OpData10              = 0x0a // 10
	OpData11              = 0x0b // 11
	OpData12              = 0x0c // 12
	OpData13              = 0x0d // 13
	OpData14              = 0x0e // 14
	OpData15              = 0x0f // 15
	OpData16              = 0x10 // 16
	OpData17              = 0x11 // 17
	OpData18              = 0x12 // 18
	OpData19              = 0x13 // 19
	OpData20              = 0x14 // 20
	OpData21              = 0x15 // 21
	OpData22              = 0x16 // 22
	OpData23              = 0x17 // 23
	OpData24              = 0x18 // 24
	OpData25              = 0x19 // 25
	OpData26              = 0x1a // 26
	OpData27              = 0x1b // 27
	OpData28              = 0x1c // 28
	OpData29              = 0x1d // 29
	OpData30              = 0x1e // 30
	OpData31              = 0x1f // 31
	OpData32              = 0x20 // 32
	OpData33              = 0x21 // 33
	OpData34              = 0x22 // 34
	OpData35              = 0x23 // 35
	OpData36              = 0x24 // 36
	OpData37              = 0x25 // 37
	OpData38              = 0x26 // 38
	OpData39              = 0x27 // 39
	OpData40              = 0x28 // 40
	OpData41              = 0x29 // 41
	OpData42              = 0x2a // 42
	OpData43              = 0x2b // 43
	OpData44              = 0x2c // 44
	OpData45              = 0x2d // 45
	OpData46              = 0x2e // 46
	OpData47              = 0x2f // 47
	OpData48              = 0x30 // 48
	OpData49              = 0x31 // 49
	OpData50              = 0x32 // 50
	OpData51              = 0x33 // 51
	OpData52              = 0x34 // 52
	OpData53              = 0x35 // 53
	OpData54              = 0x36 // 54
	OpData55              = 0x37 // 55
	OpData56              = 0x38 // 56
	OpData57              = 0x39 // 57
	OpData58              = 0x3a // 58
	OpData59              = 0x3b // 59
	OpData60              = 0x3c // 60
	OpData61              = 0x3d // 61
	OpData62              = 0x3e // 62
	OpData63              = 0x3f // 63
	OpData64              = 0x40 // 64
	OpData65              = 0x41 // 65
	OpData66              = 0x42 // 66
	OpData67              = 0x43 // 67
	OpData68              = 0x44 // 68
	OpData69              = 0x45 // 69
	OpData70              = 0x46 // 70
	OpData71              = 0x47 // 71
	OpData72              = 0x48 // 72
	OpData73              = 0x49 // 73
	OpData74              = 0x4a // 74
	OpData75              = 0x4b // 75
	OpPushData1           = 0x4c // 76
	OpPushData2           = 0x4d // 77
	OpPushData4           = 0x4e // 78
	Op1NEGATE             = 0x4f // 79
	OpReserved            = 0x50 // 80
	Op1                   = 0x51 // 81 - AKA OpTRUE
	OpTRUE                = 0x51 // 81
	Op2                   = 0x52 // 82
	Op3                   = 0x53 // 83
	Op4                   = 0x54 // 84
	Op5                   = 0x55 // 85
	Op6                   = 0x56 // 86
	Op7                   = 0x57 // 87
	Op8                   = 0x58 // 88
	Op9                   = 0x59 // 89
	Op10                  = 0x5a // 90
	Op11                  = 0x5b // 91
	Op12                  = 0x5c // 92
	Op13                  = 0x5d // 93
	Op14                  = 0x5e // 94
	Op15                  = 0x5f // 95
	Op16                  = 0x60 // 96
	OpNOP                 = 0x61 // 97
	OpVER                 = 0x62 // 98
	OpIF                  = 0x63 // 99
	OpNOTIF               = 0x64 // 100
	OpVERIF               = 0x65 // 101
	OpVERNOTIF            = 0x66 // 102
	OpELSE                = 0x67 // 103
	OpENDIF               = 0x68 // 104
	OpVERIFY              = 0x69 // 105
	OpRETURN              = 0x6a // 106
	OpTOALTSTACK          = 0x6b // 107
	OpFROMALTSTACK        = 0x6c // 108
	Op2DROP               = 0x6d // 109
	Op2DUP                = 0x6e // 110
	Op3DUP                = 0x6f // 111
	Op2OVER               = 0x70 // 112
	Op2ROT                = 0x71 // 113
	Op2SWAP               = 0x72 // 114
	OpIFDUP               = 0x73 // 115
	OpDEPTH               = 0x74 // 116
	OpDROP                = 0x75 // 117
	OpDUP                 = 0x76 // 118
	OpNIP                 = 0x77 // 119
	OpOVER                = 0x78 // 120
	OpPICK                = 0x79 // 121
	OpROLL                = 0x7a // 122
	OpROT                 = 0x7b // 123
	OpSWAP                = 0x7c // 124
	OpTUCK                = 0x7d // 125
	OpCAT                 = 0x7e // 126
	OpSUBSTR              = 0x7f // 127
	OpLEFT                = 0x80 // 128
	OpRIGHT               = 0x81 // 129
	OpSIZE                = 0x82 // 130
	OpINVERT              = 0x83 // 131
	OpAND                 = 0x84 // 132
	OpOR                  = 0x85 // 133
	OpXOR                 = 0x86 // 134
	OpEQUAL               = 0x87 // 135
	OpEQUALVERIFY         = 0x88 // 136
	OpReserved1           = 0x89 // 137
	OpReserved2           = 0x8a // 138
	Op1ADD                = 0x8b // 139
	Op1SUB                = 0x8c // 140
	Op2MUL                = 0x8d // 141
	Op2DIV                = 0x8e // 142
	OpNEGATE              = 0x8f // 143
	OpABS                 = 0x90 // 144
	OpNOT                 = 0x91 // 145
	Op0NOTEQUAL           = 0x92 // 146
	OpADD                 = 0x93 // 147
	OpSUB                 = 0x94 // 148
	OpMUL                 = 0x95 // 149
	OpDIV                 = 0x96 // 150
	OpMOD                 = 0x97 // 151
	OpLSHIFT              = 0x98 // 152
	OpRSHIFT              = 0x99 // 153
	OpBOOLAND             = 0x9a // 154
	OpBOOLOR              = 0x9b // 155
	OpNUMEQUAL            = 0x9c // 156
	OpNUMEQUALVERIFY      = 0x9d // 157
	OpNUMNOTEQUAL         = 0x9e // 158
	OpLESSTHAN            = 0x9f // 159
	OpGREATERTHAN         = 0xa0 // 160
	OpLESSTHANOREQUAL     = 0xa1 // 161
	OpGREATERTHANOREQUAL  = 0xa2 // 162
	OpMIN                 = 0xa3 // 163
	OpMAX                 = 0xa4 // 164
	OpWITHIN              = 0xa5 // 165
	OpRIPEMD160           = 0xa6 // 166
	OpSHA1                = 0xa7 // 167
	OpSHA256              = 0xa8 // 168
	OpHASH160             = 0xa9 // 169
	OpHASH256             = 0xaa // 170
	OpCODESEPARATOR       = 0xab // 171
	OpCHECKSIG            = 0xac // 172
	OpCHECKSIGVERIFY      = 0xad // 173
	OpCheckMultiSig       = 0xae // 174
	OpCheckMultiSigVerify = 0xaf // 175
	OpNOP1                = 0xb0 // 176
	OpNOP2                = 0xb1 // 177
	OpCHECKLOCKTIMEVERIFY = 0xb1 // 177 - AKA OpNOP2
	OpNOP3                = 0xb2 // 178
	OpCHECKSEQUENCEVERIFY = 0xb2 // 178 - AKA OpNOP3
	OpNOP4                = 0xb3 // 179
	OpNOP5                = 0xb4 // 180
	OpNOP6                = 0xb5 // 181
	OpNOP7                = 0xb6 // 182
	OpNOP8                = 0xb7 // 183
	OpNOP9                = 0xb8 // 184
	OpNOP10               = 0xb9 // 185
	OpUnknown186          = 0xba // 186
	OpUnknown187          = 0xbb // 187
	OpUnknown188          = 0xbc // 188
	OpUnknown189          = 0xbd // 189
	OpUnknown190          = 0xbe // 190
	OpUnknown191          = 0xbf // 191
	OpUnknown192          = 0xc0 // 192
	OpUnknown193          = 0xc1 // 193
	OpUnknown194          = 0xc2 // 194
	OpUnknown195          = 0xc3 // 195
	OpUnknown196          = 0xc4 // 196
	OpUnknown197          = 0xc5 // 197
	OpUnknown198          = 0xc6 // 198
	OpUnknown199          = 0xc7 // 199
	OpUnknown200          = 0xc8 // 200
	OpUnknown201          = 0xc9 // 201
	OpUnknown202          = 0xca // 202
	OpUnknown203          = 0xcb // 203
	OpUnknown204          = 0xcc // 204
	OpUnknown205          = 0xcd // 205
	OpUnknown206          = 0xce // 206
	OpUnknown207          = 0xcf // 207
	OpUnknown208          = 0xd0 // 208
	OpUnknown209          = 0xd1 // 209
	OpUnknown210          = 0xd2 // 210
	OpUnknown211          = 0xd3 // 211
	OpUnknown212          = 0xd4 // 212
	OpUnknown213          = 0xd5 // 213
	OpUnknown214          = 0xd6 // 214
	OpUnknown215          = 0xd7 // 215
	OpUnknown216          = 0xd8 // 216
	OpUnknown217          = 0xd9 // 217
	OpUnknown218          = 0xda // 218
	OpUnknown219          = 0xdb // 219
	OpUnknown220          = 0xdc // 220
	OpUnknown221          = 0xdd // 221
	OpUnknown222          = 0xde // 222
	OpUnknown223          = 0xdf // 223
	OpUnknown224          = 0xe0 // 224
	OpUnknown225          = 0xe1 // 225
	OpUnknown226          = 0xe2 // 226
	OpUnknown227          = 0xe3 // 227
	OpUnknown228          = 0xe4 // 228
	OpUnknown229          = 0xe5 // 229
	OpUnknown230          = 0xe6 // 230
	OpUnknown231          = 0xe7 // 231
	OpUnknown232          = 0xe8 // 232
	OpUnknown233          = 0xe9 // 233
	OpUnknown234          = 0xea // 234
	OpUnknown235          = 0xeb // 235
	OpUnknown236          = 0xec // 236
	OpUnknown237          = 0xed // 237
	OpUnknown238          = 0xee // 238
	OpUnknown239          = 0xef // 239
	OpUnknown240          = 0xf0 // 240
	OpUnknown241          = 0xf1 // 241
	OpUnknown242          = 0xf2 // 242
	OpUnknown243          = 0xf3 // 243
	OpUnknown244          = 0xf4 // 244
	OpUnknown245          = 0xf5 // 245
	OpUnknown246          = 0xf6 // 246
	OpUnknown247          = 0xf7 // 247
	OpUnknown248          = 0xf8 // 248
	OpUnknown249          = 0xf9 // 249
	OpSMALLINTEGER        = 0xfa // 250 - bitcoin core internal
	OpPUBKEYS             = 0xfb // 251 - bitcoin core internal
	OpUnknown252          = 0xfc // 252
	OpPUBKEYHASH          = 0xfd // 253 - bitcoin core internal
	OpPUBKEY              = 0xfe // 254 - bitcoin core internal
	OpINVALIDOPCODE       = 0xff // 255 - bitcoin core internal
)

// Conditional execution constants.
const (
	OpCondFalse = 0
	OpCondTrue  = 1
	OpCondSkip  = 2
)

// ParsedOpcode represents an opcode that has been parsed and includes any
// potential data associated with it.
type ParsedOpcode struct {
	Opcode *Opcode
	Data   []byte
}

// opcodes holds details about all possible opcodes such as how many bytes
// the opcode and any associated data should take, its human-readable name, and
// the handler function.
var opcodes = [256]Opcode{
	// Data push opcodes.
	OpFalse:     {OpFalse, 1},
	OpData1:     {OpData1, 2},
	OpData2:     {OpData2, 3},
	OpData3:     {OpData3, 4},
	OpData4:     {OpData4, 5},
	OpData5:     {OpData5, 6},
	OpData6:     {OpData6, 7},
	OpData7:     {OpData7, 8},
	OpData8:     {OpData8, 9},
	OpData9:     {OpData9, 10},
	OpData10:    {OpData10, 11},
	OpData11:    {OpData11, 12},
	OpData12:    {OpData12, 13},
	OpData13:    {OpData13, 14},
	OpData14:    {OpData14, 15},
	OpData15:    {OpData15, 16},
	OpData16:    {OpData16, 17},
	OpData17:    {OpData17, 18},
	OpData18:    {OpData18, 19},
	OpData19:    {OpData19, 20},
	OpData20:    {OpData20, 21},
	OpData21:    {OpData21, 22},
	OpData22:    {OpData22, 23},
	OpData23:    {OpData23, 24},
	OpData24:    {OpData24, 25},
	OpData25:    {OpData25, 26},
	OpData26:    {OpData26, 27},
	OpData27:    {OpData27, 28},
	OpData28:    {OpData28, 29},
	OpData29:    {OpData29, 30},
	OpData30:    {OpData30, 31},
	OpData31:    {OpData31, 32},
	OpData32:    {OpData32, 33},
	OpData33:    {OpData33, 34},
	OpData34:    {OpData34, 35},
	OpData35:    {OpData35, 36},
	OpData36:    {OpData36, 37},
	OpData37:    {OpData37, 38},
	OpData38:    {OpData38, 39},
	OpData39:    {OpData39, 40},
	OpData40:    {OpData40, 41},
	OpData41:    {OpData41, 42},
	OpData42:    {OpData42, 43},
	OpData43:    {OpData43, 44},
	OpData44:    {OpData44, 45},
	OpData45:    {OpData45, 46},
	OpData46:    {OpData46, 47},
	OpData47:    {OpData47, 48},
	OpData48:    {OpData48, 49},
	OpData49:    {OpData49, 50},
	OpData50:    {OpData50, 51},
	OpData51:    {OpData51, 52},
	OpData52:    {OpData52, 53},
	OpData53:    {OpData53, 54},
	OpData54:    {OpData54, 55},
	OpData55:    {OpData55, 56},
	OpData56:    {OpData56, 57},
	OpData57:    {OpData57, 58},
	OpData58:    {OpData58, 59},
	OpData59:    {OpData59, 60},
	OpData60:    {OpData60, 61},
	OpData61:    {OpData61, 62},
	OpData62:    {OpData62, 63},
	OpData63:    {OpData63, 64},
	OpData64:    {OpData64, 65},
	OpData65:    {OpData65, 66},
	OpData66:    {OpData66, 67},
	OpData67:    {OpData67, 68},
	OpData68:    {OpData68, 69},
	OpData69:    {OpData69, 70},
	OpData70:    {OpData70, 71},
	OpData71:    {OpData71, 72},
	OpData72:    {OpData72, 73},
	OpData73:    {OpData73, 74},
	OpData74:    {OpData74, 75},
	OpData75:    {OpData75, 76},
	OpPushData1: {OpPushData1, -1},
	OpPushData2: {OpPushData2, -2},
	OpPushData4: {OpPushData4, -4},
	Op1NEGATE:   {Op1NEGATE, 1},
	OpReserved:  {OpReserved, 1},
	OpTRUE:      {OpTRUE, 1},
	Op2:         {Op2, 1},
	Op3:         {Op3, 1},
	Op4:         {Op4, 1},
	Op5:         {Op5, 1},
	Op6:         {Op6, 1},
	Op7:         {Op7, 1},
	Op8:         {Op8, 1},
	Op9:         {Op9, 1},
	Op10:        {Op10, 1},
	Op11:        {Op11, 1},
	Op12:        {Op12, 1},
	Op13:        {Op13, 1},
	Op14:        {Op14, 1},
	Op15:        {Op15, 1},
	Op16:        {Op16, 1},

	// Control opcodes.
	OpNOP:                 {OpNOP, 1},
	OpVER:                 {OpVER, 1},
	OpIF:                  {OpIF, 1},
	OpNOTIF:               {OpNOTIF, 1},
	OpVERIF:               {OpVERIF, 1},
	OpVERNOTIF:            {OpVERNOTIF, 1},
	OpELSE:                {OpELSE, 1},
	OpENDIF:               {OpENDIF, 1},
	OpVERIFY:              {OpVERIFY, 1},
	OpRETURN:              {OpRETURN, 1},
	OpCHECKLOCKTIMEVERIFY: {OpCHECKLOCKTIMEVERIFY, 1},
	OpCHECKSEQUENCEVERIFY: {OpCHECKSEQUENCEVERIFY, 1},

	// Stack opcodes.
	OpTOALTSTACK:   {OpTOALTSTACK, 1},
	OpFROMALTSTACK: {OpFROMALTSTACK, 1},
	Op2DROP:        {Op2DROP, 1},
	Op2DUP:         {Op2DUP, 1},
	Op3DUP:         {Op3DUP, 1},
	Op2OVER:        {Op2OVER, 1},
	Op2ROT:         {Op2ROT, 1},
	Op2SWAP:        {Op2SWAP, 1},
	OpIFDUP:        {OpIFDUP, 1},
	OpDEPTH:        {OpDEPTH, 1},
	OpDROP:         {OpDROP, 1},
	OpDUP:          {OpDUP, 1},
	OpNIP:          {OpNIP, 1},
	OpOVER:         {OpOVER, 1},
	OpPICK:         {OpPICK, 1},
	OpROLL:         {OpROLL, 1},
	OpROT:          {OpROT, 1},
	OpSWAP:         {OpSWAP, 1},
	OpTUCK:         {OpTUCK, 1},

	// Splice opcodes.
	OpCAT:    {OpCAT, 1},
	OpSUBSTR: {OpSUBSTR, 1},
	OpLEFT:   {OpLEFT, 1},
	OpRIGHT:  {OpRIGHT, 1},
	OpSIZE:   {OpSIZE, 1},

	// Bitwise logic opcodes.
	OpINVERT:      {OpINVERT, 1},
	OpAND:         {OpAND, 1},
	OpOR:          {OpOR, 1},
	OpXOR:         {OpXOR, 1},
	OpEQUAL:       {OpEQUAL, 1},
	OpEQUALVERIFY: {OpEQUALVERIFY, 1},
	OpReserved1:   {OpReserved1, 1},
	OpReserved2:   {OpReserved2, 1},

	// Numeric related opcodes.
	Op1ADD:               {Op1ADD, 1},
	Op1SUB:               {Op1SUB, 1},
	Op2MUL:               {Op2MUL, 1},
	Op2DIV:               {Op2DIV, 1},
	OpNEGATE:             {OpNEGATE, 1},
	OpABS:                {OpABS, 1},
	OpNOT:                {OpNOT, 1},
	Op0NOTEQUAL:          {Op0NOTEQUAL, 1},
	OpADD:                {OpADD, 1},
	OpSUB:                {OpSUB, 1},
	OpMUL:                {OpMUL, 1},
	OpDIV:                {OpDIV, 1},
	OpMOD:                {OpMOD, 1},
	OpLSHIFT:             {OpLSHIFT, 1},
	OpRSHIFT:             {OpRSHIFT, 1},
	OpBOOLAND:            {OpBOOLAND, 1},
	OpBOOLOR:             {OpBOOLOR, 1},
	OpNUMEQUAL:           {OpNUMEQUAL, 1},
	OpNUMEQUALVERIFY:     {OpNUMEQUALVERIFY, 1},
	OpNUMNOTEQUAL:        {OpNUMNOTEQUAL, 1},
	OpLESSTHAN:           {OpLESSTHAN, 1},
	OpGREATERTHAN:        {OpGREATERTHAN, 1},
	OpLESSTHANOREQUAL:    {OpLESSTHANOREQUAL, 1},
	OpGREATERTHANOREQUAL: {OpGREATERTHANOREQUAL, 1},
	OpMIN:                {OpMIN, 1},
	OpMAX:                {OpMAX, 1},
	OpWITHIN:             {OpWITHIN, 1},

	// Crypto opcodes.
	OpRIPEMD160:           {OpRIPEMD160, 1},
	OpSHA1:                {OpSHA1, 1},
	OpSHA256:              {OpSHA256, 1},
	OpHASH160:             {OpHASH160, 1},
	OpHASH256:             {OpHASH256, 1},
	OpCODESEPARATOR:       {OpCODESEPARATOR, 1},
	OpCHECKSIG:            {OpCHECKSIG, 1},
	OpCHECKSIGVERIFY:      {OpCHECKSIGVERIFY, 1},
	OpCheckMultiSig:       {OpCheckMultiSig, 1},
	OpCheckMultiSigVerify: {OpCheckMultiSigVerify, 1},

	// Reserved opcodes.
	OpNOP1:  {OpNOP1, 1},
	OpNOP4:  {OpNOP4, 1},
	OpNOP5:  {OpNOP5, 1},
	OpNOP6:  {OpNOP6, 1},
	OpNOP7:  {OpNOP7, 1},
	OpNOP8:  {OpNOP8, 1},
	OpNOP9:  {OpNOP9, 1},
	OpNOP10: {OpNOP10, 1},

	// Undefined opcodes.
	OpUnknown186: {OpUnknown186, 1},
	OpUnknown187: {OpUnknown187, 1},
	OpUnknown188: {OpUnknown188, 1},
	OpUnknown189: {OpUnknown189, 1},
	OpUnknown190: {OpUnknown190, 1},
	OpUnknown191: {OpUnknown191, 1},
	OpUnknown192: {OpUnknown192, 1},
	OpUnknown193: {OpUnknown193, 1},
	OpUnknown194: {OpUnknown194, 1},
	OpUnknown195: {OpUnknown195, 1},
	OpUnknown196: {OpUnknown196, 1},
	OpUnknown197: {OpUnknown197, 1},
	OpUnknown198: {OpUnknown198, 1},
	OpUnknown199: {OpUnknown199, 1},
	OpUnknown200: {OpUnknown200, 1},
	OpUnknown201: {OpUnknown201, 1},
	OpUnknown202: {OpUnknown202, 1},
	OpUnknown203: {OpUnknown203, 1},
	OpUnknown204: {OpUnknown204, 1},
	OpUnknown205: {OpUnknown205, 1},
	OpUnknown206: {OpUnknown206, 1},
	OpUnknown207: {OpUnknown207, 1},
	OpUnknown208: {OpUnknown208, 1},
	OpUnknown209: {OpUnknown209, 1},
	OpUnknown210: {OpUnknown210, 1},
	OpUnknown211: {OpUnknown211, 1},
	OpUnknown212: {OpUnknown212, 1},
	OpUnknown213: {OpUnknown213, 1},
	OpUnknown214: {OpUnknown214, 1},
	OpUnknown215: {OpUnknown215, 1},
	OpUnknown216: {OpUnknown216, 1},
	OpUnknown217: {OpUnknown217, 1},
	OpUnknown218: {OpUnknown218, 1},
	OpUnknown219: {OpUnknown219, 1},
	OpUnknown220: {OpUnknown220, 1},
	OpUnknown221: {OpUnknown221, 1},
	OpUnknown222: {OpUnknown222, 1},
	OpUnknown223: {OpUnknown223, 1},
	OpUnknown224: {OpUnknown224, 1},
	OpUnknown225: {OpUnknown225, 1},
	OpUnknown226: {OpUnknown226, 1},
	OpUnknown227: {OpUnknown227, 1},
	OpUnknown228: {OpUnknown228, 1},
	OpUnknown229: {OpUnknown229, 1},
	OpUnknown230: {OpUnknown230, 1},
	OpUnknown231: {OpUnknown231, 1},
	OpUnknown232: {OpUnknown232, 1},
	OpUnknown233: {OpUnknown233, 1},
	OpUnknown234: {OpUnknown234, 1},
	OpUnknown235: {OpUnknown235, 1},
	OpUnknown236: {OpUnknown236, 1},
	OpUnknown237: {OpUnknown237, 1},
	OpUnknown238: {OpUnknown238, 1},
	OpUnknown239: {OpUnknown239, 1},
	OpUnknown240: {OpUnknown240, 1},
	OpUnknown241: {OpUnknown241, 1},
	OpUnknown242: {OpUnknown242, 1},
	OpUnknown243: {OpUnknown243, 1},
	OpUnknown244: {OpUnknown244, 1},
	OpUnknown245: {OpUnknown245, 1},
	OpUnknown246: {OpUnknown246, 1},
	OpUnknown247: {OpUnknown247, 1},
	OpUnknown248: {OpUnknown248, 1},
	OpUnknown249: {OpUnknown249, 1},

	// Bitcoin Core internal use opcode.  Defined here for completeness.
	OpSMALLINTEGER: {OpSMALLINTEGER, 1},
	OpPUBKEYS:      {OpPUBKEYS, 1},
	OpUnknown252:   {OpUnknown252, 1},
	OpPUBKEYHASH:   {OpPUBKEYHASH, 1},
	OpPUBKEY:       {OpPUBKEY, 1},

	OpINVALIDOPCODE: {OpINVALIDOPCODE, 1},
}
