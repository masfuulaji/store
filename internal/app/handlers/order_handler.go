package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/masfuulaji/store/internal/app/models"
	"github.com/masfuulaji/store/internal/app/repositories"
	"github.com/masfuulaji/store/internal/utils"
)

type OrderHandler interface {
	CreateOrder(w http.ResponseWriter, r *http.Request)
}

type OrderHandlerImpl struct {
	orderRepository    *repositories.OrderRepositoryImpl
	cartRepository     *repositories.CartRepositoryImpl
	cartItemRepository *repositories.CartItemRepositoryImpl
	productRepository  *repositories.ProductRepositoryImpl
}

func NewOrderHandler(db *sqlx.DB) *OrderHandlerImpl {
	return &OrderHandlerImpl{orderRepository: repositories.NewOrderRepository(db), cartRepository: repositories.NewCartRepository(db), cartItemRepository: repositories.NewCartItemRepository(db), productRepository: repositories.NewProductRepository(db)}
}

func (f *OrderHandlerImpl) CreateOrder(w http.ResponseWriter, r *http.Request) {
	order := models.Order{}
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("secret_key"), nil
	})
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	claim, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idS := fmt.Sprintf("%v", claim["id"])
	cart, err := f.cartRepository.GetCart(order.CartId)
	if cart.UserId != idS {

		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if cart.Finish.Valid && cart.Finish.Int32 == 1 {

		utils.RespondWithJSON(w, 200, map[string]string{"message": "Cart Empty"})
		return
	}

	cartItems, err := f.cartItemRepository.GetCartItemsByCart(order.CartId)

	for _, cartItem := range cartItems {
		product, err := f.productRepository.GetProduct(cartItem.ProductId)
		if err != nil {
			utils.RespondWithJSON(w, 200, map[string]string{"message": "Product Not Found"})
			return
		}
		qtyI, err := strconv.Atoi(cartItem.ProductQty)
		if err != nil {
			utils.RespondWithJSON(w, 200, map[string]string{"message": "Qty Wrong"})
			return
		}
		if product.Stock < qtyI {
			utils.RespondWithJSON(w, 200, map[string]string{"message": "Stock Not Enough"})
			return
		}

		stock := product.Stock - qtyI
		err = f.productRepository.UpdateProductStock(stock, product.ID)
	}

	order.PriceTotal = cart.PriceTotal
	id, err := f.orderRepository.CreateOrder(order)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	err = f.cartRepository.UpdateCartFinish(1, order.CartId)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	utils.RespondWithJSON(w, 200, map[string]string{"message": fmt.Sprintf("Order created successfully %d", id)})
}
