package channelling

type TurnDataCreator interface {
	CreateTurnData(*Session) *DataTurn
}
