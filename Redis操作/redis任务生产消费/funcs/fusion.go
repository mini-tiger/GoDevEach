package funcs

import (
	"context"
	"encoding/json"
	"fmt"
	mathrand "math/rand"
	"redis任务生产消费/g"
	"redis任务生产消费/modules"
	"strconv"
	"time"
)

var ctx = context.Background()

func PushMission(m *modules.Mission) {
	d, _ := json.Marshal(m)
	if err := g.RDB.LPush(ctx, "queue", d).Err(); err != nil {
		panic(err)
	}
}

func PullMission() []string {
	// use `rdb.BLPop(0, "queue")` for infinite waiting time
	result, err := g.RDB.BRPop(ctx, 1*time.Second, "queue").Result()
	if err != nil {
		panic(err)
	}

	return result
}

func Demo1() {
	for i := 0; i < 10; i++ {
		m := &modules.Mission{
			Data: strconv.Itoa(i),
			Time: time.Now().Format("2006-01-02 15:04:05"),
		}
		PushMission(m)
		fmt.Println(PullMission())
	}
}

func Demo2() {
	go func() {
		for {
			m := &modules.Mission{
				Data: strconv.Itoa(mathrand.Int()),
				Time: time.Now().Format("2006-01-02 15:04:05"),
			}
			PushMission(m)
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for {
			fmt.Println(PullMission())
			time.Sleep(1 * time.Second)
		}

	}()
}
