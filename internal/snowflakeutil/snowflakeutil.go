package snowflake

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"time"
)

var node *snowflake.Node

func Init(startTime string, machineID int64) {
	st, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		panic(fmt.Errorf("snowflake starttime parse failed: %s \n", err))
	}
	snowflake.Epoch = st.UnixMilli()
	node, err = snowflake.NewNode(machineID)
	if err != nil {
		panic(fmt.Errorf("snowflake NewNode failed: %s \n", err))
	}
}

func GenID() int64 {
	return node.Generate().Int64()
}
