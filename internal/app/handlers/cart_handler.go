package handlers

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

func (f *CartHandlerImpl) addCartItem(w http.ResponseWriter, r *http.Request) {
	cartItem := models.CartItem{}
	err := json.NewDecoder(r.Body).Decode(&cartItem)
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

	cart, err := f.cartRepository.CountCartByUserId(claim["id"].(string))
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	order := models.Cart{
		Name:   "cart 1",
		UserId: claim["id"].(string),
	}
	var cart_id string
	if cart > 0 {
		id, err := f.cartRepository.CreateCart(order)
		if err != nil {
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		cart_id = strconv.Itoa(id)

		product, err := f.productRepository.GetProduct(cartItem.ProductId)
		if err != nil {
			utils.RespondWithJSON(w, 200, map[string]string{"message": "Product Not Found"})
			return
		}

		qtyF, err := strconv.ParseFloat(cartItem.ProductQty, 64)
		if err != nil {
			utils.RespondWithJSON(w, 200, map[string]string{"message": "Qty Wrong"})
			return
		}
		totalItemPrice := product.Price * qtyF

		cartItem.PriceTotal = totalItemPrice
		cartItem.CartId = strconv.Itoa(id)
		err = f.cartItemRepository.CreateCartItem(cartItem)
	} else {
		res, err := f.cartRepository.GetCartByUserId(claim["id"].(string))
		if err != nil {
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		cart_id = res.ID

		product, err := f.productRepository.GetProduct(cartItem.ProductId)
		if err != nil {
			utils.RespondWithJSON(w, 200, map[string]string{"message": "Product Not Found"})
			return
		}

		qtyF, err := strconv.ParseFloat(cartItem.ProductQty, 64)
		if err != nil {
			utils.RespondWithJSON(w, 200, map[string]string{"message": "Qty Wrong"})
			return
		}
		totalItemPrice := product.Price * qtyF

		cartItem.PriceTotal = totalItemPrice
		cartItem.CartId = res.ID
		err = f.cartItemRepository.CreateCartItem(cartItem)
	}

	sum, err := f.cartItemRepository.SumCartItemByCart(cart_id)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	err = f.cartRepository.UpdateCartTotal(sum, cart_id)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	utils.RespondWithJSON(w, 200, map[string]string{"message": "Cart Saved"})
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
