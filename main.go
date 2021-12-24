// Copyright Key Inside Co., Ltd. 2018 All Rights Reserved.

package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("kiesnet-contract")

// Chaincode _
type Chaincode struct {
}

// Init implements shim.Chaincode interface.
func (cc *Chaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke implements shim.Chaincode interface.
func (cc *Chaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, params := stub.GetFunctionAndParameters()
	if txFn := routes[fn]; txFn != nil {
		return txFn(stub, params)
	}
	return shim.Error("unknown function: [" + fn + "]")
}

// TxFunc _
type TxFunc func(shim.ChaincodeStubInterface, []string) peer.Response

// routes is the map of invoke functions
var routes = map[string]TxFunc{
	"approve":    contractApprove,
	"cancel":     contractCancel,
	"create":     contractCreate,
	"disapprove": contractDisapprove,
	"get":        contractGet,
	"list":       contractList,
	"ver":        ver,
}

func ver(stub shim.ChaincodeStubInterface, params []string) peer.Response {
	return shim.Success([]byte("Kiesnet Contract v1.3 created by Key Inside Co., Ltd."))
}

func response(payload Payload) peer.Response {
	data, err := payload.MarshalPayload()
	if err != nil {
		logger.Debug(err.Error())
		return shim.Error("failed to marshal payload")
	}
	return shim.Success(data)
}

// If 'err' is ResponsibleError, it will add err's message to the 'msg'.
func responseError(err error, msg string) peer.Response {
	if nil != err {
		logger.Debug(err.Error())
		if _, ok := err.(ResponsibleError); ok {
			if len(msg) > 0 {
				msg = msg + "|" + err.Error()
			} else {
				msg = err.Error()
			}
		}
	}
	return shim.Error(msg)
}

func main() {
	if err := shim.Start(new(Chaincode)); err != nil {
		logger.Criticalf("failed to start chaincode|%s", err)
	}
}
