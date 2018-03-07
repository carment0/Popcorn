/**
 * @copyright Popcorn, 2018
 * @author Calvin Feng
 */

import merge from 'lodash/merge';
import { MOVIE_TRAILER_FETCH_SUCCESS } from './trailer.action';


const YouTube = 'YouTube';

/**
 * Movie trailer state is a map of movie ID to an array of YouTube video keys. The key is used for loading the video from
 * YouTube directly.
 * @param {object} state
 * @param {object} action
 */
export default function movieTrailerReducer(state = {}, action) {
  Object.freeze(state);

  switch (action.type) {
    case MOVIE_TRAILER_FETCH_SUCCESS:
      const results = action.payload.results;

      // Currently we only support movie trailers on YouTube, will consider supporting Vimeo in the future.
      const youtubeVidoeKeyList = [];
      if (results !== undefined && results.length > 0) {
        results.forEach((trailerSearchResult) => {
          if (trailerSearchResult.site === YouTube) {
            youtubeVidoeKeyList.push(trailerSearchResult.key);
          }
        });
      }

      const newState = {};
      if (youtubeVidoeKeyList.length > 0) {
        newState[action.imdbId] = youtubeVidoeKeyList;
      }

      return merge({}, state, newState);

    default:
      return state;
  }
}
