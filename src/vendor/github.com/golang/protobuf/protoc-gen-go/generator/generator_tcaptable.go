package generator

import (
	"fmt"
	_ "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	_ "github.com/golang/protobuf/protoc-gen-go/plugin"
	"log"
	"os"
	"path/filepath"
	easyapi "qqgame/baselib/file"
	"reflect"
	"strconv"
	"strings"
)

type NQGOptions struct {
	TableName      string
	CamelTableName string
	BKey           int

	BCustomerType int
	TInt8         int
	TUint8        int
	TInt16        int
	TUint16       int
	TBlob         int
}
type NQGField struct {
	Field                *descriptor.FieldDescriptorProto
	MessageName          string
	TcapColumnType       string //TcaplusDataType_EN_TCAPLUS_DATATYPE_UINT64
	TcapColumnName       string //proto.String("uin")
	TableStructFieldName string //reqTable.Uin
	TcapColumnValueName  string //key.ValueUint
	DescriptorTypeName   string

	KeyIndex    int
	KeyIndexStr string
	////request
	KeyCodeForRequest         string
	NonKeyCodeForUpdateInsert string
	NonKeyCodeForSelect       string
	NonKeyCodeForIncrease     string
	////response
	KeyCodeForResponseField    string
	NonKeyCodeForResponseField string
	////API
	APICodeForIncreaseField string

	NqgOptions NQGOptions
	//types
	protoType           string //Proto.Int32
	goType              string //int32
	tcaplusValue        string //ValueUint
	tcaplusProtoType    string //
	tcaplusGoType       string //uint64/int64
	tcaplusDefaultValue string //
	IsSignedNumber      bool
}

func upperFirst(text string) string {
	if len(text) < 1 {
		return text
	}
	head := strings.ToUpper(text[0:1])
	return head + text[1:]
}
func snake2camel(text string) string {
	camel := ""
	parts := strings.Split(text, "_")
	for _, part := range parts {
		camel = camel + upperFirst(part)
	}
	return camel
}

func messageOptionsToNQGOption(o *NQGOptions, text string) *NQGOptions {
	splitted := strings.Fields(text)
	for _, part := range splitted {
		kv := strings.Split(part, ":")
		//log.Print("part:", part, ",kv:", kv)

		if len(kv) < 2 {
			continue
		}

		//70002:"nqg_user_score"
		keyID := kv[0]

		//[tcaplusgatesvr.table_name]:"nqg_user_item"
		len0 := len(kv[0])
		keyStr := kv[0][1 : len0-1]

		val := kv[1]
		if keyID == "70002" || keyStr == "tcaplusgatesvr.table_name" {
			o.TableName = strings.Replace(val, "\"", "", -1)
			o.CamelTableName = CamelCase(o.TableName)
		} else if keyID == "60000" || keyStr == "tcaplusgatesvr.b_key" {
			o.BKey, _ = strconv.Atoi(val)

		} else if keyID == "60006" || keyStr == "tcaplusgatesvr.b_custom_type" {
			o.BCustomerType, _ = strconv.Atoi(val)

		} else if keyID == "60001" || keyStr == "tcaplusgatesvr.t_int8" {
			o.TInt8, _ = strconv.Atoi(val)

		} else if keyID == "60002" || keyStr == "tcaplusgatesvr.t_uint8" {
			o.TUint8, _ = strconv.Atoi(val)

		} else if keyID == "60003" || keyStr == "tcaplusgatesvr.t_int16" {
			o.TInt16, _ = strconv.Atoi(val)

		} else if keyID == "60004" || keyStr == "tcaplusgatesvr.t_uint16" {
			o.TUint16, _ = strconv.Atoi(val)

		} else if keyID == "60005" || keyStr == "tcaplusgatesvr.t_blob" {
			o.TBlob, _ = strconv.Atoi(val)

		} else {

		}

	}
	return o
}

func getTcapColumnType(typename string, nqgOptions *NQGOptions) string {
	ret := "undefine"
	switch typename {
	case "*uint64", "uint64":
		ret = "tcaplusgatesvr.TcaplusDataType_EN_TCAPLUS_DATATYPE_UINT64"
	case "*int64", "int64":
		ret = "tcaplusgatesvr.TcaplusDataType_EN_TCAPLUS_DATATYPE_INT64"

	case "*uint32", "uint32":
		ret = "tcaplusgatesvr.TcaplusDataType_EN_TCAPLUS_DATATYPE_UINT32"
	case "*int32", "int32":
		ret = "tcaplusgatesvr.TcaplusDataType_EN_TCAPLUS_DATATYPE_INT32"

	case "*uint16", "uint16":
		ret = "tcaplusgatesvr.TcaplusDataType_EN_TCAPLUS_DATATYPE_UINT16"
	case "*int16", "int16":
		ret = "tcaplusgatesvr.TcaplusDataType_EN_TCAPLUS_DATATYPE_INT16"

	case "*uint8", "uint8":
		ret = "tcaplusgatesvr.TcaplusDataType_EN_TCAPLUS_DATATYPE_UINT8"
	case "*int8", "int8":
		ret = "tcaplusgatesvr.TcaplusDataType_EN_TCAPLUS_DATATYPE_INT8"

	case "*string", "string":
		ret = "tcaplusgatesvr.TcaplusDataType_EN_TCAPLUS_DATATYPE_STRING"
	default:

	}

	if nqgOptions.BCustomerType == 1 {
		if nqgOptions.TInt8 == 1 {
			ret = "tcaplusgatesvr.TcaplusDataType_EN_TCAPLUS_DATATYPE_INT8"
		}

		if nqgOptions.TUint8 == 1 {
			ret = "tcaplusgatesvr.TcaplusDataType_EN_TCAPLUS_DATATYPE_UINT8"
		}

		if nqgOptions.TInt16 == 1 {
			ret = "tcaplusgatesvr.TcaplusDataType_EN_TCAPLUS_DATATYPE_INT16"
		}

		if nqgOptions.TUint16 == 1 {
			ret = "tcaplusgatesvr.TcaplusDataType_EN_TCAPLUS_DATATYPE_UINT16"
		}

		if nqgOptions.TBlob == 1 {
			ret = "tcaplusgatesvr.TcaplusDataType_EN_TCAPLUS_DATATYPE_BLOB"
		}
	}
	return ret
}
func getTcapColumnValueName(typename string) string {
	ret := "undefine"
	switch typename {
	case "*uint64", "uint64", "*uint32", "uint32", "*uint16", "uint16", "*uint8", "uint8":
		ret = "ValueUint"
	case "*int64", "int64", "*int32", "int32", "*int16", "int16", "*int8", "int8":
		ret = "ValueInt"
	case "*string", "string":
		ret = "ValueStr"
	default:

	}
	return ret
}

func (g *Generator) generateTcapAPI(file *FileDescriptor) {
	file = g.FileOf(file.FileDescriptorProto)
	for _, desc := range file.desc {
		// Don't generate virtual messages for maps.
		if desc.GetOptions().GetMapEntry() {
			continue
		}

		g.dealwithTcapTableMessage(desc)
	}

}

var patternCamelTableName = "[[CamelTableName]]"
var patternSnakeTableName = "[[SnakeTableName]]"
var patternOrgMessageName = "[[OrgMessageName]]"

var codePieceImport = `
package tcaplusgatesvrapi

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"qqgame/baselib/logs"
	"qqgame/frame"
	"qqgame/msg/protobuf/tcaplusgatesvr"
)
`

var codePieceHighLevelAPI = `
type [[CamelTableName]]TableAPI struct {
	session   *qgframe.Session
	appID     int32
	tableName string
}

func New[[CamelTableName]]TableAPI(s *qgframe.Session, appID int32) *[[CamelTableName]]TableAPI {
	return &[[CamelTableName]]TableAPI{session: s, appID: appID, tableName: "[[SnakeTableName]]"}
}

func (e *[[CamelTableName]]TableAPI) Select(uin uint64, req *tcaplusgatesvr.[[OrgMessageName]]) (rsp *tcaplusgatesvr.[[OrgMessageName]], err error) {
	return e.DoRequestPlain(uin, req, EN_TCAPLUS_GATE_UTIL_OPT_SELECT, "", 0)
}

func (e *[[CamelTableName]]TableAPI) SelectByPartKey(uin uint64, req *tcaplusgatesvr.[[OrgMessageName]]) (rsp []*tcaplusgatesvr.[[OrgMessageName]], err error) {
	return e.DoRequestPlainMultiRow(uin, req, EN_TCAPLUS_GATE_UTIL_OPT_SELECTPARTKEY)
}

func (e *[[CamelTableName]]TableAPI) Insert(uin uint64, req *tcaplusgatesvr.[[OrgMessageName]]) (rsp *tcaplusgatesvr.[[OrgMessageName]], err error) {
	return e.DoRequestPlain(uin, req, EN_TCAPLUS_GATE_UTIL_OPT_INSERT, "", 0)
}

func (e *[[CamelTableName]]TableAPI) Update(uin uint64, req *tcaplusgatesvr.[[OrgMessageName]]) (rsp *tcaplusgatesvr.[[OrgMessageName]], err error) {
	return e.DoRequestPlain(uin, req, EN_TCAPLUS_GATE_UTIL_OPT_UPDATE, "", 0)
}

func (e *[[CamelTableName]]TableAPI) UpdateInsert(uin uint64, req *tcaplusgatesvr.[[OrgMessageName]]) (rsp *tcaplusgatesvr.[[OrgMessageName]], err error) {

	rsp, err = e.Update(uin, req)
	if err == nil {
		return rsp, err
	}

	if err == ErrMsgRecordNotExist || err == ErrMsgInvalidArgument {
		return e.Insert(uin, req)
	}

	return rsp, err
}

func (e *[[CamelTableName]]TableAPI) Increase(uin uint64, req *tcaplusgatesvr.[[OrgMessageName]]) (rsp *tcaplusgatesvr.[[OrgMessageName]], err error) {
	return e.DoRequestPlain(uin, req, EN_TCAPLUS_GATE_UTIL_OPT_INCREASE, "", 0)
}
`
var codePieceLevelTwoAPI = `
func (e *[[CamelTableName]]TableAPI) DoRequestPlain(uin uint64, req *tcaplusgatesvr.[[OrgMessageName]], gateOpt int32, deltaField string, delta int32) (rsp *tcaplusgatesvr.[[OrgMessageName]], err error) {

	tcaplusgateApi := NewTcaplusGateSvrAPI(e.session)
	if tcaplusgateApi != nil {

		//data type transfer : DBNQGUserCoin-->SSMsgReqTcaplusPlainTable
		var reqGate *tcaplusgatesvr.SSMsgReqTcaplusPlainTable
		reqGate, err = e.TransMessageToTcaplusRequest(req, gateOpt, deltaField, delta)
		if err != nil {
			return rsp, err
		}

		//SSMsgRspTcaplusPlainTable
		logs.DEBUGLOG("uin:%d, gateOpt:%d, call RequestPlain start, reqGate:%+v", uin, gateOpt, reqGate)
		rspGate, errTcaplus := tcaplusgateApi.RequestPlain(reqGate)
		if errTcaplus != nil {
			logs.ERRORLOG("uin:%d,call RequestPlain failed, gateOpt:%d, errTcaplus:%s ", uin, gateOpt, errTcaplus.Error())
			return rsp, ErrMsgGateError
		}
		logs.DEBUGLOG("uin:%d, gateOpt:%d, call RequestPlain success, rspGate:%+v", uin, gateOpt, rspGate)

		//response data type transfer: SSMsgRspTcaplusPlainTable-->DBNQGUserCoin
		rsp, err = e.TransTcaplusPlainTableRspToMemssage(rspGate, gateOpt)
		if err != nil {
			//logs.ERRORLOG("uin:%d, call TransTcaplusPlainTableRspToMemssage, gateOpt:%d, err:%s ", uin, gateOpt, err.Error())
			return rsp, err
		}

	} else {
		return rsp, errors.New("NewTcaplusGateSvrAPI failed")
	}

	return rsp, nil

}
func (e *[[CamelTableName]]TableAPI) DoRequestPlainMultiRow(uin uint64, req *tcaplusgatesvr.[[OrgMessageName]], gateOpt int32) (rsp []*tcaplusgatesvr.[[OrgMessageName]], err error) {

	tcaplusgateApi := NewTcaplusGateSvrAPI(e.session)
	if tcaplusgateApi != nil {

		//data type transfer : DBNQGUserCoin-->SSMsgReqTcaplusPlainTable
		var reqGate *tcaplusgatesvr.SSMsgReqTcaplusPlainTable
		reqGate, err = e.TransMessageToTcaplusRequest(req, gateOpt, "", 0)
		if err != nil {
			return rsp, err
		}

		//SSMsgRspTcaplusPlainTable
		logs.DEBUGLOG("uin:%d, gateOpt:%d, call RequestPlain start, reqGate:%+v", gateOpt, uin, reqGate)
		rspGate, errTcaplus := tcaplusgateApi.RequestPlain(reqGate)
		if errTcaplus != nil {
			logs.ERRORLOG("uin:%d,call RequestPlain failed, gateOpt:%d, errTcaplus:%s ", uin, gateOpt, errTcaplus.Error())
			return rsp, ErrMsgGateError
		}
		logs.DEBUGLOG("uin:%d, gateOpt:%d, call RequestPlain success, rspGate:%+v", gateOpt, uin, rspGate)

		//response data type transfer: SSMsgRspTcaplusPlainTable-->DBNQGUserCoin
		rsp, err = e.TransTcaplusPlainTableRspToMemssageMultiRow(rspGate, gateOpt)
		if err != nil {
			//logs.ERRORLOG("uin:%d, call TransTcaplusPlainTableRspToMemssage, gateOpt:%d, err:%s ", uin, gateOpt, err.Error())
			return rsp, err
		}

	} else {
		return rsp, errors.New("NewTcaplusGateSvrAPI failed")
	}

	return rsp, nil

}
`

var codePieceTransTcaplusRequestHead = `
func (e *[[CamelTableName]]TableAPI) TransMessageToTcaplusRequest(
	reqTable *tcaplusgatesvr.[[OrgMessageName]], opt int32, deltaField string, delta int32) (*tcaplusgatesvr.SSMsgReqTcaplusPlainTable, error) {

	reqTcaplus := &tcaplusgatesvr.SSMsgReqTcaplusPlainTable{}
	reqTcaplus.TcaplusGateId = proto.Int32(e.appID)
	reqTcaplus.TableName = proto.String(e.tableName)
	reqTcaplus.ApiCmd = new(tcaplusgatesvr.TcaplusAPICmd)

	requestSet := &tcaplusgatesvr.DataPlainRow{}
	reqTcaplus.RequestSet = requestSet

`

var codePieceTransTcaplusRequestTail = `
	return reqTcaplus, nil
}
`
var codePieceTransTcaplusResponseHead = `
func (e *[[CamelTableName]]TableAPI) TransTcaplusPlainTableRspToMemssage(
	rspTcaplus *tcaplusgatesvr.SSMsgRspTcaplusPlainTable, gateOpt int32) (*tcaplusgatesvr.[[OrgMessageName]], error) {

	rspOnerow := &tcaplusgatesvr.[[OrgMessageName]]{}
	result := rspTcaplus.GetResult()
	if result != 0 {
		return nil, TcaplusErrorID2Error(result)
	}

	rspRows := rspTcaplus.GetResponseRows()
	if len(rspRows) == 0 {
		//return rspOnerow, ErrMsgOperationSuccButNoRecord
		return rspOnerow, nil
	}
	var pDataRow *tcaplusgatesvr.DataPlainRow
	pDataRow = rspRows[0]

	//traverse all values
	for _, pColumn := range pDataRow.Values {
`
var codePieceTransTcaplusResponseTail = `
		default:
			//

		}
	}
	_ = gateOpt
	return rspOnerow, nil
}
`

var codePieceTransTcaplusResponseMultiRowHead = `
func (e *[[CamelTableName]]TableAPI) TransTcaplusPlainTableRspToMemssageMultiRow(
	rspTcaplus *tcaplusgatesvr.SSMsgRspTcaplusPlainTable, gateOpt int32) ([]*tcaplusgatesvr.[[OrgMessageName]], error) {

	rsp := make([]*tcaplusgatesvr.[[OrgMessageName]], 0)
	result := rspTcaplus.GetResult()
	if result != 0 {
		return nil, TcaplusErrorID2Error(result)
	}

	gateRspRows := rspTcaplus.GetResponseRows()
	if len(gateRspRows) == 0 {
		//return rsp, ErrMsgOperationSuccButNoRecord
		return rsp, nil
	}

	//traverse all rows
	for _, pDataRow := range gateRspRows {
		rspOnerow := &tcaplusgatesvr.[[OrgMessageName]]{}
`
var codePieceTransTcaplusResponseMultiRowTail = `
		rsp = append(rsp, rspOnerow)
	}

	_ = gateOpt
	return rsp, nil
}
`

var codePieceKeyFieldForRequest = `
	{
		key := &tcaplusgatesvr.DataPlainColumn{}
		key.Type = new(tcaplusgatesvr.TcaplusDataType)
		*key.Type = [[TcapColumnType]]
		key.Name = proto.String("[[TcapColumnName]]")

		if reqTable.[[TableStructFieldName]] != nil { //[[TcapColumnKeyIndex]]
			key.[[TcapColumnValueName]] = [[TcaplusProtoType]]([[TcaplusGoType]](*reqTable.[[TableStructFieldName]]))
			if opt == EN_TCAPLUS_GATE_UTIL_OPT_SELECTPARTKEY {
				key.Flag = proto.Int32(int32(tcaplusgatesvr.DataColumnFlag_EN_DATA_COLUMN_FLAG_NORMAL))
			}
			requestSet.[[TcapColumnKeyIndex]] = key
		} else {
			if opt == EN_TCAPLUS_GATE_UTIL_OPT_SELECTPARTKEY {
				key.Flag = proto.Int32(int32(tcaplusgatesvr.DataColumnFlag_EN_DATA_COLUMN_FLAG_FOR_RSP_SET))
				requestSet.[[TcapColumnKeyIndex]] = key
			}
		}
	}
`

////////////////////////////////////////////
var codePieceNonKeyFieldForUpdateInsert = `
		if reqTable.[[TableStructFieldName]] != nil {
			oneColumn := &tcaplusgatesvr.DataPlainColumn{}
			oneColumn.Type = new(tcaplusgatesvr.TcaplusDataType)
			*oneColumn.Type = [[TcapColumnType]]
			oneColumn.Name = proto.String("[[TcapColumnName]]")
			oneColumn.[[TcaplusValue]] = [[TcaplusProtoType]]([[TcaplusGoType]](*reqTable.[[TableStructFieldName]]))
			allValues = append(allValues, oneColumn)
		}
`

////////////////////////////////////////////
var codePieceNonKeyFieldForSelect = `
		{
			oneColumn := &tcaplusgatesvr.DataPlainColumn{}
			oneColumn.Type = new(tcaplusgatesvr.TcaplusDataType)
			*oneColumn.Type = [[TcapColumnType]]
			oneColumn.Name = proto.String("[[TcapColumnName]]")
			oneColumn.[[TcaplusValue]] = [[TcaplusProtoType]]([[TcaplusDefaultValue]])
			tmpValues = append(tmpValues, oneColumn)
		}
`

////////////////////////////////////////////
var codePieceNonKeyFieldForIncreaseByPointer = `
			}else if reqTable.[[TableStructFieldName]] != nil {
				if increaseFieldNum > 0 { //increaseFieldNum can not bigger than 1
					return reqTcaplus, ErrMsgCmdInputError
				}
				tmpDelta := *reqTable.[[TableStructFieldName]]
				oneColumn := &tcaplusgatesvr.DataPlainColumn{}
				oneColumn.Type = new(tcaplusgatesvr.TcaplusDataType)
				*oneColumn.Type = [[TcapColumnType]]
				oneColumn.Name = proto.String("[[TcapColumnName]]")
				if tmpDelta >= 0 {
					oneColumn.[[TcaplusValue]] = [[TcaplusProtoType]]([[TcaplusGoType]](tmpDelta))
					increaseFieldOptionType = tcaplusgatesvr.OperationType_EN_OPERATION_TYPE_PLUS
				} else {
					oneColumn.[[TcaplusValue]] = [[TcaplusProtoType]]([[TcaplusGoType]](-1 * tmpDelta))
				}

				tmpValues := make([]*tcaplusgatesvr.DataPlainColumn, 1)
				tmpValues[0] = oneColumn
				requestSet.Values = tmpValues

				increaseFieldNum = increaseFieldNum + 1
				increaseFieldName = "[[TcapColumnName]]"
			}
`
var codePieceNonKeyFieldForIncrease = `
			if deltaField == "[[TcapColumnName]]" {
				oneColumn := &tcaplusgatesvr.DataPlainColumn{}
				oneColumn.Type = new(tcaplusgatesvr.TcaplusDataType)
				*oneColumn.Type = [[TcapColumnType]]
				oneColumn.Name = proto.String("[[TcapColumnName]]")
				if delta >= 0 {
					oneColumn.ValueUint = proto.Uint64(uint64(delta))
					increaseFieldOptionType = tcaplusgatesvr.OperationType_EN_OPERATION_TYPE_PLUS
				} else {
					oneColumn.ValueUint = proto.Uint64(uint64(-1 * delta))
				}

				tmpValues := make([]*tcaplusgatesvr.DataPlainColumn, 1)
				tmpValues[0] = oneColumn
				requestSet.Values = tmpValues

				increaseFieldNum = increaseFieldNum + 1
				increaseFieldName = "[[TcapColumnName]]"
				break
`
var codePieceAPICodeForFieldIncrease = `
func (e *[[CamelTableName]]TableAPI) Increase_[[TcapColumnName]](uin uint64[[ParamListForKeys]], delta int32) (rsp *tcaplusgatesvr.[[OrgMessageName]], err error) {
	req:=&tcaplusgatesvr.[[OrgMessageName]]{
[[FieldAssignForKeys]]
	}

	var errIncr error
	rsp, errIncr = e.DoRequestPlain(uin, req, EN_TCAPLUS_GATE_UTIL_OPT_INCREASE, "[[TcapColumnName]]", delta)
	if errIncr == ErrMsgRecordNotExist{
		req.[[TableStructFieldName]] = proto.[[ProtoType]]([[GoType]](delta))
		rspInsert, errInsert := e.DoRequestPlain(uin, req, EN_TCAPLUS_GATE_UTIL_OPT_INSERT, "", 0)
		return rspInsert, errInsert
	}
	return rsp, errIncr
}
`

func replacePattern(org string, nqgOption *NQGOptions, messageName string) string {

	camelTableName := nqgOption.CamelTableName
	snakeTableName := nqgOption.TableName

	ret := ""
	tmpOne := strings.Replace(org, patternCamelTableName, camelTableName, -1)
	tmpTwo := strings.Replace(tmpOne, patternSnakeTableName, snakeTableName, -1)
	ret = strings.Replace(tmpTwo, patternOrgMessageName, messageName, -1)
	return ret
}
func generateFieldCode(nqgField *NQGField,
	paramListCodeForKeys string, fieldAssignCodeForKeys string) {

	messageName := nqgField.MessageName
	fieldOption := nqgField.NqgOptions
	camelTableName := fieldOption.CamelTableName
	snakeTableName := fieldOption.TableName

	protoType := nqgField.protoType
	goType := nqgField.goType
	tcaplusValue := nqgField.tcaplusValue
	tcaplusProtoType := nqgField.tcaplusProtoType
	tcaplusGoType := nqgField.tcaplusGoType
	tcaplusDefaultValue := nqgField.tcaplusDefaultValue

	r := strings.NewReplacer("[[TcapColumnType]]", nqgField.TcapColumnType,
		"[[TcapColumnName]]", nqgField.TcapColumnName,
		"[[TableStructFieldName]]", nqgField.TableStructFieldName,
		"[[TcapColumnValueName]]", nqgField.TcapColumnValueName,
		"[[TcapColumnKeyIndex]]", nqgField.KeyIndexStr,
		"[[ProtoType]]", protoType,
		"[[GoType]]", goType,
		"[[TcaplusValue]]", tcaplusValue,
		"[[TcaplusProtoType]]", tcaplusProtoType,
		"[[TcaplusGoType]]", tcaplusGoType,
		"[[TcaplusDefaultValue]]", tcaplusDefaultValue,
		patternCamelTableName, camelTableName,
		patternOrgMessageName, messageName,
		patternSnakeTableName, snakeTableName,
		"[[ParamListForKeys]]", paramListCodeForKeys,
		"[[FieldAssignForKeys]]", fieldAssignCodeForKeys)

	if nqgField.NqgOptions.BKey == 1 {
		/////
		nqgField.KeyCodeForRequest = r.Replace(codePieceKeyFieldForRequest)
		code := `
		if pDataRow.Get[[TcapColumnKeyIndex]]() != nil {
			rspOnerow.[[TableStructFieldName]] = proto.[[ProtoType]]([[GoType]](pDataRow.Get[[TcapColumnKeyIndex]]().Get[[TcaplusValue]]()))
		}
`
		nqgField.KeyCodeForResponseField = r.Replace(code)

	} else {
		/////
		nqgField.NonKeyCodeForUpdateInsert = r.Replace(codePieceNonKeyFieldForUpdateInsert)
		nqgField.NonKeyCodeForSelect = r.Replace(codePieceNonKeyFieldForSelect)
		increaseModeByFieldName := r.Replace(codePieceNonKeyFieldForIncrease)
		increaseModeByFieldPointer := r.Replace(codePieceNonKeyFieldForIncreaseByPointer)

		if nqgField.IsSignedNumber {
			nqgField.NonKeyCodeForIncrease = increaseModeByFieldName + increaseModeByFieldPointer
		} else {
			nqgField.NonKeyCodeForIncrease = increaseModeByFieldName + `
			}
`
		}

		code := `
			case "[[TcapColumnName]]":
				if pColumn.[[TcaplusValue]] != nil {
					rspOnerow.[[TableStructFieldName]] = proto.[[ProtoType]]([[GoType]](pColumn.Get[[TcaplusValue]]()))
				}
`
		nqgField.NonKeyCodeForResponseField = r.Replace(code)

		nqgField.APICodeForIncreaseField = r.Replace(codePieceAPICodeForFieldIncrease)
	}
	return
}
func IsSignedNumber(typeName string) (isNumber bool, signed bool) {
	switch typeName {
	case "*uint64", "uint64", "*uint32", "uint32", "*uint16", "uint16", "*uint8", "uint8":
		return true, false
	case "*int64", "int64", "*int32", "int32", "*int16", "int16", "*int8", "int8":
		return true, true
	}
	return
}
func (g *Generator) protoField2NQGField(message *Descriptor,
	field *descriptor.FieldDescriptorProto,
	keyIndex int,
	tableName string,
	camelTableName string) (nqgField *NQGField) {

	nqgField = &NQGField{}
	nqgField.MessageName = message.GetName()
	nqgField.Field = field
	nqgField.KeyIndex = keyIndex
	nqgField.KeyIndexStr = fmt.Sprintf("Key_%d", keyIndex)

	fieldName := *field.Name
	typename, wiretype := g.GoType(message, field)
	_ = wiretype

	messageOptionsToNQGOption(&nqgField.NqgOptions, field.GetOptions().String())
	nqgField.NqgOptions.TableName = tableName
	nqgField.NqgOptions.CamelTableName = camelTableName

	nqgField.DescriptorTypeName = typename
	nqgField.TcapColumnType = getTcapColumnType(typename, &nqgField.NqgOptions)
	nqgField.TcapColumnName = fieldName
	nqgField.TableStructFieldName = CamelCase(fieldName)
	nqgField.TcapColumnValueName = getTcapColumnValueName(typename)

	var undefineText = "undefine"
	protoType := undefineText
	goType := undefineText
	tcaplusValue := undefineText
	tcaplusProtoType := undefineText
	tcaplusGoType := undefineText
	tcaplusDefaultValue := undefineText

	switch nqgField.DescriptorTypeName {
	case "*uint64", "uint64":
		protoType = "Uint64"
		goType = "uint64"
		tcaplusValue = "ValueUint"
		tcaplusProtoType = "proto.Uint64"
		tcaplusGoType = "uint64"
		tcaplusDefaultValue = "0"
	case "*uint32", "uint32":
		protoType = "Uint32"
		goType = "uint32"
		tcaplusValue = "ValueUint"
		tcaplusProtoType = "proto.Uint64"
		tcaplusGoType = "uint64"
		tcaplusDefaultValue = "0"

	case "*uint16", "uint16":
		protoType = "Uint16"
		goType = "uint16"
		tcaplusValue = "ValueUint"
		tcaplusProtoType = "proto.Uint64"
		tcaplusGoType = "uint64"
		tcaplusDefaultValue = "0"

	case "*uint8", "uint8":
		protoType = "Uint8"
		goType = "uint8"
		tcaplusValue = "ValueUint"
		tcaplusProtoType = "proto.Uint64"
		tcaplusGoType = "uint64"
		tcaplusDefaultValue = "0"

	case "*int64", "int64":
		protoType = "Int64"
		goType = "int64"
		tcaplusValue = "ValueInt"
		tcaplusProtoType = "proto.Int64"
		tcaplusGoType = "int64"
		tcaplusDefaultValue = "0"

	case "*int32", "int32":
		protoType = "Int32"
		goType = "int32"
		tcaplusValue = "ValueInt"
		tcaplusProtoType = "proto.Int64"
		tcaplusGoType = "int64"
		tcaplusDefaultValue = "0"

	case "*int16", "int16":
		protoType = "Int16"
		goType = "int16"
		tcaplusValue = "ValueInt"
		tcaplusProtoType = "proto.Int64"
		tcaplusGoType = "int64"
		tcaplusDefaultValue = "0"

	case "*int8", "int8":
		protoType = "Int8"
		goType = "int8"
		tcaplusValue = "ValueInt"
		tcaplusProtoType = "proto.Int64"
		tcaplusGoType = "int64"
		tcaplusDefaultValue = "0"

	case "*string", "string":
		protoType = "String"
		goType = "string"
		tcaplusValue = "ValueStr"
		tcaplusProtoType = "[]byte"
		tcaplusGoType = "string"
		tcaplusDefaultValue = "\"\""
	default:

	}

	if nqgField.NqgOptions.TBlob == 1 {
		tcaplusProtoType = "[]byte"
	}

	_, nqgField.IsSignedNumber = IsSignedNumber(nqgField.DescriptorTypeName)

	nqgField.protoType = protoType
	nqgField.goType = goType
	nqgField.tcaplusValue = tcaplusValue
	nqgField.tcaplusProtoType = tcaplusProtoType
	nqgField.tcaplusGoType = tcaplusGoType
	nqgField.tcaplusDefaultValue = tcaplusDefaultValue

	return
}

// Generate the type and default constant definitions for this Descriptor.
func (g *Generator) dealwithTcapTableMessage(message *Descriptor) {

	isTcaplusTableMessage := false

	messageName := message.GetName()
	log.Print("trace messages:", message.GetName(), ",options:", fmt.Sprintf("%+v", message.GetOptions()))
	messageOptions := message.GetOptions()
	if messageOptions == nil {
		return
	}

	var nqgOptionForMsg = &NQGOptions{}
	{ //start
		typeMsgOptions := reflect.TypeOf(*messageOptions)
		valueMsgOptions := reflect.ValueOf(*messageOptions)
		_ = typeMsgOptions
		_ = valueMsgOptions
		messageOptionsToNQGOption(nqgOptionForMsg, messageOptions.String())
		log.Print("options:", messageOptions.String(), ",nqgOptionForMsg:", fmt.Sprintf("%+v", nqgOptionForMsg))
	} //end

	if nqgOptionForMsg.TableName != "" {
		isTcaplusTableMessage = true
	}
	if !isTcaplusTableMessage {
		return
	}

	log.Print("tableName:", nqgOptionForMsg.TableName, ",camel:", nqgOptionForMsg.CamelTableName)

	fieldsOfTableKey := make([]NQGField, 0)
	fieldsOfTableNonKey := make([]NQGField, 0)
	paramListCodeForKeys := ""
	fieldAssignCodeForKeys := ""
	{
		keyIndex := 1
		for _, field := range message.Field {
			nqgField := g.protoField2NQGField(message, field, keyIndex, nqgOptionForMsg.TableName, nqgOptionForMsg.CamelTableName)

			if nqgField.NqgOptions.BKey == 1 {

				r := strings.NewReplacer("[[TableStructFieldName]]", nqgField.TableStructFieldName,
					"[[GoType]]", nqgField.goType,
					"[[ProtoType]]", nqgField.protoType)

				codePieceParam := `, [[TableStructFieldName]] [[GoType]]`
				paramListCodeForKeys = paramListCodeForKeys + r.Replace(codePieceParam)

				codePieceFieldAssign := "\t\t[[TableStructFieldName]] : proto.[[ProtoType]]([[TableStructFieldName]]),\n"
				fieldAssignCodeForKeys = fieldAssignCodeForKeys + r.Replace(codePieceFieldAssign)

				keyIndex++
			}
		}
	}

	{
		keyIndex := 1
		for i, field := range message.Field {
			_ = i
			nqgField := g.protoField2NQGField(message, field, keyIndex, nqgOptionForMsg.TableName, nqgOptionForMsg.CamelTableName)
			if nqgField.NqgOptions.BKey == 1 {
				generateFieldCode(nqgField, paramListCodeForKeys, fieldAssignCodeForKeys)
				fieldsOfTableKey = append(fieldsOfTableKey, *nqgField)
				keyIndex++
			} else {
				generateFieldCode(nqgField, paramListCodeForKeys, fieldAssignCodeForKeys)
				fieldsOfTableNonKey = append(fieldsOfTableNonKey, *nqgField)
			}
		}
	}

	//first write
	g.WriteString(codePieceImport)

	//1.codePieceHighLevelAPI
	g.WriteString(replacePattern(codePieceHighLevelAPI, nqgOptionForMsg, messageName))
	for _, nqgField := range fieldsOfTableNonKey {
		switch nqgField.DescriptorTypeName {
		case "*uint64", "uint64", "*uint32", "uint32", "*uint16", "uint16", "*uint8", "uint8":
			g.WriteString(nqgField.APICodeForIncreaseField)
		case "*int64", "int64", "*int32", "int32", "*int16", "int16", "*int8", "int8":
			g.WriteString(nqgField.APICodeForIncreaseField)
		default:

		}
	}
	g.WriteString(replacePattern(codePieceLevelTwoAPI, nqgOptionForMsg, messageName))

	//2.TransTcaplusRequest
	g.WriteString(replacePattern(codePieceTransTcaplusRequestHead, nqgOptionForMsg, messageName))
	g.WriteString("\n\t//do with keys")
	for _, nqgField := range fieldsOfTableKey {
		g.WriteString(nqgField.KeyCodeForRequest)
	}
	g.WriteString("\n\t//do with non-keys")

	//allValues
	g.WriteString(`
	allValues := make([]*tcaplusgatesvr.DataPlainColumn, 0)
	for {
		if opt != EN_TCAPLUS_GATE_UTIL_OPT_INSERT && opt != EN_TCAPLUS_GATE_UTIL_OPT_UPDATE {
			break
		}
`)
	for _, nqgField := range fieldsOfTableNonKey {
		g.WriteString(nqgField.NonKeyCodeForUpdateInsert)
	}
	g.WriteString(`
		break
	}
`)

	//switch opt
	g.WriteString(`
	switch opt {
	case EN_TCAPLUS_GATE_UTIL_OPT_SELECT,
		EN_TCAPLUS_GATE_UTIL_OPT_SELECTPARTKEY:
		if opt == EN_TCAPLUS_GATE_UTIL_OPT_SELECT {
			*reqTcaplus.ApiCmd = tcaplusgatesvr.TcaplusAPICmd_ENTCAPLUS_API_REQ_GET
		} else if opt == EN_TCAPLUS_GATE_UTIL_OPT_SELECTPARTKEY {
			*reqTcaplus.ApiCmd = tcaplusgatesvr.TcaplusAPICmd_ENTCAPLUS_API_REQ_GET_BY_PARTKEY
		}
		tmpValues := make([]*tcaplusgatesvr.DataPlainColumn, 0)
`)
	for _, nqgField := range fieldsOfTableNonKey {
		g.WriteString(nqgField.NonKeyCodeForSelect)
	}
	g.WriteString(`
		requestSet.Values = tmpValues

		reqTcaplus.Options = new(tcaplusgatesvr.DataTcaplusRequestOption)
		//reqTcaplus.Options.Flags = proto.Uint32(2) //need response field data

	case EN_TCAPLUS_GATE_UTIL_OPT_INSERT:
		*reqTcaplus.ApiCmd = tcaplusgatesvr.TcaplusAPICmd_ENTCAPLUS_API_REQ_INSERT

		requestSet.Values = allValues

		reqTcaplus.Options = new(tcaplusgatesvr.DataTcaplusRequestOption)
		reqTcaplus.Options.ResultFlag = proto.Uint32(1) //need response changed field data

	case EN_TCAPLUS_GATE_UTIL_OPT_UPDATE:
		*reqTcaplus.ApiCmd = tcaplusgatesvr.TcaplusAPICmd_ENTCAPLUS_API_REQ_UPDATE

		requestSet.Values = allValues

		//reqTcaplus.Options = new(tcaplusgatesvr.DataTcaplusRequestOption)
		//reqTcaplus.Options.ResultFlag = proto.Uint32(1) //need response changed field data

	case EN_TCAPLUS_GATE_UTIL_OPT_DELETE:
		*reqTcaplus.ApiCmd = tcaplusgatesvr.TcaplusAPICmd_ENTCAPLUS_API_REQ_DELETE

	case EN_TCAPLUS_GATE_UTIL_OPT_INCREASE:
		*reqTcaplus.ApiCmd = tcaplusgatesvr.TcaplusAPICmd_ENTCAPLUS_API_REQ_INCREASE

		//find szIncreaseFieldName
		var increaseFieldNum int32
		var increaseFieldName string
		increaseFieldOptionType := tcaplusgatesvr.OperationType_EN_OPERATION_TYPE_MINUS
		for {
`)

	for _, nqgField := range fieldsOfTableNonKey {
		switch nqgField.DescriptorTypeName {
		case "*uint64", "uint64", "*uint32", "uint32", "*uint16", "uint16", "*uint8", "uint8":
			g.WriteString(nqgField.NonKeyCodeForIncrease)
		case "*int64", "int64", "*int32", "int32", "*int16", "int16", "*int8", "int8":
			g.WriteString(nqgField.NonKeyCodeForIncrease)
		default:

		}
	}

	g.WriteString(`
			break
		}
		if increaseFieldNum <= 0 {
			return reqTcaplus, ErrMsgCmdInputError
		}

		reqTcaplus.Options = new(tcaplusgatesvr.DataTcaplusRequestOption)
		reqTcaplus.Options.ResultFlag = proto.Uint32(1) //need response changed field data

		reqTcaplus.Operation = new(tcaplusgatesvr.DataOperation)
		reqTcaplus.Operation.FieldName = proto.String(increaseFieldName)
		reqTcaplus.Operation.OperationType = new(tcaplusgatesvr.OperationType)
		*reqTcaplus.Operation.OperationType = increaseFieldOptionType

	default:
		//
	}
`)

	g.WriteString(replacePattern(codePieceTransTcaplusRequestTail, nqgOptionForMsg, messageName))

	//3.TransTcaplusResponse
	g.WriteString(replacePattern(codePieceTransTcaplusResponseHead, nqgOptionForMsg, messageName))
	for _, nqgField := range fieldsOfTableKey {
		g.WriteString(nqgField.KeyCodeForResponseField)
	}
	g.WriteString(`
		switch pColumn.GetName() {
`)
	for _, nqgField := range fieldsOfTableNonKey {
		g.WriteString(nqgField.NonKeyCodeForResponseField)
	}
	g.WriteString(replacePattern(codePieceTransTcaplusResponseTail, nqgOptionForMsg, messageName))

	//4.TransTcaplusResponseMultiRow
	g.WriteString(replacePattern(codePieceTransTcaplusResponseMultiRowHead, nqgOptionForMsg, messageName))
	for _, nqgField := range fieldsOfTableKey {
		g.WriteString(nqgField.KeyCodeForResponseField)
	}
	g.WriteString(`
		//traverse all values
		for _, pColumn := range pDataRow.Values {
			switch pColumn.GetName() {
`)
	for _, nqgField := range fieldsOfTableNonKey {
		g.WriteString(nqgField.NonKeyCodeForResponseField)
	}
	g.WriteString(`
			default:
				//

			}
		}
`)
	g.WriteString(replacePattern(codePieceTransTcaplusResponseMultiRowTail, nqgOptionForMsg, messageName))

	rootDir, errDir := filepath.Abs(filepath.Dir(os.Args[0]))
	if errDir != nil {
		log.Print("tcaptable filepath.Abs errDir:", errDir)
		return
	}
	outputDir := rootDir + "/tcapgateapi/"
	_, errStat := os.Stat(outputDir)
	if errStat != nil {
		errMK := os.Mkdir(outputDir, os.ModePerm)
		if errMK != nil {
			log.Print("tcaptable filepath.Abs errMK:", errMK)
			return
		}
	}

	log.Print("outputDir:", outputDir)
	goFileName := outputDir + nqgOptionForMsg.TableName + "_api.go"
	content := g.String()
	log.Print("goFileName:", goFileName)
	//log.Print("content:", content)

	errWrite := easyapi.SimpleWriteFile(goFileName, content, false)
	if errWrite != nil {
		log.Print("write file:", goFileName, ",fail:", errWrite)
	}
	/*rspFile := plugin.CodeGeneratorResponse_File{
		Name:    proto.String(goFileName),
		Content: proto.String(content),
	}*/

	g.Reset()
}
