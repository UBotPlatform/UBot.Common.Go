package ubot

type EventResultType int

const IgnoreEvent EventResultType = 0
const ContinueEvent EventResultType = 1
const CompleteEvent EventResultType = 2
