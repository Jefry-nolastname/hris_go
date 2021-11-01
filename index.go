package main

import (
	"encoding/base64"
	"fmt"
	"hris_go/faceUtils"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Kagami/go-face"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	mode := os.Getenv("GIN_MODE")
	gin.SetMode(mode)
}

func main() {
	route := gin.Default()
	route.POST("/matchface", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			log.Println(err.Error())
			c.JSON(400, gin.H{
				"STATUS": 400,
				"DATA":   "ERROR",
				"ERROR":  "Image not Send!",
			})
			return
		}
		test := form.Value["image"]
		if test == nil {
			log.Println(err.Error())
			c.JSON(400, gin.H{
				"STATUS": 400,
				"DATA":   "ERROR",
				"ERROR":  "Image invalid!",
			})
			return
		}
		decoded, err := base64.StdEncoding.DecodeString(test[0])
		if err != nil {
			log.Println(err.Error())
			c.JSON(400, gin.H{
				"STATUS": 400,
				"DATA":   "ERROR",
				"ERROR":  "Failed processing image!",
			})
			return
		}
		faces, err := faceUtils.FaceDetect(decoded)
		if err != nil {
			log.Println(err.Error())
			c.JSON(400, gin.H{
				"STATUS": 400,
				"DATA":   "ERROR",
				"ERROR":  "Failed finding face!",
			})
			return
		}
		var distance float64
		model_image := form.Value["model_image"]
		if model_image == nil {
			model_desc := form.Value["model_desc"]
			var desc [128]float32
			stringarr := strings.Split(model_desc[0], ",")
			for i, arg := range stringarr {
				if n, err := strconv.ParseFloat(arg, 32); err == nil {
					desc[i] = float32(n)
				}
			}
			distance = face.SquaredEuclideanDistance(faces[0].Descriptor, desc)
		} else {
			model_decoded, err := base64.StdEncoding.DecodeString(model_image[0])
			if err != nil {
				log.Println(err.Error())
				c.JSON(400, gin.H{
					"STATUS": 400,
					"DATA":   "ERROR",
					"ERROR":  "Failed processing image!",
				})
				return
			}
			model_faces, err := faceUtils.FaceDetect(model_decoded)
			if err != nil {
				log.Println(err.Error())
				c.JSON(400, gin.H{
					"STATUS": 400,
					"DATA":   "ERROR",
					"ERROR":  "Error loading model image",
				})
				return
			}
			distance = face.SquaredEuclideanDistance(faces[0].Descriptor, model_faces[0].Descriptor)
		}
		log.Println(distance)
		if distance > 0.2 {
			log.Println("Face didn't matched")
			c.JSON(400, gin.H{
				"STATUS": 400,
				"DATA":   "ERROR",
				// "DISTANCE": distance,
				"ERROR": "Face didn't matched",
			})
			return
		} else {
			log.Println("Face matched!")
			c.JSON(200, gin.H{
				"STATUS": 200,
				"DATA":   strings.Trim(strings.Join(strings.Fields(fmt.Sprint(faces[0].Descriptor)), ","), "[]"),
			})
			return
		}
	})
	route.Run(":" + os.Getenv("PORT")) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
