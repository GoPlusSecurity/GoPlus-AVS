package avs_events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	proxy "goplus/avs/contracts/bindings/TransparentUpgradeableProxy"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ScanResult struct {
	ContractName    string
	ContractAddress string
	BlockNumber     uint64
	LogIndex        uint
	EventName       string
	Data            string
}

type ContractInfo struct {
	Name string
	ABI  *abi.ABI
}

func decodeTopicValue(t abi.Type, topic common.Hash) (interface{}, error) {
	switch t.T {
	case abi.AddressTy:
		return common.HexToAddress(topic.Hex()), nil
	case abi.IntTy, abi.UintTy:
		return decodeBigInt(topic, t.Size), nil
	case abi.BoolTy:
		return topic.Big().Cmp(big.NewInt(0)) != 0, nil
	case abi.StringTy, abi.BytesTy, abi.SliceTy, abi.ArrayTy:
		// 这些类型通常不会被索引，或者如果被索引，会被哈希
		// 在这种情况下，我们无法恢复原始值，只能返回哈希
		return topic.Hex(), nil
	default:
		return nil, fmt.Errorf("unsupported indexed type: %v", t.String())
	}
}

func decodeBigInt(topic common.Hash, size int) *big.Int {
	b := topic.Big()
	if size <= 64 {
		// 对于小于等于64位的整数，检查是否应该是负数
		if b.Bit(size-1) == 1 {
			// 是负数，进行二的补码转换
			b.Sub(b, new(big.Int).Lsh(big.NewInt(1), uint(size)))
		}
	}
	return b
}

// ScanEvents 扫描 `fromBlock` 到  `toBlock` 区块范围内，由一组指定地址的合约 `contracts` 发出的事件
// 按事件发生的时间正序排列返回每一个事件。
// 如果 30 秒还没获取到结果就超时退出。
func ScanEvents(ethClient *ethclient.Client, fromBlock *big.Int, toBlock *big.Int, contracts map[common.Address]ContractInfo) ([]ScanResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var addresses []common.Address
	for k := range contracts {
		addresses = append(addresses, k)
	}
	query := ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: addresses,
	}

	logs, err := ethClient.FilterLogs(ctx, query)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		} else {
			log.Fatalf("Failed to filterLogs: %v", err)
		}
	}

	tup, err := proxy.TransparentUpgradeableProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	var resultList []ScanResult
	for _, l := range logs {
		result := ScanResult{
			ContractAddress: l.Address.Hex(),
			BlockNumber:     l.BlockNumber,
			LogIndex:        l.Index,
		}

		info, exists := contracts[l.Address]
		if !exists {
			log.Fatalf("Unknown contract address: %#v", result)
		}
		result.ContractName = info.Name

		// 先按 TUP 合约解析事件
		var abiObj *abi.ABI
		var evt *abi.Event
		if evt2, err := tup.EventByID(l.Topics[0]); err == nil {
			evt = evt2
			abiObj = tup
		} else {
			if evt3, err := info.ABI.EventByID(l.Topics[0]); err == nil {
				evt = evt3
				abiObj = info.ABI
			} else {
				log.Fatalf("Unknown event: %s, result: %#v", l.Topics[0].Hex(), result)
			}
		}

		result.EventName = evt.Name
		data := make(map[string]interface{})
		err = abiObj.UnpackIntoMap(data, evt.Name, l.Data)
		if err != nil {
			log.Fatalf("Failed to unpack event: %#v", result)
		}

		for i := 1; i < len(l.Topics) && i <= len(evt.Inputs); i++ {
			input := evt.Inputs[i-1]
			if input.Indexed {
				value, err := decodeTopicValue(input.Type, l.Topics[i])
				if err == nil {
					data[input.Name] = value
				} else {
					data[input.Name] = l.Topics[i].Hex() // 回退到十六进制
				}
			}
		}
		dataBytes, err := json.Marshal(data)
		if err != nil {
			log.Fatalf("Failed to json dump: %#v", data)
		}
		result.Data = string(dataBytes)
		resultList = append(resultList, result)
		log.Printf("%#v", result)
	}
	return resultList, nil
}
