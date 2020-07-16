// Copyright (c) 2020, The Facette Authors
//
// Licensed under the terms of the BSD 3-Clause License; a copy of the license
// is available at: https://opensource.org/licenses/BSD-3-Clause

package api

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"facette.io/facette/pkg/api"
	"facette.io/facette/pkg/store"
)

// nolint:errcheck,gosec
func (h handler) DumpStore(rw http.ResponseWriter, r *http.Request) {
	f, err := ioutil.TempFile("", "facette-dump")
	if err != nil {
		h.WriteError(rw, err)
		return
	}

	defer os.Remove(f.Name())

	err = dumpStore(h.store, f)
	if err != nil {
		h.WriteError(rw, err)
		return
	}

	rw.Header().Add("Content-Type", "application/gzip")
	f.Seek(0, 0)
	io.Copy(rw, f)
}

func dumpStore(store *store.Store, w io.Writer) error {
	gw := gzip.NewWriter(w)
	defer gw.Close() // nolint:errcheck

	tw := tar.NewWriter(gw)
	defer tw.Close() // nolint:errcheck

	ch := make(chan api.Object)
	errCh := make(chan error)

	go func(errCh chan<- error) {
		for obj := range ch {
			err := dumpObject(tw, obj)
			if err != nil {
				errCh <- err
				break
			}
		}

		close(errCh)
	}(errCh)

	err := store.Dump(ch)
	if err != nil {
		return err
	}

	return <-errCh
}

func dumpObject(tw *tar.Writer, obj api.Object) error {
	var prefix string

	switch obj.(type) {
	case *api.Chart:
		prefix = "charts/"

	case *api.Dashboard:
		prefix = "dashboards/"

	case *api.Provider:
		prefix = "providers/"
	}

	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	meta := obj.GetMeta()

	err = tw.WriteHeader(&tar.Header{
		Name:    prefix + meta.ID + ".json",
		Size:    int64(len(data)),
		Mode:    0644,
		ModTime: meta.ModifiedAt,
	})
	if err != nil {
		return err
	}

	_, err = tw.Write(data)

	return err
}

func (h handler) RestoreStore(rw http.ResponseWriter, r *http.Request) {
	mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mt != "application/gzip" {
		h.WriteError(rw, api.ErrUnsupportedType)
		return
	}

	// TODO: handle dump errors

	err = restoreStore(h.store, r.Body)
	if err != nil {
		h.WriteError(rw, err)
		return
	}

	// Reset poller jobs based on providers now being present into the back-end
	// storage
	h.poller.InsertFromStore(h.store)

	rw.WriteHeader(http.StatusNoContent)
}

func restoreStore(store *store.Store, r io.Reader) error {
	ch := make(chan api.Object)
	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error)

	go func() {
		err := readDump(r, ch)
		if err != nil {
			cancel()
			errCh <- err
		}

		close(errCh)
	}()

	err := store.Restore(ctx, ch)
	if err != nil {
		return err
	}

	return <-errCh
}

func readDump(r io.Reader, ch chan<- api.Object) error {
	defer close(ch)

	gr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}

	defer gr.Close() // nolint:errcheck

	tr := tar.NewReader(gr)

	var h *tar.Header

	for {
		h, err = tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if h.Typeflag != tar.TypeReg && h.Typeflag != tar.TypeRegA {
			continue
		}

		var (
			obj  api.Object
			data []byte
		)

		switch filepath.Dir(h.Name) {
		case "charts":
			obj = &api.Chart{}

		case "dashboards":
			obj = &api.Dashboard{}

		case "providers":
			obj = &api.Provider{}

		default:
			continue
		}

		data, err = ioutil.ReadAll(tr)
		if err != nil {
			return err
		}

		err = json.Unmarshal(data, obj)
		if err != nil {
			return err
		}

		ch <- obj
	}

	return nil
}
