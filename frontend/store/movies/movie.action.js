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
 * There are two types of recommendation fetching: (1) user is authenticated and user ID is provided, (2) user is guest
 * and does not have an account. If user is authenticated, we only need to make a request to the endpoint because server
 * will authenticate using client's session token. However, just to be safe, send the user ID anyways. If user is not
 * authenticated, then we must submit a map of movie ID to movie ratings.
 * @param {number} userId
 * @param {object} ratingMap
 * @returns {Promise}
 */
export const recommendedMoviesFetch = (session, ratings) => (dispatch) => {
  if (session.currentUser) {
    return request.get(`api/users/${session.currentUser.username}/recommend`)
      .then((res) => {
        return dispatch({
          type: RECOMMENDED_MOVIES_FETCH_SUCCESS,
          payload: res.data
        });
      })
      .catch((error) => dispatch({ type: RECOMMENDED_MOVIES_FETCH_FAIL, error }));
  }

  return request.post('api/movies/recommend', { ratings })
    .then((res) => {
      return dispatch({
        type: RECOMMENDED_MOVIES_FETCH_SUCCESS,
        payload: res.data
      });
    })
    .catch((error) => dispatch({ type: RECOMMENDED_MOVIES_FETCH_FAIL, error }));
};
