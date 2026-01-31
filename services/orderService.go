package services

import (
	"fmt"
	"go-inventory/models"
	"go-inventory/repository"
	"time"

	"github.com/google/uuid"
)

type OrderServiceImpl struct {
	OrderRepo       *repository.OrderInMemory
	LoggerRepo      *repository.LogInMemory
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
	// ini butuh userID, Slice Item
	if o.UserID == "" {
		return ErrInvalid
	}
	for _, v := range o.Item {
		if v.ProductID == "" || v.Qty <= 0 {
			return ErrInvalid
		}
	}
	newOrder := models.Order{
		ID:        "order-" + uuid.NewString(),
		UserID:    o.UserID,
		Item:      o.Item,
		Status:    models.PENDING,
		CreatedAt: time.Now().UTC(),
		//UpdatedAt
	}
	errLog := s.LoggerRepo.CreateLog(models.TransactionLog{
		ID:        "log-" + uuid.NewString(),
		EntityID:  newOrder.ID,
		Entity:    models.ORDER,
		Action:    "CREATE_ORDER",
		CreatedAt: time.Now().UTC(),
		Note:      fmt.Sprintf("Action:%s - Status:%s - Time:%s ", "Create Order", "Pending", time.Now().Format(time.DateTime)),
	})
	if errLog != nil {
		fmt.Println(errLog.Error())
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
		if err := s.ProductServices.Sell(v.ProductID, v.Qty); err != nil {
			return err
		} // jangan terminate setengah2 ketika update jalan (ini last cover harusnya sudah lewat)
	}
	errLog := s.LoggerRepo.CreateLog(models.TransactionLog{
		ID:        "log-" + uuid.NewString(),
		EntityID:  d.ID,
		Entity:    models.ORDER,
		Action:    "PAY_ORDER",
		CreatedAt: time.Now().UTC(),
		Note:      fmt.Sprintf("Action:%s - Status:%s - Time:%s ", "Pay Order", "Paid", time.Now().Format(time.DateTime)),
	})
	if errLog != nil {
		fmt.Println(errLog.Error())
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
	errLog := s.LoggerRepo.CreateLog(models.TransactionLog{
		ID:        "log-" + uuid.NewString(),
		EntityID:  d.ID,
		Entity:    models.ORDER,
		Action:    "CANCEL_ORDER",
		CreatedAt: time.Now().UTC(),
		Note:      fmt.Sprintf("Action:%s - Status:%s - Time:%s ", "Cancel Order", "Cancelled", time.Now().Format(time.DateTime)),
	})
	if errLog != nil {
		fmt.Println(errLog.Error())
	}
	return s.OrderRepo.Update(d)
}
