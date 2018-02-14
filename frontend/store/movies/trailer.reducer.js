/**
 * @copyright Popcorn, 2018
 * @author Calvin Feng
 */

import merge from 'lodash/merge';
import { MOVIE_TRAILER_FETCH_SUCCESS, MOVIE_TRAILER_FETCH_FAIL } from './trailer.action';


export default function movieTrailerReducer(state = {}, action) {
  Object.freeze(state);

  switch (action.type) {
    case MOVIE_TRAILER_FETCH_SUCCESS:
      const trailerSearchResults = action.payload.results;

      const videoSourceURLs = [];
      trailerSearchResults.forEach((trailerSearchResult) => {
        if (trailerSearchResult.site === 'YouTube') {
          videoSourceURLs.push(`https://www.youtube.com/watch?v=${trailerSearchResult.key}`);
        }
      });

      const trailerMap = {};

      if (videoSourceURLs.length > 0) {
        trailerMap[action.imdbId] = videoSourceURLs;
      }

      return merge({}, state, trailerMap);

    case MOVIE_TRAILER_FETCH_FAIL:
      console.log(`Failed to fetch movie trailer for movie ${action.imdbId}`, action.error);
      return state;

    default:
      return state;
  }
}
