package graphql

import (
	context "context"

	models "gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
)

type consolationPrizeResolver struct{ *Resolver }

func (r *consolationPrizeResolver) Guarantee(ctx context.Context, obj *models.ConsolationPrize) (string, error) {
	return obj.Guarantee.Decimal.String(), nil
}

func (r *consolationPrizeResolver) CarryIn(ctx context.Context, obj *models.ConsolationPrize) (string, error) {
	return obj.CarryIn.Decimal.String(), nil
}

func (r *consolationPrizeResolver) Allocation(ctx context.Context, obj *models.ConsolationPrize) (string, error) {
	return obj.Allocation.Decimal.String(), nil
}
