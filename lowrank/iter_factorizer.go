// Copyright (c) 2018 Popcorn
// Author(s) Calvin Feng

// Package lowrank provides tools to perform low rank factorization on latent features of movies and users.
package lowrank

// IterativeFactorizer does not use matrices at all. Instead, it holds each user's preference vector and movie's feature
// vector in a map. It does not cache predicted rating; it only computes it when it is needed. Although vectorized
// implementations lead to better time performance, it is extremely space hungry. For 20,000 users and 45,000 movies,
// the number of generated predicted rating is 900 millions. Each float64 is 8 bytes, and that is 7.2 billion bytes of
// memory.
type IterativeFactorizer struct {
}
