package main

type PDAProcessor struct { // Structure for PDA processor, we have removed jso specification since we are sending request from client to server.
	Id               string
	Name             string
	States           []string
	Input_alphabet   []string
	Stack_alphabet   []string
	Accepting_states []string
	Start_state      string
	Transitions      [][]string
	Eos              string
	Stack            []string
	Current_State    string
	Next_Pos         int
	Hold_Que         []HoldStruct
}

type PDAInfo struct { //Structure stores metadata of PDA ID and Name
	Id   string
	Name string
}

type Snapshot struct { //Structure that stores current snapshot of PDA
	Topk          []string
	Current_State string
	pdaCLOCK      int
	Hold_Queue    []HoldStruct
}

type Token struct { //structure that stores token string
	Token string
}

type HoldStruct struct { //sturcture that stores holding token and its position.
	Hold_Pos   string
	Hold_Token string
}
