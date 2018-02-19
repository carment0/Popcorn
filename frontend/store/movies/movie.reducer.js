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

function popularMovieReducer(state = new Set(), action) {
  Object.freeze(state);

  switch (action.type) {
    case POPULAR_MOVIES_FETCH_SUCCESS:
      action.payload.forEach((movie) => {
        state.add(movie.id);
      });
      return state;

    case POPULAR_MOVIES_FETCH_FAIL:
      console.log('Failed to fetch popular movies', action.error);
      return state;

    default:
      return state;
  }
}

function recommendedMovieReducer(state = new Set(), action) {
  Object.freeze(state);

  switch (action.type) {
    case RECOMMENDED_MOVIES_FETCH_SUCCESS:
      action.payload.forEach((movie) => {
        state.add(movie.id);
      });
      return state;

    case RECOMMENDED_MOVIES_FETCH_FAIL:
      console.log('Failed to fetch recommended movies', action.error);
      return state;

    default:
      return state;
  }
}

function skippedMovieReducer(state = new Set(), action) {
  Object.freeze(state);

  switch (action.type) {
    case MOVIE_SKIPPED:
      state.add(action.movieId);
      return state;

    default:
      return state;
  }
}

function ratedMovieReducer(state = new Set(), action) {
  Object.freeze(state);

  switch (action.type) {
    case MOVIE_RATING_POST_SUCCESS:
    case MOVIE_RATING_RECORDED:
      state.add(action.movieId);
      return state;

    case MOVIE_RATINGS_FETCH_SUCCESS:
      action.payloads.forEach((payload) => {
        state.add(payload.movie_id);
      });
      return state;

    default:
      return state;
  }
}

// NOTE: The keys (popular, recommended, and skipped) are holding a set while the key 'all' is holding all
// the data of the movies. This is to minimize data duplication and to improve performance of front end.
export default combineReducers({
  all: allMovieReducer,
  popular: popularMovieReducer,
  recommended: recommendedMovieReducer,
  rated: ratedMovieReducer,
  skipped: skippedMovieReducer
});
