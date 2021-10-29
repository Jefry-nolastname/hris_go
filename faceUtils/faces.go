package faceUtils

import (
	"errors"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/Kagami/go-face"
)

func FaceDetect(byte []byte) ([]face.Face, error) {
	var face_detected []face.Face
	path, err := os.Getwd()
	if err != nil {
		return face_detected, errors.New("Error loading face model!")
	}
	modelsDir := filepath.Join(path, "models")
	rec, err := face.NewRecognizer(modelsDir)
	if err != nil {
		return face_detected, errors.New("Can't init face recognizer!")
	}
	// Free the resources when you're finished.
	defer rec.Close()

	faces, err := rec.Recognize(byte)
	if err != nil {
		log.Println(err)
		return face_detected, errors.New("Failed finding face in picture!")
	}
	if len(faces) > 1 {
		return face_detected, errors.New("Please take the photo alone!")
	} else if len(faces) < 1 {

		return face_detected, errors.New("No face detected!")
	}
	face_detected = faces
	return face_detected, nil
}
func FileDetect(imageData []*multipart.FileHeader) (face.Face, error) {
	var face_detected face.Face

	file, err := imageData[0].Open()
	if err != nil {
		return face_detected, errors.New("Failed opening images!")
	}
	image, err := ioutil.ReadAll(file)
	if err != nil {
		return face_detected, errors.New("Failed reading images!")
	}
	faces, err := FaceDetect(image)
	if err != nil {
		log.Println(err)
		return face_detected, errors.New(err.Error())
	}
	face_detected = faces[0]
	return face_detected, nil
}
