/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD-3-Clause license that can be found in the LICENSE file.
 */

package ui

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

func validateLiteralData(data any, limits Limits, depth int) error {
	encoded, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("%w: literal is not JSON-compatible: %w", ErrInvalidValue, err)
	}

	decoder := json.NewDecoder(bytes.NewReader(encoded))
	decoder.UseNumber()
	var decoded any
	if err := decoder.Decode(&decoded); err != nil {
		return fmt.Errorf("%w: literal: %w", ErrInvalidValue, err)
	}

	return validateLiteralJSON(decoded, limits, depth)
}

func validateLiteralJSON(value any, limits Limits, depth int) error {
	if depth > limits.MaxExpressionDepth {
		return errors.New("expression depth limit exceeded")
	}

	switch typed := value.(type) {
	case string:
		if len(typed) > limits.MaxStringLength {
			return errors.New("expression string limit exceeded")
		}
	case []any:
		if len(typed) > limits.MaxArrayItems {
			return errors.New("expression array limit exceeded")
		}
		for _, item := range typed {
			if err := validateLiteralJSON(item, limits, depth+1); err != nil {
				return err
			}
		}
	case map[string]any:
		if len(typed) > limits.MaxObjectKeys {
			return errors.New("expression object limit exceeded")
		}
		for key, item := range typed {
			if key == "__proto__" || key == "prototype" || key == "constructor" {
				return errors.New("unsafe object key")
			}
			if err := validateLiteralJSON(item, limits, depth+1); err != nil {
				return err
			}
		}
	}

	return nil
}
