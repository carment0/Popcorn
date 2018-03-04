// Copyright (c) 2018 Popcorn
// Author(s) Calvin Feng

// Package lowrank provides tools to perform low rank factorization on latent features of movies and users.
package lowrank

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/mat"
	"math"
)

type Factorizer struct {
	UserLatent      *mat.Dense
	MovieLatent     *mat.Dense
	Rating          *mat.Dense
	MatrixConverter *MatrixConverter
}

// NewFactorizer requires either a converter or a rating matrix. Converter is needed for training phase. Rating matrix
// is needed for running new user preference approximation in the recommendation engine.
func NewFactorizer(converter *MatrixConverter, ratingMat *mat.Dense, K int) *Factorizer {
	if converter == nil && ratingMat == nil {
		return nil
	}

	if converter != nil {
		R := converter.GetRatingMatrix()
		I, J := R.Dims()
		return &Factorizer{
			UserLatent:      RandMat(I, K),
			MovieLatent:     RandMat(J, K),
			Rating:          R,
			MatrixConverter: converter,
		}
	}

	I, J := ratingMat.Dims()
	return &Factorizer{
		UserLatent:  RandMat(I, K),
		MovieLatent: RandMat(J, K),
		Rating:      ratingMat,
	}
}

func (f *Factorizer) ModelPredict() (*mat.Dense, error) {
	I, KI := f.UserLatent.Dims()
	J, KJ := f.MovieLatent.Dims()

	if KI != KJ {
		return nil, mat.ErrShape
	}

	result := mat.NewDense(I, J, nil)
	result.Mul(f.UserLatent, f.MovieLatent.T())
	return result, nil
}

func (f *Factorizer) Loss(reg float64) (float64, float64, error) {
	prediction, err := f.ModelPredict()
	if err != nil {
		return 0, 0, err
	}

	var rootMeanSqError float64
	if f.MatrixConverter != nil {
		testCount := 0.0

		for userID := range f.MatrixConverter.TestRatingMap {
			for movieID := range f.MatrixConverter.TestRatingMap[userID] {
				i := f.MatrixConverter.UserIDToIndex[userID]
				j := f.MatrixConverter.MovieIDToIndex[movieID]
				rootMeanSqError += math.Pow(prediction.At(i, j)-f.MatrixConverter.TestRatingMap[userID][movieID], 2)
				testCount += 1.0
			}
		}

		rootMeanSqError /= testCount
		rootMeanSqError = math.Sqrt(rootMeanSqError)
	}

	I, J := prediction.Dims()
	diff := mat.NewDense(I, J, nil)
	diff.Sub(prediction, f.Rating)
	diff.MulElem(diff, diff)

	// Ignore the difference generated by zero values of R matrix
	for i := 0; i < I; i += 1 {
		for j := 0; j < J; j += 1 {
			if f.Rating.At(i, j) == 0 {
				diff.Set(i, j, 0)
			}
		}
	}
	loss := 0.5 * mat.Sum(diff)

	USquared := mat.DenseCopyOf(f.UserLatent)
	USquared.MulElem(USquared, USquared)
	loss += reg * mat.Sum(USquared) / 2.0

	MSquared := mat.DenseCopyOf(f.MovieLatent)
	MSquared.MulElem(MSquared, MSquared)
	loss += reg * mat.Sum(MSquared) / 2.0

	return loss, rootMeanSqError, nil
}

func (f *Factorizer) Gradients(reg float64) (*mat.Dense, *mat.Dense, error) {
	prediction, err := f.ModelPredict()
	if err != nil {
		return nil, nil, err
	}

	I, J := prediction.Dims()
	GradR := mat.NewDense(I, J, nil)
	GradR.Sub(prediction, f.Rating)

	// Prevent the zero values of R from back-propagating gradients to M and U
	for i := 0; i < I; i += 1 {
		for j := 0; j < J; j += 1 {
			if f.Rating.At(i, j) == 0 {
				GradR.Set(i, j, 0)
			}
		}
	}

	_, K := f.UserLatent.Dims()

	GradU := mat.NewDense(I, K, nil)
	GradU.Mul(GradR, f.MovieLatent)
	RegU := mat.NewDense(I, K, nil)
	RegU.Scale(reg, f.UserLatent)
	GradU.Add(GradU, RegU)

	GradM := mat.NewDense(J, K, nil)
	GradM.Mul(GradR.T(), f.UserLatent)
	RegM := mat.NewDense(J, K, nil)
	RegM.Scale(reg, f.MovieLatent)
	GradM.Add(GradM, RegM)

	return GradU, GradM, nil
}

func (f *Factorizer) Train(steps int, epochSize int, reg float64, learnRate float64) {
	for step := 0; step < steps; step += 1 {
		if step%epochSize == 0 {
			loss, rootMeanSqError, _ := f.Loss(reg)

			var logMessage string
			if f.MatrixConverter == nil {
				logMessage = fmt.Sprintf("iteration %3d: net loss %5.2f", step, loss)
			} else {
				logMessage = fmt.Sprintf(`iteration %3d: net loss %5.2f and RMSE %1.8f`, step, loss, rootMeanSqError)
			}

			logrus.WithField("file", "lowrank.factorizer").Info(logMessage)
		}

		if GradU, GradM, err := f.Gradients(reg); err == nil {
			GradU.Scale(learnRate, GradU)
			f.UserLatent.Sub(f.UserLatent, GradU)

			GradM.Scale(learnRate, GradM)
			f.MovieLatent.Sub(f.MovieLatent, GradM)
		}
	}
}

func (f *Factorizer) ApproximateUserLatent(steps int, epochSize int, reg float64, learnRate float64) {
	I, _ := f.UserLatent.Dims()
	J, _ := f.MovieLatent.Dims()
	for step := 0; step < steps; step += 1 {
		if step%epochSize == 0 {
			loss, _, _ := f.Loss(reg)

			logMessage := fmt.Sprintf("iteration %3d: net loss %5.2f, avg loss %1.8f on %d movies",
				step, loss, loss/float64(I*J), I*J,
			)

			logrus.WithField("src", "lowrank.factorizer").Info(logMessage)
		}

		if GradU, _, err := f.Gradients(reg); err == nil {
			GradU.Scale(learnRate, GradU)
			f.UserLatent.Sub(f.UserLatent, GradU)
		}
	}
}
