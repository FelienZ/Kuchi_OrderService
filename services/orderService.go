package services

import (
	"go-inventory/models"
	"go-inventory/repository"
	"time"

	"github.com/google/uuid"
)

type OrderServiceImpl struct {
	OrderRepo       *repository.OrderInMemory
	ProductServices *ProductServiceImpl
}

func (s *OrderServiceImpl) GetByID(id string) (models.Order, error) {
	if id == "" {
		return models.Order{}, ErrInvalid
	}
	return s.OrderRepo.FindByID(id)
}

func (s *OrderServiceImpl) List() []models.Order {
	return s.OrderRepo.FindAll()
}

func (s *OrderServiceImpl) GetByUserID(userid string) ([]models.Order, error) {
	if userid == "" {
		return []models.Order{}, ErrInvalid
	}
	return s.OrderRepo.FindByUserID(userid), nil
}

func (s *OrderServiceImpl) CreateOrder(o models.Order) error {
	if o.UserID == "" {
		return ErrInvalid
	}
	for _, v := range o.Item {
		if v.ProductID == "" || v.Qty <= 0 {
			return ErrInvalid
		}
	}
	newOrder := models.Order{
		ID:        "order-" + uuid.New().String(),
		UserID:    o.UserID,
		Item:      o.Item,
		Status:    models.PENDING,
		CreatedAt: time.Now().UTC(),
		//UpdatedAt
	}
	return s.OrderRepo.Save(newOrder)
}

func (s *OrderServiceImpl) PayOrder(id string) error {
	if id == "" {
		return ErrInvalid
	}
	d, err := s.OrderRepo.FindByID(id)
	if err != nil {
		return err
	}
	// kalau pending -> update
	if d.Status == models.PENDING {
		d.Status = models.PAID
		d.UpdatedAt = time.Now().UTC()
	} else {
		return ErrConflict
	}
	for _, v := range d.Item {
		// sell validasi stock, update stock
		product, err := s.ProductServices.GetByID(v.ProductID)
		if err != nil {
			return err // terminate transaksi
		}
		if product.Stock < v.Qty {
			return ErrNotEnough // terminate transaksi
		}
	}
	// pastikan sukses validasi dulu, baru eksekusi
	for _, v := range d.Item {
		s.ProductServices.Sell(v.ProductID, v.Qty) // jangan terminate setengah2 ketika update jalan
	}
	return s.OrderRepo.Update(d)
}

func (s *OrderServiceImpl) CancelOrder(id string) error {
	if id == "" {
		return ErrInvalid
	}
	d, err := s.OrderRepo.FindByID(id)
	if err != nil {
		return err
	}
	if d.Status == models.PENDING {
		d.Status = models.CANCELLED
		d.UpdatedAt = time.Now().UTC()
	} else {
		return ErrConflict
	}
	return s.OrderRepo.Update(d)
}
