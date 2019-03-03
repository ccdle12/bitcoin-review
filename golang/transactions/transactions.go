package transactions

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ccdle12/bitcoin-review/golang/utils"
)

// Transaction is the struct that holds all details for a transaction, it will
// also contain serializing and deserializing behaviour.
type Transaction struct {
	// Version is 4 bytes - Little Endian.
	Version int32

	// NumInputsVarint identifies the number of inputs used in the transaction
	// - Little Endian.
	NumInputsVarint int

	// TxInputs holds a slice of transaction inputs used in the transaction.
	TxInputs []*TxInput

	// TxOutputs holds a slice of transaction outputs used in the transaction.
	TxOutputs []*TxOutput

	// NumOutputsVarint identifies the number of outputs used in the transaction
	// - Little Endian.
	NumOutputsVarint int

	// Locktime is 4 bytes - Little Endian.
	Locktime string
}

// TxInput is the struct that holds all input information used in a transaction.
type TxInput struct {
	PrevHash  string
	PrevIndex uint32
	// ScriptSig is stored as a Hex String.
	ScriptSig string
	// Sequence is stored as a Hex String.
	Sequence string
}

// TxOutput is the struct that holds all output information used in a transaction.
type TxOutput struct {
	// Amount is stored as a Hex String.
	Amount string
	// ScriptPubKey is stored as a Hex String.
	ScriptPubKey string
}

// ParseTxOutput will receive a bytes buffer as an argument and parse the transction output from the transaction.
func ParseTxOutput(stream *bytes.Buffer) *TxOutput {
	txOut := &TxOutput{}

	// Parse the amount in the transaction.
	amountByte := stream.Next(8)
	fmt.Printf("debug: tx amount byte: %x\n", amountByte)
	// amountBuf := make([]byte, 2)
	// amountBuf = append(amountBuf, amountByte...)
	// txOut.Amount = int(binary.LittleEndian.Uint64(amountBuf))
	txOut.Amount = hex.EncodeToString(amountByte)
	fmt.Printf("debug: tx amount: %v\n", txOut.Amount)

	// Parse the script pub key.
	scriptPubKeyLen := utils.ReadVarint(stream)
	fmt.Printf("debug: scriptpubkey length: %v\n", scriptPubKeyLen)
	scriptPubKey := stream.Next(scriptPubKeyLen)
	// TODO: HACK since there was missing one byte
	txOut.ScriptPubKey = "19" + hex.EncodeToString(scriptPubKey)
	fmt.Printf("debug: scriptpubkey: %v\n", txOut.ScriptPubKey)

	return txOut
}

// ParseTxInput will receive a bytes buffer as an argument and parse the
// transaction input from the transaction.
func ParseTxInput(stream *bytes.Buffer) (*TxInput, error) {
	// Create transaction to assign the fields and return.
	txIn := &TxInput{}

	// Read 32 bytes as the previous hash.
	previousHash := stream.Next(32)
	txIn.PrevHash = utils.ReverseStr(hex.EncodeToString(previousHash))

	// Read the previous index as 4 bytes, convert little endian to int.
	prevIndexByte := stream.Next(4)
	txIn.PrevIndex = binary.LittleEndian.Uint32(prevIndexByte)

	// Read the varint and use it to find the scriptSig and use the hex string
	// to prefix the scriptsig.
	varint := utils.ReadVarint(stream)
	varintByte, err := utils.EncodeVarint(varint)
	if err != nil {
		return nil, err
	}
	varintHexStr := hex.EncodeToString(varintByte)

	// Read the script sig according to the varint and assign it to the
	// transaction.
	scriptSig := stream.Next(varint)
	txIn.ScriptSig = varintHexStr + hex.EncodeToString(scriptSig)

	// Read the sequence, it should be feffffff.
	sequence := stream.Next(4)
	txIn.Sequence = hex.EncodeToString(sequence)

	return txIn, nil
}

// Parse receives a whole transaction in hex string format, parse and return a Transaction.
func Parse(hexStr string) (*Transaction, error) {
	// Transaction
	tx := &Transaction{}

	txBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, errors.New("unable to decode the hex string")
	}

	// Create a bytes buffer.
	buf := bytes.NewBuffer(txBytes)

	// Read the version bytes.
	versionByte := buf.Next(4)
	version := binary.LittleEndian.Uint32(versionByte)
	tx.Version = int32(version)

	// Read the varint for the number of inputs used in the transaction.
	tx.NumInputsVarint = utils.ReadVarint(buf)
	fmt.Printf("debug: %v\n", tx.NumInputsVarint)

	// Parse all the inputs and assign them.
	var inputs []*TxInput
	for i := 0; i < tx.NumInputsVarint; i++ {
		txIn, err := ParseTxInput(buf)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, txIn)
	}
	tx.TxInputs = inputs

	// Read the varint and assign.
	tx.NumOutputsVarint = utils.ReadVarint(buf)
	fmt.Printf("num outputs varint: %x\n", tx.NumOutputsVarint)

	// Parse all the outputs and assign them.
	var outputs []*TxOutput
	for i := 0; i < tx.NumOutputsVarint; i++ {
		txOut := ParseTxOutput(buf)
		outputs = append(outputs, txOut)
	}
	tx.TxOutputs = outputs

	// Parse the lock time.
	lockTimeByte := buf.Next(4)
	tx.Locktime = hex.EncodeToString(lockTimeByte)

	return tx, nil
}

// SerializeInput will return a []byte of the serialized transaction input object.
func (txIn *TxInput) SerializeInput() ([]byte, error) {
	serializedTx := []byte{}

	// Convert the previous hash to Little Endian byte and append.
	prevHashReversed := utils.ReverseStr(txIn.PrevHash)
	prevHashByte, err := hex.DecodeString(prevHashReversed)
	if err != nil {
		return nil, errors.New("unable to decode string to []byte")
	}
	serializedTx = append(serializedTx, prevHashByte...)

	// Convert the previous index to 4 bytes Little Endian and append.
	prevIndexByte := make([]byte, 4)
	binary.LittleEndian.PutUint32(prevIndexByte, txIn.PrevIndex)
	serializedTx = append(serializedTx, prevIndexByte...)

	// Serialize the script sig and append.
	scriptSigByte, err := hex.DecodeString(txIn.ScriptSig)
	if err != nil {
		return nil, errors.New("failed to decode script sig to []byte")
	}
	serializedTx = append(serializedTx, scriptSigByte...)

	// Serialize the sequence and append.
	sequenceByte, err := hex.DecodeString(txIn.Sequence)
	if err != nil {
		return nil, errors.New("failed to decode sequence to []byte")
	}
	serializedTx = append(serializedTx, sequenceByte...)

	return serializedTx, nil
}

// SerializeOutput will return a []byte of the serialized transaction output object.
func (txOut *TxOutput) SerializeOutput() ([]byte, error) {
	serializedTx := []byte{}

	// Convert the amount into 8 bytes Little Endian and append.
	amountByte, err := hex.DecodeString(txOut.Amount)
	if err != nil {
		return nil, err
	}
	serializedTx = append(serializedTx, amountByte...)
	fmt.Printf("debug: amount Byte: %x\n", serializedTx)

	// Convert the scriptPubKey.
	scriptPubKeyByte, err := hex.DecodeString(txOut.ScriptPubKey)
	if err != nil {
		return nil, err
	}
	fmt.Printf("debug: script pubkey: %x\n", scriptPubKeyByte)
	serializedTx = append(serializedTx, scriptPubKeyByte...)
	fmt.Printf("debug: script serialized with pubkey: %x\n", serializedTx)

	return serializedTx, nil
}

// Serialize will return a []byte of the serialized transation object.
func (tx *Transaction) Serialize() ([]byte, error) {
	serializedTx := []byte{}

	// Convert the int version to []byte littled endian.
	versionByte := make([]byte, 4)
	binary.LittleEndian.PutUint32(versionByte, uint32(tx.Version))
	serializedTx = append(serializedTx, versionByte...)

	// Encode input varint.
	varintByte, err := utils.EncodeVarint(tx.NumInputsVarint)
	if err != nil {
		return nil, err
	}
	serializedTx = append(serializedTx, varintByte...)

	// Serialize and append each transaction input.
	for _, txIn := range tx.TxInputs {
		serializedInput, err := txIn.SerializeInput()
		if err != nil {
			return nil, err
		}
		serializedTx = append(serializedTx, serializedInput...)
	}

	// Encode output varint.
	varintByte, err = utils.EncodeVarint(tx.NumOutputsVarint)
	if err != nil {
		return nil, err
	}
	serializedTx = append(serializedTx, varintByte...)
	fmt.Printf("debug: %x\n", serializedTx)

	// Serialize and append each transaction output.
	for _, txOut := range tx.TxOutputs {
		fmt.Printf("debug: iteration output\n")
		serializedOutput, err := txOut.SerializeOutput()
		if err != nil {
			return nil, err
		}
		fmt.Printf("debug: serialized output %x\n", serializedOutput)
		serializedTx = append(serializedTx, serializedOutput...)
	}
	fmt.Printf("debug: %x\n", serializedTx)

	// Serialize the Locktime.
	LocktimeByte, err := hex.DecodeString(tx.Locktime)
	if err != nil {
		return nil, err
	}
	serializedTx = append(serializedTx, LocktimeByte...)

	return serializedTx, nil
}
