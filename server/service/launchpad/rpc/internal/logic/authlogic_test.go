package logic

import (
	"context"
	"database/sql"
	"encoding/hex"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"o1-launchpad/common/auth"
	"o1-launchpad/common/model"
	"o1-launchpad/service/launchpad/rpc/internal/svc"
	"o1-launchpad/service/launchpad/rpc/pb/launchpad"

	"github.com/ethereum/go-ethereum/crypto"
	_ "github.com/lib/pq"
)

func TestAuthChallengeCannotBeReplayed(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL is not set")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	_, filename, _, _ := runtime.Caller(0)
	migration, err := os.ReadFile(filepath.Join(filepath.Dir(filename), "..", "..", "..", "..", "..", "migrations", "001_init.sql"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(string(migration)); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`TRUNCATE auth_nonces CASCADE`); err != nil {
		t.Fatal(err)
	}

	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	address := crypto.PubkeyToAddress(key.PublicKey)
	ctx := context.Background()
	svcCtx := &svc.ServiceContext{Model: model.New(db)}
	created, err := NewCreateAuthChallengeLogic(ctx, svcCtx).CreateAuthChallenge(&launchpad.CreateAuthChallengeRequest{
		ChainId: 84532, Address: address.Hex(), Domain: "demo.example", Uri: "https://demo.example",
	})
	if err != nil {
		t.Fatal(err)
	}
	signature, err := crypto.Sign(auth.PersonalSignHash(created.Message), key)
	if err != nil {
		t.Fatal(err)
	}
	request := &launchpad.VerifyAuthSignatureRequest{
		ChallengeId: created.ChallengeId, Address: address.Hex(), Signature: "0x" + hex.EncodeToString(signature),
	}
	if _, err := NewVerifyAuthSignatureLogic(ctx, svcCtx).VerifyAuthSignature(request); err != nil {
		t.Fatal(err)
	}
	if _, err := NewVerifyAuthSignatureLogic(ctx, svcCtx).VerifyAuthSignature(request); err == nil {
		t.Fatal("expected replay to fail")
	}
}
