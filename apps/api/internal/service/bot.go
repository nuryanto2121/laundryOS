package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nuryanto2121/laundryos/apps/api/internal/bot"
	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
	"github.com/nuryanto2121/laundryos/apps/api/internal/repository"
)

// BotService defines the contract for processing customer chat messages.
type BotService interface {
	ProcessMessage(ctx context.Context, req model.IncomingMessageRequest) (*model.BotResponse, error)
}

// BotServiceConfig bundles all dependencies for BotService.
type BotServiceConfig struct {
	CustomerService *CustomerService
	OrderService    *OrderService
	ServiceRepo     *repository.ServiceRepository
	SessionStore    bot.SessionStore
	NotificationSvc NotificationService
	GatewayURL      string
}

type botService struct {
	customerSvc  *CustomerService
	orderSvc     *OrderService
	serviceRepo  *repository.ServiceRepository
	sessionStore bot.SessionStore
	notifSvc     NotificationService
	matcher      *bot.Matcher
	gatewayURL   string
	httpClient   *http.Client
}

// NewBotService creates a new BotService instance.
func NewBotService(cfg BotServiceConfig) BotService {
	sessionStore := cfg.SessionStore
	if sessionStore == nil {
		sessionStore = bot.NewMemorySessionStore(15 * time.Minute)
	}

	matcher := bot.NewMatcher(bot.GetKnowledgeBase())
	return &botService{
		customerSvc:  cfg.CustomerService,
		orderSvc:     cfg.OrderService,
		serviceRepo:  cfg.ServiceRepo,
		sessionStore: sessionStore,
		notifSvc:     cfg.NotificationSvc,
		matcher:      matcher,
		gatewayURL:   cfg.GatewayURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// ProcessMessage matches or handles multi-turn conversation states, formats the reply,
// and asynchronously dispatches the response to the WhatsApp Gateway.
func (s *botService) ProcessMessage(ctx context.Context, req model.IncomingMessageRequest) (*model.BotResponse, error) {
	text := strings.TrimSpace(req.Text)
	phone := strings.TrimSpace(req.From)

	// 1. Check if user wants to cancel an active conversation
	if strings.EqualFold(text, "BATAL") || strings.EqualFold(text, "CANCEL") {
		if _, active := s.sessionStore.Get(phone); active {
			s.sessionStore.Delete(phone)
			reply := "❌ *Permintaan penjemputan dibatalkan.*\n\nKetik *MENU* jika membutuhkan info lainnya, atau ketik *PICKUP* kapan saja untuk mulai kembali. 🙏"
			s.dispatchAsync(phone, reply)
			return &model.BotResponse{Category: bot.CategoryPickup, Reply: reply}, nil
		}
	}

	// 2. Check for an active multi-turn session
	if session, active := s.sessionStore.Get(phone); active && session.State != bot.StateIdle {
		return s.handleActiveSession(ctx, session, req)
	}

	// 3. Not in a session — perform standard FAQ keyword matching
	category, score := s.matcher.Match(text)

	if category == bot.CategoryPickup {
		return s.startPickupFlow(ctx, req)
	}

	if category == bot.CategoryStatus {
		return s.handleCheckStatus(ctx, req)
	}

	// Audit log for unrecognized queries
	if category == bot.CategoryUnknown {
		slog.Warn("bot: unrecognized customer query",
			"from", phone,
			"push_name", req.PushName,
			"text", text,
			"score", score,
		)
	} else {
		slog.Info("bot: matched FAQ category",
			"from", phone,
			"category", category,
			"score", score,
		)
	}

	var services []model.ServiceItem
	if category == bot.CategoryPrice && s.serviceRepo != nil {
		activeServices, err := s.serviceRepo.ListActiveServices(ctx)
		if err != nil {
			slog.Warn("bot: failed to load active services from DB, falling back to default", "error", err)
		} else {
			services = activeServices
		}
	}

	reply := bot.FormatResponseWithServices(category, req.PushName, services)
	s.dispatchAsync(phone, reply)

	return &model.BotResponse{
		Category: category,
		Reply:    reply,
	}, nil
}

func (s *botService) handleCheckStatus(ctx context.Context, req model.IncomingMessageRequest) (*model.BotResponse, error) {
	phone := strings.TrimSpace(req.From)
	text := strings.ToUpper(strings.TrimSpace(req.Text))

	var orderDetail model.OrderTrackingDetail
	var err error

	// 1. Check if user specified an explicit order number in their message (e.g. "STATUS ORD261118000001")
	words := strings.Fields(text)
	var specifiedOrderNumber string
	for _, w := range words {
		clean := strings.Trim(w, `.,!?:;"'#`)
		if strings.HasPrefix(clean, "ORD") && len(clean) >= 10 {
			specifiedOrderNumber = clean
			break
		}
	}

	if specifiedOrderNumber != "" && s.orderSvc != nil {
		orderDetail, err = s.orderSvc.GetByOrderNumber(ctx, specifiedOrderNumber)
		if err != nil {
			slog.Warn("bot: failed to find order by specified number", "order_number", specifiedOrderNumber, "error", err)
		}
	}

	// 2. If no explicit order number provided or not found, lookup the customer's latest active order
	if orderDetail.OrderNumber == "" && s.orderSvc != nil && phone != "" {
		orderDetail, err = s.orderSvc.GetLatestByCustomerPhone(ctx, phone)
		if err != nil {
			slog.Warn("bot: failed to find order by customer phone", "phone", phone, "error", err)
		}
	}

	reply := bot.FormatOrderStatusResponse(&orderDetail, req.PushName)
	s.dispatchAsync(phone, reply)

	return &model.BotResponse{
		Category: bot.CategoryStatus,
		Reply:    reply,
	}, nil
}

func (s *botService) startPickupFlow(ctx context.Context, req model.IncomingMessageRequest) (*model.BotResponse, error) {
	phone := req.From
	name := strings.TrimSpace(req.PushName)
	if name == "" {
		name = "Pelanggan"
	}

	session := &bot.SessionData{
		Phone:        phone,
		CustomerName: name,
		State:        bot.StateAwaitingAddress,
	}

	var existingAddress string
	if s.customerSvc != nil {
		cust, err := s.customerSvc.GetByPhone(ctx, phone)
		if err == nil && cust.Address != "" {
			existingAddress = cust.Address
			session.Address = existingAddress
		}
	}

	s.sessionStore.Set(phone, session)

	var reply string
	if existingAddress != "" {
		reply = fmt.Sprintf(`🛵 *REQUEST PENJEMPUTAN LAUNDRYOS* 🛵

Halo Kak *%s*! Alamat penjemputan Anda yang tersimpan:
📍 *%s*

Apakah ingin dijemput di alamat ini?
• Ketik *YA* jika menggunakan alamat tersebut
• Atau langsung ketik *Alamat Baru* Anda

_(Ketik *BATAL* untuk membatalkan)_`, name, existingAddress)
	} else {
		reply = fmt.Sprintf(`🛵 *REQUEST PENJEMPUTAN LAUNDRYOS* 🛵

Halo Kak *%s*! Kurir kami siap menjemput pakaian kotor langsung ke tempat Anda.

Silakan ketik *Alamat Lengkap Penjemputan* Anda (Jalan, No Rumah, Kos/Apartemen, Patokan):

_(Ketik *BATAL* untuk membatalkan)_`, name)
	}

	s.dispatchAsync(phone, reply)
	return &model.BotResponse{Category: bot.CategoryPickup, Reply: reply}, nil
}

func (s *botService) handleActiveSession(ctx context.Context, session *bot.SessionData, req model.IncomingMessageRequest) (*model.BotResponse, error) {
	text := strings.TrimSpace(req.Text)
	phone := req.From

	switch session.State {
	case bot.StateAwaitingAddress:
		if !strings.EqualFold(text, "YA") || session.Address == "" {
			session.Address = text
		}

		// Fetch services for selection
		var serviceListText string
		if s.serviceRepo != nil {
			services, err := s.serviceRepo.ListActiveServices(ctx)
			if err == nil && len(services) > 0 {
				var sb strings.Builder
				for i, s := range services {
					durationStr := ""
					if s.DurationDay > 0 {
						durationStr = fmt.Sprintf(" (Selesai %d hari)", s.DurationDay)
					}
					sb.WriteString(fmt.Sprintf("%d. *%s* - Rp %s/Kg%s\n", i+1, s.ServiceName, formatIDR(s.PricePerKg), durationStr))
				}
				serviceListText = sb.String()
			}
		}

		if serviceListText == "" {
			serviceListText = "1. *Cuci Komplit Reguler* - Rp 8.000/Kg (Selesai 2 hari)\n2. *Cuci Kilat 24 Jam* - Rp 10.000/Kg\n3. *Cuci Express 6 Jam* - Rp 15.000/Kg\n"
		}

		session.State = bot.StateAwaitingService
		s.sessionStore.Set(phone, session)

		reply := fmt.Sprintf(`Alamat penjemputan dicatat:
📍 *%s*

Silakan pilih paket layanan laundry:
%s
Ketik angka pilihan Anda (contoh: *1*):`, session.Address, serviceListText)

		s.dispatchAsync(phone, reply)
		return &model.BotResponse{Category: bot.CategoryPickup, Reply: reply}, nil

	case bot.StateAwaitingService:
		var activeServices []model.ServiceItem
		if s.serviceRepo != nil {
			svcs, err := s.serviceRepo.ListActiveServices(ctx)
			if err == nil {
				activeServices = svcs
			}
		}

		idx, err := strconv.Atoi(text)
		if err != nil || idx < 1 || idx > len(activeServices) {
			// If services table empty fallback check
			if len(activeServices) == 0 && (idx == 1 || idx == 2 || idx == 3) {
				session.ServiceID = "a0000000-0000-0000-0000-000000000001"
				session.ServiceName = "Cuci Komplit Reguler"
			} else {
				reply := fmt.Sprintf("Pilihan nomor tidak valid. Silakan ketik angka antara 1 sampai %d:", len(activeServices))
				s.dispatchAsync(phone, reply)
				return &model.BotResponse{Category: bot.CategoryPickup, Reply: reply}, nil
			}
		} else {
			selected := activeServices[idx-1]
			session.ServiceID = selected.ID
			session.ServiceName = selected.ServiceName
		}

		session.State = bot.StateAwaitingSchedule
		s.sessionStore.Set(phone, session)

		reply := fmt.Sprintf(`Layanan dipilih: *%s* ✨

Kapan kurir kami bisa menjemput cucian Anda?
(Ketik *SEKARANG* atau sebutkan jam, contoh: *14:00 WIB*):`, session.ServiceName)

		s.dispatchAsync(phone, reply)
		return &model.BotResponse{Category: bot.CategoryPickup, Reply: reply}, nil

	case bot.StateAwaitingSchedule:
		session.Schedule = text
		session.State = bot.StateAwaitingConfirmation
		s.sessionStore.Set(phone, session)

		reply := fmt.Sprintf(`🛵 *RINGKASAN REQUEST PENJEMPUTAN* 🛵

• *Nama:* %s
• *No. WhatsApp:* %s
• *Alamat Jemput:* %s
• *Paket Layanan:* %s
• *Waktu Pickup:* %s

Ketik *OK* untuk konfirmasi penjemputan, atau ketik *BATAL* untuk membatalkan.`,
			session.CustomerName, phone, session.Address, session.ServiceName, session.Schedule,
		)

		s.dispatchAsync(phone, reply)
		return &model.BotResponse{Category: bot.CategoryPickup, Reply: reply}, nil

	case bot.StateAwaitingConfirmation:
		confirmed := strings.EqualFold(text, "OK") || strings.EqualFold(text, "YA") ||
			strings.EqualFold(text, "SIAP") || strings.EqualFold(text, "DEAL")

		if !confirmed {
			reply := "Ketik *OK* untuk konfirmasi pemesanan, atau ketik *BATAL* jika ingin membatalkan."
			s.dispatchAsync(phone, reply)
			return &model.BotResponse{Category: bot.CategoryPickup, Reply: reply}, nil
		}

		// 1. Auto-register or update customer
		var customerID string
		if s.customerSvc != nil {
			cust, err := s.customerSvc.GetOrCreateByPhone(ctx, phone, session.CustomerName, session.Address)
			if err != nil {
				slog.Error("bot: failed to auto-register customer", "error", err)
			} else {
				customerID = cust.ID
			}
		}

		// 2. Create Order with WAITING_PICKUP status
		orderNumber := "ORD" + time.Now().Format("060102150405")
		if s.orderSvc != nil && customerID != "" {
			serviceID := session.ServiceID
			if serviceID == "" {
				serviceID = "a0000000-0000-0000-0000-000000000001"
			}
			order, err := s.orderSvc.CreatePickupOrder(ctx, model.CreateOrderRequest{
				CustomerID:     customerID,
				ServiceID:      serviceID,
				PickupRequired: true,
				Notes:          fmt.Sprintf("Jadwal: %s | Alamat: %s", session.Schedule, session.Address),
			})
			if err != nil {
				slog.Error("bot: failed to create pickup order", "error", err)
			} else {
				orderNumber = order.OrderNumber
			}
		}

		// 3. Clear session
		s.sessionStore.Delete(phone)

		// 4. Notify Staff / Courier Group
		if s.notifSvc != nil {
			go func() {
				_ = s.notifSvc.NotifyPickupRequest(
					context.Background(),
					orderNumber,
					session.CustomerName,
					phone,
					session.Address,
					session.ServiceName,
					session.Schedule,
				)
			}()
		}

		reply := fmt.Sprintf(`🎉 *REQUEST PICKUP BERHASIL DIBUAT!* 🎉

No. Order: *%s*
Status: *MENUNGGU PENJEMPUTAN (WAITING_PICKUP)*

Kurir LaundryOS telah menerima pesanan ini dan akan segera menjemput pakaian kotor Anda sesuai jadwal (*%s*).

Terima kasih atas kepercayaan Anda! 🧺✨`, orderNumber, session.Schedule)

		s.dispatchAsync(phone, reply)
		return &model.BotResponse{Category: bot.CategoryPickup, Reply: reply}, nil

	default:
		s.sessionStore.Delete(phone)
		return &model.BotResponse{Category: bot.CategoryUnknown, Reply: "Sesi telah direset."}, nil
	}
}

func (s *botService) dispatchAsync(phone, message string) {
	if s.gatewayURL == "" || phone == "" {
		return
	}
	go s.dispatchReply(phone, message)
}

func (s *botService) dispatchReply(phone, message string) {
	payload := map[string]string{
		"phone":   phone,
		"message": message,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		slog.Error("bot: failed to marshal gateway dispatch payload", "error", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, s.gatewayURL+"/send-message", bytes.NewReader(body))
	if err != nil {
		slog.Error("bot: failed to construct gateway dispatch request", "error", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		slog.Error("bot: gateway dispatch failed", "error", err, "gateway_url", s.gatewayURL)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		slog.Warn("bot: gateway responded with non-2xx status", "status", resp.StatusCode)
	} else {
		slog.Info("bot: reply dispatched to gateway successfully", "phone", phone)
	}
}

func formatIDR(amount float64) string {
	n := int64(amount)
	in := fmt.Sprintf("%d", n)
	if len(in) <= 3 {
		return in
	}
	var res []string
	for len(in) > 3 {
		res = append([]string{in[len(in)-3:]}, res...)
		in = in[:len(in)-3]
	}
	if len(in) > 0 {
		res = append([]string{in}, res...)
	}
	return strings.Join(res, ".")
}
