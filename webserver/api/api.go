package api

import(
	"github.com/olivere/elastic"
	"net/http"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
	//"gopkg.in/cheggaaa/pb.v1"
	"encoding/json"
	"io"
	"bytes"
)

type ApkData struct {
	Permissions		[]string	`json:"Permissions"`
	Malicious		string		`json:"Permissions"`
}

func Query(w http.ResponseWriter, req *http.Request) {
	client, err := elastic.NewClient()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	// Count total and setup progress
	// total, err := client.Count("apks").Type("apk").Do(ctx)
	// if err != nil {
	// 	panic(err)
	// }
	//bar := pb.StartNew(int(total))

	// This example illustrates how to use goroutines to iterate
	// through a result set via ScrollService.
	//
	// It uses the excellent golang.org/x/sync/errgroup package to do so.
	//
	// The first goroutine will Scroll through the result set and send
	// individual documents to a channel.
	//
	// The second cluster of goroutines will receive documents from the channel and
	// deserialize them.
	//
	// Feel free to add a third goroutine to do something with the
	// deserialized results.
	//
	// Let's go.

	// 1st goroutine sends individual hits to channel.
	hits := make(chan json.RawMessage)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		defer close(hits)
		// Initialize scroller. Just don't call Do yet.
		scroll := client.Scroll("apks").Type("apk").Size(100)
		for {
			results, err := scroll.Do(ctx)
			if err == io.EOF {
				return nil // all results retrieved
			}
			if err != nil {
				return err // something went wrong
			}

			// Send the hits to the hits channel
			for _, hit := range results.Hits.Hits {
				select {
				case hits <- *hit.Source:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
		return nil
	})

	// 2nd goroutine receives hits and deserializes them.
	//
	// If you want, setup a number of goroutines handling deserialization in parallel.
	for i := 0; i < 10; i++ {
		g.Go(func() error {
			for hit := range hits {
				r := bytes.NewReader(hit)
				io.Copy(w, r)
				// // Deserialize
				// var a ApkData
				// err := json.Unmarshal(hit, &a)
				// if err != nil {
				// 	return err
				// }

				// // Do something with the product here, e.g. send it to another channel
				// // for further processing.
				// _ = a

				//bar.Increment()

				// Terminate early?
				select {
				default:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			return nil
		})
	}

	// Check whether any goroutines failed.
	if err := g.Wait(); err != nil {
		panic(err)
	}

	// Done.
	//bar.FinishPrint("Done")
	//return nil
}