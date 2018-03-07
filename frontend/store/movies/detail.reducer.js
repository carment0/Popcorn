/**
 * @copyright Popcorn, 2018
 * @author Calvin Feng
 */

import merge from 'lodash/merge';
import { MOVIE_DETAIL_FETCH_SUCCESS } from './detail.action';

/**
 * Movie detail state is a map of movie ID to a dictionary of details of the movie such as movie poster, movie overview,
 * cast, director and etc...
 * @param {object} state
 * @param {object} action
 */
export default function movieDetailReducer(state = {}, action) {
  Object.freeze(state);

  switch (action.type) {
    case MOVIE_DETAIL_FETCH_SUCCESS:
      const newState = {};

      // NOTE: This is the stupid part, sometimes TMDB API does not comply with its own API standard. I have to manually
      // check the existence of each field.
      let imdbId;
      if (action.payload.imdb_id) {
        imdbId = action.payload.imdb_id;
      }

      let title;
      if (action.payload.title) {
        title = action.payload.title;
      }

      let year;
      if (action.payload.year) {
        year = action.payload.year;
      }

      let plot;
      if (action.payload.overview) {
        plot = action.payload.overview;
      }

      let poster;
      if (action.payload.poster_path) {
        poster =  `https://image.tmdb.org/t/p/w300${action.payload.poster_path}`;
      }

      if (imdbId) {
        newState[imdbId] = { title, year, plot, poster };
      }

      return merge({}, state, newState);

    default:
      return state;
  }
}
