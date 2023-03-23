package api

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/ryadavDeqode/dq-vault/config"
	"github.com/ryadavDeqode/dq-vault/test/unit_test/mocks"
)

func TestPathRemoveGuardian(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mocks.NewMockStorage(ctrl)

	tErr := errors.New("test error")

	s.EXPECT().Get(context.Background(), config.StorageBasePath+"test").Return(&logical.StorageEntry{}, tErr)
	s.EXPECT().List(context.Background(), config.StorageBasePath).Return([]string{"test"}, nil).AnyTimes()
	s.EXPECT().Put(context.Background(), gomock.Any()).Return(nil).AnyTimes()
	b := backend{}
	req := logical.Request{}

	req.Storage = s

	MPatchGet("test")

	res, err := b.pathRemoveGuardian(context.Background(), &req, &framework.FieldData{})

	if tErr.Error() != err.Error() {
		t.Error("expected test error, received - ", res, err)
	}

	mpd := MPatchDecodeJSON(tErr)
	s.EXPECT().Get(context.Background(), config.StorageBasePath+"test").Return(&logical.StorageEntry{}, nil).AnyTimes()

	res, err = b.pathRemoveGuardian(context.Background(), &req, &framework.FieldData{})
	mpd.Unpatch()

	if tErr.Error() != err.Error() {
		t.Error("expected test error, received - ", res, err)
	}

	MPatchDecodeJSON(nil)
	mpjwt := MPatchVerifyJWTSignature(false, tErr.Error())

	res, err = b.pathRemoveGuardian(context.Background(), &req, &framework.FieldData{})

	if err != nil {
		t.Error(" error wasn't expected, received - ", err)
	} else {
		if res.Data["status"].(bool) {
			t.Error(" unexpected value of status,expected false, received - ", res)
		}

		if res.Data["remarks"] != tErr.Error() {
			t.Error("unexpected value of remarks, expected - ", tErr.Error(), "received - ", res)
		}
	}

	mpjwt.Unpatch()
	MPatchVerifyJWTSignature(true, tErr.Error())
	MPatchStringInSlice(true)
	MPatchNewClient()
	MPatchGetPubSub(tErr.Error(), nil)

	res, err = b.pathRemoveGuardian(context.Background(), &req, &framework.FieldData{})

	if err != nil {
		t.Error("Expected no error, received - ", err)
	} else if !res.Data["status"].(bool) {

		t.Error(" unexpected value of status,expected true, received - ", res)
	}

}
