package archive

import (
	"os"
	"path/filepath"

	"github.com/itchio/savior"

	"github.com/go-errors/errors"
	"github.com/itchio/butler/butlerd"
	"github.com/itchio/butler/installer"
	"github.com/itchio/butler/installer/archive/intervalsaveconsumer"
	"github.com/itchio/butler/installer/bfs"
)

func (m *Manager) Install(params *installer.InstallParams) (*installer.InstallResult, error) {
	consumer := params.Consumer

	var res = installer.InstallResult{
		Files: []string{},
	}

	f := params.File

	archiveInfo := params.InstallerInfo.ArchiveInfo
	if archiveInfo.Features.ResumeSupport == savior.ResumeSupportNone {
		consumer.Infof("Forcing local for %s", archiveInfo.Features)
		localFile, err := installer.AsLocalFile(f)
		if err != nil {
			return nil, errors.Wrap(err, 0)
		}
		f = localFile
	}

	ex, err := archiveInfo.GetExtractor(f, consumer)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	ex.SetConsumer(consumer)

	statePath := filepath.Join(params.StageFolderPath, "install-state.dat")
	sc := intervalsaveconsumer.New(statePath, intervalsaveconsumer.DefaultInterval, consumer, params.Context)
	ex.SetSaveConsumer(sc)

	cancelled := false
	defer func() {
		if !cancelled {
			consumer.Infof("Clearing archive install state")
			os.Remove(statePath)
		}
	}()

	checkpoint, err := sc.Load()
	if err != nil {
		consumer.Warnf("Could not load checkpoint: %s", err.Error())
	}

	sink := &savior.FolderSink{
		Directory: params.InstallFolderPath,
		Consumer:  consumer,
	}

	aRes, err := ex.Resume(checkpoint, sink)
	if err != nil {
		if errors.Is(err, savior.ErrStop) {
			cancelled = true
			return nil, &butlerd.ErrCancelled{}
		}
		return nil, errors.Wrap(err, 0)
	}

	err = sink.Close()
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	aRes, err = archiveInfo.ApplyStageTwo(consumer, aRes, params.InstallFolderPath)
	if err != nil {
		if errors.Is(err, savior.ErrStop) {
			cancelled = true
			return nil, &butlerd.ErrCancelled{}
		}
		return nil, errors.Wrap(err, 0)
	}

	for _, entry := range aRes.Entries {
		res.Files = append(res.Files, entry.CanonicalPath)
	}

	err = bfs.BustGhosts(&bfs.BustGhostsParams{
		Folder:   params.InstallFolderPath,
		NewFiles: res.Files,
		Receipt:  params.ReceiptIn,

		Consumer: params.Consumer,
	})
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return &res, nil
}
