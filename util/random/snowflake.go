package rd

import (
	"sync"

	"github.com/bwmarrin/snowflake"
)

var snowflaker = make(map[int64]*snowflake.Node)
var rwlk sync.RWMutex

func getNode(node int64) (*snowflake.Node, error) {
	n := snowflaker[node]
	if n == nil {
		err := setNode(node)
		if err != nil {
			return nil, err
		}
	}
	return snowflaker[node], nil
}

func setNode(node int64) error {
	rwlk.Lock()
	defer func() {
		rwlk.Unlock()
	}()
	if snowflaker[node] == nil {
		n, err := snowflake.NewNode(node)
		if err != nil {
			return err
		}
		snowflaker[node] = n
	}
	return nil
}

func SnowFlake(node int64) (int64, error) {
	n, err := getNode(node)
	if err != nil {
		return 0, err
	}
	return n.Generate().Int64(), nil
}
