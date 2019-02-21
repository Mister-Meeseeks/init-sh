package initsh

import "fmt"
import "os"

type PathIngester interface {
	ingestPath(path string, info os.FileInfo) error
}

type IngestMultiplexer struct {
	subs []PathIngester
}

func (p IngestMultiplexer) ingestPath (path string, info os.FileInfo) error {
	for _, s := range p.subs {
		err := s.ingestPath(path, info)
		if err != nil {
			return err
		}
	}
	return nil
}

type IngestFilter interface {
	doIngest(path string, info os.FileInfo) (bool, error)
}

type IngestConditional struct {
	acceptor IngestFilter
	wrapped PathIngester
}

func (p IngestConditional) ingestPath (path string, info os.FileInfo) error {
	isAccept, err := p.acceptor.doIngest(path, info)
	if err != nil {
		return err
	} else if (isAccept) {
		return p.wrapped.ingestPath(path, info)
	} else {
		return nil
	}
}

type PathPrinter struct { }

func (p PathPrinter) ingestPath (path string, info os.FileInfo) error {
	fmt.Printf("Descent path %q - base %q - isDir= %q - Mode %o\n",
		path, info.Name(), info.IsDir(), info.Mode().Perm())
	return nil
}

