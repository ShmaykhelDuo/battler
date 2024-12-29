package db

type TxIsolation int

const (
	TxIsolationDefault TxIsolation = iota
	TxIsolationReadUncommitted
	TxIsolationReadCommitted
	TxIsolationRepeatableRead
	TxIsolationSerializable
)
