package storable

import (
	"database/sql"
	"strconv"
	"time"

	types2 "github.com/kekDAO/kekBackend/types"

	"github.com/kekDAO/kekBackend/utils"

	"github.com/lib/pq"

	"github.com/sirupsen/logrus"

	"github.com/alethio/web3-go/types"
)

var log = logrus.WithField("module", "storable(block)")

type Block struct {
	RawBlock             types.Block
	Number               int64
	BlockHash            string
	ParentBlockHash      string
	BlockCreationTime    types2.DatetimeToJSONUnix
	BlockGasLimit        string
	BlockGasUsed         string
	BlockDifficulty      string
	TotalBlockDifficulty string
	BlockExtraData       types2.ByteArray
	BlockMixHash         types2.ByteArray
	BlockNonce           types2.ByteArray
	BlockSize            int64
	BlockLogsBloom       types2.ByteArray
	IncludesUncle        types2.JSONStringArray
	HasBeneficiary       types2.ByteArray
	HasReceiptsTrie      types2.ByteArray
	HasTxTrie            types2.ByteArray
	Sha3Uncles           types2.ByteArray
	NumberOfUncles       int32
	NumberOfTxs          int32
}

func NewStorableBlock(block types.Block) *Block {
	return &Block{RawBlock: block}
}

func (sb *Block) ToDB(tx *sql.Tx) error {
	log.Trace("storing block")
	start := time.Now()
	defer func() { log.WithField("duration", time.Since(start)).Debug("done storing block") }()

	err := sb.enhance()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(pq.CopyIn("blocks", "number", "block_hash", "parent_block_hash", "block_creation_time", "block_gas_limit", "block_gas_used", "block_difficulty", "total_block_difficulty", "block_extra_data", "block_mix_hash", "block_nonce", "block_size", "block_logs_bloom", "includes_uncle", "has_beneficiary", "has_receipts_trie", "has_tx_trie", "sha3_uncles", "number_of_uncles", "number_of_txs"))
	if err != nil {
		return err
	}

	_, err = stmt.Exec(sb.Number, sb.BlockHash, sb.ParentBlockHash, sb.BlockCreationTime, sb.BlockGasLimit, sb.BlockGasUsed, sb.BlockDifficulty, sb.TotalBlockDifficulty, sb.BlockExtraData, sb.BlockMixHash, sb.BlockNonce, sb.BlockSize, sb.BlockLogsBloom, sb.IncludesUncle, sb.HasBeneficiary, sb.HasReceiptsTrie, sb.HasTxTrie, sb.Sha3Uncles, sb.NumberOfUncles, sb.NumberOfTxs)
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return nil
}

// enhance processes all the raw data of a block and does the necessary transformations,
// resulting in an object that's ready for inserting into the database
func (sb *Block) enhance() error {
	b := sb.RawBlock

	if b.Miner == "" {
		b.Miner = b.Author
	}

	sb.BlockHash = utils.Trim0x(b.Hash)
	sb.ParentBlockHash = utils.Trim0x(b.ParentHash)
	sb.BlockExtraData = types2.ByteArray(utils.Trim0x(b.ExtraData))
	sb.BlockMixHash = types2.ByteArray(utils.Trim0x(b.MixHash))
	sb.BlockNonce = types2.ByteArray(utils.Trim0x(b.Nonce))
	sb.BlockLogsBloom = types2.ByteArray(utils.Trim0x(b.LogsBloom))
	sb.IncludesUncle = types2.JSONStringArray(b.Uncles)
	sb.HasBeneficiary = types2.ByteArray(utils.Trim0x(b.Miner))
	sb.HasReceiptsTrie = types2.ByteArray(utils.Trim0x(b.ReceiptsRoot))
	sb.HasTxTrie = types2.ByteArray(utils.Trim0x(b.TransactionsRoot))
	sb.Sha3Uncles = types2.ByteArray(utils.Trim0x(b.Sha3Uncles))

	// -- ints
	number, err := strconv.ParseInt(b.Number, 0, 64)
	if err != nil {
		log.Error(err)
		return err
	}
	sb.Number = number

	size, err := strconv.ParseInt(b.Size, 0, 64)
	if err != nil {
		log.Error(err)
		return err
	}
	sb.BlockSize = size

	// --hexes
	gasLimit, err := utils.HexStrToBigIntStr(b.GasLimit)
	if err != nil {
		log.Error(err)
		return err
	}
	sb.BlockGasLimit = gasLimit

	gasUsed, err := utils.HexStrToBigIntStr(b.GasUsed)
	if err != nil {
		log.Error(err)
		return err
	}
	sb.BlockGasUsed = gasUsed

	difficulty, err := utils.HexStrToBigIntStr(b.Difficulty)
	if err != nil {
		log.Error(err)
		return err
	}
	sb.BlockDifficulty = difficulty

	totalDifficulty, err := utils.HexStrToBigIntStr(b.TotalDifficulty)
	if err != nil {
		log.Error(err)
		return err
	}
	sb.TotalBlockDifficulty = totalDifficulty

	// --timestamp
	timestamp, err := strconv.ParseInt(b.Timestamp, 0, 64)
	if err != nil {
		log.Error(err)
		return err
	}
	sb.BlockCreationTime = types2.DatetimeToJSONUnix(time.Unix(timestamp, 0))

	// -- computed
	sb.NumberOfTxs = int32(len(b.Transactions))
	sb.NumberOfUncles = int32(len(b.Uncles))

	return nil
}
