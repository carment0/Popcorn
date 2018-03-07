/**
 * @copyright Popcorn, 2018
 * @author Carmen To
 */

import request from 'axios';


export const CLEAR_SESSION_ERRORS = 'CLEAR_SESSION_ERRORS';
export const clearSessionErrors = () => ({ type: CLEAR_SESSION_ERRORS });

export const SIGNUP_SUCCESS = 'SIGNUP_SUCCESS';
export const SIGNUP_FAIL = 'SIGNUP_FAIL';

export const signup = (user) => (dispatch) => {
  return request.post('api/users/register', user).then((res) => {
    return dispatch({
      type: SIGNUP_SUCCESS,
      currentUser: res.data
    });
  }).catch((err) => dispatch({ type: SIGNUP_FAIL, error: err.response.data.error }));
};

export const LOGIN_SUCCESS = 'LOGIN_SUCCESS';
export const LOGIN_FAIL = 'LOGIN_FAIL';

export const login = (user) => (dispatch) => {
  return request.post('api/users/login', user).then((res) => {
    return dispatch({
      type: LOGIN_SUCCESS,
      currentUser: res.data
    });
  }).catch((err) => dispatch({ type: LOGIN_FAIL, error: err.response.data.error }));
};

export const LOGOUT_SUCCESS = 'LOGOUT_SUCCESS';
export const LOGOUT_FAIL = 'LOGOUT_FAIL';

export const logout = () => (dispatch) => {
  return request.delete('api/users/logout').then(() => {
    return dispatch({
      type: LOGOUT_SUCCESS,
      currentUser: null
    });
  }).catch((err) => dispatch({ type: LOGOUT_FAIL, error: err.response.data.error }));
};

export const AUTH_SUCCESS = 'AUTH_SUCCESS';
export const AUTH_FAIL = 'AUTH_FAIL';

export const tokenAuthenticate = () => (dispatch) => {
  return request.get('api/users/authenticate').then((res) => {
    return dispatch({
      type: AUTH_SUCCESS,
      currentUser: res.data
    });
  }).catch(() => {
    console.log('There is no current user');
  });
};
