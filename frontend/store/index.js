/**
 * @copyright Consilium, 2018
 * @author Calvin Feng
 */

// Thirdparty imports
import { combineReducers } from 'redux';
import { createStore, applyMiddleware } from 'redux';
import thunk from 'redux-thunk';
import logger from 'redux-logger';

import UserPreferenceReducer from './users/preference.reducer';

const rootReducer = combineReducers({
  userPreference: UserPreferenceReducer
});

const preloadedState = {};

if (window.currentUser) {
  preloadedState.sessions = { currentUser: window.currentUser };
  delete window.currentUser;
}

export default createStore(rootReducer, preloadedState, applyMiddleware(thunk, logger));