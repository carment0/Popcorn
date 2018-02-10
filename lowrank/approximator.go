// Copyright (c) 2018 Popcorn
// Author(s) Calvin Feng

// Package lowrank provides tools to perform low rank approximation on latent features of movies and users.
package lowrank

import (
	"gonum.org/v1/gonum/mat"
	"github.com/sirupsen/logrus"
)

type Approximator struct {
	UserLatent  *mat.Dense
	MovieLatent *mat.Dense
	Rating      *mat.Dense
}

func NewApproximator(R *mat.Dense, K int) *Approximator {
	I, J := R.Dims()
	return &Approximator{
		UserLatent:  RandMat(I, K),
		MovieLatent: RandMat(J, K),
		Rating:      R,
	}
}

func (a *Approximator) ModelPredict() (*mat.Dense, error) {
	I, KI := a.UserLatent.Dims()
	J, KJ := a.MovieLatent.Dims()

	if KI != KJ {
		return nil, mat.ErrShape
	}

	result := mat.NewDense(I, J, nil)
	result.Mul(a.UserLatent, a.MovieLatent.T())
	return result, nil
}

func (a *Approximator) Loss(reg float64) (float64, float64, error) {
	prediction, err := a.ModelPredict()
	if err != nil {
		return 0, 0, err
	}

	I, J := prediction.Dims()
	diff := mat.NewDense(I, J, nil)
	diff.Sub(prediction, a.Rating)
	avgDiscrepancy := AbsAverage(diff)
	diff.MulElem(diff, diff)

	loss := 0.5 * mat.Sum(diff)

	USquared := mat.DenseCopyOf(a.UserLatent)
	USquared.MulElem(USquared, USquared)
	loss += reg * mat.Sum(USquared) / 2.0

	MSquared := mat.DenseCopyOf(a.MovieLatent)
	MSquared.MulElem(MSquared, MSquared)
	loss += reg * mat.Sum(MSquared) / 2.0

	return loss, avgDiscrepancy, nil
}

func (a *Approximator) Gradients(reg float64) (*mat.Dense, *mat.Dense, error) {
	prediction, err := a.ModelPredict()
	if err != nil {
		return nil, nil, err
	}

	I, J := prediction.Dims()
	GradR := mat.NewDense(I, J, nil)
	GradR.Sub(prediction, a.Rating)

	_, K := a.UserLatent.Dims()

	GradU := mat.NewDense(I, K, nil)
	GradU.Mul(GradR, a.MovieLatent)
	RegU := mat.NewDense(I, K, nil)
	RegU.Scale(reg, a.UserLatent)
	GradU.Add(GradU, RegU)

	GradM := mat.NewDense(J, K, nil)
	GradM.Mul(GradR.T(), a.UserLatent)
	RegM := mat.NewDense(J, K, nil)
	RegM.Scale(reg, a.MovieLatent)
	GradM.Add(GradM, RegM)

	return GradU, GradM, nil
}

func (a *Approximator) Train(steps int, epochSize int, reg float64, learnRate float64) {
	I, _ := a.UserLatent.Dims()
	J, _ := a.MovieLatent.Dims()
	for step := 0; step < steps; step += 1 {
		if step%epochSize == 0 {
			loss, avgDiscrepancy, _ := a.Loss(reg)
			logrus.Infof("%d: net loss %.2f, avg loss %.8f, and average discrepancy from true value %.8f \n",
				step, loss, loss/float64(I*J), avgDiscrepancy,
			)
		}

		if GradU, GradM, err := a.Gradients(reg); err == nil {
			GradU.Scale(learnRate, GradU)
			a.UserLatent.Sub(a.UserLatent, GradU)

			GradM.Scale(learnRate, GradM)
			a.MovieLatent.Sub(a.MovieLatent, GradM)
		}
	}
}

func (a *Approximator) ApproximateUserLatent(steps int, epochSize int, reg float64, learnRate float64) {
	I, _ := a.UserLatent.Dims()
	J, _ := a.MovieLatent.Dims()
	for step := 0; step < steps; step += 1 {
		if step%epochSize == 0 {
			loss, avgDiscrepancy, _ := a.Loss(reg)
			logrus.Infof("%d: net loss %.2f, avg loss %.8f, and average discrepancy from true value %.8f \n",
				step, loss, loss/float64(I*J), avgDiscrepancy,
			)
		}

		if GradU, _, err := a.Gradients(reg); err == nil {
			GradU.Scale(learnRate, GradU)
			a.UserLatent.Sub(a.UserLatent, GradU)
		}
	}
}