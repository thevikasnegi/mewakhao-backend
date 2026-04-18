package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	categoryEntity "ecom/internal/category/entity"
	productEntity "ecom/internal/product/entity"
	"ecom/pkg/cloudinary"
	"ecom/pkg/config"
	"ecom/pkg/dbs"
)

func main() {
	cfg := config.LoadConfig()

	db, err := dbs.NewDatabase(cfg.DatabaseURI)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	gdb := db.GetDB()

	err = db.AutoMigrate(
		&categoryEntity.Category{},
		&productEntity.Product{},
		&productEntity.ProductVariant{},
		&productEntity.NutritionalInfo{},
	)
	if err != nil {
		log.Fatal("Failed to migrate:", err)
	}

	assetsDir := "assets"
	ctx := context.Background()

	categoryDirs, err := os.ReadDir(assetsDir)
	if err != nil {
		log.Fatal("Cannot read assets directory:", err)
	}

	for _, catEntry := range categoryDirs {
		if !catEntry.IsDir() {
			continue
		}

		catName := catEntry.Name()
		catSlug := slugify(catName)
		catDir := filepath.Join(assetsDir, catName)

		// Upsert category
		var cat categoryEntity.Category
		gdb.Where("slug = ?", catSlug).First(&cat)

		isNewCat := cat.ID == ""
		if isNewCat {
			cat.Name = catName
			cat.Slug = catSlug
			cat.Description = categoryDescription(catName)
		}

		// Upload category image from first product's first image (if not set)
		if cat.Image == "" {
			cat.Image = uploadFirstImage(ctx, catDir, "mewa-khao/categories")
		}

		if isNewCat {
			if err := gdb.Create(&cat).Error; err != nil {
				log.Printf("✗ Failed to create category %s: %v", catName, err)
				continue
			}
			log.Printf("✓ Category: %s", catName)
		} else if cat.Image != "" {
			gdb.Model(&cat).Update("image", cat.Image)
			log.Printf("~ Category exists, image updated: %s", catName)
		} else {
			log.Printf("~ Category exists: %s", catName)
		}

		// Seed products inside this category
		productDirs, err := os.ReadDir(catDir)
		if err != nil {
			log.Printf("  Cannot read category dir %s: %v", catDir, err)
			continue
		}

		for _, prodEntry := range productDirs {
			if !prodEntry.IsDir() {
				continue
			}

			prodName := prodEntry.Name()
			prodSlug := slugify(prodName)
			prodDir := filepath.Join(catDir, prodName)

			// Skip if product already exists
			var existing productEntity.Product
			if gdb.Where("slug = ?", prodSlug).First(&existing).Error == nil {
				log.Printf("  ~ Product exists: %s", prodName)
				continue
			}

			// Upload all product images
			imageURLs := uploadAllImages(ctx, prodDir, "mewa-khao/products/"+prodSlug)
			if len(imageURLs) == 0 {
				log.Printf("  ⚠ No images for %s, skipping", prodName)
				continue
			}

			basePrice := categoryBasePrice(catName)
			nutri := categoryNutritionalInfo(catName)
			nutri.ProductID = "" // will be set by GORM after product insert

			prod := productEntity.Product{
				Name:             prodName,
				Slug:             prodSlug,
				Description:      fmt.Sprintf("Premium %s — %s. Sourced from the finest farms, carefully processed to deliver exceptional taste and quality.", catName, prodName),
				ShortDescription: fmt.Sprintf("High quality %s", prodName),
				CategoryID:       cat.ID,
				Images:           productEntity.StringArray(imageURLs),
				BasePrice:        basePrice,
				Stock:            200,
				Featured:         false,
				BestSeller:       false,
				Variants:         []productEntity.ProductVariant{},
				NutritionalInfo:  &nutri,
			}

			if err := gdb.Create(&prod).Error; err != nil {
				log.Printf("  ✗ Failed to create product %s: %v", prodName, err)
				continue
			}
			log.Printf("  ✓ Product: %s (%d images)", prodName, len(imageURLs))
		}
	}

	log.Println("\n✅ Seeding complete!")
}

// uploadFirstImage finds the first image file in any immediate subfolder and uploads it.
func uploadFirstImage(ctx context.Context, dir string, folder string) string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		subDir := filepath.Join(dir, entry.Name())
		files, err := os.ReadDir(subDir)
		if err != nil {
			continue
		}
		for _, f := range files {
			if f.IsDir() || !isImageFile(f.Name()) {
				continue
			}
			url := uploadFile(ctx, filepath.Join(subDir, f.Name()), folder)
			if url != "" {
				return url
			}
		}
	}
	return ""
}

// uploadAllImages uploads every image in the given directory to Cloudinary.
func uploadAllImages(ctx context.Context, dir string, folder string) []string {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var urls []string
	for _, f := range files {
		if f.IsDir() || !isImageFile(f.Name()) {
			continue
		}
		url := uploadFile(ctx, filepath.Join(dir, f.Name()), folder)
		if url != "" {
			urls = append(urls, url)
			log.Printf("    ↑ %s", f.Name())
		} else {
			log.Printf("    ✗ failed: %s", f.Name())
		}
	}
	return urls
}

func uploadFile(ctx context.Context, path string, folder string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()
	url, err := cloudinary.UploadImage(ctx, f, folder)
	if err != nil {
		return ""
	}
	return url
}

func isImageFile(name string) bool {
	lower := strings.ToLower(name)
	return strings.HasSuffix(lower, ".jpg") ||
		strings.HasSuffix(lower, ".jpeg") ||
		strings.HasSuffix(lower, ".png") ||
		strings.HasSuffix(lower, ".webp")
}

func slugify(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == ' ' || r == '_':
			b.WriteRune('-')
		}
	}
	return strings.Trim(b.String(), "-")
}

func round2(v float64) float64 {
	return float64(int(v*100+0.5)) / 100
}

func categoryDescription(name string) string {
	m := map[string]string{
		"Almonds": "Premium quality almonds sourced from the finest orchards",
		"Cashew":  "Creamy and delicious cashews with exceptional taste",
		"Makhana": "Light and nutritious lotus seeds, a healthy snacking choice",
		"Pista":   "Naturally opened green pistachios bursting with flavour",
		"Walnut":  "Brain-boosting walnuts rich in omega-3 fatty acids",
	}
	if d, ok := m[name]; ok {
		return d
	}
	return fmt.Sprintf("Premium quality %s", strings.ToLower(name))
}

func categoryBasePrice(name string) float64 {
	m := map[string]float64{
		"Almonds": 14.99,
		"Cashew":  16.99,
		"Makhana": 9.99,
		"Pista":   19.99,
		"Walnut":  13.99,
	}
	if p, ok := m[name]; ok {
		return p
	}
	return 12.99
}

func categoryNutritionalInfo(name string) productEntity.NutritionalInfo {
	m := map[string]productEntity.NutritionalInfo{
		"Almonds": {Calories: "579 kcal", Protein: "21g", Fat: "50g", Carbs: "22g", Fiber: "12g"},
		"Cashew":  {Calories: "553 kcal", Protein: "18g", Fat: "44g", Carbs: "30g", Fiber: "3g"},
		"Makhana": {Calories: "347 kcal", Protein: "9g", Fat: "0.1g", Carbs: "77g", Fiber: "14g"},
		"Pista":   {Calories: "562 kcal", Protein: "20g", Fat: "45g", Carbs: "28g", Fiber: "10g"},
		"Walnut":  {Calories: "654 kcal", Protein: "15g", Fat: "65g", Carbs: "14g", Fiber: "7g"},
	}
	if n, ok := m[name]; ok {
		return n
	}
	return productEntity.NutritionalInfo{}
}
