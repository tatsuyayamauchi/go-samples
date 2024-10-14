package handler

import (
	"context"
	"fmt"
	"strconv"

	"github.com/tatsuyayamauchi/go-samples/echo-strict-server/spec"
)

func NewHandler() *Handler {
	return &Handler{}

}

type Handler struct {
}

func (h *Handler) ListPets(ctx context.Context, request spec.ListPetsRequestObject) (spec.ListPetsResponseObject, error) {
	fmt.Printf("called ListPets\n")

	return &spec.ListPets200JSONResponse{
		Body: spec.Pets{
			{
				Id:   1,
				Name: "taro",
			},
			{
				Id:   2,
				Name: "jiro",
			},
		},
		Headers: spec.ListPets200ResponseHeaders{
			XNext: "next",
		},
	}, nil
}

func (h *Handler) CreatePets(ctx context.Context, request spec.CreatePetsRequestObject) (spec.CreatePetsResponseObject, error) {
	fmt.Printf("called CreatePets\n")

	return &spec.CreatePets201Response{}, nil
}

func (h *Handler) ShowPetById(ctx context.Context, request spec.ShowPetByIdRequestObject) (spec.ShowPetByIdResponseObject, error) {
	fmt.Printf("called ShowPetById\n")

	petId, err := strconv.Atoi(request.PetId)
	if err != nil {
		return nil, err
	}
	return &spec.ShowPetById200JSONResponse{
		Id:   int64(petId),
		Name: "taro",
	}, nil
}
