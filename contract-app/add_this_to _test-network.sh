##Add this to test-network folder and run it in there
export PATH=$PATH:/usr/local/go/bin

./network.sh up createChannel -ca -s couchdb

./network.sh deployCC -ccn "KYRAL" -ccp ../../shf_continuity/contract/ -ccl go -ccep "OR('Org1MSP.peer','Org2MSP.peer')"

