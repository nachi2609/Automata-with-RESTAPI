printf "\n<<<<<<<<<<<<< Create PDA with id 100 >>>>>>>>>>>>\n" 
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
}' http://localhost:8080/pdas/100

printf "\n<<<<<<<<<<<<<<<ALL PDAS on Server>>>>>>>>>>\n" 
curl -X GET http://localhost:8080/pdas


printf "\n<<<<<<<<<<<<<<<< OPEN >>>>>>>>>>>>>>\n" 
curl -x GET http://localhost:8080/pdas/100/Open


printf "\n<<<<<<<<<<<<< feed tokens >>>>>>>>>>>>>>>>>\n" 

curl -X PUT -H "Content-Type: application/json" -d '{"token": "0"}' http://localhost:8080/pdas/100/tokens/1


printf "\n"

curl -X PUT -H "Content-Type: application/json" -d '{"token": "0"}' http://localhost:8080/pdas/100/feed/tokens/2


printf "\n<<<<<<<<<<<<<<<<<<<<<<<<< Current state of the PDA >>>>>>>>>>>>\n" 

curl -X GET http://localhost:8080/pdas/100/Control


printf "\n<<<<<<<<<<<<<<<<<< Continue feeding other tokens>>>>>>>>>>>>>>>>>>>>>>>\n" 

curl -X PUT -H "Content-Type: application/json" -d '{"token": "1"}' http://localhost:8080/pdas/100/feed/token/3

printf "\n"

curl -X PUT -H "Content-Type: application/json" -d '{"token": "1"}' http://localhost:8080/pdas/100/feed/token/4

printf "\n"

curl -X PUT -H "Content-Type: application/json" -d '{"token": "1"}' http://localhost:8080/pdas/100/feed/token/5

printf "\n<<<<<<<<<<<<<<<< Queue >>>>>>>>>>>>>>>>>>\n" 

curl -X GET http://localhost:8080/pdas/100/queue


printf "\n <<<<<<<<<<<<<<<<<< Put token at position 0 >>>>>>>>>>>>>>>>>>\n"

curl -X PUT -H "Content-Type: application/json" -d '{"token": "0"}' http://localhost:8080/pdas/100/feed/token/0

printf "\n<<<<<<<<<<<<<<<<<<<< Snapshot >>>>>>>>>>>>>>>>>>>>>\n" 

curl -X GET http://localhost:8080/pdas/100/snapshot/3


printf "\n<<<<<<<<<<<<<<<< API noMore >>>>>>>>>>>>>>>>>\n" 

curl http://localhost:8080/pdas/100/noMore/6

printf "\n<<<<<<<<<<<<<<<< API isAccepted() >>>>>>>>>>>>>>>>>>>>>\n" 

curl http://localhost:8080/pdas/100/isAccepted
