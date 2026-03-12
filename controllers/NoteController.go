package controllers

import (
	"context"
	"go-app/db"
	"go-app/cache"
	"go-app/models"
	u "go-app/utils"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var NoteCreate = func(w http.ResponseWriter, r *http.Request) {
	Note := &models.Note{}
	err := json.NewDecoder(r.Body).Decode(Note)

	if err != nil {
		u.HandleBadRequest(w, err)
		return
	}

	db := db.GetDB()
	err = db.Create(Note).Error

	if err != nil {
		u.HandleBadRequest(w, err)
	} else {
		res, _ := json.Marshal(Note)
		u.RespondJSON(w, res)
	}
}

var NoteRetrieve = func(w http.ResponseWriter, r *http.Request) {
	Note := &models.Note{}

	params := mux.Vars(r)
	id := params["id"]

	redisClient := cache.GetRedis()
	cacheKey := cache.NoteKey(id)
	if redisClient != nil {
		ctx, cancel := context.WithTimeout(r.Context(), 500*time.Millisecond)
		defer cancel()

		if cached, err := redisClient.Get(ctx, cacheKey).Bytes(); err == nil {
			log.Infof("[REDIS] GET %s: %s", cacheKey, string(cached))
			u.RespondJSON(w, cached)
			return
		} else if err != redis.Nil {
			log.Warnf("[REDIS] GET %s failed: %v", cacheKey, err)
		}
	}

	db := db.GetDB()
	err := db.First(&Note, id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			u.HandleNotFound(w)
		} else {
			u.HandleBadRequest(w, err)
		}
		return
	}

	res, err := json.Marshal(Note)
	if err != nil {
		u.HandleBadRequest(w, err)
	} else if Note.ID == 0 {
		u.HandleNotFound(w)
	} else {
		if redisClient != nil {
			ctx, cancel := context.WithTimeout(r.Context(), 500*time.Millisecond)
			defer cancel()

			if err := redisClient.Set(ctx, cacheKey, res, cache.TTL()).Err(); err != nil {
				log.Warnf("[REDIS] SET %s failed: %v", cacheKey, err)
			}
		}
		u.RespondJSON(w, res)
	}
}

var NoteUpdate = func(w http.ResponseWriter, r *http.Request) {
	Note := &models.Note{}

	params := mux.Vars(r)
	id := params["id"]

	db := db.GetDB()
	err := db.First(&Note, id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			u.HandleNotFound(w)
		} else {
			u.HandleBadRequest(w, err)
		}
		return
	}

	newNote := &models.Note{}
	err = json.NewDecoder(r.Body).Decode(newNote)

	if err != nil {
		u.HandleBadRequest(w, err)
		return
	}

	err = db.Model(&Note).Updates(newNote).Error

	if err != nil {
		u.HandleBadRequest(w, err)
	} else {
		if redisClient := cache.GetRedis(); redisClient != nil {
			cacheKey := cache.NoteKey(id)
			ctx, cancel := context.WithTimeout(r.Context(), 500*time.Millisecond)
			defer cancel()

			if err := redisClient.Del(ctx, cacheKey).Err(); err != nil {
				log.Warnf("[REDIS] DEL %s failed: %v", cacheKey, err)
			}
		}
		u.Respond(w, u.Message(true, "OK"))
	}
}

var NoteDelete = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	db := db.GetDB()
	err := db.Delete(&models.Note{}, id).Error

	if err != nil {
		u.HandleBadRequest(w, err)
	} else {
		if redisClient := cache.GetRedis(); redisClient != nil {
			cacheKey := cache.NoteKey(id)
			ctx, cancel := context.WithTimeout(r.Context(), 500*time.Millisecond)
			defer cancel()

			if err := redisClient.Del(ctx, cacheKey).Err(); err != nil {
				log.Warnf("[REDIS] DEL %s failed: %v", cacheKey, err)
			}
		}
		u.Respond(w, u.Message(true, "OK"))
	}
}

var NoteQuery = func(w http.ResponseWriter, r *http.Request) {
	var notes []models.Note
	db := db.GetDB()

	query, ok := r.URL.Query()["query"]
	if !ok || len(query[0]) < 1 {
		err := db.Find(&notes).Error
		if err != nil {
			u.HandleBadRequest(w, err)
			return
		}
	} else {
		q := "%" + query[0] + "%"
		err := db.Where("title LIKE ? OR content LIKE ?", q, q).Find(&notes).Error
		if err != nil {
			u.HandleBadRequest(w, err)
			return
		}
	}

	res, err := json.Marshal(notes)
	if err != nil {
		u.HandleBadRequest(w, err)
	} else {
		u.RespondJSON(w, res)
	}
}
