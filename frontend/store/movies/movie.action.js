/**
 * @copyright Popcorn, 2018
 * @author Calvin Feng
 */

import request from 'axios';


export const MOVIE_SKIPPED = 'MOVIE_SKIPPED';
export const movieSkip = (movieId) => (dispatch) => dispatch({ type: MOVIE_SKIPPED, movieId });

export const ALL_MOVIES_FETCH_SUCCESS = 'ALL_MOVIES_FETCH_SUCCESS';
export const ALL_MOVIES_FETCH_FAIL = 'ALL_MOVIES_FETCH_FAIL';

/**
 * Fetch the list of all movies from the database.
 * @returns {Promise}
 */
export const allMoviesFetch = () => (dispatch) => {
  return request.get('api/movies')
    .then((res) => {
      return dispatch({
        type: ALL_MOVIES_FETCH_SUCCESS,
        payload: res.data
      });
    })
    .catch((error) => dispatch({ type: ALL_MOVIES_FETCH_FAIL, error }));
};

export const POPULAR_MOVIES_FETCH_SUCCESS = 'POPULAR_MOVIES_FETCH_SUCCESS';
export const POPULAR_MOVIES_FETCH_FAIL = 'POPULAR_MOVIES_FETCH_FAIL';

/**
 * Fetch movies ranked by popularity from the server. Popularity is determined by number of views and average rating
 * score. The current endpoint returns all movie and this may become a problem. It is better to provide a limit.
 * @returns {Promise}
 */
export const popularMoviesFetch = () => (dispatch) => {
  return request.get('api/movies/popular')
    .then((res) => {
      return dispatch({
        type: POPULAR_MOVIES_FETCH_SUCCESS,
        payload: res.data
      });
    })
    .catch((error) => dispatch({ type: POPULAR_MOVIES_FETCH_FAIL, error }));
};

export const RECOMMENDED_MOVIES_FETCH_SUCCESS = 'RECOMMENDED_MOVIES_FETCH_SUCCESS';
export const RECOMMENDED_MOVIES_FETCH_FAIL = 'RECOMMENDED_MOVIES_FETCH_FAIL';

/**
 * General recommendations are for anonymouse users who are new to our site. This function requires a map of movie ID to
 * rating value submitted by the user.
 * @param {object} ratingMap
 * @returns {Promise}
 */
export const recommendedMoviesFetch = (session, ratings) => (dispatch) => {
  return request.post('api/movies/recommend', { ratings })
    .then((res) => {
      return dispatch({
        type: RECOMMENDED_MOVIES_FETCH_SUCCESS,
        payload: res.data
      });
    })
    .catch((error) => dispatch({ type: RECOMMENDED_MOVIES_FETCH_FAIL, error }));
};

/**
 * Personalized recommendations are for authenticated users who have rated a good number of movies
 * @param {object} session
 * @param {object} yearRange
 * @param {number} percentile
 * @param {Set} skipped
 * @returns {Promise}
 */
export const personalizedRecommendedMoviesFetch = (session, yearRange, percentile, skipped) => (dispatch) => {
  const payload = {
    max: yearRange.maxYear,
    min: yearRange.minYear,
    percent: percentile,
    skipped: [...skipped]
  };

  return request.post(`api/users/${session.currentUser.id}/recommend`, payload)
    .then((res) => {
      return dispatch({
        type: RECOMMENDED_MOVIES_FETCH_SUCCESS,
        payload: res.data
      });
    })
    .catch((error) => dispatch({ type: RECOMMENDED_MOVIES_FETCH_FAIL, error }));
};
