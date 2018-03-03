// Copyright (c) 2018 Popcorn
// Author(s) Carmen To
package kmeans

import (
  "encoding/csv"
  "os"
  "strconv"
  "io"
)

type Movie struct {
  MovieID string
  Feature []float64
}

func ReadFromCSV(filePath string) ([]*Movie, error) {
  var movies []*Movie

  csvFile, fileError := os.Open(filePath)

  if fileError != nil {
    return nil, fileError
  }

  reader := csv.NewReader(csvFile)
  reader.Read()

  for {
    if row, err := reader.Read(); err != nil {
      if err == io.EOF {
        break
      }
    } else {
      var feature []float64
      for i := 1; i < len(row); i += 1 {
        stringFeat := row[i]
        integer, err := strconv.ParseFloat(stringFeat, 64)

        if err != nil {
          return nil, err
        }
        feature = append(feature, integer)
      }

      movies = append(movies, &Movie {
        MovieID: row[0],
        Feature: feature,
      })
    }
  }
  return movies, nil
}

func WriteToCSV(filepath string, clustData []*MovieAssignments) error {
  csvFile, fileError := os.Create(filepath)
  if fileError != nil {
    return fileError
  }

  writer := csv.NewWriter(csvFile)
  header := []string{"movieId", "centGroup", "close1", "close2", "close3", "close4", "far1", "far2", "far3", "far4",}

  if err := writer.Write(header); err != nil {
    return err
  }

  for _, movie := range clustData {
    row := []string{movie.Movie.MovieID, strconv.Itoa(movie.Centroid.ClusterID)}
    for _, closest := range movie.ClosestClusters {
      row = append(row, strconv.Itoa(closest.ClusterID))
    }

    for _, farthest := range movie.FarthestClusters {
      row = append(row, strconv.Itoa(farthest.ClusterID))
    }
    writer.Write(row)
  }
  writer.Flush()
  return nil
}
