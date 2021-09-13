package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

var wg sync.WaitGroup
var cache = make(map[string]PDAProcessor)

func peekServer(w http.ResponseWriter, r *http.Request) { //Provides top k elements of stack
	fmt.Println("\nApi Hit====>PEEK")
	var vars = mux.Vars(r)
	var id = vars["id"]
	var kstring = vars["k"]
	k, _ := strconv.Atoi(kstring) //parsing string to int conversion
	proc := cache[id]
	top := peekClient(&proc, k)
	json.NewEncoder(w).Encode(top)

}
func reset(w http.ResponseWriter, r *http.Request) { //resets the pda to its default state
	fmt.Println("\nApi Hit====>RESET")
	var vars = mux.Vars(r)
	var id = vars["id"]
	p := cache[id]
	resetClient(&p)
	cache[id] = p
}

func create(w http.ResponseWriter, r *http.Request) { //creates PDA processor with given id on server
	fmt.Println("\n Api Hit====>create")
	var p PDAProcessor
	var vars = mux.Vars(r)
	var id = vars["id"]
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	created := openClient(id, p)
	if created {
		json.NewEncoder(w).Encode("PDA with given id is created.")
	} else {
		json.NewEncoder(w).Encode("Cant create already exists")
	}
}

func openServer(w http.ResponseWriter, r *http.Request) { //opens the pda with given id
	fmt.Println("\n Api Hit====>open")
	var vars = mux.Vars(r)
	var id = vars["id"]
	if id != "" {
		json.NewEncoder(w).Encode("PDA with given id is opened.")
	} else {
		json.NewEncoder(w).Encode("PDA with given id is not created so cannot open.")
	}
}

func AllPdas(w http.ResponseWriter, r *http.Request) { //Prints all the pdas on server
	fmt.Println("Api Hit====>All Pdas available at the server:")
	var pdalist []PDAInfo
	for key, value := range cache {
		_ = key
		info := PDAInfo{
			Name: value.Name,
		}
		pdalist = append(pdalist, info)
	}
	json.NewEncoder(w).Encode(pdalist)
}

// Function to check if the input string has been accepted by the pda
func is_accepted(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n Api Hit====>isAccepted")
	var vars = mux.Vars(r)
	var id = vars["id"]
	proc := cache[id]
	flag := isAccepted(proc)
	if flag {
		json.NewEncoder(w).Encode("Input tokens successfully Accepted")
	} else {
		json.NewEncoder(w).Encode("Input tokens Rejected by the PDA")
	}
}

// The done returns the final status of the current state and the stack after the input string is processed.
func done(proc PDAProcessor, is_accepted bool, transition_count int) {
	fmt.Println("pda = ", proc.Name, "::total_clock = ", transition_count, "::method = is_accepted = ", is_accepted, "::Current State = ", proc.Current_State)
	fmt.Println("Current_state: ", proc.Current_State)
	fmt.Println(proc.Stack)
}

// Returns the current state of the PDA
func control(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n Api Hit====>control")
	var vars = mux.Vars(r)
	var id = vars["id"]
	proc := cache[id]
	json.NewEncoder(w).Encode(proc.Current_State)

}

func feed(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n Api Hit====>feed")
	var vars = mux.Vars(r)
	var id = vars["id"]
	var position = vars["position"]

	var t Token
	json.NewDecoder(r.Body).Decode(&t)
	var token = t.Token

	proc := cache[id]
	make_transition(&proc, 1)
	pos_int, _ := strconv.Atoi(position)

	token_processed := false
	token_blocked := -1
	var hold_back_flag = false
	if proc.Next_Pos == pos_int {
		fmt.Println("Calling put in client")

		token_processed = putClient(proc, token)
		if token_processed {
			wg.Add(1)
			go func() {
				token_blocked = que_tokens(proc)
			}()
			wg.Wait()
		} else {
			token_blocked = pos_int
		}

	} else if proc.Next_Pos < pos_int {
		var duplicate_token = false

		for _, v := range proc.Hold_Que {
			hold_back_pos_int, _ := strconv.Atoi(v.Hold_Pos)
			if hold_back_pos_int == pos_int {
				duplicate_token = true
			}
		}
		if !duplicate_token {
			var hold_back HoldStruct
			hold_back_flag = true
			hold_back.Hold_Token = token
			hold_back.Hold_Pos = position

			proc.Hold_Que = append(proc.Hold_Que, hold_back)
			sort.Slice(proc.Hold_Que, func(i, j int) bool {
				return proc.Hold_Que[i].Hold_Pos > proc.Hold_Que[j].Hold_Pos
			})
			cache[proc.Id] = proc
			json.NewEncoder(w).Encode("Token kept in Queue")
		} else {
			json.NewEncoder(w).Encode("Duplicate token received")
		}

	} else {
		hold_back_flag = true
		json.NewEncoder(w).Encode("Conflict, token at the same position.")
	}
	if token_blocked == -1 && !hold_back_flag {
		json.NewEncoder(w).Encode("Token consumed")
	} else if !hold_back_flag {
		flag := isAccepted(proc)
		if flag {
			json.NewEncoder(w).Encode("Input successfully Accepted")
		} else {
			json.NewEncoder(w).Encode("Input Rejected by the PDA")
		}
	}
}

func pda_queue(w http.ResponseWriter, r *http.Request) { //returns the queue called on given id
	fmt.Println("\n Api Hit====>queue")
	var vars = mux.Vars(r)
	var id = vars["id"]

	proc := cache[id]

	for j := 0; j < len(proc.Hold_Que)-1; j++ {
		fmt.Println("Queued token :", proc.Hold_Que[j].Hold_Token, " At position :", proc.Hold_Que[j].Hold_Pos)
	}
	json.NewEncoder(w).Encode(proc.Hold_Que)
}

func snapshot(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n API Hit====>Snapshot")
	var vars = mux.Vars(r)
	var id = vars["id"]

	proc := cache[id]
	var snap Snapshot
	snap.Topk = make([]string, 0)
	snap.Current_State = proc.Current_State
	snap.pdaCLOCK = proc.Next_Pos
	snap.Hold_Queue = proc.Hold_Que
	snap.Topk = peekClient(&proc, 5)
	json.NewEncoder(w).Encode(snap)
}

// Performs the last transition to move the Automata to accepting state after the input
// string has been successfully parsed.
func noMore(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n API Hit====>NoMore")
	var vars = mux.Vars(r)
	var id = vars["id"]
	var position = vars["position"]
	pos_int, _ := strconv.Atoi(position)
	proc := cache[id]
	length_of_stack := len(proc.Stack)
	transitions := proc.Transitions
	target_state := ""
	top_of_stack := ""
	var currentStackSymbol = ""
	var top = peekClient(&proc, 1)

	if len(top) >= 1 {
		currentStackSymbol = top[0]
	}
	for j := 0; j < len(transitions); j++ {
		var current_state = transitions[j][0]
		top_of_stack = transitions[j][2]

		if current_state == proc.Current_State && top_of_stack == currentStackSymbol {
			target_state = transitions[j][3]
			break
		}
	}
	if currentStackSymbol == proc.Eos && pos_int == proc.Next_Pos {

		fmt.Println("END OF INPUT REACHED STACK is Empty!")
		proc.Current_State = target_state
		if length_of_stack > 0 {
			pop(&proc)
		}
	}

	cache[id] = proc
}

func stack(w http.ResponseWriter, r *http.Request) { //stack operations
	var vars = mux.Vars(r)
	var id = vars["id"]

	proc := cache[id]
	var l = len(proc.Stack)

	json.NewEncoder(w).Encode(l)
}

func singlecontain(a []string, b string) bool { //To check whether starting state of PDA is valid
	var flag bool
	for i := range a {
		if a[i] == b {
			flag = true
			break
		}
	}
	if flag {
		return true
	} else {
		return false
	}
}

func multicontain(a []string, bc []string) bool { //To check whether accepting states are part of PDA states
	var flag []int
	for i := range bc {
		for j := range a {
			if bc[i] == a[j] {
				flag = append(flag, 1)
				break
			}
		}

	}
	for i := range flag {
		if flag[i] == 1 {
			continue
		} else {
			return false
		}
	}
	return true
}

func triplecontain(a []string, b [][]string) bool { //For validating Transitions with input and stack alphabets
	var flag []int
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			for k := 0; k < len(a); k++ {
				if b[i][j] == a[k] {
					flag = append(flag, 1)
					break
				}
			}
		}
	}
	for i := range flag {
		if flag[i] == 1 {
			continue
		} else {
			return false
		}
	}
	return true
}

func isValid(w http.ResponseWriter, r *http.Request) { //Checks validity of Json specification of PDa with given conditions
	fmt.Println("\n API Hit====>isValid")
	var vars = mux.Vars(r)
	var id = vars["id"]

	p := cache[id]

	if singlecontain(p.States, p.Start_state) && multicontain(p.States, p.Accepting_states) && triplecontain(p.States, p.Transitions) || triplecontain(p.Stack_alphabet, p.Transitions) || triplecontain(p.Input_alphabet, p.Transitions) {
		fmt.Println("Given JSON specifications are not valid for PDA")
	} else {
		fmt.Println("Given JSON specifications are not valid for PDA")
	}
}

func drop(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n API Hit====>drop") //Delete pda processor with given id
	var vars = mux.Vars(r)
	var id = vars["id"]
	_, found := cache[id]
	if found {
		delete(cache, id)
		json.NewEncoder(w).Encode("Pda deleted.")
	} else {
		json.NewEncoder(w).Encode("Pda not found.")

	}
}

func source(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n API Hit====>Source")
	var vars = mux.Vars(r)
	var id = vars["id"]

	pdasource := cache[id]
	json.NewEncoder(w).Encode(pdasource)

}
func close(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Success. No resources to clean.")
}
func clock(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n API Hit====>Clock")
	var vars = mux.Vars(r)
	var id = vars["id"]
	pdasource := cache[id]
	var clk = pdasource.Next_Pos
	fmt.Printf("PDA clock is %d ", clk)
}

func serverRequest() {
	PDARouter := mux.NewRouter().StrictSlash(true)
	PDARouter.HandleFunc("/pdas", AllPdas)
	PDARouter.HandleFunc("/pdas/{id}", create)
	PDARouter.HandleFunc("/pdas/{id}/Open", openServer)
	PDARouter.HandleFunc("/pdas/{id}/reset", reset)
	PDARouter.HandleFunc("/pdas/{id}/isValid", isValid)
	PDARouter.HandleFunc("/pdas/{id}/feed/token/{position}", feed)
	PDARouter.HandleFunc("/pdas/{id}/noMore/{position}", noMore)
	PDARouter.HandleFunc("/pdas/{id}/isAccepted", is_accepted)
	PDARouter.HandleFunc("/pdas/{id}/source", source)
	PDARouter.HandleFunc("/pdas/{id}/clock", clock)
	PDARouter.HandleFunc("/pdas/{id}/peek/{k}", peekServer)
	PDARouter.HandleFunc("/pdas/{id}/stack/len", stack)
	PDARouter.HandleFunc("/pdas/{id}/Control", control)
	PDARouter.HandleFunc("/pdas/{id}/queue", pda_queue)
	PDARouter.HandleFunc("/pdas/{id}/snapshot/{k}", snapshot)
	PDARouter.HandleFunc("/pdas/{id}/close", close)
	PDARouter.HandleFunc("/pdas/{id}/drop", drop)

	log.Fatal(http.ListenAndServe(":8080", PDARouter))
}

func main() {
	fmt.Println("Server started. Listening at port 8080")
	serverRequest()
}
