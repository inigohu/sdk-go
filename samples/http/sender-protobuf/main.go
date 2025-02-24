//
//  Copyright 2021 The CloudEvents Authors
//  SPDX-License-Identifier: Apache-2.0
//

package main

import (
	"context"
	"log"

	pbcloudevents "github.com/cloudevents/sdk-go/binding/format/protobuf/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
)

func main() {
	ctx := cloudevents.ContextWithTarget(context.Background(), "http://localhost:8080/")

	p, err := cloudevents.NewHTTP()
	if err != nil {
		log.Fatalf("failed to create protocol: %s", err.Error())
	}

	c, err := cloudevents.NewClient(p, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	for i := 0; i < 10; i++ {
		data := &Sample{Value: "sample"}
		e := cloudevents.NewEvent()
		e.SetType("com.cloudevents.sample.sent")
		e.SetSource("https://github.com/cloudevents/sdk-go/v2/samples/http/sender-protobuf")
		e.SetDataSchema("my-schema-registry://" + string(data.ProtoReflect().Descriptor().FullName()))
		_ = e.SetData(pbcloudevents.ContentTypeProtobuf, data)

		res := c.Send(ctx, e)
		if cloudevents.IsUndelivered(res) {
			log.Printf("Failed to send: %v", res)
		} else {
			var httpResult *cehttp.Result
			if cloudevents.ResultAs(res, &httpResult) {
				log.Printf("Sent %d with status code %d", i, httpResult.StatusCode)
			}
			log.Printf("Send did not return an HTTP response: %s", res)
		}
	}
}
