package main

import (
	"fmt"
	"strconv"
)

func push(p *PDAProcessor, val string) { //Pushes the token to provided stack position
	p.Stack = append(p.Stack, val)
}
func pop(p *PDAProcessor) { //pops the stack top upon encountering 1 in token string
	p.Stack = p.Stack[:len(p.Stack)-1]
}
func peekClient(p *PDAProcessor, k int) []string { //provides peek of top k elements of stack.

	Stacktop := []string{}
	l := len(p.Stack)
	if l <= k {
		Stacktop = p.Stack
	} else if k == 1 {
		Stacktop = append(Stacktop, p.Stack[l-1])
	} else {
		Stacktop = p.Stack[l-k : l-1]
	}
	return Stacktop
}
func resetClient(p *PDAProcessor) { //This resets PDA and sets it to default
	p.Stack = make([]string, 0)
	p.Current_State = p.Start_state
	p.Next_Pos = 0
	p.Hold_Que = make([]HoldStruct, 0)
}
func openClient(id string, p PDAProcessor) bool { //This opens PDA and takes json input from curl-PUT
	p.Id = id
	_, found := cache[id]
	if !found {
		resetClient(&p)
		cache[id] = p
		return true
	}
	return false
}
func isAccepted(proc PDAProcessor) bool { //This function checks whether token string is accepted by PDA.
	flag := false
	accepting_states := proc.Accepting_states
	cs := proc.Current_State
	if len(proc.Stack) == 0 && len(proc.Hold_Que) == 0 {
		for i := 0; i < len(accepting_states); i++ {
			if cs == accepting_states[i] {
				flag = true
				break
			}
		}
	}
	return flag
}
func que_tokens(proc PDAProcessor) int { //This function ensures whether token is inserted in given position and holds a token if the order is mismatched.
	defer wg.Done() //Defering the execution of whjere wait counter is decremented by 1 to ensure smooth client server interaction.
	token_blocked := -1
	token_processed := false
	for {
		if len(proc.Hold_Que) == 0 {
			break
		}
		proc = cache[proc.Id]
		hold_back := proc.Hold_Que[len(proc.Hold_Que)-1]
		pos_int, _ := strconv.Atoi(hold_back.Hold_Pos) //converting to integer type
		if proc.Next_Pos == pos_int {
			token_processed = putClient(proc, hold_back.Hold_Token)
			if !token_processed {
				token_blocked = pos_int
				break
			} else {
				proc = cache[proc.Id]
				proc.Hold_Que = proc.Hold_Que[:len(proc.Hold_Que)-1]
				cache[proc.Id] = proc
			}
		} else {
			break
		}
	}
	return token_blocked
}
func putClient(proc PDAProcessor, token string) bool { //This function takes input token string and depending upon current input and state of PDA stack operations are performed.
	transitions := proc.Transitions
	tran_len := len(transitions)
	token_processed := false
	for j := 0; j < tran_len; j++ {
		var current_state = transitions[j][0]
		var input = transitions[j][1]
		var top_of_stack = transitions[j][2]
		var target_state = transitions[j][3]
		var action_item = transitions[j][4]
		var currentStackSymbol = ""
		var top = peekClient(&proc, 1)
		if len(top) >= 1 {
			currentStackSymbol = top[0]
		}
		if input == "null" && current_state == proc.Current_State && top_of_stack == "null" && action_item == "null" {
			fmt.Printf("Current State-->%v, Stack-->%v, Transitoned State-->%v ", proc.Current_State, proc.Stack, target_state)
			fmt.Println("No stack operations performed--->>dead transition(RPC SEMANTIC)")
			proc.Current_State = target_state
		}
		if current_state == proc.Current_State && input == token {
			if action_item != "null" && top_of_stack == "null" { //PUSH to Stack
				fmt.Printf("Current State-->%v, Stack-->%v, Transitoned State-->%v ", proc.Current_State, proc.Stack, target_state)
				fmt.Println("Pushed-->", action_item, " to stack")
				proc.Next_Pos = proc.Next_Pos + 1
				token_processed = true
				proc.Current_State = target_state
				push(&proc, action_item)
				break
			} else if action_item != "null" && top_of_stack == currentStackSymbol { //PUSH to stack
				fmt.Printf("Current State-->%v, Stack-->%v, Transitoned State-->%v ", proc.Current_State, proc.Stack, target_state)
				fmt.Println("Pushed-->", action_item, " to stack")
				proc.Next_Pos = proc.Next_Pos + 1
				token_processed = true
				proc.Current_State = target_state
				push(&proc, action_item)
				break

			} else if action_item == "null" && top_of_stack == currentStackSymbol { //POP the stack
				pop(&proc)
				fmt.Printf("Current State-->%v, Stack-->%v, Transitoned State-->%v ", proc.Current_State, proc.Stack, target_state)
				fmt.Println("===>PDA Stack has been popped")
				proc.Next_Pos = proc.Next_Pos + 1
				token_processed = true
				proc.Current_State = target_state
				break

			} else if top_of_stack == "null" { //When no token element on top of stack.
				fmt.Printf("Current State-->%v, Stack-->%v, Transitoned State-->%v ", proc.Current_State, proc.Stack, target_state)
				fmt.Println("No stack operations performed====>>Consumed input token from URL")
				proc.Current_State = target_state
				proc.Next_Pos = proc.Next_Pos + 1
				token_processed = true
				break
			}
		}
	}
	cache[proc.Id] = proc
	return token_processed
}
func make_transition(proc *PDAProcessor, transition_count int) { //Pushes EOS into stack to initiate the state transition.
	transitions := proc.Transitions
	target_state := ""
	input := ""
	top_of_stack := ""
	action_item := ""
	for j := 0; j < len(transitions); j++ {
		if transitions[j][0] == proc.Current_State {
			input = transitions[j][1]
			top_of_stack = transitions[j][2]
			target_state = transitions[j][3]
			action_item = transitions[j][4]
			break
		}
	}
	if input == "null" && top_of_stack == "null" {
		fmt.Println("Current State ", proc.Current_State)
		fmt.Println("Pushed $ in the stack")
		push(proc, action_item)
		proc.Current_State = target_state
		fmt.Println("Transitioning to State--->", proc.Current_State)
		transition_count = transition_count + 1
		fmt.Println()
	}
}
