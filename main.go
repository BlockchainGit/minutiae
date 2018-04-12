package main

import (
  "fmt"
  "log"
  "net/http"
  "os"

  "cloud.google.com/go/datastore"
  "golang.org/x/net/context"
  "google.golang.org/appengine"
)

type note struct {
    Addr    string
    Cost    uint32
    Key     string
    Status  string
    Value   uint32
}

var datastoreClient *datastore.Client

func main() {
  ctx := context.Background()
  projectID := os.Getenv("GCLOUD_DATASET_ID")
  var err error
  datastoreClient, err = datastore.NewClient(ctx, projectID)
  if err != nil {
    log.Fatal(err)
  }

  http.HandleFunc("/", handle)
  appengine.Main()
}

func handle(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
    http.NotFound(w, r)
    return
  }
  ctx := context.Background()

  notes, err := listNotes(ctx, 4)
  if err != nil {
    msg := fmt.Sprintf("Could not get list of notes: %v", err)
    http.Error(w, msg, http.StatusInternalServerError)
    return
  }
  fmt.Fprintln(w, "Notes:")
  for _, nt := range notes {
    fmt.Fprintf(w, "Address = %v\nValue = %d satoshi\nCost = %0.2f NZD\n\n", nt.Addr, nt.Value, nt.Cost/100)
  }

}

func listNotes(ctx context.Context, limit int) ([]*note, error) {
  q := datastore.NewQuery("Note").
    Order("=Value").
    Limit(limit)
  notes := make([]*note, 0)
  _, err := datastoreClient.GetAll(ctx, q, &notes)
  return notes, err
}
