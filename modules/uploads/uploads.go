package uploads

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"learnt.io/core/db"
	"learnt.io/core/rest"
)

type Upload struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	Owner      primitive.ObjectID `json:"owner"`
	Name       string             `json:"name"`
	Mime       string             `json:"mime"`
	CreatedAt  time.Time          `json:"created_at"`
	LastRender time.Time          `json:"last_render"`
}

func Service() *rest.Service {
	s := rest.NewService("uploads", "0.0.0")
	s.Post("/", rest.GuardAuth(upload))
	s.Get("/{id}", rest.GuardAuth(render))
	return s
}

func upload(c *rest.Context) {

	owner, err := primitive.ObjectIDFromHex(c.User())

	if err != nil {
		c.Write(err, http.StatusUnauthorized)
		return
	}

	if err := c.Request().ParseMultipartForm(32 << 20); err != nil {
		c.Write(err, http.StatusBadRequest)
		return
	}

	file, header, err := c.Request().FormFile("upload")

	if err != nil {
		c.Write(err, http.StatusBadRequest)
		return
	}

	defer file.Close()

	data, err := io.ReadAll(file)

	if err != nil {
		c.Write(err, http.StatusInternalServerError)
		return
	}

	upload := Upload{
		ID:        primitive.NewObjectID(),
		Owner:     owner,
		Name:      header.Filename,
		Mime:      http.DetectContentType(data),
		CreatedAt: time.Now(),
	}

	filePath := fmt.Sprintf(
		"%s/%s/%s",
		os.Getenv("UPLOADS_DIR"),
		owner.Hex(),
		upload.ID.Hex(),
	)

	if err := os.MkdirAll(
		filepath.Dir(filePath),
		os.ModePerm,
	); err != nil {
		c.Write(err, http.StatusInternalServerError)
		return
	}

	dst, err := os.Create(filePath)

	if err != nil {
		c.Write(err, http.StatusInternalServerError)
		return
	}

	if _, err := dst.Write(data); err != nil {
		c.Write(err, http.StatusInternalServerError)
		return
	}

	if err := dst.Close(); err != nil {
		c.Write(err, http.StatusInternalServerError)
		return
	}

	_, err = db.C("uploads").InsertOne(
		context.TODO(),
		upload,
	)

	if err != nil {
		c.Write(err, http.StatusInternalServerError)
		return
	}

	c.Write(upload, http.StatusCreated)
}

func render(i *rest.Context) {

	id := i.Param("id")

	if id == "" {
		i.Status(http.StatusNotFound)
		return
	}

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		i.Write(err, http.StatusBadRequest)
		return
	}

	res := db.C("uploads").FindOne(
		context.TODO(),
		bson.M{"_id": oid},
	)

	var upload Upload

	if err := res.Decode(&upload); err != nil {
		i.Write(err, http.StatusInternalServerError)
		return
	}

	path := fmt.Sprintf(
		"%s/%s/%s",
		os.Getenv("UPLOADS_DIR"),
		upload.Owner.Hex(),
		id,
	)

	data, err := os.ReadFile(path)

	if err != nil {
		i.Write(err, http.StatusInternalServerError)
		return
	}

	i.ResponseWriter().Header().Add("Content-Type", upload.Mime)

	if _, err := i.ResponseWriter().Write(data); err != nil {
		i.Write(err, http.StatusInternalServerError)
		return
	}
}
