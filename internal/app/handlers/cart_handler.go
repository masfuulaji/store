package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/masfuulaji/store/internal/app/models"
	"github.com/masfuulaji/store/internal/app/repositories"
	"github.com/masfuulaji/store/internal/utils"
)

type CartHandler interface {
	AddCartItem(w http.ResponseWriter, r *http.Request)
	DeleteCart(w http.ResponseWriter, r *http.Request)
}

type CartHandlerImpl struct {
	cartRepository     *repositories.CartRepositoryImpl
	cartItemRepository *repositories.CartItemRepositoryImpl
	productRepository  *repositories.ProductRepositoryImpl
}

func NewCartHandler(db *sqlx.DB) *CartHandlerImpl {
	return &CartHandlerImpl{cartRepository: repositories.NewCartRepository(db), cartItemRepository: repositories.NewCartItemRepository(db), productRepository: repositories.NewProductRepository(db)}
}

func (f *CartHandlerImpl) AddCartItem(w http.ResponseWriter, r *http.Request) {
	cartItem := models.CartItem{}
	err := json.NewDecoder(r.Body).Decode(&cartItem)
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

	cart, err := f.cartRepository.CountCartByUserId(idS)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"message": "No Cart"})
		return
	}

	order_data := models.Cart{
		Name:       "cart 1",
		UserId:     "1",
		PriceTotal: 10,
	}
	var cart_id string
	if cart < 1 {
		id, err := f.cartRepository.CreateCart(order_data)
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
		qtyI, err := strconv.Atoi(cartItem.ProductQty)
		if err != nil {
			utils.RespondWithJSON(w, 200, map[string]string{"message": "Qty Wrong"})
			return
		}
		if product.Stock < qtyI {
			utils.RespondWithJSON(w, 200, map[string]string{"message": "Stock Not Enough"})
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
		if err != nil {
			utils.RespondWithJSON(w, 200, map[string]string{"message": "fail"})
			return
		}
	} else {
		res, err := f.cartRepository.GetCartByUserId(idS)
		fmt.Println(idS)
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

		qtyI, err := strconv.Atoi(cartItem.ProductQty)
		if err != nil {
			utils.RespondWithJSON(w, 200, map[string]string{"message": "Qty Wrong"})
			return
		}
		if product.Stock < qtyI {
			utils.RespondWithJSON(w, 200, map[string]string{"message": "Stock Not Enough"})
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
		oldCartItem, err := f.cartItemRepository.GetCartItemByCart(res.ID)
		fmt.Println(err)
		if err != nil {
			err = f.cartItemRepository.CreateCartItem(cartItem)
		} else {
			oldCartItem.ProductQty = cartItem.ProductQty
			oldCartItem.PriceTotal = cartItem.PriceTotal
			err = f.cartItemRepository.UpdateCartItem(oldCartItem, oldCartItem.ID)
		}
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
	cart, err := f.cartRepository.GetCart(id)
	if cart.UserId != idS {

		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	err = f.cartRepository.DeleteCart(id)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	utils.RespondWithJSON(w, 200, map[string]string{"message": "Cart deleted successfully"})
}

func (f *CartHandlerImpl) ReadCart(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

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
	cart, err := f.cartRepository.GetCart(id)
	if cart.UserId != idS {

		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	cartItem, err := f.cartItemRepository.GetCartItemsByCart(id)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	type cartResponse struct {
		Cart      models.Cart       `json:"cart"`
		CartItems []models.CartItem `json:"cart_items"`
	}
	response := cartResponse{
		Cart:      cart,
		CartItems: cartItem,
	}

	utils.RespondWithJSON(w, 200, response)
}
