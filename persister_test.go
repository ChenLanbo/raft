package raft

import (
  "testing"
)

func TestInit(t *testing.T) {
  p := NewPersister("/tmp/raft_persister")
  defer p.Close()
}
