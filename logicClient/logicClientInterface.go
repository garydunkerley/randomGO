package logic

import "fmt"

type LogicClient interface {
	RequestBoard(BoardParams) BoardTop
	StartGame(BoardTop) error
	SendMove(MoveInput) GameUpdate
	SuggestDeadStones() DeadStoneProposal
	SendScoreInput(MoveInput) 
}

type BoardParams struct {


type LogicServer interface {
	
}

