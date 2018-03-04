package main

import (
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/mat"
	"popcorn/lowrank"
	"popcorn/model"
)

type RecommendEngine struct {
	// Database connection is necessary for updating models. Also gorm is thread safe, we can have up to N connections
	// open from separate go routines.
	DBConn *gorm.DB

	// Connection map handles the mapping of user ID to web socket connection to whichever client who initiated the long
	// running task. We definitely need a sync/mutex for writing to ConnMap in handlers.
	ConnMap map[uint]*websocket.Conn
}

func NewRecommendEngine(db *gorm.DB, connMap map[uint]*websocket.Conn) *RecommendEngine {
	return &RecommendEngine{
		DBConn:  db,
		ConnMap: connMap,
	}
}

type Notification struct {
	UserID  uint   `json:"user_id"`
	Message string `json:"message"`
}

func (re *RecommendEngine) ListenToInbound(queue chan *model.User) {
	for user := range queue {
		var movies []model.Movie
		if err := re.DBConn.Order("id asc").Find(&movies).Error; err != nil {
			logrus.WithField(
				"src", "engine.go",
			).Error("RecommendEngine has failed to load all movies from database", err)
			continue
		}

		// In case that database was not seeded properly and no movies are found in the database, we should not proceed
		// with the algorithm.
		if len(movies) == 0 {
			logrus.WithField("src", "main.engine").Error("No movies are found in the database!")
			continue
		}

		featureDim := len(movies[0].Feature)

		// Construct a user rating map, mapping movie ID to user submitted movie rating value.
		ratingMapByID := make(map[uint]float64)
		for _, rating := range user.Ratings {
			ratingMapByID[rating.MovieID] = rating.Value
		}

		// Allocate 0 length and K * M capacity for latent feature slice. Note: K is feature dimension and M is number of
		// rated movies by the user. Also create a rating matrix, which has a dimension of (1, M). Because we are only
		// computing latent preference for one user.
		M := len(ratingMapByID)
		movieFeatureData := make([]float64, 0, featureDim*M)
		ratingMat := mat.NewDense(1, M, nil)
		j := 0
		for _, movie := range movies {
			if val, ok := ratingMapByID[movie.ID]; ok {
				movieFeatureData = append(movieFeatureData, movie.Feature...)
				ratingMat.Set(0, j, val)
				j += 1
			}
		}

		approximator := lowrank.NewFactorizer(nil, ratingMat, featureDim)
		approximator.MovieLatent = mat.NewDense(M, featureDim, movieFeatureData)
		approximator.ApproximateUserLatent(300, 1, 0, 0.01)

		if len(approximator.UserLatent.RawRowView(0)) == featureDim {
			user.Preference = approximator.UserLatent.RawRowView(0)
			if err := re.DBConn.Save(user).Error; err != nil {
				logrus.WithField(
					"src", "main.engine",
				).Error("failed to save preference to user model", err)
			} else {
				logrus.WithField(
					"src", "main.engine",
				).Infof("preference for user %s is saved", user.Username)
				// re.ConnMap[user.ID].WriteJSON(Notification{UserID: user.ID, Message: "Preference is ready!"})
			}
		} else {
			logrus.WithField(
				"src", "main.engine",
			).Errorf(`something went wrong with approximator, user latent preference vector does not have
				the correct length; it has %d but expected %d`,
				len(approximator.UserLatent.RawRowView(0)),
				featureDim,
			)
		}
	}
}

func (re *RecommendEngine) approximateUserPreference(user *model.User) error {
	// Refactor the above logic into here.
	return nil
}
