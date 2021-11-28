package item

import (
	"github.com/artheus/go-minecraft/core/block"
	"github.com/artheus/go-minecraft/core/texture"
	"log"
)

var (
	Tex = NewItemHub()
)

func MakeFaceTexture(idx int) texture.FaceTexture {
	const textureColums = 16
	var m = 1 / float32(textureColums)
	dx, dy := float32(idx%textureColums)*m, float32(idx/textureColums)*m
	n := float32(1 / 2048.0)
	m -= n
	return [6][2]float32{
		{dx + n, dy + n},
		{dx + m, dy + n},
		{dx + m, dy + m},
		{dx + m, dy + m},
		{dx + n, dy + m},
		{dx + n, dy + n},
	}
}

type Hub struct {
	tex map[string]*texture.BlockTexture
}

func NewItemHub() *Hub {
	return &Hub{
		tex: make(map[string]*texture.BlockTexture),
	}
}

func (h *Hub) Tex() map[string]*texture.BlockTexture {
	return h.tex
}

func (h *Hub) AddTexture(w string, l, r, u, d, f, b int) {
	h.tex[w] = &texture.BlockTexture{
		Left:  MakeFaceTexture(l),
		Right: MakeFaceTexture(r),
		Up:    MakeFaceTexture(u),
		Down:  MakeFaceTexture(d),
		Front: MakeFaceTexture(f),
		Back:  MakeFaceTexture(b),
	}
}

func (h *Hub) Texture(w string) *texture.BlockTexture {
	t, ok := h.tex[w]
	if !ok {
		log.Printf("%d not found", w)
		return h.tex[block.AirID]
	}
	return t
}

func LoadTextureDesc() error {
	for w, f := range itemDesc {
		Tex.AddTexture(w, f[0], f[1], f[2], f[3], f[4], f[5])
	}
	return nil
}

// w => left, right, top, bottom, front, back
var itemDesc = map[string][6]int{
	"core:air":  {0, 0, 0, 0, 0, 0},      // air
	"core:grass_block":  {16, 16, 32, 0, 16, 16}, // grass block
	"core:sand":  {1, 1, 1, 1, 1, 1},      // sand block
	"core:stone_bricks":  {2, 2, 2, 2, 2, 2},      // stone brick
	"core:bricks":  {3, 3, 3, 3, 3, 3},      // bricks
	"core:wood":  {20, 20, 36, 4, 20, 20}, // wood
	"core:stone":  {5, 5, 5, 5, 5, 5},
	"core:dirt":  {6, 6, 6, 6, 6, 6},      // dirt block
	"core:planks":  {7, 7, 7, 7, 7, 7},      // planks
	"core:grass_block_snowy":  {24, 24, 40, 8, 24, 24},
	"core:glass": {9, 9, 9, 9, 9, 9},
	"core:cobblestone": {10, 10, 10, 10, 10, 10},
	//12: {11, 11, 11, 11, 11, 11},
	//13: {12, 12, 12, 12, 12, 12},
	//14: {13, 13, 13, 13, 13, 13},
	"core:leaves": {14, 14, 14, 14, 14, 14}, // leaves
	"core:cloud": {15, 15, 15, 15, 15, 15},
	"core:grass": {48, 48, 0, 0, 48, 48},   // grass
	"core:dandelion": {49, 49, 0, 0, 49, 49},
	"core:tulip": {50, 50, 0, 0, 50, 50},
	"core:lily": {51, 51, 0, 0, 51, 51},
	"core:sunflower": {52, 52, 0, 0, 52, 52},
	"core:cotton": {53, 53, 0, 0, 53, 53},
	"core:violet": {54, 54, 0, 0, 54, 54},
	//24: {0, 0, 0, 0, 0, 0},
	//25: {0, 0, 0, 0, 0, 0},
	//26: {0, 0, 0, 0, 0, 0},
	//27: {0, 0, 0, 0, 0, 0},
	//28: {0, 0, 0, 0, 0, 0},
	//29: {0, 0, 0, 0, 0, 0},
	//30: {0, 0, 0, 0, 0, 0},
	//31: {0, 0, 0, 0, 0, 0},
	//32: {176, 176, 176, 176, 176, 176},
	//33: {177, 177, 177, 177, 177, 177},
	//34: {178, 178, 178, 178, 178, 178},
	//35: {179, 179, 179, 179, 179, 179},
	//36: {180, 180, 180, 180, 180, 180},
	//37: {181, 181, 181, 181, 181, 181},
	//38: {182, 182, 182, 182, 182, 182},
	//39: {183, 183, 183, 183, 183, 183},
	//40: {184, 184, 184, 184, 184, 184},
	//41: {185, 185, 185, 185, 185, 185},
	//42: {186, 186, 186, 186, 186, 186},
	//43: {187, 187, 187, 187, 187, 187},
	//44: {188, 188, 188, 188, 188, 188},
	//45: {189, 189, 189, 189, 189, 189},
	//46: {190, 190, 190, 190, 190, 190},
	//47: {191, 191, 191, 191, 191, 191},
	//48: {192, 192, 192, 192, 192, 192},
	//49: {193, 193, 193, 193, 193, 193},
	//50: {194, 194, 194, 194, 194, 194},
	//51: {195, 195, 195, 195, 195, 195},
	//52: {196, 196, 196, 196, 196, 196},
	//53: {197, 197, 197, 197, 197, 197},
	//54: {198, 198, 198, 198, 198, 198},
	//55: {199, 199, 199, 199, 199, 199},
	//56: {200, 200, 200, 200, 200, 200},
	//57: {201, 201, 201, 201, 201, 201},
	//58: {202, 202, 202, 202, 202, 202},
	//59: {203, 203, 203, 203, 203, 203},
	//60: {204, 204, 204, 204, 204, 204},
	//61: {205, 205, 205, 205, 205, 205},
	//62: {206, 206, 206, 206, 206, 206},
	//63: {207, 207, 207, 207, 207, 207},
	"core:player": {226, 224, 241, 209, 227, 225},
}

var AvailableItems = []int{
	1,
	2,
	3,
	4,
	5,
	6,
	7,
	8,
	9,
	10,
	11,
	12,
	13,
	14,
	15,
	16,
	17,
	18,
	19,
	20,
	21,
	22,
	23,
	32,
	33,
	34,
	35,
	36,
	37,
	38,
	39,
	40,
	41,
	42,
	43,
	44,
	45,
	46,
	47,
	48,
	49,
	50,
	51,
	52,
	53,
	54,
	55,
	56,
	57,
	58,
	59,
	60,
	61,
	62,
	63,
	64,
}
