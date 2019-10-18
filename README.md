# Kiesnet Contract Chaincode

## Requirement
- kiesnet-id chaincode (devmode: kiesnet-cc-id)

#

## API

method __`func`__ [arg1, _arg2_, ... ] {trs1, _trs2_, ... }
- method : __query__ or __invoke__
- func : function name
- [arg] : mandatory argument
- [_arg_] : optional argument
- {trs} : mandatory transient
- {_trs_} : optional transient

#

> invoke __`approve`__ [contract_id] {_"kiesnet-id/pin"_}
- Approve the contract
- If all signers have approved the contract, it invokes 'contract/execute' callback.

> invoke __`cancel`__ [contract_id] {_"kiesnet-id/pin"_}
- Cancel the contract
- It invokes 'contract/cancel' callback.

> invoke __`create`__ [document, expiry, signers...] {_"kiesnet-id/pin"_}
- Create a contract
- [document] : contract document JSON string, it will be passed to callbacks
- [expiry] : duration(seconds) represented by int64, if it's less than 10 minutes, default expiry will be set (15 days)
- [signers...] : KIDs of signers (exclude invoker, max 127)

> invoke __`disapprove`__ [contract_id] {_"kiesnet-id/pin"_}
- Disapprove the contract
- It invokes 'contract/cancel' callback.

> query __`get`__ [contract_id]
- Get the contract

> query __`list`__ [ccid, _option_, _bookmark_]
- Get contracts list of the invoker
- [ccid] : chaincode ID created a contract
- [option] : 1 of [finished, unfinished, approved, unsigned, all], default unsigned

> query __`ver`__
- Get version

#

## Callbacks
Invoker chaincodes must implement callbacks.

#

> invoke __`contract/execute`__ [contract_id, document] {_"kiesnet-id/pin"_}
- Execute the contract

> invoke __`contract/cancel`__ [contract_id, document] {_"kiesnet-id/pin"_}
- Cancel the contract
