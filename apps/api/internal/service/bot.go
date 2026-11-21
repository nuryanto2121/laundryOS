package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sort"
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
	CustomerRepo    *repository.CustomerRepository
	OrderService    *OrderService
	ServiceRepo     *repository.ServiceRepository
	MasterRepo      *repository.MasterRepository
	SessionStore    bot.SessionStore
	NotificationSvc NotificationService
	GatewayURL      string
	StaffPhone      string
}

type botService struct {
	customerSvc  *CustomerService
	customerRepo *repository.CustomerRepository
	orderSvc     *OrderService
	serviceRepo  *repository.ServiceRepository
	masterRepo   *repository.MasterRepository
	sessionStore bot.SessionStore
	notifSvc     NotificationService
	matcher      *bot.Matcher
	gatewayURL   string
	staffPhone   string
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
		customerRepo: cfg.CustomerRepo,
		orderSvc:     cfg.OrderService,
		serviceRepo:  cfg.ServiceRepo,
		masterRepo:   cfg.MasterRepo,
		sessionStore: sessionStore,
		notifSvc:     cfg.NotificationSvc,
		matcher:      matcher,
		gatewayURL:   cfg.GatewayURL,
		staffPhone:   strings.TrimSpace(cfg.StaffPhone),
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// ProcessMessage matches or handles multi-turn conversation states, formats the reply,
// and asynchronously dispatches the response to the WhatsApp Gateway.
func (s *botService) ProcessMessage(ctx context.Context, req model.IncomingMessageRequest) (*model.BotResponse, error) {
	// 0. Handle Group Messages (e.g. Courier/Staff claiming an order from the staff group)
	if req.IsGroup {
		return s.handleGroupMessage(ctx, req)
	}

	text := strings.TrimSpace(req.Text)
	phone := strings.TrimSpace(req.From)

	// Extract customer name if they introduced themselves (e.g. "saya nana teman'a nini")
	if extractedName := bot.ExtractCustomerName(text); extractedName != "" {
		req.PushName = extractedName
	}

	// 1. Check if user wants to cancel an active conversation
	if strings.EqualFold(text, "BATAL") || strings.EqualFold(text, "CANCEL") {
		if _, active := s.sessionStore.Get(phone); active {
			s.sessionStore.Delete(phone)
			reply := "❌ *Permintaan penjemputan dibatalkan.*\n\nKetik *MENU* atau angka *1-4* jika membutuhkan info lainnya, atau kirim pesan kapan saja untuk mulai kembali. 🙏"
			s.dispatchAsync(phone, reply)
			return &model.BotResponse{Category: bot.CategoryPickup, Reply: reply}, nil
		}
	}

	// Fetch active services (if available)
	var activeServices []model.ServiceItem
	if s.serviceRepo != nil {
		if svcs, err := s.serviceRepo.ListActiveServices(ctx); err == nil {
			activeServices = sortServicesByCategory(svcs)
		}
	}

	// 2. Check for an active multi-turn session
	if session, active := s.sessionStore.Get(phone); active && session.State != bot.StateIdle {
		return s.handleActiveSession(ctx, session, req, activeServices)
	}

	// 3. Check if customer is selecting a service directly (e.g. "cuci regular iah", "mau kilat", "express ya")
	if matchedSvc := bot.MatchServiceIntent(text, activeServices); matchedSvc != nil {
		return s.handleServiceSelection(ctx, req, matchedSvc)
	}

	// 4. Not in a session & not picking a service — perform standard FAQ keyword matching
	category, score := s.matcher.Match(text)

	if category == bot.CategoryPickup {
		return s.startPickupFlow(ctx, req)
	}

	if category == bot.CategoryStatus || bot.ExtractOrderNumber(req.Text) != "" {
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
	if (category == bot.CategoryPrice || category == bot.CategoryService) && len(activeServices) > 0 {
		services = activeServices
	}

	reply := bot.FormatResponseWithServices(category, req.PushName, services)
	s.dispatchAsync(phone, reply)

	return &model.BotResponse{
		Category: category,
		Reply:    reply,
	}, nil
}

func (s *botService) handleServiceSelection(ctx context.Context, req model.IncomingMessageRequest, svc *model.ServiceItem) (*model.BotResponse, error) {
	phone := req.From
	name := strings.TrimSpace(req.PushName)
	if extracted := bot.ExtractCustomerName(req.Text); extracted != "" {
		name = extracted
	}
	if name == "" {
		name = "Pelanggan"
	}

	session := &bot.SessionData{
		Phone:        phone,
		CustomerName: name,
		ServiceName:  svc.ServiceName,
		ServiceID:    svc.ID,
		LastActive:   time.Now(),
	}

	var existingAddress string
	if s.customerSvc != nil {
		cust, err := s.customerSvc.GetByPhone(ctx, phone)
		if err == nil && cust.Address != "" {
			existingAddress = cust.Address
			session.Address = existingAddress
		}
	}

	if isSelfService(svc.ServiceName) {
		session.State = bot.StateAwaitingDetergent
		s.sessionStore.Set(phone, session)
		reply := fmt.Sprintf("Siap Kak *%s*, pesanan *%s* (Rp %s/Tabung max 8 Kg) dicatat! 🧺✨\n\n"+
			"Apakah Kakak sudah menyiapkan sabun & pewangi sendiri, atau ingin sekalian beli di FAN Laundry?\n"+
			"• *Sabun / Pewangi:* Rp 1.000 / Pcs\n\n"+
			"👉 *Pilihan:*\n"+
			"1️⃣ Ketik *1* jika *Bawa Sendiri*\n"+
			"2️⃣ Ketik *2* jika *Beli Sabun & Pewangi* (1 Sabun + 1 Pewangi = Rp 2.000)\n"+
			"_(Atau sebutkan kebutuhan Anda, contoh: 'Beli 2 sabun 2 pewangi')_\n\n"+
			"_(Ketik *BATAL* untuk membatalkan)_", name, svc.ServiceName, formatIDR(svc.PricePerKg))
		s.dispatchAsync(phone, reply)
		return &model.BotResponse{Category: bot.CategoryPickup, Reply: reply}, nil
	}

	var reply string
	if existingAddress != "" {
		session.State = bot.StateAwaitingSchedule
		s.sessionStore.Set(phone, session)
		reply = fmt.Sprintf(`Siap Kak *%s*, pesanan *%s* (Rp %s/Kg) dicatat! 🧺✨

Alamat penjemputan Anda yang tersimpan:
📍 *%s*

Mau dijemput jam berapa nih Kak?
(Contoh: *Sekarang*, *Nanti jam 2 siang*, atau *Besok pagi*):

_(Ketik *BATAL* untuk membatalkan)_`, name, svc.ServiceName, formatIDR(svc.PricePerKg), existingAddress)
	} else {
		session.State = bot.StateAwaitingAddress
		s.sessionStore.Set(phone, session)
		reply = fmt.Sprintf(`Siap Kak *%s*, pesanan *%s* (Rp %s/Kg) dicatat! 🧺✨

Mau dijemput di mana dan jam berapa nih Kak?
Silakan kirim *Alamat Lengkap Penjemputan* Anda:
(Contoh: _Jl. Mawar No. 12, jemput jam 2 siang ya_)

_(Ketik *BATAL* untuk membatalkan)_`, name, svc.ServiceName, formatIDR(svc.PricePerKg))
	}

	s.dispatchAsync(phone, reply)
	return &model.BotResponse{Category: bot.CategoryPickup, Reply: reply}, nil
}

func (s *botService) handleCheckStatus(ctx context.Context, req model.IncomingMessageRequest) (*model.BotResponse, error) {
	phone := strings.TrimSpace(req.From)
	text := strings.ToUpper(strings.TrimSpace(req.Text))

	var orderDetail model.OrderTrackingDetail
	var err error

	// 1. Check if user specified an explicit order number in their message
	// Supporting "#261121000001", "261121000001", "ORD261121000001"
	specifiedOrderNumber := bot.ExtractOrderNumber(req.Text)
	if specifiedOrderNumber == "" {
		words := strings.Fields(text)
		for _, w := range words {
			clean := strings.Trim(w, `.,!?:;"'#`)
			if strings.HasPrefix(clean, "ORD") && len(clean) >= 10 {
				specifiedOrderNumber = clean
				break
			}
			// Check for 10-14 digit numeric order number
			if len(clean) >= 10 && len(clean) <= 14 {
				isDigits := true
				for _, r := range clean {
					if r < '0' || r > '9' {
						isDigits = false
						break
					}
				}
				if isDigits {
					specifiedOrderNumber = "ORD" + clean
					break
				}
			}
		}
	}

	if specifiedOrderNumber != "" && s.orderSvc != nil {
		orderDetail, err = s.orderSvc.GetByOrderNumber(ctx, specifiedOrderNumber)
		if err != nil || orderDetail.OrderNumber == "" {
			// Also try without "ORD" prefix if specified had it
			rawNum := strings.TrimPrefix(specifiedOrderNumber, "ORD")
			if altDetail, altErr := s.orderSvc.GetByOrderNumber(ctx, rawNum); altErr == nil && altDetail.OrderNumber != "" {
				orderDetail = altDetail
			}
		}
		if err != nil && orderDetail.OrderNumber == "" {
			slog.Warn("bot: failed to find order by specified number", "order_number", specifiedOrderNumber, "error", err)
		}
	}

	// 2. If user explicitly specified an order number BUT it was not found:
	// Do NOT silently fall back to an unrelated previous order! Tell the customer clearly.
	if specifiedOrderNumber != "" && orderDetail.OrderNumber == "" {
		dispNum := bot.FormatDisplayOrderNumber(specifiedOrderNumber)
		name := req.PushName
		if name == "" {
			name = "Pelanggan"
		}
		reply := fmt.Sprintf("🧺 *STATUS CUCIAN (FAN Laundry)* 🧺\n\nMaaf Kak *%s*, order dengan nomor *%s* tidak ditemukan di sistem kami.\n\nMohon periksa kembali nomor order pada nota Anda, atau ketik *STATUS* untuk mengecek cucian aktif terakhir Anda. 🙏\n\nKetik *MENU* untuk pilihan info lainnya.", name, dispNum)
		s.dispatchAsync(phone, reply)
		return &model.BotResponse{
			Category: bot.CategoryStatus,
			Reply:    reply,
		}, nil
	}

	// 3. If no explicit order number provided, lookup the customer's latest active order by phone
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
	if extracted := bot.ExtractCustomerName(req.Text); extracted != "" {
		name = extracted
	}
	if name == "" {
		name = "Pelanggan"
	}

	session := &bot.SessionData{
		Phone:        phone,
		CustomerName: name,
		State:        bot.StateAwaitingAddress,
		LastActive:   time.Now(),
	}

	var existingAddress string
	var existingCustomerID string
	if s.customerSvc != nil {
		cust, err := s.customerSvc.GetByPhone(ctx, phone)
		if err == nil && cust.ID != "" {
			existingCustomerID = cust.ID
			if cust.Address != "" {
				existingAddress = cust.Address
				session.Address = existingAddress
			}
		}
	}

	// 1. Fetch active outlets dynamically
	var activeOutlets []model.Outlet
	if s.masterRepo != nil {
		all, _ := s.masterRepo.ListOutletsWithGPS(ctx)
		for _, o := range all {
			if o.IsActive {
				activeOutlets = append(activeOutlets, o)
			}
		}
	}

	// Determine dynamic brand name (e.g. "FAN Laundry" or "LaundryOS")
	brandName := "LaundryOS"
	if len(activeOutlets) > 0 && activeOutlets[0].Name != "" {
		brandName = activeOutlets[0].Name
	}

	// 2. Hybrid Branch Detection:
	// A. Check if customer already has a known branch preference
	var matchedOutlet *model.Outlet
	if existingCustomerID != "" && s.customerRepo != nil {
		preferredOutletID := s.customerRepo.GetCustomerDefaultOutletID(ctx, existingCustomerID)
		if preferredOutletID != "" {
			for _, o := range activeOutlets {
				if o.ID == preferredOutletID {
					matchedOutlet = &o
					break
				}
			}
		}
	}

	// If only 1 outlet exists in the system, automatically assign it
	if matchedOutlet == nil && len(activeOutlets) == 1 {
		matchedOutlet = &activeOutlets[0]
	}

	// B. If branch is already matched, proceed straight to Address!
	if matchedOutlet != nil {
		session.OutletID = matchedOutlet.ID
		session.OutletName = matchedOutlet.Name
		session.TargetGroupJID = matchedOutlet.WhatsAppGroupID
		session.State = bot.StateAwaitingAddress
		s.sessionStore.Set(phone, session)

		var reply string
		if existingAddress != "" {
			reply = fmt.Sprintf(`🛵 *REQUEST PENJEMPUTAN %s* 🛵

Halo Kak *%s*! Cabang: *%s*
Alamat penjemputan Anda yang tersimpan:
📍 *%s*

Apakah ingin dijemput di alamat ini?
• Ketik *YA* jika menggunakan alamat tersebut
• Atau langsung ketik *Alamat Baru* Anda

_(Ketik *BATAL* untuk membatalkan)_`, strings.ToUpper(session.OutletName), name, session.OutletName, existingAddress)
		} else {
			reply = fmt.Sprintf(`🛵 *REQUEST PENJEMPUTAN %s* 🛵

Halo Kak *%s*! Cabang: *%s*
Kurir kami siap meluncur menjemput pakaian kotor Anda. 🧺✨

Silakan ketik *Alamat Lengkap Penjemputan* Anda (Jalan, No Rumah, Kos/Apartemen, Patokan):
_(Atau ketik *1* jika ingin cek daftar harga terlebih dahulu)_

_(Ketik *BATAL* untuk membatalkan)_`, strings.ToUpper(session.OutletName), name, session.OutletName)
		}

		s.dispatchAsync(phone, reply)
		return &model.BotResponse{Category: bot.CategoryPickup, Reply: reply}, nil
	}

	// C. If customer has NO branch and there are multiple branches: Ask to choose!
	if len(activeOutlets) > 1 {
		session.State = bot.StateAwaitingBranchSelection
		s.sessionStore.Set(phone, session)

		reply := fmt.Sprintf("🛵 *REQUEST PENJEMPUTAN %s* 🛵\n\n"+
			"Halo Kak *%s*! Kurir kami siap meluncur menjemput pakaian kotor Anda. 🧺✨\n\n"+
			"Silakan pilih cabang laundry terdekat untuk penjemputan:\n",
			strings.ToUpper(brandName), name,
		)
		for i, o := range activeOutlets {
			addrInfo := ""
			if o.Address != "" {
				addrInfo = " (" + o.Address + ")"
			}
			reply += fmt.Sprintf("\n%d️⃣ *%s*%s", i+1, o.Name, addrInfo)
		}
		reply += "\n\n👉 *Ketik angka pilihan cabang Anda* (contoh: *1*):\n_(Ketik *BATAL* untuk membatalkan)_"

		s.dispatchAsync(phone, reply)
		return &model.BotResponse{Category: bot.CategoryPickup, Reply: reply}, nil
	}

	// Fallback if no outlets registered yet
	s.sessionStore.Set(phone, session)
	reply := fmt.Sprintf(`🛵 *REQUEST PENJEMPUTAN %s* 🛵

Halo Kak *%s*! Kurir kami siap meluncur menjemput pakaian kotor Anda. 🧺✨

Silakan ketik *Alamat Lengkap Penjemputan* Anda:
_(Ketik *BATAL* untuk membatalkan)_`, strings.ToUpper(brandName), name)

	s.dispatchAsync(phone, reply)
	return &model.BotResponse{Category: bot.CategoryPickup, Reply: reply}, nil
}

func (s *botService) handleActiveSession(ctx context.Context, session *bot.SessionData, req model.IncomingMessageRequest, activeServices []model.ServiceItem) (*model.BotResponse, error) {
	text := strings.TrimSpace(req.Text)
	phone := req.From

	// If customer explicitly introduced name, update it
	if extracted := bot.ExtractCustomerName(text); extracted != "" {
		session.CustomerName = extracted
	}

	switch session.State {
	case bot.StateAwaitingBranchSelection:
		var activeOutlets []model.Outlet
		if s.masterRepo != nil {
			all, _ := s.masterRepo.ListOutletsWithGPS(ctx)
			for _, o := range all {
				if o.IsActive {
					activeOutlets = append(activeOutlets, o)
				}
			}
		}

		idx, err := strconv.Atoi(text)
		if err != nil || idx < 1 || idx > len(activeOutlets) {
			reply := fmt.Sprintf("Pilihan tidak valid. Silakan ketik angka antara 1 sampai %d sesuai cabang terdekat Anda:\n_(Ketik *BATAL* untuk membatalkan)_", len(activeOutlets))
			s.dispatchAsync(phone, reply)
			return &model.BotResponse{Category: bot.CategoryPickup, Reply: reply}, nil
		}

		chosen := activeOutlets[idx-1]
		session.OutletID = chosen.ID
		session.OutletName = chosen.Name
		session.TargetGroupJID = chosen.WhatsAppGroupID
		session.State = bot.StateAwaitingAddress

		// Persist customer preferred outlet for future orders
		if s.customerSvc != nil && s.customerRepo != nil {
			cust, err := s.customerSvc.GetByPhone(ctx, phone)
			if err == nil && cust.ID != "" {
				_ = s.customerRepo.SetCustomerOutletID(ctx, cust.ID, chosen.ID)
			}
		}

		s.sessionStore.Set(phone, session)

		reply := fmt.Sprintf("Cabang dipilih: *%s* ✅\n\n"+
			"Silakan ketik *Alamat Lengkap Penjemputan* Anda (Jalan, No Rumah, Kos/Apartemen, Patokan):\n"+
			"_(Ketik *BATAL* untuk membatalkan)_", chosen.Name)
		s.dispatchAsync(phone, reply)
		return &model.BotResponse{Category: bot.CategoryPickup, Reply: reply}, nil

	case bot.StateAwaitingAddress:
		// If user asks "Harga" or "1" while being asked for address, show price and re-prompt address without resetting
		if strings.EqualFold(text, "1") || strings.EqualFold(text, "HARGA") {
			priceResp := bot.FormatResponseWithServices(bot.CategoryPrice, session.CustomerName, activeServices)
			reply := priceResp + "\n\nSilakan ketik *Alamat Penjemputan* Anda untuk melanjutkan:\n_(Ketik *BATAL* untuk membatalkan)_"
			s.dispatchAsync(phone, reply)
			return &model.BotResponse{Category: bot.CategoryPrice, Reply: reply}, nil
		}

		// Check if user is picking a service in this step instead of giving address
		if svc := bot.MatchServiceIntent(text, activeServices); svc != nil {
			session.ServiceName = svc.ServiceName
			session.ServiceID = svc.ID

			var reply string
			if session.Address != "" {
				session.State = bot.StateAwaitingSchedule
				s.sessionStore.Set(phone, session)
				reply = fmt.Sprintf(`Siap Kak *%s*, pesanan *%s* (Rp %s/Kg) dicatat! 🧺✨

Alamat penjemputan:
📍 *%s*

Mau dijemput jam berapa nih Kak?
(Contoh: *Sekarang*, *Nanti jam 2 siang*, atau *Besok pagi*):

_(Ketik *BATAL* untuk membatalkan)_`, session.CustomerName, session.ServiceName, formatIDR(svc.PricePerKg), session.Address)
			} else {
				s.sessionStore.Set(phone, session)
				reply = fmt.Sprintf(`Siap Kak *%s*, pesanan *%s* (Rp %s/Kg) dicatat! 🧺✨

Boleh share *Alamat Lengkap Penjemputan* Anda Kak?
(Contoh: _Jl. Mawar No. 12, seberang masjid_)

_(Ketik *BATAL* untuk membatalkan)_`, session.CustomerName, session.ServiceName, formatIDR(svc.PricePerKg))
			}
			s.dispatchAsync(phone, reply)
			return &model.BotResponse{Category: bot.CategoryPickup, Reply: reply}, nil
		}

		if !strings.EqualFold(text, "YA") || session.Address == "" {
			session.Address = text
		}

		// If service is already known, advance straight to schedule!
		if session.ServiceName != "" {
			session.State = bot.StateAwaitingSchedule
			s.sessionStore.Set(phone, session)

			reply := fmt.Sprintf(`Alamat penjemputan dicatat:
📍 *%s*

Paket: *%s* ✨

Mau dijemput jam berapa nih Kak?
(Contoh: *Sekarang*, *Nanti jam 2 siang*, atau *Besok pagi*):

_(Ketik *BATAL* untuk membatalkan)_`, session.Address, session.ServiceName)
			s.dispatchAsync(phone, reply)
			return &model.BotResponse{Category: bot.CategoryPickup, Reply: reply}, nil
		}

		// Otherwise, prompt for service selection
		var serviceListText string
		if len(activeServices) > 0 {
			serviceListText = formatGroupedServiceMenu(activeServices)
		}

		if serviceListText == "" {
			serviceListText = "\n🧺 *CUCI SETRIKA*\n1. *Cuci Setrika Reguler* - Rp 7.000/Kg (3 hari)\n2. *Cuci Setrika Express* - Rp 9.000/Kg (2 hari)\n3. *Cuci Setrika Kilat* - Rp 11.000/Kg (1 hari)\n"
		}

		session.State = bot.StateAwaitingService
		s.sessionStore.Set(phone, session)

		reply := fmt.Sprintf("Alamat penjemputan dicatat:\n"+
			"📍 *%s*\n\n"+
			"Silakan pilih paket layanan laundry:\n"+
			"%s\n"+
			"Ketik angka pilihan Anda atau sebutkan nama paketnya (contoh: *1* atau *'Cuci Reguler'*):", session.Address, serviceListText)

		s.dispatchAsync(phone, reply)
		return &model.BotResponse{Category: bot.CategoryPickup, Reply: reply}, nil

	case bot.StateAwaitingService:
		// Accept both numeric choice AND natural text name
		var matchedService *model.ServiceItem
		idx, err := strconv.Atoi(text)
		if err == nil && idx >= 1 && idx <= len(activeServices) {
			item := activeServices[idx-1]
			matchedService = &item
		} else if err == nil && len(activeServices) == 0 && (idx == 1 || idx == 2 || idx == 3) {
			matchedService = &model.ServiceItem{
				ID:          "a0000000-0000-0000-0000-000000000001",
				ServiceName: "Cuci Komplit Reguler",
			}
		} else {
			// Check natural text
			matchedService = bot.MatchServiceIntent(text, activeServices)
		}

		if matchedService == nil {
			reply := fmt.Sprintf("Paket tidak dikenali. Silakan ketik angka antara 1 sampai %d, atau sebutkan nama paketnya (contoh: *'Reguler'*):", len(activeServices))
			s.dispatchAsync(phone, reply)
			return &model.BotResponse{Category: bot.CategoryPickup, Reply: reply}, nil
		}

		session.ServiceID = matchedService.ID
		session.ServiceName = matchedService.ServiceName

		if isSelfService(matchedService.ServiceName) {
			session.State = bot.StateAwaitingDetergent
			s.sessionStore.Set(phone, session)
			reply := fmt.Sprintf("🧺 Layanan dipilih: *%s* (Rp %s/Tabung max 8 Kg) ✨\n\n"+
				"Apakah Kakak sudah menyiapkan sabun & pewangi sendiri, atau ingin sekalian beli di FAN Laundry?\n"+
				"• *Sabun / Pewangi:* Rp 1.000 / Pcs\n\n"+
				"👉 *Pilihan:*\n"+
				"1️⃣ Ketik *1* jika *Bawa Sendiri*\n"+
				"2️⃣ Ketik *2* jika *Beli Sabun & Pewangi* (1 Sabun + 1 Pewangi = Rp 2.000)\n"+
				"_(Atau sebutkan kebutuhan Anda, contoh: 'Beli 2 sabun 2 pewangi')_\n\n"+
				"_(Ketik *BATAL* untuk membatalkan)_", session.ServiceName, formatIDR(matchedService.PricePerKg))
			s.dispatchAsync(phone, reply)
			return &model.BotResponse{Category: bot.CategoryPickup, Reply: reply}, nil
		}

		session.State = bot.StateAwaitingSchedule
		s.sessionStore.Set(phone, session)

		reply := fmt.Sprintf(`Layanan dipilih: *%s* ✨

Kapan kurir kami bisa menjemput cucian Anda?
(Contoh: *Sekarang*, *Nanti jam 2 siang*, atau *Besok pagi*):

_(Ketik *BATAL* untuk membatalkan)_`, session.ServiceName)

		s.dispatchAsync(phone, reply)
		return &model.BotResponse{Category: bot.CategoryPickup, Reply: reply}, nil

	case bot.StateAwaitingDetergent:
		norm := strings.ToLower(text)
		if norm == "1" || strings.Contains(norm, "bawa") || strings.Contains(norm, "sendiri") || strings.Contains(norm, "ada") {
			session.DetergentAddon = "Bawa Sabun Sendiri"
			session.DetergentPrice = 0
		} else if norm == "2" || strings.Contains(norm, "beli") || strings.Contains(norm, "sabun") || strings.Contains(norm, "pewangi") {
			count := 2 // default 1 sabun + 1 pewangi
			words := strings.Fields(norm)
			for _, w := range words {
				if n, err := strconv.Atoi(w); err == nil && n >= 1 && n <= 20 {
					count = n
					break
				}
			}
			session.DetergentAddon = fmt.Sprintf("Beli Sabun/Pewangi (%d pcs)", count)
			session.DetergentPrice = float64(count * 1000)
		} else {
			session.DetergentAddon = "Bawa Sabun Sendiri"
			session.DetergentPrice = 0
		}

		if session.Address == "" {
			session.State = bot.StateAwaitingAddress
			s.sessionStore.Set(phone, session)
			reply := fmt.Sprintf("Pilihan sabun dicatat: *%s* ✨\n\n"+
				"Silakan kirim *Alamat Lengkap Penjemputan* Anda:\n"+
				"(Contoh: _Jl. Mawar No. 12, Kalibaru_)\n\n"+
				"_(Ketik *BATAL* untuk membatalkan)_", session.DetergentAddon)
			s.dispatchAsync(phone, reply)
			return &model.BotResponse{Category: bot.CategoryPickup, Reply: reply}, nil
		}

		session.State = bot.StateAwaitingSchedule
		s.sessionStore.Set(phone, session)
		reply := fmt.Sprintf("Pilihan sabun dicatat: *%s* ✨\n\n"+
			"Kapan kurir kami bisa menjemput cucian Anda?\n"+
			"(Contoh: *Sekarang*, *Nanti jam 2 siang*, atau *Besok pagi*):\n\n"+
			"_(Ketik *BATAL* untuk membatalkan)_", session.DetergentAddon)
		s.dispatchAsync(phone, reply)
		return &model.BotResponse{Category: bot.CategoryPickup, Reply: reply}, nil

	case bot.StateAwaitingSchedule:
		session.Schedule = text
		session.State = bot.StateAwaitingConfirmation
		s.sessionStore.Set(phone, session)

		addonInfo := ""
		if session.DetergentAddon != "" {
			addonInfo = fmt.Sprintf("\n• *Sabun & Pewangi:* %s (Rp %s)", session.DetergentAddon, formatIDR(session.DetergentPrice))
		}

		reply := fmt.Sprintf("🛵 *RINGKASAN REQUEST PENJEMPUTAN* 🛵\n\n"+
			"• *Nama:* Kak %s\n"+
			"• *No. WhatsApp:* %s\n"+
			"• *Alamat Jemput:* %s\n"+
			"• *Paket Layanan:* %s%s\n"+
			"• *Waktu Pickup:* %s\n\n"+
			"Apakah data di atas sudah benar?\n"+
			"👉 Ketik *YA* (atau *OKE*) untuk konfirmasi penjemputan, atau ketik *BATAL* untuk membatalkan.",
			session.CustomerName, phone, session.Address, session.ServiceName, addonInfo, session.Schedule,
		)

		s.dispatchAsync(phone, reply)
		return &model.BotResponse{Category: bot.CategoryPickup, Reply: reply}, nil

	case bot.StateAwaitingConfirmation:
		confirmed := isAffirmative(text)
		if !confirmed {
			reply := "Ketik *YA* atau *OKE* untuk konfirmasi pemesanan, atau ketik *BATAL* jika ingin membatalkan."
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
			notes := fmt.Sprintf("Jadwal: %s | Alamat: %s", session.Schedule, session.Address)
			if session.DetergentAddon != "" {
				notes += fmt.Sprintf(" | Sabun: %s (Rp %s)", session.DetergentAddon, formatIDR(session.DetergentPrice))
			}
			order, err := s.orderSvc.CreatePickupOrder(ctx, model.CreateOrderRequest{
				CustomerID:     customerID,
				ServiceID:      serviceID,
				OutletID:       session.OutletID,
				PickupRequired: true,
				Notes:          notes,
			})
			if err != nil {
				slog.Error("bot: failed to create pickup order", "error", err)
			} else {
				orderNumber = order.OrderNumber
			}
		}

		// 3. Clear session
		s.sessionStore.Delete(phone)

		// 4. Notify Staff / Courier Group for the specific branch
		if s.notifSvc != nil {
			targetGroup := session.TargetGroupJID
			outletTitle := session.OutletName
			go func() {
				_ = s.notifSvc.NotifyPickupRequestToGroup(
					context.Background(),
					targetGroup,
					orderNumber,
					session.CustomerName,
					phone,
					session.Address,
					session.ServiceName,
					session.Schedule,
					outletTitle,
				)
			}()
		}

		branchNote := "Kurir kami"
		if session.OutletName != "" {
			branchNote = fmt.Sprintf("Kurir %s", session.OutletName)
		}

		reply := fmt.Sprintf("🎉 *REQUEST PICKUP BERHASIL DIBUAT!* 🎉\n\n"+
			"• *No. Order:* %s\n"+
			"• *Status:* *MENUNGGU PENJEMPUTAN (WAITING_PICKUP)*\n\n"+
			"%s telah menerima pesanan ini dan akan segera menjemput pakaian kotor Anda sesuai jadwal (*%s*).\n\n"+
			"Terima kasih atas kepercayaannya! 🧺✨", bot.FormatDisplayOrderNumber(orderNumber), branchNote, session.Schedule)

		s.dispatchAsync(phone, reply)
		return &model.BotResponse{Category: bot.CategoryPickup, Reply: reply}, nil

	case bot.StateAwaitingOverloadAction:
		norm := strings.ToLower(text)
		isOption1 := norm == "1" || strings.Contains(norm, "sisih") || strings.Contains(norm, "satu") || strings.Contains(norm, "1 tabung")
		isOption2 := norm == "2" || strings.Contains(norm, "tambah") || strings.Contains(norm, "dua") || strings.Contains(norm, "2 tabung")

		if !isOption1 && !isOption2 {
			reply := "Silakan ketik angka *1* (Sisihkan kelebihan) atau *2* (Jadikan 2 tabung) yaa Kak. 🙏"
			s.dispatchAsync(phone, reply)
			return &model.BotResponse{Category: bot.CategoryPickup, Reply: reply}, nil
		}

		orderNumber := session.OrderNumber
		var order model.OrderTrackingDetail
		if s.orderSvc != nil && orderNumber != "" {
			order, _ = s.orderSvc.GetByOrderNumber(ctx, orderNumber)
		}

		tabungCount := 1
		decisionText := fmt.Sprintf("Pilihan Kakak: *Sisihkan %.1f Kg*. 1 tabung diproses", session.OverloadExcess)
		if isOption2 {
			tabungCount = 2
			decisionText = "Pilihan Kakak: *Jadikan 2 Tabung*. Cucian diproses menjadi 2 tabung"
		}

		baseTabungPrice := 20000.0
		if strings.Contains(strings.ToLower(order.ServiceName), "basah") {
			baseTabungPrice = 10000.0
		}
		totalPrice := (float64(tabungCount) * baseTabungPrice) + session.DetergentPrice

		sla := bot.CalculateCompletionEstimate(order.ServiceName, order.DurationDay, time.Now())
		if s.orderSvc != nil && order.ID != "" {
			newNotes := order.Notes + fmt.Sprintf(" | %s (Total %d Tabung)", decisionText, tabungCount)
			_, _ = s.orderSvc.Update(ctx, order.ID, model.UpdateOrderRequest{
				CustomerID:     order.CustomerID,
				ServiceID:      order.ServiceID,
				TotalWeight:    session.OverloadWeight,
				TotalPrice:     totalPrice,
				PickupRequired: order.PickupRequired,
				Notes:          newNotes,
			})
			_, _ = s.orderSvc.UpdateStatus(ctx, order.ID, model.OrderStatusInProgress)
		}

		s.sessionStore.Delete(phone)

		custAck := fmt.Sprintf("🧺 *NOTA & RINCIAN CUCIAN (FAN Laundry)* 🧺\n\n"+
			"Halo Kak *%s*! %s:\n\n"+
			"• *No. Order:* %s\n"+
			"• *Paket Layanan:* %s\n"+
			"• *Kapasitas:* %d Tabung (Berat total: %.1f Kg)\n"+
			"• *Total Tagihan:* *Rp %s*\n"+
			"• *Status:* *SEDANG DICUCI (IN_PROGRESS)*\n"+
			"• *Estimasi Selesai:* *%s*\n\n"+
			"Pakaian Anda sedang kami proses dengan mesin higienis. Terima kasih! ✨",
			session.CustomerName, decisionText, bot.FormatDisplayOrderNumber(orderNumber), order.ServiceName, tabungCount, session.OverloadWeight, formatIDR(totalPrice), sla)
		s.dispatchAsync(phone, custAck)

		if s.staffPhone != "" {
			groupAck := fmt.Sprintf("✅ *KONFIRMASI PELANGGAN DITERIMA*\n\n"+
				"No. Order %s (Kak *%s*):\n"+
				"• %s\n"+
				"• Total Tagihan: *Rp %s*\n"+
				"• Status: *IN_PROGRESS* (Estimasi: %s)\n"+
				"Nota telah otomatis diperbarui dan dikirimkan ke pelanggan via WhatsApp. 🧾✨",
				bot.FormatDisplayOrderNumber(orderNumber), session.CustomerName, decisionText, formatIDR(totalPrice), sla)
			s.dispatchAsync(s.staffPhone, groupAck)
		}

		return &model.BotResponse{Category: bot.CategoryPickup, Reply: custAck}, nil

	default:
		s.sessionStore.Delete(phone)
		return &model.BotResponse{Category: bot.CategoryUnknown, Reply: "Sesi telah direset."}, nil
	}
}

func isAffirmative(text string) bool {
	norm := strings.TrimSpace(strings.ToUpper(text))
	affirmatives := []string{"OK", "OKE", "YA", "Y", "IYA", "SIAP", "DEAL", "SIP", "GAS", "BENAR", "BOLEH", "COCOK"}
	for _, a := range affirmatives {
		if norm == a {
			return true
		}
	}
	return false
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

func isSelfService(name string) bool {
	lower := strings.ToLower(name)
	return strings.Contains(lower, "self service") || strings.Contains(lower, "cuci basah") || strings.Contains(lower, "cuci kering") || strings.Contains(lower, "tabung")
}

func getServiceCategory(name string) int {
	if isSelfService(name) {
		return 0 // Self Service at the very top!
	}
	lower := strings.ToLower(name)
	if strings.Contains(lower, "cuci setrika") {
		return 1
	}
	if strings.Contains(lower, "setrika") && !strings.Contains(lower, "cuci") {
		return 2
	}
	if strings.Contains(lower, "karpet") || strings.Contains(lower, "permadani") {
		return 3
	}
	if strings.Contains(lower, "bad cover") || strings.Contains(lower, "bed cover") {
		return 4
	}
	if strings.Contains(lower, "boneka") {
		return 5
	}
	return 6
}

func sortServicesByCategory(services []model.ServiceItem) []model.ServiceItem {
	sorted := make([]model.ServiceItem, len(services))
	copy(sorted, services)
	sort.SliceStable(sorted, func(i, j int) bool {
		catI := getServiceCategory(sorted[i].ServiceName)
		catJ := getServiceCategory(sorted[j].ServiceName)
		if catI != catJ {
			return catI < catJ
		}
		return sorted[i].PricePerKg < sorted[j].PricePerKg
	})
	return sorted
}

func formatGroupedServiceMenu(services []model.ServiceItem) string {
	var sb strings.Builder
	currentCat := -1

	for i, s := range services {
		cat := getServiceCategory(s.ServiceName)
		if cat != currentCat {
			if currentCat == 0 {
				sb.WriteString("_(Sabun & Pewangi tersedia: Rp 1.000/pcs)_\n")
			}
			currentCat = cat
			switch cat {
			case 0:
				sb.WriteString("\n🧺 *SELF SERVICE (Max. 8 Kg/Tabung)*\n")
			case 1:
				sb.WriteString("\n🧺 *CUCI SETRIKA (Min. 4 Kg)*\n")
			case 2:
				sb.WriteString("\n👔 *SETRIKA SAJA (Min. 4 Kg)*\n")
			case 3:
				sb.WriteString("\n🛋️ *KARPET PERMADANI (Pengerjaan 7-8 Hari)*\n")
			case 4:
				sb.WriteString("\n🛏️ *BAD COVER*\n")
			case 5:
				sb.WriteString("\n🧸 *BONEKA (Pengerjaan 3 Hari)*\n")
			default:
				sb.WriteString("\n✨ *LAYANAN LAINNYA*\n")
			}
		}

		unit := "/Kg"
		if cat == 0 {
			unit = "/Tabung"
		} else if cat == 3 {
			unit = "/Mtr"
		} else if cat == 4 || cat == 5 {
			unit = "/Pcs"
		}

		durationStr := ""
		lowerName := strings.ToLower(s.ServiceName)
		if cat == 0 {
			durationStr = " (3 jam)"
		} else if strings.Contains(lowerName, "turbo") {
			durationStr = " (6 jam)"
		} else if strings.Contains(lowerName, "super kilat") {
			durationStr = " (12 jam)"
		} else if s.DurationDay > 0 {
			durationStr = fmt.Sprintf(" (%d hari)", s.DurationDay)
		}

		sb.WriteString(fmt.Sprintf("%d. *%s* - Rp %s%s%s\n", i+1, s.ServiceName, formatIDR(s.PricePerKg), unit, durationStr))
	}

	if currentCat == 0 {
		sb.WriteString("_(Sabun & Pewangi tersedia: Rp 1.000/pcs)_\n")
	}

	return sb.String()
}

func (s *botService) handleGroupMessage(ctx context.Context, req model.IncomingMessageRequest) (*model.BotResponse, error) {
	groupJID := strings.TrimSpace(req.GroupJID)
	if groupJID == "" {
		groupJID = strings.TrimSpace(req.From)
	}

	// 1. Extract Order Number from Text or QuotedText
	orderNumber := bot.ExtractOrderNumber(req.Text)
	if orderNumber == "" {
		orderNumber = bot.ExtractOrderNumber(req.QuotedText)
	}

	// If no order number found, ignore to avoid spamming the group
	if orderNumber == "" {
		return &model.BotResponse{Category: "GROUP_IGNORED", Reply: ""}, nil
	}

	if s.orderSvc == nil {
		return &model.BotResponse{Category: "GROUP_IGNORED", Reply: ""}, nil
	}

	order, err := s.orderSvc.GetByOrderNumber(ctx, orderNumber)
	if err != nil || order.OrderNumber == "" {
		slog.Warn("bot: group reply order not found", "order_number", orderNumber, "error", err)
		return &model.BotResponse{Category: "ORDER_NOT_FOUND", Reply: ""}, nil
	}

	staffName := strings.TrimSpace(req.PushName)
	if staffName == "" {
		if req.Participant != "" {
			staffName = "Kurir (" + req.Participant + ")"
		} else {
			staffName = "Kurir FAN Laundry"
		}
	}

	// Fetch customer phone & name
	customerName := order.CustomerName
	var customerPhone string
	if s.customerSvc != nil && order.CustomerID != "" {
		cust, err := s.customerSvc.Get(ctx, order.CustomerID)
		if err == nil {
			customerPhone = cust.Phone
			if cust.FullName != "" {
				customerName = cust.FullName
			}
		}
	}

	displayOrderNum := bot.FormatDisplayOrderNumber(orderNumber)

	// 2. CHECK IF THIS IS A WEIGHING COMMAND (TIMBANG / BERAT)
	if bot.IsWeighCommand(req.Text) {
		weight, found := bot.ParseWeightCommand(req.Text)
		if !found {
			reply := fmt.Sprintf("⚠️ Format timbang tidak valid. Contoh: *TIMBANG %s 5.5*", displayOrderNum)
			s.dispatchAsync(groupJID, reply)
			return &model.BotResponse{Category: "GROUP_STAFF_WEIGH_ERROR", Reply: reply}, nil
		}

		// Check if Self Service
		if isSelfService(order.ServiceName) {
			if weight > 8.5 {
				excess := weight - 8.0
				custPrompt := fmt.Sprintf("⚠️ *KONFIRMASI KAPASITAS MESIN (FAN Laundry)* ⚠️\n\n"+
					"Halo Kak *%s*! Cucian Anda telah tiba di outlet FAN Laundry dan ditimbang: *%.1f Kg*.\n"+
					"Kapasitas normal 1 tabung adalah 8 Kg (toleransi max 8.5 Kg). Terdapat kelebihan *%.1f Kg*.\n\n"+
					"Mohon konfirmasi pilihan Kakak:\n"+
					"1️⃣ *Sisihkan %.1f Kg* (hanya 8 Kg yang dicuci 1 tabung, kelebihan dikembalikan)\n"+
					"2️⃣ *Jadikan 2 Tabung* (bisa tambah pakaian kotor lagi agar pas 2 tabung / 16 Kg)\n\n"+
					"👉 Balas *1* atau *2* yaa Kak.",
					customerName, weight, excess, excess)
				if customerPhone != "" {
					s.dispatchAsync(customerPhone, custPrompt)
					s.sessionStore.Set(customerPhone, &bot.SessionData{
						Phone:          customerPhone,
						CustomerName:   customerName,
						State:          bot.StateAwaitingOverloadAction,
						OrderNumber:    order.OrderNumber,
						OverloadWeight: weight,
						OverloadExcess: excess,
						ServiceID:      order.ServiceID,
						ServiceName:    order.ServiceName,
						LastActive:     time.Now(),
					})
				}

				groupAck := fmt.Sprintf("⚠️ *MENUNGGU KONFIRMASI PELANGGAN*\n\n"+
					"No. Order %s (Kak *%s*) ditimbang *%.1f Kg* (melebihi kapasitas 1 tabung max 8.5 Kg).\n"+
					"Pertanyaan opsi 1 tabung (sisihkan) atau 2 tabung telah otomatis dikirimkan ke WhatsApp pelanggan. 📲",
					displayOrderNum, customerName, weight)
				s.dispatchAsync(groupJID, groupAck)
				return &model.BotResponse{Category: "GROUP_STAFF_WEIGH_OVERLOAD", Reply: groupAck}, nil
			}

			// Weight <= 8.5 kg (Fits 1 Tabung)
			basePrice := 20000.0
			if strings.Contains(strings.ToLower(order.ServiceName), "basah") {
				basePrice = 10000.0
			}
			sla := bot.CalculateCompletionEstimate(order.ServiceName, order.DurationDay, time.Now())
			newNotes := order.Notes + fmt.Sprintf(" | Ditimbang: %.1f Kg (1 Tabung)", weight)
			_, _ = s.orderSvc.Update(ctx, order.ID, model.UpdateOrderRequest{
				CustomerID:     order.CustomerID,
				ServiceID:      order.ServiceID,
				TotalWeight:    weight,
				TotalPrice:     basePrice,
				PickupRequired: order.PickupRequired,
				Notes:          newNotes,
			})
			_, _ = s.orderSvc.UpdateStatus(ctx, order.ID, model.OrderStatusInProgress)

			if customerPhone != "" {
				customerInvoice := fmt.Sprintf("🧺 *NOTA & RINCIAN CUCIAN (FAN Laundry)* 🧺\n\n"+
					"Halo Kak *%s*! Cucian Anda telah tiba di outlet FAN Laundry dan telah selesai ditimbang:\n\n"+
					"• *No. Order:* %s\n"+
					"• *Paket Layanan:* %s\n"+
					"• *Berat Cucian:* %.1f Kg (1 Tabung)\n"+
					"• *Total Tagihan:* *Rp %s*\n"+
					"• *Status:* *SEDANG DICUCI (IN_PROGRESS)*\n"+
					"• *Estimasi Selesai:* *%s*\n\n"+
					"Pakaian Anda sedang kami proses dengan mesin higienis. Terima kasih! ✨",
					customerName, displayOrderNum, order.ServiceName, weight, formatIDR(basePrice), sla)
				s.dispatchAsync(customerPhone, customerInvoice)
			}

			groupAck := fmt.Sprintf("✅ *PENIMBANGAN BERHASIL DICATAT*\n\n"+
				"• *No. Order:* %s (Kak *%s*)\n"+
				"• *Berat Cucian:* %.1f Kg (1 Tabung)\n"+
				"• *Total Tagihan:* Rp %s\n"+
				"• *Status:* *IN_PROGRESS* (Estimasi: %s)\n\n"+
				"Nota dan rincian biaya telah otomatis dikirimkan ke pelanggan via WhatsApp. 🧾✨",
				displayOrderNum, customerName, weight, formatIDR(basePrice), sla)
			s.dispatchAsync(groupJID, groupAck)
			return &model.BotResponse{Category: "GROUP_STAFF_WEIGH_SUCCESS", Reply: groupAck}, nil
		}

		// Kiloan or other services
		ratePerKg := 7000.0
		if s.serviceRepo != nil {
			if svcs, err := s.serviceRepo.ListActiveServices(ctx); err == nil {
				for _, svc := range svcs {
					if svc.ID == order.ServiceID || strings.EqualFold(svc.ServiceName, order.ServiceName) {
						if svc.PricePerKg > 0 {
							ratePerKg = svc.PricePerKg
						}
						break
					}
				}
			}
		}

		chargeWeight := weight
		if chargeWeight < 4.0 && strings.Contains(strings.ToLower(order.ServiceName), "setrika") {
			chargeWeight = 4.0
		}
		totalPrice := chargeWeight * ratePerKg
		sla := bot.CalculateCompletionEstimate(order.ServiceName, order.DurationDay, time.Now())
		newNotes := order.Notes + fmt.Sprintf(" | Ditimbang: %.2f Kg", weight)
		_, _ = s.orderSvc.Update(ctx, order.ID, model.UpdateOrderRequest{
			CustomerID:     order.CustomerID,
			ServiceID:      order.ServiceID,
			TotalWeight:    weight,
			TotalPrice:     totalPrice,
			PickupRequired: order.PickupRequired,
			Notes:          newNotes,
		})
		_, _ = s.orderSvc.UpdateStatus(ctx, order.ID, model.OrderStatusInProgress)

		if customerPhone != "" {
			customerInvoice := fmt.Sprintf("🧺 *NOTA & RINCIAN CUCIAN (FAN Laundry)* 🧺\n\n"+
				"Halo Kak *%s*! Cucian Anda telah tiba di outlet FAN Laundry dan telah selesai ditimbang:\n\n"+
				"• *No. Order:* %s\n"+
				"• *Paket Layanan:* %s\n"+
				"• *Berat Cucian:* %.2f Kg\n"+
				"• *Tarif:* Rp %s/Kg\n"+
				"• *Total Tagihan:* *Rp %s*\n"+
				"• *Status:* *SEDANG DICUCI (IN_PROGRESS)*\n"+
				"• *Estimasi Selesai:* *%s*\n\n"+
				"Pakaian Anda sedang kami proses dengan deterjen higienis dan pewangi premium. Terima kasih! ✨",
				customerName, displayOrderNum, order.ServiceName, weight, formatIDR(ratePerKg), formatIDR(totalPrice), sla)
			s.dispatchAsync(customerPhone, customerInvoice)
		}

		groupAck := fmt.Sprintf("✅ *PENIMBANGAN BERHASIL DICATAT*\n\n"+
			"• *No. Order:* %s (Kak *%s*)\n"+
			"• *Berat Cucian:* %.2f Kg\n"+
			"• *Total Tagihan:* Rp %s\n"+
			"• *Status:* *IN_PROGRESS* (Estimasi: %s)\n\n"+
			"Nota dan rincian biaya telah otomatis dikirimkan ke pelanggan via WhatsApp. 🧾✨",
			displayOrderNum, customerName, weight, formatIDR(totalPrice), sla)
		s.dispatchAsync(groupJID, groupAck)
		return &model.BotResponse{Category: "GROUP_STAFF_WEIGH_SUCCESS", Reply: groupAck}, nil
	}

	// 3. CHECK IF THIS IS A COMPLETION COMMAND (SELESAI / DONE)
	if bot.IsCompleteCommand(req.Text) {
		newNotes := order.Notes + fmt.Sprintf(" | Selesai oleh: %s", staffName)
		_, _ = s.orderSvc.Update(ctx, order.ID, model.UpdateOrderRequest{
			CustomerID:     order.CustomerID,
			ServiceID:      order.ServiceID,
			TotalWeight:    order.TotalWeight,
			TotalPrice:     order.TotalPrice,
			PickupRequired: order.PickupRequired,
			Notes:          newNotes,
		})
		_, _ = s.orderSvc.UpdateStatus(ctx, order.ID, model.OrderStatusCompleted)

		if customerPhone != "" {
			custCompleted := fmt.Sprintf("✨ *CUCIAN ANDA SUDAH SELESAI!* ✨\n\n"+
				"Halo Kak *%s*! Kabar gembira, cucian Anda untuk No. Order %s sudah selesai dicuci, dikeringkan, dan disetrika dengan rapi & wangi. 🧺🌸\n\n"+
				"• *Paket Layanan:* %s\n"+
				"• *Total Biaya:* Rp %s\n"+
				"• *Status Terkini:* *SELESAI & SIAP DIANTAR (COMPLETED)*\n\n"+
				"Kurir kami akan segera mengantarkan cucian kembali ke lokasi Anda. Terima kasih telah mempercayakan cucian Anda di FAN Laundry! 🛵💨",
				customerName, displayOrderNum, order.ServiceName, formatIDR(order.TotalPrice))
			s.dispatchAsync(customerPhone, custCompleted)
		}

		groupAck := fmt.Sprintf("✅ *ORDER COMPLETED*\n\n"+
			"No. Order %s (Kak *%s*) telah selesai diproses oleh *%s*.\n"+
			"Notifikasi cucian siap antar telah dikirimkan ke pelanggan. Kurir dapat segera menjadwalkan pengantaran. 🛵✨",
			displayOrderNum, customerName, staffName)
		s.dispatchAsync(groupJID, groupAck)
		return &model.BotResponse{Category: "GROUP_STAFF_COMPLETE_SUCCESS", Reply: groupAck}, nil
	}

	// 4. DEFAULT: COURIER CLAIMING THE ORDER
	if customerPhone != "" {
		customerMsg := fmt.Sprintf("Halo Kak *%s*! 🛵💨\n\n"+
			"Penjemputan cucian Anda untuk No. Order %s akan segera diambil oleh kurir kami:\n"+
			"👤 *%s*\n\n"+
			"Mohon disiapkan pakaian kotornya yaa. Terima kasih telah mempercayakan cucian Anda di FAN Laundry! 🧺✨",
			customerName, displayOrderNum, staffName)
		s.dispatchAsync(customerPhone, customerMsg)
	}

	groupAck := fmt.Sprintf("✅ *PENJEMPUTAN DICONFIRM*\n\n"+
		"No. Order %s (Kak *%s*) telah diambil oleh *%s*.\n"+
		"Notifikasi penjemputan telah otomatis dikirimkan ke pelanggan via WhatsApp. Semangat bertugas! 🛵💨",
		displayOrderNum, customerName, staffName)
	s.dispatchAsync(groupJID, groupAck)

	courierTag := fmt.Sprintf(" | Kurir: %s", staffName)
	if !strings.Contains(order.Notes, "Kurir:") {
		_, _ = s.orderSvc.Update(ctx, order.ID, model.UpdateOrderRequest{
			CustomerID:     order.CustomerID,
			ServiceID:      order.ServiceID,
			TotalWeight:    order.TotalWeight,
			TotalPrice:     order.TotalPrice,
			PickupRequired: order.PickupRequired,
			Notes:          order.Notes + courierTag,
		})
	}
	_, _ = s.orderSvc.UpdateStatus(ctx, order.ID, model.OrderStatusConfirmed)

	return &model.BotResponse{
		Category: "GROUP_STAFF_CLAIM",
		Reply:    groupAck,
	}, nil
}
