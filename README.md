# lambda-libs
Lambda-libs is a dependency library for lambda projects, mostly from [Ethereum](https://github.com/ethereum/go-ethereum) and used by [lambda-vm](https://github.com/LambdaIM/lambda-vm)

## abi
ABI(Application Binary Interface) is the standard way to interact with contracts in the Lambda-vm.
Parsing of multidimensional arrays is not supported now.

## common
Common contains some base libraries (math, hexutil, etc.).

## crypto
Crypto is the crypto library of Ethereum, currently used mainly by lambda-vm

## ethdb
Physical storage used by Lambda-vm (leveldb)

## rlp
RLP (Recursive Length Prefix) is to encode arbitrarily nested arrays of binary data used to serialize objects in Lambda-vm.

## state
State contains StateDB and stateobject.  
StateDB stores stateObject, and a stateObject represents a Lambda account. 
Stateobject contains status information such as account address, balance, nonce, contract code hash, and so on. 

## trie
Trie(Merkle Patricia Trie) are used to one global state trie, and it updates over time.
