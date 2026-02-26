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
	LoggerRepo      models.LoggerService
	ProductServices models.ProductService
}

func (s *OrderServiceImpl) GetByID(id string) (models.Order, error) {
	if id == "" {
		return models.Order{}, ErrOrderInvalid
	}
	d, err := s.OrderRepo.FindByID(id)
	if err == repository.ErrNotFound {
		return models.Order{}, ErrOrderNotFound
	}
	if err != nil {
		return models.Order{}, err
	}
	return d, nil
}

func (s *OrderServiceImpl) FilterByStatus(status models.Status, orders []models.Order) []models.Order {
	res := []models.Order{}
	for _, v := range orders {
		if v.Status == status {
			res = append(res, v)
		}
	}
	return res
}

func (s *OrderServiceImpl) PaginateList(orders []models.Order, limit, offset int) []models.Order {
	res := []models.Order{}
	/* Ini kalau soal page, lebih dynamic dominan UI
	// misal end pada paginasi per 10, di halaman 2 berarti di item ke 20 (dari 11-20) anggap mulai dari 1
	start := 1 + ((offset + 1) * limit) - limit //1, 11, 21
	end := start + limit - 1                    // 10, 20, 30
	// fmt.Println("CEK START & END: ", start, end)
	*/
	start := offset       // 0, 10 (ambil item dari idx 0 atau 10 dsb) as start
	end := offset + limit //offset 0, limit 10 -> end 10 , off 10, limit 10 -> end 20
	// kenapa -1 mulainya karena index order dari 0
	// fmt.Println("Cek start, end: ", start, end, len(orders))
	if end > len(orders) || start > len(orders) {
		// out of bound -> kosong
		return res
	}
	for i := start; i < end; i++ {
		res = append(res, orders[i])
	}
	return res
}

func (s *OrderServiceImpl) List(filter models.GetOrderParameter) []models.Order {
	list := s.OrderRepo.FindAll()
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	if filter.Status.String() != "" {
		list = s.FilterByStatus(filter.Status, list)
		// default if > 10, limit 10 (max) else len (min)
		filter.Limit = min(len(list), 10)

	}
	return s.PaginateList(list, filter.Limit, filter.Offset)
}

func (s *OrderServiceImpl) GetByUserID(userid string) ([]models.Order, error) {
	if userid == "" {
		return []models.Order{}, ErrOrderInvalid
	}
	return s.OrderRepo.FindByUserID(userid), nil
}

func (s *OrderServiceImpl) CreateOrder(o models.Order) error {
	// ini butuh userID, Slice Item
	if o.UserID == "" {
		return ErrOrderInvalid
	}
	for _, v := range o.Item {
		if v.ProductID == "" || v.Qty <= 0 {
			return ErrOrderInvalid
		}
		if _, err := s.ProductServices.GetByID(v.ProductID); err != nil {
			return err
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
	err := s.OrderRepo.Save(newOrder)
	if err == repository.ErrConflict {
		return ErrOrderConflict
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *OrderServiceImpl) PayOrder(id string) error {
	if id == "" {
		return ErrOrderInvalid
	}
	d, err := s.OrderRepo.FindByID(id)
	if err == repository.ErrNotFound {
		return ErrOrderNotFound
	}
	if err != nil {
		return err
	}
	snapshot := []models.ProductSnapshot{}
	for _, v := range d.Item {
		// gagal -> snapshot, rollback
		product, err := s.ProductServices.GetByID(v.ProductID)
		if err != nil {
			return err
		} // sudah dihandle di product service
		if product.Stock < v.Qty {
			return ErrProductNotEnough
		}
		snapshot = append(snapshot, models.ProductSnapshot{
			ProductID: product.ID,
			Stock:     product.Stock,
		})
	}
	if d.Status != models.PENDING {
		return ErrOrderConflict
	}
	for i, v := range d.Item {
		if err := s.ProductServices.Sell(v.ProductID, v.Qty); err != nil {
			// semisal ada error, rollback semua (tolak transaksi awal-akhir)
			for j := 0; j < i; j++ {
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
	err = s.OrderRepo.Update(d)
	if err == repository.ErrNotFound {
		return ErrOrderNotFound
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *OrderServiceImpl) CancelOrder(id string) error {
	if id == "" {
		return ErrOrderInvalid
	}
	d, err := s.OrderRepo.FindByID(id)
	if err == repository.ErrNotFound {
		return ErrOrderNotFound
	}
	if err != nil {
		return err
	}
	if d.Status == models.PENDING {
		d.Status = models.CANCELED
		d.UpdatedAt = time.Now().UTC()
	} else {
		return ErrOrderConflict
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
	err = s.OrderRepo.Update(d)
	if err == repository.ErrNotFound {
		return ErrOrderNotFound
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *OrderServiceImpl) DeleteOrder(id string) error {
	if id == "" {
		return ErrOrderInvalid
	}
	orderData, err := s.OrderRepo.FindByID(id)
	if err == repository.ErrNotFound {
		return ErrOrderNotFound
	}
	if err != nil {
		return err
	}
	if orderData.Status == models.PAID {
		return ErrOrderConflict
	}
	if errDelete := s.OrderRepo.Delete(orderData.ID); errDelete != nil {
		if errDelete == repository.ErrNotFound {
			return ErrOrderNotFound
		}
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
