// Copyright Key Inside Co., Ltd. 2018 All Rights Reserved.

package main

import "fmt"

// ResponsibleError is the interface used to distinguish responsible errors
type ResponsibleError interface {
	IsReponsible() bool
}

// ResponsibleErrorImpl _
type ResponsibleErrorImpl struct{}

// IsReponsible _
func (e ResponsibleErrorImpl) IsReponsible() bool {
	return true
}

// NotExistedContractError _
type NotExistedContractError struct {
	ResponsibleErrorImpl
	id     string
	signer string
}

// Error implements error interface
func (e NotExistedContractError) Error() string {
	return fmt.Sprintf("the contract [%s] for the signer [%s] is not exists", e.id, e.signer)
}
