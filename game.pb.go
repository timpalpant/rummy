// Code generated by protoc-gen-go.
// source: game.proto
// DO NOT EDIT!

/*
Package rummy is a generated protocol buffer package.

It is generated from these files:
	game.proto
	service.proto

It has these top-level messages:
	Meld
	PlayerState
	GameState
	GameEvent
	CreateGameRequest
	CreateGameResponse
	JoinGameRequest
	JoinGameResponse
	StartGameRequest
	StartGameResponse
	GetGameStateRequest
	GetHandCardsRequest
	GetHandCardsResponse
	SubscribeGameRequest
	PickUpStockRequest
	PickUpStockResponse
	PickUpDiscardRequest
	PickUpDiscardResponse
	PlayCardsRequest
	PlayCardsResponse
	DiscardCardRequest
	DiscardCardResponse
	CallRummyRequest
	CallRummyResponse
*/
package rummy

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import deck "github.com/timpalpant/rummy/deck"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type GameState_TurnState int32

const (
	GameState_TURN_START      GameState_TurnState = 0
	GameState_PICKED_UP_CARDS GameState_TurnState = 1
	GameState_PLAYED_CARDS    GameState_TurnState = 2
)

var GameState_TurnState_name = map[int32]string{
	0: "TURN_START",
	1: "PICKED_UP_CARDS",
	2: "PLAYED_CARDS",
}
var GameState_TurnState_value = map[string]int32{
	"TURN_START":      0,
	"PICKED_UP_CARDS": 1,
	"PLAYED_CARDS":    2,
}

func (x GameState_TurnState) String() string {
	return proto.EnumName(GameState_TurnState_name, int32(x))
}
func (GameState_TurnState) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{2, 0} }

type GameEvent_Type int32

const (
	GameEvent_UNKNOWN_TYPE    GameEvent_Type = 0
	GameEvent_TURN_START      GameEvent_Type = 1
	GameEvent_PICK_UP_STOCK   GameEvent_Type = 2
	GameEvent_PICK_UP_DISCARD GameEvent_Type = 3
	GameEvent_PLAY_CARDS      GameEvent_Type = 4
	GameEvent_DISCARD         GameEvent_Type = 5
	GameEvent_GAME_OVER       GameEvent_Type = 6
)

var GameEvent_Type_name = map[int32]string{
	0: "UNKNOWN_TYPE",
	1: "TURN_START",
	2: "PICK_UP_STOCK",
	3: "PICK_UP_DISCARD",
	4: "PLAY_CARDS",
	5: "DISCARD",
	6: "GAME_OVER",
}
var GameEvent_Type_value = map[string]int32{
	"UNKNOWN_TYPE":    0,
	"TURN_START":      1,
	"PICK_UP_STOCK":   2,
	"PICK_UP_DISCARD": 3,
	"PLAY_CARDS":      4,
	"DISCARD":         5,
	"GAME_OVER":       6,
}

func (x GameEvent_Type) String() string {
	return proto.EnumName(GameEvent_Type_name, int32(x))
}
func (GameEvent_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{3, 0} }

type Meld struct {
	Cards []*deck.Card `protobuf:"bytes,1,rep,name=cards" json:"cards,omitempty"`
}

func (m *Meld) Reset()                    { *m = Meld{} }
func (m *Meld) String() string            { return proto.CompactTextString(m) }
func (*Meld) ProtoMessage()               {}
func (*Meld) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Meld) GetCards() []*deck.Card {
	if m != nil {
		return m.Cards
	}
	return nil
}

type PlayerState struct {
	Id             int32        `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Name           string       `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Melds          []*Meld      `protobuf:"bytes,3,rep,name=melds" json:"melds,omitempty"`
	Rummies        []*deck.Card `protobuf:"bytes,4,rep,name=rummies" json:"rummies,omitempty"`
	NumCardsInHand int32        `protobuf:"varint,5,opt,name=num_cards_in_hand,json=numCardsInHand" json:"num_cards_in_hand,omitempty"`
	CurrentScore   int32        `protobuf:"varint,6,opt,name=current_score,json=currentScore" json:"current_score,omitempty"`
}

func (m *PlayerState) Reset()                    { *m = PlayerState{} }
func (m *PlayerState) String() string            { return proto.CompactTextString(m) }
func (*PlayerState) ProtoMessage()               {}
func (*PlayerState) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *PlayerState) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *PlayerState) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *PlayerState) GetMelds() []*Meld {
	if m != nil {
		return m.Melds
	}
	return nil
}

func (m *PlayerState) GetRummies() []*deck.Card {
	if m != nil {
		return m.Rummies
	}
	return nil
}

func (m *PlayerState) GetNumCardsInHand() int32 {
	if m != nil {
		return m.NumCardsInHand
	}
	return 0
}

func (m *PlayerState) GetCurrentScore() int32 {
	if m != nil {
		return m.CurrentScore
	}
	return 0
}

type GameState struct {
	NumCardsInStock   int32               `protobuf:"varint,1,opt,name=num_cards_in_stock,json=numCardsInStock" json:"num_cards_in_stock,omitempty"`
	DiscardPile       []*deck.Card        `protobuf:"bytes,2,rep,name=discard_pile,json=discardPile" json:"discard_pile,omitempty"`
	AggregatedMelds   []*Meld             `protobuf:"bytes,3,rep,name=aggregated_melds,json=aggregatedMelds" json:"aggregated_melds,omitempty"`
	Players           []*PlayerState      `protobuf:"bytes,4,rep,name=players" json:"players,omitempty"`
	Turn              int32               `protobuf:"varint,6,opt,name=turn" json:"turn,omitempty"`
	CurrentPlayerTurn int32               `protobuf:"varint,7,opt,name=current_player_turn,json=currentPlayerTurn" json:"current_player_turn,omitempty"`
	TurnState         GameState_TurnState `protobuf:"varint,8,opt,name=turn_state,json=turnState,enum=rummy.GameState_TurnState" json:"turn_state,omitempty"`
	GameOver          bool                `protobuf:"varint,9,opt,name=game_over,json=gameOver" json:"game_over,omitempty"`
}

func (m *GameState) Reset()                    { *m = GameState{} }
func (m *GameState) String() string            { return proto.CompactTextString(m) }
func (*GameState) ProtoMessage()               {}
func (*GameState) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *GameState) GetNumCardsInStock() int32 {
	if m != nil {
		return m.NumCardsInStock
	}
	return 0
}

func (m *GameState) GetDiscardPile() []*deck.Card {
	if m != nil {
		return m.DiscardPile
	}
	return nil
}

func (m *GameState) GetAggregatedMelds() []*Meld {
	if m != nil {
		return m.AggregatedMelds
	}
	return nil
}

func (m *GameState) GetPlayers() []*PlayerState {
	if m != nil {
		return m.Players
	}
	return nil
}

func (m *GameState) GetTurn() int32 {
	if m != nil {
		return m.Turn
	}
	return 0
}

func (m *GameState) GetCurrentPlayerTurn() int32 {
	if m != nil {
		return m.CurrentPlayerTurn
	}
	return 0
}

func (m *GameState) GetTurnState() GameState_TurnState {
	if m != nil {
		return m.TurnState
	}
	return GameState_TURN_START
}

func (m *GameState) GetGameOver() bool {
	if m != nil {
		return m.GameOver
	}
	return false
}

type GameEvent struct {
	PlayerId int32          `protobuf:"varint,1,opt,name=player_id,json=playerId" json:"player_id,omitempty"`
	Type     GameEvent_Type `protobuf:"varint,2,opt,name=type,enum=rummy.GameEvent_Type" json:"type,omitempty"`
	Cards    []*deck.Card   `protobuf:"bytes,3,rep,name=cards" json:"cards,omitempty"`
	Score    int32          `protobuf:"varint,4,opt,name=score" json:"score,omitempty"`
}

func (m *GameEvent) Reset()                    { *m = GameEvent{} }
func (m *GameEvent) String() string            { return proto.CompactTextString(m) }
func (*GameEvent) ProtoMessage()               {}
func (*GameEvent) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *GameEvent) GetPlayerId() int32 {
	if m != nil {
		return m.PlayerId
	}
	return 0
}

func (m *GameEvent) GetType() GameEvent_Type {
	if m != nil {
		return m.Type
	}
	return GameEvent_UNKNOWN_TYPE
}

func (m *GameEvent) GetCards() []*deck.Card {
	if m != nil {
		return m.Cards
	}
	return nil
}

func (m *GameEvent) GetScore() int32 {
	if m != nil {
		return m.Score
	}
	return 0
}

func init() {
	proto.RegisterType((*Meld)(nil), "rummy.Meld")
	proto.RegisterType((*PlayerState)(nil), "rummy.PlayerState")
	proto.RegisterType((*GameState)(nil), "rummy.GameState")
	proto.RegisterType((*GameEvent)(nil), "rummy.GameEvent")
	proto.RegisterEnum("rummy.GameState_TurnState", GameState_TurnState_name, GameState_TurnState_value)
	proto.RegisterEnum("rummy.GameEvent_Type", GameEvent_Type_name, GameEvent_Type_value)
}

func init() { proto.RegisterFile("game.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 576 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x93, 0xc1, 0x6e, 0xda, 0x40,
	0x10, 0x86, 0x63, 0x63, 0x07, 0x3c, 0x24, 0xc4, 0x6c, 0x1a, 0xc9, 0x4a, 0x2f, 0xae, 0xdb, 0x83,
	0xa3, 0xb6, 0xae, 0x44, 0xa5, 0x4a, 0x3d, 0x52, 0xb0, 0x52, 0x44, 0x01, 0xcb, 0x36, 0xad, 0x72,
	0x5a, 0xb9, 0xec, 0x8a, 0x5a, 0xc1, 0x06, 0xad, 0x0d, 0x12, 0x52, 0x5f, 0xa0, 0x0f, 0xd3, 0x17,
	0xe9, 0x53, 0x55, 0xbb, 0x6b, 0x43, 0x9a, 0x2a, 0x17, 0x58, 0xfe, 0xf9, 0x76, 0x66, 0xfe, 0xd9,
	0x01, 0x60, 0x99, 0x64, 0xd4, 0xdb, 0xb0, 0x75, 0xb9, 0x46, 0x3a, 0xdb, 0x66, 0xd9, 0xfe, 0xfa,
	0x4a, 0x7c, 0xbd, 0x23, 0x74, 0x71, 0x2f, 0x3e, 0x64, 0xd4, 0x71, 0x41, 0x9b, 0xd0, 0x15, 0x41,
	0x36, 0xe8, 0x8b, 0x84, 0x91, 0xc2, 0x52, 0xec, 0x86, 0xdb, 0xee, 0x81, 0x27, 0x98, 0x41, 0xc2,
	0x48, 0x28, 0x03, 0xce, 0x1f, 0x05, 0xda, 0xc1, 0x2a, 0xd9, 0x53, 0x16, 0x95, 0x49, 0x49, 0x51,
	0x07, 0xd4, 0x94, 0x58, 0x8a, 0xad, 0xb8, 0x7a, 0xa8, 0xa6, 0x04, 0x21, 0xd0, 0xf2, 0x24, 0xa3,
	0x96, 0x6a, 0x2b, 0xae, 0x11, 0x8a, 0x33, 0x7a, 0x01, 0x7a, 0x46, 0x57, 0xa4, 0xb0, 0x1a, 0x22,
	0x6b, 0xdb, 0x13, 0x4d, 0x78, 0xbc, 0x62, 0x28, 0x23, 0xe8, 0x15, 0x34, 0xb9, 0x98, 0xd2, 0xc2,
	0xd2, 0xfe, 0x2b, 0x5d, 0x87, 0xd0, 0x0d, 0x74, 0xf3, 0x6d, 0x86, 0x45, 0x27, 0x38, 0xcd, 0xf1,
	0x8f, 0x24, 0x27, 0x96, 0x2e, 0x6a, 0x77, 0xf2, 0x6d, 0xc6, 0xe1, 0x62, 0x94, 0x7f, 0x4e, 0x72,
	0x82, 0x5e, 0xc2, 0xf9, 0x62, 0xcb, 0x18, 0xcd, 0x4b, 0x5c, 0x2c, 0xd6, 0x8c, 0x5a, 0xa7, 0x02,
	0x3b, 0xab, 0xc4, 0x88, 0x6b, 0xce, 0xef, 0x06, 0x18, 0xb7, 0x49, 0x46, 0xa5, 0x95, 0xd7, 0x80,
	0xfe, 0xc9, 0x5e, 0x94, 0xeb, 0xc5, 0x7d, 0x65, 0xed, 0xe2, 0x98, 0x3e, 0xe2, 0x32, 0x7a, 0x0b,
	0x67, 0x24, 0x2d, 0x38, 0x8b, 0x37, 0xe9, 0x8a, 0xfb, 0x7d, 0xdc, 0x75, 0xbb, 0x8a, 0x07, 0xe9,
	0x8a, 0xa2, 0x0f, 0x60, 0x26, 0xcb, 0x25, 0xa3, 0xcb, 0xa4, 0xa4, 0x04, 0x3f, 0x39, 0x8d, 0x8b,
	0x23, 0x34, 0x11, 0x73, 0x79, 0x03, 0xcd, 0x8d, 0x98, 0x76, 0x3d, 0x17, 0x54, 0xe1, 0x0f, 0xde,
	0x20, 0xac, 0x11, 0x3e, 0xfc, 0x72, 0xcb, 0xf2, 0xca, 0xab, 0x38, 0x23, 0x0f, 0x2e, 0xeb, 0x41,
	0x48, 0x0c, 0x0b, 0xa4, 0x29, 0x90, 0x6e, 0x15, 0x92, 0xd9, 0x62, 0xce, 0x7f, 0x04, 0xe0, 0x00,
	0x2e, 0x78, 0x6a, 0xab, 0x65, 0x2b, 0x6e, 0xa7, 0x77, 0x5d, 0x15, 0x3d, 0xcc, 0xca, 0xe3, 0xa8,
	0x2c, 0x6e, 0x94, 0xf5, 0x11, 0x3d, 0x07, 0x83, 0x6f, 0x1c, 0x5e, 0xef, 0x28, 0xb3, 0x0c, 0x5b,
	0x71, 0x5b, 0x61, 0x8b, 0x0b, 0xb3, 0x1d, 0x65, 0xce, 0x27, 0x30, 0x0e, 0x97, 0x50, 0x07, 0x20,
	0x9e, 0x87, 0x53, 0x1c, 0xc5, 0xfd, 0x30, 0x36, 0x4f, 0xd0, 0x25, 0x5c, 0x04, 0xa3, 0xc1, 0xd8,
	0x1f, 0xe2, 0x79, 0x80, 0x07, 0xfd, 0x70, 0x18, 0x99, 0x0a, 0x32, 0xe1, 0x2c, 0xf8, 0xd2, 0xbf,
	0xf3, 0x87, 0x95, 0xa2, 0x3a, 0xbf, 0x54, 0xf9, 0x5e, 0xfe, 0x8e, 0xe6, 0x25, 0x2f, 0x57, 0x39,
	0x3a, 0x6c, 0x60, 0x4b, 0x0a, 0x23, 0x82, 0x6e, 0x40, 0x2b, 0xf7, 0x1b, 0xb9, 0x87, 0x9d, 0xde,
	0xd5, 0x03, 0x03, 0xe2, 0xb2, 0x17, 0xef, 0x37, 0x34, 0x14, 0xc8, 0x71, 0xe9, 0x1b, 0x4f, 0x2c,
	0x3d, 0x7a, 0x06, 0xba, 0x5c, 0x22, 0x4d, 0x54, 0x91, 0x3f, 0x9c, 0x9f, 0xa0, 0xf1, 0x2c, 0xbc,
	0xcf, 0xf9, 0x74, 0x3c, 0x9d, 0x7d, 0x9b, 0xe2, 0xf8, 0x2e, 0xf0, 0xcd, 0x93, 0x47, 0xf6, 0x14,
	0xd4, 0x85, 0x73, 0x6e, 0x8f, 0x9b, 0x8b, 0xe2, 0xd9, 0x60, 0x6c, 0xaa, 0xb5, 0x63, 0x2e, 0x0d,
	0x47, 0x11, 0x37, 0x68, 0x36, 0xf8, 0x3d, 0xee, 0xb8, 0xf2, 0xab, 0xa1, 0x36, 0x34, 0xeb, 0xa0,
	0x8e, 0xce, 0xc1, 0xb8, 0xed, 0x4f, 0x7c, 0x3c, 0xfb, 0xea, 0x87, 0xe6, 0xe9, 0xf7, 0x53, 0xf1,
	0xcf, 0x7d, 0xff, 0x37, 0x00, 0x00, 0xff, 0xff, 0xdc, 0x76, 0x6e, 0x23, 0xe5, 0x03, 0x00, 0x00,
}
