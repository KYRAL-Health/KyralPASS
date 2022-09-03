# How to connect to the Kyral Beta Network
If you have already setup your own hyperledger fabric peer you can connect by submitting a contact form on our website stating you want to connect.  
Please readily have the following needed to be able to join the channel:
* anchor peer dns name/ip address and port
* zip of your msp dir containing only public certificates and should have the following folders:
    * cacert
    * config.yaml with NodeOUs. Required OUs are client, peer, admin and orderer.
    * tlscacerts

## Example config.yaml
```yaml
NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/ca-beta-8054-ca-beta.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/ca-beta-8054-ca-beta.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/ca-beta-8054-ca-beta.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/ca-beta-8054-ca-beta.pem
    OrganizationalUnitIdentifier: orderer
```

Upon receiving your form to join the network we will send hostname and ports of our orderers and genesis block so you can join the channel.

After joining the channel all you have to do is run the latest version of the chaincode and you can interact with other members in the network.
