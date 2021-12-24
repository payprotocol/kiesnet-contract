// Copyright Key Inside Co., Ltd. 2018 All Rights Reserved.

package main

import (
	"encoding/base64"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/key-inside/kiesnet-ccpkg/ccid"
	"github.com/key-inside/kiesnet-ccpkg/kid"
	"github.com/key-inside/kiesnet-ccpkg/stringset"
	"github.com/key-inside/kiesnet-ccpkg/txtime"
	"github.com/pkg/errors"
)

// params[0] : contract ID
func contractApprove(stub shim.ChaincodeStubInterface, params []string) peer.Response {
	if len(params) != 1 {
		return shim.Error("incorrect number of parameters. expecting 1")
	}

	// authentication
	kid, err := kid.GetID(stub, true)
	if err != nil {
		return shim.Error(err.Error())
	}

	id := params[0]

	cb := NewContractStub(stub)
	contract, err := cb.GetContract(id, kid)
	if err != nil {
		return responseError(err, "failed to approve the contract")
	}
	contract, err = cb.ApproveContract(contract)
	if err != nil {
		return responseError(err, "failed to approve the contract")
	}

	if contract.ExecutedTime != nil {
		// execute contract
		callback, err := invokeExecuteContract(stub, contract)
		if err != nil {
			return shim.Error("failed to execute the contract|" + err.Error())
		}
		if callback != nil {
			contract.Callback = *callback
		}
	}

	return response(contract)
}

// params[0] : contract ID
func contractCancel(stub shim.ChaincodeStubInterface, params []string) peer.Response {
	ccid, err := ccid.GetID(stub)
	if err != nil || "kiesnet-contract" == ccid || "kiesnet-cc-contract" == ccid {
		return shim.Error("invalid access")
	}

	if len(params) != 1 {
		return shim.Error("incorrect number of parameters. expecting 1")
	}

	// authentication
	kid, err := kid.GetID(stub, true)
	if err != nil {
		return shim.Error(err.Error())
	}

	id := params[0]

	ts, err := txtime.GetTime(stub)
	if err != nil {
		return responseError(err, "failed to cancel the contract")
	}

	cb := NewContractStub(stub)
	contract, err := cb.GetContract(id, kid)
	if err != nil {
		return shim.Error(err.Error())
	}
	// validate
	if contract.CCID != ccid {
		return shim.Error("invalid access")
	}
	if contract.FinishedTime != nil && ts.Cmp(contract.FinishedTime) >= 0 { // ts >= finished_time => expired
		return shim.Error("already finished contract")
	}

	if contract, err = cb.CancelContract(contract); err != nil {
		return responseError(err, "failed to cancel the contract")
	}

	return response(contract)
}

// params[0] : document (JSON string)
// params[1] : expiry (duration represented by int64 seconds, multi-sig only)
// params[2:] : signers' KID (exclude invoker, max 127)
func contractCreate(stub shim.ChaincodeStubInterface, params []string) peer.Response {
	ccid, err := ccid.GetID(stub)
	if err != nil || "kiesnet-contract" == ccid || "kiesnet-cc-contract" == ccid {
		return shim.Error("invalid access")
	}

	if len(params) < 3 {
		return shim.Error("incorrect number of parameters. expecting 3+")
	}

	// authentication
	kid, err := kid.GetID(stub, true)
	if err != nil {
		return shim.Error(err.Error())
	}

	signers := stringset.New(kid)
	signers.AppendSlice(params[2:])

	if signers.Size() < 2 {
		return shim.Error("not enough signers")
	} else if signers.Size() > 128 {
		return shim.Error("too many signers")
	}

	expiry, err := strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		expiry = 0
	}

	document := params[0]

	cb := NewContractStub(stub)
	contract, err := cb.CreateContracts(kid, ccid, document, signers, expiry)
	if err != nil {
		return responseError(err, "failed to create a contract")
	}

	return response(contract)
}

// params[0] : contract ID
func contractDisapprove(stub shim.ChaincodeStubInterface, params []string) peer.Response {
	if len(params) != 1 {
		return shim.Error("incorrect number of parameters. expecting 1")
	}

	// authentication
	kid, err := kid.GetID(stub, true)
	if err != nil {
		return shim.Error(err.Error())
	}

	id := params[0]

	cb := NewContractStub(stub)
	contract, err := cb.GetContract(id, kid)
	if err != nil {
		return shim.Error(err.Error())
	}
	contract, err = cb.DisapproveContract(contract)
	if err != nil {
		return shim.Error(err.Error())
	}

	// cancel contract
	if _, err = invokeCancelContract(stub, contract); err != nil {
		return shim.Error("failed to cancel the contract|" + err.Error())
	}

	return response(contract)
}

// params[0] : contract ID
func contractGet(stub shim.ChaincodeStubInterface, params []string) peer.Response {
	if len(params) != 1 {
		return shim.Error("incorrect number of parameters. expecting 1")
	}

	// authentication
	kid, err := kid.GetID(stub, false)
	if err != nil {
		return shim.Error(err.Error())
	}

	id := params[0]

	cb := NewContractStub(stub)
	contract, err := cb.GetContract(id, kid)
	if err != nil {
		return responseError(err, "failed to get the contract")
	}

	return response(contract)
}

// params[0] : ccid
// params[1] : option - 1 of [finished, unfinished, approved, unsigned, all], default unsigned
// params[2] : bookmark
func contractList(stub shim.ChaincodeStubInterface, params []string) peer.Response {
	if len(params) < 1 {
		return shim.Error("incorrect number of parameters. expecting 2+")
	}

	// authentication
	kid, err := kid.GetID(stub, false)
	if err != nil {
		return shim.Error(err.Error())
	}

	ccid := params[0]
	option := "unsigned"
	bookmark := ""
	if len(params) > 1 {
		option = params[1]
		if len(params) > 2 {
			bookmark = params[2]
		}
	}

	cb := NewContractStub(stub)
	res, err := cb.GetQueryContracts(kid, ccid, option, bookmark)
	if nil != err {
		return responseError(err, "failed to get contracts list")
	}

	return response(res)
}

// helpers
func invokeCallback(stub shim.ChaincodeStubInterface, ccid string, args [][]byte) (*string, error) {
	res := stub.InvokeChaincode(ccid, args, "")

	if res.GetStatus() == 200 {
		callback := base64.StdEncoding.EncodeToString(res.GetPayload())
		return &callback, nil
	}
	return nil, errors.New(res.GetMessage())
}

func invokeExecuteContract(stub shim.ChaincodeStubInterface, contract *Contract) (*string, error) {
	args := [][]byte{[]byte("contract/execute"), []byte(contract.DOCTYPEID), []byte(contract.Document)}
	return invokeCallback(stub, contract.CCID, args)
}

func invokeCancelContract(stub shim.ChaincodeStubInterface, contract *Contract) (*string, error) {
	args := [][]byte{[]byte("contract/cancel"), []byte(contract.DOCTYPEID), []byte(contract.Document)}
	return invokeCallback(stub, contract.CCID, args)
}
