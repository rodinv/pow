#POW Server & Client
The implementation of the “Word of Wisdom” tcp server and client with DDOS attacks protection
* TCP server is protected from DDOS attacks with the Prof of Work ( https://en.wikipedia.org/wiki/Proof_of_work ) concept
* The Hashcash ( https://en.wikipedia.org/wiki/Hashcash ) is used as the PoW algorithm, because it is the most popular challenge-response algorithm for preventing DDOS attacks
* After Prof Of Work verification, client prints a random quote from “word of wisdom” book quotes  

## Install & Run
```
docker-compose up -d
```
or
* setup env variables
```
POW_HOST (default "localhost")
POW_PORT (default "8081")
POW_BITS (default "24")
```
* run
```
make mod
make server
make client
```
## Additional information
Not implemented (because of PoC):
* tests
* connection timeouts
* persistent hash storage, real resource validation
* connection encryption