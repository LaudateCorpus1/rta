package collector

import (
	"context"
	"github.com/google/uuid"
	"github.com/viant/afs"
	"github.com/viant/afs/file"
	"github.com/viant/rta/collector/config"
	"github.com/viant/rta/shared"
	tconfig "github.com/viant/tapper/config"
	"github.com/viant/tapper/log"
	"github.com/viant/toolbox"
	"time"
)

type Batch struct {
	ID string
	*tconfig.Stream
	Accumulator map[interface{}]interface{}
	Started     time.Time
	Count       int32
	logger      *log.Logger
	PendingURL  string
}

func (b Batch) IsActive(batch *config.Batch) bool {
	return toolbox.AsInt(len(b.Accumulator)) < batch.MaxElements && time.Now().Sub(b.Started) < batch.MaxDuration()

}

func NewBatch(stream *tconfig.Stream, fs afs.Service) (*Batch, error) {

	UUID := uuid.New()
	URL := stream.URL + "_" + UUID.String()
	pendingURL := URL + shared.PendingSuffix

	batchSteam := &tconfig.Stream{
		FlushMod: stream.FlushMod,
		Codec:    stream.Codec,
		URL:      URL,
	}
	if err := fs.Create(context.Background(), pendingURL, file.DefaultFileOsMode, false); err != nil {
		return nil, err
	}
	logger, err := log.New(batchSteam, "", fs)
	if err != nil {
		return nil, err
	}
	return &Batch{
		PendingURL:  pendingURL,
		ID:          UUID.String(),
		Stream:      batchSteam,
		Accumulator: make(map[interface{}]interface{}),
		Started:     time.Now(),
		Count:       0,
		logger:      logger,
	}, nil
}
