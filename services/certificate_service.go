package services

import (
	"backend-brevet/config"
	"backend-brevet/helpers"
	"backend-brevet/models"
	"backend-brevet/repository"
	"backend-brevet/utils"
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"

	"github.com/phpdave11/gofpdf/contrib/gofpdi"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

// ICertificateService interface
type ICertificateService interface {
	generatePDF(userName string, batch *models.Batch, certNumber string, listOfMeeting []string, qrPath string) (string, error)
	generateQRCode(certID uuid.UUID) (string, error)
	generateCertificateNumber(ctx context.Context, batchID uuid.UUID) (string, error)
	checkUserAccess(ctx context.Context, user *utils.Claims, batchID uuid.UUID) (bool, error)
	EnsureCertificate(ctx context.Context, batchID uuid.UUID, user *utils.Claims) (*models.Certificate, error)
	GetCertificateByBatchUser(ctx context.Context, batchID uuid.UUID, user *utils.Claims) (*models.Certificate, error)
	GetCertificatesByBatch(ctx context.Context, batchID uuid.UUID, user *utils.Claims) ([]*models.Certificate, error)
	GetCertificateDetail(ctx context.Context, certID uuid.UUID, claims *utils.Claims) (*models.Certificate, error)
	VerifyCertificate(ctx context.Context, certID uuid.UUID) (*models.Certificate, error)
	GetByNumber(ctx context.Context, number string) (*models.Certificate, error)
}

// CertificateService provides methods for managing courses
type CertificateService struct {
	certRepo        repository.ICertificateRepository
	userRepo        repository.IUserRepository
	batchRepo       repository.IBatchRepository
	attendanceRepo  repository.IAttendanceRepository
	meetingRepo     repository.IMeetingRepository
	purchaseService IPurchaseService
	batchService    IBatchService
	fileService     IFileService
}

// NewCertificateService creates a new instance of CertificateService
func NewCertificateService(
	certRepo repository.ICertificateRepository,
	userRepo repository.IUserRepository,
	batchRepo repository.IBatchRepository,
	attendanceRepo repository.IAttendanceRepository,
	meetingRepo repository.IMeetingRepository,
	purchaseService IPurchaseService,
	batchService IBatchService,
	fileService IFileService,
) ICertificateService {
	return &CertificateService{
		certRepo:        certRepo,
		userRepo:        userRepo,
		batchRepo:       batchRepo,
		attendanceRepo:  attendanceRepo,
		meetingRepo:     meetingRepo,
		purchaseService: purchaseService,
		batchService:    batchService,
		fileService:     fileService,
	}
}

// LearningMaterial Misal data materi sudah ada
type LearningMaterial struct {
	No       string
	Material string
}

// generatePDF membuat PDF sertifikat dari template gambar
func (s *CertificateService) generatePDF(userName string, batch *models.Batch, certNumber string, listOfMeeting []string, qrPath string) (string, error) {
	// pastikan folder certificates ada
	if _, err := os.Stat("./certificates"); os.IsNotExist(err) {
		if err := os.Mkdir("./certificates", os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create certificates directory: %w", err)
		}
	}

	templatePath := "./templates/sertifikat.pdf"
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return "", fmt.Errorf("template file not found: %s", templatePath)
	}

	pdf := gofpdf.New("L", "mm", "A4", "")

	// Add fonts with try-catch mechanism
	safeAddFont := func(name, style, path string) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("[WARNING] Failed to add font %s: %v\n", path, r)
			}
		}()
		if _, err := os.Stat(path); err == nil {
			pdf.AddUTF8Font(name, style, path)
		}
	}

	safeAddFont("Lucida", "", "./fonts/lucida.ttf")
	safeAddFont("Cambria", "I", "./fonts/cambriai.ttf")
	safeAddFont("Cambria", "", "./fonts/cambria.ttf")
	safeAddFont("Cambria", "IB", "./fonts/cambriaib.ttf")
	safeAddFont("Cambria", "B", "./fonts/cambriab.ttf")

	// Helper function to safely set font with fallback
	safeSetFont := func(fontFamily, style string, size float64) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("[WARNING] Font setting failed (%s %s): %v, using Arial\n", fontFamily, style, r)
				pdf.SetFont("Arial", style, size)
			}
		}()
		pdf.SetFont(fontFamily, style, size)
	}

	// Safe SplitLines function
	safeSplitLines := func(text string, width float64) [][]byte {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("[WARNING] SplitLines failed for text: %q, width: %f, error: %v\n", text, width, r)
			}
		}()

		if width <= 0 {
			fmt.Printf("[WARNING] Invalid width %f for text: %q\n", width, text)
			return [][]byte{[]byte(text)}
		}

		if len(text) == 0 {
			return [][]byte{[]byte("")}
		}

		// Fallback for very long texts or problematic characters
		if len(text) > 500 {
			fmt.Printf("[WARNING] Text too long (%d chars), truncating: %q\n", len(text), text[:50]+"...")
			text = text[:500] + "..."
		}

		lines := pdf.SplitLines([]byte(text), width)
		if len(lines) == 0 {
			fmt.Printf("[WARNING] SplitLines returned empty for: %q\n", text)
			return [][]byte{[]byte(text)}
		}

		return lines
	}

	importer := gofpdi.NewImporter()

	// Import halaman 1 - halaman utama sertifikat
	tpl1 := importer.ImportPage(pdf, templatePath, 1, "/MediaBox")

	pdf.AddPage()
	importer.UseImportedTemplate(pdf, tpl1, 0, 0, 297, 210)

	// Add QR code
	if _, err := os.Stat(qrPath); err == nil {
		pageW, pageH := pdf.GetPageSize()
		qrSize := 40.0
		margin := 20.0
		x := pageW - qrSize - margin
		y := pageH - qrSize - margin
		pdf.Image(qrPath, x, y, qrSize, 0, false, "", 0, "")
	} else {
		fmt.Printf("[WARNING] QR code file not found: %s\n", qrPath)
	}

	// Tambahkan text di halaman 1
	pdf.SetTextColor(0, 0, 0)
	safeSetFont("Lucida", "", 24)
	pdf.SetXY(100, 95)
	pdf.CellFormat(100, 10, userName, "", 0, "C", false, 0, "")

	safeSetFont("Cambria", "", 12)
	pdf.Ln(11) // jarak vertikal 10mm
	pdf.CellFormat(0, 10, "No. "+certNumber, "", 1, "C", false, 0, "")

	// Pindah baris dulu (biar teks berikutnya di bawahnya)
	pdf.Ln(12.5) // jarak vertikal 10mm
	safeSetFont("Cambria", "I", 12)
	// Tambahin teks di tengah
	pdf.CellFormat(0, 10, "Training Programs of The Applied "+batch.Course.Title, "", 1, "C", false, 0, "")

	// Pindah ke bawah
	pdf.Ln(0.2)

	// Format tanggal batch
	startDate := batch.StartAt.Format("January 02, 2006")
	endDate := batch.EndAt.Format("January 02, 2006")

	period := fmt.Sprintf("%s – %s, Gunadarma University", startDate, endDate)
	pdf.CellFormat(0, 10, period, "", 1, "C", false, 0, "")

	// Ambil tanggal cetak (sekarang)
	printDate := time.Now()

	// Format dengan suffix
	day := printDate.Day()
	suffix := helpers.OrdinalSuffix(day)
	month := printDate.Format("January")
	year := printDate.Format("2006")

	formattedPrintDate := fmt.Sprintf("Jakarta, %s %d%s, %s", month, day, suffix, year)

	// Atur posisi (misalnya di tengah bawah halaman)
	pdf.Ln(0.3) // jarak ke bawah
	pdf.CellFormat(0, 10, formattedPrintDate, "", 1, "C", false, 0, "")

	// Import halaman 2 - halaman tambahan (misal: terms & conditions, info tambahan, dll)
	// Coba import halaman 2
	tpl2 := importer.ImportPage(pdf, templatePath, 2, "/MediaBox")
	if tpl2 == 0 {
		fmt.Println("[WARNING] Failed to import page 2 from template, skipping second page")
	} else {

		// Bangun data materials dari listOfMeeting
		materials := make([]LearningMaterial, 0, len(listOfMeeting))
		for i, m := range listOfMeeting {
			// Sanitize meeting name
			sanitized := strings.TrimSpace(m)
			if len(sanitized) == 0 {
				sanitized = fmt.Sprintf("Meeting %d", i+1)
			}
			materials = append(materials, LearningMaterial{
				No:       fmt.Sprintf("%d", i+1),
				Material: sanitized,
			})
		}

		pdf.AddPage()
		importer.UseImportedTemplate(pdf, tpl2, 0, 0, 297, 210)

		startX := 25.0
		rowHeight := 8.0
		colWidthNo := 15.0
		colWidthMaterial := 150.0
		padding := 1.0
		pageW, pageH := pdf.GetPageSize()

		// ===== Hitung tinggi header dengan safe method =====
		headerNo := "No."
		headerMateri := "Learning Materials"

		safeColWidthNo := colWidthNo - 2*padding
		safeColWidthMaterial := colWidthMaterial - 2*padding

		linesHeaderNo := safeSplitLines(headerNo, safeColWidthNo)
		linesHeaderMateri := safeSplitLines(headerMateri, safeColWidthMaterial)

		maxLinesHeader := len(linesHeaderNo)
		if len(linesHeaderMateri) > maxLinesHeader {
			maxLinesHeader = len(linesHeaderMateri)
		}

		cellHHeader := float64(maxLinesHeader)*rowHeight + 2*padding

		// ===== Hitung total tinggi body dengan safe method =====
		tableBodyH := 0.0
		for _, m := range materials {

			linesMaterial := safeSplitLines(m.Material, safeColWidthMaterial)

			cellH := float64(len(linesMaterial))*rowHeight + 2*padding
			tableBodyH += cellH

		}

		// ===== Total tinggi tabel =====
		tableW := colWidthNo + colWidthMaterial
		tableH := cellHHeader + tableBodyH

		// ===== Hitung startX, startY agar tabel center =====
		startX = (pageW - tableW) / 2
		startY := (pageH - tableH) / 2

		safeSetFont("Cambria", "B", 12)

		// ===== Gambar header =====
		pdf.SetXY(startX+padding, startY+padding)
		pdf.MultiCell(colWidthNo-2*padding, rowHeight, headerNo, "1", "C", false)

		safeSetFont("Cambria", "IB", 12)
		pdf.SetXY(startX+colWidthNo+padding, startY+padding)
		pdf.MultiCell(colWidthMaterial-2*padding, rowHeight, headerMateri, "1", "C", false)

		// ===== Baris data =====
		pdf.SetY(startY + cellHHeader)

		for _, m := range materials {
			yStart := pdf.GetY()

			linesMaterial := safeSplitLines(m.Material, safeColWidthMaterial)

			cellH := float64(len(linesMaterial))*rowHeight + 2*padding

			safeSetFont("Cambria", "", 12)
			// No
			pdf.SetXY(startX+padding, yStart+padding)
			pdf.CellFormat(colWidthNo-2*padding, cellH-2*padding, m.No+".", "1", 0, "C", false, 0, "")

			safeSetFont("Cambria", "I", 12)
			// Materi
			pdf.SetXY(startX+colWidthNo+padding, yStart+padding)
			pdf.MultiCell(colWidthMaterial-2*padding, rowHeight, m.Material, "1", "C", false)

			pdf.SetY(yStart + cellH)
		}
	}

	buf := bytes.NewBuffer(nil)
	if err := pdf.Output(buf); err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	fileName := fmt.Sprintf("certificate_%s_%s.pdf",
		uuid.New().String()[:8],
		time.Now().Format("20060102"))

	// pakai FileService untuk simpan
	publicURL, err := s.fileService.SaveGeneratedFile("certificates", fileName, buf.Bytes())
	if err != nil {
		return "", err
	}

	return publicURL, nil
}

// generateQRCode Generate QR code base64 dari teks
func (s *CertificateService) generateQRCode(certID uuid.UUID) (string, error) {
	// pastikan folder ./qrcodes ada
	err := os.MkdirAll("./qrcodes", os.ModePerm)
	if err != nil {
		return "", err
	}

	filePath := fmt.Sprintf("./qrcodes/%s.png", uuid.New().String())

	baseURL := config.GetEnv("FRONTEND_URL", "http://localhost")
	url := fmt.Sprintf("%s/sertifikat/%s", baseURL, certID.String())

	err = qrcode.WriteFile(url, qrcode.Medium, 256, filePath)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func (s *CertificateService) generateCertificateNumber(ctx context.Context, batchID uuid.UUID) (string, error) {
	// prefix bisa statis, misalnya kode lembaga
	prefix := "20100112"

	// ambil nomor urut terakhir untuk batch ini
	lastSeq, err := s.certRepo.GetLastSequenceByBatch(ctx, batchID)
	if err != nil {
		return "", err
	}

	// increment
	nextSeq := lastSeq + 1

	// misal batch keberapa, bisa ambil dari batch table
	batch, err := s.batchRepo.FindByID(ctx, batchID)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s-%s %04d", prefix, batch.ID.String()[:8], nextSeq), nil
}

func (s *CertificateService) checkUserAccess(ctx context.Context, user *utils.Claims, batchID uuid.UUID) (bool, error) {
	// Cari batch info dari meetingID
	batch, err := s.batchRepo.FindByID(ctx, batchID) // balikin batchSlug & batchID
	if err != nil {
		return false, err
	}

	// Kalau role teacher, cek apakah dia mengajar batch ini
	if user.Role == string(models.RoleTypeGuru) {
		return s.meetingRepo.IsBatchOwnedByUser(ctx, user.UserID, batch.Slug)
	}

	// Kalau student, cek pembayaran
	if user.Role == string(models.RoleTypeSiswa) {
		return s.purchaseService.HasPaid(ctx, user.UserID, batch.ID)
	}

	// Role lain tidak diizinkan
	return false, nil
}

// EnsureCertificate memastikan sertifikat dibuat untuk user
func (s *CertificateService) EnsureCertificate(ctx context.Context, batchID uuid.UUID, user *utils.Claims) (*models.Certificate, error) {

	allowed, err := s.checkUserAccess(ctx, user, batchID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, fmt.Errorf("forbidden")
	}

	// 1. cek progress assignment & quiz
	progress, err := s.batchService.CalculateProgress(ctx, batchID, user.UserID)
	if err != nil {
		return nil, err
	}
	if progress < 100 {
		return nil, fmt.Errorf("progress belum 100%%")
	}

	// 2. cek attendance
	totalMeetings, err := s.batchRepo.CountMeetings(ctx, batchID)
	if err != nil {
		return nil, err
	}
	attendedMeetings, err := s.attendanceRepo.CountByBatchUser(ctx, batchID, user.UserID)
	if err != nil {
		return nil, err
	}
	if totalMeetings > 0 && attendedMeetings < totalMeetings {
		return nil, fmt.Errorf("murid belum hadir di semua pertemuan (%d/%d)", attendedMeetings, totalMeetings)
	}

	// ✅ 3. ambil data batch untuk cek end_at
	batch, err := s.batchRepo.GetBatchWithCourse(ctx, batchID)
	if err != nil {
		return nil, err
	}
	if batch == nil {
		return nil, fmt.Errorf("batch tidak ditemukan")
	}

	// ✅ 4. validasi periode cetak sertifikat (H+1 setelah end_at)
	now := time.Now()
	allowedPrintDate := batch.EndAt.Add(24 * time.Hour) // H+1
	if now.Before(allowedPrintDate) {
		return nil, fmt.Errorf("periode cetak sertifikat dimulai H+1 setelah tanggal selesai batch (%s)", batch.EndAt.Format("2006-01-02"))
	}

	// cek sertifikat existing
	_, err = s.certRepo.GetByBatchUser(ctx, batchID, user.UserID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == nil {
		return nil, fmt.Errorf("sertifikat sudah ada untuk user %s di batch %s", user.UserID, batchID)
	}

	certNumber, err := s.generateCertificateNumber(ctx, batchID)
	if err != nil {
		return nil, err
	}

	cert := &models.Certificate{
		BatchID: batchID,
		Number:  certNumber,
		UserID:  user.UserID,
	}
	err = s.certRepo.Create(ctx, cert)
	if err != nil {
		return nil, err
	}

	qrPath, err := s.generateQRCode(cert.ID)
	if err != nil {
		return nil, err
	}
	defer os.Remove(qrPath)

	listOfMeeting, err := s.meetingRepo.GetMeetingNamesByBatchID(ctx, batch.ID)
	if err != nil {
		return nil, err
	}

	pdfURL, err := s.generatePDF(user.Name, batch, certNumber, listOfMeeting, qrPath)
	if err != nil {
		return nil, err
	}

	cert.URL = pdfURL
	cert.QRCode = fmt.Sprintf("/certificates/%s", cert.ID)

	err = s.certRepo.Update(ctx, cert)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

// GetCertificateByBatchUser get certificate by batch id & user id
func (s *CertificateService) GetCertificateByBatchUser(ctx context.Context, batchID uuid.UUID, user *utils.Claims) (*models.Certificate, error) {
	allowed, err := s.checkUserAccess(ctx, user, batchID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, fmt.Errorf("forbidden")
	}

	cert, err := s.certRepo.GetByBatchUser(ctx, batchID, user.UserID)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

// GetCertificatesByBatch for get
func (s *CertificateService) GetCertificatesByBatch(ctx context.Context, batchID uuid.UUID, user *utils.Claims) ([]*models.Certificate, error) {
	allowed, err := s.checkUserAccess(ctx, user, batchID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, fmt.Errorf("forbidden")
	}

	return s.certRepo.GetByBatch(ctx, batchID)
}

// GetCertificateDetail detail certificate
func (s *CertificateService) GetCertificateDetail(ctx context.Context, certID uuid.UUID, claims *utils.Claims) (*models.Certificate, error) {
	// Ambil certificate dulu berdasarkan certID
	cert, err := s.certRepo.GetByID(ctx, certID)
	if err != nil {
		return nil, err
	}

	// baru ambil batchID dari hasil query
	allowed, err := s.checkUserAccess(ctx, claims, cert.BatchID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, fmt.Errorf("forbidden")
	}

	return cert, nil
}

// VerifyCertificate verifies certificate authenticity by ID (no auth required)
func (s *CertificateService) VerifyCertificate(ctx context.Context, certID uuid.UUID) (*models.Certificate, error) {
	return s.certRepo.GetByID(ctx, certID)
}

// GetByNumber get certificate by number
func (s *CertificateService) GetByNumber(ctx context.Context, number string) (*models.Certificate, error) {
	return s.certRepo.GetByNumber(ctx, number)
}
