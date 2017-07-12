package raft

import (
  "io/ioutil"
  "log"
  "os"
  "strings"
  "sync"

  "github.com/golang/protobuf/proto"
)

const (
  MetadataFileName = "raft-metadata"
)

type Persister struct {
  logDir string
  metadataFilePath string
  logDirMetadata *RaftLogDirMetadata

  currentTerm int64
  commitIndex int64
  mu synx.Mutex
}

func (p *Persister) init() {
  _, err := os.Stat(p.logDir)
  if err != nil {
    if os.IsNotExist(err) {
      err = os.Mkdir(p.logDir, 0755)
      if err != nil {
        log.Fatal("Error at creating raft dir:", err)
      }
    } else {
      log.Fatal("Error at opening raft dir:", err)
    }
  }

  _, err = os.Stat(p.metadataFilePath)
  if err != nil {
    if os.IsNotExist(err) {
      p.WriteLogDirMetadata()
    } else {
      log.Fatal(err)
    }
  }

  p.ReadLogDirMetadata()
}

func (p *Persister) WriteLogDirMetadata() {
  out, err := proto.Marshal(p.logDirMetadata)
  if err != nil {
    log.Fatal(err)
  }

  if err := ioutil.WriteFile(p.metadataFilePath, out, 0644); err != nil {
    log.Fatal(err)
  }
}

func (p *Persister) ReadLogDirMetadata() {
  in, err := ioutil.ReadFile(p.metadataFilePath) 
  if err != nil {
    log.Fatal(err)
  }

  p.logDirMetadata = &RaftLogDirMetadata{}
  if err := proto.Unmarshal(in, p.logDirMetadata); err != nil {
    log.Fatal(err)
  }
}

func (p *Persister) Close() {
}

func NewPersister(logDir string) *Persister {
  p := &Persister{}
  p.logDir = logDir
  p.metadataFilePath = strings.Join([]string{p.logDir, MetadataFileName}, "/")
  p.logDirMetadata = &RaftLogDirMetadata{}
  p.init()
  return p
}
