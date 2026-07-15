// Copyright (C) 2026 Eneo Tecnologia S.L.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.

package kafkaclient

import (
	"testing"
)

func TestIsTopicEmpty_EmptyBrokers(t *testing.T) {
	_, err := IsTopicEmpty([]string{}, "test-topic")
	if err == nil {
		t.Error("expected error when no brokers are provided, got nil")
	}
}
