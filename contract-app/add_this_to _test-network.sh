##Add this to test-network folder and run it in there
###(assume you have go in the bin directory) export PATH=$PATH:/usr/local/go/bin

./network.sh up createChannel -ca -s couchdb

./network.sh deployCC -ccn "KYRAL" -ccp ../../KYRAL-Health/KyralTEST/ -ccl go -ccep "OR('Org1MSP.peer','Org2MSP.peer')"

