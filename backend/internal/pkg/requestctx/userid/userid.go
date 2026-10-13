package userid

import (
	"context"

	"fortunefinger/backend/internal/pkg/requestctx"
)

const Header = "X-User-ID"

func Get(ctx context.Context) string {
	return requestctx.GetHeader(ctx, Header)
}
