/**
 * @copyright Popcorn, 2018
 * @author Calvin Feng, Carmen To
 */

import { SET_MOVIE_POPULARITY_PERCENTILE } from './percentile.action';


/**
 * Movie popularity percentile state is simple a number, ranging from 0 to 100
 * @param {number} state
 * @param {object} action
 */
export default function moviePopularityPercentileReducer(state = 100, action) {
  switch (action.type) {
    case SET_MOVIE_POPULARITY_PERCENTILE:
      return action.percentile;
    default:
      return state;
  }
}
