package domain

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOrderItem_JSONRoundtripWithSelectedAttributes(t *testing.T) {
	item := OrderItem{
		ProductID:          101,
		Quantity:           2,
		SelectedAttributes: map[string]string{"Swing": "250mm", "Color": "Blue"},
	}
	data, err := json.Marshal(item)
	require.NoError(t, err)

	var decoded OrderItem
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	require.Equal(t, item.ProductID, decoded.ProductID)
	require.Equal(t, item.Quantity, decoded.Quantity)
	require.Equal(t, item.SelectedAttributes, decoded.SelectedAttributes)
}

func TestOrderItem_JSONRoundtripWithoutSelectedAttributes(t *testing.T) {
	item := OrderItem{
		ProductID: 102,
		Quantity:  1,
	}
	data, err := json.Marshal(item)
	require.NoError(t, err)

	var decoded OrderItem
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	require.Equal(t, item.ProductID, decoded.ProductID)
	require.Equal(t, item.Quantity, decoded.Quantity)
	require.Nil(t, decoded.SelectedAttributes)
}

func TestOrderItem_JSONUnmarshalFromCartSnapshot(t *testing.T) {
	// Simulates cart_snapshot payload from frontend
	raw := `[{"product_id":1,"quantity":2,"price":99.50,"selected_attributes":{"Swing":"300mm"}}]`
	var items []OrderItem
	err := json.Unmarshal([]byte(raw), &items)
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, int64(1), items[0].ProductID)
	require.Equal(t, 2, items[0].Quantity)
	require.Equal(t, map[string]string{"Swing": "300mm"}, items[0].SelectedAttributes)
}
