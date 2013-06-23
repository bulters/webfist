package main

import (
  "encoding/json"
  "log"
  "net/http"
  "strings"

  "github.com/bradfitz/webfist"
)

func (s *server) Lookup(w http.ResponseWriter, r *http.Request) {
  if r.ParseForm() != nil {
    http.Error(w, "Bad request", http.StatusBadRequest)
    return
  }
  resource := r.Form.Get("resource")
  if resource == "" {
    http.Error(w, "'resource' missing", http.StatusBadRequest)
    return
  }
  emailLikeId := strings.TrimPrefix(resource, "acct:")
  if (emailLikeId == resource) {
    http.Error(w, "'resource' must start with 'acct:'", http.StatusBadRequest)
    return
  }

  foundData, err := s.lookup.WebFinger(emailLikeId)
  if err != nil {
    log.Printf("Error looking up resource: %s -- %v", emailLikeId, err)
    http.Error(w, "Error doing lookup", http.StatusInternalServerError)
    return
  }
  if foundData == nil {
    log.Printf("Not found: %s", emailLikeId)
    http.NotFound(w, r)
    return
  }
  b, err := json.Marshal(foundData)
  if err != nil {
    log.Printf("Bad data for resource: %s -- %v", emailLikeId, err)
    http.Error(w, "Bad data from lookup", http.StatusInternalServerError)
    return
  }

  log.Printf("Found user %s -- %+v", emailLikeId, foundData)
  w.Write(b)
}

type emailLookup struct {
  storage webfist.Storage
}

func (l *emailLookup) WebFinger(emailAddr string) (*webfist.WebFingerResponse, error) {
  resp := &webfist.WebFingerResponse{
  	Subject: "foo@bar.com",
  }
  email := webfist.NewEmailAddr(emailAddr)
  emailList, _ := l.storage.Emails(email)
  // TODO: Make this actuall parse emails
  log.Printf("Email list: %v", emailList)
  return resp, nil
}

func NewLookup(storage webfist.Storage) webfist.Lookup {
  return &emailLookup{
    storage: storage,
  }
}
