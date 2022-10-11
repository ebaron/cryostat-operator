// Copyright The Cryostat Authors
//
// The Universal Permissive License (UPL), Version 1.0
//
// Subject to the condition set forth below, permission is hereby granted to any
// person obtaining a copy of this software, associated documentation and/or data
// (collectively the "Software"), free of charge and under any and all copyright
// rights in the Software, and any and all patent rights owned or freely
// licensable by each licensor hereunder covering either (i) the unmodified
// Software as contributed to or provided by such licensor, or (ii) the Larger
// Works (as defined below), to deal in both
//
// (a) the Software, and
// (b) any piece of software and/or hardware listed in the lrgrwrks.txt file if
// one is included with the Software (each a "Larger Work" to which the Software
// is contributed by such licensors),
//
// without restriction, including without limitation the rights to copy, create
// derivative works of, display, perform, and distribute the Software and make,
// use, sell, offer for sale, import, export, have made, and have sold the
// Software and the Larger Work(s), and to sublicense the foregoing rights on
// either these or other terms.
//
// This license is subject to the following condition:
// The above copyright notice and either this complete permission notice or at
// a minimum a reference to the UPL must be included in all copies or
// substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package test

import (
	"context"
	"time"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// TODO When using envtest instead of fake client, this is probably no longer needed
type timestampClient struct {
	ctrlclient.Client
}

func NewClientWithTimestamp(client ctrlclient.Client) ctrlclient.Client {
	return &timestampClient{
		Client: client,
	}
}

func (c *timestampClient) Create(ctx context.Context, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
	err := SetCreationTimestamp(obj)
	if err != nil {
		return err
	}
	return c.Client.Create(ctx, obj, opts...)
}

var creationTimestamp = metav1.NewTime(time.Unix(1664573254, 0))

func SetCreationTimestamp(objs ...runtime.Object) error {
	for _, obj := range objs {
		metaObj, err := meta.Accessor(obj)
		if err != nil {
			return err
		}
		metaObj.SetCreationTimestamp(creationTimestamp)
	}
	return nil
}

type clientUpdateError struct {
	ctrlclient.Client
	failObj ctrlclient.Object
	err     *kerrors.StatusError
}

// NewClientWithUpdateError wraps a Client by returning an error when updating
// a specified object
func NewClientWithUpdateError(client ctrlclient.Client, failObj ctrlclient.Object,
	err *kerrors.StatusError) ctrlclient.Client {
	return &clientUpdateError{
		Client:  client,
		failObj: failObj,
		err:     err,
	}
}

func (c *clientUpdateError) Update(ctx context.Context, obj ctrlclient.Object,
	opts ...ctrlclient.UpdateOption) error {
	if obj.GetName() == c.failObj.GetName() && obj.GetNamespace() == c.failObj.GetNamespace() {
		// Look up Kind and compare against object to fail on
		match, err := c.matchesKind(obj)
		if err != nil {
			return err
		}
		if *match {
			return c.err
		}
	}
	return c.Client.Update(ctx, obj, opts...)
}

func (c *clientUpdateError) matchesKind(obj ctrlclient.Object) (*bool, error) {
	match := false
	failKinds, _, err := c.Scheme().ObjectKinds(c.failObj)
	if err != nil {
		return nil, err
	}
	kinds, _, err := c.Scheme().ObjectKinds(obj)
	if err != nil {
		return nil, err
	}

	for _, failKind := range failKinds {
		for _, kind := range kinds {
			if failKind == kind {
				match = true
				return &match, nil
			}
		}
	}
	return &match, nil
}
