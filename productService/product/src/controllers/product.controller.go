package controllers

import (
	"context"
	"mime/multipart"
	"net/http"
	"supernova/productService/product/src/db"
	"supernova/productService/product/src/dto"
	"supernova/productService/product/src/models"
	"supernova/productService/product/src/services"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const maxImages = 5
const fileFieldName = "images"

// CreateProduct handles product creation and image upload to Cloudinary
func CreateProduct(c *gin.Context) {
    var productDTO dto.ProductDTO
    
    if err := c.ShouldBind(&productDTO); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }


    form, err := c.MultipartForm()
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read multipart form"})
        return
    }

    
    files := form.File["images"]
    if len(files) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No images uploaded"})
        return
    }
    productDTO.Images = files

   

    var images []models.Image

   var wg sync.WaitGroup // WaitGroup for synchronization
    imageResults := make(chan struct {
        image models.Image
        err   error
    }, len(productDTO.Images)) // Channel to collect results

    for _, fileHeader := range productDTO.Images {
        // Validation (keep this sequential for immediate response)
        contentType := fileHeader.Header.Get("Content-Type")
        if contentType != "image/jpeg" && contentType != "image/png" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Only JPEG and PNG images are allowed"})
            return
        }

        // Open uploaded file (must be done before the Goroutine)
        file, err := fileHeader.Open()
        if err != nil {
             // Handle error
            return
        }
        
        wg.Add(1) // Increment the WaitGroup counter

        // Start a Goroutine for concurrent upload
        go func(file multipart.File) {
            defer wg.Done()   // Decrement the counter when the Goroutine finishes
            defer file.Close() // Close the file stream

            // Upload to Cloudinary
            res, err := services.UploadImage(file)
            
            // Send result (or error) to the channel
            if err != nil {
                imageResults <- struct {
                    image models.Image
                    err   error
                }{image: models.Image{}, err: err}
                return
            }
            
            // Append to images slice
            image := models.Image{
                URL:       res.SecureURL,
                Thumbnail: res.SecureURL,
                ID:        res.PublicID,
            }
            imageResults <- struct {
                image models.Image
                err   error
            }{image: image, err: nil}

        }(file) // Pass the file handle to the Goroutine
    }

    // Wait for all uploads to complete
    wg.Wait()
    close(imageResults) // Close the channel after all Goroutines are done

    // Collect results from the channel
    for result := range imageResults {
        if result.err != nil {
            // Handle the first error and abort
            c.JSON(http.StatusInternalServerError, gin.H{"error": result.err.Error()})
            return
        }
        images = append(images, result.image)
    }
    sellerID , exists := c.Get("UserID")
    if !exists {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get seller ID from context"})
        return
    }
    sellerIDStr, ok := sellerID.(string)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid seller ID format"})
        return
    }

    // Create Product object
    product := models.Product{
        Title:       productDTO.Title,
        Description: productDTO.Description,
        Price: models.Price{
            Amount:   productDTO.Price.Amount,
            Currency: productDTO.Price.Currency,
        },
        Images: images,
        Stock:  productDTO.Stock,
        SellerID: sellerIDStr,
    }

     collection := db.GetProductCollection()

	 _,err = collection.InsertOne(c, product)
	 if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	 }

    c.JSON(http.StatusOK, gin.H{
        "message": "Product created successfully",
        "product": product,
    })
    return
}


func GetProducts(c *gin.Context) {
    // Query parameters
    q := c.Query("q")
    minPriceStr := c.DefaultQuery("minprice", "")
    maxPriceStr := c.DefaultQuery("maxprice", "")
    skipStr := c.DefaultQuery("skip", "0")
    limitStr := c.DefaultQuery("limit", "20")

    // Convert pagination values to integers
    skip, _ := strconv.Atoi(skipStr)
    limit, _ := strconv.Atoi(limitStr)

    // Build MongoDB filter
    filter := bson.M{}

    // Text search (title OR description)
    if q != "" {
        filter["$or"] = []bson.M{
            {"title": bson.M{"$regex": q, "$options": "i"}},
            {"description": bson.M{"$regex": q, "$options": "i"}},
        }
    }

    // Price filtering
    priceFilter := bson.M{}
    if minPriceStr != "" {
        if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
            priceFilter["$gte"] = minPrice
        }
    }
    if maxPriceStr != "" {
        if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
            priceFilter["$lte"] = maxPrice
        }
    }
    if len(priceFilter) > 0 {
        filter["price.amount"] = priceFilter
    }

    // Context for MongoDB
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Query options (pagination)
    findOptions := options.Find().
        SetSkip(int64(skip)).
        SetLimit(int64(limit))

    cursor, err := db.GetProductCollection().Find(ctx, filter, findOptions)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer cursor.Close(ctx)

    // Decode all products
    var products []models.Product
    if err := cursor.All(ctx, &products); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Send response
    c.JSON(http.StatusOK, gin.H{
        "count":    len(products),
        "skip":     skip,
        "limit":    limit,
        "products": products,
    })
    return
}

func GetProductByID(c *gin.Context) {
    id := c.Param("id")

    
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
        return
    }

    // Context for MongoDB
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var product models.Product
    err = db.GetProductCollection().FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }
    c.JSON(http.StatusOK, product)
}

func UpdateProduct(c *gin.Context) {
    id := c.Param("id")

    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
        return
    }
    var productDTO dto.ProductDTO
    
    if err := c.ShouldBind(&productDTO); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    form, err := c.MultipartForm()
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read multipart form"})
        return
    }
    files := form.File["images"]
    if len(files) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No images uploaded"})
        return
    }
    productDTO.Images = files

    var images []models.Image
    var wg sync.WaitGroup // WaitGroup for synchronization  
    imageResults := make(chan struct {
        image models.Image
        err   error
    }, len(productDTO.Images)) // Channel to collect results
    for _, fileHeader := range productDTO.Images {
        // Validation (keep this sequential for immediate response)
        contentType := fileHeader.Header.Get("Content-Type")
        if contentType != "image/jpeg" && contentType != "image/png" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Only JPEG and PNG images are allowed"})
            return
        }
        // Open uploaded file (must be done before the Goroutine)
        file, err := fileHeader.Open()
        if err != nil {
             // Handle error
            return
        }
        wg.Add(1) // Increment the WaitGroup counter
        // Start a Goroutine for concurrent upload
        go func(file multipart.File) {  
            defer wg.Done()   // Decrement the counter when the Goroutine finishes
            defer file.Close() // Close the file stream
            // Upload to Cloudinary
            res, err := services.UploadImage(file)
            // Send result (or error) to the channel
            if err != nil {
                imageResults <- struct {
                    image models.Image
                    err   error
                }{image: models.Image{}, err: err}
                return
            }
            // Append to images slice
            image := models.Image{
                URL:       res.SecureURL,
                Thumbnail: res.SecureURL,
                ID:        res.PublicID,
            }
            imageResults <- struct {    
                image models.Image
                err   error
            }{image: image, err: nil}
        }(file) // Pass the file handle to the Goroutine
    }
    // Wait for all uploads to complete
    wg.Wait()
    close(imageResults) // Close the channel after all Goroutines are done
    // Collect results from the channel
    for result := range imageResults {
        if result.err != nil {
            // Handle the first error and abort
            c.JSON(http.StatusInternalServerError, gin.H{"error": result.err.Error()})
            return
        }
        images = append(images, result.image)
    }

    // Create Product object
    product := models.Product{
        Title:       productDTO.Title,
        Description: productDTO.Description,
        Price: models.Price{
            Amount:   productDTO.Price.Amount,
            Currency: productDTO.Price.Currency,
        },
        Images: images,
    }   
    collection := db.GetProductCollection()
    update := bson.M{
        "$set": bson.M{
            "title":       product.Title,
            "description": product.Description,
            "price":       product.Price,
            "images":      product.Images,
        },
    }
    result, err := collection.UpdateByID(c, objectID, update)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
        return
    }
    if result.MatchedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "message": "Product updated successfully",
        "product": product,
    })
    return
}