/**
 * @copyright Popcorn, 2018
 * @author Calvin Feng
 */

import merge from 'lodash/merge';
import { combineReducers } from 'redux';
import { ALL_MOVIES_FETCH_SUCCESS, ALL_MOVIES_FETCH_FAIL } from './movie.action';
import { POPULAR_MOVIES_FETCH_SUCCESS, POPULAR_MOVIES_FETCH_FAIL } from './movie.action';
import { RECOMMENDED_MOVIES_FETCH_SUCCESS, RECOMMENDED_MOVIES_FETCH_FAIL } from './movie.action';
import { MOVIE_SKIPPED } from './movie.action';
import { MOVIE_RATING_POST_SUCCESS, MOVIE_RATING_RECORDED, MOVIE_RATINGS_FETCH_SUCCESS } from './rating.action';


/**
 * All movie state is an object that maps movie ID to a dictionary of movie information such as IMDB ID, year, average
 * ratings and etc...
 * @param {Set} state
 * @param {object} action
 */
function allMovieReducer(state = {}, action) {
  Object.freeze(state);

  switch (action.type) {
    case ALL_MOVIES_FETCH_SUCCESS:
      const newState = {};
      action.payload.forEach((movie) => {
        newState[movie.id] = movie;
      });
      return merge({}, state, newState);

    case ALL_MOVIES_FETCH_FAIL:
      console.log('Failed to fetch all movies', action.error);
      return state;

    default:
      return state;
  }
}

/**
 * Popular movie state is a set of movie ID's
 * @param {Set} state
 * @param {object} action
 */
function popularMovieReducer(state = new Set(), action) {
  Object.freeze(state);

  switch (action.type) {
    case POPULAR_MOVIES_FETCH_SUCCESS:
      const newState = new Set(state);
      action.payload.forEach((movie) => {
        newState.add(movie.id);
      });
      return newState;

    case POPULAR_MOVIES_FETCH_FAIL:
      console.log('Failed to fetch popular movies', action.error);
      return state;

    default:
      return state;
  }
}

/**
 * Recommended movie state is a set of movie ID's
 * @param {Set} state
 * @param {object} action
 */
function recommendedMovieReducer(state = new Set(), action) {
  Object.freeze(state);

  switch (action.type) {
    case RECOMMENDED_MOVIES_FETCH_SUCCESS:
      const newState = new Set();
      action.payload.forEach((movie) => {
        newState.add(movie.id);
      });
      return newState;

    case RECOMMENDED_MOVIES_FETCH_FAIL:
      console.log('Failed to fetch recommended movies', action.error);
      return state;

    default:
      return state;
  }
}

/**
 * Skipped movie state is a set of movie ID's
 * @param {Set} state
 * @param {object} action
 */
function skippedMovieReducer(state = new Set(), action) {
  Object.freeze(state);

  switch (action.type) {
    case MOVIE_SKIPPED:
      const newState = new Set(state);
      newState.add(action.movieId);
      return newState;

    default:
      return state;
  }
}

/**
 * Rated movie is a set of movie ID's
 * @param {Set} state
 * @param {object} action
 */
function ratedMovieReducer(state = new Set(), action) {
  Object.freeze(state);

  let newState;
  switch (action.type) {
    case MOVIE_RATING_POST_SUCCESS:
    case MOVIE_RATING_RECORDED:
      newState = new Set(state);
      newState.add(action.movieId);
      return newState;

    case MOVIE_RATINGS_FETCH_SUCCESS:
      newState = new Set(state);
      action.payloads.forEach((payload) => {
        newState.add(payload.movie_id);
      });
      return newState;

    default:
      return state;
  }
}

export default combineReducers({
  all: allMovieReducer,
  popular: popularMovieReducer,
  recommended: recommendedMovieReducer,
  rated: ratedMovieReducer,
  skipped: skippedMovieReducer
});
