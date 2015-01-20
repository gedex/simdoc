package upload

import (
	"errors"
	"fmt"
	"sync"

	"github.com/gedex/simdoc/pkg/model"
)

type FileResult struct {
	*File
	Versions map[string]*FileVersion `json:"versions"`
}

type FileVersion struct {
	*model.DocumentFileVersion
	Error error `json:"error"`
}

type ProcessManager interface {
	Add(p Processor)
	Run(ap AfterProcessFn) *FileResult
}

type Processor interface {
	Process(f *File) (*File, error)
	GetName() string
	GetSource() string
	CanProcess(baseMime string) bool
}

var SourceOriginal = ":original:"

type processManager struct {
	mu  sync.RWMutex
	src *File                    // Original source
	pe  map[string]*processEntry // key is process name provided by processor implementor
}

type processEntry struct {
	src   *File
	proc  Processor
	downs []*processEntry // Downstreams
}

type AfterProcessFn func(out *File, err error) (*File, error)

func ProcessFile(f *File, ap AfterProcessFn, procs ...Processor) *FileResult {
	pm := new(processManager)
	pm.src = f
	pm.pe = make(map[string]*processEntry, 0)

	for _, p := range procs {
		pm.Add(p)
	}

	return pm.Run(ap)
}

func (pm *processManager) Add(p Processor) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, exists := pm.pe[p.GetName()]; exists {
		panic("upload: processor name already exists")
	}
	pname := p.GetName()
	psrc := p.GetSource()

	pe := new(processEntry)
	if psrc == SourceOriginal {
		pe.src = pm.src
	} else {
		// If source processor doesn't exists panic, otherwise set source processor
		// to upstream ups and set current processor as downstream downs of
		// the source processor.
		if u, ok := pm.pe[psrc]; !ok {
			panic(fmt.Sprintf("upload: source processor %s does not exists", psrc))
		} else {
			u.downs = append(u.downs, pe)
		}
	}
	pe.proc = p

	// Adds process entry to manager.
	pm.pe[pname] = pe
}

func (pm *processManager) Run(ap AfterProcessFn) *FileResult {
	ver := make(map[string]*FileVersion, len(pm.pe))
	fr := &FileResult{pm.src, ver}
	for pname, pe := range pm.pe {
		if pe.src == nil {
			continue
		}

		if !pe.proc.CanProcess(pe.src.Type) {
			continue
		}

		fr.Versions[pname] = &FileVersion{new(model.DocumentFileVersion), nil}

		res, err := pe.proc.Process(pe.src)

		res, err = ap(res, err)

		if err != nil {
			fr.Versions[pname].Error = errors.New(fmt.Sprintf("upload.processor.Run %s.Process error: %s", pname, err))
			continue
		}

		fr.Versions[pname].Filepath = res.Filepath
		fr.Versions[pname].URL = res.URL
		fr.Versions[pname].Meta = &model.DocumentFileMeta{
			Type: res.Type,
			Mime: res.Mime,
			Size: res.Size,
		}

		// Supplies src to downstreams.
		for _, d := range pe.downs {
			d.src = res
		}
	}

	return fr
}
