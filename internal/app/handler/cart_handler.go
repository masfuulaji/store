package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/securecookie"
	"github.com/jmoiron/sqlx"
	"github.com/masfuulaji/store/internal/app/models"
	"github.com/masfuulaji/store/internal/app/repositories"
	"github.com/masfuulaji/store/internal/utils"
)

type CartHandler interface {
	CreateCart(w http.ResponseWriter, r *http.Request)
	UpdateCart(w http.ResponseWriter, r *http.Request)
	DeleteCart(w http.ResponseWriter, r *http.Request)
	GetCart(w http.ResponseWriter, r *http.Request)
	GetCategories(w http.ResponseWriter, r *http.Request)
}

type CartHandlerImpl struct {
	cartRepository     *repositories.CartRepositoryImpl
	cartItemRepository *repositories.CartItemRepositoryImpl
	productRepository  *repositories.ProductRepositoryImpl
}

func NewCartHandler(db *sqlx.DB) *CartHandlerImpl {
	return &CartHandlerImpl{cartRepository: repositories.NewCartRepository(db), cartItemRepository: repositories.NewCartItemRepository(db), productRepository: repositories.NewProductRepository(db)}
}

func (f *CartHandlerImpl) CreateCart(w http.ResponseWriter, r *http.Request) {
	order := models.Order{}
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	cookie, err := r.Cookie("jwt")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var tokenString string

	err = securecookie.New([]byte("secret"), nil).Decode("jwt", cookie.Value, &tokenString)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("secret"), nil
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

	order.Cart.UserId = claim["id"].(string)

	id, err := f.cartRepository.CreateCart(order.Cart)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	var totalPrice float64
	for _, items := range order.CartItem {
		product, err := f.productRepository.GetProduct(items.ProductId)
		if err != nil {
			utils.RespondWithJSON(w, 200, map[string]string{"message": "Product Not Found"})
			return
		}

		qtyF, err := strconv.ParseFloat(items.ProductQty, 64)
		if err != nil {
			utils.RespondWithJSON(w, 200, map[string]string{"message": "Qty Wrong"})
			return
		}
		totalItemPrice := product.Price * qtyF

		totalPrice += totalItemPrice

		items.PriceTotal = totalItemPrice
		items.CartId = strconv.Itoa(id)
		err = f.cartItemRepository.CreateCartItem(items)
	}

	order.Cart.PriceTotal = totalPrice
	err = f.cartRepository.UpdateCart(order.Cart, strconv.Itoa(id))
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	utils.RespondWithJSON(w, 200, map[string]string{"message": "Cart created successfully"})
}

func (f *CartHandlerImpl) UpdateCart(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	cart := models.Cart{}
	err := json.NewDecoder(r.Body).Decode(&cart)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	err = f.cartRepository.UpdateCart(cart, id)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	utils.RespondWithJSON(w, 200, map[string]string{"message": "Cart updated successfully"})
}

func (f *CartHandlerImpl) DeleteCart(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := f.cartRepository.DeleteCart(id)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	utils.RespondWithJSON(w, 200, map[string]string{"message": "Cart deleted successfully"})
}

type ProductCart struct {
	Cart     models.Cart       `json:"cart"`
	CartItem []models.CartItem `json:"cart_item"`
}

func (f *CartHandlerImpl) GetCart(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	cart, err := f.cartRepository.GetCart(id)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	cartItems, err := f.cartItemRepository.GetCartItemsByCart(cart.ID)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	response := ProductCart{
		Cart:     cart,
		CartItem: cartItems,
	}
	json.NewEncoder(w).Encode(response)
}

func (f *CartHandlerImpl) GetCartByUserId(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var tokenString string

	err = securecookie.New([]byte("secret"), nil).Decode("jwt", cookie.Value, &tokenString)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("secret"), nil
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

	id := claim["id"].(string)

	cart, err := f.cartRepository.GetCartByUserId(id)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	cartItems, err := f.cartItemRepository.GetCartItemsByCart(cart.ID)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	response := ProductCart{
		Cart:     cart,
		CartItem: cartItems,
	}
	json.NewEncoder(w).Encode(response)
}

type deleteCartItem struct {
	CardID     string `json:"card_id"`
	CardItemID string `json:"card_item_id"`
}

func (f *CartHandlerImpl) DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	deleteCartItem := deleteCartItem{}
	err := json.NewDecoder(r.Body).Decode(&deleteCartItem)

	err = f.cartItemRepository.DeleteCartItem(deleteCartItem.CardItemID)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	cartItems, err := f.cartItemRepository.GetCartItemsByCart(deleteCartItem.CardID)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	var totalPrice float64
	for _, items := range cartItems {
		product, err := f.productRepository.GetProduct(items.ProductId)
		if err != nil {
			utils.RespondWithJSON(w, 200, map[string]string{"message": "Product Not Found"})
			return
		}

		qtyF, err := strconv.ParseFloat(items.ProductQty, 64)
		if err != nil {
			utils.RespondWithJSON(w, 200, map[string]string{"message": "Qty Wrong"})
			return
		}
		totalItemPrice := product.Price * qtyF

		totalPrice += totalItemPrice

		items.PriceTotal = totalItemPrice
	}

	cart, err := f.cartRepository.GetCart(deleteCartItem.CardID)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	cart.PriceTotal = totalPrice
	err = f.cartRepository.UpdateCart(cart, cart.ID)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	response := ProductCart{
		Cart:     cart,
		CartItem: cartItems,
	}
	json.NewEncoder(w).Encode(response)
}

func (f *CartHandlerImpl) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := f.cartRepository.GetCarts()
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	json.NewEncoder(w).Encode(categories)
}
