package main

import (
  "encoding/csv"
  "fmt"
  "os"
  "strconv"
  "io"
  "math/rand"
  "math"
  "time"
)

func ReadFromCSV(filepath string) (map[string][]float64, error) {
  var featureMap = make(map[string][]float64)
  csvFile, fileError := os.Open(filepath)

  if fileError != nil {
    return featureMap, fileError
  }

  reader := csv.NewReader(csvFile)
  reader.Read()

  for {
    if row, err := reader.Read(); err != nil {
      if err == io.EOF {
        break
      }
    } else {
      var features []float64
      for i := 1; i <= 10; i += 1 {
        string := row[i]
        intager, err := strconv.ParseFloat(string, 64)
        if err != nil {
          return featureMap, err
        }
        features = append(features, intager)
      }
      featureMap[row[0]] = features
    }
  }
  return featureMap, nil
}

// func WriteToCSV(filepath string, clustData map[string]int) error {
//   csvFile, fileError := os.Create(filepath)
//   if fileError != nil {
//     return fileError
//   }
//
//   writer := csv.NewWriter(csvFile)
//   header := []string{"movieId", "centGroup"}
//
//   if err := writer.Write(header); err != nil {
//     return err
//   }
//
//   for k, v := range clustData {
//     row := []string{k, strconv.Itoa(v)}
//     writer.Write(row)
//   }
//   writer.Flush()
//   return nil
// }

func main() {
  featureMap, err := ReadFromCSV("../dataset/features.csv")
  if err != nil {
    fmt.Println("Failed to read CSV", err)
  } else {
    fmt.Println("Done reading")
    movieIdKeys := movieIdKeys(featureMap)
    centroids := initCentroids(featureMap, movieIdKeys, 10)
    findKMeans(featureMap, centroids)
  }
}

func findKMeans(featureMap map[string][]float64, centroids map[int][]float64) {
  var currentCent = centroids
  // var changeCent = true;
  // var updatedCent map[int][]float64
  var movieCentAssignment map[string]int

  movieCentAssignment = assigningClosetCent(featureMap, currentCent)
  fmt.Println(movieCentAssignment)
  // for changeCent == true {
  // }
}

func assigningClosetCent(featureMap map[string][]float64, centroids map[int][]float64) map[string]int {
  var assigned map[string]int
  var minDist float64
  var centID int

  assigned = make(map[string]int)
  for k1, v1 := range featureMap {
    for k2, v2 := range centroids {
      dist := getDistance(v1, v2)
      if minDist == 0 && centID == 0 || minDist > dist {
        minDist = dist
        centID = k2
      }
    }

    assigned[k1] = centID
  }
  return assigned
}

func getDistance(movieFeat, centPoint []float64) float64 {
  var dist float64
  for idx, val := range movieFeat {
    x := val - centPoint[idx]
    dist += math.Pow(x, 2)
  }
  return math.Sqrt(dist)
}

func initCentroids(featureMap map[string][]float64, keys []string, centCount int) map[int][]float64 {
  var centroids = make(map[int][]float64)
  for i := 0; i < centCount; i += 1 {
    rand.Seed(int64(time.Now().Nanosecond()))
    idx := rand.Intn(100)
    movieId := keys[idx]
    feature := featureMap[movieId]
    centroids[i] = feature
  }
  return centroids
}

func movieIdKeys(featureMap map[string][]float64) []string {
  var keys []string
  for k := range featureMap {
      keys = append(keys, k)
  }
  return keys
}