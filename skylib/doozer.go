package skylib

import (
	"bytes"
	"github.com/4ad/doozer"
	"log"
  "sync"
)

type DoozerServer struct {
	Key      string
	Id       int
	Addr     string
}

type DoozerConnection struct {
	Connection *doozer.Conn
	Log        *log.Logger
	Discover   bool
  Uri        string
  BootUri    string

  connectionMutex sync.Mutex
}

func (d *DoozerConnection) Connect() {
	if d.Uri == "" && d.BootUri == "" {
		d.Log.Panic("Must supply a Doozer server or a BootUri to connect to")
	}

  success, err := d.dial()

  if success == false {
    d.Log.Panic("Failed to connect to any of the supplied Doozer Servers: " + err.Error())
  }
}

func (d *DoozerConnection) dial()  (bool, error) {
  d.connectionMutex.Lock()
  defer d.connectionMutex.Unlock()

	var err error

  if d.Uri != "" {
    d.Connection, err = doozer.DialUri(d.Uri, d.BootUri)
  } else {
    d.Connection, err = doozer.DialUri(d.BootUri, "")
  }

	if err != nil {
		return false, err
	}

  d.Log.Println("Connected to Doozer")

  return true, nil
}

func (d *DoozerConnection) GetCurrentRevision() (rev int64) {
	defer func() {
		if err := recover(); err != nil {
      d.recoverFromError(err)

      rev = d.GetCurrentRevision()
		}
	}()

	revision, err := d.Connection.Rev()

	if err != nil {
		d.Log.Panic(err.Error())
	}

	return revision
}

func (d *DoozerConnection) Set(file string, rev int64, body []byte) (newRev int64, err error) {
	defer func() {
		if err := recover(); err != nil {
      d.recoverFromError(err)

      newRev, err = d.Set(file, rev, body)
		}
	}()

	return d.Connection.Set(file, rev, body)
}

func (d *DoozerConnection) Del(path string, rev int64) (err error) {
	defer func() {
		if err := recover(); err != nil {
      d.recoverFromError(err)

      err = d.Del(path, rev)
		}
	}()

	return d.Connection.Del(path, rev)
}

func (d *DoozerConnection) Get(file string, rev *int64) (data []byte, revision int64, err error) {
	defer func() {
		if err := recover(); err != nil {
      d.recoverFromError(err)

      data, revision, err = d.Get(file, rev)
		}
	}()

	return d.Connection.Get(file, rev)
}

func (d *DoozerConnection) Rev() (rev int64, err error) {
	defer func() {
		if err := recover(); err != nil {
      d.recoverFromError(err)

      rev, err = d.Rev()
		}
	}()

	return d.Connection.Rev()
}

func (d *DoozerConnection) recoverFromError(err interface{}){
  if err == "EOF" {
    d.Log.Println("Lost connection to Doozer: Reconnecting...")

    success, _ := d.dial()

    if success == true {
      return
    }

    // If we made it here we didn't find a server
    d.Log.Panic("Unable to find a Doozer instance to connect to")

  } else {
    // Don't know how to handle, go ahead and panic
    d.Log.Panic(err)
  }
}

func (d *DoozerConnection) getDoozerServer(key string) (*DoozerServer){
  rev := d.GetCurrentRevision()
  data, _, err := d.Get("/ctl/node/" + key + "/addr", &rev)
  buf := bytes.NewBuffer(data)

  if err == nil {
    return &DoozerServer {
      Addr: buf.String(),
      Key:  key,
    }
  }

  return nil
}

func basename(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[i+1:]
		}
	}
	return path
}
