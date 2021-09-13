# Push Down Automata-with-REST-API


  •	pda_client.go is a client file which has client methods of the PDA. \
  •	pda_server.go has server methods and definitions of all api endpoints. \
  •	Please refer pda_def.go for understanding of structure definitions required for processing of PDA. \
  •	Refer the comments for function operations.
## Firstly, open command prompt and run the following command:
  ``` javascript
  go get -u github.com/gorilla/mux
  ```
### This will install required router library which is required to implement REST API with client server.
  •	Now open any bash executor. \
  •	Start one server_script.sh on a separate bash executor with the command: 
``` javascript
bash server_script.sh
```
### Note: Keep server script running ALL THE TIME! 
  •	Run the following commands for individual scripts
``` javascript
bash client1.sh
```
``` javascript
bash clientt2.sh
```
  •	Standard Output is been recorded in respective files.\
  •	Pda_server.go implements a pda router that runs on port 8080 to provide rest api for all the given tasks.
