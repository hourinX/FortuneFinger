package snowflake

import (
	"fmt"
	"sync"
	"time"
)

const (
	customEpoch    int64 = 1704067200000
	machineBits          = 10
	sequenceBits         = 12
	maxMachineID         = (1 << machineBits) - 1
	maxSequence          = (1 << sequenceBits) - 1
	timestampShift       = machineBits + sequenceBits
	machineShift         = sequenceBits
	maxWait              = 200 * time.Millisecond
)

type Generator struct {
	mu            sync.Mutex
	machineID     int64
	lastTimestamp int64
	sequence      int64
}

func New(machineID int64) (*Generator, error) {
	if machineID < 0 || machineID > maxMachineID {
		return nil, fmt.Errorf("snowflake: MachineID必须在0到%d之间，当前值为「%d」", maxMachineID, machineID)
	}
	return &Generator{machineID: machineID}, nil
}

func (g *Generator) Next() (uint64, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now().UnixMilli()
	if now < g.lastTimestamp {
		wait := time.Duration(g.lastTimestamp-now) * time.Millisecond
		if wait > maxWait {
			return 0, fmt.Errorf("snowflake: 时钟回拨时间「%s」超过允许范围", wait)
		}
		time.Sleep(wait)
		now = time.Now().UnixMilli()
		if now < g.lastTimestamp {
			return 0, fmt.Errorf("snowflake: 等待后系统时钟仍未恢复")
		}
	}

	if now == g.lastTimestamp {
		g.sequence = (g.sequence + 1) & maxSequence
		if g.sequence == 0 {
			for now <= g.lastTimestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		g.sequence = 0
	}
	g.lastTimestamp = now

	if now < customEpoch {
		return 0, fmt.Errorf("snowflake: 系统时间早于自定义纪元")
	}
	timestamp := uint64(now - customEpoch)
	return (timestamp << timestampShift) | (uint64(g.machineID) << machineShift) | uint64(g.sequence), nil
}
