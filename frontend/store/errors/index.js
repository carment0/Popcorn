/**
 * @copyright Popcorn, 2018
 * @author Carmen To, Calvin Feng
 */

import { combineReducers } from 'redux';
import SessionErrorReducer from './session_error';
import MovieDetailErrorReducer from './movie_detail_error';
import MovieTrailerErrorReducer from './movie_trailer_error';

export default combineReducers({
  session: SessionErrorReducer,
  movieDetails: MovieDetailErrorReducer,
  movieTrailers: MovieTrailerErrorReducer
});
