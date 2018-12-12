/**
 * Copyright (c) 2018, 2019 National Digital ID COMPANY LIMITED
 *
 * This file is part of NDID software.
 *
 * NDID is the free software: you can redistribute it and/or modify it under
 * the terms of the Affero GNU General Public License as published by the
 * Free Software Foundation, either version 3 of the License, or any later
 * version.
 *
 * NDID is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 * See the Affero GNU General Public License for more details.
 *
 * You should have received a copy of the Affero GNU General Public License
 * along with the NDID source code. If not, see https://www.gnu.org/licenses/agpl.txt.
 *
 * Please contact info@ndid.co.th for any further questions
 *
 */

package did

import (
	"encoding/base64"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/ndidplatform/smart-contract/abci/code"
	"github.com/ndidplatform/smart-contract/abci/utils"
	"github.com/ndidplatform/smart-contract/abci/version"
	"github.com/sirupsen/logrus"
	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/abci/types"

	protoTm "github.com/ndidplatform/smart-contract/protos/tendermint"
)

var (
	stateKey        = []byte("stateKey")
	kvPairPrefixKey = []byte("kvPairKey:")
)

type State struct {
	db *iavl.MutableTree
}

func prefixKey(key []byte) []byte {
	return append(kvPairPrefixKey, key...)
}

var _ types.Application = (*DIDApplication)(nil)

type DIDApplication struct {
	types.BaseApplication
	state              State
	checkTxTempState   map[string][]byte
	ValUpdates         []types.ValidatorUpdate
	logger             *logrus.Entry
	Version            string
	AppProtocolVersion uint64
	CurrentBlock       int64
	CurrentChain       string
}

func NewDIDApplication(logger *logrus.Entry, tree *iavl.MutableTree) *DIDApplication {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("%s", identifyPanic())
			panic(r)
		}
	}()
	var state State
	state.db = tree

	ABCIVersion := version.Version
	ABCIProtocolVersion := version.AppProtocolVersion
	logger.Infof("Start ABCI version: %s", ABCIVersion)
	return &DIDApplication{
		state:              state,
		checkTxTempState:   make(map[string][]byte),
		logger:             logger,
		Version:            ABCIVersion,
		AppProtocolVersion: ABCIProtocolVersion,
	}
}

func (app *DIDApplication) SetStateDB(key, value []byte) {
	app.state.db.Set(prefixKey(key), value)
}

func (app *DIDApplication) SetStateDBWithOutPrefix(key, value []byte) {
	app.state.db.Set(key, value)
}

func (app *DIDApplication) DeleteStateDB(key []byte) {
	app.state.db.Remove(prefixKey(key))
}

func (app *DIDApplication) Info(req types.RequestInfo) (resInfo types.ResponseInfo) {
	var res types.ResponseInfo
	res.Version = app.Version
	res.LastBlockHeight = app.state.db.Version()
	res.LastBlockAppHash = app.state.db.Hash()
	res.AppVersion = app.AppProtocolVersion
	app.CurrentBlock = app.state.db.Version()
	return res
}

// Save the validators in the merkle tree
func (app *DIDApplication) InitChain(req types.RequestInitChain) types.ResponseInitChain {
	for _, v := range req.Validators {
		r := app.updateValidator(v)
		if r.IsErr() {
			app.logger.Error("Error updating validators", "r", r)
		}
	}
	return types.ResponseInitChain{}
}

// Track the block hash and header information
func (app *DIDApplication) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	app.logger.Infof("BeginBlock: %d, Chain ID: %s", req.Header.Height, req.Header.ChainID)
	var dbDir = getEnv("ABCI_DB_DIR_PATH", "./DID")
	utils.WriteEventLogBeginBlock(dbDir, time.Now(), "start_BeginBlock", req.Header.Height, req.Header.NumTxs)
	app.CurrentBlock = req.Header.Height
	app.CurrentChain = req.Header.ChainID
	// reset valset changes
	app.ValUpdates = make([]types.ValidatorUpdate, 0)
	utils.WriteEventLogBeginBlock(dbDir, time.Now(), "stop_BeginBlock", req.Header.Height, req.Header.NumTxs)
	return types.ResponseBeginBlock{}
}

// Update the validator set
func (app *DIDApplication) EndBlock(req types.RequestEndBlock) types.ResponseEndBlock {
	var dbDir = getEnv("ABCI_DB_DIR_PATH", "./DID")
	utils.WriteEventLog(dbDir, time.Now(), "start_EndBlock")
	app.logger.Infof("EndBlock: %d", req.Height)
	utils.WriteEventLog(dbDir, time.Now(), "start_EndBlock")
	return types.ResponseEndBlock{ValidatorUpdates: app.ValUpdates}
}

func (app *DIDApplication) DeliverTx(tx []byte) (res types.ResponseDeliverTx) {
	// Recover when panic
	defer func() {
		if r := recover(); r != nil {
			app.logger.Errorf("Recovered in %s, %s", r, identifyPanic())
			res = app.ReturnDeliverTxLog(code.UnknownError, "Unknown error", "")
		}
	}()

	var txObj protoTm.Tx
	err := proto.Unmarshal(tx, &txObj)
	if err != nil {
		app.logger.Error(err.Error())
	}

	method := txObj.Method
	param := txObj.Params
	nonce := txObj.Nonce
	signature := txObj.Signature
	nodeID := txObj.NodeId

	var dbDir = getEnv("ABCI_DB_DIR_PATH", "./DID")
	nonceBase64 := base64.StdEncoding.EncodeToString(nonce)
	utils.WriteEventLogTx(dbDir, time.Now(), "start_DeliverTx", method, nonceBase64)

	// ---- Check duplicate nonce ----
	nonceDup := app.isDuplicateNonce(nonce)
	if nonceDup {
		utils.WriteEventLogTx(dbDir, time.Now(), "end_DeliverTx", method, nonceBase64)
		return app.ReturnDeliverTxLog(code.DuplicateNonce, "Duplicate nonce", "")
	}

	app.logger.Infof("DeliverTx: %s, NodeID: %s", method, nodeID)

	if method != "" {
		result := app.DeliverTxRouter(method, param, nonce, signature, nodeID)
		app.logger.Infof(`DeliverTx response: {"code":%d,"log":"%s","tags":[{"key":"%s","value":"%s"}]}`, result.Code, result.Log, string(result.Tags[0].Key), string(result.Tags[0].Value))
		utils.WriteEventLogTx(dbDir, time.Now(), "end_DeliverTx", method, nonceBase64)
		return result
	}
	utils.WriteEventLogTx(dbDir, time.Now(), "end_DeliverTx", method, nonceBase64)
	return app.ReturnDeliverTxLog(code.MethodCanNotBeEmpty, "method can not be empty", "")
}

func (app *DIDApplication) CheckTx(tx []byte) (res types.ResponseCheckTx) {
	// Recover when panic
	defer func() {
		if r := recover(); r != nil {
			app.logger.Errorf("Recovered in %s, %s", r, identifyPanic())
			res = ReturnCheckTx(code.UnknownError, "Unknown error")
		}
	}()

	var txObj protoTm.Tx
	err := proto.Unmarshal(tx, &txObj)
	if err != nil {
		app.logger.Error(err.Error())
	}

	method := txObj.Method
	param := txObj.Params
	nonce := txObj.Nonce
	signature := txObj.Signature
	nodeID := txObj.NodeId

	var dbDir = getEnv("ABCI_DB_DIR_PATH", "./DID")
	nonceBase64 := base64.StdEncoding.EncodeToString(nonce)
	utils.WriteEventLogTx(dbDir, time.Now(), "start_CheckTx", method, nonceBase64)

	// ---- Check duplicate nonce ----
	nonceDup := app.isDuplicateNonce(nonce)
	if nonceDup {
		res.Code = code.DuplicateNonce
		res.Log = "Duplicate nonce"
		utils.WriteEventLogTx(dbDir, time.Now(), "end_CheckTx", method, nonceBase64)
		return res
	}

	// Check duplicate nonce in checkTx stateDB
	_, exist := app.checkTxTempState[nonceBase64]
	if !exist {
		app.checkTxTempState[nonceBase64] = []byte("1")
	} else {
		res.Code = code.DuplicateNonce
		res.Log = "Duplicate nonce"
		utils.WriteEventLogTx(dbDir, time.Now(), "end_CheckTx", method, nonceBase64)
		return res
	}

	app.logger.Infof("CheckTx: %s, NodeID: %s", method, nodeID)

	if method != "" && param != "" && nonce != nil && signature != nil && nodeID != "" {
		// Check has function in system
		if IsMethod[method] {
			result := app.CheckTxRouter(method, param, nonce, signature, nodeID)
			utils.WriteEventLogTx(dbDir, time.Now(), "end_CheckTx", method, nonceBase64)
			return result
		}
		res.Code = code.UnknownMethod
		res.Log = "Unknown method name"
		utils.WriteEventLogTx(dbDir, time.Now(), "end_CheckTx", method, nonceBase64)
		return res
	}
	res.Code = code.InvalidTransactionFormat
	res.Log = "Invalid transaction format"
	utils.WriteEventLogTx(dbDir, time.Now(), "end_CheckTx", method, nonceBase64)
	return res
}

func (app *DIDApplication) Commit() types.ResponseCommit {
	var dbDir = getEnv("ABCI_DB_DIR_PATH", "./DID")
	utils.WriteEventLog(dbDir, time.Now(), "start_Commit")
	app.logger.Infof("Commit")
	app.state.db.SaveVersion()
	app.checkTxTempState = make(map[string][]byte)
	utils.WriteEventLog(dbDir, time.Now(), "end_Commit")
	return types.ResponseCommit{Data: app.state.db.Hash()}
}

func (app *DIDApplication) Query(reqQuery types.RequestQuery) (res types.ResponseQuery) {

	// Recover when panic
	defer func() {
		if r := recover(); r != nil {
			app.logger.Errorf("Recovered in %s, %s", r, identifyPanic())
			res = app.ReturnQuery(nil, "Unknown error", app.state.db.Version())
		}
	}()

	var query protoTm.Query
	err := proto.Unmarshal(reqQuery.Data, &query)
	if err != nil {
		app.logger.Error(err.Error())
	}

	method := query.Method
	param := query.Params

	var dbDir = getEnv("ABCI_DB_DIR_PATH", "./DID")
	utils.WriteEventLogQuery(dbDir, time.Now(), "start_Query", method)

	app.logger.Infof("Query: %s", method)

	height := reqQuery.Height
	if height == 0 {
		height = app.state.db.Version()
	}

	if method != "" {
		utils.WriteEventLogQuery(dbDir, time.Now(), "end_Query", method)
		return app.QueryRouter(method, param, height)
	}
	utils.WriteEventLogQuery(dbDir, time.Now(), "end_Query", method)
	return app.ReturnQuery(nil, "method can't empty", app.state.db.Version())
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultValue
	}
	return value
}

func identifyPanic() string {
	var name, file string
	var line int
	var pc [16]uintptr

	n := runtime.Callers(3, pc[:])
	for _, pc := range pc[:n] {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		file, line = fn.FileLine(pc)
		name = fn.Name()
		if !strings.HasPrefix(name, "runtime.") {
			break
		}
	}

	switch {
	case name != "":
		return fmt.Sprintf("%v:%v", name, line)
	case file != "":
		return fmt.Sprintf("%v:%v", file, line)
	}

	return fmt.Sprintf("pc:%x", pc)
}
