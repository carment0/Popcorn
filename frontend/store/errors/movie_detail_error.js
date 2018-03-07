/**
 * @copyright Popcorn, 2018
 * @author Calvin Feng
 */
import merge from 'lodash/merge';
import { MOVIE_DETAIL_FETCH_FAIL } from '../movies/detail.action';


export default (state = {}, action) => {
  Object.freeze(state);
  switch (action.type) {
    case MOVIE_DETAIL_FETCH_FAIL:
      const newState = {};
      newState[action.imdbId] = action.error;
      return merge({}, state, newState);

    default:
      return state;
  }
};
