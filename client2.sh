printf "\n################ Create PDA with id 121 ################" 
curl -X PUT -H "Content-Type: application/json" -d '{
    "name": "0n1n",
    "states": ["q1", "q2", "q3", "q4"],
    "input_alphabet": [ "0", "1" ],
    "stack_alphabet" : [ "0", "1" ],
    "accepting_states": ["q1", "q4"],
    "start_state": "q1",
    "transitions": [
        ["q1", "null", "null", "q2", "$"],
        ["q2", "0", "null", "q2", "0"],
        ["q2", "0", "0", "q2", "0"],
        ["q2", "1", "0", "q3", "null"],
        ["q3", "1", "0", "q3", "null"],
        ["q3", "null", "$", "q4", "null"]
    ],
    "eos": "$"
}' http://localhost:8080/pdas/121

printf "\n PDA CREATED"

curl -X GET http://localhost:8080/pdas/121/OPEN

printf "\n PDA OPENED"

printf "\n################ FEED tokens ################\n" 

curl -X PUT -H "Content-Type: application/json" -d '{"token": "0"}' http://localhost:8080/pdas/121/tokens/0

printf "\n"

curl -X PUT -H "Content-Type: application/json" -d '{"token": "0"}' http://localhost:8080/pdas/121/tokens/1


printf "\n################## Current state of the PDA #######################\n" 

curl -X GET http://localhost:8080/pdas/121/control

printf "\n###################### Queue ##################\n" 

curl -X GET http://localhost:8080/pdas/121/queue

printf "\n<<<<<<<<<<<<<<<<<< Peek of 3 >>>>>>>>>>>>>>>\n"

curl --X GET http://localhost:8080/pdas/121/peek/3

printf "\n######################### Continue feeding other tokens #################\n" 

curl -X PUT -H "Content-Type: application/json" -d '{"token": "1"}' http://localhost:8080/pdas/101/tokens/2

printf "\n"

curl -X PUT -H "Content-Type: application/json" -d '{"token": "0"}' http://localhost:8080/pdas/101/tokens/3

printf "\n###################### Queue ##################\n" 

curl -X GET http://localhost:8080/pdas/121/queue

printf "\n###################### clock ##################\n" 

curl -X GET http://localhost:8080/pdas/121/clock


printf "\n######################### Snapshot ####################\n" 

curl -X GET http://localhost:8080/pdas/121/snapshot/3

printf "\n##################### noMore ####################\n" 

curl http://localhost:8080/pdas/121/noMore/6

printf "\n#################### Call API isAccepted() ###################\n" 

curl http://localhost:8080/pdas/121/isAccepted


printf "\n################ Reset ####################\n" 

curl -X GET http://localhost:8080/pdas/121/reset

printf "\n################## Snapshot ##################\n" 

curl -X GET http://localhost:8080/pdas/101/snapshot/3

printf "\n###################### All Pdas on server############\n" 

curl -X GET http://localhost:8080/pdas

printf "\n################## Drop pda with id 101 ################\n" 

curl -X GET http://localhost:8080/pdas/101/drop

printf "\n#################### Show all Pdas after dropping  ###################\n" 

curl -X GET http://localhost:8080/pdas

