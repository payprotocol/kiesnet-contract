// Copyright Key Inside Co., Ltd. 2018 All Rights Reserved.

package main

import (
	"fmt"

	"github.com/key-inside/kiesnet-ccpkg/txtime"
)

// QueryContractsByID _
const QueryContractsByID = `{
	"selector": {
		"@contract": "%s"
	},
	"use_index": ["contract", "id"]
}`

// CreateQueryContractsByID _
func CreateQueryContractsByID(id string) string {
	return fmt.Sprintf(QueryContractsByID, id)
}

// QueryContractsBySigner _
const QueryContractsBySigner = `{
	"selector": {
		"@contract": {
			"$exists": true
		},
		"sign.signer": "%s",
		"ccid": "%s"
	},
	"sort": [{"sign.signer": "desc"}, {"ccid": "desc"}, {"created_time": "desc"}],
	"use_index": ["contract", "created-time"]
}`

// CreateQueryContractsBySigner _
func CreateQueryContractsBySigner(kid, ccid string) string {
	return fmt.Sprintf(QueryContractsBySigner, kid, ccid)
}

// QueryFinishedContractsBySigner _
const QueryFinishedContractsBySigner = `{
	"selector": {
		"@contract": {
			"$exists": true
		},
		"sign.signer": "%s",
		"ccid": "%s",
		"finished_time": {
			"$lte": "%s"
		}
	},
	"sort": [{"sign.signer": "desc"}, {"ccid": "desc"}, {"finished_time": "desc"}],
	"use_index": ["contract", "finished-time"]
}`

// CreateQueryFinishedContractsBySigner _
func CreateQueryFinishedContractsBySigner(kid, ccid string, ts *txtime.Time) string {
	return fmt.Sprintf(QueryFinishedContractsBySigner, kid, ccid, ts.String())
}

// QueryUnfinishedContractsBySigner _
const QueryUnfinishedContractsBySigner = `{
	"selector": {
		"@contract": {
			"$exists": true
		},
		"sign.signer": "%s",
		"ccid": "%s",
		"finished_time": {
			"$gt": "%s"
		}
	},
	"sort": ["sign.signer", "ccid", "finished_time"],
	"use_index": ["contract", "finished-time"]
}`

// CreateQueryUnfinishedContractsBySigner _
func CreateQueryUnfinishedContractsBySigner(kid, ccid string, ts *txtime.Time) string {
	return fmt.Sprintf(QueryUnfinishedContractsBySigner, kid, ccid, ts.String())
}

// QueryApprovedContractsBySigner - unfinished, approved
const QueryApprovedContractsBySigner = `{
	"selector": {
		"$and": [
			{
				"@contract": {
					"$exists": true
				}
			},
			{
				"sign.approved_time": {
					"$exists": true
				}
			},
			{
				"executed_time": {
					"$exists": false
				}
			},
			{
				"canceled_time": {
					"$exists": false
				}
			}
		],
		"sign.signer": "%s",
		"ccid": "%s",
		"expiry_time": {
			"$gt": "%s"
		}
	},
	"sort": ["sign.signer", "ccid", "expiry_time"],
	"use_index": ["contract", "approved-expiry-time"]
}`

// CreateQueryApprovedContractsBySigner _
func CreateQueryApprovedContractsBySigner(kid, ccid string, ts *txtime.Time) string {
	return fmt.Sprintf(QueryApprovedContractsBySigner, kid, ccid, ts.String())
}

// QueryUnsignedContractsBySigner - unfinished, unsigned
const QueryUnsignedContractsBySigner = `{
	"selector": {
		"$and": [
			{
				"@contract": {
					"$exists": true
				}
			},
			{
				"sign.approved_time": {
					"$exists": false
				}
			},
			{
				"sign.disapproved_time": {
					"$exists": false
				}
			},
			{
				"executed_time": {
					"$exists": false
				}
			},
			{
				"canceled_time": {
					"$exists": false
				}
			}
		],
		"sign.signer": "%s",
		"ccid": "%s",
		"expiry_time": {
			"$gt": "%s"
		}
	},
	"sort": ["sign.signer", "ccid", "expiry_time"],
	"use_index": ["contract", "unsigned-expiry-time"]
}`

// CreateQueryUnsignedContractsBySigner _
func CreateQueryUnsignedContractsBySigner(kid, ccid string, ts *txtime.Time) string {
	return fmt.Sprintf(QueryUnsignedContractsBySigner, kid, ccid, ts.String())
}
