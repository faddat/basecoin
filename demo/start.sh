#! /bin/bash
set -eu

# DEFINE COLORS
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
# RED_BOLD='\033[0;1;31m'
NC='\033[0m' # No Color

cd $GOPATH/src/github.com/tendermint/basecoin/demo

function removeQuotes() {
	temp="${1%\"}"
	temp="${temp#\"}"
	echo "$temp"
}

# grab the chain ids
CHAIN_ID1=$(cat ./data/chain1/basecoin/genesis.json | jq .[1])
CHAIN_ID1=$(removeQuotes $CHAIN_ID1)
CHAIN_ID2=$(cat ./data/chain2/basecoin/genesis.json | jq .[1])
CHAIN_ID2=$(removeQuotes $CHAIN_ID2)
echo "CHAIN_ID1: $CHAIN_ID1"
echo "CHAIN_ID2: $CHAIN_ID2"

# make reusable chain flags
CHAIN_FLAGS1="--chain_id $CHAIN_ID1 --from ./data/chain1/basecoin/priv_validator.json"
CHAIN_FLAGS2="--chain_id $CHAIN_ID2 --from ./data/chain2/basecoin/priv_validator.json --node tcp://localhost:36657"

echo ""
printf "    ${RED}STARTING CHAINS CHAIN1 & CHAIN2${NC}\n"
printf "    ${YELLOW}STARTING CHAINS CHAIN1 & CHAIN2${NC}\n"
printf "    ${NC}STARTING CHAINS CHAIN1 & CHAIN2${NC}\n"
echo ""
# start the first node
TMROOT=./data/chain1/tendermint tendermint node &> chain1_tendermint.log &
basecoin start --ibc-plugin --dir ./data/chain1/basecoin &> chain1_basecoin.log &

# start the second node
TMROOT=./data/chain2/tendermint tendermint node --node_laddr tcp://localhost:36656 --rpc_laddr tcp://localhost:36657 --proxy_app tcp://localhost:36658 &> chain2_tendermint.log &
basecoin start --address tcp://localhost:36658 --ibc-plugin --dir ./data/chain2/basecoin &> chain2_basecoin.log &

echo "... waiting for chains to start"
sleep 5
read -p "... press enter"

echo ""
printf "    ${RED}REGISTERING CHAIN1 ON CHAIN2${NC}\n"
printf "    ${YELLOW}REGISTERING CHAIN1 ON CHAIN2${NC}\n"
printf "    ${NC}REGISTERING CHAIN1 ON CHAIN2${NC}\n"
echo ""
# register chain1 on chain2
basecoin ibc --amount 10 $CHAIN_FLAGS2 register --chain_id $CHAIN_ID1 --genesis ./data/chain1/tendermint/genesis.json
read -p "... press enter"

echo ""
printf "    ${RED}CREATING EGRESS PACKET ON CHAIN1${NC}\n"
printf "    ${YELLOW}CREATING EGRESS PACKET ON CHAIN1${NC}\n"
printf "    ${NC}CREATING EGRESS PACKET ON CHAIN1${NC}\n"
echo ""
# create a packet on chain1 destined for chain2
PAYLOAD="DEADBEEF" #TODO
basecoin ibc --amount 10 $CHAIN_FLAGS1 packet create --from $CHAIN_ID1 --to $CHAIN_ID2 --type coin --payload $PAYLOAD --sequence 1
read -p "... press enter"

echo ""
printf "    ${RED}QUERYING FOR PACKET DATA ON CHAIN1${NC}\n"
printf "    ${YELLOW}QUERYING FOR PACKET DATA ON CHAIN1${NC}\n"
printf "    ${NC}QUERYING FOR PACKET DATA ON CHAIN1${NC}\n"
echo ""
# query for the packet data and proof
QUERY_RESULT=$(basecoin query ibc,egress,$CHAIN_ID1,$CHAIN_ID2,1)
HEIGHT=$(echo $QUERY_RESULT | jq .height)
PACKET=$(echo $QUERY_RESULT | jq .value)
PROOF=$(echo $QUERY_RESULT | jq .proof)
PACKET=$(removeQuotes $PACKET)
PROOF=$(removeQuotes $PROOF)
echo ""
echo "QUERY_RESULT: $QUERY_RESULT"
printf "${YELLOW}HEIGHT${NC}: $HEIGHT\n"
printf "${YELLOW}PACKET${NC}: $PACKET\n"
printf "${YELLOW}PROOF${NC}: $PROOF\n"
echo "... waiting for more blocks to be mined"
sleep 5
read -p "... press enter"

echo ""
printf "    ${RED}QUERYING FOR BLOCK DATA ON CHAIN1${NC}\n"
printf "    ${YELLOW}QUERYING FOR BLOCK DATA ON CHAIN1${NC}\n"
printf "    ${NC}QUERYING FOR BLOCK DATA ON CHAIN1${NC}\n"
echo ""
# get the header and commit for the height
HEADER_AND_COMMIT=$(basecoin block $HEIGHT) 
HEADER=$(echo $HEADER_AND_COMMIT | jq .hex.header)
HEADER=$(removeQuotes $HEADER)
COMMIT=$(echo $HEADER_AND_COMMIT | jq .hex.commit)
COMMIT=$(removeQuotes $COMMIT)
echo ""
printf "${YELLOW}HEADER_AND_COMMIT${NC}: $HEADER_AND_COMMIT\n"
printf "${YELLOW}HEADER${NC}: $HEADER\n"
printf "${YELLOW}COMMIT${NC}: $COMMIT\n"
read -p "... press enter"

echo ""
printf "    ${RED}UPDATING STATE OF CHAIN1 ON CHAIN2${NC}\n"
printf "    ${YELLOW}UPDATING STATE OF CHAIN1 ON CHAIN2${NC}\n"
printf "    ${NC}UPDATING STATE OF CHAIN1 ON CHAIN2${NC}\n"
echo ""
# update the state of chain1 on chain2
basecoin ibc --amount 10 $CHAIN_FLAGS2 update --header 0x$HEADER --commit 0x$COMMIT
read -p "... press enter"

echo ""
printf "    ${RED}POSTING PACKET FROM CHAIN1 ON CHAIN2${NC}\n"
printf "    ${YELLOW}POSTING PACKET FROM CHAIN1 ON CHAIN2${NC}\n"
printf "    ${NC}POSTING PACKET FROM CHAIN1 ON CHAIN2${NC}\n"
echo ""
# post the packet from chain1 to chain2
basecoin ibc --amount 10 $CHAIN_FLAGS2 packet post --from $CHAIN_ID1 --height $((HEIGHT + 1)) --packet 0x$PACKET --proof 0x$PROOF
read -p "... press enter"

echo ""
printf "    ${RED}CHECKING IF THE PACKET IS PRESENT ON CHAIN2${NC}\n"
printf "    ${YELLOW}CHECKING IF THE PACKET IS PRESENT ON CHAIN2${NC}\n"
printf "    ${NC}CHECKING IF THE PACKET IS PRESENT ON CHAIN2${NC}\n"
echo ""
# query for the packet on chain2 !
basecoin query --node tcp://localhost:36657 ibc,ingress,test_chain_2,test_chain_1,1
