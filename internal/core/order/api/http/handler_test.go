package http

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToOrderItemInputs_WithSelectedAttributes(t *testing.T) {
	price := 99.50
	items := []OrderItemRequest{
		{
			ProductID:          1,
			Quantity:           2,
			Price:              &price,
			SelectedAttributes: map[string]string{"Swing": "250mm"},
		},
	}
	result, err := toOrderItemInputs(items)
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, int64(1), result[0].ProductID)
	require.Equal(t, 2, result[0].Quantity)
	require.Equal(t, map[string]string{"Swing": "250mm"}, result[0].SelectedAttributes)
}

func TestToOrderItemInputs_WithoutSelectedAttributes(t *testing.T) {
	price := 50.0
	items := []OrderItemRequest{
		{ProductID: 2, Quantity: 1, Price: &price},
	}
	result, err := toOrderItemInputs(items)
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, int64(2), result[0].ProductID)
	require.Nil(t, result[0].SelectedAttributes)
}

func TestToOrderItemInputs_UseIDWhenProductIDZero(t *testing.T) {
	items := []OrderItemRequest{
		{ID: 42, ProductID: 0, Quantity: 1},
	}
	result, err := toOrderItemInputs(items)
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, int64(42), result[0].ProductID)
}

func TestToOrderItemInputs_ErrorWhenNoProductID(t *testing.T) {
	items := []OrderItemRequest{
		{Quantity: 1},
	}
	_, err := toOrderItemInputs(items)
	require.Error(t, err)
	require.Contains(t, err.Error(), "product_id")
}

func TestOrderItemResponse_HasSelectedAttributesField(t *testing.T) {
	// Ensure OrderItemResponse can carry selected_attributes to clients
	resp := OrderItemResponse{
		ProductID:          1,
		Quantity:           1,
		SelectedAttributes: map[string]string{"Swing": "300mm"},
	}
	require.Equal(t, map[string]string{"Swing": "300mm"}, resp.SelectedAttributes)
}
