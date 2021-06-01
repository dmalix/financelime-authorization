package generator

import (
	"context"
	"strings"
)

func GenerateRequestID(_ context.Context, short bool) (string, error) {

	var requestID string

	prefix := StringRand(4, 4, false)
	prefixArr := strings.Split(prefix, "")
	for _, v := range prefixArr {
		requestID = requestID + v + StringRand(2, 2, false)
	}
	requestID = prefix + requestID
	if !short {
		requestID = requestID + StringRand(48, 48, false)
	}

	return requestID, nil
}
