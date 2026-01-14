package mock

import (
	"math/rand"
	"time"

	"abacatepay-cli/internal/utils"

	v1 "github.com/almeidazs/go-abacate-types/v1"
	"github.com/brianvoe/gofakeit/v7"
)

func CreatePixQRCodeMock() *v1.RESTPostCreateQRCodePixBody {
	expires := 15 * 30
	desc := "salve"
	gofakeit.Seed(0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return &v1.RESTPostCreateQRCodePixBody{
		Amount:      gofakeit.Number(100, 10000),
		ExpiresIn:   &expires,
		Description: &desc,
		Customer: &v1.APICustomerMetadata{
			Name:      gofakeit.Name(),
			Email:     gofakeit.Email(),
			TaxID:     utils.GenerateValidCPF(r),
			Cellphone: "11999999999",
		},
	}
}

