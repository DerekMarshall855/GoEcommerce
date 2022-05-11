package database

import "errors"

var (
	ErrCantFindProduct   = errors.New("Product Cannot Be Found")
	ErrCantDecodeProduct = errors.New("Product Cannot Be Decoded")
	ErrCantUpdateUser    = errors.New("User Cannot Be Updated")
	ErrCantRemoveItem    = errors.New("Item Cannot Be Removed")
	ErrCantBuyItem       = errors.New("Item Cannot Be Bought")
)

func AddToCart() {

}

func RemoveFromCart() {

}

func BuyItemFromCart() {

}

func InstantBuy() {

}
