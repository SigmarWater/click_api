package scripts

import (
	"context"
	"github.com/ClickHouse/clickhouse-go/v2"
	"log"
	"math/rand"
	"time"
)

var events = []string{
	"login",
	"logout",
	"purchase",
	"view_product",
	"add_to_cart",
}

func InsertBatchToDb(ctx context.Context, conn clickhouse.Conn) error {
	query := "INSERT INTO events (id, user_id, event_name, create_at)"

	batch, err := conn.PrepareBatch(ctx, query)
	if err != nil {
		log.Printf("failed create batch: %v\n", err)
		return err
	}

	for i := 1; i <= 100000; i++ {
		err := batch.Append(
			uint64(i),
			uint64(rand.Intn(10000)),
			events[rand.Intn(len(events))],
			time.Now(),
		)

		if err != nil {
			log.Printf("failed append into batch: %v\n", err)
			return err
		}
	}

	if err := batch.Send(); err != nil {
		log.Printf("failed insert into db: %v\n", err)
		return err
	}

	return nil
}
