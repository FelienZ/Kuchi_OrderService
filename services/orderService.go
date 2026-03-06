package services

import (
	"context"
	"fmt"
	"go-inventory/models"
	"go-inventory/repository/database"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderServiceImpl struct {
	Db            *pgxpool.Pool
	OrderRepoDB   models.OrderRepository
	ProductRepoDB models.ProductRepository
	LoggerRepo    models.LoggerService
}

func (s *OrderServiceImpl) GetByID(ctx context.Context, id string) (models.Order, error) {
	if id == "" {
		return models.Order{}, ErrOrderInvalid
	}
	d, err := s.OrderRepoDB.FindByID(ctx, s.Db, id)
	if err == database.ErrNoRows {
		return models.Order{}, ErrOrderNotFound
	}
	if err != nil {
		return models.Order{}, err
	}
	return d, nil
}

func (s *OrderServiceImpl) List(ctx context.Context, filter models.GetOrderParameter) ([]models.Order, error) {
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	return s.OrderRepoDB.FindAll(ctx, s.Db, filter)
}

func (s *OrderServiceImpl) GetByUserID(ctx context.Context, userid string) ([]models.Order, error) {
	if userid == "" {
		return []models.Order{}, ErrOrderInvalid
	}
	p, err := s.OrderRepoDB.FindByUserID(ctx, s.Db, userid)
	if err == database.ErrNoRows {
		return []models.Order{}, nil
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *OrderServiceImpl) CreateOrder(ctx context.Context, o models.OrderPayload) (string, error) {
	// ini butuh userID, Slice Item
	if o.UserID == "" {
		return "", ErrOrderInvalid
	}
	tx, errTx := s.Db.Begin(ctx)
	if errTx != nil {
		return "", errTx
	}
	defer tx.Rollback(ctx)
	newOrder := models.Order{
		ID:        uuid.NewString(),
		UserID:    o.UserID,
		Item:      []models.OrderItem{},
		Status:    models.PENDING,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	for _, v := range o.Item {
		if v.ProductID == "" || v.Qty <= 0 {
			return "", ErrOrderInvalid
		}
		p, err := s.ProductRepoDB.FindByID(ctx, tx, v.ProductID)
		if err == database.ErrNoRows {
			return "", ErrProductNotFound
		}
		if err != nil {
			// fmt.Println("cek err create order find: ", err)
			return "", err
		}
		newOrder.Item = append(newOrder.Item, models.OrderItem{ProductID: v.ProductID, Qty: v.Qty,
			TotalPrice: p.Price * v.Qty})
	}
	id, err := s.OrderRepoDB.Save(ctx, tx, newOrder)
	if err != nil {
		if errViolation := ErrorOrderDomainTranslator(err); errViolation != nil {
			return "", errViolation
		}
		// fmt.Println("cek err create order: ", err)
		return "", err
	}
	errLog := s.LoggerRepo.CreateLog(ctx, models.TransactionLog{
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
	return id, tx.Commit(ctx)
}

func (s *OrderServiceImpl) PayOrder(ctx context.Context, id string) error {
	if id == "" {
		return ErrOrderInvalid
	}
	tx, errTx := s.Db.Begin(ctx)
	// fmt.Println("mulai payOrder")
	if errTx != nil {
		return errTx
	}
	defer tx.Rollback(ctx)
	// run 3 query
	// update dulu, race -> rollback invalid status
	if errUpdate := s.OrderRepoDB.UpdateOrderStatus(ctx, tx, models.PENDING, models.PAID, id); errUpdate != nil {
		// fmt.Println("Cek error invalid status payOrder: ", errUpdate)
		if errViolation := ErrorOrderDomainTranslator(errUpdate); errViolation != nil {
			return errViolation
		}
		return errUpdate
	}
	order, err := s.OrderRepoDB.FindByID(ctx, tx, id)
	// fmt.Println("Cek error invalid find Order Pay (err): ", err)
	if err == database.ErrNoRows {
		return ErrOrderNotFound
	}
	if err != nil {
		return err
	}
	for _, v := range order.Item {
		if err = s.ProductRepoDB.DecreaseStock(ctx, tx, v.ProductID, v.Qty); err != nil {
			// fmt.Println("Cek error invalid status decreasestok: ", err)
			if errViolation := ErrorProductDomainTranslator(err); errViolation != nil {
				return errViolation
			}
			if err == database.ErrNoUpdate {
				return ErrProductInvalid
			}
			return err
		}
	}
	errLog := s.LoggerRepo.CreateLog(ctx, models.TransactionLog{
		ID:        "log-" + uuid.NewString(),
		EntityID:  id,
		Entity:    models.ORDER,
		Action:    "PAY_ORDER",
		CreatedAt: time.Now().UTC(),
		Note:      fmt.Sprintf("Action:%s - Status:%s - Time:%s ", "Pay Order", "Paid", time.Now().Format(time.DateTime)),
	})
	if errLog != nil {
		fmt.Println(errLog.Error())
	}
	return tx.Commit(ctx)
}

func (s *OrderServiceImpl) CancelOrder(ctx context.Context, id string) error {
	if id == "" {
		return ErrOrderInvalid
	}
	tx, errTx := s.Db.Begin(ctx)
	if errTx != nil {
		return errTx
	}
	defer tx.Rollback(ctx)
	if errUpdate := s.OrderRepoDB.UpdateOrderStatus(ctx, tx, models.PENDING, models.CANCELED, id); errUpdate != nil {
		if errViolation := ErrorOrderDomainTranslator(errUpdate); errViolation != nil {
			return errViolation
		}
		if errUpdate == database.ErrNoUpdate {
			return ErrOrderInvalid
		}
		return errUpdate
	}
	errLog := s.LoggerRepo.CreateLog(ctx, models.TransactionLog{
		ID:        "log-" + uuid.NewString(),
		EntityID:  id,
		Entity:    models.ORDER,
		Action:    "CANCEL_ORDER",
		CreatedAt: time.Now().UTC(),
		Note:      fmt.Sprintf("Action:%s - Status:%s - Time:%s ", "Cancel Order", "Cancelled", time.Now().Format(time.DateTime)),
	})
	if errLog != nil {
		fmt.Println(errLog.Error())
	}
	return tx.Commit(ctx)
}

func (s *OrderServiceImpl) DeleteOrder(ctx context.Context, id string) error {
	if id == "" {
		return ErrOrderInvalid
	}
	tx, errTx := s.Db.Begin(ctx)
	if errTx != nil {
		return errTx
	}
	defer tx.Rollback(ctx)
	if errDelete := s.OrderRepoDB.Delete(ctx, tx, id); errDelete != nil {
		if errViolation := ErrorOrderDomainTranslator(errDelete); errViolation != nil {
			return errViolation
		}
		if errDelete == database.ErrNoDelete {
			return ErrOrderInvalid
		}
		return errDelete
	}
	errLog := s.LoggerRepo.CreateLog(ctx, models.TransactionLog{
		ID:        "log-" + uuid.NewString(),
		EntityID:  id,
		Entity:    models.ORDER,
		Action:    "DELETE_ORDER",
		CreatedAt: time.Now().UTC(),
		Note:      fmt.Sprintf("Action:%s - Status:%s - Time:%s ", "Delete Order", "Deleted", time.Now().Format(time.DateTime)),
	})
	if errLog != nil {
		fmt.Println(errLog.Error())
	}
	return tx.Commit(ctx)
}
