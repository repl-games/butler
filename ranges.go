package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/go-errors/errors"
	"github.com/itchio/butler/comm"
	"github.com/itchio/wharf/pools/blockpool"
	"github.com/itchio/wharf/pwr"
	"github.com/itchio/wharf/pwr/genie"
	"github.com/itchio/wharf/tlc"
	"github.com/itchio/wharf/wire"
)

type loggingSink struct {
	Container *tlc.Container
}

func (ls *loggingSink) Store(loc blockpool.BlockLocation, data []byte) error {
	comm.Logf("storing %v", loc)
	return nil
}

func (ls *loggingSink) GetContainer() *tlc.Container {
	return ls.Container
}

var _ blockpool.Sink = (*loggingSink)(nil)

func ranges(manifest string, patch string) {
	must(doRanges(manifest, patch))
}

func doRanges(manifest string, patch string) error {
	patchStats, err := os.Lstat(patch)
	if err != nil {
		return errors.Wrap(err, 1)
	}

	patchReader, err := os.Open(patch)
	if err != nil {
		return errors.Wrap(err, 1)
	}

	bigBlockSize := *appArgs.bigBlockSize

	g := &genie.Genie{
		BlockSize: bigBlockSize,
	}
	err = g.ParseHeader(patchReader)
	if err != nil {
		return errors.Wrap(err, 1)
	}

	targetContainer := g.TargetContainer
	sourceContainer := g.SourceContainer

	comm.Opf("Showing ranges for %s patch", humanize.IBytes(uint64(patchStats.Size())))
	comm.Statf("Old version: %s in %s", humanize.IBytes(uint64(targetContainer.Size)), targetContainer.Stats())
	comm.Statf("New version: %s in %s", humanize.IBytes(uint64(sourceContainer.Size)), sourceContainer.Stats())
	deltaOp := "+"
	if sourceContainer.Size < targetContainer.Size {
		deltaOp = "-"
	}
	delta := math.Abs(float64(sourceContainer.Size - targetContainer.Size))
	comm.Statf("Delta: %s%s (%s%.2f%%)", deltaOp, humanize.IBytes(uint64(delta)), deltaOp, delta/float64(targetContainer.Size)*100.0)
	comm.Log("")

	requiredOldBlocks := make([]map[int64]bool, len(targetContainer.Files))
	for i := 0; i < len(targetContainer.Files); i++ {
		requiredOldBlocks[i] = make(map[int64]bool)
	}

	requiredOldBlocksList := []blockpool.BlockLocation{}

	freshNewBlocks := make([]map[int64]bool, len(sourceContainer.Files))
	for i := 0; i < len(sourceContainer.Files); i++ {
		freshNewBlocks[i] = make(map[int64]bool)
	}

	comps := make(chan *genie.Composition)

	go func() {
		for comp := range comps {
			reuse := false

			if len(comp.Origins) == 1 {
				switch origin := comp.Origins[0].(type) {
				case *genie.BlockOrigin:
					if origin.Offset%bigBlockSize == 0 {
						reuse = true
					}
				}
			}

			if reuse {
				// comm.Logf("file %d, block %d is a re-use", comp.FileIndex, comp.BlockIndex)
			} else {
				comm.Logf("%s", comp.String())
				freshNewBlocks[comp.FileIndex][comp.BlockIndex] = true
				for _, anyOrigin := range comp.Origins {
					switch origin := anyOrigin.(type) {
					case *genie.BlockOrigin:
						blockStart := origin.Offset / bigBlockSize
						blockEnd := (origin.Offset + origin.Size + bigBlockSize - 1) / bigBlockSize
						for j := blockStart; j < blockEnd; j++ {
							if !requiredOldBlocks[origin.FileIndex][j] {
								requiredOldBlocksList = append(requiredOldBlocksList, blockpool.BlockLocation{FileIndex: origin.FileIndex, BlockIndex: j})
							}
							requiredOldBlocks[origin.FileIndex][j] = true
						}
					}
				}
			}
		}
	}()

	err = g.ParseContents(comps)
	if err != nil {
		return err
	}

	totalBlocks := 0
	partialBlocks := 0
	neededBlocks := 0
	neededBlockSize := int64(0)

	blockAddresses := make(blockpool.BlockAddressMap)

	for i, blockMap := range requiredOldBlocks {
		f := targetContainer.Files[i]
		fileNumBlocks := (f.Size + bigBlockSize - 1) / bigBlockSize
		for j := int64(0); j < fileNumBlocks; j++ {
			totalBlocks++
			if blockMap[j] {
				size := bigBlockSize
				if (j+1)*bigBlockSize > f.Size {
					partialBlocks++
					size = f.Size % bigBlockSize
				}
				neededBlockSize += size
				neededBlocks++
			}
		}
	}
	comm.Statf("Total old blocks: %d, needed: %d (of which %d are smaller than %s)", totalBlocks, neededBlocks, partialBlocks, humanize.IBytes(uint64(bigBlockSize)))
	comm.Statf("Needed block size: %s (%.2f%% of full old build size)", humanize.IBytes(uint64(neededBlockSize)), float64(neededBlockSize)/float64(targetContainer.Size)*100.0)

	freshBlocks := 0
	freshBlocksSize := int64(0)

	for i, blockMap := range freshNewBlocks {
		f := sourceContainer.Files[i]
		fileNumBlocks := (f.Size + bigBlockSize - 1) / bigBlockSize
		for j := int64(0); j < fileNumBlocks; j++ {
			if blockMap[j] {
				size := bigBlockSize
				if (j+1)*bigBlockSize > f.Size {
					size = f.Size % bigBlockSize
				}
				freshBlocksSize += size
				freshBlocks++
			}
		}
	}
	comm.Statf("Fresh blocks: %d, %s total", freshBlocks, humanize.IBytes(uint64(freshBlocksSize)))
	comm.Statf("Required old blocks order: %v", requiredOldBlocksList)

	pathToFileIndex := make(map[string]int)
	for i, f := range targetContainer.Files {
		pathToFileIndex[f.Path] = i
	}

	manifestReader, err := os.Open(manifest)
	if err != nil {
		return err
	}

	rawManWire := wire.NewReadContext(manifestReader)
	err = rawManWire.ExpectMagic(pwr.ManifestMagic)
	if err != nil {
		return err
	}

	mh := &pwr.ManifestHeader{}
	err = rawManWire.ReadMessage(mh)
	if err != nil {
		return err
	}

	if mh.Algorithm != pwr.HashAlgorithm_SHAKE128_32 {
		return fmt.Errorf("Manifest has unsupported hash algorithm %d, expected %d", mh.Algorithm, pwr.HashAlgorithm_SHAKE128_32)
	}

	manWire, err := pwr.DecompressWire(rawManWire, mh.GetCompression())
	if err != nil {
		return err
	}

	manContainer := &tlc.Container{}
	err = manWire.ReadMessage(manContainer)
	if err != nil {
		return err
	}

	sh := &pwr.SyncHeader{}
	mbh := &pwr.ManifestBlockHash{}

	for i, f := range manContainer.Files {
		sh.Reset()
		err = manWire.ReadMessage(sh)
		if err != nil {
			return err
		}

		if int64(i) != sh.FileIndex {
			return fmt.Errorf("manifest format error: expected file %d, got %d", i, sh.FileIndex)
		}

		fileIndex := int64(pathToFileIndex[f.Path])
		numBlocks := int64(math.Ceil(float64(f.Size) / float64(bigBlockSize)))
		for j := int64(0); j < numBlocks; j++ {
			mbh.Reset()
			err = manWire.ReadMessage(mbh)
			if err != nil {
				return err
			}

			size := bigBlockSize
			if (j+1)*bigBlockSize > f.Size {
				size = f.Size % bigBlockSize
			}

			address := fmt.Sprintf("shake128-32/%x/%d", mbh.Hash, size)
			blockAddresses.Set(blockpool.BlockLocation{FileIndex: fileIndex, BlockIndex: j}, address)
		}
	}

	var source blockpool.Source

	source = &blockpool.DiskSource{
		BasePath:       "./blocks",
		BlockAddresses: blockAddresses,

		Container: targetContainer,
	}

	if *rangesArgs.latency > 0 {
		source = &blockpool.DelayedSource{
			Latency: time.Duration(*rangesArgs.latency) * time.Millisecond,
			Source:  source,
		}
	}

	targetPool := &blockpool.BlockPool{
		Container: targetContainer,
		BlockSize: bigBlockSize,

		Upstream: source,

		Consumer: comm.NewStateConsumer(),
	}

	actx := &pwr.ApplyContext{
		Consumer:   comm.NewStateConsumer(),
		TargetPool: targetPool,
	}

	actx.OutputPool = &blockpool.BlockPool{
		Container: sourceContainer,
		BlockSize: bigBlockSize,

		Downstream: &loggingSink{},

		Consumer: comm.NewStateConsumer(),
	}

	_, err = patchReader.Seek(0, os.SEEK_SET)
	if err != nil {
		return err
	}

	startTime := time.Now()

	comm.StartProgress()
	err = actx.ApplyPatch(patchReader)
	if err != nil {
		return err
	}
	comm.EndProgress()

	totalTime := time.Since(startTime)
	comm.Statf("Processed in %s (%s/s)", totalTime, humanize.IBytes(uint64(float64(targetContainer.Size)/totalTime.Seconds())))

	return nil
}