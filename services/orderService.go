package services

import (
	"fmt"
	"go-inventory/models"
	"go-inventory/repository"
	"time"

	"github.com/google/uuid"
)

type OrderServiceImpl struct {
	OrderRepo       models.OrderRepository
	LoggerRepo      *repository.LogInMemory
	ProductServices models.ProductService
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
	snapshot := []models.ProductSnapshot{}
	for _, v := range d.Item {
		// gagal -> snapshot, rollback
		product, err := s.ProductServices.GetByID(v.ProductID)
		if err != nil {
			return err
		}
		if product.Stock < v.Qty {
			return ErrNotEnough
		}
		snapshot = append(snapshot, models.ProductSnapshot{
			ProductID: product.ID,
			Stock:     product.Stock,
		})
	}
	if d.Status != models.PENDING {
		return ErrConflict
	}
	for i, v := range d.Item {
		if err := s.ProductServices.Sell(v.ProductID, v.Qty); err != nil {
			// semisal ada error, rollback semua (tolak transaksi awal-akhir)
			for j := range i {
				snapshot := snapshot[j]
				_ = s.ProductServices.RecoverStock(snapshot.ProductID, snapshot.Stock)
			}
			return err
		}
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
	d.Status = models.PAID
	d.UpdatedAt = time.Now().UTC()
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

func (s *OrderServiceImpl) DeleteOrder(id string) error {
	if id == "" {
		return ErrInvalid
	}
	orderData, err := s.OrderRepo.FindByID(id)
	if err != nil {
		return err
	}
	if orderData.Status == models.PAID {
		return ErrConflict
	}
	if errDelete := s.OrderRepo.Delete(orderData.ID); errDelete != nil {
		return errDelete
	}
	errLog := s.LoggerRepo.CreateLog(models.TransactionLog{
		ID:        "log-" + uuid.NewString(),
		EntityID:  orderData.ID,
		Entity:    models.ORDER,
		Action:    "DELETE_ORDER",
		CreatedAt: time.Now().UTC(),
		Note:      fmt.Sprintf("Action:%s - Status:%s - Time:%s ", "Delete Order", "Deleted", time.Now().Format(time.DateTime)),
	})
	if errLog != nil {
		fmt.Println(errLog.Error())
	}
	return nil
}
