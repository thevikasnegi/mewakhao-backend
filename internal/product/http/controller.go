package http

import (
	"encoding/json"
	categoryservice "ecom/internal/category/service"
	"ecom/internal/product/dto"
	"ecom/internal/product/entity"
	productservice "ecom/internal/product/service"
	"ecom/pkg/cloudinary"
	"ecom/pkg/response"
	"ecom/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"
)

type ProductController struct {
	srv       *productservice.ProductService
	catSvc    *categoryservice.CategoryService
	validator validation.Validation
}

func NewProductController(srv *productservice.ProductService, catSvc *categoryservice.CategoryService, validator validation.Validation) *ProductController {
	return &ProductController{
		srv:       srv,
		catSvc:    catSvc,
		validator: validator,
	}
}

func (ctrl *ProductController) List(ctx *gin.Context) {
	categoryID := ctx.Query("category")
	search := ctx.Query("search")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))

	products, total, err := ctrl.srv.List(ctx, categoryID, search, page, limit)
	if err != nil {
		logger.Error("Failed to list products", err)
		response.Error(ctx, http.StatusInternalServerError, err, "Failed to list products")
		return
	}

	var productResList []dto.ProductRes
	for _, p := range products {
		var res dto.ProductRes
		res.ID = p.ID
		res.Name = p.Name
		res.Slug = p.Slug
		res.Description = p.Description
		res.ShortDescription = p.ShortDescription
		res.Images = []string(p.Images)
		res.BasePrice = p.BasePrice
		res.Stock = p.Stock
		res.Rating = p.Rating
		res.ReviewCount = p.ReviewCount
		res.Featured = p.Featured
		res.BestSeller = p.BestSeller
		res.CreatedAt = p.CreatedAt

		if p.Category != nil {
			utils.Copy(&res.Category, p.Category)
		}
		for _, v := range p.Variants {
			var vr dto.VariantRes
			utils.Copy(&vr, &v)
			res.Variants = append(res.Variants, vr)
		}
		if res.Variants == nil {
			res.Variants = []dto.VariantRes{}
		}
		if res.Images == nil {
			res.Images = []string{}
		}
		if p.NutritionalInfo != nil {
			var nr dto.NutritionalInfoRes
			utils.Copy(&nr, p.NutritionalInfo)
			res.NutritionalInfo = &nr
		}
		productResList = append(productResList, res)
	}

	if productResList == nil {
		productResList = []dto.ProductRes{}
	}

	response.JSON(ctx, http.StatusOK, dto.ProductListRes{
		Products: productResList,
		Total:    total,
		Page:     page,
		Limit:    limit,
	})
}

func (ctrl *ProductController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	product, err := ctrl.srv.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to get product", err)
		response.Error(ctx, http.StatusNotFound, err, "Product not found")
		return
	}

	res := ctrl.mapProductToRes(product)
	response.JSON(ctx, http.StatusOK, res)
}

func (ctrl *ProductController) GetBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	product, err := ctrl.srv.GetBySlug(ctx, slug)
	if err != nil {
		logger.Error("Failed to get product by slug", err)
		response.Error(ctx, http.StatusNotFound, err, "Product not found")
		return
	}

	res := ctrl.mapProductToRes(product)
	response.JSON(ctx, http.StatusOK, res)
}

func (ctrl *ProductController) Create(ctx *gin.Context) {
	// Parse simple form fields
	name := ctx.PostForm("name")
	description := ctx.PostForm("description")
	shortDescription := ctx.PostForm("short_description")
	categoryID := ctx.PostForm("category_id")
	basePriceStr := ctx.PostForm("base_price")
	featuredStr := ctx.PostForm("featured")
	bestSellerStr := ctx.PostForm("best_seller")
	variantsJSON := ctx.PostForm("variants")
	nutritionalInfoJSON := ctx.PostForm("nutritional_info")

	basePrice, _ := strconv.ParseFloat(basePriceStr, 64)

	var variants []dto.VariantReq
	if variantsJSON != "" {
		if err := json.Unmarshal([]byte(variantsJSON), &variants); err != nil {
			logger.Error("Failed to parse variants", err)
			response.Error(ctx, http.StatusBadRequest, err, "Invalid variants format")
			return
		}
	}

	var nutritionalInfo *dto.NutritionalInfoReq
	if nutritionalInfoJSON != "" {
		if err := json.Unmarshal([]byte(nutritionalInfoJSON), &nutritionalInfo); err != nil {
			logger.Error("Failed to parse nutritional_info", err)
			response.Error(ctx, http.StatusBadRequest, err, "Invalid nutritional_info format")
			return
		}
	}

	req := &dto.CreateProductReq{
		Name:             name,
		Description:      description,
		ShortDescription: shortDescription,
		CategoryID:       categoryID,
		BasePrice:        basePrice,
		Variants:         variants,
		NutritionalInfo:  nutritionalInfo,
		Featured:         featuredStr == "true",
		BestSeller:       bestSellerStr == "true",
	}

	if err := ctrl.validator.ValidateStruct(req); err != nil {
		logger.Error("Validation failed", err)
		response.Error(ctx, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	// Fetch category to build Cloudinary folder path
	category, err := ctrl.catSvc.GetByID(ctx, categoryID)
	if err != nil {
		logger.Error("Category not found", err)
		response.Error(ctx, http.StatusBadRequest, err, "Category not found")
		return
	}

	// Upload images to Cloudinary under products/<category-slug>/
	form, err := ctx.MultipartForm()
	if err != nil || len(form.File["images"]) == 0 {
		response.Error(ctx, http.StatusBadRequest, err, "At least one image is required")
		return
	}

	folder := "products/" + category.Slug
	var imageURLs []string
	for _, fileHeader := range form.File["images"] {
		file, err := fileHeader.Open()
		if err != nil {
			logger.Error("Failed to open image file", err)
			response.Error(ctx, http.StatusInternalServerError, err, "Failed to process image")
			return
		}
		defer file.Close()

		url, err := cloudinary.UploadImage(ctx, file, folder)
		if err != nil {
			logger.Error("Failed to upload image to Cloudinary", err)
			response.Error(ctx, http.StatusInternalServerError, err, "Failed to upload image")
			return
		}
		imageURLs = append(imageURLs, url)
	}
	req.Images = imageURLs

	product, err := ctrl.srv.Create(ctx, req)
	if err != nil {
		logger.Error("Failed to create product", err)
		response.Error(ctx, http.StatusConflict, err, "Failed to create product")
		return
	}

	res := ctrl.mapProductToRes(product)
	response.JSON(ctx, http.StatusCreated, res)
}

func (ctrl *ProductController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dto.UpdateProductReq
	if err := ctx.ShouldBindJSON(&req); ctx.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(ctx, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	product, err := ctrl.srv.Update(ctx, id, &req)
	if err != nil {
		logger.Error("Failed to update product", err)
		response.Error(ctx, http.StatusInternalServerError, err, "Failed to update product")
		return
	}

	res := ctrl.mapProductToRes(product)
	response.JSON(ctx, http.StatusOK, res)
}

func (ctrl *ProductController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := ctrl.srv.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete product", err)
		response.Error(ctx, http.StatusInternalServerError, err, "Failed to delete product")
		return
	}

	response.JSON(ctx, http.StatusOK, nil)
}

func (ctrl *ProductController) UpdateInventory(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dto.UpdateInventoryReq
	if err := ctx.ShouldBindJSON(&req); ctx.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(ctx, http.StatusBadRequest, err, "Invalid parameters")
		return
	}
	if err := ctrl.validator.ValidateStruct(req); err != nil {
		logger.Error("Validation failed", err)
		response.Error(ctx, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	if err := ctrl.srv.UpdateInventory(ctx, id, &req); err != nil {
		logger.Error("Failed to update inventory", err)
		response.Error(ctx, http.StatusInternalServerError, err, "Failed to update inventory")
		return
	}

	response.JSON(ctx, http.StatusOK, nil)
}

func (ctrl *ProductController) mapProductToRes(p *entity.Product) dto.ProductRes {
	var res dto.ProductRes
	res.ID = p.ID
	res.Name = p.Name
	res.Slug = p.Slug
	res.Description = p.Description
	res.ShortDescription = p.ShortDescription
	res.Images = []string(p.Images)
	res.BasePrice = p.BasePrice
	res.Stock = p.Stock
	res.Rating = p.Rating
	res.ReviewCount = p.ReviewCount
	res.Featured = p.Featured
	res.BestSeller = p.BestSeller
	res.CreatedAt = p.CreatedAt

	if p.Category != nil {
		utils.Copy(&res.Category, p.Category)
	}
	for _, v := range p.Variants {
		var vr dto.VariantRes
		utils.Copy(&vr, &v)
		res.Variants = append(res.Variants, vr)
	}
	if res.Variants == nil {
		res.Variants = []dto.VariantRes{}
	}
	if res.Images == nil {
		res.Images = []string{}
	}
	if p.NutritionalInfo != nil {
		var nr dto.NutritionalInfoRes
		utils.Copy(&nr, p.NutritionalInfo)
		res.NutritionalInfo = &nr
	}
	return res
}
