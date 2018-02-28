import csv
import numpy as np


DIR = '../datasets/26m/'


def cosine_similarity(a, b):
    a_mag = (a * a).sum()
    b_mag = (b * b).sum()
    return a.dot(b) / (a_mag * b_mag)


def geometric_distance(a, b):
    diff = a - b
    return diff.dot(diff)


def main():
    movies = dict()
    with open(DIR + 'movies.csv', 'rb') as csvfile:
        reader = csv.reader(csvfile, delimiter=',')
        next(reader, None)  # Skip the header
        for row in reader:
            movie_id, title, _ = row
            movies[movie_id] = title

    features = dict()
    with open('../datasets/production/features.csv', 'rb') as csvfile:
        reader = csv.reader(csvfile, delimiter=',')
        next(reader, None)  # Skip the header
        for row in reader:
            movie_id = row[0]
            features[movie_id] = np.array(row[1:], dtype=float)

    chosen_id = '1'

    distances = []
    for key in features:
        if key == chosen_id:
            continue

        distances.append((key, geometric_distance(features[key], features[chosen_id])))

    distances = sorted(distances, key=lambda x: x[1], reverse=False)

    print "Choosen movie is %s" % movies[chosen_id]
    for item in distances[0:100]:
        print "Nearest neigbor: %s with distance %.3f" % (movies[item[0]], item[1])


if __name__ == '__main__':
    main()
