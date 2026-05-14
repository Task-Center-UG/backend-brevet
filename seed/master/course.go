package master

import (
	"backend-brevet/models"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type courseSeed struct {
	Slug              string
	Title             string
	ShortDescription  string
	Description       string
	LearningOutcomes  string
	Achievements      string
	ImageURL          string
	BatchSlug         string
	BatchTitle        string
	BatchDescription  string
	BatchThumbnailURL string
	CourseType        models.CourseType
	Days              []models.DayType
	StartTime         string
	EndTime           string
	Room              string
	Quota             int
}

// SeedCourses seeds demo courses, batches, meetings, and teacher assignments.
func SeedCourses(db *gorm.DB) error {
	seeds := []courseSeed{
		{
			Slug:              "brevet-pajak-a",
			Title:             "Brevet Pajak A",
			ShortDescription:  "Dasar perpajakan untuk pemula dan calon praktisi pajak.",
			Description:       "Pelajari konsep dasar PPh orang pribadi, KUP, administrasi pajak, dan praktik pengisian dokumen pajak.",
			LearningOutcomes:  "Memahami dasar KUP; Menghitung PPh orang pribadi; Menyiapkan dokumen administrasi pajak.",
			Achievements:      "Sertifikat Brevet A; Studi kasus pajak dasar; Portofolio latihan perhitungan pajak.",
			ImageURL:          "/brevet/brevet-1.jpg",
			BatchSlug:         "brevet-pajak-a-angkatan-dev",
			BatchTitle:        "Brevet Pajak A - Angkatan Dev",
			BatchDescription:  "Kelas dev untuk materi Brevet Pajak A.",
			BatchThumbnailURL: "/brevet/brevet-1.jpg",
			CourseType:        models.CourseTypeOnline,
			Days:              []models.DayType{models.Monday, models.Wednesday},
			StartTime:         "19:00",
			EndTime:           "21:00",
			Room:              "Zoom Dev A",
			Quota:             40,
		},
		{
			Slug:              "brevet-pajak-b",
			Title:             "Brevet Pajak B",
			ShortDescription:  "Pendalaman PPh badan, PPN, pemotongan, dan pelaporan pajak.",
			Description:       "Kursus lanjutan untuk memahami kewajiban pajak perusahaan, PPN, rekonsiliasi fiskal, dan pelaporan masa.",
			LearningOutcomes:  "Menghitung PPh badan; Menyusun rekonsiliasi fiskal; Memahami PPN dan withholding tax.",
			Achievements:      "Sertifikat Brevet B; Studi kasus perusahaan; Simulasi pelaporan pajak.",
			ImageURL:          "/brevet/brevet-2.jpg",
			BatchSlug:         "brevet-pajak-b-angkatan-dev",
			BatchTitle:        "Brevet Pajak B - Angkatan Dev",
			BatchDescription:  "Kelas dev untuk materi Brevet Pajak B.",
			BatchThumbnailURL: "/brevet/brevet-2.jpg",
			CourseType:        models.CourseTypeOnline,
			Days:              []models.DayType{models.Tuesday, models.Thursday},
			StartTime:         "19:00",
			EndTime:           "21:00",
			Room:              "Zoom Dev B",
			Quota:             35,
		},
		{
			Slug:              "pph-orang-pribadi",
			Title:             "PPh Orang Pribadi",
			ShortDescription:  "Perhitungan dan pelaporan PPh untuk wajib pajak orang pribadi.",
			Description:       "Bahas penghasilan, biaya, PTKP, kredit pajak, dan praktik penyusunan SPT Tahunan orang pribadi.",
			LearningOutcomes:  "Mengidentifikasi objek PPh; Menghitung pajak terutang; Menyusun SPT orang pribadi.",
			Achievements:      "Template perhitungan PPh OP; Studi kasus SPT; Sertifikat pelatihan.",
			ImageURL:          "/brevet/brevet-3.jpg",
			BatchSlug:         "pph-orang-pribadi-angkatan-dev",
			BatchTitle:        "PPh Orang Pribadi - Angkatan Dev",
			BatchDescription:  "Kelas dev untuk PPh orang pribadi.",
			BatchThumbnailURL: "/brevet/brevet-3.jpg",
			CourseType:        models.CourseTypeOffline,
			Days:              []models.DayType{models.Saturday},
			StartTime:         "09:00",
			EndTime:           "12:00",
			Room:              "Lab Pajak 1",
			Quota:             30,
		},
		{
			Slug:              "pph-badan",
			Title:             "PPh Badan",
			ShortDescription:  "PPh badan, koreksi fiskal, dan penyusunan SPT Tahunan badan.",
			Description:       "Materi fokus pada rekonsiliasi fiskal, tarif pajak badan, kredit pajak, dan dokumentasi pendukung.",
			LearningOutcomes:  "Membuat rekonsiliasi fiskal; Menghitung PPh badan; Menyiapkan SPT badan.",
			Achievements:      "Studi kasus PPh badan; Workbook koreksi fiskal; Sertifikat pelatihan.",
			ImageURL:          "/brevet/brevet-4.jpg",
			BatchSlug:         "pph-badan-angkatan-dev",
			BatchTitle:        "PPh Badan - Angkatan Dev",
			BatchDescription:  "Kelas dev untuk PPh badan.",
			BatchThumbnailURL: "/brevet/brevet-4.jpg",
			CourseType:        models.CourseTypeOnline,
			Days:              []models.DayType{models.Monday, models.Thursday},
			StartTime:         "18:30",
			EndTime:           "20:30",
			Room:              "Zoom Dev C",
			Quota:             35,
		},
		{
			Slug:              "ppn-dan-ppnbm",
			Title:             "PPN dan PPnBM",
			ShortDescription:  "Konsep, faktur pajak, kredit pajak, dan pelaporan PPN.",
			Description:       "Pelajari objek PPN, faktur pajak, mekanisme pajak masukan dan keluaran, serta PPnBM.",
			LearningOutcomes:  "Memahami objek PPN; Mengelola faktur pajak; Menghitung PPN kurang atau lebih bayar.",
			Achievements:      "Simulasi faktur pajak; Studi kasus PPN; Sertifikat pelatihan.",
			ImageURL:          "/brevet/brevet-5.jpg",
			BatchSlug:         "ppn-dan-ppnbm-angkatan-dev",
			BatchTitle:        "PPN dan PPnBM - Angkatan Dev",
			BatchDescription:  "Kelas dev untuk PPN dan PPnBM.",
			BatchThumbnailURL: "/brevet/brevet-5.jpg",
			CourseType:        models.CourseTypeOnline,
			Days:              []models.DayType{models.Tuesday, models.Friday},
			StartTime:         "18:30",
			EndTime:           "20:30",
			Room:              "Zoom Dev D",
			Quota:             35,
		},
		{
			Slug:              "akuntansi-pajak",
			Title:             "Akuntansi Pajak",
			ShortDescription:  "Jembatan akuntansi komersial menuju laporan fiskal.",
			Description:       "Kursus membahas beda tetap, beda waktu, jurnal pajak, dan penyajian beban pajak.",
			LearningOutcomes:  "Memahami beda fiskal; Menyusun jurnal pajak; Menghubungkan laporan komersial dan fiskal.",
			Achievements:      "Workbook akuntansi pajak; Studi kasus laporan fiskal; Sertifikat pelatihan.",
			ImageURL:          "/brevet/brevet-1.jpg",
			BatchSlug:         "akuntansi-pajak-angkatan-dev",
			BatchTitle:        "Akuntansi Pajak - Angkatan Dev",
			BatchDescription:  "Kelas dev untuk akuntansi pajak.",
			BatchThumbnailURL: "/brevet/brevet-1.jpg",
			CourseType:        models.CourseTypeOffline,
			Days:              []models.DayType{models.Saturday},
			StartTime:         "13:00",
			EndTime:           "16:00",
			Room:              "Lab Pajak 2",
			Quota:             25,
		},
		{
			Slug:              "e-faktur-dan-e-bupot",
			Title:             "e-Faktur dan e-Bupot",
			ShortDescription:  "Praktik penggunaan aplikasi pajak elektronik.",
			Description:       "Latihan membuat faktur, bukti potong, validasi data, dan persiapan pelaporan elektronik.",
			LearningOutcomes:  "Mengoperasikan e-Faktur; Mengelola e-Bupot; Menghindari error umum pelaporan.",
			Achievements:      "Checklist pelaporan elektronik; Simulasi dokumen; Sertifikat workshop.",
			ImageURL:          "/brevet/brevet-2.jpg",
			BatchSlug:         "e-faktur-dan-e-bupot-angkatan-dev",
			BatchTitle:        "e-Faktur dan e-Bupot - Angkatan Dev",
			BatchDescription:  "Kelas dev untuk e-Faktur dan e-Bupot.",
			BatchThumbnailURL: "/brevet/brevet-2.jpg",
			CourseType:        models.CourseTypeOnline,
			Days:              []models.DayType{models.Wednesday},
			StartTime:         "19:00",
			EndTime:           "21:00",
			Room:              "Zoom Dev E",
			Quota:             45,
		},
		{
			Slug:              "tax-planning",
			Title:             "Tax Planning",
			ShortDescription:  "Strategi perencanaan pajak yang patuh dan efisien.",
			Description:       "Bahas tax planning, manajemen risiko pajak, dokumentasi, dan pengambilan keputusan bisnis.",
			LearningOutcomes:  "Menganalisis risiko pajak; Menyusun rencana pajak; Memahami batas kepatuhan.",
			Achievements:      "Casebook tax planning; Simulasi strategi pajak; Sertifikat pelatihan.",
			ImageURL:          "/brevet/brevet-3.jpg",
			BatchSlug:         "tax-planning-angkatan-dev",
			BatchTitle:        "Tax Planning - Angkatan Dev",
			BatchDescription:  "Kelas dev untuk tax planning.",
			BatchThumbnailURL: "/brevet/brevet-3.jpg",
			CourseType:        models.CourseTypeOnline,
			Days:              []models.DayType{models.Thursday},
			StartTime:         "19:00",
			EndTime:           "21:00",
			Room:              "Zoom Dev F",
			Quota:             35,
		},
		{
			Slug:              "pajak-internasional",
			Title:             "Pajak Internasional",
			ShortDescription:  "Dasar transaksi lintas negara dan isu pajak internasional.",
			Description:       "Pelajari P3B, BUT, transfer pricing dasar, dan pajak atas transaksi internasional.",
			LearningOutcomes:  "Memahami P3B; Mengidentifikasi BUT; Mengenali risiko transfer pricing.",
			Achievements:      "Studi kasus internasional; Ringkasan P3B; Sertifikat pelatihan.",
			ImageURL:          "/brevet/brevet-4.jpg",
			BatchSlug:         "pajak-internasional-angkatan-dev",
			BatchTitle:        "Pajak Internasional - Angkatan Dev",
			BatchDescription:  "Kelas dev untuk pajak internasional.",
			BatchThumbnailURL: "/brevet/brevet-4.jpg",
			CourseType:        models.CourseTypeOnline,
			Days:              []models.DayType{models.Friday},
			StartTime:         "19:00",
			EndTime:           "21:00",
			Room:              "Zoom Dev G",
			Quota:             30,
		},
		{
			Slug:              "pemeriksaan-dan-sengketa-pajak",
			Title:             "Pemeriksaan dan Sengketa Pajak",
			ShortDescription:  "Persiapan pemeriksaan, keberatan, dan sengketa pajak.",
			Description:       "Kursus membahas alur pemeriksaan, dokumen pembuktian, keberatan, banding, dan strategi respons.",
			LearningOutcomes:  "Menyiapkan dokumen pemeriksaan; Memahami proses keberatan; Membaca risiko sengketa.",
			Achievements:      "Template respons pemeriksaan; Studi kasus sengketa; Sertifikat pelatihan.",
			ImageURL:          "/brevet/brevet-5.jpg",
			BatchSlug:         "pemeriksaan-dan-sengketa-pajak-angkatan-dev",
			BatchTitle:        "Pemeriksaan dan Sengketa Pajak - Angkatan Dev",
			BatchDescription:  "Kelas dev untuk pemeriksaan dan sengketa pajak.",
			BatchThumbnailURL: "/brevet/brevet-5.jpg",
			CourseType:        models.CourseTypeOffline,
			Days:              []models.DayType{models.Sunday},
			StartTime:         "09:00",
			EndTime:           "12:00",
			Room:              "Lab Pajak 3",
			Quota:             25,
		},
	}

	return db.Transaction(func(tx *gorm.DB) error {
		var teachers []models.User
		if err := tx.Where("role_type = ?", models.RoleTypeGuru).Order("email asc").Find(&teachers).Error; err != nil {
			return fmt.Errorf("failed loading teachers: %w", err)
		}
		if len(teachers) == 0 {
			return errors.New("no teacher users found")
		}

		now := time.Now()
		for i, seed := range seeds {
			course, err := upsertCourse(tx, seed)
			if err != nil {
				return err
			}

			batch, err := upsertBatch(tx, seed, course.ID, now.AddDate(0, i, 0))
			if err != nil {
				return err
			}

			if err := upsertMeetings(tx, batch.ID, seed.Title, now.AddDate(0, i, 7), teachers, i); err != nil {
				return err
			}
		}

		return nil
	})
}

func upsertCourse(tx *gorm.DB, seed courseSeed) (*models.Course, error) {
	course := models.Course{
		Slug:             seed.Slug,
		Title:            seed.Title,
		ShortDescription: seed.ShortDescription,
		Description:      seed.Description,
		LearningOutcomes: seed.LearningOutcomes,
		Achievements:     seed.Achievements,
	}

	var existing models.Course
	err := tx.Where("slug = ?", seed.Slug).First(&existing).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		course.ID = uuid.New()
		if err := tx.Create(&course).Error; err != nil {
			return nil, fmt.Errorf("failed creating course %s: %w", seed.Slug, err)
		}
	case err != nil:
		return nil, fmt.Errorf("failed finding course %s: %w", seed.Slug, err)
	default:
		course.ID = existing.ID
		if err := tx.Model(&existing).Updates(map[string]any{
			"title":             seed.Title,
			"short_description": seed.ShortDescription,
			"description":       seed.Description,
			"learning_outcomes": seed.LearningOutcomes,
			"achievements":      seed.Achievements,
		}).Error; err != nil {
			return nil, fmt.Errorf("failed updating course %s: %w", seed.Slug, err)
		}
	}

	if err := tx.Where("course_id = ?", course.ID).Delete(&models.CourseImage{}).Error; err != nil {
		return nil, fmt.Errorf("failed resetting course images %s: %w", seed.Slug, err)
	}

	image := models.CourseImage{
		ID:       uuid.New(),
		CourseID: course.ID,
		ImageURL: seed.ImageURL,
	}
	if err := tx.Create(&image).Error; err != nil {
		return nil, fmt.Errorf("failed creating course image %s: %w", seed.Slug, err)
	}

	return &course, nil
}

func upsertBatch(tx *gorm.DB, seed courseSeed, courseID uuid.UUID, startAt time.Time) (*models.Batch, error) {
	registrationStart := time.Now().AddDate(0, 0, -7)
	registrationEnd := startAt.AddDate(0, 0, -1)
	endAt := startAt.AddDate(0, 1, 0)

	batch := models.Batch{
		Slug:                seed.BatchSlug,
		CourseID:            courseID,
		Title:               seed.BatchTitle,
		Description:         seed.BatchDescription,
		BatchThumbnail:      seed.BatchThumbnailURL,
		StartAt:             startAt,
		EndAt:               endAt,
		StartTime:           seed.StartTime,
		EndTime:             seed.EndTime,
		RegistrationStartAt: registrationStart,
		RegistrationEndAt:   registrationEnd,
		Room:                seed.Room,
		Quota:               seed.Quota,
		CourseType:          seed.CourseType,
	}

	var existing models.Batch
	err := tx.Where("slug = ?", seed.BatchSlug).First(&existing).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		batch.ID = uuid.New()
		if err := tx.Create(&batch).Error; err != nil {
			return nil, fmt.Errorf("failed creating batch %s: %w", seed.BatchSlug, err)
		}
	case err != nil:
		return nil, fmt.Errorf("failed finding batch %s: %w", seed.BatchSlug, err)
	default:
		batch.ID = existing.ID
		if err := tx.Model(&existing).Updates(map[string]any{
			"course_id":             courseID,
			"title":                 seed.BatchTitle,
			"description":           seed.BatchDescription,
			"batch_thumbnail":       seed.BatchThumbnailURL,
			"start_at":              startAt,
			"end_at":                endAt,
			"start_time":            seed.StartTime,
			"end_time":              seed.EndTime,
			"registration_start_at": registrationStart,
			"registration_end_at":   registrationEnd,
			"room":                  seed.Room,
			"quota":                 seed.Quota,
			"course_type":           seed.CourseType,
		}).Error; err != nil {
			return nil, fmt.Errorf("failed updating batch %s: %w", seed.BatchSlug, err)
		}
	}

	if err := tx.Where("batch_id = ?", batch.ID).Delete(&models.BatchDay{}).Error; err != nil {
		return nil, fmt.Errorf("failed resetting batch days %s: %w", seed.BatchSlug, err)
	}
	if err := tx.Where("batch_id = ?", batch.ID).Delete(&models.BatchGroup{}).Error; err != nil {
		return nil, fmt.Errorf("failed resetting batch groups %s: %w", seed.BatchSlug, err)
	}

	for _, day := range seed.Days {
		if err := tx.Create(&models.BatchDay{
			ID:      uuid.New(),
			BatchID: batch.ID,
			Day:     day,
		}).Error; err != nil {
			return nil, fmt.Errorf("failed creating batch day %s: %w", seed.BatchSlug, err)
		}
	}

	for _, groupType := range []models.GroupType{
		models.MahasiswaGunadarma,
		models.MahasiswaNonGunadarma,
		models.Umum,
	} {
		if err := tx.Create(&models.BatchGroup{
			ID:        uuid.New(),
			BatchID:   batch.ID,
			GroupType: groupType,
		}).Error; err != nil {
			return nil, fmt.Errorf("failed creating batch group %s: %w", seed.BatchSlug, err)
		}
	}

	return &batch, nil
}

func upsertMeetings(tx *gorm.DB, batchID uuid.UUID, courseTitle string, startAt time.Time, teachers []models.User, seedIndex int) error {
	meetingTitles := []string{
		fmt.Sprintf("%s - Orientasi dan Konsep Dasar", courseTitle),
		fmt.Sprintf("%s - Praktik dan Studi Kasus", courseTitle),
	}

	for i, title := range meetingTitles {
		meetingStart := startAt.AddDate(0, 0, i*7)
		meeting := models.Meeting{
			BatchID:     batchID,
			Title:       title,
			Description: fmt.Sprintf("Pertemuan dev %d untuk %s.", i+1, courseTitle),
			Type:        models.BasicMeeting,
			StartAt:     meetingStart,
			EndAt:       meetingStart.Add(2 * time.Hour),
			IsOpen:      true,
		}

		var existing models.Meeting
		err := tx.Where("batch_id = ? AND title = ?", batchID, title).First(&existing).Error
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			meeting.ID = uuid.New()
			if err := tx.Create(&meeting).Error; err != nil {
				return fmt.Errorf("failed creating meeting %s: %w", title, err)
			}
		case err != nil:
			return fmt.Errorf("failed finding meeting %s: %w", title, err)
		default:
			meeting.ID = existing.ID
			if err := tx.Model(&existing).Updates(map[string]any{
				"description": meeting.Description,
				"type":        meeting.Type,
				"start_at":    meeting.StartAt,
				"end_at":      meeting.EndAt,
				"is_open":     meeting.IsOpen,
			}).Error; err != nil {
				return fmt.Errorf("failed updating meeting %s: %w", title, err)
			}
		}

		teacher := teachers[(seedIndex+i)%len(teachers)]
		if err := tx.Where("meeting_id = ?", meeting.ID).Delete(&models.MeetingTeacher{}).Error; err != nil {
			return fmt.Errorf("failed resetting meeting teachers %s: %w", title, err)
		}

		meetingTeacher := models.MeetingTeacher{
			MeetingID: meeting.ID,
			UserID:    teacher.ID,
		}
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&meetingTeacher).Error; err != nil {
			return fmt.Errorf("failed assigning teacher to meeting %s: %w", title, err)
		}
	}

	return nil
}
