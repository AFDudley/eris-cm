package chainmaker

import (
	"github.com/eris-ltd/eris-chainmaker/definitions"
	"github.com/eris-ltd/eris-chainmaker/util"

	log "github.com/eris-ltd/eris-chainmaker/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

func MakeMintChain(name string, accounts []*definitions.Account) error {
	genesis := definitions.BlankGenesis()
	genesis.ChainID = name
	for _, account := range accounts {
		log.WithFields(log.Fields{
			"name":    account.Name,
			"address": account.Address,
			"tokens":  account.Tokens,
			"perms":   account.MintPermissions.MintBase.MintPerms,
		}).Debug("Making a Mint Account")

		thisAct := MakeMintAccount(account)
		genesis.Accounts = append(genesis.Accounts, thisAct)

		if account.Validator {
			thisVal := MakeMintValidator(account)
			genesis.Validators = append(genesis.Validators, thisVal)
		}
	}
	for _, account := range accounts {
		if err := util.WritePrivVals(genesis.ChainID, account, len(accounts) == 1); err != nil {
			return err
		}
		if err := util.WriteGenesisFile(genesis.ChainID, genesis, account, len(accounts) == 1); err != nil {
			return err
		}
	}
	return nil
}

func MakeMintAccount(account *definitions.Account) *definitions.MintAccount {
	mintAct := &definitions.MintAccount{}
	mintAct.Address = account.Address
	mintAct.Amount = account.Tokens
	mintAct.Name = account.Name
	mintAct.Permissions = account.MintPermissions
	return mintAct
}

func MakeMintValidator(account *definitions.Account) *definitions.MintValidator {
	mintVal := &definitions.MintValidator{}
	mintVal.Name = account.Name
	mintVal.Amount = account.ToBond
	mintVal.UnbondTo = append(mintVal.UnbondTo, &definitions.MintTxOutput{
		Address: account.Address,
		Amount:  account.ToBond,
	})
	mintVal.PubKey = append(mintVal.PubKey, 1)
	mintVal.PubKey = append(mintVal.PubKey, account.PubKey)
	return mintVal
}
